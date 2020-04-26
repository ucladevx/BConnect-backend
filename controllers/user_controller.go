package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/ucladevx/BConnect-backend/models"
)

// Claims used to determine auth token passed in header
type Claims struct {
	jwt.MapClaims
	UUID           string
	Email          string
	StandardClaims *jwt.StandardClaims
}

// RefreshClaims used to determine auth token passed in header
type RefreshClaims struct {
	jwt.MapClaims
	UUID           string
	ID             string
	StandardClaims *jwt.StandardClaims
}

// UserService abstract user-side functionality in case we switch from whatever current db scheme we are using
type UserService interface {
	Login(username string, password string) (map[string]interface{}, string, string, time.Time, time.Time)
	Update(user *models.User) (map[string]interface{}, error)
	Signup(email string, password string, firstName string, lastName string) (bool, error)
	DeleteUser(username string, password string) (bool, error)
	RefreshToken(uuid string) (map[string]interface{}, string, time.Time)
	FriendRequest(currUUID string, friendUUID string, optionalMsg string) (*models.Friends, error)
	AcceptFriendRequest(currUUID string, friendUUID string) (*models.Friends, error)
	GetFriends(currUUID string) map[string]interface{}
	Leave(currUUID string)
	Filter(finder models.Finder, filters map[string]models.Filterer, args map[string][]string) map[string]interface{}
}

// Filterers abstracts filters
type Filterers interface {
	NameFilter(curr *gorm.DB, names []string) *gorm.DB
	MajorFilter(curr *gorm.DB, majors []string) *gorm.DB
	GradYearFilter(curr *gorm.DB, gradYear []string) *gorm.DB
	InterestsFilter(curr *gorm.DB, interests []string) *gorm.DB
	LocationRadiusFilter(curr *gorm.DB, radius []string) *gorm.DB
	FinalFilter(filters map[string]models.Filterer, args map[string][]string) map[string]interface{}
}

// UserController abstract server-side authentication
type UserController struct {
	UserService UserService
	Filters     Filterers
}

//NewUserController uc
func NewUserController(userService UserService, filters Filterers) *UserController {
	return &UserController{
		UserService: userService,
		Filters:     filters,
	}
}

// Setup sets up handlers
func (uc *UserController) Setup(r *mux.Router) {
	r.HandleFunc("/login", uc.Login).Methods("POST")
	r.HandleFunc("/signup", uc.Signup).Methods("POST")
	r.HandleFunc("/refresh", uc.Refresh).Methods("GET")
}

// AuthSetup sets up auth handlers
func (uc *UserController) AuthSetup(r *mux.Router) {
	r.HandleFunc("/change", uc.Update).Methods("POST")
	r.HandleFunc("/delete", uc.DeleteUser).Methods("GET")
	r.HandleFunc("/addfriend", uc.AddFriend).Methods("GET")
	r.HandleFunc("/acceptfriend", uc.AcceptFriend).Methods("GET")
	r.HandleFunc("/getfriends", uc.GetFriend).Methods("GET")
	r.HandleFunc("/filter/{filterOne}", uc.Filter).Methods("GET")
}

// Login login users and provides authentication token for user
func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var userInfo models.User
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&userInfo)
	resp, token, _, _, _ := uc.UserService.Login(userInfo.Email, userInfo.Password)
	if token != "" {
		json.NewEncoder(w).Encode(resp)
	}
}

// Signup signs up users and provides auth token
func (uc *UserController) Signup(w http.ResponseWriter, r *http.Request) {
	var userInfo models.User
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&userInfo)
	var resp = map[string]interface{}{"status": false, "user": userInfo}

	status, _ := uc.UserService.Signup(userInfo.Email, userInfo.Password, userInfo.FirstName, userInfo.LastName)
	if status != true {
		http.Error(w, "Error signing up", 500)
		return
	}
	print(resp["user"])
	resp, token, _, _, _ := uc.UserService.Login(userInfo.Email, userInfo.Password)
	if token != "" {
		json.NewEncoder(w).Encode(resp)
	}
}

// Update changes users in DB
func (uc *UserController) Update(w http.ResponseWriter, r *http.Request) {
	var userInfo models.User

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userInfo)
	if err != nil {
		print(err.Error)
	}
	resp, err := uc.UserService.Update(&userInfo)
	if err != nil {
		print(err.Error)
	}
	json.NewEncoder(w).Encode(resp)
}

// AddFriend generates friendRequest to specified UUID
func (uc *UserController) AddFriend(w http.ResponseWriter, r *http.Request) {
	claims := uc.getCurrentUserFromTokenProvided(w, r)
	friend, err := uc.UserService.FriendRequest(claims.UUID, r.URL.Query().Get("friend_uuid"), r.URL.Query().Get("message"))
	if err != nil {
		print("Hey")
	}
	var resp = map[string]interface{}{"added": true, "friend": friend}
	json.NewEncoder(w).Encode(resp)
}

// AcceptFriend accepts friendRequest from specified UUID
func (uc *UserController) AcceptFriend(w http.ResponseWriter, r *http.Request) {
	claims := uc.getCurrentUserFromTokenProvided(w, r)
	friend, err := uc.UserService.AcceptFriendRequest(claims.UUID, r.URL.Query().Get("friend_uuid"))
	if err != nil {

	}
	var resp = map[string]interface{}{"added": true, "friend": friend}
	json.NewEncoder(w).Encode(resp)
}

// GetFriend gets a list of user friends
func (uc *UserController) GetFriend(w http.ResponseWriter, r *http.Request) {
	claims := uc.getCurrentUserFromTokenProvided(w, r)

	resp := uc.UserService.GetFriends(claims.UUID)
	json.NewEncoder(w).Encode(resp)
}

//DeleteUser deletes users and removes them from DB
func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var userInfo models.User

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&userInfo)

	uc.UserService.DeleteUser(userInfo.UserID, userInfo.Password)
}

// Logout logs out users
func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	claims := uc.getCurrentUserFromTokenProvided(w, r)
	uc.UserService.Leave(claims.UUID)
}

// Refresh generates a new authentication token for the current user and sends it
func (uc *UserController) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshClaim := uc.getUUIDFromRefreshToken(w, r)
	resp, token, _ := uc.UserService.RefreshToken(refreshClaim.UUID)
	if token != "" {
		json.NewEncoder(w).Encode(resp)
	}
}

//Filter filters
func (uc *UserController) Filter(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	filterOne := params["filterOne"]
	if filterOne == "" {

	}
	funcMapper := map[string]models.Filterer{
		"name":      uc.Filters.NameFilter,
		"major":     uc.Filters.MajorFilter,
		"gradyear":  uc.Filters.GradYearFilter,
		"interests": uc.Filters.InterestsFilter,
		"radius":    uc.Filters.LocationRadiusFilter,
	}
	categories := map[string][]string{
		"name":      strings.Split(r.URL.Query().Get("names"), ","),
		"major":     strings.Split(r.URL.Query().Get("majors"), ","),
		"gradyear":  strings.Split(r.URL.Query().Get("gradyears"), ","),
		"interests": strings.Split(r.URL.Query().Get("interests"), ","),
		"radius":    strings.Split(r.URL.Query().Get("radius"), ","),
	}

	var resp = uc.UserService.Filter(uc.Filters.FinalFilter, funcMapper, categories)
	json.NewEncoder(w).Encode(resp)
}

func (uc *UserController) getCurrentUserFromTokenProvided(w http.ResponseWriter, r *http.Request) Claims {
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

func (uc *UserController) getUUIDFromRefreshToken(w http.ResponseWriter, r *http.Request) RefreshClaims {
	header := strings.TrimSpace(r.Header.Get("x-access-token"))

	refreshClaims := RefreshClaims{}

	header = strings.Replace(header, "Bearer ", "", -1)
	_, err := jwt.ParseWithClaims(header, &refreshClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		panic(err)
	}
	return refreshClaims
}