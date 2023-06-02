package util

var (
	PORT int = 8085

	REDIS_HOST     string = "localhost"
	REDIS_PORT     int    = 6379
	REDIS_PASSWORD string = ""
	REDIS_DB       int    = 0

	RABBITMQ_HOST string = "localhost"
	RABBITMQ_PORT int    = 5672

	ELOG_PATH    string = "./tmp"
	ELOG_MAXSIZE int    = 100
	ELOG_BACKUPS int    = 10
	ELOG_MAXAGE  int    = 30
)
