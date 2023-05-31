package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Temctl/E-Notification/outputWorker/helper"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

func init() {
	file, _ := os.Create("./log/output.log")
	log.SetOutput(file)
	file.Close()
	log.SetFlags(log.Ldate | log.Lshortfile)

	err := godotenv.Load(".config/env")
	if err != nil {
		panic(err)
	}
	fmt.Println("Error initializing Firebase app:")
}

func main() {

	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
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

	// print consumed messages from queue
	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			helper.Push_notif(msg.Body, "dIMtXp4UUkdZoj1D4M8wwD:APA91bFzD_WEW2cvd6QaXRk9cllEbr_ECrREZ2KzlbjbbWpW-7I5gNYgpgZOLGUu4HpNtc_hjyPG6YYceUbjhniqQmafV-DXV5__ezlMo07-Wq1m0trdJ5H7UWPe9SgxeFmjwN8HwmBO")
			fmt.Printf("Received Message: %s\n", msg.Body)
		}
	}()

	fmt.Println("Waiting for messages...")
	<-forever
}
