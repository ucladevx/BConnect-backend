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

	controllers.Setup(context.Background(), []byte(`{
		"type": "service_account",
		"project_id": "connect-b",
		"private_key_id": "1f86e5f09bb86c86d41eaac7ae97a333e714d7a8",
		"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCvLHKTDyYRL/HJ\nK6oDxiZo6TIQOx7IYOxqS5GUILHpInef0volDuPp26b5gvrAm7W7nGcNDWjowEOb\ndnHaz7ANWS7xU34nBsXQdqOEVtqrj0NaS7eichVUBDG5chWn3Z+0HPlLiWggfXc1\nRNQcfyoVPql4HVRiVzUcmQkV0BcvshbQ+DiC2VXFZMUTBBEJMT0HW/SRHD9I0+JX\nilBDj+nHONOUKlin15jkayFsmA8Zyo4a6VXJ02juyxdy71Ptef7TAijUNUMtKPwh\nWxvAQh4ecgLG5uGEf5E7SUIj4AJnOzANPO9zL/XqDkfJkF/3m5cW5/AZc1HsvmML\n8y6Pg8e7AgMBAAECggEAHdOsopumuLB4M+/KYCAiNCTY0Giwoh/WwXaikB1NDcw/\noCgVTbAOYKh08vE5bwhUVmCYL2HEJjVujY8Kbd6FNJCl7JNx3IVLs2YwC32aeDHJ\nnxtbZj6UoRyhttjVFSMoaUuxqwn3f+i9HoXctAl3Ce0EgB5GWwRMV9U69crb19Pr\nGzlT/eONQbOuF6TWP9VpHfv25M1Kei6yR9063KSx7CJ/ZCREUZqDSwzi54mPMJdL\npLFOIglCoRFAYJ1cZ7s3t9DzcCJluaddGOyhLbqrsrPySarsKeZEWAPG4jc99hpM\nms3EGQ+UEkqh5aQfpX2zP7SmewNWwEd0G8ijIVf1GQKBgQD2XStj1KjbHphcWoxW\n1i1rZy/1xEI9MU41A6mfdqxk873Yl27fJKKNYseXiZeH/obR2LuvIpZuF9fvEjfA\n9ZY2+d2ky0aXFV9hl7jxMY0m6dpQZ74XgtzlIy/7U0fJ/aNuDG1pbxl/eyYztsmt\nrgZBgtGBgTGG8xkpj/j74GwdlwKBgQC2BnPn8iylBJsuoE/zK23gdaMbwdvKimJn\n6vqqkDJawYfiW1OQUaiCjtNq886o9tKtGD49zfQsIfJ1k+ecwxaEIA+t7GW3ihQB\n51qt6PvSK347egqgGNeLw3014mq/OsOVA2T8t6HlT4enDQR4ZJO8AFN/PBtNlndb\nOakOpNXzfQKBgQCkTx1KNqHriis5pRZmL/AY6rgc5Kj84O51Ax1VAjQWHE3or469\nelSDkXbmhViv1bybJj5+nsXPZU2Z3/+ZTPHDdsAxUXWh/BoiH6u0CUVHx73X1Gj+\ni0PB+sbciv2dJPGjytwJ7pQF5t9irC00DZWUiagrBDxA8c89Xg1EB7hzvwKBgAiM\niQkNwLcHXlp0QQ0EryZpn/1/v1jl4vVKffdgylyk1kL1UxlmHGn1V0ygosYgwYhL\nqTCx6ZPhDEglaC1epEIUV2gtwCE8pO/p88JTPPCEBmu4saMPR6BS0CClv6m3ktP+\n0tFjtoDUmWRpIhqdbqrXwRQquIWOWZC5Ro1fhITpAoGBAMEcPaU//7NpR1VAOZk+\n8asBWDAJgeeb/sMhtMKDn/kHU0R0OB6OKH5UYUOgH1B/e6RNzA3K1mXDdbFsJq9B\npAf6VuEqZERj1Nt02YUBemJa4LKMJeGNlb93w/iN4jh5uUYQx+Ep6aQcj3ktg1uZ\nSMYyrpU7P1dGoHU3BktTo7t0\n-----END PRIVATE KEY-----\n",
		"client_email": "firebase-adminsdk-e01ut@connect-b.iam.gserviceaccount.com",
		"client_id": "111920755554416487917",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-e01ut%40connect-b.iam.gserviceaccount.com"
	  }`))
}
