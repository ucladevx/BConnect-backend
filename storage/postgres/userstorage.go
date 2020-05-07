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
func (us *UserStorage) NewUser(email string, password string, firstname string, lastname string) (bool, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	user := &models.User{}
	user.Password = string(pass)
	user.Email = email
	user.FirstName = firstname
	user.LastName = lastname
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

	return &user
}

func (us *UserStorage) AddFriend(userID string, friend *models.User) *models.User {
	var user models.User

	if err := us.client.Where(&models.User{UserID: userID}).First(&user); err == nil {
		us.client.Model(&user).Association("Friends").Append(friend)
	}

	return &user
}

func (us *UserStorage) AddInterest(userID string, interest *models.Interest) *models.User {
	var user models.User

	if err := us.client.Where(&models.User{UserID: userID}).First(&user); err == nil {
		us.client.Model(&user).Association("Interests").Append(interest)
	}

	return &user
}

func (us *UserStorage) AddClub(userID string, club *models.Club) *models.User {
	var user models.User

	if err := us.client.Where(&models.User{UserID: userID}).First(&user); err == nil {
		us.client.Model(&user).Association("Clubs").Append(club)
	}

	return &user
}

func (us *UserStorage) GetInterests(userID string) []*models.Interest {
	var user models.User

	if err := us.client.Where(&models.User{UserID: userID}).First(&user); err == nil {
		interests := make([]*models.Interest, us.client.Model(&user).Association("Interests").Count())
		us.client.Model(&user).Association("Interests").Find(&interests)
		return interests
	}

	return nil
}

func (us *UserStorage) GetClubs(userID string) []*models.Club {
	var user models.User

	if err := us.client.Where(&models.User{UserID: userID}).First(&user); err == nil {
		clubs := make([]*models.Club, us.client.Model(&user).Association("Clubs").Count())
		us.client.Model(&user).Association("Clubs").Find(&clubs)
		return clubs
	}

	return nil
}

func (us *UserStorage) GetFriends(userID string) []*models.User {
	var user models.User

	if err := us.client.Where(&models.User{UserID: userID}).First(&user); err == nil {
		friends := make([]*models.User, us.client.Model(&user).Association("Friends").Count())
		us.client.Model(&user).Association("Friends").Find(&friends)
		return friends
	}

	return nil
}

// DeleteUser dels clients
func (us *UserStorage) DeleteUser(key string, password string) (bool, error) {
	/*TODO: Implement*/
	return true, nil
}

//GetFromID gets user from uuid
func (us *UserStorage) GetFromID(uuid string) *models.User {
	user := &models.User{}
	if err := us.client.Where("UUID = ?", uuid).First(user).Error; err != nil {
		return nil
	}

	return user
}

// Leave dummy function for logs users out
func (us *UserStorage) Leave(currUUID string) {

}
