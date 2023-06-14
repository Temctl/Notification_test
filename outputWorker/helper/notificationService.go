package helper

import (
	"firebase.google.com/go/messaging"
	"github.com/Temctl/E-Notification/util/model"
	"github.com/go-redis/redis"
)

func SendRegularNotif(civilId string, content string, notifificationType model.NotificationType, redis *redis.Client, client *messaging.Client) {
	userConf, err := redis.HGetAll("conf:" + civilId).Result()
	var push1 model.PushNotificationModel
	push1.Body = "regular notif test"
	push1.Title = "regular notif test"

	if err != nil {
		panic(err)
	}
	if isPush, ok := userConf["isPush"]; ok && isPush == "true" {
		userDeviceTokens, err := redis.LRange("deviceTokens:"+civilId, 0, -1).Result()
		if err != nil {
			panic(err)
		} else {
			PushToTokens(push1, userDeviceTokens, client)
		}
	}
	if isNationalEmail, ok := userConf["isNationalEmail"]; ok && isNationalEmail == "true" {
		SendNatEmail(civilId)
	}
	if isEmail, ok := userConf["isEmail"]; ok && isEmail == "true" {
		if emailAddress, ok := userConf["emailAddress"]; ok || emailAddress != "" {
			SendPrivEmail(emailAddress)
		}
	}
	if isSocial, ok := userConf["social"]; ok && isSocial == "true" {
		SendSocial(civilId)
	}
}
