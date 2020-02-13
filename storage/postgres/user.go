package postgres

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/ucladevx/BConnect-backend/models"
)

// UserActions holds the actions a user can do
type UserActions struct {
	client *gorm.DB
	friend *gorm.DB
}

func (client *UserActions) create() {
	client.friend.AutoMigrate(&Friends{})
}

// NewUserActions builds UserActions
func NewUserActions(client *gorm.DB, friendClient *gorm.DB) *UserActions {
	return &UserActions{
		client: client,
		friend: friendClient,
	}
}

// ADD adds friends
func (user *UserActions) ADD(currUUID string, friendUUID string, optionalMsg string) (*models.FriendRequest, error) {
	var friendRequest models.FriendRequest
	friendRequest.Message = optionalMsg
	friendRequest.Sender = currUUID
	friendRequest.Timestamp = time.Now().Unix()

	fReq := &Friends{}
	fReq.UUID = currUUID
	fReq.FUUID = friendUUID
	user.friend.Create(fReq)

	return &friendRequest, nil
}

// ACCEPT accepts added friends
func (user *UserActions) ACCEPT(currUUID string, friendUUID string) (*models.FriendRequest, error) {
	var friendRequest models.FriendRequest
	friendRequest.Sender = friendUUID

	fReq := &Friends{}
	fReq.UUID = currUUID
	fReq.UUID = friendUUID
	user.friend.Create(fReq)

	return &friendRequest, nil
}
