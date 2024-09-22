package domain

type Room struct {
	Name    string
	Creator *User
}

var defaultRoom = &Room{
	Name: "general",
}
