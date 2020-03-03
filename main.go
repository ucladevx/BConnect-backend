package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ucladevx/BConnect-backend/bconnecthandlers"
	"github.com/ucladevx/BConnect-backend/server/userauth"
	"github.com/ucladevx/BConnect-backend/storage/postgres"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func startServerAndServices(config Config) {
	db := postgres.Connect(config.Storage.UserHost,
		config.Storage.UserUsername,
		config.Storage.Username,
		config.Storage.UserPassword)

	friendDB := postgres.Connect(config.Storage.FriendHost,
		config.Storage.FriendUsername,
		config.Storage.Friendname,
		config.Storage.FriendPassword)

	auth := postgres.NewPostgresClient(db)
	userActions := postgres.NewUserActions(db, friendDB)

	postgres.CreatePostgresTables(auth, userActions)

	var userController = userauth.UserController{
		Service: auth,
		Actions: userActions,
	}

	r := mux.NewRouter()
	http.Handle("/", r)
	userController.Setup(r)
	s := r.PathPrefix("/auth").Subrouter()
	s.Use(bconnecthandlers.VerifyToken)
	userController.AuthSetup(s)

	log.Printf("Listening on %s%s", config.Server.Host, config.Server.Port)

	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)))
}

func main() {
	conf := Conf()
	startServerAndServices(conf)
}
