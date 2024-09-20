package chat

import "log"

type Broadcaster struct {
	clients map[Client]bool

	register   chan Client
	unregister chan Client

	message chan []byte

	logger *log.Logger
}

func NewBroadcaster(logger *log.Logger) *Broadcaster {
	return &Broadcaster{
		clients:    make(map[Client]bool),
		register:   make(chan Client),
		unregister: make(chan Client),
		message:    make(chan []byte),
		logger:     logger,
	}
}

// Supposed to be run as goroutine.
func (b *Broadcaster) Start() {
	for {
		select {
		case client := <-b.register:
			b.clients[client] = true
			b.logger.Printf("Client %s registered with broadcaster.\n", client.ID())
		case client := <-b.unregister:
			if _, ok := b.clients[client]; ok {
				delete(b.clients, client)
				close(client.Send())
			}

		case message := <-b.message:
			for client := range b.clients {
				select {
				case client.Send() <- message:
				default:
					close(client.Send())
					delete(b.clients, client)
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

func (b *Broadcaster) Register() chan Client {
	return b.register
}

func (b *Broadcaster) Broadcast() chan []byte {
	return b.message
}
