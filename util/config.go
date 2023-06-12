package util

import "os"

var (
	PRIVATE_EMAIL_NUM  = "private_email_number"
	NATIONAL_EMAIL_NUM = "national_email_number"
	PUSH_NOTIF_NUM     = "push_notif_number"
	SOCIAL_NOTIF_NUM   = "social_notif_number"

	REGULARNOTIFKEY   = "REGULARNOTIFKEY"
	ATTENTIONNOTIFKEY = "ATTENTIONNOTIFKEY"
	XYPNOTIFKEY       = "XYPNOTIFKEY"
	GROUPNOTIFKEY     = "GROUPNOTIFKEY"

	ENV  string = os.Getenv("ENV")
	PORT int

	REDIS_HOST     string
	REDIS_PORT     int
	REDIS_PASSWORD string
	REDIS_DB       int

	RABBITMQ_HOST string
	RABBITMQ_PORT int

	ELOG_MAXSIZE int
	ELOG_BACKUPS int
	ELOG_MAXAGE  int

	RABBITMQURL string

	AWS_SES_USER     = "AKIA25YVKNUDIE6WUJAP"
	AWS_SES_PASSWORD = "BH/k4+kFhzOVLoKV8SG4bRPGrkMy1tM3cFsHunMAgFX2"
	AWS_SMTP         = "email-smtp.ap-southeast-1.amazonaws.com"
	FROM_EMAIL       = "Notification system <notification@e-mongolia.mn>"
)

func init() {
	if ENV == "" {
		PORT = 8085

		REDIS_HOST = "localhost"
		REDIS_PORT = 6379
		REDIS_PASSWORD = ""
		REDIS_DB = 0

		RABBITMQ_HOST = "localhost"
		RABBITMQ_PORT = 5672

		ELOG_MAXSIZE = 100
		ELOG_BACKUPS = 10
		ELOG_MAXAGE = 30

		RABBITMQURL = "amqp://guest:guest@localhost:5672/"
	} else if ENV == "prod" {
		PORT = 8085

		REDIS_HOST = "localhost"
		REDIS_PORT = 6379
		REDIS_PASSWORD = ""
		REDIS_DB = 0

		RABBITMQ_HOST = "localhost"
		RABBITMQ_PORT = 5672

		ELOG_MAXSIZE = 100
		ELOG_BACKUPS = 10
		ELOG_MAXAGE = 30

		RABBITMQURL = "amqp://guest:guest@localhost:5672/"
	}
}

func WriteDbLog() bool {
	return true
}
