package controllers

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

// RefreshClaims used to determine auth token passed in header
type RefreshClaims struct {
	jwt.MapClaims
	UUID           string
	ID             string
	StandardClaims *jwt.StandardClaims
}

//InterestsForm interests
type InterestsForm struct {
	Interests string `json:"interests"`
}

// UserService abstract user-side functionality in case we switch from whatever current db scheme we are using
type UserService interface {
	//first set of parentheses is the input, second set of parens is the outputs
	Login(username string, password string) (map[string]interface{}, string, error)
	Update(user *models.User, interestsString []string) (map[string]interface{}, error)
	Signup(user *models.User) (bool, error)
	DeleteUser(username string, password string) (bool, error)
	RefreshToken(uuid string) (map[string]interface{}, string, time.Time)
	FriendRequest(currUUID string, friendUUID string, optionalMsg string) (*models.Friends, error)
	AcceptFriendRequest(currUUID string, friendUUID string) (*models.Friends, error)
	GetFriends(currUUID string) []models.Friends
	Leave(currUUID string)
	Filter(finder models.Finder, currentUser *models.User, filters map[string]models.Filterer, args map[string][]string) []models.User
	CurrentUser(uuid string) *models.User
	AddEmail(email string) (bool, error)
}

// Filterers abstracts filters
type Filterers interface {
	NameFilter(curr *models.FilterReturn, currentUser *models.User, names []string) *models.FilterReturn
	MajorFilter(curr *models.FilterReturn, currentUser *models.User, majors []string) *models.FilterReturn
	GradYearFilter(curr *models.FilterReturn, currentUser *models.User, gradYear []string) *models.FilterReturn
	InterestsFilter(curr *models.FilterReturn, currentUser *models.User, interests []string) *models.FilterReturn
	LocationRadiusFilter(curr *models.FilterReturn, currentUser *models.User, radius []string) *models.FilterReturn
	FinalFilter(filters map[string]models.Filterer, currentUser *models.User, args map[string][]string) []models.User
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
	r.HandleFunc("/email", uc.AddEmail).Methods("POST")
}

// AuthSetup sets up auth handlers
func (uc *UserController) AuthSetup(r *mux.Router) {
	r.HandleFunc("/change", uc.Update).Methods("POST")
	r.HandleFunc("/delete", uc.DeleteUser).Methods("GET")
	r.HandleFunc("/addfriend", uc.AddFriend).Methods("GET")
	r.HandleFunc("/acceptfriend", uc.AcceptFriend).Methods("GET")
	r.HandleFunc("/getfriends", uc.GetFriend).Methods("GET")
	r.HandleFunc("/filter", uc.Filter).Methods("GET")
}

// Login login users and provides authentication token for user
func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var userInfo models.User
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&userInfo)
	resp, token, err := uc.UserService.Login(userInfo.Email, userInfo.Password)
	if err != nil {

	}
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

	status, _ := uc.UserService.Signup(&userInfo)
	if status != true {
		http.Error(w, "Error signing up", 500)
		return
	}
	resp, token, err := uc.UserService.Login(userInfo.Email, userInfo.Password)
	if err != nil {

	}
	if token != "" {
		json.NewEncoder(w).Encode(resp)
	}
}

// Update changes users in DB
func (uc *UserController) Update(w http.ResponseWriter, r *http.Request) {
	var formData models.ChangeData

	claim := uc.getCurrentUserFromTokenProvided(w, r)

	decoder := json.NewDecoder(r.Body)

	decoder.Decode(&formData)
	userInfo, interests := formData.SplitIntoUserAndInterests()

	userInfo.UserID = claim.UUID
	userInfo.Email = claim.Email

	resp, err := uc.UserService.Update(userInfo, interests)
	if err != nil {
		print(err.Error())
	}
	json.NewEncoder(w).Encode(resp)
}

// AddFriend generates friendRequest to specified UUID
func (uc *UserController) AddFriend(w http.ResponseWriter, r *http.Request) {
	claims := uc.getCurrentUserFromTokenProvided(w, r)
	friend, err := uc.UserService.FriendRequest(claims.UUID, r.URL.Query().Get("friend_uuid"), r.URL.Query().Get("message"))
	if err != nil {

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

	friends := uc.UserService.GetFriends(claims.UUID)
	var resp = map[string]interface{}{"num_friends": len(friends), "friends": friends}
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
	claim := uc.getCurrentUserFromTokenProvided(w, r)
	currentUser := uc.UserService.CurrentUser(claim.UUID)

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

	users := uc.UserService.Filter(uc.Filters.FinalFilter, currentUser, funcMapper, categories)
	var resp = map[string]interface{}{"users": users}
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

//AddEmail adds email
func (uc *UserController) AddEmail(w http.ResponseWriter, r *http.Request) {
	var emailInfo models.Email
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&emailInfo)

	status, _ := uc.UserService.AddEmail(emailInfo.Email)
	if status != true {
		http.Error(w, "Error adding email", 500)
		return
	}

	var resp = map[string]interface{}{"status": false, "message": "Email registered!", "email": emailInfo}
	json.NewEncoder(w).Encode(resp)
}

/*func (uc *UserController) Email(w http.ResponseWriter, r *http.Request) {
	var userInfo models.User
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&userInfo)
	var resp = map[string]interface{}{"status": false, "user": userInfo}

	status, _ := uc.UserService.Email(userInfo.Email)
	if status != true {
		http.Error(w, "Error signing up", 500)
		return
	}
	print(resp["user"])
	resp, token, _, _, _ := uc.UserService.Login(userInfo.Email, userInfo.Password)
	if token != "" {
		json.NewEncoder(w).Encode(resp)
	}
}*/ //contents are still copy pasted from the login portion, need to change to be email specific
