package controllers

import (
	"net/http"

	firebase "firebase.google.com/go"

	fire "github.com/ucladevx/bconnect-backend/storage/firebase"
)

func Setup(fb *firebase.App) {
	app := fire.NewFirebaseApp(
		fb,
	)
}

func GetUser(w http.ResponseWriter, r *http.Request) {

}

func CreateUser(w http.ResponseWriter, r *http.Request) {

}
