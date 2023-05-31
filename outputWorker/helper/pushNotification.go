package helper

import (
	"fmt"

	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/Temctl/E-Notification/util/model"
	"google.golang.org/api/option"
)

func Push_notif(request model.PushNotificationModel, token string) {
	// Initialize the Firebase app
	opt := option.WithCredentialsFile("../config/firebase.json")
	config := &firebase.Config{ProjectID: "mgov-12390"}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		fmt.Println("Error initializing Firebase app:", err)
		return
	}

	// Get the FCM client
	client, err := app.Messaging(context.Background())
	if err != nil {
		fmt.Println("Error getting FCM client:", err)
		return
	}

	// Create a new notification message
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: request.Title,
			Body:  request.Body,
		},
		Data:  request.Data,
		Token: token, // Replace with your desired topic or use a device token
	}

	// Send the push notification
	response, err := client.Send(context.Background(), message)
	if err != nil {
		fmt.Println("Error sending push notification:", err)
		return
	}

	fmt.Println("Successfully sent push notification:", response)
}
