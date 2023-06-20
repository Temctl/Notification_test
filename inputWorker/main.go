package main

import (
	"sync"

	"github.com/Temctl/E-Notification/inputWorker/database"
	"github.com/Temctl/E-Notification/inputWorker/middleware"
	"github.com/Temctl/E-Notification/inputWorker/worker"
	"github.com/Temctl/E-Notification/util/elog"
)

func main() {
	middleware.PrintZ()
	elog.Info().Println("SERVER STARTED...")

	database.ConnectPostgreSQL()
	// ----------------------------------------------------------------------
	// WORKER START ---------------------------------------------------------
	// ----------------------------------------------------------------------
	var wg sync.WaitGroup
	wg.Add(3)
	// ----------------------------------------------------------------------
	// CRON JOB -------------------------------------------------------------
	// ----------------------------------------------------------------------
	go func() {
		defer wg.Done()
		worker.AttentionNotificationEveryday()
	}()
	go func() {
		defer wg.Done()
		worker.SocialSend()
	}()
	// ----------------------------------------------------------------------
	// XYP WORKER -----------------------------------------------------------
	// ----------------------------------------------------------------------
	go func() {
		defer wg.Done()
		worker.XypWorker()
	}()

	// Wait for all goroutines to complete
	wg.Wait()

}
