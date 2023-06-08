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
		"XYPNOTIF",
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
		"ATTENTIONNOTIF",
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
	for i := 0; i < 1200; i++ {
		tmp = append(tmp, strconv.FormatInt(int64(i), 10))
	}

	// print consumed messages from queue
	forever := make(chan bool)
	go func() {
		for msg := range xypNotifs {
			err := json.Unmarshal(msg.Body, &model.XypNotification)
			if err == nil {
				helper.PushToTokens(notif, tmp, client)
				fmt.Printf("Received Message: %s\n", msg.Body)
			} else {
				panic(err)
			}

		}
	}()

	go func() {
		for msg := range attentionNotifs {
			err := json.Unmarshal(msg.Body, &notif)
			if err == nil {
				helper.PushToTokens(notif, tmp, client)
				fmt.Printf("Received Message: %s\n", msg.Body)
			} else {
				panic(err)
			}

		}
	}()

	fmt.Println("Waiting for messages...")
	<-forever

}
