package chat

import (
	"log"

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
			b.logger.Printf("Broadcasting message %v", message)
			for socket := range b.sockets {
				select {
				case socket.outbound <- message:
				default:
					close(socket.outbound)
					delete(b.sockets, socket)
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
