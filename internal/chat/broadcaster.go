package chat

import (
	"log"

	"github.com/lennylebedinsky/chatter/internal/domain"
)

type Broadcaster struct {
	sockets map[*UserSocket]bool

	register   chan *UserSocket
	unregister chan *UserSocket

	message chan []byte

	repo domain.Repository

	logger *log.Logger
}

func NewBroadcaster(repo domain.Repository, logger *log.Logger) *Broadcaster {
	return &Broadcaster{
		sockets:    make(map[*UserSocket]bool),
		register:   make(chan *UserSocket),
		unregister: make(chan *UserSocket),
		message:    make(chan []byte),
		repo:       repo,
		logger:     logger,
	}
}

// Supposed to be run as goroutine.
func (b *Broadcaster) Start() {
	for {
		select {
		case socket := <-b.register:
			b.sockets[socket] = true
			b.logger.Printf("User %s registered with broadcaster.\n", socket.ID())
		case socket := <-b.unregister:
			if _, ok := b.sockets[socket]; ok {
				delete(b.sockets, socket)
				close(socket.Send())
				b.logger.Printf("User %s unregistered from broadcaster.\n", socket.ID())
				if err := b.repo.Unregister(socket.user.Name); err != nil {
					b.logger.Printf("Error: %v\n", err)
				}
			}

		case message := <-b.message:
			for socket := range b.sockets {
				select {
				case socket.Send() <- message:
				default:
					close(socket.Send())
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

func (b *Broadcaster) Unregister() chan *UserSocket {
	return b.unregister
}

func (b *Broadcaster) Broadcast() chan []byte {
	return b.message
}
