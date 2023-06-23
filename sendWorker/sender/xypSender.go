package sender

import (
	"context"
	"log"
	"sync"

	"github.com/Temctl/E-Notification/util/connections"
	"github.com/Temctl/E-Notification/util/model"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func XypFromDb(mongoClient *mongo.Client, redis *redis.Client) {
	var collectionName model.Collections
	collectionName = model.XYPNOTIFICATION
	collection, client, err := connections.GetMongoCollection(collectionName)
	if err != nil {
		//error log
	}
	// Construct the aggregation pipeline
	// pipeline := mongo.Pipeline{
	// 	{{Key: "$group", Value: bson.D{
	// 		{"_id", bson.D{
	// 			{"regnum", "$regnum"},
	// 			{"civilId", "$civilId"},
	// 			{"clientid", "$clientid"},
	// 		}},
	// 		{"dt", bson.D{
	// 			{"$push", bson.D{
	// 				{"$concat", bson.A{
	// 					"$orgName", "|", "$serviceDesc", "|", "$date", "|", "$serviceName", "|", "$requestId",
	// 				}},
	// 			}},
	// 		}},
	// 		{"max_id", bson.D{
	// 			{"$max", "$id"},
	// 		}},
	// 	}}},
	// }

	for{
		// Execute the aggregation query
		cursor, err := collection.Aggregate(context.Background(), pipeline)
		if err != nil {
			log.Fatal(err)
		}

		// Iterate over the cursor and delete each document
		for cursor.Next(context.Background()) {
			var request model.XypNotification
			if err := cursor.Decode(&request); err != nil {
				log.Fatal(err)
			}else {
				wg.Add(1)
				go func(notif model.XypNotification) {
					defer wg.Done()
					if notif.CivilId == "" && notif.Regnum == "" {

					} else {
						if notif.CivilId == "" {
							exists, err := redis.Exists("getByReg:" + notif.Regnum).Result()
							if err != nil {
								panic(err)
							} else if exists == 1 {
								notif.CivilId, err = redis.Get("getByReg:" + notif.Regnum).Result() // if civil id is not sent, get it using regnum from redis conf
								if err != nil {
									panic(err)
								}
							}
						}
						userConf, err := redis.HGetAll("conf:" + request.CivilId).Result()
						if err != nil {
							panic(err)
						}
						sendMq(notif, userConf)
					}
				}(request)
			}

			// Define the delete filter
			filter := bson.M{
				"regnum":   request.Regnum,
				"civilId":  request.CivilId,
				"clientid": request.ClientId,
			}

			// Delete the document
			_, err = collection.DeleteMany(context.Background(), filter)
			if err != nil {
				log.Fatal(err)
			}
		}
		// Wait for all goroutines to complete
		wg.Wait()
		
		cursor.Close(context.Background())
	}
	

}


func sendMq(request model.XypNotification, userConf map[string]string) {
	if isPush, ok := userConf["isPush"]; ok && isPush == "1" {
		deviceTokensExists, err := redis.Exists("deviceTokens:" + civilId).Result()
		if err != nil {
			panic(err)
		} else if deviceTokensExists == 1 { //if it exists
			data, err := redis.Get("deviceTokens:" + civilId).Result() //get device token from redis
			fmt.Println(data)
			if err != nil {
				panic(err)
			} else {
				var deviceTokens []string
				err := json.Unmarshal([]byte(data), &deviceTokens) //redis is returning string, so turn the string array into slice
				if err != nil {
					panic(err)
				}
				request := generatePushRequest("regular notif test", "medeeleld haldsan baina", notifificationType)
				PushToTokens(request, deviceTokens, client)
			}
		} // TODO when is push is open but no device token found
	}
	if isNationalEmail, ok := userConf["isNationalEmail"]; ok && isNationalEmail == "1" {
		result += XypNatEmail(civilId, content)
	}
	if isEmail, ok := userConf["isEmail"]; ok && isEmail == "1" {
		if emailAddress, ok := userConf["emailAddress"]; ok || emailAddress != "" {
			result += RegularPrivEmail(civilId, content)
		}
	}
	if isSocial, ok := userConf["social"]; ok && isSocial == "1" {
		SendSocial(civilId)
	}
	} else {
	deviceTokensExists, err := redis.Exists("deviceTokens:" + civilId).Result()
	if err != nil {
		panic(err)
	} else if deviceTokensExists == 1 {
		userDeviceTokens, err := redis.LRange("deviceTokens:"+civilId, 0, -1).Result()
		if err != nil {
			panic(err)
		} else {
			request := generatePushRequest("regular notif test", "medeeleld haldsan baina", notifificationType)
			PushToTokens(request, userDeviceTokens, client)
		}
	}
}

func getOrgInfo(clientId int, redis *redis.Client){
	orgType := "GOVERNMENT_COMPANY"
	companyName := strings.Split(dataArr[0], "|")[0]
	isDefault := true
	url := "https://10.10.18.28/api/service/notifications"
	var orgInfo map[string]interface{}

	if redisClient.Get(clientID) != nil {
		log.Println("already in redis")
		orgInfo = json.loads(redisClient.Get(clientID))
		if orgInfo["client_is_soap"] == "0" {
			orgType = "USER"
		} else if orgInfo["organization"]["org_type_id"] != "1" {
			orgType = "PRIVATE_COMPANY"
		}
		if !checkNull(org_info){
			isDefault = false
		}
	} else {
		payload := fmt.Sprintf("client_id=%d", clientID)
		headers := map[string]string{
			"Authorization": "Bearer 101|rCb9yOIYDxUooLyuaPL1VrdgEnbcUA2AMUoCuOXo",
			"Content-Type":  "application/x-www-form-urlencoded",
		}

		req, err := http.NewRequest("POST", url, strings.NewReader(payload))
		if err != nil {
			log.Fatal(err)
		}
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			orgInfo = string(body)
			err = json.Unmarshal([]byte(orgInfo), &orgInfo)
			if err != nil {
				log.Fatal(err)
			}

			isDefault = checkNull(orgInfo)
			if orgInfo["success"] {
				redisLocal.Set(clientID, json.dumps(orgInfo))
				if orgInfo["client_is_soap"] == "0" {
					orgType = "USER"
				} else if orgInfo["organization"]["org_type_id"] != "1" {
					orgType = "PRIVATE_COMPANY"
				}
				isDefault = false
			}
		}
	}
	if clientID == 61 {
		orgType = "KIOSK"
	}
}

func generateXypHtml(request model.XypNotification){
	
}

func emailNationalDefault(pool *Pool, dataArr []string, civilID, baseHTML, regNumber, seqID, clientID string) (bool, error) {
	dataObj := make(map[string][][]string)
	for _, d := range dataArr {
		data := strings.Split(d, "|")
		key := data[0]
		if _, ok := dataObj[key]; ok {
			dataObj[key] = append(dataObj[key], data[1:])
		} else {
			dataObj[key] = [][]string{data[1:]}
		}
	}

	rows := ""
	table := ""
	companyName := ""
	for key := range dataObj {
		companyName = key
		tableRows := ""
		for _, subCont := range dataObj[key] {
			tableRows += fmt.Sprintf(SERVICE_HTML_TOGGLE_DEFAULT, subCont[0], subCont[1], subCont[3])
		}
		table += fmt.Sprintf(TABLE_DEFAULT, companyName, tableRows)
		rows += table
	}

	sendMsg := strings.ReplaceAll(baseHTML, "{TEXT_REPLACE}", DEFAULT_TEXT)
	sendMsg = strings.ReplaceAll(sendMsg, "{REGNUMBER_REPLACE}", civilID)
	sendMsg = strings.ReplaceAll(sendMsg, "{TABLES}", rows)
	sendMsg = strings.ReplaceAll(sendMsg, "{FOOTER_TEXT}", FOOTER_TEXT_DEFAULT)

	// f, err := os.Create("test.html")
	// if err != nil {
	// 	return false, err
	// }
	// f.WriteString(sendMsg)
	// f.Close()

	result, err := sendEmail(pool, sendMsg, regNumber, civilID, seqID, dataArr, clientID)
	if err != nil {
		return false, err
	}

	if result {
		err = writeLog(pool, dataArr, civilID, regNumber, "NATIONAL_EMAIL", "DEFAULT", companyName)
		if err != nil {
			return false, err
		}
	}

	return result, nil
}


func emailNational(pool *Pool, dataArr []string, civilID, baseHTML, regNumber, seqID, clientID string, orgInfo map[string]interface{}, isDefault bool, redisClient *RedisClient) (bool, error) {
	isDAN := false
	isPrivClient := false
	companyName := ""
	dataObj := make(map[string][][]string)

	if isDefault {
		return emailNationalDefault(pool, dataArr, civilID, baseHTML, regNumber, seqID, clientID)
	}

	orgSuccess := orgInfo["success"].(bool)
	if orgSuccess {
		clientSOAP := orgInfo["client_is_soap"].(string)
		if clientSOAP == "0" {
			isDAN = true
		} else if orgInfo["organization"].(map[string]interface{})["org_type_id"].(string) != "1" {
			isPrivClient = true
		}
		orgInfo = orgInfo["organization"].(map[string]interface{})
	} else {
		return emailNationalDefault(pool, dataArr, civilID, baseHTML, regNumber, seqID, clientID)
	}

	for _, d := range dataArr {
		data := strings.Split(d, "|")
		key := data[0]
		if _, ok := dataObj[key]; ok {
			dataObj[key] = append(dataObj[key], data[1:])
		} else {
			dataObj[key] = [][]string{data[1:]}
		}
	}

	rows := ""
	count := 1
	for key := range dataObj {
		companyName = key
		tableRows := ""
		for _, subCont := range dataObj[key] {
			if orgInfo != nil {
				table := ""
				notifs := orgInfo["notifs"].([]interface{})
				for _, notif := range notifs {
					if notif.(map[string]interface{})["ws_operation_name"].(string) == subCont[2] {
						text := "Тухайн мэдээллийг мэдээллийн эзний зөвшөөрлийг авахгүйгээр ашиглах хууль зүйн үндэслэл"
						description := ""
						if !isPrivClient && !isDAN {
							description = fmt.Sprintf(DESCRIPTION_ROW, text, notif.(map[string]interface{})["description"].(string))
						}
						if isPrivClient {
							text = "Гэрээний дугаар"
							description = fmt.Sprintf(DESCRIPTION_ROW, text, notif.(map[string]interface{})["description"].(string))
						}
						tableRows += fmt.Sprintf(SERVICE_HTML_TOGGLE, count, subCont[0], count, description, notif.(map[string]interface{})["reason"].(string), subCont[1], subCont[3])
						count++
					}
				}
			}
		}
		table += fmt.Sprintf(TABLE, companyName, tableRows, orgInfo["contact_phone"].(string), orgInfo["contact_email"].(string))
		rows += table
	}

	orgType := "GOVERNMENT_COMPANY"
	text := ""
	footer := FOOTER_TEXT_NATIONAL
	if clientID == "61" {
		text = TUTS_TEXT
		footer = FOOTER_TEXT_DEFAULT
		orgType = "KIOSK"
	} else {
		if isDAN {
			text = strings.ReplaceAll(USER_TEXT, "{COMPANY_NAME_REPLACE}", companyName)
			footer = FOOTER_TEXT_USER
			orgType = "USER"
		} else if isPrivClient {
			text = strings.ReplaceAll(PRIV_COMPANY, "{COMPANY_NAME_REPLACE}", companyName)
			footer = FOOTER_TEXT_PRIV
			orgType = "PRIVATE_COMPANY"
		} else {
			text = GOV_COMPANY
		}
	}
	sendMsg := strings.ReplaceAll(baseHTML, "{TEXT_REPLACE}", text)
	sendMsg = strings.ReplaceAll(sendMsg, "{REGNUMBER_REPLACE}", civilID)
	sendMsg = strings.ReplaceAll(sendMsg, "{TABLES}", rows)
	sendMsg = strings.ReplaceAll(sendMsg, "{FOOTER_TEXT}", footer)

	result, err := sendEmail(pool, sendMsg, regNumber, civilID, seqID, dataArr, clientID)
	if err != nil {
		return false, err
	}

	if result {
		err = writeLog(pool, dataArr, civilID, regNumber, "NATIONAL_EMAIL", orgType, companyName)
		if err != nil {
			return false, err
		}
	}

	return result, nil
}

func checkNull(data map[string]interface{}) bool {
	for _, value := range data {
		switch val := value.(type) {
		case string:
			if val == "" {
				return true
			}
		case map[string]interface{}:
			isNull := checkNull(val)
			if isNull {
				return true
			}
		case []interface{}:
			if len(val) > 0 {
				for _, v := range val {
					isNull := checkNull(v.(map[string]interface{}))
					if isNull {
						return true
					}
				}
			} else {
				return true
			}
		}
	}
	return false
}
