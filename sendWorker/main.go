package sendworker

import (
	"log"
	"os"
	"sync"

	"github.com/Temctl/E-Notification/sendWorker/sender"
	"github.com/Temctl/E-Notification/util/connections"
)

func init() {
	file, _ := os.Create("./log/output.log")
	log.SetOutput(file)
	file.Close()
	log.SetFlags(log.Ldate | log.Lshortfile)

}

func main() {

	redis, err := connections.ConnectionRedis()
	if err != nil {
		//log error
	}
	//get mongodb connection
	mongoClient, err := connections.ConnectMongoDB()
	if err != nil {
		//log
	}

	var wg sync.WaitGroup
	wg.Add(3)
	// ----------------------------------------------------------------------
	// CRON JOB -------------------------------------------------------------
	// ----------------------------------------------------------------------
	go func() {
		defer wg.Done()
		worker.AttentionNotificationEveryday()
	}()
	// ----------------------------------------------------------------------
	// XYP WORKER -----------------------------------------------------------
	// ----------------------------------------------------------------------
	go func() {
		defer wg.Done()
		sender.XypFromDb(mongoClient, redis)
	}()

	// Wait for all goroutines to complete
	wg.Wait()
}

// func main() {

// 	// Example JSON data
// 	jsonData := []byte(`
// {
//     "success": true,
//     "client_is_soap": "1",
//     "organization": {
//         "org_id": "22",
//         "org_phone": "70117676",
//         "contact_phone": "43",
//         "contact_email": "43",
//         "contact_web": "43",
//         "contract_id": "43",
//         "org_type_id": "22",
//         "notifs": [
//         ]
//     }
// }`)

// 	var jsonDataMap map[string]interface{}
// 	err := json.Unmarshal(jsonData, &jsonDataMap)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Check for null or empty values
// 	hasNullEmpty := checkNullEmpty(jsonDataMap)
// 	fmt.Printf("Has null or empty values: %v\n", hasNullEmpty)
// }

// func checkNullEmpty(data interface{}) bool {
// 	switch value := data.(type) {
// 	case nil:
// 		return true
// 	case string:
// 		return value == ""
// 	case []interface{}:
// 		if len(value) == 0 {
// 			return true
// 		}
// 		for _, item := range value {
// 			if checkNullEmpty(item) {
// 				return true
// 			}
// 		}
// 	case map[string]interface{}:
// 		if len(value) == 0 {
// 			return true
// 		}
// 		for _, v := range value {
// 			if checkNullEmpty(v) {
// 				return true
// 			}
// 		}
// 	case []string:
// 		if len(value) == 0 {
// 			return true
// 		}
// 		for _, item := range value {
// 			return item == ""
// 		}
// 	}
// 	return false
// }
