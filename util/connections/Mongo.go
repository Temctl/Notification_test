package connections

import (
	"context"
	"fmt"

	"github.com/Temctl/E-Notification/util"
	"github.com/Temctl/E-Notification/util/elog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB() (*mongo.Client, error) {
	// Set connection options
	clientOptions := options.Client().ApplyURI(util.MONGO_URL)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		elog.Error().Println(err)
		return nil, err
	}

	// Ping MongoDB to check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		elog.Error().Println(err)
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")
	return client, nil
}
