package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Temctl/E-Notification/util/connections"
	"github.com/Temctl/E-Notification/util/elog"
)

type DeviceTokens struct {
	CivilId      string   `json:"civilId"`
	DeviceTokens []string `json:"deviceTokens"`
}

func DeviceTokenConfig(w http.ResponseWriter, r *http.Request) {
	var devicetokens DeviceTokens

	err := json.NewDecoder(r.Body).Decode(&devicetokens)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	// -------------------------------------------------------
	// CONNECTION REDIS CLIENT -------------------------------
	// -------------------------------------------------------
	client, err := connections.ConnectionRedis()
	if err != nil {
		elog.Error().Panic(err)
		return
	}
	// -------------------------------------------------------
	// REDIS luu data oruulah --------------------------------
	// -------------------------------------------------------

	clientErr := client.RPush("deviceTokens:"+devicetokens.CivilId, devicetokens.DeviceTokens).Err()
	if clientErr != nil {
		elog.Error().Println("redis setlehed aldaa garlaa", clientErr)
		return
	} else {
		elog.Info().Println("Successful...")
	}

	// -------------------------------------------------------
	// Close the Redis client --------------------------------
	// -------------------------------------------------------
	closeErr := client.Close()
	if closeErr != nil {
		elog.Error().Println("Error closing Redis client:", closeErr)
		return
	} else {
		elog.Info().Println("Redis client closed successfully")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(devicetokens.CivilId))
}
