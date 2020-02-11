package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	firebase "firebase.google.com/go"
	"github.com/ucladevx/BConnect-backend/controllers"
	"google.golang.org/api/option"
)

func main() {
	type Configuration struct {
		project_id           string
		private_key_id       string
		private_key          string
		client_email         string
		client_id            string
		auth_uri             string
		token_uri            string
		auth_provider        string
		client_x509_cert_url string
	}

	file, _ := os.Open("credentials/connect_b-sdk.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}

	config := &firebase.Config{
		DatabaseURL: "https://connect-b.firebaseio.com/",
	}

	opt := option.WithCredentialsFile("./credentials/connect_b-sdk.json")
	fb, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		print(err.Error())
	}
	controllers.Setup(context.Background(), fb)
}
