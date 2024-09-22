package chat

import "time"

type Message struct {
	User           string    `json:"user"`
	Room           string    `json:"room"`
	IsNotification bool      `json:"isNotification"`
	Value          string    `json:"value"`
	ServerTime     time.Time `json:"serverTime"`
}

const CreateRoomEvent = "create-room"

func NewNotification(user, room, event string) *Message {
	return &Message{
		User:           user,
		Room:           room,
		Value:          event,
		IsNotification: true,
	}
}
