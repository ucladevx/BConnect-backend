package controllers

import (
	"context"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/ucladevx/BConnect-backend/server/actions"
	fire "github.com/ucladevx/BConnect-backend/storage/firebase"
	"google.golang.org/api/option"
)

// Setup inits server for setup
func Setup(ctx context.Context, service []byte) {
	config := &firebase.Config{
		DatabaseURL: "https://connect-b.firebaseio.com/",
	}

	opt := option.WithCredentialsJSON(service)
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

	if app == nil {
		return
	}

	var userController = actions.UserController{
		Service: app,
	}

	r := mux.NewRouter()
	http.Handle("/", r)
	userController.Setup(r)

	log.Printf("Listening on %s%s", os.Getenv("HOST"), os.Getenv("PORT"))

	log.Fatal(http.ListenAndServe(":3000", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)))
}
