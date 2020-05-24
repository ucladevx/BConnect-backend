package postgres

import (

	"github.com/ucladevx/BConnect-backend/models"
	"github.com/jinzhu/gorm"

)

// EmailStorage user storage
type EmailStorage struct {
	client *gorm.DB
}

func (es *EmailStorage) create() {
	es.client.AutoMigrate(&models.Email{})
}

func NewEmailStorage(client *gorm.DB) *EmailStorage {
	return &EmailStorage{
		client: client,
	}
}

func (es *EmailStorage) AddEmail(email string) (bool, error) {
	user := &models.Email{}
	user.Email = email
	
	//returns false if the email already exists
	if !es.client.NewRecord(user) {
		return false, nil
	}
	es.client.Create(user)

	//returns false and nil if the email wasn't properly created
	if es.client.NewRecord(user) {
		return false, nil
	}
	return true, nil
}
