package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Temctl/E-Notification/util/elog"
	"github.com/joho/godotenv"
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

func generate1(wg *sync.WaitGroup) {
	forever := make(chan bool)
	go func() {
		fmt.Println("3")
	}()

	fmt.Println("Waiting for messages...")
	<-forever
}

func generate2(wg *sync.WaitGroup) {
	forever := make(chan bool)
	go func() {
		fmt.Println("5")
	}()

	fmt.Println("Waiting for messages...")
	<-forever
}
func main() {
	var wg sync.WaitGroup

	wg.Add(2)
	go generate1(&wg)
	go generate2(&wg)

	fmt.Println("Waiting for goroutines to finish...")
	wg.Wait()
	fmt.Println("Done!")
	// // Initialize the Firebase app
	// opt := option.WithCredentialsFile("config/firebase.json")
	// config := &firebase.Config{ProjectID: "mgov-12390"}
	// app, err := firebase.NewApp(context.Background(), config, opt)
	// if err != nil {
	// 	fmt.Println("Error initializing Firebase app:", err)
	// 	return
	// }

	// // Get the FCM client
	// client, err := app.Messaging(context.Background())
	// if err != nil {
	// 	fmt.Println("Error getting FCM client:", err)
	// 	return
	// }

	// connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	// if err != nil {
	// 	elog.ErrorLogger.Println(err)
	// }

	// channel, err := connection.Channel()
	// if err != nil {
	// 	elog.ErrorLogger.Println(err)
	// }
	// fmt.Print(os.Getenv("PUSHNOTIFCHANNELKEY"))

	// msgs, err := channel.Consume(
	// 	os.Getenv("PUSHNOTIFCHANNELKEY"),
	// 	"",
	// 	true,
	// 	false,
	// 	false,
	// 	false,
	// 	nil,
	// )
	// if err != nil {
	// 	panic(err)
	// }
	// var tmp []string
	// for i := 0; i < 1200; i++ {
	// 	tmp = append(tmp, strconv.FormatInt(int64(i), 10))
	// }

	// // print consumed messages from queue
	// var notif model.PushNotificationModel
	// forever := make(chan bool)
	// go func() {
	// 	for msg := range msgs {
	// 		err := json.Unmarshal(msg.Body, &notif)
	// 		if err == nil {
	// 			helper.PushToTokens(notif, tmp, client)
	// 			fmt.Printf("Received Message: %s\n", msg.Body)
	// 		} else {
	// 			panic(err)
	// 		}

	// 	}
	// }()

	// fmt.Println("Waiting for messages...")
	// <-forever

}
