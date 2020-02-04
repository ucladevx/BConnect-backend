package firebase

import (
	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
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
func (fb *App) GET(username string, password string) {

}

// SET sets new user in database
func (fb *App) SET(username string, password string) {

}
