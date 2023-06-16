package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

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

	// get redis for config
	notifRedis, err := connections.ConnectionRedis()
	if err != nil {
		fmt.Println("Error connecting to redis:", err)
	} else {
		fmt.Println("NotifRedis connected succesfully")
	}

	//get the four queues that will listen
	xypNotifs, attentionNotifs, regularNotifs, groupNotifs := getQueues()

	notificationType := model.NotificationType(rune(0))

	// send the rconsumed messages
	forever := make(chan bool)
	var xypModel model.XypNotification
	go func() {
		for msg := range xypNotifs { // send xyp notifs
			err := json.Unmarshal(msg.Body, &xypModel)
			if err == nil {
				helper.SendXypNotif(xypModel, notificationType, notifRedis, client)
			} else {
				panic(err)
			}

		}
	}()

	var attentionModel model.AttentionNotification
	go func() {
		for msg := range attentionNotifs { // send attention notifs
			err := json.Unmarshal(msg.Body, &attentionModel)
			if err == nil {
				go func(request model.AttentionNotification) {
					helper.SendAttentionNotif(attentionModel.CivilId, attentionModel.Regnum, attentionModel.Content, attentionModel.Type, notifRedis, client)
				}(attentionModel)

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
				go func(request model.RegularNotification, msg amqp.Delivery) {
					helper.SendRegularNotif(request.CivilId, request.Regnum, request.Content, notificationType, notifRedis, client)

					fmt.Printf("Received Message: %s\n", msg.Body)
				}(regularModel, msg)

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
				if len(groupModel.CivilIds) == 0 {
					for _, regnum := range groupModel.Regnums {
						go func(regnum string, content string) {
							helper.SendRegularNotif("", regnum, content, notificationType, notifRedis, client)
						}(regnum, groupModel.Content)
					}
				} else {
					for _, civilId := range groupModel.CivilIds {
						go func(civilId string, content string) {
							helper.SendRegularNotif(civilId, "", content, notificationType, notifRedis, client)
						}(civilId, groupModel.Content)
					}
				}
				fmt.Printf("Received group Message: %s\n", msg.Body)

			} else {
				panic(err)
			}

		}
	}()
	fmt.Println("Waiting for messages...")
	<-forever

}
