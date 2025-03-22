package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/anggara-26/pomodoro-backend.git/app/model"
	"github.com/anggara-26/pomodoro-backend.git/db"
	"github.com/anggara-26/pomodoro-backend.git/pkg/router"
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

	log.Println("Server is running on port " + os.Getenv("PORT"))

	app := fiber.New()
	router.CreateRouter(app)
	app.Listen(":" + os.Getenv("PORT"))
}
