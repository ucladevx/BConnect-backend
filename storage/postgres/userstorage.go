package postgres

import (
	"fmt"
	"github.com/ucladevx/BConnect-backend/models"
	"github.com/ucladevx/BConnect-backend/utils/uuid"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// UserStorage user storage
type UserStorage struct {
	client *gorm.DB
}

// NewUserStorage new UserStorage
func NewUserStorage(client *gorm.DB) *UserStorage {
	return &UserStorage{
		client: client,
	}
}

func (us *UserStorage) create() {
	/* TODO: Maybe change migration model to maybe define DB relationships */
	us.client.AutoMigrate(&models.User{})
}

// GetUser gets user for login
func (us *UserStorage) GetUser(key string, password string) (*models.User, error) {
	user, err := us.findUser(key, password)
	return user, err
}

func (us *UserStorage) findUser(email, password string) (*models.User, error) {
	user := &models.User{}

	if err := us.client.Where("Email = ?", email).First(user).Error; err != nil {
		return nil, err
	}

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword {
		return nil, errf
	}

	return user, nil
}

// NewUser puts user into postgres
func (us *UserStorage) NewUser(user *models.User) (bool, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	newUser := *user
	newUser.Password = string(pass)
	newUser.UserID = uuid.UUID()
	if !us.client.NewRecord(newUser) {
		return false, nil
	}
	us.client.Create(&newUser)
	if us.client.NewRecord(newUser) {
		return false, nil
	}

	return true, nil
}

// ModifyUser sets updated fields
func (us *UserStorage) ModifyUser(userModded *models.User) (*models.User, error) {
	var user models.User
	if res := us.client.Where("EMAIL = ? AND USER_ID = ?", userModded.Email, userModded.UserID).Find(&user); res.Error != nil {
		return nil, res.Error
	}

	fname := userModded.FirstName
	lname := userModded.LastName
	major := userModded.Major
	gradYear := userModded.GradYear
	interests := userModded.Interests
	bio := userModded.Bio
	clubs := userModded.Clubs
	lat := userModded.Lat
	lon := userModded.Lon

	if userModded.FirstName == "" {
		fname = user.FirstName
	}
	if userModded.LastName == "" {
		lname = user.LastName
	}
	if userModded.Major == "" {
		major = user.Major
	}
	if userModded.GradYear == "" {
		gradYear = user.GradYear
	}
	if userModded.Bio == "" {
		bio = user.Bio
	}
	if userModded.Lat == 0 {
		lat = user.Lat
	}
	if userModded.Lon == 0 {
		lon = user.Lon
	}

	user.FirstName = fname
	user.LastName = lname
	user.Major = major
	user.GradYear = gradYear
	user.Interests = interests
	user.Bio = bio
	user.Clubs = clubs
	user.Lat = lat
	user.Lon = lon
	us.client.Save(&user)

	if err := us.client.Where("EMAIL = ? AND USER_ID = ?", userModded.Email, userModded.UserID).Find(&user); err != nil {
	}

	return &user, nil
}

func (us *UserStorage) AddFriend(userID string, friendID string, msg string) (*models.User, error) {
	var user models.User
	var friend models.User

	if res := us.client.Where(&models.User{UserID: userID}).First(&user); res.Error != nil {
		return nil, res.Error
	}
	if res := us.client.Where(&models.User{UserID: friendID}).First(&friend); res.Error != nil {
		return nil, res.Error
	}

	us.client.Model(&user).Association("Friends").Append(friend)
	return &user, nil
}

func (us *UserStorage) AddInterest(userID string, interestString string) (*models.User, error) {
	var user models.User
	var interestStruct models.Interest

	if res := us.client.Where(&models.User{UserID: userID}).First(&user); res.Error != nil {
		return nil, res.Error
	}
	if res := us.client.Where(&models.Interest{Interest: interestString}).First(&interestStruct); res.Error != nil {
		return nil, res.Error
	}

	interest := &models.Interest{Interest: interestString}
	us.client.Model(&user).Association("Interests").Append([]*models.Interest{interest})
	return &user, nil
}

func (us *UserStorage) AddClub(userID string, clubString string) (*models.User, error){
	var user models.User
	var clubStruct models.Club

	if res := us.client.Where(&models.User{UserID: userID}).First(&user); res.Error != nil {
		return nil, res.Error
	}
	if res := us.client.Where(&models.Club{Club: clubString}).First(&clubStruct); res.Error != nil {
		return nil, res.Error
	}

	us.client.Model(&user).Association("user_clubs").Append(clubStruct)
	return &user, nil
}

func (us *UserStorage) GetInterests(userID string) (map[string]interface{}, error) {
	var user models.User

	if res := us.client.Where(&models.User{UserID: userID}).First(&user); res.Error != nil {
		return nil, res.Error
	}

	numInterests := us.client.Model(&user).Association("Interests").Count()
	interests := make([]*models.Interest, numInterests)
	interestList := make([]string, numInterests)
	us.client.Model(&user).Related(&interests, "Interests")

	for _, val := range interests {
		interestList = append(interestList, val.Interest)
	}

	return map[string]interface{}{"num_interests": numInterests, "interests": interestList}, nil
}

func (us *UserStorage) GetClubs(userID string) (map[string]interface{}, error) {
	var user models.User

	if res := us.client.Where(&models.User{UserID: userID}).First(&user); res.Error != nil {
		return nil, res.Error
	}

	numClubs := us.client.Model(&user).Association("Clubs").Count()
	clubs := make([]*models.Club, numClubs)
	clubList := make([]string, numClubs)
	us.client.Model(&user).Association("Clubs").Find(&clubs)

	for _, val := range clubs {
		clubList = append(clubList, val.Club)
	}

	return map[string]interface{}{"num_clubs": numClubs, "clubs": clubList}, nil
}

func (us *UserStorage) GetFriends(userID string) (map[string]interface{}, error) {
	var user models.User

	if res := us.client.Where(&models.User{UserID: userID}).First(&user); res.Error != nil {
		return nil, res.Error
	}

	numFriends := us.client.Model(&user).Association("Friends").Count()
	friends := make([]*models.User, numFriends)
	us.client.Model(&user).Association("Friends").Find(&friends)

	type friend struct {
		ID string
		FirstName string
		LastName string
		GradYear string
	}

	friendList := make([]friend, numFriends)
	for _, val := range friends {
		friendList = append(friendList, friend{
			ID: val.UserID,
			FirstName: val.FirstName,
			LastName: val.LastName,
			GradYear: val.GradYear,
		})
	}

	return map[string]interface{}{"num_friends": numFriends, "friends": friendList}, nil
}

// DeleteUser dels clients
func (us *UserStorage) DeleteUser(key string, password string) (bool, error) {
	/*TODO: Implement*/
	return true, nil
}

//GetFromID gets user from uuid
func (us *UserStorage) GetFromID(uuid string) (*models.User, error){
	user := &models.User{}
	if res := us.client.Where("UUID = ?", uuid).First(user); res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

// Leave dummy function for logs users out
func (us *UserStorage) Leave(currUUID string) {

}

// Filter filters based on categories
func (us *UserStorage) Filter(finder models.Finder, filters map[string]models.Filterer, args map[string][]string) map[string]interface{} {
	return finder(filters, args)
}
