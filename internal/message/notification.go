package message

const CreateRoomEvent = "create-room"

func NewNotification(user, room, event string) *Message {
	return &Message{
		User:           user,
		Room:           room,
		Value:          event,
		IsNotification: true,
	}
}
