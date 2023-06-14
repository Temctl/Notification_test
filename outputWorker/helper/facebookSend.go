package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var SOCIAL_URL = "https://enterprise.chatbot.mn/api/bots/fb2120ef7cb32a80270409d9f97978fd/user/notification/sendNotification?token=c875809bbef0d18801032b21fe5140ad4128322c99b03ec6f10453c89ea2cbfb"

func SendSocial(civilId string) int {
	// Create JSON data
	data := map[string]interface{}{
		"message": map[string]interface{}{
			"text": "test",
		},
		"ref":     civilId,
		"channel": "messenger",
	}

	// Convert JSON data to byte slice
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error:", err)
		return 0
	}

	// Send POST request
	response, err := http.Post(SOCIAL_URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error:", err)
		return 0
	}
	defer response.Body.Close()

	// Check response status code
	if response.StatusCode != http.StatusOK {
		fmt.Println("Error: Unexpected status code:", response.StatusCode)
		return 0
	}

	// Read response body
	var result map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error:", err)
		return 0
	}

	// Process the response
	fmt.Println("Response:", result)
	return 1
}
