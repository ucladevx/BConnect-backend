package users

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ucladevx/BConnect-backend/models"
)

//UserStorage user store
type UserStorage interface {
	GetUser(username string, password string) (*models.User, error)
	NewUser(email string, uuid string, firstname string, lastname string) (bool, error)
	ModifyUser(user *models.User) (*models.User, error)
	AddFriend(userID string, friend string, msg string) (*models.User, error)
	AddInterest(userID string, interestString string) (*models.User, error)
	AddClub(userID string, clubString string) (*models.User, error)
	GetInterests(userID string) (map[string]interface{}, error)
	GetClubs(userID string) (map[string]interface{}, error)
	GetFriends(userID string) (map[string]interface{}, error)
	DeleteUser(username string, password string) (bool, error)
	GetFromID(uuid string) (*models.User, error)
	Leave(currUUID string)
	Filter(finder models.Finder, filters map[string]models.Filterer, args map[string][]string) map[string]interface{}
}

//UserService holds services for users
type UserService struct {
	userStore   UserStorage
}

//NewUserService constructs new user service
func NewUserService(userStore UserStorage) *UserService {
	return &UserService{
		userStore:   userStore,
	}
}

//Login login
func (us *UserService) Login(username string, password string) (map[string]interface{}, string, string, time.Time, time.Time) {
	user, err := us.userStore.GetUser(username, password)
	if err != nil {

	}
	if user == nil {
		var resp = map[string]interface{}{"status": false, "message": "Invalid login credentials"}
		return resp, "", "", time.Time{}, time.Time{}
	}

	expiresAt := time.Now().Add(time.Minute * 15)

	tk := &models.Token{
		UUID:  user.UserID,
		Email: user.Email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}

	refreshString, stringErr := generateRandomString(32)
	if stringErr != nil {

	}

	refreshExpiresAt := time.Now().Add(time.Hour * 144)

	refreshTk := &models.RefreshToken{
		UUID: user.UserID,
		ID:   refreshString,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: refreshExpiresAt.Unix(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), refreshTk)

	refreshTokenString, error := refreshToken.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}

	var resp = map[string]interface{}{"status": false, "message": "logged in"}
	resp["token"] = tokenString
	resp["refresh"] = refreshTokenString
	resp["user"] = user

	return resp, tokenString, refreshTokenString, expiresAt, refreshExpiresAt
}

func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

//Signup signs user in
func (us *UserService) Signup(email string, uuid string, firstname string, lastname string) (bool, error) {
	return us.userStore.NewUser(email, uuid, firstname, lastname)
}

//Update sets categories
func (us *UserService) Update(user *models.User) (map[string]interface{}, error) {
	var resp = map[string]interface{}{"status": false, "message": "logged in"}
	user, err := us.userStore.ModifyUser(user)
	resp["mod_user"] = user

	return resp, err
}

//DeleteUser delete user
func (us *UserService) DeleteUser(username string, password string) (bool, error) {
	return us.userStore.DeleteUser(username, password)
}

//Leave dummy function
func (us *UserService) Leave(currUUID string) {
	us.userStore.Leave(currUUID)
}

//RefreshToken refresh token
func (us *UserService) RefreshToken(uuid string) (map[string]interface{}, string, time.Time) {
	user, err := us.userStore.GetFromID(uuid)
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "UUID not found"}
		return resp, "", time.Time{}
	}

	expiresAt := time.Now().Add(time.Minute * 15)

	tk := &models.Token{
		UUID:  user.UserID,
		Email: user.Email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}

	var resp = map[string]interface{}{"status": false, "message": "logged in"}
	resp["token"] = tokenString
	resp["user"] = user

	return resp, tokenString, expiresAt
}

//FriendRequest add user recipe
func (us *UserService) FriendRequest(currUUID string, friendUUID string, optionalMsg string) (*models.User, error) {
	return us.userStore.AddFriend(currUUID, friendUUID, optionalMsg)
}

//GetFriends removes user recipe
func (us *UserService) GetFriends(currUUID string) map[string]interface{} {
	res, err := us.userStore.GetFriends(currUUID)
	if err != nil {
		// handle error
	}

	return res
}

func (us *UserService) GetInterests(currUUID string) map[string]interface{} {
	res, err := us.userStore.GetInterests(currUUID)
	if err != nil {
		// handle error
	}

	return res
}

func (us *UserService) GetClubs(currUUID string) map[string]interface{} {
	res, err := us.userStore.GetClubs(currUUID)
	if err != nil {
		// handle error
	}

	return res
}

//Filter filters
func (us *UserService) Filter(finder models.Finder, filters map[string]models.Filterer, args map[string][]string) map[string]interface{} {
	return us.userStore.Filter(finder, filters, args)
}
