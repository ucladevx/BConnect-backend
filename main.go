package main

import (
	"log"
	"net/http"

	"github.com/ucladevx/BConnect-backend/bconnecthandlers"
	"github.com/ucladevx/BConnect-backend/server/userauth"
	"github.com/ucladevx/BConnect-backend/storage/postgres"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func startServerAndServices(config Config) {
	db := postgres.Connect(config.Storage.Host,
		config.Storage.Username,
		config.Storage.Name,
		config.Storage.Password)

	friendDB := postgres.Connect(config.Storage.Host,
		config.Storage.Username,
		config.Storage.Friends,
		config.Storage.Password)

	app := postgres.NewPostgresClient(db)
	userActions := postgres.NewUserActions(db, friendDB)

	postgres.CreatePostgresTables(app, userActions)

	var userController = userauth.UserController{
		Service: app,
		Actions: userActions,
	}

	r := mux.NewRouter()
	http.Handle("/", r)
	userController.Setup(r)
	s := r.PathPrefix("/auth").Subrouter()
	s.Use(bconnecthandlers.VerifyToken)
	userController.AuthSetup(s)

	log.Printf("Listening on %s%s", config.Server.Host, config.Server.Port)

	log.Fatal(http.ListenAndServe(config.Server.Port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)))
}

func main() {
	conf := Conf()
	startServerAndServices(conf)
}
