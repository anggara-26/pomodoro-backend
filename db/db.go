package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func GetDBCollection(col string) *mongo.Collection {
	return db.Collection(col)
}

func ConnectToMongo() (*mongo.Client, error) {
	// getting username, password, and host from .env
	username := os.Getenv("MONGODB_USERNAME")
	password := os.Getenv("MONGODB_PASSWORD")
	host := os.Getenv("MONGODB_HOST")
	dbName := os.Getenv("MONGODB")

	// MongoDb connection string
	clientOptions := options.Client().ApplyURI("mongodb+srv://" + username + ":" + password + "@" + host + "/" + dbName)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	db = client.Database(dbName)

	return client, nil
}
