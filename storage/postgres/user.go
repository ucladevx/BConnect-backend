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
	client.friend.AutoMigrate(&models.Friends{})
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

	fReq := &models.Friends{}
	fReq.UUID = currUUID
	fReq.FUUID = friendUUID
	fReq.FReqMess = optionalMsg
	fReq.TimeStamp = friendRequest.Timestamp
	fReq.Status = 0
	user.friend.Create(fReq)

	return &friendRequest, nil
}

// ACCEPT accepts added friends
func (user *UserActions) ACCEPT(currUUID string, friendUUID string) (*models.FriendRequest, error) {
	var friendRequest models.FriendRequest
	friendRequest.Sender = friendUUID

	fRequest := &models.Friends{}
	fRequest.UUID = currUUID
	fRequest.FUUID = friendUUID
	user.acceptReq(currUUID, friendUUID)

	return &friendRequest, nil
}

func (user *UserActions) acceptReq(currUUID string, friendUUID string) {
	var friend models.Friends
	backFriend := &models.Friends{}

	if err := user.friend.Where("F_UUID = ? AND UUID = ?", currUUID, friendUUID).Find(&friend); err != nil {
	}
	friend.Status = 1
	friend.TimeAccepted = time.Now().Unix()
	user.friend.Save(&friend)

	backFriend.UUID = currUUID
	backFriend.FUUID = friendUUID
	backFriend.Status = 1
	backFriend.TimeAccepted = time.Now().Unix()
	user.friend.Create(backFriend)
}

// GET gets users friends
func (user *UserActions) GET(currUUID string) map[string]interface{} {
	var friends []models.Friends
	if err := user.friend.Where("UUID = ?", currUUID).Where("Status = ?", 1).Find(&friends); err != nil {

	}
	type FriendResponse struct {
		FUUID        string
		TimeAccepted int64
	}
	var fuuids []FriendResponse
	var currFriendResponse FriendResponse
	for _, friend := range friends {
		currFriendResponse.FUUID = friend.FUUID
		currFriendResponse.TimeAccepted = friend.TimeAccepted
		fuuids = append(fuuids, currFriendResponse)
	}
	var resp = map[string]interface{}{"num_friends": len(friends), "friends": fuuids}
	return resp
}

// LEAVE logs users out
func (user *UserActions) LEAVE(currUUID string) {

}
