package firebase

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"

	"github.com/ucladevx/BConnect-backend/models"
)

// App struct to hold firebase functionality
type App struct {
	app    *firebase.App
	client *db.Client
}

// NewFirebaseApp struct to hold firebase functionality
func NewFirebaseApp(app *firebase.App, client *db.Client) *App {
	return &App{
		app:    app,
		client: client,
	}
}

// GET gets user by username and password from database
func (fb *App) GET(key string, password string) (*models.User, bool) {
	var user models.User
	users := fb.client.NewRef("users").OrderByChild("key").EqualTo(key).EqualTo(password).Get(context.Background(), &user)
	if users != nil {
		return &user, true
	}
	return nil, false
}

// PUT sets new user in database
func (fb *App) PUT(uuid string, password string) (bool, error) {
	users := fb.client.NewRef("users")
	err := users.Set(context.Background(), map[string]*models.User{
		uuid: {
			FirstName: "name",
		},
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

// SETS changes user in database
func (fb *App) SET(key string, password string) (bool, error) {
	users := fb.client.NewRef("users")
	err := users.Set(context.Background(), map[string]*models.User{
		key: {
			FirstName: "name",
		},
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

// DEL deletes user in database
func (fb *App) DEL(key string, password string) (bool, error) {
	users := fb.client.NewRef("users")
	err := users.Set(context.Background(), map[string]*models.User{
		key: {
			FirstName: "name",
		},
	})
	if err != nil {
		return false, err
	}
	return true, nil
}
