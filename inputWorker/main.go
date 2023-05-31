package main

import (
	"fmt"
	"log"
	"net/http"
	"notifson/controller"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
func start() {
	router := mux.NewRouter()

	router.HandleFunc("/notification", controller.Notification).Methods("GET")
	router.HandleFunc("/config", controller.Config).Methods("POST")

	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, router)
}

func main() {
	fmt.Println("Start server ...")
	start()
}
