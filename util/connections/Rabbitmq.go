package connections

import (
	"fmt"

	"github.com/Temctl/E-Notification/util"
	"github.com/streadway/amqp"
)

func GetRabbitmqChannel() (*amqp.Channel, error) {

	connection, err := amqp.Dial(util.RABBITMQURL)

	if err != nil {
		fmt.Println("Error initializing Firebase app:", err)
		return nil, err
	}
	// defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	} else {
		_, err = channel.QueueDeclare(
			util.XYPNOTIFKEY,
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			panic(err)
		}
		_, err = channel.QueueDeclare(
			util.ATTENTIONNOTIFKEY,
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			panic(err)
		}
		_, err = channel.QueueDeclare(
			util.REGULARNOTIFKEY,
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			panic(err)
		}
		_, err = channel.QueueDeclare(
			util.GROUPNOTIFKEY,
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			panic(err)
		}
	}

	return channel, err
}
