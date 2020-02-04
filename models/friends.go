package models

// Friends structure of friend
type Friends struct {
	FirstName    string   `json:"firstname"`
	LastName     string   `json:"lastname"`
	UserLocation Location `json:"location"`
	UUID         string   `json:"uuid"`
}
