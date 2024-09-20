package chat

import (
	"log"

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

	broadcaster *Broadcaster

	send chan []byte

	logger *log.Logger
}

func NewUserSocket(
	user *domain.User,
	conn *websocket.Conn,
	broadcaster *Broadcaster,
	logger *log.Logger) *UserSocket {
	return &UserSocket{
		user:        user,
		conn:        conn,
		broadcaster: broadcaster,
		send:        make(chan []byte),
		logger:      logger,
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

func (s *UserSocket) ReadLoop() {
	defer func() {
		s.broadcaster.Unregister() <- s
		s.conn.Close()
	}()
	for {
		_, _, err := s.conn.ReadMessage()
		if err != nil {
			if closeErr, ok := err.(*websocket.CloseError); ok {
				s.logger.Printf("Connection closed for client %s: %v\n", s.ID(), closeErr)
			} else {
				s.logger.Printf("Error reading message for client %s: %v\n", s.ID(), err)
			}
			break
		}
		//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		//c.hub.broadcast <- message
	}
}
