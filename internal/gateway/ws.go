package gateway

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/lennylebedinsky/chatter/internal/chat"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (g *Gateway) serveUserWs(w http.ResponseWriter, r *http.Request) {
	userName := strings.ToLower(mux.Vars(r)["username"])
	user, err := g.repo.RegisterUser(userName)
	if err != nil {
		g.logError(err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		g.logError(err)
		err = g.repo.UnregisterUser(user.Name)
		if err != nil {
			g.logError(err)
		}
		return
	}
	userSocket := chat.NewUserSocket(user, conn, g.broadcaster, g.logger)
	g.broadcaster.Register() <- userSocket

	go userSocket.ReadLoop()
	go userSocket.WriteLoop()
}
