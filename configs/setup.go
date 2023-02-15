package configs

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func ConnectDB(URI string) *mongo.Client {
	// we can use directly connection string => "mongodb://localhost:27017"
	client, err := mongo.NewClient(options.Client().ApplyURI(URI))

	if err != nil {
		log.Fatalln(err)
	}
	// If don't connect within 20 seconds, give us an error
	var ctx, _ = context.WithTimeout(context.Background(), 20*time.Second)
	err = client.Connect(ctx)

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalln(err)
	}

	return client
}
