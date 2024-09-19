package gateway

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Gateway struct {
	router *mux.Router
	logger *log.Logger
}

func New(logger *log.Logger) *Gateway {
	g := &Gateway{
		router: mux.NewRouter(),
		logger: logger,
	}

	g.registerRoutes()

	return g
}

func (g *Gateway) Router() *mux.Router {
	return g.router
}

func (g *Gateway) registerRoutes() {
	//	s.router.HandleFunc("/hello", s.handleHelloWorld).Queries("name", "{name}")
	g.router.Use(g.loggingMiddleware)
	http.Handle("/", g.router)
}

func (g *Gateway) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Logging.
		g.logger.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
