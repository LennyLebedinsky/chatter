package chat

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/lennylebedinsky/chatter/internal/domain"
	"github.com/lennylebedinsky/chatter/internal/message"
)

type Broadcaster struct {
	sockets map[*UserSocket]bool

	register   chan *UserSocket
	unregister chan *UserSocket

	message chan *message.Message

	repo         domain.Repository
	messageStore message.Store

	logger *log.Logger
}

func NewBroadcaster(repo domain.Repository, messageStore message.Store, logger *log.Logger) *Broadcaster {
	return &Broadcaster{
		sockets:      make(map[*UserSocket]bool),
		register:     make(chan *UserSocket),
		unregister:   make(chan *UserSocket),
		message:      make(chan *message.Message),
		repo:         repo,
		messageStore: messageStore,
		logger:       logger,
	}

}

// Start listens to messages coming from user sockets and dispatches them.
// It is supposed to run as goroutine.
// Only one broadcaster runs for the whole service.
func (b *Broadcaster) Start(ctx context.Context) {
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
			}
		case msg := <-b.message:
			if err := b.validate(ctx, msg); err != nil {
				// Not fatal, just log and continue listening for other messages.
				b.logger.Printf("Message is not accepted by broadcaster: %v\n", err)
			} else {
				b.accept(ctx, msg)
				b.logger.Printf("Broadcasting message %v", msg)
				destination, err := b.dispatch(ctx, msg)
				if err == nil {
					for _, socket := range destination {
						select {
						case socket.outbound <- msg:
						default:
							// If send buffer is full, assume client is disconnected or hanged.
							close(socket.outbound)
							delete(b.sockets, socket)
						}
					}
				}
			}
		case <-ctx.Done():
			b.logger.Printf("Context canceled, stopping broadcaster...")
			return
		}
	}

}

func (b *Broadcaster) Register() chan *UserSocket {
	return b.register
}

func (b *Broadcaster) Message() chan *message.Message {
	return b.message
}

func (b *Broadcaster) IsRegistered(user *domain.User) bool {
	for socket := range b.sockets {
		if socket.user == user {
			return true
		}
	}

	return false
}

// validate checks if message is considered valid for broadcasting.
func (b *Broadcaster) validate(_ context.Context, msg *message.Message) error {
	// Notifications potentially could have user or room missed.
	if msg.IsNotification {
		return nil
	}

	if msg.User == "" {
		return errors.New("message does not have an author")
	}

	if msg.Room == "" {
		return errors.New("message does not have room destination")
	}

	return nil
}

// accept marks that message is allowed into system.
func (b *Broadcaster) accept(ctx context.Context, msg *message.Message) {
	// Setup server timestamp.
	msg.ServerTime = time.Now()
	// TODO: assign unique ID, possibly logical clock.
	// Add message to persistent storage.
	if err := b.messageStore.SaveMessage(ctx, msg.Room, msg); err != nil {
		// Not fatal, just continue without message retention.
		b.logger.Printf("Message could not be stored: %v\n", err)
	}
}

// dispatch determines only those users to whom message will be broadcasted.
func (b *Broadcaster) dispatch(ctx context.Context, msg *message.Message) ([]*UserSocket, error) {
	sockets := []*UserSocket{}

	// Notifications are going to everyone.
	if msg.IsNotification {
		for socket := range b.sockets {
			sockets = append(sockets, socket)
		}
		return sockets, nil
	}

	// Main rule for this chat: message is broadcasted only to users who joined the same room.
	usersInSameRoom, err := b.repo.ListParticipants(ctx, msg.Room)
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
