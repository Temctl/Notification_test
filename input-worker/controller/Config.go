package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"notifson/model"
	"os"

	"github.com/redis/go-redis/v9"
)

func WriteRedis(info model.ConfigInfo) {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	addr := redisHost + ":" + redisPort

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisPassword, // no password set
		DB:       0,             // use default DB
	})

	data, err := json.Marshal(info)
	if err != nil {
		log.Fatal("Failed to serialize struct:", err)
	}

	ctx := context.Background()
	clientErr := client.Set(ctx, info.Email, string(data), 0).Err()
	if clientErr != nil {
		panic(clientErr)
	}
	fmt.Println("Success...")
}

func Config(w http.ResponseWriter, r *http.Request) {
	fmt.Println("config")
	var info model.ConfigInfo
	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	// Close the request body to prevent resource leaks
	defer r.Body.Close()

	// Redis write heseg
	WriteRedis(info)

	// Send a response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(info.Email))
}
