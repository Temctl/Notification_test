package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Temctl/E-Notification/inputWorker/middleware"
	"github.com/Temctl/E-Notification/util"
	"github.com/Temctl/E-Notification/util/connections"
	"github.com/Temctl/E-Notification/util/elog"
	"github.com/robfig/cron"
)

func IdcardExpire() {
	// Create the request body
	requestBody := middleware.RequestBody{
		ServiceCode: "WS101137_citizenInfoLogByDate",
		CitizenAuthData: middleware.CitizenAuthData{
			Otp: "",
		},
		CustomFields: middleware.CustomFields{
			ObjectCode:  "GET_IDCARD_DATE_OF_EXPIRY_LIST",
			OrgCode:     "10001001",
			OrgName:     "E-mongolia",
			OrgPassword: "aaa1",
			OrgToken:    "aaaaa",
		},
		AuthData: middleware.AuthData{},
		SignData: middleware.SignData{},
	}

	// Convert the request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		elog.Error().Panic(err)
	}
	// Make the POST request
	resp, err := http.Post("https://st-sso.e-mongolia.mn/xyp-api/api/xyp/get-data-no-auth", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Request failed with status code: %d", resp.StatusCode)
		return
	}

	var responseData middleware.ResponseBody
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	db, err := connections.ConnectPostgreSQL()
	if err != nil {
		elog.Error().Panic(err)
	}
	stmt, err := db.Prepare("INSERT INTO attentionnotification (type, regnum, civilid, expiredate) " +
		"VALUES($1, $2, $3, $4)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for _, value := range responseData.Data.Listdata {
		_, err := stmt.Exec("IDCARDGOINGTOEXPIRE", "", value.CivilId, responseData.Data.DateOfExpiry)
		if err != nil {
			elog.Error().Println("Error:", err)
			continue
		}
	}
	stmt.Close()
	db.Close()

	fmt.Println("POST request succeeded")
}
func PassportExpire() {
	// Create the request body
	requestBody := middleware.RequestBody{
		ServiceCode: "WS101137_citizenInfoLogByDate",
		CitizenAuthData: middleware.CitizenAuthData{
			Otp: "",
		},
		CustomFields: middleware.CustomFields{
			ObjectCode:  "GET_PASSPORT_DATE_OF_EXPIRY_LIST",
			OrgCode:     "10001001",
			OrgName:     "E-mongolia",
			OrgPassword: "aaa1",
			OrgToken:    "aaaaa",
		},
		AuthData: middleware.AuthData{},
		SignData: middleware.SignData{},
	}

	// Convert the request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		elog.Error().Panic(err)
	}
	// Make the POST request
	resp, err := http.Post("https://st-sso.e-mongolia.mn/xyp-api/api/xyp/get-data-no-auth", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Request failed with status code: %d", resp.StatusCode)
		return
	}

	var responseData middleware.ResponseBody
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	db, err := connections.ConnectPostgreSQL()
	if err != nil {
		elog.Error().Panic(err)
	}
	stmt, err := db.Prepare("INSERT INTO attentionnotification (type, regnum, civilid, expiredate) " +
		"VALUES($1, $2, $3, $4)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for _, value := range responseData.Data.Listdata {
		_, err := stmt.Exec("IDCARDGOINGTOEXPIRE", "", value.CivilId, responseData.Data.DateOfExpiry)
		if err != nil {
			elog.Error().Println("Error:", err)
			continue
		}
	}
	stmt.Close()
	db.Close()

	fmt.Println("POST request succeeded")
}

func CronJob() {
	go IdcardExpire()
}

func AttentionNotificationEveryday() {
	elog.Info().Println("AttentionNotificationEveryday CRON JOB STARTED...")
	location, err := util.GetTZ()
	if err != nil {
		elog.Error().Panic(err)
	}
	c := cron.NewWithLocation(location)

	// Define the cron job function
	// ----------------------------------------------------------------------
	// Add the cron job to the cron scheduler -------------------------------
	// ----------------------------------------------------------------------
	c.AddFunc("0 6 11 * *", CronJob) // Runs the job at 10:18 AM in GMT+8
	// Start the cron scheduler
	c.Start()

	// Block the program from exiting
	// Use a channel to prevent the main goroutine from exiting
	select {}
}
