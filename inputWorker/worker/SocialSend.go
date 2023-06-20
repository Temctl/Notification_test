package worker

import (
	"fmt"
	"time"

	"github.com/Temctl/E-Notification/util"
	"github.com/Temctl/E-Notification/util/elog"
	"github.com/robfig/cron"
)

func SocialSend() {
	elog.Info().Println("SocialSend CRON JOB STARTED...")
	location, err := util.GetTZ()
	if err != nil {
		elog.Error().Panic(err)
	}
	c := cron.NewWithLocation(location)

	// Define the cron job function
	cronJob := func() {
		// Perform the task or action you want to execute on the specified schedule
		fmt.Println("SocialSend Cron job ran at", time.Now())
	}
	// ----------------------------------------------------------------------
	// Add the cron job to the cron scheduler -------------------------------
	// ----------------------------------------------------------------------
	c.AddFunc("0 32 13 * *", cronJob) // Runs the job at 10:18 AM in GMT+8
	// Start the cron scheduler
	c.Start()

	// Block the program from exiting
	// Use a channel to prevent the main goroutine from exiting
	select {}
}
