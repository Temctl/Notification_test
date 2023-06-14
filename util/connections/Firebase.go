package connections

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

func GetFCMClient() (*messaging.Client, error) {
	// Initialize the Firebase app
	opt := option.WithCredentialsFile("config/firebase.json")
	config := &firebase.Config{ProjectID: "mgov-12390"}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		fmt.Println("Error initializing Firebase app:", err)
	}

	// return the FCM client
	return app.Messaging(context.Background())
}
