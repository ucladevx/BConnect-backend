package postgres

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

// Location type definition
type Location int

// User user struct
type User struct {
	gorm.Model
	FirstName    string   `json:"firstname"`
	LastName     string   `json:"lastname"`
	Email        string   `json:"email"`
	Password     string   `json:"password"`
	PhoneNumber  string   `json:"phonenumber"`
	ProfilePic   string   `json:"profilepic"`
	UUID         string   `json:"uuid"`
	UserLocation Location `json:"location"`
}

// Friends structure of friend
type Friends struct {
	gorm.Model
	UUID  string `json:"uuid"`
	FUUID string `json:"fuuid"`
}

// Token tokens
type Token struct {
	jwt.Claims
	UUID           string
	FirstName      string
	Email          string
	StandardClaims *jwt.StandardClaims
}
