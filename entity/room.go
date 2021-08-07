package entity

type User struct {
	Name string `json:"name"`
	Hand *bool  `json:"hand,omitempty"`
}

type RoomData struct {
	Status string          `json:"status"`
	Owner  string          `json:"owner"`
	Users  map[string]User `json:"users"`
}
