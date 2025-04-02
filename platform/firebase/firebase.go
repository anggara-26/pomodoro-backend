package firebase

import (
	"context"
	"log"
	"sync"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var (
	authClient *auth.Client
	once       sync.Once
)

func InitFirebase() *auth.Client {
	once.Do(func() {
		opt := option.WithCredentialsFile("serviceAccountKey.json")
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			log.Fatalf("Error initializing Firebase App: %v", err)
		}

		authClient, err = app.Auth(context.Background())
		if err != nil {
			log.Fatalf("Error initializing Firebase Auth Client: %v", err)
		}
	})

	return authClient
}
