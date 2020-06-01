package postgres

import (
	"testing"

	"github.com/ucladevx/BConnect-backend/models"
)

/* Simple tests for storage functionality - more precise testing can be implemented later */

//UserStoreTest Tests
func TestUserStore(t *testing.T) {
	testDB := Connect("localhost",
		"connect_b",
		"connect_b_test",
		"connect_b")
	testDB.DropTableIfExists(&models.User{})
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

	user, loginErr := userStore.GetUser(testEmail, testPassword)
	if loginErr != nil {
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

	user, loginErr = userStore.GetUser(testEmail, testPassword)
	if loginErr != nil {
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

	user1, loginErr1 := userStore.GetUser(testEmail, testPassword)
	if loginErr1 != nil {
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

	user2, loginErr2 := userStore.GetUser(testEmail2, testPassword2)
	if loginErr2 != nil {
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

func TestFilters(t *testing.T) {
	testDB := Connect("localhost",
		"connect_b",
		"connect_b_test",
		"connect_b")
	userStore := NewUserStorage(testDB)
	filterer := NewFilterers(testDB)
	CreatePostgresTables(userStore)
	if !testDB.HasTable(&models.User{}) {
		t.Errorf("Table `user` not created")
	}

	testEmail := "test@ucla.edu"
	testPassword := "test_password"
	testFirstName := "Test"
	testLastName := "Test"
	testLatitude := 34.2
	testLongitude := 34.2

	testUser1 := &models.User{
		Email:     testEmail,
		Password:  testPassword,
		FirstName: testFirstName,
		LastName:  testLastName,
		Lat:       testLatitude,
		Lon:       testLongitude,
		GradYear:  "2022",
		Major:     "Philosophy",
	}

	userStore.NewUser(testUser1)

	user1, loginErr1 := userStore.GetUser(testEmail, testPassword)
	if loginErr1 != nil {
		t.Errorf("Error retrieving user")
	}

	testEmail2 := "test2@ucla.edu"
	testPassword2 := "test2_password"
	testFirstName2 := "Test2"
	testLastName2 := "Test2"
	testLatitude2 := 34.21
	testLongitude2 := 34.21

	testUser2 := &models.User{
		Email:     testEmail2,
		Password:  testPassword2,
		FirstName: testFirstName2,
		LastName:  testLastName2,
		Lat:       testLatitude2,
		Lon:       testLongitude2,
		GradYear:  "2022",
		Major:     "Art",
	}

	userStore.NewUser(testUser2)

	curr := &models.FilterReturn{
		Filter: filterer.client,
	}
	var user []models.User

	filterer.LocationRadiusFilter(curr, user1, []string{"0.070"}).Filter.Find(&user)
	if user[0].FirstName != "Test2" {
		t.Errorf("Error on filtering by location")
	}

	filterer.NameFilter(curr, user1, []string{"Test2"}).Filter.Find(&user)
	if user[0].FirstName != "Test2" {
		t.Errorf("Error on filtering by first name")
	}

	filterer.MajorFilter(curr, user1, []string{"Art"}).Filter.Find(&user)
	if user[0].FirstName != "Test2" {
		t.Errorf("Error on filtering by major")
	}

	filterer.GradYearFilter(curr, user1, []string{"2022"}).Filter.Find(&user)
	if user[0].FirstName != "Test2" {
		t.Errorf("Error on filtering by grad year")
	}

	filterer.NameFilter(filterer.LocationRadiusFilter(curr, user1, []string{"0.070"}), user1, []string{"Test2"}).Filter.Find(&user)
	if user[0].FirstName != "Test2" {
		t.Errorf("Error on filtering by grad year")
	}

	testDB.DropTableIfExists(&models.User{})
	testDB.Close()
}

func TestEmailStore(t *testing.T) {
	testDB := Connect("localhost",
		"connect_b",
		"connect_b_test",
		"connect_b")
	emailStore := NewEmailStorage(testDB)
	CreatePostgresTables(emailStore)

	if _, err := emailStore.AddEmail("test@gmail.com"); err != nil {
		t.Errorf("Error adding email")
	}
	var checkEmail models.Email
	testDB.Where("EMAIL = ?", "test@gmail.com").Find(&checkEmail)
	if checkEmail.Email != "test@gmail.com" {
		t.Errorf("Error adding email")
	}

	testDB.DropTableIfExists(&models.Email{})
	testDB.Close()
}
