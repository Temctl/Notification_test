package helper

import (
	"context"
	"fmt"
	"sync"

	"firebase.google.com/go/messaging"
	"github.com/Temctl/E-Notification/util/model"
)

func getPushTokenRegnum(regnum string) string {
	return "dIMtXp4UUkdZoj1D4M8wwD:APA91bFzD_WEW2cvd6QaXRk9cllEbr_ECrREZ2KzlbjbbWpW-7I5gNYgpgZOLGUu4HpNtc_hjyPG6YYceUbjhniqQmafV-DXV5__ezlMo07-Wq1m0trdJ5H7UWPe9SgxeFmjwN8HwmBO"
}
func getPushTokenCivilId(civilId string) string {
	return "dIMtXp4UUkdZoj1D4M8wwD:APA91bFzD_WEW2cvd6QaXRk9cllEbr_ECrREZ2KzlbjbbWpW-7I5gNYgpgZOLGUu4HpNtc_hjyPG6YYceUbjhniqQmafV-DXV5__ezlMo07-Wq1m0trdJ5H7UWPe9SgxeFmjwN8HwmBO"
}

func PushToAll(request model.PushNotificationModel) {
	//TODO send all with async
}

func PushToGroupCivilId(request model.PushNotificationModel, civilIds []string, client *messaging.Client) {
	//TODO
}
func PushToGroupRegnum(request model.PushNotificationModel, regnums []string, client *messaging.Client) {
	//TODO
}

func batchTokens(tokens []string, batchSize int) [][]string {
	var batches [][]string

	for batchSize < len(tokens) {
		tokens, batches = tokens[batchSize:], append(batches, tokens[0:batchSize:batchSize])
	}

	batches = append(batches, tokens)

	return batches
}

func PushToTokens(request model.PushNotificationModel, deviceTokens []string, client *messaging.Client) {

	tokenBatches := batchTokens(deviceTokens, 500) // Split tokens into batches of 500

	var wg sync.WaitGroup
	wg.Add(len(tokenBatches))

	for _, tokenBatch := range tokenBatches {
		fmt.Println(tokenBatch)
		go func() {
			defer wg.Done()

			message := &messaging.MulticastMessage{
				Notification: &messaging.Notification{
					Title: request.Title,
					Body:  request.Body,
				},
				Android: &messaging.AndroidConfig{
					Priority: "normal",
					Notification: &messaging.AndroidNotification{
						Title:      request.Title,
						Body:       request.Body,
						Visibility: messaging.VisibilityPublic,
					},
				},
				APNS: &messaging.APNSConfig{
					Payload: &messaging.APNSPayload{
						Aps: &messaging.Aps{
							ContentAvailable: true,
						},
					},
				},
				Data:   request.Data,
				Tokens: tokenBatch,
			}
			// Send the push notification
			response, err := client.SendMulticast(context.Background(), message)
			if err != nil {
				fmt.Println("Error sending push notification:", err)
				return
			}

			fmt.Printf("Successful count: %d\n", response.SuccessCount)
			fmt.Printf("Failed count: %d\n", response.FailureCount)
		}()
	}

	wg.Wait()
}

func PushToNonToken(request model.PushNotificationModel, regnum string, client *messaging.Client) {
	//TODO
	PushToToken(request, getPushTokenRegnum(regnum), client)
}

func PushToToken(request model.PushNotificationModel, deviceToken string, client *messaging.Client) {
	// Initialize the Firebase app
	// opt := option.WithCredentialsFile("../config/firebase.json")
	// config := &firebase.Config{ProjectID: "mgov-12390"}
	// app, err := firebase.NewApp(context.Background(), config, opt)
	// if err != nil {
	// 	fmt.Println("Error initializing Firebase app:", err)
	// 	return
	// }

	// // Get the FCM client
	// client, err := app.Messaging(context.Background())
	// if err != nil {
	// 	fmt.Println("Error getting FCM client:", err)
	// 	return
	// }

	// Create a new notification message
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: request.Title,
			Body:  request.Body,
		},
		Android: &messaging.AndroidConfig{
			Priority: "normal",
			Notification: &messaging.AndroidNotification{
				Title:      request.Title,
				Body:       request.Body,
				Visibility: messaging.VisibilityPublic,
			},
		},
		APNS: &messaging.APNSConfig{
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					ContentAvailable: true,
				},
			},
		},
		Data:  request.Data,
		Token: deviceToken, // Replace with your desired topic or use a device token
	}

	// Send the push notification
	response, err := client.Send(context.Background(), message)
	if err != nil {
		fmt.Println("Error sending push notification:", err)
		return
	}

	fmt.Println("Successfully sent push notification:", response)
}
