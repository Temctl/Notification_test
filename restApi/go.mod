module github.com/Temctl/E-Notification/restApi

go 1.20

require (
	github.com/Temctl/E-Notification/inputWorker v0.0.0-20230616075337-2612745fb402
	github.com/Temctl/E-Notification/util v0.0.0-20230616075337-2612745fb402
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/streadway/amqp v1.0.0
)

require (
	github.com/go-redis/redis v6.15.9+incompatible // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
