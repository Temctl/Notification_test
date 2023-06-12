package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func SendSocial() {
	civilID := "your_civil_id" // Replace with the actual value
	content := "Your content"  // Replace with the actual value

	text := fmt.Sprintf("%s: %s", civilID, content)

	SOCIAL_DATA["message"]["text"] = text
	SOCIAL_DATA["ref"] = civilID

	logger.Debug("socialdata is --------------------------------", SOCIAL_DATA)

	jsonData, err := json.Marshal(SOCIAL_DATA)
	if err != nil {
		log.Fatal("Failed to marshal SOCIAL_DATA:", err)
	}

	request, err := http.NewRequest("POST", SOCIAL_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Failed to create HTTP request:", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+AUTH_TOKEN)

	client := http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal("Failed to send HTTP request:", err)
	}
	defer response.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		log.Fatal("Failed to decode response:", err)
	}

	if code, ok := data["code"].(int); ok && code == 1000 {
		logger.Info("send social success")

		now := time.Now()
		dtString := now.Format("02/01/2006 15:04:05")

		sql := fmt.Sprintf(`
			INSERT INTO public.notif_log
			(regnum, civilid, "content", date_pushed, sent_type, org_name, org_type, notif_type)
			VALUES ('%s', %s, '%s', (to_timestamp('%s', 'dd/mm/yyyy HH24:MI:SS')), 'FACEBOOK', '', '', '%s')
		`, regnum, civilID, text, dtString, notifType)

		// Execute the SQL query using your database library

		redisClient := redis_local() // Replace with your Redis client initialization code
		redisClient.Incr(SOCIAL_NOTIF_NUM)

		fmt.Println("Social notification sent successfully.")
		return true
	} else {
		logger.Error("error sending to social", data["message"])
		return false
	}
}
