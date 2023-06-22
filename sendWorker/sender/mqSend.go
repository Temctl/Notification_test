package sender

import (
	"encoding/json"
	"fmt"

	"github.com/Temctl/E-Notification/util"
	"github.com/Temctl/E-Notification/util/connections"
	"github.com/Temctl/E-Notification/util/model"
	"github.com/streadway/amqp"
)

func MqPush(request model.RegularNotificationModel) {
	channel, err := connections.GetRabbitmqChannel()
	if err != nil {
		fmt.Println("Error connecting to amqp channel:", err)
	}
	data, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	err = channel.Publish(
		"",
		util.PUSHNOTIFICATIONKEY,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
	if err != nil {
		panic(err)
	}
}

func MqNatEmail(request model.EmailModel) {
	channel, err := connections.GetRabbitmqChannel()
	if err != nil {
		fmt.Println("Error connecting to amqp channel:", err)
	}
	data, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	err = channel.Publish(
		"",
		util.NATEMAILKEY,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
	if err != nil {
		panic(err)
	}
}

func MqPrivEmail(request model.EmailModel) {
	channel, err := connections.GetRabbitmqChannel()
	if err != nil {
		fmt.Println("Error connecting to amqp channel:", err)
	}
	data, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	err = channel.Publish(
		"",
		util.PRIVEMAILKEY,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
	if err != nil {
		panic(err)
	}
}

func MqMessenger(request model.MessengerModel) {
	channel, err := connections.GetRabbitmqChannel()
	if err != nil {
		fmt.Println("Error connecting to amqp channel:", err)
	}
	data, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	err = channel.Publish(
		"",
		util.MESSENGERKEY,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
	if err != nil {
		panic(err)
	}
}
