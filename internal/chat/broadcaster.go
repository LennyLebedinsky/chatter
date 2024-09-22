package chat

import (
	"errors"
	"log"
	"time"

	"github.com/lennylebedinsky/chatter/internal/domain"
)

type Broadcaster struct {
	sockets map[*UserSocket]bool

	register   chan *UserSocket
	unregister chan *UserSocket

	message chan *Message

	repo domain.Repository

	logger *log.Logger
}

func NewBroadcaster(repo domain.Repository, logger *log.Logger) *Broadcaster {
	return &Broadcaster{
		sockets:    make(map[*UserSocket]bool),
		register:   make(chan *UserSocket),
		unregister: make(chan *UserSocket),
		message:    make(chan *Message),
		repo:       repo,
		logger:     logger,
	}

}

// Supposed to be run as goroutine.
func (b *Broadcaster) Start() {
	defer func() {
		b.logger.Println("Message broadcaster stopped.")
	}()

	b.logger.Println("Message broadcaster started.")
	for {
		select {
		case socket := <-b.register:
			b.sockets[socket] = true
			b.logger.Printf("User %s registered with broadcaster.\n", socket.user.Name)
		case socket := <-b.unregister:
			if _, ok := b.sockets[socket]; ok {
				delete(b.sockets, socket)
				close(socket.outbound)
				b.logger.Printf("User %s unregistered from broadcaster.\n", socket.user.Name)
				if err := b.repo.Unregister(socket.user.Name); err != nil {
					b.logger.Printf("Error: %v\n", err)
				}
			}
		case message := <-b.message:
			if err := b.validate(message); err != nil {
				b.logger.Printf("Message is not accepted by broadcaster: %v\n", err)
			} else {
				b.accept(message)
				b.logger.Printf("Broadcasting message %v", message)
				destination, err := b.dispatch(message)
				if err == nil {
					for _, socket := range destination {
						select {
						case socket.outbound <- message:
						default:
							close(socket.outbound)
							delete(b.sockets, socket)
						}
					}
				}
			}
			/*
				case _, ok := <-stop:
					if !ok {
						b.logger.Println("Stopping broadcaster...")
						return
					}*/
		}
	}

}

func (b *Broadcaster) Register() chan *UserSocket {
	return b.register
}

// validate checks if message is considered valid for broadcasting.
func (b *Broadcaster) validate(message *Message) error {
	if message.User == "" {
		return errors.New("message does not have an author")
	}

	if message.Room == "" {
		return errors.New("message does not have room destination")
	}

	return nil
}

// accept marks that message is allowed into system.
func (b *Broadcaster) accept(message *Message) {
	// Setup server timestamp.
	message.ServerTime = time.Now()
	// TODO: assign unique ID, possibly logical clock.
	// TODO: add message to persistent storage.
}

// dispatch determines only those users to whom message will be broadcasted.
func (b *Broadcaster) dispatch(message *Message) ([]*UserSocket, error) {
	sockets := []*UserSocket{}

	// Main rule for this chat: message is broadcasted only to users who joined the same room.
	usersInSameRoom, err := b.repo.ListParticipants(message.Room)
	if err != nil {
		return nil, err
	}
	index := map[*domain.User]bool{}
	for _, user := range usersInSameRoom {
		index[user] = true
	}

	// Dispatch message to active users participating in the same room.
	for socket := range b.sockets {
		if _, ok := index[socket.user]; ok {
			sockets = append(sockets, socket)
		}
	}
	return sockets, nil
}
