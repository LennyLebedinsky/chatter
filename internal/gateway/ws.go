package gateway

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/lennylebedinsky/chatter/internal/chat"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// NB: Just for example allowing local clients to reach server, should take precaution in real environment.
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// serveUserWs initiates Websocket upgrade when user logs in to chat server.
// Socket is registered with the broadcaster, and read and write loops are started.
func (g *Gateway) serveUserWs(w http.ResponseWriter, r *http.Request) {
	userName := strings.ToLower(mux.Vars(r)["username"])
	user := g.repo.FindUser(r.Context(), userName)
	var err error
	if user == nil {
		user, err = g.repo.CreateUser(r.Context(), userName)
		if err != nil {
			g.logError(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if g.broadcaster.IsRegistered(user) {
		err = fmt.Errorf("User %s is already registered with broadcaster.", user.Name)
		g.logError(err)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		g.logError(err)
		return
	}
	userSocket := chat.NewUserSocket(user, conn, g.broadcaster, g.logger)
	g.broadcaster.Register() <- userSocket

	go userSocket.ReadLoop()
	go userSocket.WriteLoop()
}
