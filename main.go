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
	"strconv"
	"strings"
	"encoding/json"
	"math/rand"
	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v8/linebot"
	_ "github.com/line/line-bot-sdk-go/v8/linebot/channel_access_token"
	_ "github.com/line/line-bot-sdk-go/v8/linebot/insight"
	_ "github.com/line/line-bot-sdk-go/v8/linebot/liff"
	_ "github.com/line/line-bot-sdk-go/v8/linebot/manage_audience"
	_ "github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	_ "github.com/line/line-bot-sdk-go/v8/linebot/module"
	_ "github.com/line/line-bot-sdk-go/v8/linebot/module_attach"
	_ "github.com/line/line-bot-sdk-go/v8/linebot/shop"
	_ "github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"time"
	"net/url"
)

var bot *linebot.Client
var registeredIDs map[string]bool = make(map[string]bool)

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
func luhnAlgorithm(cardNumber string) bool {
    // this function implements the luhn algorithm
    // it takes as argument a cardnumber of type string
    // and it returns a boolean (true or false) if the
    // card number is valid or not

    // initialise a variable to keep track of the total sum of digits
    total := 0
    // Initialize a flag to track whether the current digit is the second digit from the right.
    isSecondDigit := false

    // iterate through the card number digits in reverse order
    for i := len(cardNumber) - 1; i >= 0; i-- {
        // conver the digit character to an integer
        digit := int(cardNumber[i] - '0')

        if isSecondDigit {
            // double the digit for each second digit from the right
            digit *= 2
            if digit > 9 {
                // If doubling the digit results in a two-digit number,
                //subtract 9 to get the sum of digits.
                digit -= 9
            }
        }

        // Add the current digit to the total sum
        total += digit

        //Toggle the flag for the next iteration.
        isSecondDigit = !isSecondDigit
    }

    // return whether the total sum is divisible by 10
    // making it a valid luhn number
    return total%10 == 0
}
func check_credit_card(input string) string {
	result := luhnAlgorithm(input)
	var output string

	if result {
		
		output = "信用卡號正確"
		fmt.Printf("%s",output)
	} else{
		output = "信用卡號錯誤 請重新輸入"
		fmt.Printf("%s",output)
	}
	return output
}
type Output struct {
    Title string
    Id  string
	ChannelTitle string
	LikeCount string
	ViewCount string
	PublishedAt string
	CommentCount string
}


func check_yt_imformation(youtubeURL string) string {
	// TODO: Get API token from .env file
	// TODO: Get video ID from URL query `v`
	// TODO: Get video information from YouTube API
	// TODO: Parse the JSON response and store the information into a struct
	// TODO: Display the information in an HTML page through `template`
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	youtubeAPIKey := os.Getenv("YOUTUBE_API_KEY")	
	baseURL := "https://www.googleapis.com/youtube/v3/videos"
	parsedURL, _ := url.Parse(youtubeURL)
	videoID := parsedURL.Query().Get("v")
	if videoID == "" {
		// http.ServeFile(w, r, "error.html")
		return "error"
	}	
	url := fmt.Sprintf("%s?part=statistics,snippet&id=%s&key=%s", baseURL, videoID, youtubeAPIKey)
	resp, err := http.Get(url)
	if err != nil {
		// http.ServeFile(w, r, "error.html")
		return "error" 
	}
	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		// http.ServeFile(w, r, "error.html")
		return "error"
	}

	items, ok := data["items"].([]interface{})
	if !ok || len(items) == 0 {
		// http.ServeFile(w, r, "error.html")
		return "error"
	}
	statistics, ok := items[0].(map[string]interface{})["statistics"].(map[string]interface{})
	if !ok {
		// http.ServeFile(w, r, "error.html")
		return "error"
	}
	views := formatNumber(statistics["viewCount"].(string))
	likes := formatNumber(statistics["likeCount"].(string))
	comments := formatNumber(statistics["commentCount"].(string))
	title := items[0].(map[string]interface{})["snippet"].(map[string]interface{})["title"].(string)
	channelTitle := items[0].(map[string]interface{})["snippet"].(map[string]interface{})["channelTitle"].(string)
	publishedAt := items[0].(map[string]interface{})["snippet"].(map[string]interface{})["publishedAt"].(string)
	parsedTime, err := time.Parse(time.RFC3339, publishedAt)
	if err != nil {
		// http.ServeFile(w, r, "error.html")
		return "error"
	}
	formattedDate := parsedTime.Format("2006年01月02日")
	var output = Output{
		Title :title,
		Id :videoID,
		ChannelTitle :channelTitle,
		LikeCount :likes,
		ViewCount :views,
		PublishedAt :formattedDate,
		CommentCount :comments,
	}
	outputString := fmt.Sprintf("Information:\n"+
		"Title: %s\n"+
		"Channel Title: %s\n"+
		"Like Count: %s\n"+
		"View Count: %s\n"+
		"Published At: %s\n"+
		"Comment Count: %s\n"+
		"------------------------------",
		output.Title, output.ChannelTitle, output.LikeCount, output.ViewCount, output.PublishedAt, output.CommentCount)

	// 印出結果
	fmt.Println(outputString)
	return outputString
}
func formatNumber(number string) string {
	// Format the number with commas every 3 digits
	parts := strings.Split(number, "")
	result := ""
	for i := len(parts) - 1; i >= 0; i-- {
		result = parts[i] + result
		if (len(parts)-i)%3 == 0 && i != 0 {
			result = "," + result
		}
	}
	return result
}
func drawLottery(members []string) (string, error) {
	// 檢查是否有成員可供抽籤
	if len(members) == 0 {
		return "", fmt.Errorf("沒有可供抽籤的成員")
	}

	winnerIndex := rand.Intn(len(members))
	winner := members[winnerIndex]

	return winner, nil
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
				if message.Text == "!功能表" {
					result = "1. 計算\n2. 驗證信用卡\n3. 抽籤\n4. 查詢影片資訊\n5. 登記\n6. 取消登記"
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(result)).Do(); err != nil {
						log.Print(err)
					}
				} else if message.Text == "!登記" {

                    registeredIDs[event.Source.UserID] = true

					userProfile, err := bot.GetProfile(event.Source.UserID).Do()
					if err != nil {
						log.Print(err)
						return
					}				
                	userName := userProfile.DisplayName
					result := fmt.Sprintf("%s 成功登記", userName)
                    if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(result)).Do(); err != nil {
                        log.Print(err)
                    }
                } else if message.Text == "!取消登記" {

                    registeredIDs[event.Source.UserID] = false

					userProfile, err := bot.GetProfile(event.Source.UserID).Do()
					if err != nil {
						log.Print(err)
						return
					}				
                	userName := userProfile.DisplayName
					result := fmt.Sprintf("%s 取消登記", userName)
                    if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(result)).Do(); err != nil {
                        log.Print(err)
                    }
                }
				feature := strings.Split(message.Text, " ")
				fmt.Printf("切割後的結果: %v\n", feature[0])
				if feature[0] == "!計算" {
					
					result = calculate(feature[1])			
					
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(result)).Do(); err != nil {
						log.Print(err)
					}				
				} else if feature[0] == "!驗證信用卡"{
					result = check_credit_card(feature[1])
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(result)).Do(); err != nil {
						log.Print(err)
					}	

				}else if feature[0] == "!查詢影片資訊"{
					result = check_yt_imformation(feature[1])
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(result)).Do(); err != nil {
						log.Print(err)
					}	

				}else if feature[0] == "!抽籤"{
					members := make([]string, 0, len(registeredIDs))
                    for id := range registeredIDs {
                        members = append(members, id)
                    }

					// 執行群組抽籤
					winner, err := drawLottery(members)
					if err != nil {
						fmt.Println("抽籤時發生錯誤:", err)
						return
					}
					userProfile, err := bot.GetProfile(winner).Do()
					if err != nil {
						log.Print(err)
						return
					}				
                	winnerName := userProfile.DisplayName
					result = fmt.Sprintf("抽中的成員是：%s\n", winnerName)
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
