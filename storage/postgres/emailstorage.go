package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/ucladevx/BConnect-backend/models"
)

// EmailStorage user storage
type EmailStorage struct {
	client *gorm.DB
}

func (es *EmailStorage) create() {
	es.client.AutoMigrate(&models.Email{})
}

//NewEmailStorage generates new email storage
func NewEmailStorage(client *gorm.DB) *EmailStorage {
	return &EmailStorage{
		client: client,
	}
}

//AddEmail adds new email to storage
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
