package gateway

import (
	"context"
	"log"
	"sync/atomic"

	"github.com/gorilla/mux"
	"github.com/lennylebedinsky/chatter/internal/chat"
	"github.com/lennylebedinsky/chatter/internal/domain"
	"github.com/lennylebedinsky/chatter/internal/message"
)

// Gateway orchestrates HTTP and Websocket communications and data persistence.
type Gateway struct {
	router             *mux.Router
	repo               domain.Repository
	messageStore       message.Store
	broadcaster        *chat.Broadcaster
	broadcasterStarted atomic.Bool

	logger *log.Logger
}

func New(repo domain.Repository, messageStore message.Store, logger *log.Logger) *Gateway {
	g := &Gateway{
		router:       mux.NewRouter(),
		repo:         repo,
		messageStore: messageStore,
		broadcaster:  chat.NewBroadcaster(repo, messageStore, logger),
		logger:       logger,
	}

	g.broadcasterStarted.Store(false)

	g.registerRoutes()

	return g
}

func (g *Gateway) StartBroadcaster(ctx context.Context) {
	// Just one broadcaster goroutine should run for the gateway.
	if g.broadcasterStarted.CompareAndSwap(false, true) {
		go g.broadcaster.Start(ctx)
	}
}

func (g *Gateway) Router() *mux.Router {
	return g.router
}
