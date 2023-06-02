package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Temctl/E-Notification/inputWorker/model"
	"github.com/Temctl/E-Notification/util/elog"
	"github.com/redis/go-redis/v9"
)

// redis luu medeelel oruulah heseg
func WriteRedis(info model.ConfigInfo) bool {
	elog.Info("fdsfsf")
	log.Println("write redis")
	// get ENV info
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	// create redis client
	addr := redisHost + ":" + redisPort
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisPassword, // no password set
		DB:       0,             // use default DB
	})

	// json marshal
	data, err := json.Marshal(info)
	if err != nil {
		log.Println("Failed to serialize struct:", err)
		return false
	}

	// set info in redis
	ctx := context.Background()
	clientErr := client.Set(ctx, info.Email, string(data), 0).Err()
	if clientErr != nil {
		log.Println(clientErr)
		return false
	} else {
		log.Println("Successful")
	}

	// Close the Redis client
	closeErr := client.Close()
	if closeErr != nil {
		fmt.Println("Error closing Redis client:", closeErr)
		return false
	} else {
		fmt.Println("Redis client closed successfully")
	}
	return true
}

// CONFIG api controller code
func Config(w http.ResponseWriter, r *http.Request) {
	log.Println("config")
	var info model.ConfigInfo
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	// Close the request body to prevent resource leaks
	defer r.Body.Close()
	log.Println("write redis")
	// Redis write heseg
	result := WriteRedis(info)

	log.Println("?fdfd ez")
	// Send a response
	fmt.Println(result)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(info.Email))
}
