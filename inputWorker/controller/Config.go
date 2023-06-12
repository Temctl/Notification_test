package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Temctl/E-Notification/util/elog"
	umodel "github.com/Temctl/E-Notification/util/model"
	"github.com/Temctl/E-Notification/util/redis"
)

// -------------------------------------------------------
// SET REDIS ---------------------------------------------
// -------------------------------------------------------
func SetRedis(key string, data string) bool {

	elog.Info().Println("set redis...")

	// -------------------------------------------------------
	// CONNECTION REDIS CLIENT -----------------------------------
	// -------------------------------------------------------
	client := redis.ConnectionRedis()
	if client == nil {
		elog.Error().Println("Failed to create Redis client")
		return false
	}
	// -------------------------------------------------------
	// REDIS luu data oruulah --------------------------------
	// -------------------------------------------------------

	ctx := context.Background()
	clientErr := client.HSet(ctx, "conf:"+key, data).Err()
	if clientErr != nil {
		elog.Error().Println("redis setlehed aldaa garlaa", clientErr)
		return false
	} else {
		elog.Info().Println("Successful...")
	}

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

// -------------------------------------------------------
// GET REDIS ---------------------------------------------
// -------------------------------------------------------
func GetRedis(key string) string {
	elog.Info().Println("get redis...")
	// -------------------------------------------------------
	// CONNECTION REDIS CLIENT -----------------------------------
	// -------------------------------------------------------
	client := redis.ConnectionRedis()
	if client == nil {
		fmt.Println("Failed to create Redis client")
		return ""
	}

	// -------------------------------------------------------
	// REDIS ees data avah -----------------------------------
	// -------------------------------------------------------

	ctx := context.Background()
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		elog.Error().Println("redis client get data: ", err)
		return ""
	}

	// -------------------------------------------------------
	// Close the Redis client --------------------------------
	// -------------------------------------------------------

	closeErr := client.Close()
	if closeErr != nil {
		elog.Error().Println("Error closing Redis client:", closeErr)
		return ""
	} else {
		elog.Info().Println("Redis client closed successfully")
	}

	return val
}

// ----------------------------------------------------------------------------------
// CONFIG api controller code -------------------------------------------------------
// ----------------------------------------------------------------------------------

// Redist medeelel hiih
func Config(w http.ResponseWriter, r *http.Request) {
	elog.Info().Println("Config api start..")
	// var info model.ConfigInfo
	var redisConfigData umodel.RedisConfigNotification
	err := json.NewDecoder(r.Body).Decode(&redisConfigData)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	// Close the request body to prevent resource leaks
	defer r.Body.Close()
	// json marshal
	data, err := json.Marshal(redisConfigData)
	if err != nil {
		log.Println("Failed to serialize struct:", err)
	}
	elog.Info().Println("Writing redis")
	// Redis write heseg
	result := SetRedis(redisConfigData.CivilId, string(data))

	// Send a response
	fmt.Println(result)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(redisConfigData.CivilId))
}

// Redis medeelel avah
func ConfigGet(w http.ResponseWriter, r *http.Request) {

	elog.Info().Println("Config api start..")

	// Get the value of the "name" parameter from the query string
	keys := r.FormValue("keys")
	// Check if the parameter is present
	if keys == "" {
		elog.Warning().Println("Medeelel alga...")
	}

	result := GetRedis(keys)

	// Send a response
	fmt.Println(result)
	w.WriteHeader(http.StatusOK)
}
