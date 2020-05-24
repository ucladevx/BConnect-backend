package postgres

import (
	"testing"

	"github.com/ucladevx/BConnect-backend/models"
)

//UserStoreTest Tests
func TestUserStore(t *testing.T) {
	testDB := Connect("localhost",
		"connect_b",
		"connect_b_test",
		"connect_b")
	userStore := NewUserStorage(testDB)
	CreatePostgresTables(userStore)
	if !testDB.HasTable(&models.User{}) {
		t.Errorf("Table `user` not created")
	}

	testEmail := "test@ucla.edu"
	testPassword := "test_password"
	testFirstName := "Test"
	testLastName := "Test"
	status, err := userStore.NewUser(testEmail, testPassword, testFirstName, testLastName)
	if err != nil {
		t.Errorf("Error creating new user")
	}
	if !status {
		t.Errorf("Error creating new user")
	}

	user, strErr := userStore.GetUser(testEmail, testPassword)
	if strErr != "" {
		t.Errorf("Error retrieving user")
	}
	if user.Email != "test@ucla.edu" {
		t.Errorf("Error retrieving user")
	}

	moddedUser := models.User{
		UserID:    user.UserID,
		Email:     user.Email,
		FirstName: "ModificationTest",
	}
	userStore.ModifyUser(&moddedUser)

	user, strErr = userStore.GetUser(testEmail, testPassword)
	if strErr != "" {
		t.Errorf("Error retrieving user")
	}
	if user.FirstName != "ModificationTest" {
		t.Errorf("Error modifying user")
	}

	currUUID := user.UserID
	uuidUser := userStore.GetFromID(currUUID)
	if *user != *uuidUser {
		t.Errorf("Error getting user by UUID")
	}

	testDB.DropTableIfExists(&models.User{})
	testDB.Close()
}

func TestFriendStore(t *testing.T) {
	testDB := Connect("localhost",
		"connect_b",
		"connect_b_test",
		"connect_b")
	userStore := NewUserStorage(testDB)
	friendStore := NewFriendStorage(testDB)
	CreatePostgresTables(userStore, friendStore)

	testEmail := "test@ucla.edu"
	testPassword := "test_password"
	testFirstName := "Test"
	testLastName := "Test"
	userStore.NewUser(testEmail, testPassword, testFirstName, testLastName)

	user1, strErr1 := userStore.GetUser(testEmail, testPassword)
	if strErr1 != "" {
		t.Errorf("Error retrieving user")
	}

	testEmail2 := "test2@ucla.edu"
	testPassword2 := "test2_password"
	testFirstName2 := "Test2"
	testLastName2 := "Test2"
	userStore.NewUser(testEmail2, testPassword2, testFirstName2, testLastName2)

	user2, strErr2 := userStore.GetUser(testEmail2, testPassword2)
	if strErr2 != "" {
		t.Errorf("Error retrieving user")
	}

	_, err := friendStore.AddFriend(user1.UserID, user2.UserID, "")
	if err != nil {
		t.Errorf("Error adding friend")
	}

	_, err = friendStore.AcceptFriend(user2.UserID, user1.UserID)
	if err != nil {
		t.Errorf("Error accepting friend")
	}

	friends := friendStore.GetFriends(user1.UserID)
	if friends[0].FriendID != user2.UserID {
		t.Errorf("Error getting friends")
	}

	testDB.DropTableIfExists(&models.User{}, &models.Friends{})
	testDB.Close()
}
