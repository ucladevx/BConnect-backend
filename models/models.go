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

// ChangeData ChangeData struct
type ChangeData struct {
	gorm.Model
	FirstName   string   `json:"firstname"`
	LastName    string   `json:"lastname"`
	Age         int      `json:"age"`
	GradYear    string   `json:"gradyear"`
	CurrentJob  string   `json:"currentjob"`
	Gender      string   `json:"gender"`
	Email       string   `json:"username"`
	Major       string   `json:"major"`
	Password    string   `json:"password"`
	PhoneNumber string   `json:"phonenumber"`
	ProfilePic  string   `json:"profilepic"`
	Lat         float64  `json:"lat"`
	Lon         float64  `json:"lon"`
	Bio         string   `json:"bio"`
	Interests   []string `json:"interests"`
}

//SplitIntoUserAndInterests splits
func (cd *ChangeData) SplitIntoUserAndInterests() (*User, []string) {
	user := &User{
		FirstName:   cd.FirstName,
		LastName:    cd.LastName,
		Age:         cd.Age,
		GradYear:    cd.GradYear,
		CurrentJob:  cd.CurrentJob,
		Gender:      cd.Gender,
		Email:       cd.Email,
		Major:       cd.Major,
		Password:    cd.Password,
		PhoneNumber: cd.PhoneNumber,
		ProfilePic:  cd.ProfilePic,
		Lat:         cd.Lat,
		Lon:         cd.Lon,
		Bio:         cd.Bio,
	}

	return user, cd.Interests
}

// User user struct
type User struct {
	gorm.Model
	UserID      string  `json:"userid" gorm:"unique;not null"`
	FirstName   string  `json:"firstname" gorm:"not null"`
	LastName    string  `json:"lastname" gorm:"not null"`
	Age         int     `json:"age"`
	GradYear    string  `json:"gradyear"`
	CurrentJob  string  `json:"currentjob"`
	Gender      string  `json:"gender"`
	Email       string  `json:"username" gorm:"unique; not null"`
	Major       string  `json:"major"`
	Password    string  `json:"password"`
	PhoneNumber string  `json:"phonenumber"`
	ProfilePic  string  `json:"profilepic"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
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

//ChatRoom chat room
type ChatRoom struct {
	UserID     string `json:"userid"`
	ChatRoomID string `json:"chatid"`
}

// Message message object
type Message struct {
	Message     []byte `json:"message"`
	MessageRoom string
	MessageID   string
}

// Interests interests
type Interests struct {
	UserID   string `json:"userid"`
	Interest string `json:"interests"`
}

// Filterer filters
type Filterer func(*FilterReturn, *User, []string) *FilterReturn

// Finder finds
type Finder func(map[string]Filterer, *User, map[string][]string) []User

// Email is just an email
type Email struct {
	gorm.Model
	Email string `json:"email"`
}
