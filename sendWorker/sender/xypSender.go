package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/Temctl/E-Notification/util/connections"
	"github.com/Temctl/E-Notification/util/model"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GroupResult struct {
	ID          XypUserInfo        `bson:"_id"`
	Count       int                `bson:"count"`
	ContentData []model.XypContent `bson:"data"`
}

type XypUserInfo struct {
	Regnum         string `json:"regnum"`
	OperatorRegnum string `json:"operatorRegnum"`
	CivilId        string `json:"civilId"`
	ClientId       int    `json:"clientId"`
}

func XypFromDb(mongoClient *mongo.Client, redis *redis.Client, emailsFile []byte) {
	var collectionName model.Collections
	collectionName = model.XYPNOTIFICATION
	collection, _, err := connections.GetMongoCollection(collectionName)
	if err != nil {
		//error log
	}
	// Construct the aggregation pipeline
	// MongoDB aggregation pipeline
	pipeline := []bson.M{
		{"$limit": 500},
		{
			"$group": bson.M{
				"_id": bson.M{
					"regnum":   "$regnum",
					"civilId":  "$civilId",
					"clientId": "$clientId",
				},
				"count":       bson.M{"$sum": 1},
				"contentData": bson.M{"$push": "$contentData"},
			},
		},
	}

	for {
		// Execute the aggregation query
		cursor, err := collection.Aggregate(context.Background(), pipeline)
		if err != nil {
			log.Fatal(err)
		}

		// Iterate over the cursor and delete each document
		for cursor.Next(context.Background()) {
			var wg sync.WaitGroup
			var request GroupResult
			if err := cursor.Decode(&request); err != nil {
				log.Fatal(err)
			} else {
				wg.Add(1)
				go func(notif GroupResult, wg *sync.WaitGroup) {
					defer wg.Done()
					if notif.ID.CivilId == "" && notif.ID.Regnum == "" {

					} else {
						if notif.ID.CivilId == "" {
							exists, err := redis.Exists("getByReg:" + notif.ID.Regnum).Result()
							if err != nil {
								panic(err)
							} else if exists == 1 {
								notif.ID.CivilId, err = redis.Get("getByReg:" + notif.ID.Regnum).Result() // if civil id is not sent, get it using regnum from redis conf
								if err != nil {
									panic(err)
								}
							}
						}
						userConf, err := redis.HGetAll("conf:" + notif.ID.CivilId).Result()
						if err != nil {
							panic(err)
						}
						sendMq(notif, userConf, redis, emailsFile)
					}
				}(request, &wg)
			}

			// Define the delete filter
			filter := bson.M{
				"regnum":   request.ID.Regnum,
				"civilId":  request.ID.CivilId,
				"clientid": request.ID.ClientId,
			}

			// Delete the document
			_, err = collection.DeleteMany(context.Background(), filter)
			if err != nil {
				log.Fatal(err)
			}
			// Wait for all goroutines to complete
			wg.Wait()
		}

		cursor.Close(context.Background())
	}

}

func sendMq(request GroupResult, userConf map[string]string, redis *redis.Client, emailsFile []byte) {
	if isPush, ok := userConf["isPush"]; ok && isPush == "1" {
		deviceTokensExists, err := redis.Exists("deviceTokens:" + request.ID.CivilId).Result()
		if err != nil {
			panic(err)
		} else if deviceTokensExists == 1 { //if it exists
			data, err := redis.Get("deviceTokens:" + request.ID.CivilId).Result() //get device token from redis
			fmt.Println(data)
			if err != nil {
				panic(err)
			} else {
				var deviceTokens []string
				err := json.Unmarshal([]byte(data), &deviceTokens) //redis is returning string, so turn the string array into slice
				if err != nil {
					panic(err)
				}
				request := generatePushRequest("regular notif test", "medeeleld haldsan baina", request.ID.CivilId, request.ID.Regnum, deviceTokens)
				MqPush(request)
			}
		} // TODO when is push is open but no device token found
	}
	if isNationalEmail, ok := userConf["isNationalEmail"]; ok && isNationalEmail == "1" {
		XypNatEmail(civilId, content)
	}
	if isEmail, ok := userConf["isEmail"]; ok && isEmail == "1" {
		if emailAddress, ok := userConf["emailAddress"]; ok || emailAddress != "" {
			result += RegularPrivEmail(civilId, content)
		}
	}
	if isSocial, ok := userConf["social"]; ok && isSocial == "1" {
		SendSocial(civilId)
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
}

func generatePushRequest(title string, body string, civilId string, regnum string, deviceTokens []string) model.RegularNotificationModel {
	var request model.RegularNotificationModel
	request.Body = body
	request.Title = title
	request.Type = model.NotificationType(0)
	request.CivilIds = append(request.CivilIds, civilId)
	request.Regnums = append(request.Regnums, regnum)
	request.Tokens = deviceTokens
	return request
}

func getOrgInfo(clientId int, redis *redis.Client) (map[string]interface{}, string) {
	orgType := "GOVERNMENT_COMPANY"
	url := "https://10.10.18.28/api/service/notifications"
	var orgInfo map[string]interface{}

	if redis.Get("orgInfo:"+string(clientId)) != nil {
		log.Println("already in redis")
		orgInfo = json.loads(redis.Get("orgInfo:" + string(clientId)))
		if orgInfo["client_is_soap"] == "0" {
			orgType = "USER"
		} else if orgInfo["organization"]["org_type_id"] != "1" {
			orgType = "PRIVATE_COMPANY"
		}
		if !checkNull(org_info) {
			isDefault = false
		}
	} else {
		payload := fmt.Sprintf("client_id=%d", clientId)
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

			isDefault = checkNullEmpty(orgInfo)
			if orgInfo["success"] {
				redisLocal.Set(clientId, json.dumps(orgInfo))
				if orgInfo["client_is_soap"] == "0" {
					orgType = "USER"
				} else if orgInfo["organization"]["org_type_id"] != "1" {
					orgType = "PRIVATE_COMPANY"
				}
				isDefault = false
			}
		}
	}
	if clientId == 61 {
		orgType = "KIOSK"
	}
	return orgInfo, orgType
}

func generateXypHtml(request model.XypNotification) {

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

func checkNullEmpty(data interface{}) bool {
	if data == nil {
		return true
	}
	switch value := data.(type) {
	case nil:
		return true
	case string:
		return value == ""
	case []interface{}:
		if len(value) == 0 {
			return true
		}
		for _, item := range value {
			if checkNullEmpty(item) {
				return true
			}
		}
	case map[string]interface{}:
		if len(value) == 0 {
			return true
		}
		for _, v := range value {
			if checkNullEmpty(v) {
				return true
			}
		}
	case []string:
		if len(value) == 0 {
			return true
		}
		for _, item := range value {
			return item == ""
		}
	}
	return false
}
