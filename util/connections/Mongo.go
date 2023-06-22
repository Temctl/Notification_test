package connections

import (
	"context"

	"github.com/Temctl/E-Notification/util"
	"github.com/Temctl/E-Notification/util/elog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB(collectionName string) (*mongo.Collection, *mongo.Client, error) {
	// Set connection options
	clientOptions := options.Client().ApplyURI(util.MONGO_URL)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		elog.Error().Panic(err)
		return nil, nil, err
	}

	// Ping MongoDB to check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		elog.Error().Panic(err)
		return nil, nil, err
	}

	collection := client.Database("notification").Collection(collectionName)
	return collection, client, nil
}
