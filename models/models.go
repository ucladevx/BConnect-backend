package models

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
	UUID         string `json:"uuid"`
	FUUID        string `json:"fuuid"`
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
	FirstName      string
	Email          string
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
}
