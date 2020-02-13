package postgres

import (
	"fmt"
	"time"

	"github.com/ucladevx/BConnect-backend/utils/uuid"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Client Postgresql client
type Client struct {
	client *gorm.DB
}

// NewPostgresClient Postgresql client
func NewPostgresClient(client *gorm.DB) *Client {
	return &Client{
		client: client,
	}
}

func (client *Client) create() {
	client.client.AutoMigrate(&User{})
}

// GET gets user for login
func (client *Client) GET(key string, password string) (map[string]interface{}, string, time.Time) {
	resp, token, expirationTime := client.findOne(key, password)
	return resp, token, expirationTime
}

func (client *Client) findOne(email, password string) (map[string]interface{}, string, time.Time) {
	user := &User{}

	if err := client.client.Where("Email = ?", email).First(user).Error; err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Email address not found"}
		return resp, "", time.Time{}
	}

	expiresAt := time.Now().Add(time.Minute * 100000)

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		var resp = map[string]interface{}{"status": false, "message": "Invalid login credentials. Please try again"}
		return resp, "", time.Time{}
	}

	tk := &Token{
		UUID:      user.UUID,
		FirstName: user.FirstName,
		Email:     user.Email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}

	var resp = map[string]interface{}{"status": false, "message": "logged in"}
	resp["token"] = tokenString
	resp["user"] = user

	return resp, tokenString, expiresAt
}

// PUT puts user into postgres
func (client *Client) PUT(email string, password string, firstname string, lastname string) (bool, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	user := &User{}
	user.Password = string(pass)
	user.Email = email
	user.FirstName = firstname
	user.LastName = lastname
	user.UUID = uuid.UUID()
	client.client.Create(user)
	return true, nil
}

// SET sets updated fields
func (client *Client) SET(key string, password string) (bool, error) {
	return true, nil
}

// DEL dels clients
func (client *Client) DEL(key string, password string) (bool, error) {
	return true, nil
}
