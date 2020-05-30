package postgres

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/ucladevx/BConnect-backend/models"
)

type ClubStorage struct {
	client *gorm.DB
}

func NewClubStorage(client *gorm.DB) *ClubStorage {
	return &ClubStorage{
		client: client,
	}
}

func (clubStore *ClubStorage) create() {
	clubStore.client.AutoMigrate(&models.Club{})
}

func (clubStore *ClubStorage) NewClubFromString(clubStr string) (*models.Club, error) {
	club := models.Club{
		Club: clubStr,
	}

	if !clubStore.client.NewRecord(club) {
		clubStore.client.Create(&club)
		return &club, nil
	}

	return nil, errors.New("club already exists")
}

func (clubStore ClubStorage) NewClub(club *models.Club) (*models.Club, error) {
	if !clubStore.client.NewRecord(club) {
		clubStore.client.Create(club)
		return club, nil
	}

	return nil, errors.New("club already exists")
}

func (clubStore *ClubStorage) GetClubFromString(clubStr string) (*models.Club, error) {
	var club models.Club

	if res := clubStore.client.Where(&models.Club{Club: clubStr}).First(&club); res.Error != nil {
		return nil, res.Error
	}

	return &club, nil
}

func (clubStore *ClubStorage) GetAllClubs() ([]*models.Club, error) {
	var clubs []*models.Club

	if res := clubStore.client.Find(&clubs); res.Error != nil {
		return nil, res.Error
	}

	return clubs, nil
}

func (clubStore *ClubStorage) AddUser(club *models.Club, user *models.User) (*models.Club, error) {
	if res := clubStore.client.Model(club).Association("Users").Append(*user); res.Error != nil {
		return nil, res.Error
	}

	return club, nil
}

func (clubStore *ClubStorage) GetUsers(club *models.Club) ([]*models.User, error) {
	numUsers := clubStore.client.Model(club).Association("Users").Count()
	users := make([]*models.User, numUsers)

	if res := clubStore.client.Model(club).Association("Users").Find(&users); res.Error != nil {
		return nil, res.Error
	}

	return users, nil
}

