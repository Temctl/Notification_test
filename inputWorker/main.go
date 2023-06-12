package main

import (
	"net/http"

	"github.com/Temctl/E-Notification/inputWorker/auth"
	"github.com/Temctl/E-Notification/inputWorker/controller"
	"github.com/Temctl/E-Notification/util/elog"
	"github.com/gorilla/mux"
)

func startRouter() {
	// -------------------------------------------------------
	// START -------------------------------------------------
	// -------------------------------------------------------

	elog.Info().Println("Start Router ...")
	router := mux.NewRouter()

	// -------------------------------------------------------
	// API ---------------------------------------------------
	// -------------------------------------------------------

	router.HandleFunc("/input", controller.Input).Methods("GET")

	router.HandleFunc("/config", controller.Config).Methods("POST")
	router.HandleFunc("/config", auth.AuthMiddleware(controller.ConfigGet)).Methods("GET")

	// Template dd
	router.HandleFunc("/login", auth.Login).Methods("POST")
	router.HandleFunc("/login", controller.LoginTemplateHandler).Methods("GET")
	router.HandleFunc("/", controller.HomeTemplateHandler).Methods("GET")
	// -------------------------------------------------------
	// LISTEN ------------------------------------------------
	// -------------------------------------------------------

	http.ListenAndServe(":8085", router)
}

func main() {
	elog.Info().Println("Start server ...")
	startRouter()
}
