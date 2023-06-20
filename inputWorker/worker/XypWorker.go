package worker

import (
	"github.com/Temctl/E-Notification/util"
	"github.com/Temctl/E-Notification/util/connections"
	"github.com/Temctl/E-Notification/util/elog"
	"github.com/streadway/amqp"
)

func XypWorker() {
	elog.Info().Println("XYP NOTIF WORKER STARTED...")
	// ----------------------------------------------------------------------
	// RABBITMQ CONNECTION --------------------------------------------------
	// ----------------------------------------------------------------------

	queue, rErr := connections.GetRabbitmqChannel()
	if rErr != nil {
		elog.Error().Panic(rErr)
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
		result, err := redisClient.BLPop(0, "XYPNOTIF").Result()

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
			err = queue.Publish(
				"",
				util.PUSHNOTIFICATIONKEY,
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
