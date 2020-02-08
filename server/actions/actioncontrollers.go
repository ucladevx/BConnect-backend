package actions

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ucladevx/BConnect-backend/models"
)

// AuthService abstract server-side authentication in case we switch from whatever current auth scheme we are using
type AuthService interface {
	GET(username string, password string) (*models.User, bool)
	SET(username string, password string) (bool, error)
	PUT(username string, password string) (bool, error)
	DEL(username string, password string) (bool, error)
}

// UserController abstract server-side authentication
type UserController struct {
	Service AuthService
}

// Body credentials necessary for username/password auth
type Body struct {
	UUID     string `json:"uuid"`
	Password string `json:"password"`
}

// Auth credentials necessary for username/password auth
type Auth struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

// CurrUser refers to current user and corresponding JWT token
type CurrUser struct {
	Username string `json:"Username"`
	Token    string `json:"Token"`
}

// Setup sets up handlers
func (auth *UserController) Setup(r *mux.Router) {
	r.HandleFunc("/login", auth.Login).Methods("POST", "GET", "OPTIONS")
	r.HandleFunc("/signup", auth.Signup).Methods("GET", "OPTIONS")
	r.HandleFunc("/delete", auth.Delete).Methods("GET", "OPTIONS")
}

// Login login users and provides auth token
func (auth *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var userInfo Body

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&userInfo)

	auth.Service.GET(userInfo.UUID, userInfo.Password)
}

// Signup signs up users and provides auth token
func (auth *UserController) Signup(w http.ResponseWriter, r *http.Request) {
	var userInfo Body

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&userInfo)

	auth.Service.PUT(userInfo.UUID, userInfo.Password)
}

// Set changes users in DB
func (auth *UserController) Set(w http.ResponseWriter, r *http.Request) {
	var userInfo Body

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&userInfo)

	auth.Service.SET(userInfo.UUID, userInfo.Password)
}

// Delete deletes users and removes them from DB
func (auth *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	var userInfo Body

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&userInfo)

	auth.Service.DEL(userInfo.UUID, userInfo.Password)
}
