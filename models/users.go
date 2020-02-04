package models

import (
	"time"
)

// Location type definition
type Location int

// User structure of user
type User struct {
	FirstName    string    `json:"firstname"`
	LastName     string    `json:"lastname"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	PhoneNumber  string    `json:"phonenumber"`
	ProfilePic   string    `json:"profilepic"`
	CreatedAt    time.Time `json:"uploadtime"`
	UUID         string    `json:"uuid"`
	UserLocation Location  `json:"location"`
	Friends      []Friends `json:"friends"`
}
