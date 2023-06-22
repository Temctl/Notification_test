package sender

import (
	"context"
	"log"
	"sync"

	"github.com/Temctl/E-Notification/util/connections"
	"github.com/Temctl/E-Notification/util/model"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func XypFromDb(mongoClient *mongo.Client, redis *redis.Client) {
	var collectionName model.Collections
	collectionName = model.XYPNOTIFICATION
	collection, client, err := connections.ConnectMongoDB(collectionName)
	if err != nil {
		//error log
	}
	defer client.Disconnect(context.Background())

	findOptions := options.Find().SetLimit(500)

	cursor, err := collection.Find(context.Background(), nil, findOptions)
	if err != nil {
		//log error
	}

	// Iterate over the returned documents
	defer cursor.Close(context.TODO())
	var wg sync.WaitGroup
	for cursor.Next(context.TODO()) {
		var notif model.XypNotification
		if err := cursor.Decode(&notif); err != nil {
			log.Fatal(err)
		} else {
			wg.Add(1)
			go func(notif model.XypNotification) {
				defer wg.Done()
				if notif.CivilId == "" && notif.Regnum == "" {

				} else {
					if notif.CivilId == "" {
						exists, err := redis.Exists("getByReg:" + notif.Regnum).Result()
						if err != nil {
							panic(err)
						} else if exists == 1 {
							notif.CivilId, err = redis.Get("getByReg:" + notif.Regnum).Result() // if civil id is not sent, get it using regnum from redis conf
							if err != nil {
								panic(err)
							}
						}
					}
					sendMq(notif)
				}
			}(notif)
		}
	}

	// Wait for all goroutines to complete
	wg.Wait()
}

func sendMq(request model.XypNotification) {

}
