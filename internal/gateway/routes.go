package gateway

import "net/http"

func (g *Gateway) registerRoutes() {
	g.router.HandleFunc("/rooms", g.handleListRooms).Methods("GET")
	g.router.Use(g.loggingMiddleware)
	http.Handle("/", g.router)
}

func (g *Gateway) handleListRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := g.repo.ListRooms()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = encode(w, r, http.StatusOK, rooms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (g *Gateway) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		g.logger.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
