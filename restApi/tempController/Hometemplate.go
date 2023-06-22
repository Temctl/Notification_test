package tempController

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/Temctl/E-Notification/util"
	"github.com/Temctl/E-Notification/util/connections"
	"github.com/Temctl/E-Notification/util/elog"
)

type WorkerInfo struct {
	WorkerName string
	Status     interface{}
}

func HomeTemplateHandler(w http.ResponseWriter, r *http.Request) {
	// Load and parse the template file
	tmpl, err := template.ParseFiles("./template/Home.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	client, err := connections.ConnectionRedis()
	if err != nil {
		elog.Error().Panic(err)
	}
	keys := []string{util.PUSHWORKER, util.NATEMAILWORKER, util.PRIVEMAILWORKER, util.MESSENGERWORKER, util.SMSWORKER}
	// Get the values for the keys
	result, err := client.MGet(keys...).Result()
	if err != nil {
		fmt.Println("Failed to get values:", err)
	}
	myArray := []WorkerInfo{}

	// Print the retrieved values
	for i, val := range result {
		if val == nil {

			fmt.Printf("Key '%s' not found\n", keys[i])
		} else {
			myArray = append(myArray, WorkerInfo{
				WorkerName: keys[i],
				Status:     val,
			})
			fmt.Printf("Value for key '%s': %s\n", keys[i], val)
		}
	}
	data := struct {
		Title   string
		Workers []WorkerInfo
	}{
		Title:   "User List",
		Workers: myArray,
	}
	// Render the template with the provided data
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
