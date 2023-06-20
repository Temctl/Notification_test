package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/Temctl/E-Notification/util"
	"github.com/Temctl/E-Notification/util/connections"
	"github.com/Temctl/E-Notification/util/elog"
	"github.com/Temctl/E-Notification/util/model"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"google.golang.org/api/option"
)

var currentWorker string

func init() {
	file, _ := os.Create("./log/output.log")
	log.SetOutput(file)
	file.Close()
	log.SetFlags(log.Ldate | log.Lshortfile)

	err := godotenv.Load(".env")
	if err != nil {
		elog.ErrorLogger.Println(err)
	}

	currentWorker = os.Getenv("workername")
}

func getQueues() (<-chan amqp.Delivery, <-chan amqp.Delivery, <-chan amqp.Delivery, <-chan amqp.Delivery) {
	channel, err := connections.GetRabbitmqChannel()
	if err != nil {
		fmt.Println("Error connecting to amqp channel:", err)
	} else {
		fmt.Println("RabbitMQ started succesfully")
	}

	pushNotif, err := channel.Consume(
		util.PUSHNOTIFICATIONKEY,
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

	natEmail, err := channel.Consume(
		util.NATEMAILKEY,
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

	privEmail, err := channel.Consume(
		util.PRIVEMAILKEY,
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

	messegeNotif, err := channel.Consume(
		util.MESSENGERKEY,
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
	return pushNotif, natEmail, privEmail, messegeNotif
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
	// client, err := GetFCMClient()
	// if err != nil {
	// 	fmt.Println("Error initializing Firebase app:", err)
	// } else {
	// 	fmt.Println("FireBase started succesfully")
	// }

	//get the four queues that will listen
	pushNotif, natEmail, privEmail, messegeNotif := getQueues()
	// send the rconsumed messages
	forever := make(chan bool)
	if currentWorker == util.PUSHWORKER {
		var pushRequest model.RegularNotificationModel
		go func() {
			for msg := range pushNotif { // send xyp notifs
				if connections.IsWorkerOn(util.PUSHWORKER) == 1 {
					err := json.Unmarshal(msg.Body, &pushRequest)
					if err == nil {
						// helper.PushToTokens(pushRequest, client)
						fmt.Println("push")
					} else {
						panic(err)
					}
				} else {
					fmt.Println("worker is turned off")
				}
			}
		}()
	} else if currentWorker == util.NATEMAILWORKER {
		var natEmailRequest model.EmailModel
		go func() {
			for msg := range natEmail {
				if connections.IsWorkerOn(util.NATEMAILWORKER) == 1 {
					err := json.Unmarshal(msg.Body, &natEmailRequest)
					if err == nil {
						go func(request model.EmailModel) {
							// helper.SendNatEmail(request.CivilId, request.Body)
							fmt.Println("natemail")
						}(natEmailRequest)

						fmt.Printf("Received Message: %s\n", msg.Body)
					} else {
						panic(err)
					}
				} else {
					fmt.Println("worker is turned off")
				}
			}
		}()
	} else if currentWorker == util.PRIVEMAILWORKER {
		var privEmailRequest model.EmailModel
		go func() {
			for msg := range privEmail {
				if connections.IsWorkerOn(util.PRIVEMAILWORKER) == 1 {
					err := json.Unmarshal(msg.Body, &privEmailRequest)
					if err == nil {
						go func(request model.EmailModel) {
							// helper.SendPrivEmail(request.CivilId, request.Body)
							fmt.Println("privemail")
						}(privEmailRequest)
						fmt.Printf("Received Message: %s\n", msg.Body)
					} else {
						panic(err)
					}
				} else {
					fmt.Println("worker is turned off")
				}
			}
		}()
	} else if currentWorker == util.MESSENGERWORKER {
		var messengerRequest model.MessengerModel
		go func() {
			for msg := range messegeNotif {
				if connections.IsWorkerOn(util.MESSENGERWORKER) == 1 {
					err := json.Unmarshal(msg.Body, &messengerRequest)
					if err == nil {
						go func(request model.MessengerModel) {
							// helper.SendMessenger(request.CivilId, request.Body)
							fmt.Println("messenger")
						}(messengerRequest)
						fmt.Printf("Received Message: %s\n", msg.Body)
					} else {
						panic(err)
					}
				} else {
					fmt.Println("worker is turned off")
				}
			}
		}()
	} else {
		fmt.Println("re wokrer")
	}

	fmt.Println("Waiting for messages...")
	<-forever

}
