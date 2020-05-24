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
func (us *UserStorage) GetUser(key string, password string) (*models.User, string) {
	user, err := us.findUser(key, password)
	return user, err
}

func (us *UserStorage) findUser(email, password string) (*models.User, string) {
	user := &models.User{}

	if err := us.client.Where("Email = ?", email).First(user).Error; err != nil {
		return nil, ""
	}

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword {
		return nil, ""
	}

	return user, ""
}

// NewUser puts user into postgres
func (us *UserStorage) NewUser(user *models.User) (bool, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	user.Password = string(pass)
	user.UserID = uuid.UUID()
	if !us.client.NewRecord(user) {
		return false, nil
	}
	us.client.Create(user)
	if us.client.NewRecord(user) {
		return false, nil
	}
	return true, nil
}

// ModifyUser sets updated fields
func (us *UserStorage) ModifyUser(userModded *models.User) *models.User {
	var user models.User
	if err := us.client.Where("EMAIL = ? AND USER_ID = ?", userModded.Email, userModded.UserID).Find(&user); err != nil {
	}

	fname := userModded.FirstName
	lname := userModded.LastName
	age := userModded.Age
	currentJob := userModded.CurrentJob
	gender := userModded.Gender
	major := userModded.Major
	gradYear := userModded.GradYear
	bio := userModded.Bio
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
	if userModded.Gender == "" {
		gender = user.Gender
	}
	if userModded.CurrentJob == "" {
		currentJob = user.CurrentJob
	}
	if userModded.Age == 0 {
		age = user.Age
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
	user.Bio = bio
	user.Age = age
	user.Gender = gender
	user.CurrentJob = currentJob
	user.Lat = lat
	user.Lon = lon
	us.client.Save(&user)

	if err := us.client.Where("EMAIL = ? AND USER_ID = ?", userModded.Email, userModded.UserID).Find(&user); err != nil {
	}

	return &user
}

// DeleteUser dels clients
func (us *UserStorage) DeleteUser(key string, password string) (bool, error) {
	/*TODO: Implement*/
	return true, nil
}

//GetFromID gets user from uuid
func (us *UserStorage) GetFromID(uuid string) *models.User {
	user := &models.User{}
	if err := us.client.Where("USER_ID = ?", uuid).First(user).Error; err != nil {
		return nil
	}

	return user
}

// Leave dummy function for logs users out
func (us *UserStorage) Leave(currUUID string) {

}
