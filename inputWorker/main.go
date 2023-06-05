package main

import (
	"log"
	"net/http"

	"github.com/Temctl/E-Notification/inputWorker/controller"
	"github.com/Temctl/E-Notification/util/elog"
	"github.com/gorilla/mux"
)

func startRouter() {
	// -------------------------------------------------------
	// START -------------------------------------------------
	// -------------------------------------------------------

	elog.Info("Start Router ...")
	router := mux.NewRouter()

	// -------------------------------------------------------
	// API ---------------------------------------------------
	// -------------------------------------------------------

	router.HandleFunc("/input", controller.Input).Methods("GET")
	router.HandleFunc("/config", controller.Config).Methods("POST")
	router.HandleFunc("/config", controller.ConfigGet).Methods("GET")

	// -------------------------------------------------------
	// LISTEN ------------------------------------------------
	// -------------------------------------------------------

	http.ListenAndServe(":8085", router)
}

func main() {
	elog.Info("Start server ...")
	log.Println("Fsdfsdf")
	startRouter()
}
