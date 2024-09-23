package domain

// User represents chat room for users to exchange messages.
// TODO: Add unique ID (e.g., GUID).
type Room struct {
	Name    string
	Creator *User
}

var defaultRoom = &Room{
	Name: "general",
}
