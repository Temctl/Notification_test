module github.com/Temctl/E-Notification/util

go 1.20

require gopkg.in/natefinch/lumberjack.v2 v2.2.1

require github.com/go-redis/redis v6.15.9+incompatible // indirect

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/redis/go-redis/v9 v9.0.5
	github.com/streadway/amqp v1.0.0
)
