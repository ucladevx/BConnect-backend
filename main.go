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
	friendStore := postgres.NewFriendStorage(db)
	filters := postgres.NewFilterers(db)
	emailStore := postgres.NewEmailStorage(db)

	postgres.CreatePostgresTables(userStore, friendStore, emailStore)
	userService := users.NewUserService(userStore, friendStore, emailStore)

	userController := controllers.NewUserController(userService, filters)

	r := mux.NewRouter()
	CORSMWR := mux.CORSMethodMiddleware(r)
	r.Use(CORSMWR)

	http.Handle("/", r)
	userController.Setup(r)

	s := r.PathPrefix("/auth").Subrouter()
	CORSMWS := mux.CORSMethodMiddleware(s)
	s.Use(CORSMWS)

	s.Use(middleware.VerifyToken)
	userController.AuthSetup(s)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	log.Printf("Listening on %s:%s", config.Server.Host, port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Access-Control-Allow-Origin"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedHeaders([]string{"x-access-token"}), handlers.AllowedOrigins([]string{"*"}))(r)))
}

func main() {
	conf := Conf()
	startServerAndServices(conf)

}
