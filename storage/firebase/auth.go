package firebase

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"

	"encoding/json"

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
func (fb *App) GET(key string, password string) (map[string]interface{}, bool) {
	var user models.User
	var resp = make(map[string]interface{})
	users := fb.client.NewRef("users").OrderByChild("key").EqualTo(key).EqualTo(password).Get(context.Background(), &user)
	if users != nil {
		j, _ := json.Marshal(users)
		json.Unmarshal(j, &resp)
		return resp, true
	}
	return nil, false
}

func (fb *App) setTokens() {

}

// PUT sets new user in database
func (fb *App) PUT(uuid string, password string) (bool, error) {
	users := fb.client.NewRef("users")
	err := users.Set(context.Background(), map[string]*models.User{
		uuid: {
			FirstName: uuid,
			Password:  password,
		},
	})
	if err != nil {
		print(err.Error())
		return false, err
	}
	return true, nil
}

// SET changes user in database
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
