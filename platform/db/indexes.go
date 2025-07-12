package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateIndexes creates database indexes for better performance
func CreateIndexes(db *mongo.Database) error {
	ctx := context.Background()

	// Users collection indexes
	usersCollection := db.Collection("users")
	usersIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "firebase_uid", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "email", Value: 1}},
		},
	}

	_, err := usersCollection.Indexes().CreateMany(ctx, usersIndexes)
	if err != nil {
		log.Printf("Error creating users indexes: %v", err)
		return err
	}

	// Tasks collection indexes
	tasksCollection := db.Collection("tasks")
	tasksIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "status", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "assigned_at", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "created_at", Value: -1},
			},
		},
		{
			Keys: bson.D{{Key: "title", Value: "text"}},
		},
	}

	_, err = tasksCollection.Indexes().CreateMany(ctx, tasksIndexes)
	if err != nil {
		log.Printf("Error creating tasks indexes: %v", err)
		return err
	}

	// Sessions collection indexes
	sessionsCollection := db.Collection("sessions")
	sessionsIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "task_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "status", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "started_at", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "type", Value: 1},
				{Key: "status", Value: 1},
				{Key: "started_at", Value: -1},
			},
		},
	}

	_, err = sessionsCollection.Indexes().CreateMany(ctx, sessionsIndexes)
	if err != nil {
		log.Printf("Error creating sessions indexes: %v", err)
		return err
	}

	log.Println("Database indexes created successfully")
	return nil
}
