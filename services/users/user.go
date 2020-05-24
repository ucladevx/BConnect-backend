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
	GetUser(username string, password string) (*models.User, string)
	NewUser(user *models.User, interests []models.Interests) (bool, error)
	ModifyUser(user *models.User) *models.User
	DeleteUser(username string, password string) (bool, error)
	GetFromID(uuid string) *models.User
	Leave(currUUID string)
}

//FriendStorage user friend store
type FriendStorage interface {
	AddFriend(currUUID string, friendUUID string, optionalMsg string) (*models.Friends, error)
	AcceptFriend(currUUID string, friendUUID string) (*models.Friends, error)
	GetFriends(currUUID string) []models.Friends
	Filter(finder models.Finder, filters map[string]models.Filterer, args map[string][]string) []models.User
}

//EmailStorage email store
type EmailStorage interface {
	AddEmail(email string) (bool, error)
}

//UserService holds services for users
type UserService struct {
	userStore   UserStorage
	friendStore FriendStorage
	emailStore  EmailStorage
}

//NewUserService constructs new user service
func NewUserService(userStore UserStorage, friendStore FriendStorage, emailStore EmailStorage) *UserService {
	return &UserService{
		userStore:   userStore,
		friendStore: friendStore,
		emailStore:  emailStore,
	}
}

//Login login
func (us *UserService) Login(username string, password string) (map[string]interface{}, string, string, time.Time, time.Time) {
	user, err := us.userStore.GetUser(username, password)
	if err != "" {

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
func (us *UserService) Signup(user *models.User, interestsList []string) (bool, error) {
	var interests []models.Interests

	for _, interest := range interestsList {
		interests = append(interests, models.Interests{
			Interest: interest,
		})
	}

	return us.userStore.NewUser(user, interests)
}

//Update sets categories
func (us *UserService) Update(user *models.User) (map[string]interface{}, error) {
	var resp = map[string]interface{}{"status": false, "message": "logged in"}
	resp["mod_user"] = us.userStore.ModifyUser(user)
	return resp, nil
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
	user := us.userStore.GetFromID(uuid)
	if user == nil {
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

//FriendRequest adds friend
func (us *UserService) FriendRequest(currUUID string, friendUUID string, optionalMsg string) (*models.Friends, error) {
	return us.friendStore.AddFriend(currUUID, friendUUID, optionalMsg)
}

//AcceptFriendRequest adds friend
func (us *UserService) AcceptFriendRequest(currUUID string, friendUUID string) (*models.Friends, error) {
	return us.friendStore.AcceptFriend(currUUID, friendUUID)
}

//GetFriends gets friends
func (us *UserService) GetFriends(currUUID string) []models.Friends {
	return us.friendStore.GetFriends(currUUID)
}

//Filter filters
func (us *UserService) Filter(finder models.Finder, filters map[string]models.Filterer, args map[string][]string) []models.User {
	return us.friendStore.Filter(finder, filters, args)
}

//AddEmail adds email to storage
func (us *UserService) AddEmail(email string) (bool, error) {
	return us.emailStore.AddEmail(email)
}
