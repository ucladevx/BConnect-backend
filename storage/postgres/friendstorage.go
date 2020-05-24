package postgres

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/ucladevx/BConnect-backend/models"
)

// FriendStorage holds the actions friends can do
type FriendStorage struct {
	client *gorm.DB
}

func (fs *FriendStorage) create() {
	fs.client.AutoMigrate(&models.User{}, &models.Friends{})
}

// NewFriendStorage builds FriendStorage
func NewFriendStorage(client *gorm.DB) *FriendStorage {
	return &FriendStorage{
		client: client,
	}
}

// AddFriend adds friends
func (fs *FriendStorage) AddFriend(currUUID string, friendUUID string, optionalMsg string) (*models.Friends, error) {
	friendRequest := &models.Friends{}
	friendRequest.UserID = currUUID
	friendRequest.FriendID = friendUUID
	friendRequest.FReqMess = optionalMsg
	friendRequest.TimeStamp = time.Now().Unix()
	friendRequest.Status = 0
	fs.client.Create(friendRequest)

	return friendRequest, nil
}

//AcceptFriend accepts added friends
func (fs *FriendStorage) AcceptFriend(currUUID string, friendUUID string) (*models.Friends, error) {
	friendRequest := &models.Friends{}
	friendRequest.UserID = currUUID
	friendRequest.FriendID = friendUUID

	fs.acceptRequest(currUUID, friendUUID)

	return friendRequest, nil
}

func (fs *FriendStorage) acceptRequest(currUUID string, friendUUID string) {
	var friend models.Friends
	backFriend := &models.Friends{}

	if err := fs.client.Where("FRIEND_ID = ? AND USER_ID = ?", currUUID, friendUUID).Find(&friend); err != nil {
	}
	friend.Status = 1
	friend.TimeAccepted = time.Now().Unix()
	fs.client.Save(&friend)

	backFriend.UserID = currUUID
	backFriend.FriendID = friendUUID
	backFriend.Status = 1
	backFriend.TimeAccepted = time.Now().Unix()
	fs.client.Create(backFriend)
}

// GetFriends gets users friends
func (fs *FriendStorage) GetFriends(currUUID string) []models.Friends {
	var friends []models.Friends
	if err := fs.client.Where("USER_ID = ?", currUUID).Where("Status = ?", 1).Find(&friends); err != nil {

	}
	return friends
}

// Filter filters based on categories
func (fs *FriendStorage) Filter(finder models.Finder, filters map[string]models.Filterer, args map[string][]string) []models.User {
	return finder(filters, args)
}
