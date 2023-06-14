package connections

import (
	"context"
	"strconv"

	"github.com/Temctl/E-Notification/util"
	"github.com/Temctl/E-Notification/util/elog"
	"github.com/redis/go-redis/v9"
)

func ConnectionRedis() *redis.Client {

	// -------------------------------------------------------
	// GET UTIL CONFIG ---------------------------------------
	// -------------------------------------------------------

	host := util.REDIS_HOST
	port := strconv.Itoa(util.REDIS_PORT)

	// -------------------------------------------------------
	// CREATE REDIS CLIENT -----------------------------------
	// -------------------------------------------------------

	addr := host + ":" + port
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// -------------------------------------------------------
	// Ping the Redis server to check the connection ---------
	// -------------------------------------------------------

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		elog.Info().Println("Failed to ping Redis server:", err)
		return nil
	}
	elog.Info().Println("Redis server response: " + pong)

	return client
}
