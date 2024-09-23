package chat

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/lennylebedinsky/chatter/internal/domain"
	"github.com/lennylebedinsky/chatter/internal/message"
)

type UserSocket struct {
	user *domain.User
	conn *websocket.Conn

	broadcaster *Broadcaster

	outbound chan *message.Message

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
		outbound:    make(chan *message.Message),
		logger:      logger,
	}
}

// ReadLoop listens to messages coming from client's side of Websocket connection
// and redirects them to broadcaster.
// It is supposed to run as goroutine, one read loop per client.
func (s *UserSocket) ReadLoop() {
	defer func() {
		s.broadcaster.unregister <- s
		s.conn.Close()
	}()
	for {
		msg := &message.Message{}
		err := s.conn.ReadJSON(msg)
		if err != nil {
			// If connection had been closed from client's side, break the loop.
			if closeErr, ok := err.(*websocket.CloseError); ok {
				s.logger.Printf("Connection closed for user %s: %v\n", s.user.Name, closeErr)
				return
			} else {
				s.logger.Printf("Error reading message for user %s: %v\n", s.user.Name, err)
			}
		}
		s.logger.Printf("Received message: %v\n", msg)
		s.broadcaster.message <- msg
	}
}

// Write listens to messages coming from broadcaster and
// redirects them  client's side of Websocket connection.
// It is supposed to run as goroutine, one read loop per client.
func (s *UserSocket) WriteLoop() {
	defer func() {
		s.conn.Close()
	}()
	for {
		select {
		case message, ok := <-s.outbound:
			// If broadcaster closed channel from its side, initiate closing handshake.
			if !ok {
				s.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := s.conn.WriteJSON(message)
			if err != nil {
				// If connection had been closed from client's side, break the loop.
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
