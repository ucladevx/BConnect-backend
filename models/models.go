package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
)

// Location type definition
type Location int

//FilterReturn return type of filter
type FilterReturn struct {
	Filter *gorm.DB
}

// TODO - implement for privacy-forward implementation of this

// UserLocation attempts to create a less obvious mapping of users to locations for privacy
type UserLocation struct {
	UUID         string   `json:"uuid"`
	UserLocation Location `json:"location"`
}

// User user struct
type User struct {
	gorm.Model
	FirstName   string  `json:"fname"`
	LastName    string  `json:"lname"`
	Email       string  `json:"username"`
	Password    string  `json:"password"`
	PhoneNumber string  `json:"phonenumber"`
	ProfilePic  string  `json:"profilepic"`
	UserID      string  `json:"userid"`
	Major       string  `json:"degree"`
	GradYear    string  `json:"year"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Interests   string  `json:"interests"`
	Clubs       string  `json:"clubs"`
	Bio         string  `json:"bio"`
}

// Friends structure of friend
type Friends struct {
	gorm.Model
	UserID       string `json:"uuid"`
	FriendID     string `json:"fuuid"`
	FReqMess     string `json:"msg"`
	TimeStamp    int64  `json:"timestamp"`
	Status       int    `json:"status"`
	TimeAccepted int64  `json:"timeaccepted"`
}

// FriendRequest friend request
type FriendRequest struct {
	gorm.Model
	Sender    string `json:"uuid"`
	Receiver  string `json:"fuuid"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

// Token tokens
type Token struct {
	jwt.Claims
	UUID           string
	Email          string
	StandardClaims *jwt.StandardClaims
}

// RefreshToken refresh tokens
type RefreshToken struct {
	jwt.Claims
	UUID           string
	ID             string
	StandardClaims *jwt.StandardClaims
}

// Chats chats
type Chats struct {
}

//ChatRooms user specific chat rooms
type ChatRooms struct {
}

// Message message object
type Message struct {
	Sender   string `json:"sender"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

// Interests interests
type Interests struct {
	UserID   string `json:"userid"`
	Interest string `json:"interest"`
}

// Filterer filters
type Filterer func(*FilterReturn, []string) *FilterReturn

// Finder finds
type Finder func(map[string]Filterer, map[string][]string) map[string]interface{}
