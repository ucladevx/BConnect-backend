package controllers

import (
	"context"
	"net/http"

	firebase "firebase.google.com/go"
	fire "github.com/ucladevx/bconnect-backend/storage/firebase"
	"google.golang.org/api/option"
)

func Setup(ctx context.Context) {
	config := &firebase.Config{
		DatabaseURL: "https://connect-b.firebaseio.com/",
	}
	opt := option.WithCredentialsFile("../storage/firebase/credentials/connect_b-sdk.json")
	fb, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {

	}
	client, err := fb.Database(ctx)
	if err != nil {

	}
	app := fire.NewFirebaseApp(
		fb,
		client,
	)
}

func GetUser(w http.ResponseWriter, r *http.Request) {

}

func CreateUser(w http.ResponseWriter, r *http.Request) {

}
