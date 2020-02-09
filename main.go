package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ucladevx/BConnect-backend/controllers"
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

	controllers.Setup(context.Background(), []byte(``))
}
