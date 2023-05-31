package model

type PushNotificationModel struct {
	Title    string
	Body     string
	ImageUrl string
	Data     map[string]string
}
