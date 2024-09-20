package chat

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/lennylebedinsky/chatter/internal/domain"
)

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

func (s *UserSocket) Send() chan []byte {
	return s.send
}

// Supposed to be run as goroutine.
func (s *UserSocket) ReadLoop() {
	defer func() {
		s.broadcaster.unregister <- s
		s.conn.Close()
	}()
	for {
		_, _, err := s.conn.ReadMessage()
		if err != nil {
			if closeErr, ok := err.(*websocket.CloseError); ok {
				s.logger.Printf("Connection closed for user %s: %v\n", s.user.Name, closeErr)
			} else {
				s.logger.Printf("Error reading message for user %s: %v\n", s.user.Name, err)
			}
			break
		}
		//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		//c.hub.broadcast <- message
	}
}
