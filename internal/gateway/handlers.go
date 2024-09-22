package gateway

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gorilla/mux"
	"github.com/lennylebedinsky/chatter/internal/chat"
	"github.com/lennylebedinsky/chatter/internal/domain"
)

func (g *Gateway) handleListRooms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

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

func (g *Gateway) handleListRoomsWithUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	userName := strings.ToLower(mux.Vars(r)["username"])
	roomsParticipation, err := g.repo.ListParticipantsForAllRooms()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type roomWithUser struct {
		Room              *domain.Room
		UserIsParticipant bool
	}

	response := make([]*roomWithUser, len(roomsParticipation))
	for i, roomParticipation := range roomsParticipation {
		response[i] = &roomWithUser{
			Room: roomParticipation.Room,
			UserIsParticipant: slices.IndexFunc(roomParticipation.Participants, func(u *domain.User) bool {
				return u.Name == userName
			}) >= 0,
		}
	}

	err = encode(w, r, http.StatusOK, response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (g *Gateway) handleJoinRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	roomName := strings.ToLower(mux.Vars(r)["roomname"])
	userName := strings.ToLower(mux.Vars(r)["username"])
	if err := g.repo.JoinRoom(userName, roomName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	g.logger.Printf("User %s joined room %s.\n", userName, roomName)
}

func (g *Gateway) handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	roomName := strings.ToLower(mux.Vars(r)["roomname"])
	userName := strings.ToLower(mux.Vars(r)["username"])
	if _, err := g.repo.CreateRoom(roomName, userName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	g.logger.Printf("User %s created and joined room %s.\n", userName, roomName)

	// Notify other clients about room creation so that they could update room list.
	g.broadcaster.Message() <- chat.NewNotification(userName, roomName, chat.CreateRoomEvent)
}
