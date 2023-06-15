package helper

import (
	"firebase.google.com/go/messaging"
	"github.com/Temctl/E-Notification/util/model"
	"github.com/go-redis/redis"
)

func SendRegularNotif(civilId string, content string, notifificationType model.NotificationType, redis *redis.Client, client *messaging.Client) {
	result := 0
	exists, err := redis.Exists("conf:" + civilId).Result()
	if err != nil {
		panic(err)
	}
	if exists == 1 {
		userConf, err := redis.HGetAll("conf:" + civilId).Result()
		if err != nil {
			panic(err)
		}
		if isPush, ok := userConf["isPush"]; ok && isPush == "true" {
			deviceTokensExists, err := redis.Exists("deviceTokens:" + civilId).Result()
			if err != nil {
				panic(err)
			} else if deviceTokensExists == 1 {
				userDeviceTokens, err := redis.LRange("deviceTokens:"+civilId, 0, -1).Result()
				if err != nil {
					panic(err)
				} else {
					request := generatePushRequest("regular notif test", content, notifificationType)
					PushToTokens(request, userDeviceTokens, client)
				}
			} // TODO when is push is open but no device token found
		}
		if isNationalEmail, ok := userConf["isNationalEmail"]; ok && isNationalEmail == "true" {
			result += RegularNatEmail(civilId, content)
		}
		if isEmail, ok := userConf["isEmail"]; ok && isEmail == "true" {
			if emailAddress, ok := userConf["emailAddress"]; ok || emailAddress != "" {
				result += RegularPrivEmail(civilId, content)
			}
		}
		if isSocial, ok := userConf["social"]; ok && isSocial == "true" {
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
				request := generatePushRequest("regular notif test", content, notifificationType)
				PushToTokens(request, userDeviceTokens, client)
			}
		}
		result += RegularNatEmail(civilId, content)
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
			request := generatePushRequest("attention notif test", content, notifificationType)
			PushToTokens(request, userDeviceTokens, client)
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

func generatePushRequest(title string, body string, notificationType model.NotificationType) model.PushNotificationModel {
	var request model.PushNotificationModel
	request.Body = body
	request.Title = title
	request.Type = notificationType
	return request
}
