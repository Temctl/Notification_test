package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Temctl/E-Notification/outputWorker/helper"
	"github.com/Temctl/E-Notification/util/elog"
	"github.com/Temctl/E-Notification/util/model"
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
		elog.ErrorLogger.Println(err)
	}
}

func main() {

	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		elog.Warning(err)
	}

	channel, err := connection.Channel()
	if err != nil {
		elog.Warning(err)
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
