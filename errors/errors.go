package errors

import "fmt"

//InvalidTokenError invalid token
type InvalidTokenError struct {
	token string
}

func (e *InvalidTokenError) Error() string {
	return fmt.Sprintf("Invalid token: %s", e.token)
}
