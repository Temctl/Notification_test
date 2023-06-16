package helper

import (
	"context"
	"fmt"
	"sync"

	"firebase.google.com/go/messaging"
	"github.com/Temctl/E-Notification/util/model"
)

func getTokensFromCivilIds(civilIds []string) []string {
	//TODO return tokens from redis
	var tokens []string
	return tokens
}

func batchTokens(tokens []string, batchSize int) [][]string {
	var batches [][]string

	for batchSize < len(tokens) {
		tokens, batches = tokens[batchSize:], append(batches, tokens[0:batchSize:batchSize])
	}

	batches = append(batches, tokens)

	return batches
}

func PushToTokens(request model.RegularNotificationModel, client *messaging.Client) {
	tokenBatches := batchTokens(request.Tokens, 500) // Split tokens into batches of 500

	var successCount int
	var notSuccessCount int
	var wg sync.WaitGroup
	wg.Add(len(tokenBatches))

	for _, tokenBatch := range tokenBatches {
		fmt.Println(tokenBatch)
		go func(tokens []string, successCount int, notSuccessCount int) {

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
				Tokens: tokens,
			}
			// Send the push notification
			response, err := client.SendMulticast(context.Background(), message)
			if err != nil {
				fmt.Println("Error sending push notification:", err)
			}

			successCount += response.SuccessCount
			notSuccessCount += response.FailureCount
			defer wg.Done()
		}(tokenBatch, successCount, notSuccessCount)
	}

	wg.Wait()
	fmt.Printf("Successful count: %d\n", successCount)
	fmt.Printf("Failed count: %d\n", notSuccessCount)
	// TODO write log
	if successCount == 0 {
		//TODO when nothing is sent
	}
}

func PushToToken(request model.PushNotificationModel, deviceToken string, client *messaging.Client) {
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
