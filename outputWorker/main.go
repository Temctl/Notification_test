package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/Temctl/E-Notification/outputWorker/helper"
	"github.com/Temctl/E-Notification/util"
	"github.com/Temctl/E-Notification/util/connections"
	"github.com/Temctl/E-Notification/util/elog"
	"github.com/Temctl/E-Notification/util/model"
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

func getQueues() (<-chan amqp.Delivery, <-chan amqp.Delivery, <-chan amqp.Delivery, <-chan amqp.Delivery) {
	channel, err := connections.GetRabbitmqChannel()
	if err != nil {
		fmt.Println("Error connecting to amqp channel:", err)
	} else {
		fmt.Println("RabbitMQ started succesfully")
	}

	xypNotifs, err := channel.Consume(
		util.XYPNOTIFKEY,
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
		util.ATTENTIONNOTIFKEY,
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
		util.REGULARNOTIFKEY,
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

	groupNotifs, err := channel.Consume(
		util.GROUPNOTIFKEY,
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
	return xypNotifs, attentionNotifs, regularNotifs, groupNotifs
}

func GetFCMClient() (*messaging.Client, error) {
	// Initialize the Firebase app
	opt := option.WithCredentialsFile("config/firebase.json")
	config := &firebase.Config{ProjectID: "mgov-12390"}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		fmt.Println("Error initializing Firebase app:", err)
		return nil, err
	}

	// Get the FCM client
	return app.Messaging(context.Background())
}

func main() {
	// Get the FCM client
	client, err := GetFCMClient()
	if err != nil {
		fmt.Println("Error initializing Firebase app:", err)
	} else {
		fmt.Println("FireBase started succesfully")
	}

	notifRedis, err := connections.ConnectionRedis()
	if err != nil {
		fmt.Println("Error connecting to redis:", err)
	} else {
		fmt.Println("NotifRedis connected succesfully")
	}
	xypNotifs, attentionNotifs, regularNotifs, groupNotifs := getQueues()

	// print consumed messages from queue
	forever := make(chan bool)
	var xypModel model.XypNotification
	go func() {
		for msg := range xypNotifs {
			err := json.Unmarshal(msg.Body, &xypModel)
			if err == nil {
				var civilId string
				if xypModel.CivilId == "" {
					civilId, err = notifRedis.Get("getByReg:" + xypModel.Regnum).Result()
					if err != nil {
						panic(err)
					}
				} else {
					civilId = xypModel.CivilId
				}
				userConf, err := notifRedis.HGetAll("conf:" + civilId).Result()
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
					userDeviceTokens, err := notifRedis.LRange("deviceTokens:"+civilId, 0, -1).Result()
					if err != nil {
						panic(err)
					} else {
						helper.PushToTokens(push1, userDeviceTokens, client)
					}

				}
				if isNationalEmail, ok := userConf["isNationalEmail"]; ok && isNationalEmail == "true" {
					// helper.SendNatEmail(civilId)
				}
				if isEmail, ok := userConf["isEmail"]; ok && isEmail == "true" {
					if emailAddress, ok := userConf["emailAddress"]; ok || emailAddress != "" {
						// helper.SendPrivEmail(emailAddress)
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
					civilId, err = notifRedis.Get("getByReg:" + attentionModel.Regnum).Result()
					if err != nil {
						panic(err)
					}
				} else {
					civilId = attentionModel.CivilId
				}
				helper.SendAttentionNotif(civilId, attentionModel.Content, attentionModel.Type, notifRedis, client)

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
					civilId, err = notifRedis.Get("getByReg:" + regularModel.Regnum).Result()
					if err != nil {
						panic(err)
					}
				} else {
					civilId = regularModel.CivilId
				}
				var notificationType model.NotificationType
				helper.SendRegularNotif(civilId, regularModel.Content, notificationType, notifRedis, client)

				fmt.Printf("Received Message: %s\n", msg.Body)
			} else {
				panic(err)
			}

		}
	}()

	var groupModel model.GroupNotification
	go func() {
		for msg := range groupNotifs {
			err := json.Unmarshal(msg.Body, &groupModel)
			if err == nil {
				var notificationType model.NotificationType
				if len(groupModel.CivilIds) == 0 {
					for _, regnum := range groupModel.Regnums {
						civilId, err := notifRedis.Get("getByReg:" + regnum).Result()
						if err != nil {
							panic(err)
						}
						helper.SendRegularNotif(civilId, groupModel.Content, notificationType, notifRedis, client)
					}

				} else {
					for _, civilId := range groupModel.CivilIds {
						helper.SendRegularNotif(civilId, groupModel.Content, notificationType, notifRedis, client)
					}
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
