package gateway

import "net/http"

func (g *Gateway) registerRoutes() {
	g.router.HandleFunc("/rooms", g.handleListRooms).Methods(http.MethodGet)
	g.router.HandleFunc("/ws/{username}", g.serveUserWs)
	g.router.Use(g.broadcasterStartMiddleware)
	g.router.Use(g.loggingMiddleware)
	http.Handle("/", g.router)
}
