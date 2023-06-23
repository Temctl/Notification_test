package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Temctl/E-Notification/restApi/model"
	"github.com/Temctl/E-Notification/util/elog"
	umodel "github.com/Temctl/E-Notification/util/model"
)

func PushNotification(w http.ResponseWriter, r *http.Request) {
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
	var configData umodel.RegularNotificationModel
	err := json.NewDecoder(r.Body).Decode(&configData)
	if err != nil {
		elog.Error().Println(err)
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	// Close the request body to prevent resource leaks
	defer r.Body.Close()

	// configDataJson, err := json.Marshal(configData)
	// if err != nil {
	// 	elog.Error().Println(err)
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// // ----------------------------------------------------------------------
	// // RABBITMQ CONNECTION --------------------------------------------------
	// // ----------------------------------------------------------------------
	// queue, rErr := connections.GetRabbitmqChannel()
	// if rErr != nil {
	// 	elog.Error().Println(rErr)
	// }
	// err = queue.Publish(
	// 	"",
	// 	util.NATEMAILKEY,
	// 	false,
	// 	false,
	// 	amqp.Publishing{
	// 		ContentType: "text/plain",
	// 		Body:        []byte(configDataJson),
	// 	},
	// )
	// if err != nil {
	// 	elog.Error().Println("Publish error", err)
	// }
	// elog.Info().Println("RABBITMQ: Successfully Publishing message")

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
