package gateway

import (
	"net/http"
)

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
