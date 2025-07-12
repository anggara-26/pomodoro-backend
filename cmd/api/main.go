package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/anggara-26/pomodoro-backend.git/app/model"
	"github.com/anggara-26/pomodoro-backend.git/pkg/middleware"
	"github.com/anggara-26/pomodoro-backend.git/pkg/router"
	"github.com/anggara-26/pomodoro-backend.git/platform/db"
	"github.com/gofiber/fiber/v2"
)

type Application struct {
	Models model.Models
}

// @title           Pomodoro API
// @version         1.0
// @description     This is the API for Pomodoro App
func main() {
	mongoClient, err := db.ConnectToMongo()
	if err != nil {
		log.Panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	log.Println("Connected to MongoDB!")

	// Create database indexes
	if err := db.CreateIndexes(mongoClient.Database(os.Getenv("MONGODB"))); err != nil {
		log.Printf("Warning: Failed to create indexes: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server is running on port " + port)

	config := middleware.CORSConfig()
	app := fiber.New(config)
	router.CreateRouter(app)
	log.Fatal(app.Listen(":" + port))
}
