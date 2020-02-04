package actions

// ActionService abstract server-side authentication in case we switch from whatever current auth scheme we are using
type ActionService interface {
	GET(key string, request string) (bool, error)
	SET(key string, request string) error
	PUT(key string, request string) error
	DEL(key string, request string) error
}
