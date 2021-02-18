package fbauth

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	"os"
)

var (
	firebaseConfigFile = os.Getenv("FIREBASE_CONFIG_FILE")
)

func InitAuth() (*auth.Client, error) {
	opt := option.WithCredentialsFile(firebaseConfigFile)
	config := &firebase.Config{ProjectID: "tassk-2bff1"}

	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return nil, errors.Wrap(err, "error initializing firebase auth (create firebase app)")
	}

	client, errAuth := app.Auth(context.Background())
	if errAuth != nil {
		return nil, errors.Wrap(err, "error initializing firebase auth (creating client)")
	}

	return client, nil
}