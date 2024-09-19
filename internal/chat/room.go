package chat

type Room struct {
	Name    string
	Creator *User
}

var defaultRoom = &Room{
	Name: "General",
}
