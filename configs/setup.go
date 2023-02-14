package configs

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func ConnectDB(config map[string]string) *mongo.Client {
	// we can use connection string => "mongodb://localhost:27017" => env.go -> EnvMongoURI method
	client, err := mongo.NewClient(options.Client().ApplyURI(config["MONGOURI"]))

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
