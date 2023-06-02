package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Temctl/E-Notification/inputWorker/controller"
	"github.com/Temctl/E-Notification/util/elog"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func startRouter() {
	router := mux.NewRouter()

	router.HandleFunc("/input", controller.Notification).Methods("GET")
	router.HandleFunc("/config", controller.Config).Methods("POST")

	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, router)
}

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero error")
	}
	return a / b, nil
}

func main() {
	log.Println("Start server ...")
	elog.Info("gogo...")
	numerator := 10
	denominator := 0

	result, err := divide(numerator, denominator)
	if err != nil {
		elog.Error("tom aldaa", err)
	} else {
		fmt.Println("Result:", result)
	}
	log.Println("Start Router")
	startRouter()
}
