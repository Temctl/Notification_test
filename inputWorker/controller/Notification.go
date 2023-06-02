package controller

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Temctl/E-Notification/inputWorker/utils"
	"github.com/streadway/amqp"
)

func Notification(w http.ResponseWriter, r *http.Request) {
	// GET ENV
	rabbitHost := os.Getenv("RABBITMQ_HOST")
	rabbitPort := os.Getenv("RABBITMQ_PORT")

	// RABBITMQ connection example
	conn, err := amqp.Dial("amqp://guest:guest@" + rabbitHost + ":" + rabbitPort + "/")
	utils.HandleError(err, "Cannot connect")
	defer conn.Close()

	fmt.Fprintln(w, "Succesful connected...")

	// CREATE CHANNEL
	amqpChannel, err := conn.Channel()
	utils.HandleError(err, "Cannot create amqp channel")
	defer amqpChannel.Close()

	// QUEUE DECLARE
	queue, err := amqpChannel.QueueDeclare("queue2", false, false, false, false, nil)
	utils.HandleError(err, "couldn't declare add queue")

	// err = amqpChannel.Qos(1, 0, false)
	// handleError(err, "could notconfig Qos")
	fmt.Fprintln(w, queue)
	err = amqpChannel.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("push2"),
		},
	)
	utils.HandleError(err, "Publish error")
	fmt.Fprintln(w, "Successfully Publishing message")
}
