package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Temctl/E-Notification/inputWorker/model"
	"github.com/Temctl/E-Notification/util/elog"
)

// -------------------------------------------------------
// CONFIG api controller code ----------------------------
// -------------------------------------------------------
func Config(w http.ResponseWriter, r *http.Request) {

	elog.Info("Config api start..")
	var info model.ConfigInfo
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	// Close the request body to prevent resource leaks
	defer r.Body.Close()
	// json marshal
	data, err := json.Marshal(info)
	if err != nil {
		log.Println("Failed to serialize struct:", err)
	}
	elog.Info("Writing redis")
	// Redis write heseg
	result := SetRedis("key1", string(data))

	re := GetRedis("key1")

	// Send a response
	fmt.Println(result)
	fmt.Println(re)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(info.Email))
}
