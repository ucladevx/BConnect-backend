package errors

import "fmt"

//LoginError login error wrapper
type LoginError struct {
	Issue error
}

func (e *LoginError) Error() string {
	return fmt.Sprintf("Login error; error due to: %s", e.Issue.Error())
}

//SignupError singup error wrapper
type SignupError struct {
	Issue error
}

func (e *SignupError) Error() string {
	return fmt.Sprintf("Signup error; error due to: %s", e.Issue.Error())
}

//InvalidTokenError invalid token
type InvalidTokenError struct {
	token string
}

func (e *InvalidTokenError) Error() string {
	return fmt.Sprintf("Invalid token: %s", e.token)
}

//ChatWriteError write to chat error
type ChatWriteError struct {
}

func (e *ChatWriteError) Error() string {
	return fmt.Sprintf("Unable to write to chat")
}

//ChatReadError read from chat error
type ChatReadError struct {
}

func (e *ChatReadError) Error() string {
	return fmt.Sprintf("Unable to read from chat")
}

//UpgradeError read from chat error
type UpgradeError struct {
	Issue error
}

func (e *UpgradeError) Error() string {
	return fmt.Sprintf("Upgrade error; error due to: %s", e.Issue.Error())
}
