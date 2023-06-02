package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Temctl/E-Notification/outputWorker/helper"
	"github.com/Temctl/E-Notification/util/model"
	"github.com/Temctl/E-Notification/util/elog"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func init() {
	file, _ := os.Create("./log/output.log")
	log.SetOutput(file)
	file.Close()
	log.SetFlags(log.Ldate | log.Lshortfile)

	err := godotenv.Load("./env")
	if err != nil {
		panic(err)
	}
	fmt.Println("Error initializing Firebase app:")
}

func getPushToken(regnum string) string {

	return "dIMtXp4UUkdZoj1D4M8wwD:APA91bFzD_WEW2cvd6QaXRk9cllEbr_ECrREZ2KzlbjbbWpW-7I5gNYgpgZOLGUu4HpNtc_hjyPG6YYceUbjhniqQmafV-DXV5__ezlMo07-Wq1m0trdJ5H7UWPe9SgxeFmjwN8HwmBO"
}

func main() {

	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		elog.Warning(err)
	}

	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}

	msgs, err := channel.Consume(
		os.Getenv("PUSHNOTIFCHANNELKEY"),
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err == nil {
		// print consumed messages from queue
		var notif model.PushNotificationModel
		forever := make(chan bool)
		go func() {
			for msg := range msgs {
				err := json.Unmarshal(msg.Body, &notif)
				if err == nil {
					helper.Push_notif(notif, getPushToken(notif.Regnum))
					fmt.Printf("Received Message: %s\n", msg.Body)
				} else {
					panic(err)
				}

			}
		}()

		fmt.Println("Waiting for messages...")
		<-forever
	}
}
