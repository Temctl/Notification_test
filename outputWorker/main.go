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
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"google.golang.org/api/option"
)

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

	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		elog.ErrorLogger.Println(err)
	}

	channel, err := connection.Channel()
	if err != nil {
		elog.ErrorLogger.Println(err)
	}

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

	var tmp []string
	tmp = append(tmp, "dIMtXp4UUkdZoj1D4M8wwD:APA91bFzD_WEW2cvd6QaXRk9cllEbr_ECrREZ2KzlbjbbWpW-7I5gNYgpgZOLGUu4HpNtc_hjyPG6YYceUbjhniqQmafV-DXV5__ezlMo07-Wq1m0trdJ5H7UWPe9SgxeFmjwN8HwmBO")
	for i := 0; i < 1200; i++ {
		tmp = append(tmp, strconv.FormatInt(int64(i), 10))
	}

	// print consumed messages from queue
	forever := make(chan bool)
	var xypModel model.XypNotification
	go func() {
		for msg := range xypNotifs {
			err := json.Unmarshal(msg.Body, &xypModel)
			if err == nil {
				var push model.PushNotificationModel
				push.Body = "xyp notif test"
				push.Title = "xyp notif test"
				helper.PushToTokens(push, tmp, client)
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
				var push1 model.PushNotificationModel
				push1.Body = "attention notif test"
				push1.Title = "attention notif test"
				helper.PushToTokens(push1, tmp, client)
				fmt.Printf("Received Message: %s\n", msg.Body)
			} else {
				panic(err)
			}

		}
	}()

	var regularModel model.RegularNotification
	go func() {
		for msg := range attentionNotifs {
			err := json.Unmarshal(msg.Body, &regularModel)
			if err == nil {
				var push1 model.PushNotificationModel
				push1.Body = "regular notif test"
				push1.Title = "regular notif test"
				helper.PushToTokens(push1, tmp, client)
				helper.SendEmail()
				fmt.Printf("Received Message: %s\n", msg.Body)
			} else {
				panic(err)
			}

		}
	}()

	fmt.Println("Waiting for messages...")
	<-forever

}
