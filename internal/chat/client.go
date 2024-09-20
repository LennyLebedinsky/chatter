package chat

import (
	"github.com/gorilla/websocket"
	"github.com/lennylebedinsky/chatter/internal/domain"
)

type Client interface {
	ID() string
	Send() chan []byte
}

type UserSocket struct {
	user *domain.User
	conn *websocket.Conn

	send chan []byte
}

func NewUserSocket(user *domain.User, conn *websocket.Conn) *UserSocket {
	return &UserSocket{
		user: user,
		conn: conn,
		send: make(chan []byte),
	}
}

func (s *UserSocket) ID() string {
	if s.user != nil {
		return s.user.Name
	}

	return ""
}

func (s *UserSocket) Send() chan []byte {
	return s.send
}
