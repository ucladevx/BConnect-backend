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

func (intStore *InterestStorage) NewInterest(interestStr string) (*models.Interest, error) {
	interest := models.Interest{
		Interest: interestStr,
	}

	if !intStore.client.NewRecord(interest) {
		intStore.client.Create(&interest)
		return &interest, nil
	}

	return nil, errors.New("interest already exists")
}

func (intStore *InterestStorage) GetInterestFromString(interestStr string) (*models.Interest, error) {
	var interest models.Interest

	if res := intStore.client.Where(&models.Interest{Interest: interestStr}).First(&interest); res.Error != nil {
		return nil, res.Error
	}

	return &interest, nil
}

func (intStore *InterestStorage) GetAllInterests() ([]*models.Interest, error) {
	var interests []*models.Interest

	if res := intStore.client.Find(&interests); res.Error != nil {
		return nil, res.Error
	}

	return interests, nil
}

func (intStore *InterestStorage) AddUser(interest *models.Interest, user *models.User) (*models.Interest, error) {
	if res := intStore.client.Model(interest).Association("Users").Append(*user); res.Error != nil {
		return nil, res.Error
	}

	return interest, nil
}

func (intStore *InterestStorage) GetUsers(interest *models.Interest) ([]*models.User, error) {
	numUsers := intStore.client.Model(interest).Association("Users").Count()
	users := make([]*models.User, numUsers)

	if res := intStore.client.Model(interest).Association("Users").Find(&users); res.Error != nil {
		return nil, res.Error
	}

	return users, nil
}

