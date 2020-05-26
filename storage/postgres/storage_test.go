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

	testUser := &models.User{
		Email:     testEmail,
		Password:  testPassword,
		FirstName: testFirstName,
		LastName:  testLastName,
	}
	status, err := userStore.NewUser(testUser)
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
	testInterests := []models.Interests{
		models.Interests{
			Interest: "Diving",
		},
	}

	userStore.ModifyUser(&moddedUser, testInterests)

	interests := userStore.GetInterestsFromID(user.UserID)
	if interests[0].Interest != "Diving" {
		t.Errorf("Error retrieving interests")
	}

	user, strErr = userStore.GetUser(testEmail, testPassword)
	if strErr != "" {
		t.Errorf("Error retrieving user")
	}
	if user.FirstName != "ModificationTest" {
		t.Errorf("Error modifying user")
	}
	if user.LastName != "Test" {
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

	testUser1 := &models.User{
		Email:     testEmail,
		Password:  testPassword,
		FirstName: testFirstName,
		LastName:  testLastName,
	}

	userStore.NewUser(testUser1)

	user1, strErr1 := userStore.GetUser(testEmail, testPassword)
	if strErr1 != "" {
		t.Errorf("Error retrieving user")
	}

	testEmail2 := "test2@ucla.edu"
	testPassword2 := "test2_password"
	testFirstName2 := "Test2"
	testLastName2 := "Test2"

	testUser2 := &models.User{
		Email:     testEmail2,
		Password:  testPassword2,
		FirstName: testFirstName2,
		LastName:  testLastName2,
	}
	userStore.NewUser(testUser2)

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
