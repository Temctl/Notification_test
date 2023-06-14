package main

import (
	"github.com/Temctl/E-Notification/inputWorker/middleware"
	"github.com/Temctl/E-Notification/inputWorker/router"
	"github.com/Temctl/E-Notification/util"
	"github.com/Temctl/E-Notification/util/connections"
	"github.com/Temctl/E-Notification/util/elog"
	"github.com/streadway/amqp"
)

func XypWorker() {
	elog.Info().Println("XYP NOTIF WORKER STARTED...")
	// ----------------------------------------------------------------------
	// GET UTIL CONFIG ------------------------------------------------------
	// ----------------------------------------------------------------------

	host := util.RABBITMQURL

	// -----------------------------------------------------------
	// RABBITMQ connection ---------------------------------------
	// -----------------------------------------------------------

	conn, err := amqp.Dial(host)
	if err != nil {
		elog.Error().Println("Cannot connect", err)
	}
	defer conn.Close()
	elog.Info().Println("RABBITMQ: Succesful connected...")

	// -----------------------------------------------------------
	// RABBITMQ CREATE CHANNEL -----------------------------------
	// -----------------------------------------------------------

	amqpChannel, err := conn.Channel()
	if err != nil {
		elog.Error().Println("Cannot create amqp channel", err)
	}
	defer amqpChannel.Close()

	// -----------------------------------------------------------
	// RABBITMQ QUEUE DECLARE ------------------------------------
	// -----------------------------------------------------------

	queue, err := amqpChannel.QueueDeclare("XYPNOTIFtest", false, false, false, false, nil)
	if err != nil {
		elog.Error().Println("couldn't declare add queue", err)
	}

	// ----------------------------------------------------------------------
	// REDIS ----------------------------------------------------------------
	// ----------------------------------------------------------------------

	redisClient, err := connections.ConnectionRedis()
	if err != nil {
		elog.Error().Panic(err)
	}

	// ------------------------------------------------------------
	// Infinite loop to continuously pop items from the list ------
	// ------------------------------------------------------------
	for {

		// --------------------------------------------------------
		// Pop an item from the list using the BLPOP command ------
		// --------------------------------------------------------
		result, err := redisClient.BLPop(0, "queue:queue").Result()

		if err != nil {
			elog.Error().Println("Error:", err)
			continue
		}
		// --------------------------------------------------------
		// Check if an item was successfully popped ---------------
		// --------------------------------------------------------
		//
		// [queue:queue, value] len 2 bol zuv gesen ug
		// [queue:queue, value1]
		// [queue:queue, value2]
		if len(result) == 2 {
			value := result[1]
			// --------------------------------------------------------
			// RABBITMQ QUEUE DEE PUBLISH HIIH ------------------------
			// --------------------------------------------------------
			err = amqpChannel.Publish(
				"",
				queue.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(value),
				},
			)
			if err != nil {
				elog.Error().Println("Publish error", err)
			}
			elog.Info().Println("RABBITMQ: Successfully Publishing message")
		}
	}
}

func main() {
	middleware.PrintZ()
	elog.Info().Println("SERVER STARTED...")

	// ----------------------------------------------------------------------
	// WORKER START ---------------------------------------------------------
	// ----------------------------------------------------------------------

	go XypWorker()

	// ----------------------------------------------------------------------
	// REST API START -------------------------------------------------------
	// ----------------------------------------------------------------------

	router.RESTAPI()
}
