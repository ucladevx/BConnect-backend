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

func (us *UserStorage) AddFriend(user *models.User, friend *models.User, msg string) (*models.User, error) {
	if res := us.client.Model(user).Association("Friends").Append(friend); res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

func (us *UserStorage) AddInterest(user *models.User, interest *models.Interest) (*models.User, error) {
	if res := us.client.Model(user).Association("Interests").Append(interest); res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

func (us *UserStorage) GetFriends(user *models.User) ([]*models.User, error) {
	numFriends := us.client.Model(user).Association("Friends").Count()
	friends := make([]*models.User, numFriends)

	if res := us.client.Model(user).Association("Friends").Find(&friends); res.Error != nil {
		return nil, res.Error
	}

	return friends, nil
}

func (us *UserStorage) GetInterests(user *models.User) ([]*models.Interest, error) {
	numInterests := us.client.Model(user).Association("Interests").Count()
	interests := make([]*models.Interest, numInterests)

	if res := us.client.Model(user).Association("Interests").Find(&interests); res.Error != nil {
		return nil, res.Error
	}

	return interests, nil
}

// DeleteUser dels clients
func (us *UserStorage) DeleteUser(key string, password string) (bool, error) {
	/*TODO: Implement*/
	return true, nil
}

//GetFromID gets user from uuid
func (us *UserStorage) GetFromID(uuid string) (*models.User, error) {
	user := &models.User{}
	if res := us.client.Where("UserID = ?", uuid).First(user); res.Error != nil {
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
