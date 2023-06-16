package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Temctl/E-Notification/inputWorker/model"
	"github.com/Temctl/E-Notification/util"
	"github.com/Temctl/E-Notification/util/connections"
	"github.com/Temctl/E-Notification/util/elog"
	umodel "github.com/Temctl/E-Notification/util/model"
	"github.com/streadway/amqp"
)

func SinglePushNotif(w http.ResponseWriter, r *http.Request) {
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
	var configData umodel.RegularNotification
	err := json.NewDecoder(r.Body).Decode(&configData)
	if err != nil {
		elog.Error().Println(err)
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	fmt.Println(configData)
	// Close the request body to prevent resource leaks
	defer r.Body.Close()

	configDataJson, err := json.Marshal(configData)
	if err != nil {
		elog.Error().Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// ----------------------------------------------------------------------
	// RABBITMQ CONNECTION --------------------------------------------------
	// ----------------------------------------------------------------------
	queue, rErr := connections.GetRabbitmqChannel()
	if rErr != nil {
		elog.Error().Println(rErr)
	}
	err = queue.Publish(
		"",
		util.NATEMAILKEY,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(configDataJson),
		},
	)
	if err != nil {
		elog.Error().Println("Publish error", err)
	}
	elog.Info().Println("RABBITMQ: Successfully Publishing message")

	// RESPONSE SETUP
	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Set the response header content type
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJson)
}
