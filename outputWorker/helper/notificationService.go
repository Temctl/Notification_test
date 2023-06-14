package helper

import (
	"firebase.google.com/go/messaging"
	"github.com/Temctl/E-Notification/util/model"
	"github.com/go-redis/redis"
)

func SendRegularNotif(civilId string, content string, notifificationType model.NotificationType, redis *redis.Client, client *messaging.Client) {
	userConf, err := redis.HGetAll("conf:" + civilId).Result()
	if err != nil {
		panic(err)
	}
	if isPush, ok := userConf["isPush"]; ok && isPush == "true" {
		userDeviceTokens, err := redis.LRange("deviceTokens:"+civilId, 0, -1).Result()
		if err != nil {
			panic(err)
		} else {
			var push model.PushNotificationModel
			push.Body = content
			push.Title = "regular notif test"
			push.Type = notifificationType
			PushToTokens(push, userDeviceTokens, client)
		}
	}
	if isNationalEmail, ok := userConf["isNationalEmail"]; ok && isNationalEmail == "true" {

	}
	if isEmail, ok := userConf["isEmail"]; ok && isEmail == "true" {
		if emailAddress, ok := userConf["emailAddress"]; ok || emailAddress != "" {

		}
	}
	if isSocial, ok := userConf["social"]; ok && isSocial == "true" {
		SendSocial(civilId)
	}
}

func SendAttentionNotif(civilId string, content string, notifificationType model.NotificationType, redis *redis.Client, client *messaging.Client) {
	result := 0
	userConf, err := redis.HGetAll("conf:" + civilId).Result()
	if err != nil {
		panic(err)
	}
	if isPush, ok := userConf["isPush"]; ok && isPush == "true" {
		userDeviceTokens, err := redis.LRange("deviceTokens:"+civilId, 0, -1).Result()
		if err != nil {
			panic(err)
		} else {
			var push model.PushNotificationModel
			push.Body = content
			push.Title = "regular notif test"
			push.Type = notifificationType
			PushToTokens(push, userDeviceTokens, client)
		}
	}
	if isNationalEmail, ok := userConf["isNationalEmail"]; ok && isNationalEmail == "true" {
		result += AttentionNatEmail(civilId, content, notifificationType)
	}
	if isEmail, ok := userConf["isEmail"]; ok && isEmail == "true" {
		if emailAddress, ok := userConf["emailAddress"]; ok || emailAddress != "" {
			result += AttentionPrivEmail(emailAddress, content, notifificationType)
		}
	}
	if isSocial, ok := userConf["social"]; ok && isSocial == "true" {
		result += SendSocial(civilId)
	}
}
