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
	DB *mongo.Client
)

func Connect() func() {
	// Set client options
	clientOptions := options.Client().ApplyURI(fmt.Sprintf(
		"mongodb+srv://social-sum:%s@socialsumcluster.eprbbju.mongodb.net/?retryWrites=true&w=majority",
		os.Getenv("MONGO_PASSWORD"),
	))

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Connection to MongoDB successful")

	DB = client

	return func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
		slog.Info("Connection to MongoDB closed")
	}
}
