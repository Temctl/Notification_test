package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
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
			var sentCount int
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
						exists, err := redis.Exists("conf:" + notif.ID.CivilId).Result()
						if err != nil {
							panic(err)
						} else if exists == 1 {
							userConf, err := redis.HGetAll("conf:" + notif.ID.CivilId).Result()
							if err != nil {
								panic(err)
							}
							sentCount = sendMq(notif, userConf, redis, emailsFile)
						} else {
							//onlly push and national email
							sentCount = sendMq(notif, nil, redis, emailsFile)
						}

					}
				}(request, &wg)
			}

			if sentCount > 0 {
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
			}

			// Wait for all goroutines to complete
			wg.Wait()
		}

		cursor.Close(context.Background())
	}

}

func sendMq(request GroupResult, userConf map[string]string, redis *redis.Client, emailsFile []byte) int {
	sentCount := 0
	text := "Иргэн таны мэдээлэлд \"" + request.ContentData[0].OrgName + "\" байгууллага хандсаныг үүгээр мэдэгдэж байна. Дэлгэрэнгүй мэдээлэл харахыг хүсвэл e-Mongolia системийн и-мэйл цэснээс харна уу."
	if userConf != nil{
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
					request := generatePushRequest(text, "medeeleld haldsan baina", request.ID.CivilId, request.ID.Regnum, deviceTokens)
					MqPush(request)
					sentCount += 1
				}
			} // TODO when is push is open but no device token found
		}
		if isNationalEmail, ok := userConf["isNationalEmail"]; ok && isNationalEmail == "1" {
			orgInfo, orgType, isDefault := getOrgInfo(request.ID.ClientId, redis)
			html := generateXypHtml(request, orgInfo, orgType, isDefault)
			var email model.EmailModel
			email.Body = html
			email.CivilId = request.ID.CivilId
			email.Regnum = request.ID.Regnum
			email.Subject = ""
			MqNatEmail(email)
			sentCount += 1
		}
		if isEmail, ok := userConf["isEmail"]; ok && isEmail == "1" {
			if emailAddress, ok := userConf["emailAddress"]; ok || emailAddress != "" {
				orgInfo, orgType, isDefault := getOrgInfo(request.ID.ClientId, redis)
				html := generateXypHtml(request, orgInfo, orgType, isDefault)
				var email model.EmailModel
				email.Body = html
				email.CivilId = request.ID.CivilId
				email.Regnum = request.ID.Regnum
				email.Destination = emailAddress
				email.Subject = ""
				MqPrivEmail(email)
				sentCount += 1
			}
		}
		if isSocial, ok := userConf["social"]; ok && isSocial == "1" {
			var messengerRequest model.MessengerModel
			messengerRequest.Body = text
			messengerRequest.CivilId = request.ID.CivilId
			messengerRequest.Regnum = request.ID.Regnum
			MqMessenger(messengerRequest)
			sentCount += 1
	}else {
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
				request := generatePushRequest(text, "medeeleld haldsan baina", request.ID.CivilId, request.ID.Regnum, deviceTokens)
				MqPush(request)
				sentCount += 1
			}
		}
		if isNationalEmail, ok := userConf["isNationalEmail"]; ok && isNationalEmail == "1" {
			orgInfo, orgType, isDefault := getOrgInfo(request.ID.ClientId, redis)
			html := generateXypHtml(request, orgInfo, orgType, isDefault)
			var email model.EmailModel
			email.Body = html
			email.CivilId = request.ID.CivilId
			email.Regnum = request.ID.Regnum
			email.Subject = ""
			MqNatEmail(email)
			sentCount += 1
		}
	}
	return sentCount
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

func getOrgInfo(clientId int, redis *redis.Client) (model.OrgInfoModel, string, bool) {
	orgType := "GOVERNMENT_COMPANY"
	isDefault := false
	url := "https://10.10.18.28/api/service/notifications"
	var orgInfo model.OrgInfoModel

	if redis.Get("orgInfo:"+string(clientId)) != nil {
		log.Println("already in redis:" + string(clientId))
		redisString, err := redis.Get("orgInfo:" + string(clientId)).Result()
		if err != nil {
			//error log
			isDefault = true
		} else {
			err = json.Unmarshal([]byte(redisString), &orgInfo)
			if err != nil {
				//error log
				isDefault = true
			} else {
				if isOrgInfoModelEmpty(orgInfo) {
					isDefault = true
				} else {
					if orgInfo.Client_is_soap == "0" {
						orgType = "USER"
					} else if orgInfo.Organization.Org_type_id != "1" {
						orgType = "PRIVATE_COMPANY"
					}
				}
			}
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

			err = json.Unmarshal(body, &orgInfo)
			if err != nil {
				log.Fatal(err)
			} else {
				if isOrgInfoModelEmpty(orgInfo) {
					isDefault = true
				} else {
					if orgInfo.Success {
						temp, err := json.Marshal(orgInfo)
						if err != nil {
							// error log
							isDefault = true
						} else {
							_, err = redis.Set("orgInfo:"+string(clientId), string(temp), 0).Result()
							if err != nil {
								//error log
							}
							if orgInfo.Client_is_soap == "0" {
								orgType = "USER"
							} else if orgInfo.Organization.Org_type_id != "1" {
								orgType = "PRIVATE_COMPANY"
							}
						}

					}
				}
			}
		}
	}
	if clientId == 61 {
		orgType = "KIOSK"
	}
	return orgInfo, orgType, isDefault
}

func generateXypHtml(request GroupResult, orgInfo model.OrgInfoModel, orgType string, isDefault bool) string {
	var resulthtml string
	return resulthtml
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

func isNullOrEmpty(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.String() == ""
	case reflect.Slice, reflect.Map, reflect.Array:
		return value.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return value.IsNil()
	default:
		return false
	}
}

func isOrgInfoModelEmpty(model model.OrgInfoModel) bool {
	modelValue := reflect.ValueOf(model)
	modelType := modelValue.Type()

	for i := 0; i < modelValue.NumField(); i++ {
		field := modelValue.Field(i)
		fieldType := modelType.Field(i)

		if isNullOrEmpty(field) && fieldType.Tag.Get("json") != "-" {
			return true
			// Alternatively, you can return false if you want to check if any field is non-empty
		}
	}

	return false
}
