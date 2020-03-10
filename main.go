package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/ucladevx/BConnect-backend/middleware/bconnecthandlers"
	"github.com/ucladevx/BConnect-backend/server/usercontroller"
	"github.com/ucladevx/BConnect-backend/storage/postgres"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func startServerAndServices(config Config) {
	var db *gorm.DB
	var friendDB *gorm.DB
	_, ok := os.LookupEnv("DATABASE_URL")
	if ok {
		db = postgres.HerokuConnect("DATABASE_URL")
		friendDB = postgres.HerokuConnect("HEROKU_POSTGRESQL_GOLD_URL")
	}
	if !ok {
		db = postgres.Connect(config.Storage.UserHost,
			config.Storage.UserUsername,
			config.Storage.Username,
			config.Storage.UserPassword)

		friendDB = postgres.Connect(config.Storage.FriendHost,
			config.Storage.FriendUsername,
			config.Storage.Friendname,
			config.Storage.FriendPassword)
	}

	auth := postgres.NewPostgresClient(db)
	userActions := postgres.NewUserActions(db, friendDB)
	filters := postgres.NewFilterers(db, friendDB)

	postgres.CreatePostgresTables(auth, userActions)

	var userController = usercontroller.UserController{
		Service: auth,
		Actions: userActions,
		Filters: filters,
	}

	r := mux.NewRouter()
	http.Handle("/", r)
	userController.Setup(r)
	s := r.PathPrefix("/auth").Subrouter()
	s.Use(bconnecthandlers.VerifyToken)
	userController.AuthSetup(s)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	log.Printf("Listening on %s:%s", config.Server.Host, port)

	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Access-Control-Allow-Origin"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)))
}

func main() {
	conf := Conf()
	startServerAndServices(conf)
}
