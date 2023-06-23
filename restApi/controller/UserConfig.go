package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Temctl/E-Notification/inputWorker/model"
	"github.com/Temctl/E-Notification/util/connections"
	"github.com/Temctl/E-Notification/util/elog"
	umodel "github.com/Temctl/E-Notification/util/model"
)

// -------------------------------------------------------
// SET REDIS ---------------------------------------------
// -------------------------------------------------------
func SetRedis(regnum string, key string, data map[string]interface{}) error {

	elog.Info().Println("SET REDIS...")

	// -------------------------------------------------------
	// CONNECTION REDIS CLIENT -------------------------------
	// -------------------------------------------------------
	client, err := connections.ConnectionRedis()
	if err != nil {
		elog.Error().Println("Failed to create Redis client")
		return err
	}
	// -------------------------------------------------------
	// REDIS luu data oruulah --------------------------------
	// -------------------------------------------------------
	result, existErr := client.Exists("conf:" + key).Result()
	if existErr != nil {
		elog.Error().Println("Register shalgahad aldaa garlaa", existErr)
		return existErr
	} else if result == 0 {
		regErr := client.Set("getByReg:"+regnum, key, 0).Err()
		if regErr != nil {
			elog.Error().Println(regErr)
			return regErr
		}
		elog.Info().Println("Successful...")
	}
	clientErr := client.HMSet("conf:"+key, data).Err()
	if clientErr != nil {
		elog.Error().Println("redis setlehed aldaa garlaa", clientErr)
		return clientErr
	}
	elog.Info().Println("Successful...")

	// -------------------------------------------------------
	// Close the Redis client --------------------------------
	// -------------------------------------------------------

	closeErr := client.Close()
	if closeErr != nil {
		elog.Error().Println("Error closing Redis client:", closeErr)
		return closeErr
	}

	elog.Info().Println("Redis client closed successfully")

	return nil

}

// ----------------------------------------------------------------------------------
// CONFIG api controller code -------------------------------------------------------
// ----------------------------------------------------------------------------------

// Redist medeelel hiih
func UserConfig(w http.ResponseWriter, r *http.Request) {
	elog.Info().Println("USERCONFIG API STARTED...")
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
	var redisConfigData umodel.UserConfigNotification
	err := json.NewDecoder(r.Body).Decode(&redisConfigData)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	// Close the request body to prevent resource leaks
	defer r.Body.Close()

	// Create a map with field-value pairs
	fields := map[string]interface{}{
		"civilId":         redisConfigData.CivilId,
		"regnum":          redisConfigData.Regnum,
		"isSms":           redisConfigData.IsSms,
		"isEmail":         redisConfigData.IsEmail,
		"isPush":          redisConfigData.IsPush,
		"isNationalEmail": redisConfigData.IsNationalEmail,
		"emailAddress":    redisConfigData.EmailAddress,
		"social":          redisConfigData.Social,
	}

	elog.Info().Println("WRITING REDIS")
	// Redis write heseg
	resultErr := SetRedis(redisConfigData.Regnum, redisConfigData.CivilId, fields)

	if resultErr != nil {
		response.Message = "АМЖИЛТГҮЙ: " + resultErr.Error()
		response.Status = http.StatusBadRequest
	}
	// -----------------------------
	// RESPONSE SETUP
	// -----------------------------
	responseJson, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Set the response header content type
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJson)
}
