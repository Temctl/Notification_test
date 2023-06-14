package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Temctl/E-Notification/util/connections"
	"github.com/Temctl/E-Notification/util/elog"
	umodel "github.com/Temctl/E-Notification/util/model"
)

// -------------------------------------------------------
// SET REDIS ---------------------------------------------
// -------------------------------------------------------
func SetRedis(key string, data map[string]interface{}) bool {

	elog.Info().Println("SET REDIS...")

	// -------------------------------------------------------
	// CONNECTION REDIS CLIENT -----------------------------------
	// -------------------------------------------------------
	client, err := connections.ConnectionRedis()
	if err != nil {
		elog.Error().Println("Failed to create Redis client")
		return false
	}
	// -------------------------------------------------------
	// REDIS luu data oruulah --------------------------------
	// -------------------------------------------------------

	clientErr := client.HMSet("conf:"+key, data).Err()
	if clientErr != nil {
		elog.Error().Println("redis setlehed aldaa garlaa", clientErr)
		return false
	}
	elog.Info().Println("Successful...")

	// -------------------------------------------------------
	// Close the Redis client --------------------------------
	// -------------------------------------------------------

	closeErr := client.Close()
	if closeErr != nil {
		elog.Error().Println("Error closing Redis client:", closeErr)
		return false
	} else {
		elog.Info().Println("Redis client closed successfully")
	}
	return true
}

// ----------------------------------------------------------------------------------
// CONFIG api controller code -------------------------------------------------------
// ----------------------------------------------------------------------------------

// Redist medeelel hiih
func UserConfig(w http.ResponseWriter, r *http.Request) {
	elog.Info().Println("USERCONFIG API STARTED...")

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
	result := SetRedis(redisConfigData.CivilId, fields)

	if !result {
		http.Error(w, "REDIS failed", http.StatusBadRequest)
		return
	}
	// Send a response
	fmt.Println(result)
	fmt.Println(fields)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(redisConfigData.CivilId))
}
