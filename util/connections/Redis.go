package connections

import (
	"strconv"

	"github.com/Temctl/E-Notification/util"
	"github.com/Temctl/E-Notification/util/elog"
	"github.com/go-redis/redis"
)

func ConnectionRedis() (*redis.Client, error) {

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

	ping, err := client.Ping().Result()
	if err != nil {
		elog.Info().Println("Failed to ping Redis server:", err)
	}
	elog.Info().Println("Redis server response: " + ping)

	return client, err
}
