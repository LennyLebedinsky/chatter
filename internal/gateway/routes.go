package gateway

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (g *Gateway) registerRoutes() {
	g.router.HandleFunc("/rooms", g.handleListRooms).Methods(http.MethodGet, http.MethodOptions)
	g.router.HandleFunc("/rooms/{username}", g.handleListRoomsWithUser).Methods(http.MethodGet, http.MethodOptions)
	g.router.HandleFunc("/join-room/{roomname}/{username}", g.handleJoinRoom).Methods(http.MethodPost, http.MethodOptions)
	g.router.HandleFunc("/ws/{username}", g.serveUserWs)
	g.router.Use(g.broadcasterStartMiddleware)
	g.router.Use(g.loggingMiddleware)
	g.router.Use(mux.CORSMethodMiddleware(g.router))
	http.Handle("/", g.router)
}
