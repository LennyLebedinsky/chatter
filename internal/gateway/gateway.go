package gateway

import (
	"log"

	"github.com/gorilla/mux"
	"github.com/lennylebedinsky/chatter/internal/chat"
)

type Gateway struct {
	repo chat.Repository

	router *mux.Router
	logger *log.Logger
}

func New(logger *log.Logger) *Gateway {
	g := &Gateway{
		repo:   chat.NewInMemoryRepository(),
		router: mux.NewRouter(),
		logger: logger,
	}

	g.registerRoutes()

	return g
}

func (g *Gateway) Router() *mux.Router {
	return g.router
}
