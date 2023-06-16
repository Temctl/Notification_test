package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Temctl/E-Notification/inputWorker/model"
	"github.com/Temctl/E-Notification/util/connections"
	"github.com/Temctl/E-Notification/util/elog"
)

type DeviceTokens struct {
	CivilId      string   `json:"civilId"`
	DeviceTokens []string `json:"deviceTokens"`
}

func DeviceTokenConfig(w http.ResponseWriter, r *http.Request) {
	// -------------------------------------------------------
	// RESPONSE ----------------------------------------------
	// -------------------------------------------------------
	response := model.ApiResponse{
		Status:  http.StatusOK,
		Message: "Success",
	}
	// -------------------------------------------------------
	// PROGRESS ----------------------------------------------
	// -------------------------------------------------------
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
		response.Message = err.Error()
		response.Status = 400
	}
	// -------------------------------------------------------
	// REDIS luu data oruulah --------------------------------
	// -------------------------------------------------------
	jsonData, err := json.Marshal(devicetokens.DeviceTokens)
	if err != nil {
		elog.Error().Panic(err)
		response.Message = err.Error()
		response.Status = 400
	}
	clientErr := client.Set("deviceTokens:"+devicetokens.CivilId, jsonData, 0).Err()
	if clientErr != nil {
		elog.Error().Println("redis setlehed aldaa garlaa", clientErr)
		response.Message = "redis setlehed aldaa garlaa"
		response.Status = 400
	}
	elog.Info().Println("Successful...")

	// -------------------------------------------------------
	// Close the Redis client --------------------------------
	// -------------------------------------------------------
	closeErr := client.Close()
	if closeErr != nil {
		elog.Error().Println("Error closing Redis client:", closeErr)
		response.Message = "Error closing Redis client:" + closeErr.Error()
		response.Status = 400
	}
	elog.Info().Println("Redis client closed successfully")
	// RESPONSE SETUP
	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Set the response header content type
	w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusBadRequest)
	w.Write(responseJson)
}
