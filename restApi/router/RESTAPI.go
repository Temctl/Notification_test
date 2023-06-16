package router

import (
	"net/http"

	"github.com/Temctl/E-Notification/restApi/auth"
	"github.com/Temctl/E-Notification/restApi/controller"
	"github.com/Temctl/E-Notification/restApi/tempController"
	"github.com/Temctl/E-Notification/util/elog"
	"github.com/gorilla/mux"
)

func RESTAPI() {
	// -----------------------------------------------------------
	// START -----------------------------------------------------
	// -----------------------------------------------------------

	elog.Info().Println("ROUTE STARTED...")
	router := mux.NewRouter()

	// -----------------------------------------------------------
	// API -------------------------------------------------------
	// -----------------------------------------------------------

	// router.HandleFunc("/input", controller.Input).Methods("GET")

	router.HandleFunc("/config", controller.UserConfig).Methods("POST")
	router.HandleFunc("/devicetoken", controller.DeviceTokenConfig).Methods("POST")
	router.HandleFunc("/pushnotif/group", controller.PushNotification).Methods("POST")
	// router.HandleFunc("/config", controller.ConfigGet).Methods("GET")

	// -----------------------------------------------------------
	// TEMPLATE SECTION ------------------------------------------
	// -----------------------------------------------------------

	router.HandleFunc("/login", auth.Login).Methods("POST")
	router.HandleFunc("/login", tempController.LoginTemplateHandler).Methods("GET")
	router.HandleFunc("/", auth.AuthMiddleware(tempController.HomeTemplateHandler)).Methods("GET")

	// -----------------------------------------------------------
	// LISTEN ----------------------------------------------------
	// -----------------------------------------------------------

	http.ListenAndServe(":8085", router)
}
