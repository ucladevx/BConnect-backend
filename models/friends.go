package models

// Friends structure of friend
type Friends struct {
	UUID  string `json:"uuid"`
	FUUID string `json:"fuuid"`
}

// FriendRequest struct to hold content of friend requests
type FriendRequest struct {
	Sender    string `json:"sender"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}
