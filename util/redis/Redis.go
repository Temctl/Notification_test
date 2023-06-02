package redis

import (
	"context"
	"strconv"

	"github.com/Temctl/E-Notification/util"
	"github.com/Temctl/E-Notification/util/elog"
	"github.com/redis/go-redis/v9"
)

// -------------------------------------------------------
// SET REDIS ---------------------------------------------
// -------------------------------------------------------
func SetRedis(key string, data string) bool {

	elog.Info("set redis...")
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
	// REDIS luu data oruulah --------------------------------
	// -------------------------------------------------------

	ctx := context.Background()
	clientErr := client.Set(ctx, key, data, 0).Err()
	if clientErr != nil {
		elog.Error("redis setlehed aldaa garlaa", clientErr)
		return false
	} else {
		elog.Info("Successful...")
	}

	// -------------------------------------------------------
	// Close the Redis client --------------------------------
	// -------------------------------------------------------

	closeErr := client.Close()
	if closeErr != nil {
		elog.Error("Error closing Redis client:", closeErr)
		return false
	} else {
		elog.Info("Redis client closed successfully")
	}
	return true
}

// -------------------------------------------------------
// GET REDIS ---------------------------------------------
// -------------------------------------------------------
func GetRedis(key string) string {
	elog.Info("get redis...")
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
	// REDIS ees data avah -----------------------------------
	// -------------------------------------------------------

	ctx := context.Background()
	val, err := client.Get(ctx, key).Result()
	if err != nil {
		elog.Error("redis client get data: ", err)
		return ""
	}
	return val
}
