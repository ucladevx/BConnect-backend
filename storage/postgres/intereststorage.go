package postgres

import (
	"errors"

	"github.com/ucladevx/BConnect-backend/models"
	"github.com/jinzhu/gorm"
)

type InterestStorage struct {
	client *gorm.DB
}

func NewInterestStorage(client *gorm.DB) *InterestStorage {
	return &InterestStorage{
		client: client,
	}
}

func (intStore *InterestStorage) create() {
	intStore.client.AutoMigrate(&models.Interest{})
}

func (intStore *InterestStorage) NewInterest(interestStr string) error {
	interest := models.Interest{
		Interest: interestStr,
	}

	if intStore.client.NewRecord(interest) {
		intStore.client.Create(&interest)
		return nil
	}

	return errors.New("interest already exists")
}

func (intStore *InterestStorage) GetInterest(interestStr string) (*models.Interest, error) {
	var interest models.Interest

	if res := intStore.client.Where(&models.Interest{Interest: interestStr}).First(&interest); res.Error != nil {
		return nil, res.Error
	}

	return &interest, nil
}

func (intStore *InterestStorage) AddUser(interestStr string, user *models.User) (*models.Interest, error) {
	interest, err := intStore.GetInterest(interestStr)
	if err != nil {
		return nil, err
	}

	intStore.client.Model(interest).Association("Users").Append(*user)
	return interest, nil
}

func (intStore *InterestStorage) GetUsers(interestStr string) (map[string]interface{}, error) {
	interest, err := intStore.GetInterest(interestStr)
	if err != nil {
		return nil, err
	}

	numUsers := intStore.client.Model(interest).Association("Users").Count()
	users := make([]*models.User, numUsers)
	userList := make([]string, numUsers)
	intStore.client.Model(interest).Association("Users").Find(&users)

	for _, val := range users {
		userList = append(userList, val.UserID)
	}

	return map[string]interface{}{"num_users": numUsers, "users": userList}, nil
}

func (intStore *InterestStorage) ListInterests() (map[string]interface{}, error) {
	
}
