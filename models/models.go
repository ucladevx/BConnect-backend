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
	FirstName   string      `json:"fname"`
	LastName    string      `json:"lname"`
	Email       string      `json:"username"`
	Password    string      `json:"password"`
	PhoneNumber string      `json:"phonenumber"`
	ProfilePic  string      `json:"profilepic"`
	UserID      string      `json:"userid"`
	Major       string      `json:"degree"`
	GradYear    string      `json:"year"`
	Lat         float64     `json:"lat"'`
	Lon         float64     `json:"lon"`
	Bio         string      `json:"bio"`
	Clubs       []*Club     `gorm:"many2many:user_clubs;" json:"clubs"`
	Interests   []*Interest `gorm:"many2many:user_interests;" json:"interests"`
	Friends     []*User     `gorm:"many2many:user_friends;association_jointable_foreignkey:friend_id"`
}

// FriendRequest friend request
type FriendRequest struct {
	gorm.Model
	Sender string `json:"uuid"`
	Receiver string `json:"fuuid"`
	Message  string `json:"message"`
}

// Interest struct
type Interest struct {
	Interest string  `gorm:"primary_key" json:"interest"`
	Users    []*User `gorm:"many2many:user_interests;"`
}

type Club struct {
	Club  string  `gorm:"primary_key"`
	Users []*User `gorm:"many2many:user_clubs;"`
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

// Filterer filters
type Filterer func(*gorm.DB, []string) *gorm.DB

// Finder finds
type Finder func(map[string]Filterer, map[string][]string) map[string]interface{}
