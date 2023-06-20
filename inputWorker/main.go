package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/Temctl/E-Notification/inputWorker/middleware"
	"github.com/Temctl/E-Notification/util"
	"github.com/Temctl/E-Notification/util/connections"
	"github.com/Temctl/E-Notification/util/elog"
	"github.com/robfig/cron"
	"github.com/streadway/amqp"
)

func XypWorker() {
	fmt.Println("Fsdfds")
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

func AsyncCronJob() {
	fmt.Println("Fsdf")
	location, err := time.LoadLocation("Asia/Ulaanbaatar")
	if err != nil {
		fmt.Println("Error loading time zone:", err)
		return
	}
	c := cron.NewWithLocation(location)

	// Define the cron job function
	cronJob := func() {
		// Perform the task or action you want to execute on the specified schedule
		fmt.Println("Cron job ran at", time.Now())
	}
	// ----------------------------------------------------------------------
	// Add the cron job to the cron scheduler -------------------------------
	// ----------------------------------------------------------------------
	c.AddFunc("0 7 13 * *", cronJob) // Runs the job at 10:18 AM in GMT+8
	// Start the cron scheduler
	c.Start()

	// Block the program from exiting
	// Use a channel to prevent the main goroutine from exiting
	select {}
}

func main() {
	middleware.PrintZ()
	elog.Info().Println("SERVER STARTED...")

	// ----------------------------------------------------------------------
	// WORKER START ---------------------------------------------------------
	// ----------------------------------------------------------------------
	var wg sync.WaitGroup
	wg.Add(2)

	// ----------------------------------------------------------------------
	// CRON JOB -------------------------------------------------------------
	// ----------------------------------------------------------------------
	go func() {
		defer wg.Done()
		AsyncCronJob()
	}()

	// ----------------------------------------------------------------------
	// XYP WORKER -----------------------------------------------------------
	// ----------------------------------------------------------------------
	go func() {
		defer wg.Done()
		XypWorker()
	}()

	// Wait for all goroutines to complete
	wg.Wait()

}
