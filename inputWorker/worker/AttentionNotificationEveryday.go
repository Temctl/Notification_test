package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Temctl/E-Notification/inputWorker/middleware"
	"github.com/Temctl/E-Notification/util"
	"github.com/Temctl/E-Notification/util/connections"
	"github.com/Temctl/E-Notification/util/elog"
	"github.com/Temctl/E-Notification/util/model"
	"github.com/robfig/cron"
)

func IdcardExpire() {
	// Create the request body
	requestBody := middleware.RequestBody{
		ServiceCode: util.ATTENTION_SERVICENAME,
		CitizenAuthData: middleware.CitizenAuthData{
			Otp: "",
		},
		CustomFields: middleware.CustomFields{
			ObjectCode:  util.OBJECTCODE,
			OrgCode:     util.ORGCODE,
			OrgName:     util.ORGNAME,
			OrgPassword: util.ORGPASSWORD,
			OrgToken:    util.ORGTOKEN,
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
	resp, err := http.Post(util.ATTENTION_URL, "application/json", bytes.NewBuffer(jsonBody))
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
	client, err := connections.ConnectMongoDB()
	if err != nil {
		elog.Error().Panic(err)
	}
	collection := client.Database("notification").Collection("attentionnotification")
	for _, value := range responseData.Data.Listdata {
		data := model.AttentionNotification{
			Regnum:     "",
			CivilId:    value.CivilId,
			ExpireDate: responseData.Data.DateOfExpiry,
			Type:       "IDCARDGOINGTOEXPIRE",
			Passport:   value.Passport,
		}
		// Insert the document
		_, insertErr := collection.InsertOne(context.Background(), data)
		if insertErr != nil {
			elog.Error().Panic(insertErr)
		}
	}

	client.Disconnect(context.Background())
	fmt.Println("POST request succeeded")
}
func PassportExpire() {
	// Create the request body
	requestBody := middleware.RequestBody{
		ServiceCode: util.ATTENTION_SERVICENAME,
		CitizenAuthData: middleware.CitizenAuthData{
			Otp: "",
		},
		CustomFields: middleware.CustomFields{
			ObjectCode:  util.OBJECTCODE,
			OrgCode:     util.ORGCODE,
			OrgName:     util.ORGNAME,
			OrgPassword: util.ORGPASSWORD,
			OrgToken:    util.ORGTOKEN,
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
	resp, err := http.Post(util.ATTENTION_URL, "application/json", bytes.NewBuffer(jsonBody))
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
	client, err := connections.ConnectMongoDB()
	if err != nil {
		elog.Error().Panic(err)
	}
	collection := client.Database("notification").Collection("attentionnotification")
	for _, value := range responseData.Data.Listdata {
		data := model.AttentionNotification{
			Regnum:     "",
			CivilId:    value.CivilId,
			ExpireDate: responseData.Data.DateOfExpiry,
			Type:       "INTPASSPORTGOINGTOEXPIRE",
			Passport:   value.Passport,
		}
		// Insert the document
		_, insertErr := collection.InsertOne(context.Background(), data)
		if insertErr != nil {
			elog.Error().Panic(insertErr)
		}
	}
	client.Disconnect(context.Background())
	fmt.Println("POST request succeeded")
}

func DriverLicenseExpire() {
	// Create the request body
	// Create the request body
	requestBody := middleware.RequestBody{
		ServiceCode: util.ATTENTION_SERVICENAME2,
		CitizenAuthData: middleware.CitizenAuthData{
			Otp: "",
		},
		CustomFields: middleware.CustomFields{},
		AuthData:     middleware.AuthData{},
		SignData:     middleware.SignData{},
	}

	// Convert the request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		elog.Error().Panic(err)
	}
	// Make the POST request
	resp, err := http.Post(util.ATTENTION_URL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Request failed with status code: %d", resp.StatusCode)
		return
	}

	var responseData middleware.DResponseBody
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	client, err := connections.ConnectMongoDB()
	if err != nil {
		elog.Error().Panic(err)
	}
	collection := client.Database("notification").Collection("attentionnotification")
	for _, value := range responseData.Data.Listdata {
		data := model.AttentionNotification{
			Regnum:     value.Regnum,
			CivilId:    "",
			ExpireDate: value.ExpirationDate,
			Type:       "DRIVERLICENSEEXPIRED",
			Passport:   "",
		}
		// Insert the document
		_, insertErr := collection.InsertOne(context.Background(), data)
		if insertErr != nil {
			elog.Error().Panic(insertErr)
		}
	}

	client.Disconnect(context.Background())
	fmt.Println("POST request succeeded")
}

func CronJob() {
	go IdcardExpire()
	go PassportExpire()
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
	c.AddFunc("0 31 14 * *", CronJob) // Runs the job at 10:18 AM in GMT+8
	c.AddFunc("0 42 14 21 *", DriverLicenseExpire)
	// Start the cron scheduler
	c.Start()

	// Block the program from exiting
	// Use a channel to prevent the main goroutine from exiting
	select {}
}
