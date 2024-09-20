package gateway

import (
	"log"

	"github.com/gorilla/mux"
	"github.com/lennylebedinsky/chatter/internal/chat"
	"github.com/lennylebedinsky/chatter/internal/domain"
)

type Gateway struct {
	repo domain.Repository

	router *mux.Router
	logger *log.Logger

	broadcaster *chat.Broadcaster
}

func New(broadcaster *chat.Broadcaster, repo domain.Repository, logger *log.Logger) *Gateway {
	g := &Gateway{
		router:      mux.NewRouter(),
		broadcaster: broadcaster,
		repo:        repo,
		logger:      logger,
	}

	g.registerRoutes()

	return g
}

func (g *Gateway) Router() *mux.Router {
	return g.router
}

func (g *Gateway) Broadcaster() *chat.Broadcaster {
	return g.broadcaster
}
