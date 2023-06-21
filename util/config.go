package util

import (
	"fmt"
	"os"
	"time"
)

var (
	PRIVATE_EMAIL_NUM  = "private_email_number"
	NATIONAL_EMAIL_NUM = "national_email_number"
	PUSH_NOTIF_NUM     = "push_notif_number"
	SOCIAL_NOTIF_NUM   = "social_notif_number"

	PUSHNOTIFICATIONKEY = "PUSHNOTIFICATIONKEY"
	NATEMAILKEY         = "NATEMAILKEY"
	PRIVEMAILKEY        = "PRIVEMAILKEY"
	MESSENGERKEY        = "MESSENGERKEY"
	SMSKEY              = "SMSKEY"

	PUSHWORKER      = "PUSHWORKER"
	NATEMAILWORKER  = "NATEMAILWORKER"
	PRIVEMAILWORKER = "PRIVEMAILWORKER"
	MESSENGERWORKER = "MESSENGERWORKER"
	SMSWORKER       = "SMSWORKER"

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

	DB_HOST     string
	DB_PORT     string
	DB_USERNAME string
	DB_PASSWORD string
	DB_DBNAME   string

	MONGO_URL string

	ATTENTION_URL          string
	ATTENTION_SERVICENAME  string
	ATTENTION_SERVICENAME2 string
	OBJECTCODE             string
	ORGCODE                string
	ORGNAME                string
	ORGPASSWORD            string
	ORGTOKEN               string

	AWS_SES_USER     = "AKIA25YVKNUDIE6WUJAP"
	AWS_SES_PASSWORD = "BH/k4+kFhzOVLoKV8SG4bRPGrkMy1tM3cFsHunMAgFX2"
	AWS_SMTP         = "email-smtp.ap-southeast-1.amazonaws.com"
	FROM_EMAIL       = "notification@e-mongolia.mn"
)

func init() {
	if ENV == "" {
		PORT = 8085

		REDIS_HOST = "172.72.0.11"
		REDIS_PORT = 6379
		REDIS_PASSWORD = ""
		REDIS_DB = 0

		RABBITMQ_HOST = "localhost"
		RABBITMQ_PORT = 5672

		ELOG_MAXSIZE = 100
		ELOG_BACKUPS = 10
		ELOG_MAXAGE = 30

		RABBITMQURL = "amqp://guest:guest@172.72.0.11:5672/"

		DB_HOST = "172.72.0.11"
		DB_PORT = "5432"
		DB_USERNAME = "postgres"
		DB_PASSWORD = "changeme"
		DB_DBNAME = "postgres"

		MONGO_URL = "mongodb://admin:VRuAd2Nvmp4ELHh5@172.72.0.11:27017"

		ATTENTION_URL = "https://st-sso.e-mongolia.mn/xyp-api/api/xyp/get-data-no-auth"
		OBJECTCODE = "GET_IDCARD_DATE_OF_EXPIRY_LIST"
		ORGCODE = "10001001"
		ORGNAME = "E-mongolia"
		ORGPASSWORD = "aaa1"
		ORGTOKEN = "aaaaa"
		ATTENTION_SERVICENAME = "WS101137_citizenInfoLogByDate"
		ATTENTION_SERVICENAME2 = "WS100443_driverLicenseExpiredLog"

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

		DB_HOST = "172.72.0.11"
		DB_PORT = "5432"
		DB_USERNAME = "postgres"
		DB_PASSWORD = "changeme"
		DB_DBNAME = "postgres"

		ATTENTION_URL = "https://st-sso.e-mongolia.mn/xyp-api/api/xyp/get-data-no-auth"
		OBJECTCODE = "GET_IDCARD_DATE_OF_EXPIRY_LIST"
		ORGCODE = "10001001"
		ORGNAME = "E-mongolia"
		ORGPASSWORD = "aaa1"
		ORGTOKEN = "aaaaa"

		MONGO_URL = "mongodb://admin:VRuAd2Nvmp4ELHh5@172.72.0.11:27017"
	}
}

func GetTZ() (*time.Location, error) {
	location, err := time.LoadLocation("Asia/Ulaanbaatar")
	if err != nil {
		fmt.Println("Error loading time zone:", err)
	}
	return location, err
}

func WriteDbLog() bool {
	return true
}
