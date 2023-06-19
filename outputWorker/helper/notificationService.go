package helper

// import (
// 	"encoding/json"
// 	"fmt"

// 	"firebase.google.com/go/messaging"
// 	"github.com/Temctl/E-Notification/util/model"
// 	"github.com/go-redis/redis"
// )

// func SendRegularNotif(civilId string, regnum string, content string, notifificationType model.NotificationType, redis *redis.Client, client *messaging.Client) {
// 	if civilId == "" {
// 		exists, err := redis.Exists("getByReg:" + regnum).Result()
// 		if err != nil {
// 			panic(err)
// 		} else if exists == 1 {
// 			civilId, err = redis.Get("getByReg:" + regnum).Result() // if civil id is not sent, get it using regnum from redis conf
// 			if err != nil {
// 				panic(err)
// 			}
// 		}
// 	} else {
// 		// TODO get regnum from redis
// 	}
// 	result := 0
// 	fmt.Printf("CivilId: %s\n", civilId)
// 	fmt.Printf("Content: %s\n", content)
// 	exists, err := redis.Exists("conf:" + civilId).Result()
// 	if err != nil {
// 		panic(err)
// 	}
// 	if exists == 1 {
// 		userConf, err := redis.HGetAll("conf:" + civilId).Result()
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Println(userConf)
// 		if isPush, ok := userConf["isPush"]; ok && isPush == "1" {
// 			deviceTokensExists, err := redis.Exists("deviceTokens:" + civilId).Result()
// 			if err != nil {
// 				panic(err)
// 			} else if deviceTokensExists == 1 { //if it exists
// 				data, err := redis.Get("deviceTokens:" + civilId).Result() //get device token from redis
// 				fmt.Println(data)
// 				if err != nil {
// 					panic(err)
// 				} else {
// 					var deviceTokens []string
// 					err := json.Unmarshal([]byte(data), &deviceTokens) //redis is returning string, so turn the string array into slice
// 					if err != nil {
// 						panic(err)
// 					}
// 					request := generatePushRequest("regular notif test", content, notifificationType)
// 					PushToTokens(request, deviceTokens, client)
// 				}
// 			} // TODO when is push is open but no device token found
// 		}
// 		if isNationalEmail, ok := userConf["isNationalEmail"]; ok && isNationalEmail == "1" {
// 			result += RegularNatEmail(civilId, content)
// 		}
// 		if isEmail, ok := userConf["isEmail"]; ok && isEmail == "1" {
// 			if emailAddress, ok := userConf["emailAddress"]; ok || emailAddress != "" {
// 				result += RegularPrivEmail(emailAddress, content)
// 			}
// 		}
// 		if isSocial, ok := userConf["social"]; ok && isSocial == "1" {
// 			SendSocial(civilId)
// 		}
// 	} else {
// 		deviceTokensExists, err := redis.Exists("deviceTokens:" + civilId).Result()
// 		if err != nil {
// 			panic(err)
// 		} else if deviceTokensExists == 1 {
// 			userDeviceTokens, err := redis.LRange("deviceTokens:"+civilId, 0, -1).Result()
// 			if err != nil {
// 				panic(err)
// 			} else {
// 				request := generatePushRequest("regular notif test", content, notifificationType)
// 				PushToTokens(request, userDeviceTokens, client)
// 			}
// 		}
// 		result += RegularNatEmail(civilId, content)
// 	}

// }

// func SendAttentionNotif(civilId string, regnum string, content string, notifificationType model.NotificationType, redis *redis.Client, client *messaging.Client) {
// 	if civilId == "" {
// 		exists, err := redis.Exists("getByReg:" + regnum).Result()
// 		if err != nil {
// 			panic(err)
// 		} else if exists == 1 {
// 			civilId, err = redis.Get("getByReg:" + regnum).Result() // if civil id is not sent, get it using regnum from redis conf
// 			if err != nil {
// 				panic(err)
// 			}
// 		}
// 	} else {
// 		// TODO get regnum from redis
// 	}
// 	result := 0
// 	exists, err := redis.Exists("conf:" + civilId).Result()
// 	if err != nil {
// 		panic(err)
// 	} else if exists == 1 {
// 		userConf, err := redis.HGetAll("conf:" + civilId).Result()
// 		if err != nil {
// 			panic(err)
// 		}
// 		if isPush, ok := userConf["isPush"]; ok && isPush == "1" {
// 			deviceTokensExists, err := redis.Exists("deviceTokens:" + civilId).Result()
// 			if err != nil {
// 				panic(err)
// 			} else if deviceTokensExists == 1 { //if it exists
// 				data, err := redis.Get("deviceTokens:" + civilId).Result() //get device token from redis
// 				fmt.Println(data)
// 				if err != nil {
// 					panic(err)
// 				} else {
// 					var deviceTokens []string
// 					err := json.Unmarshal([]byte(data), &deviceTokens) //redis is returning string, so turn the string array into slice
// 					if err != nil {
// 						panic(err)
// 					}
// 					request := generatePushRequest("regular notif test", content, notifificationType)
// 					PushToTokens(request, deviceTokens, client)
// 				}
// 			} // TODO when is push is open but no device token found
// 		}
// 		if isNationalEmail, ok := userConf["isNationalEmail"]; ok && isNationalEmail == "1" {
// 			result += AttentionNatEmail(civilId, content, notifificationType)
// 		}
// 		if isEmail, ok := userConf["isEmail"]; ok && isEmail == "1" {
// 			if emailAddress, ok := userConf["emailAddress"]; ok || emailAddress != "" {
// 				result += AttentionPrivEmail(emailAddress, content, notifificationType)
// 			}
// 		}
// 		if isSocial, ok := userConf["social"]; ok && isSocial == "1" {
// 			result += SendSocial(civilId)
// 		}
// 	} else {
// 		deviceTokensExists, err := redis.Exists("deviceTokens:" + civilId).Result()
// 		if err != nil {
// 			panic(err)
// 		} else if deviceTokensExists == 1 {
// 			userDeviceTokens, err := redis.LRange("deviceTokens:"+civilId, 0, -1).Result()
// 			if err != nil {
// 				panic(err)
// 			} else {
// 				request := generatePushRequest("regular notif test", content, notifificationType)
// 				PushToTokens(request, userDeviceTokens, client)
// 			}
// 		}
// 		result += AttentionNatEmail(civilId, content, notifificationType)
// 	}

// }

// func SendXypNotif(request model.XypNotification, notifificationType model.NotificationType, redis *redis.Client, client *messaging.Client) {
// 	// var civilId string
// 	// var regnum string
// 	// if request.CivilId == "" {
// 	// 	regnum = request.Regnum
// 	// 	exists, err := redis.Exists("getByReg:" + regnum).Result()
// 	// 	if err != nil {
// 	// 		panic(err)
// 	// 	} else if exists == 1 {
// 	// 		civilId, err = redis.Get("getByReg:" + regnum).Result() // if civil id is not sent, get it using regnum from redis conf
// 	// 		if err != nil {
// 	// 			panic(err)
// 	// 		}
// 	// 	}
// 	// } else {
// 	// 	// TODO get regnum from redis
// 	// }

// 	// result := 0
// 	// fmt.Printf("CivilId: %s\n", civilId)
// 	// exists, err := redis.Exists("conf:" + civilId).Result()
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	// if exists == 1 {
// 	// 	userConf, err := redis.HGetAll("conf:" + civilId).Result()
// 	// 	if err != nil {
// 	// 		panic(err)
// 	// 	}
// 	// 	fmt.Println(userConf)
// 	// 	if isPush, ok := userConf["isPush"]; ok && isPush == "1" {
// 	// 		deviceTokensExists, err := redis.Exists("deviceTokens:" + civilId).Result()
// 	// 		if err != nil {
// 	// 			panic(err)
// 	// 		} else if deviceTokensExists == 1 { //if it exists
// 	// 			data, err := redis.Get("deviceTokens:" + civilId).Result() //get device token from redis
// 	// 			fmt.Println(data)
// 	// 			if err != nil {
// 	// 				panic(err)
// 	// 			} else {
// 	// 				var deviceTokens []string
// 	// 				err := json.Unmarshal([]byte(data), &deviceTokens) //redis is returning string, so turn the string array into slice
// 	// 				if err != nil {
// 	// 					panic(err)
// 	// 				}
// 	// 				request := generatePushRequest("regular notif test", "medeeleld haldsan baina", notifificationType)
// 	// 				PushToTokens(request, deviceTokens, client)
// 	// 			}
// 	// 		} // TODO when is push is open but no device token found
// 	// 	}
// 	// 	if isNationalEmail, ok := userConf["isNationalEmail"]; ok && isNationalEmail == "1" {
// 	// 		result += XypNatEmail(civilId, content)
// 	// 	}
// 	// 	if isEmail, ok := userConf["isEmail"]; ok && isEmail == "1" {
// 	// 		if emailAddress, ok := userConf["emailAddress"]; ok || emailAddress != "" {
// 	// 			result += RegularPrivEmail(civilId, content)
// 	// 		}
// 	// 	}
// 	// 	if isSocial, ok := userConf["social"]; ok && isSocial == "1" {
// 	// 		SendSocial(civilId)
// 	// 	}
// 	// } else {
// 	// 	deviceTokensExists, err := redis.Exists("deviceTokens:" + civilId).Result()
// 	// 	if err != nil {
// 	// 		panic(err)
// 	// 	} else if deviceTokensExists == 1 {
// 	// 		userDeviceTokens, err := redis.LRange("deviceTokens:"+civilId, 0, -1).Result()
// 	// 		if err != nil {
// 	// 			panic(err)
// 	// 		} else {
// 	// 			request := generatePushRequest("regular notif test", "medeeleld haldsan baina", notifificationType)
// 	// 			PushToTokens(request, userDeviceTokens, client)
// 	// 		}
// 	// 	}
// 	// 	result += RegularNatEmail(civilId, content)
// 	// }

// }

// func generatePushRequest(title string, body string, notificationType model.NotificationType) model.PushNotificationModel {
// 	var request model.PushNotificationModel
// 	request.Body = body
// 	request.Title = title
// 	request.Type = notificationType
// 	return request
// }
