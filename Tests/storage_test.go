package Tests
import (
	"github.com/ucladevx/BConnect-backend/storage/postgres"
	"testing"

	"github.com/ucladevx/BConnect-backend/models"
)

//UserStoreTest Tests
func TestUserStore(t *testing.T) {
	testDB := postgres.Connect("localhost",
		"postgres",
		"connect_b_test",
		"")

	userStore := postgres.NewUserStorage(testDB)
	intStore := postgres.NewInterestStorage(testDB)

	postgres.CreatePostgresTables(userStore, intStore)
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
		Interests: []*models.Interest{
			{Interest: "Diving"},
		},
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

	//moddedUser := models.User{
	//	UserID:    user.UserID,
	//	Email:     user.Email,
	//	FirstName: "ModificationTest",
	//}
	//testInterests := []models.Interest{
	//	models.Interest{
	//		Interest: "Diving",
	//	},
	//}
	//
	//for _, interest := range testInterests {
	//	userStore.NewInterest(&interest)
	//}

	//userStore.AddInterest(user.UserID, "Diving")

	interests, err := userStore.GetInterests(user.UserID)
	if interests["interests"] != "Diving" {
		t.Errorf("Error retrieving interests")
	}

	user, strErr = userStore.GetUser(testEmail, testPassword)
	if strErr != nil {
		t.Errorf("Error retrieving user")
	}
	if user.FirstName != "ModificationTest" {
		t.Errorf("Error modifying user")
	}

	currUUID := user.UserID
	uuidUser, err := userStore.GetFromID(currUUID)
	if user.UserID != uuidUser.UserID {
		t.Errorf("Error getting user by UUID")
	}

	testDB.DropTableIfExists(&models.User{})
	testDB.Close()
}

