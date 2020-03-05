package postgres

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/ucladevx/BConnect-backend/models"
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
	/* TODO: Maybe change migration model to maybe define DB relationships */
	client.client.AutoMigrate(&models.User{})
}

// GET gets user for login
func (client *Client) GET(key string, password string) (map[string]interface{}, string, string, time.Time, time.Time) {
	resp, token, refreshToken, expirationTime, refreshTokenExpirationTime := client.findOne(key, password)
	return resp, token, refreshToken, expirationTime, refreshTokenExpirationTime
}

func (client *Client) findOne(email, password string) (map[string]interface{}, string, string, time.Time, time.Time) {
	user := &models.User{}

	if err := client.client.Where("Email = ?", email).First(user).Error; err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Email address not found"}
		return resp, "", "", time.Time{}, time.Time{}
	}

	expiresAt := time.Now().Add(time.Minute * 15)

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword {
		var resp = map[string]interface{}{"status": false, "message": "Invalid login credentials"}
		return resp, "", "", time.Time{}, time.Time{}
	}

	tk := &models.Token{
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

	refreshString, err := generateRandomString(32)
	if err != nil {

	}

	refreshExpiresAt := time.Now().Add(time.Hour * 144)

	refreshTk := &models.RefreshToken{
		UUID: user.UUID,
		ID:   refreshString,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: refreshExpiresAt.Unix(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), refreshTk)

	refreshTokenString, error := refreshToken.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}

	var resp = map[string]interface{}{"status": false, "message": "logged in"}
	resp["token"] = tokenString
	resp["refresh"] = refreshTokenString
	resp["user"] = user

	return resp, tokenString, refreshTokenString, expiresAt, refreshExpiresAt
}

func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// PUT puts user into postgres
func (client *Client) PUT(email string, password string, firstname string, lastname string) (bool, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	user := &models.User{}
	user.Password = string(pass)
	user.Email = email
	user.FirstName = firstname
	user.LastName = lastname
	user.UUID = uuid.UUID()
	client.client.Create(user)
	return true, nil
}

// SET sets updated fields
func (client *Client) SET(email string, password string, firstName string, lastName string) (bool, error) {
	var user models.User
	if err := client.client.Where("EMAIL = ? AND PASSWORD = ?", email, password).Find(&user); err != nil {
	}
	user.FirstName = firstName
	user.LastName = lastName
	client.client.Save(&user)

	return true, nil
}

// DEL dels clients
func (client *Client) DEL(key string, password string) (bool, error) {
	return true, nil
}

//REFRESH generates a new refresh token for the authenticated client
func (client *Client) REFRESH() {

}
