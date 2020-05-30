package Tests

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/ucladevx/BConnect-backend/storage/postgres"
	"testing"

	"github.com/ucladevx/BConnect-backend/models"
)

func initTestDB() (*gorm.DB, *postgres.UserStorage, *postgres.InterestStorage, *postgres.ClubStorage) {
	testDB := postgres.Connect("localhost",
		"postgres",
		"connect_b_test",
		"")

	userStore := postgres.NewUserStorage(testDB)
	intStore := postgres.NewInterestStorage(testDB)
	clubStore := postgres.NewClubStorage(testDB)

	postgres.CreatePostgresTables(userStore, intStore, clubStore)

	return testDB, userStore, intStore, clubStore
}

func cleanup(testDB *gorm.DB) {
	testDB.DropTableIfExists(&models.User{}, &models.Interest{})
	testDB.Close()
}

func compareInterests(interest1 *models.Interest, interest2 *models.Interest) bool {
	if interest1 == nil || interest2 == nil {
		return false
	}

	return interest1.Interest == interest2.Interest
}

func compareClubs(club1 *models.Club, club2 *models.Club) bool {
	if club1 == nil || club2 == nil {
		return false
	}

	return club1.Club == club2.Club
}

func compareUsers(user1 *models.User, user2 *models.User) bool {
	if user1 == nil || user2 == nil {
		return false
	}

	return user1.UserID == user2.UserID
}

//UserStoreTest Tests
func TestUserStore(t *testing.T) {
	testDB, userStore, _, _ := initTestDB()

	if !testDB.HasTable(&models.User{}) {
		t.Errorf("Table `user` not created")
	}

	testEmail := "test@ucla.edu"
	testPassword := "test_password"
	testFirstName := "Test"
	testLastName := "Test"

	testInterest1 := &models.Interest{
		Interest: "Being annoyed by Go syntax",
	}

	testInterest2 := &models.Interest{
		Interest: "Sleeping",
	}

	testUser := &models.User{
		Email:     testEmail,
		Password:  testPassword,
		FirstName: testFirstName,
		LastName:  testLastName,
		UserID:    "1234",
		Interests: []*models.Interest{
			testInterest1,
		},
	}

	testFriend := &models.User{
		Email:     "bob@bob.com",
		Password:  "securepassword",
		FirstName: "Bob",
		LastName:  "Bob",
		UserID:    "BobSpelledBackwards",
	}

	status, err := userStore.NewUser(testUser)
	if err != nil {
		t.Errorf("Error creating new user")
	}
	if !status {
		t.Errorf("Error creating new user")
	}

	user, strErr := userStore.GetUser(testEmail, testPassword)
	if strErr != nil {
		t.Errorf("Error retrieving user")
	}
	if user.Email != "test@ucla.edu" {
		t.Errorf("Error retrieving user")
	}

	user, err = userStore.AddInterest(user, testInterest2)
	if err != nil {
		t.Error("error adding interest")
	}

	interests, err := userStore.GetInterests(user)
	if err != nil {
		t.Error("error getting interests")
	}
	for _, val := range interests {
		if !(compareInterests(val, testInterest1) || compareInterests(val, testInterest2)) {
			t.Error("get interests: interests do not match")
		}
	}

	user, err = userStore.AddFriend(user, testFriend, "")
	if err != nil {
		t.Error("error adding friend")
	}

	friends, err := userStore.GetFriends(user)
	if err != nil {
		t.Error("error getting friends")
	}
	for _, val := range friends {
		if !compareUsers(testFriend, val) {
			t.Error("get friends: friends do not match")
		}
	}

	//currUUID := user.UserID
	//uuidUser, err := userStore.GetFromID(currUUID)
	//if err != nil {
	//	t.Errorf("Error getting user by UUID")
	//}
	//if compareUsers(user, uuidUser) {
	//	t.Error("getfromid: user does not match")
	//}

	cleanup(testDB)
}

func TestInterestStore(t *testing.T) {
	testDB, _, intStore, _ := initTestDB()

	testInterest1 := &models.Interest{
		Interest: "Being annoyed by Go syntax",
	}

	testUser := &models.User{
		Email:     "bob@bob.com",
		Password:  "thisisapassword",
		FirstName: "Bob",
		LastName:  "Bob",
		UserID:    "1234",
	}

	interest, err := intStore.NewInterestFromString("Eating")
	if err != nil {
		t.Error(err, "error creating new interest from string")
	}
	if interest == nil || interest.Interest != "Eating" {
		t.Error("newinterestfromstring: interest does not match")
	}

	interest, err = intStore.NewInterest(testInterest1)
	if err != nil {
		t.Error("error creating new interest")
	}
	if interest == nil || !compareInterests(interest, testInterest1) {
		t.Error("newinterest: interest does not match")
	}

	interest, err = intStore.AddUser(testInterest1, testUser)
	if err != nil {
		t.Error("error adding user")
	}

	users, err := intStore.GetUsers(testInterest1)
	if err != nil {
		t.Error("error getting users")
	}
	for _, val := range users {
		if !compareUsers(val, testUser) {
			t.Error("getusers: users do not match")
		}
	}

	interest, err = intStore.GetInterestFromString("Eating")
	if err != nil {
		t.Error("error getting interest from string")
	}
	if interest != nil && interest.Interest != "Eating" {
		t.Error("getinterestfromstring: interest does not match")
	}

	interests, err := intStore.GetAllInterests()
	if err != nil {
		t.Error("error getting all interests")
	}
	if len(interests) != 2 {
		fmt.Print("num interests: ", len(interests))
		t.Error("getallinterests: wrong number of interests")
	}
	for _, val := range interests {
		if !(compareInterests(val, testInterest1) || compareInterests(val, interest)) {
			t.Error("getallinterests: interests do not match")
		}
	}

	cleanup(testDB)
}

func TestClubStore(t *testing.T) {
	testDB, _, _, clubStore := initTestDB()

	testClub1 := &models.Club{
		Club: "Being annoyed by Go syntax",
	}

	testUser := &models.User{
		Email:     "bob@bob.com",
		Password:  "thisisapassword",
		FirstName: "Bob",
		LastName:  "Bob",
		UserID:    "1234",
	}

	club, err := clubStore.NewClubFromString("Eating")
	if err != nil {
		t.Error(err, "error creating new club from string")
	}
	if club == nil || club.Club != "Eating" {
		t.Error("newclubfromstring: club does not match")
	}

	club, err = clubStore.NewClub(testClub1)
	if err != nil {
		t.Error("error creating new club")
	}
	if club == nil || !compareClubs(club, testClub1) {
		t.Error("newclub: club does not match")
	}

	club, err = clubStore.AddUser(testClub1, testUser)
	if err != nil {
		t.Error("error adding user")
	}

	users, err := clubStore.GetUsers(testClub1)
	if err != nil {
		t.Error("error getting users")
	}
	for _, val := range users {
		if !compareUsers(val, testUser) {
			t.Error("getusers: users do not match")
		}
	}

	club, err = clubStore.GetClubFromString("Eating")
	if err != nil {
		t.Error("error getting club from string")
	}
	if club != nil && club.Club != "Eating" {
		t.Error("getclubfromstring: club does not match")
	}

	clubs, err := clubStore.GetAllClubs()
	if err != nil {
		t.Error("error getting all clubs")
	}
	if len(clubs) != 2 {
		fmt.Print("num clubs: ", len(clubs))
		t.Error("getallclubs: wrong number of clubs")
	}
	for _, val := range clubs {
		if !(compareClubs(val, testClub1) || compareClubs(val, club)) {
			t.Error("getallclubs: clubs do not match")
		}
	}

	cleanup(testDB)
}
