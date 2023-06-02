package controller

import (
	"net/http"
	"strconv"

	"github.com/Temctl/E-Notification/util"
	"github.com/Temctl/E-Notification/util/elog"
	"github.com/streadway/amqp"
)

func Input(w http.ResponseWriter, r *http.Request) {
	// -------------------------------------------------------
	// GET UTIL CONFIG ---------------------------------------
	// -------------------------------------------------------

	host := util.RABBITMQ_HOST
	port := strconv.Itoa(util.RABBITMQ_PORT)

	// -------------------------------------------------------
	// RABBITMQ connection -----------------------------------
	// -------------------------------------------------------

	conn, err := amqp.Dial("amqp://guest:guest@" + host + ":" + port + "/")
	elog.Error("Cannot connect", err)

	defer conn.Close()

	elog.Info("RABBITMQ: Succesful connected...")

	// -------------------------------------------------------
	// CREATE CHANNEL ----------------------------------------
	// -------------------------------------------------------

	amqpChannel, err := conn.Channel()

	elog.Error("Cannot create amqp channel", err)

	defer amqpChannel.Close()

	// -------------------------------------------------------
	// QUEUE DECLARE -----------------------------------------
	// -------------------------------------------------------

	queue, err := amqpChannel.QueueDeclare("queue2", false, false, false, false, nil)
	elog.Error("couldn't declare add queue", err)
	// err = amqpChannel.Qos(1, 0, false)
	// handleError(err, "could notconfig Qos")

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
	elog.Error("Publish error", err)
	elog.Info("RABBITMQ: Successfully Publishing message")
}
