package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	firebase "firebase.google.com/go"
	"github.com/Temctl/E-Notification/outputWorker/helper"
	"github.com/Temctl/E-Notification/util/elog"
	"github.com/Temctl/E-Notification/util/model"
	"github.com/Temctl/E-Notification/util/redis"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"google.golang.org/api/option"
)

var SOCIAL_URL = "https://enterprise.chatbot.mn/api/bots/fb2120ef7cb32a80270409d9f97978fd/user/notification/sendNotification?token=c875809bbef0d18801032b21fe5140ad4128322c99b03ec6f10453c89ea2cbfb"

func init() {
	file, _ := os.Create("./log/output.log")
	log.SetOutput(file)
	file.Close()
	log.SetFlags(log.Ldate | log.Lshortfile)

	err := godotenv.Load(".env")
	if err != nil {
		elog.ErrorLogger.Println(err)
	}
}

func main() {

	// Initialize the Firebase app
	opt := option.WithCredentialsFile("config/firebase.json")
	config := &firebase.Config{ProjectID: "mgov-12390"}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		fmt.Println("Error initializing Firebase app:", err)
		return
	}

	// Get the FCM client
	client, err := app.Messaging(context.Background())
	if err != nil {
		fmt.Println("Error getting FCM client:", err)
		return
	}

	connection, err := amqp.Dial(util.)
	if err != nil {
		fmt.Println("Error connection to rabbitmq:", err)
	}

	channel, err := connection.Channel()
	if err != nil {
		fmt.Println("Error connecting to amqp channel:", err)
	}

	notifRedis := redis.ConnectionRedis()

	xypNotifs, err := channel.Consume(
		"XYPNOTIFtest",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	attentionNotifs, err := channel.Consume(
		"ATTENTIONNOTIFtest",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	regularNotifs, err := channel.Consume(
		"REGULARNOTIFtest",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	// print consumed messages from queue
	forever := make(chan bool)
	var xypModel model.XypNotification
	go func() {
		for msg := range xypNotifs {
			err := json.Unmarshal(msg.Body, &xypModel)
			if err == nil {
				var civilId string
				if xypModel.CivilId == "" {
					civilId, err = notifRedis.Get(context.Background(), "getByReg:"+xypModel.Regnum).Result()
					if err != nil {
						panic(err)
					}
				} else {
					civilId = xypModel.CivilId
				}
				userConf, err := notifRedis.HGetAll(context.Background(), "conf:"+civilId).Result()
				var push1 model.PushNotificationModel
				push1.Body = "regular notif test"
				push1.Title = "regular notif test"

				if err != nil {
					panic(err)
				}
				if isPush, ok := userConf["isPush"]; ok && isPush == "true" {
					var tmp []string
					tmp = append(tmp, "dIMtXp4UUkdZoj1D4M8wwD:APA91bFzD_WEW2cvd6QaXRk9cllEbr_ECrREZ2KzlbjbbWpW-7I5gNYgpgZOLGUu4HpNtc_hjyPG6YYceUbjhniqQmafV-DXV5__ezlMo07-Wq1m0trdJ5H7UWPe9SgxeFmjwN8HwmBO")
					for i := 0; i < 700; i++ {
						tmp = append(tmp, strconv.FormatInt(int64(i), 10))
					}
					userDeviceTokens, err := notifRedis.LRange(context.Background(), "deviceTokens:"+civilId, 0, -1).Result()
					if err != nil {
						panic(err)
					} else {
						helper.PushToTokens(push1, userDeviceTokens, client)
					}

				}
				if isNationalEmail, ok := userConf["isNationalEmail"]; ok && isNationalEmail == "true" {
					helper.SendNatEmail(civilId)
				}
				if isEmail, ok := userConf["isEmail"]; ok && isEmail == "true" {
					if emailAddress, ok := userConf["emailAddress"]; ok || emailAddress != "" {
						helper.SendPrivEmail(emailAddress)
					}
				}
				if isSocial, ok := userConf["social"]; ok && isSocial == "true" {
					helper.SendSocial(civilId)
				}

				fmt.Printf("Received Message: %s\n", msg.Body)
			} else {
				panic(err)
			}

		}
	}()

	var attentionModel model.AttentionNotification
	go func() {
		for msg := range attentionNotifs {
			err := json.Unmarshal(msg.Body, &attentionModel)
			if err == nil {
				var civilId string
				if attentionModel.CivilId == "" {
					civilId, err = notifRedis.Get(context.Background(), "getByReg:"+attentionModel.Regnum).Result()
					if err != nil {
						panic(err)
					}
				} else {
					civilId = attentionModel.CivilId
				}
				userConf, err := notifRedis.HGetAll(context.Background(), "conf:"+civilId).Result()
				var push1 model.PushNotificationModel
				push1.Body = "regular notif test"
				push1.Title = "regular notif test"

				if err != nil {
					panic(err)
				}
				if isPush, ok := userConf["isPush"]; ok && isPush == "true" {
					var tmp []string
					tmp = append(tmp, "dIMtXp4UUkdZoj1D4M8wwD:APA91bFzD_WEW2cvd6QaXRk9cllEbr_ECrREZ2KzlbjbbWpW-7I5gNYgpgZOLGUu4HpNtc_hjyPG6YYceUbjhniqQmafV-DXV5__ezlMo07-Wq1m0trdJ5H7UWPe9SgxeFmjwN8HwmBO")
					for i := 0; i < 700; i++ {
						tmp = append(tmp, strconv.FormatInt(int64(i), 10))
					}
					userDeviceTokens, err := notifRedis.LRange(context.Background(), "deviceTokens:"+civilId, 0, -1).Result()
					if err != nil {
						panic(err)
					} else {
						helper.PushToTokens(push1, userDeviceTokens, client)
					}

				}
				if isNationalEmail, ok := userConf["isNationalEmail"]; ok && isNationalEmail == "true" {
					helper.SendNatEmail(civilId)
				}
				if isEmail, ok := userConf["isEmail"]; ok && isEmail == "true" {
					if emailAddress, ok := userConf["emailAddress"]; ok || emailAddress != "" {
						helper.SendPrivEmail(emailAddress)
					}
				}
				if isSocial, ok := userConf["social"]; ok && isSocial == "true" {
					helper.SendSocial(civilId)
				}

				fmt.Printf("Received Message: %s\n", msg.Body)
			} else {
				panic(err)
			}

		}
	}()

	var regularModel model.RegularNotification
	go func() {
		for msg := range regularNotifs {
			err := json.Unmarshal(msg.Body, &regularModel)
			if err == nil {
				var civilId string
				if regularModel.CivilId == "" {
					civilId, err = notifRedis.Get(context.Background(), "getByReg:"+regularModel.Regnum).Result()
					if err != nil {
						panic(err)
					}
				} else {
					civilId = regularModel.CivilId
				}
				userConf, err := notifRedis.HGetAll(context.Background(), "conf:"+civilId).Result()
				var push1 model.PushNotificationModel
				push1.Body = "regular notif test"
				push1.Title = "regular notif test"

				if err != nil {
					panic(err)
				}
				if isPush, ok := userConf["isPush"]; ok && isPush == "true" {
					var tmp []string
					tmp = append(tmp, "dIMtXp4UUkdZoj1D4M8wwD:APA91bFzD_WEW2cvd6QaXRk9cllEbr_ECrREZ2KzlbjbbWpW-7I5gNYgpgZOLGUu4HpNtc_hjyPG6YYceUbjhniqQmafV-DXV5__ezlMo07-Wq1m0trdJ5H7UWPe9SgxeFmjwN8HwmBO")
					for i := 0; i < 700; i++ {
						tmp = append(tmp, strconv.FormatInt(int64(i), 10))
					}
					userDeviceTokens, err := notifRedis.LRange(context.Background(), "deviceTokens:"+civilId, 0, -1).Result()
					if err != nil {
						panic(err)
					} else {
						helper.PushToTokens(push1, userDeviceTokens, client)
					}

				}
				if isNationalEmail, ok := userConf["isNationalEmail"]; ok && isNationalEmail == "true" {
					helper.SendNatEmail(civilId)
				}
				if isEmail, ok := userConf["isEmail"]; ok && isEmail == "true" {
					if emailAddress, ok := userConf["emailAddress"]; ok || emailAddress != "" {
						helper.SendPrivEmail(emailAddress)
					}
				}
				if isSocial, ok := userConf["social"]; ok && isSocial == "true" {
					helper.SendSocial(civilId)
				}

				fmt.Printf("Received Message: %s\n", msg.Body)
			} else {
				panic(err)
			}

		}
	}()

	fmt.Println("Waiting for messages...")
	<-forever

}
