package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/anggara-26/pomodoro-backend.git/db"
	"github.com/anggara-26/pomodoro-backend.git/handlers"
	"github.com/anggara-26/pomodoro-backend.git/services"
)

type Application struct {
	Models services.Models
}

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

	services.New(mongoClient)

	log.Println("Connected to MongoDB!")
	// running on port
	log.Println("Server is running on port " + os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), handlers.CreateRouter()))
}
