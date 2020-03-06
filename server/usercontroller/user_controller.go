package usercontroller

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/ucladevx/BConnect-backend/models"
)

// Claims used to determine auth token passed in header
type Claims struct {
	jwt.MapClaims
	UUID           string
	Email          string
	StandardClaims *jwt.StandardClaims
}

// AuthService abstract server-side authentication in case we switch from whatever current auth scheme we are using
type AuthService interface {
	GET(username string, password string) (map[string]interface{}, string, string, time.Time, time.Time)
	SET(email string, uuid string, userModded *models.User) (map[string]interface{}, error)
	PUT(email string, password string, firstName string, lastName string) (bool, error)
	DEL(username string, password string) (bool, error)
	REFRESH()
}

// UserService abstract user-side functionality in case we switch from whatever current db scheme we are using
type UserService interface {
	ADD(currUUID string, friendUUID string, optionalMsg string) (*models.FriendRequest, error)
	ACCEPT(currUUID string, friendUUID string) (*models.FriendRequest, error)
	GET(currUUID string) map[string]interface{}
	LEAVE(currUUID string)
}

// UserController abstract server-side authentication
type UserController struct {
	Service AuthService
	Actions UserService
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

// NewUser creds necessary to create new user
type NewUser struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	FName    string `json:"FName"`
	LName    string `json:"LName"`
}

// CurrUser refers to current user and corresponding JWT token
type CurrUser struct {
	Username string `json:"Username"`
	Token    string `json:"Token"`
}

// Setup sets up handlers
func (auth *UserController) Setup(r *mux.Router) {
	r.HandleFunc("/login", auth.Login).Methods("POST")
	r.HandleFunc("/signup", auth.Signup).Methods("POST")
}

// AuthSetup sets up auth handlers
func (auth *UserController) AuthSetup(r *mux.Router) {
	r.HandleFunc("/change", auth.Set).Methods("POST")
	r.HandleFunc("/delete", auth.Delete).Methods("GET")
	r.HandleFunc("/addfriend", auth.AddFriend).Methods("GET")
	r.HandleFunc("/acceptfriend", auth.AcceptFriend).Methods("GET")
	r.HandleFunc("/getfriends", auth.GetFriend).Methods("GET")
	r.HandleFunc("/refresh", auth.Refresh).Methods("GET")
}

// Login login users and provides authentication token for user
func (auth *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var userInfo Auth
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&userInfo)
	resp, token, refreshToken, expirationTime, refreshExpirationTime := auth.Service.GET(userInfo.Username, userInfo.Password)
	if token != "" {
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   token,
			Expires: expirationTime,
		})
		http.SetCookie(w, &http.Cookie{
			Name:    "refresh",
			Value:   refreshToken,
			Expires: refreshExpirationTime,
		})
		json.NewEncoder(w).Encode(resp)
	}
}

// Signup signs up users and provides auth token
func (auth *UserController) Signup(w http.ResponseWriter, r *http.Request) {
	var userInfo NewUser
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&userInfo)

	auth.Service.PUT(userInfo.Username, userInfo.Password, userInfo.FName, userInfo.LName)
	resp, token, refreshToken, expirationTime, refreshExpirationTime := auth.Service.GET(userInfo.Username, userInfo.Password)
	if token != "" {
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   token,
			Expires: expirationTime,
		})
		http.SetCookie(w, &http.Cookie{
			Name:    "refresh",
			Value:   refreshToken,
			Expires: refreshExpirationTime,
		})
		json.NewEncoder(w).Encode(resp)
	}
}

// Set changes users in DB
func (auth *UserController) Set(w http.ResponseWriter, r *http.Request) {
	var userInfo models.User

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userInfo)
	if err != nil {
		print(err.Error)
	}
	claims := auth.getCurrentUserFromTokenProvided(w, r)

	resp, err := auth.Service.SET(claims.Email, claims.UUID, &userInfo)
	if err != nil {
		print(err.Error)
	}
	json.NewEncoder(w).Encode(resp)
}

// AddFriend generates friendRequest to specified UUID
func (auth *UserController) AddFriend(w http.ResponseWriter, r *http.Request) {
	claims := auth.getCurrentUserFromTokenProvided(w, r)
	auth.Actions.ADD(claims.UUID, r.URL.Query().Get("friend_uuid"), r.URL.Query().Get("message"))
}

// AcceptFriend accepts friendRequest from specified UUID
func (auth *UserController) AcceptFriend(w http.ResponseWriter, r *http.Request) {
	claims := auth.getCurrentUserFromTokenProvided(w, r)
	auth.Actions.ACCEPT(claims.UUID, r.URL.Query().Get("friend_uuid"))
}

// GetFriend gets a list of user friends
func (auth *UserController) GetFriend(w http.ResponseWriter, r *http.Request) {
	claims := auth.getCurrentUserFromTokenProvided(w, r)

	resp := auth.Actions.GET(claims.UUID)
	json.NewEncoder(w).Encode(resp)
}

// Delete deletes users and removes them from DB
func (auth *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	var userInfo Body

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&userInfo)

	auth.Service.DEL(userInfo.UUID, userInfo.Password)
}

// Logout logs out users
func (auth *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	claims := auth.getCurrentUserFromTokenProvided(w, r)
	auth.Actions.LEAVE(claims.UUID)
}

// Refresh generates a new authentication token for the current user and sends it
func (auth *UserController) Refresh(w http.ResponseWriter, r *http.Request) {

}

func (auth *UserController) getCurrentUserFromTokenProvided(w http.ResponseWriter, r *http.Request) Claims {
	header := strings.TrimSpace(r.Header.Get("x-access-token"))

	claims := Claims{}

	header = strings.Replace(header, "Bearer ", "", -1)
	_, err := jwt.ParseWithClaims(header, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		panic(err)
	}
	return claims
}