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

	outbound chan *Message

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
		outbound:    make(chan *Message),
		logger:      logger,
	}
}

// Supposed to be run as goroutine.
func (s *UserSocket) ReadLoop() {
	defer func() {
		s.broadcaster.unregister <- s
		s.conn.Close()
	}()
	for {
		message := &Message{}
		err := s.conn.ReadJSON(message)
		if err != nil {
			if closeErr, ok := err.(*websocket.CloseError); ok {
				s.logger.Printf("Connection closed for user %s: %v\n", s.user.Name, closeErr)
				return
			} else {
				s.logger.Printf("Error reading message for user %s: %v\n", s.user.Name, err)
			}
		}
		s.logger.Printf("Received message: %v\n", message)
		s.broadcaster.message <- message
	}
}

func (s *UserSocket) WriteLoop() {
	defer func() {
		s.conn.Close()
	}()
	for {
		select {
		case message, ok := <-s.outbound:
			if !ok {
				s.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := s.conn.WriteJSON(message)
			if err != nil {
				if closeErr, ok := err.(*websocket.CloseError); ok {
					s.logger.Printf("Connection closed for user %s: %v\n", s.user.Name, closeErr)
					return
				} else {
					s.logger.Printf("Error writing message for user %s: %v\n", s.user.Name, err)
				}
			}
			s.logger.Printf("Sent message %v to user %s\n", message, s.user.Name)
		}
	}
}
