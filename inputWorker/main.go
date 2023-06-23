package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/Temctl/E-Notification/inputWorker/middleware"
	"github.com/Temctl/E-Notification/inputWorker/worker"
	"github.com/Temctl/E-Notification/util/connections"
	"github.com/Temctl/E-Notification/util/elog"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateCollection(client *mongo.Client, databaseName string, collectionName string) error {
	// Create a new collection
	err := client.Database(databaseName).CreateCollection(context.Background(), collectionName)
	if err != nil {
		return err
	}

	fmt.Printf("Collection '%s' created!\n", collectionName)
	return nil
}
func DeleteCollection(client *mongo.Client, dbName, collectionName string) error {
	// Access the specified database and collection
	collection := client.Database(dbName).Collection(collectionName)

	// Drop the collection
	err := collection.Drop(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("Collection deleted successfully.")
	return nil
}

func Cre() {
	// Connect to MongoDB
	client, err := connections.ConnectMongoDB()
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	}
	defer client.Disconnect(context.Background())

	err = CreateCollection(client, "notification", "attentionnotification")
	if err != nil {
		fmt.Println("Error creating collection:", err)
		return
	}

	fmt.Println("Database and collection created successfully!")
}

func main() {
	middleware.PrintZ()
	elog.Info().Println("SERVER STARTED...")

	// Cre()
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
