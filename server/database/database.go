package database

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

)

var (
	Client *mongo.Client
)

func Connect() func() {
	// Set client options
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://jesuloba:%s@socialsumcluster.eprbbju.mongodb.net/?retryWrites=true&w=majority", os.Getenv("MONGO_PASSWORD")))

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	slog.Info("Connection to MongoDB successful")

	if err != nil {
		log.Fatal(err)
	}

	Client = client

	return func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
		slog.Info("Connection to MongoDB closed")
	}
}
