package actions

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// AuthService abstract server-side authentication in case we switch from whatever current auth scheme we are using
type AuthService interface {
	GET(username string, password string) (bool, error)
	SET() error
	PUT(key string) error
	DEL(key string) error
}

// AuthController abstract server-side authentication
type AuthController struct {
	service AuthService
}

// Auth credentials necessary for username/password auth
type Auth struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

// CurrUser refers to current user and corresponding JWT token
type CurrUser struct {
	Username string `json:"Username"`
	Token    string `json:"Token"`
}

// Setup sets up handlers
func (auth *AuthController) Setup(g echo.Group) {

}

// Login login users and provides auth token
func (auth *AuthController) Login(c echo.Context, cred *Auth) error {
	isValid, _ := auth.service.GET(cred.Username, cred.Password)
	if isValid {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = cred.Username
		claims["admin"] = false
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}
	return echo.ErrUnauthorized
}
