package gateway

import (
	"log"
	"sync/atomic"

	"github.com/gorilla/mux"
	"github.com/lennylebedinsky/chatter/internal/chat"
	"github.com/lennylebedinsky/chatter/internal/domain"
)

type Gateway struct {
	router             *mux.Router
	repo               domain.Repository
	broadcaster        *chat.Broadcaster
	broadcasterStarted atomic.Bool

	logger *log.Logger
}

func New(repo domain.Repository, logger *log.Logger) *Gateway {
	g := &Gateway{
		router:      mux.NewRouter(),
		repo:        repo,
		broadcaster: chat.NewBroadcaster(repo, logger),
		logger:      logger,
	}

	g.broadcasterStarted.Store(false)

	g.registerRoutes()

	return g
}

func (g *Gateway) Router() *mux.Router {
	return g.router
}
