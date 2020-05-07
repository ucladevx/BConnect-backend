package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/ucladevx/BConnect-backend/controllers"
	"github.com/ucladevx/BConnect-backend/middleware"
	"github.com/ucladevx/BConnect-backend/services/users"
	"github.com/ucladevx/BConnect-backend/storage/postgres"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func startServerAndServices(config Config) {
	var db *gorm.DB
	_, ok := os.LookupEnv("DATABASE_URL")
	if ok {
		db = postgres.HerokuConnect("DATABASE_URL")
	}
	if !ok {
		db = postgres.Connect(config.Storage.UserHost,
			config.Storage.UserUsername,
			config.Storage.Username,
			config.Storage.UserPassword)
	}

	userStore := postgres.NewUserStorage(db)
	filters := postgres.NewFilterers(db)

	postgres.CreatePostgresTables(userStore)

	userService := users.NewUserService(userStore)
	userController := controllers.NewUserController(userService, filters)

	r := mux.NewRouter()
	http.Handle("/", r)
	userController.Setup(r)
	s := r.PathPrefix("/auth").Subrouter()
	s.Use(middleware.VerifyToken)
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
