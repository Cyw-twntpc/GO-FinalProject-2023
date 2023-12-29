// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"strconv"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var bot *linebot.Client

func main() {
	err := godotenv.Load() // Load environment variable from .env file
	if err != nil{
		log.Print(err)
	}

	ChannelAccessToken := os.Getenv("CHANNEL_ACCESS_TOKEN")
	ChannelSecret := os.Getenv("CHANNEL_SECRET")

	bot, err = linebot.New(ChannelSecret, ChannelAccessToken)
	if err != nil{
		log.Print("Bot Successfully Bulid!")
	}
	http.HandleFunc("/callback", callbackHandler)
	http.ListenAndServe(":64503", nil)
}
func calculate(input string) string{
	operators := []string{"+", "-", "*", "/"}

	// 查找輸入字串中的運算符
	var operator string
	for _, op := range operators {
		if strings.Contains(input, op) {
			operator = op
			break
		}
	}

	// 如果未找到運算符，輸出錯誤信息
	if operator == "" {
		output := "無效的輸入"
		return output
	}

	// 使用 strings.Split 分割輸入字串，獲得數字和運算符
	parts := strings.Split(input, operator)
	if len(parts) != 2 {
		output := "無效的輸入"
		return output
	}

	// 解析數字
	num1, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		output := fmt.Sprintf("無效的數字: %s", parts[0])
		return output
	}

	num2, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		output := fmt.Sprintf("無效的數字: %s", parts[1])
		return output
	}

	// 根據運算符執行相應的操作
	var result float64
	switch operator {
	case "+":
		result = num1 + num2
	case "-":
		result = num1 - num2
	case "*":
		result = num1 * num2
	case "/":
		if num2 == 0 {
			output := "除數不能為零"
			return output
		}
		result = num1 / num2
	default:
		output := fmt.Sprintf("不支持的運算符: %s", operator)
		return output
	}

	// 輸出結果
	output := fmt.Sprintf("結果: %v %s %v = %v", num1, operator, num2, result)
	return output
}
func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			// Handle only on text message
			case *linebot.TextMessage:
				// GetMessageQuota: Get how many remain free tier push message quota you still have this month. (maximum 500)
				// quota, err := bot.GetMessageQuota().Do()
				if err != nil {
					log.Println("Quota err:", err)
				}
				// message.ID: Msg unique ID
				// message.Text: Msg text
				// if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
				// 	log.Print(err)
				// }
				var result string
				if message.Text == "功能表" {
					result = "1. 計算\n2. 驗證信用卡\n3. 查詢推文\n4. 查詢影片資訊"
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(result)).Do(); err != nil {
						log.Print(err)
					}
				}
				feature := strings.Split(message.Text, " ")
				fmt.Printf("切割後的結果: %v\n", feature[0])
				if feature[0] == "計算" {
					
					result = calculate(feature[1])			
					
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(result)).Do(); err != nil {
						log.Print(err)
					}				
				}
				// length := len(message.Text)
				// fmt.Printf("%s %d",message.Text,length)
			// Handle only on Sticker message
			case *linebot.StickerMessage:
				var kw string
				for _, k := range message.Keywords {
					kw = kw + "," + k
				}

				outStickerResult := fmt.Sprintf("收到貼圖訊息: %s, pkg: %s kw: %s  text: %s", message.StickerID, message.PackageID, kw, message.Text)
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(outStickerResult)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}
