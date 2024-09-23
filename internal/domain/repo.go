package domain

import (
	"context"
	"errors"
	"slices"
	"sync"
)

type RoomParticipation struct {
	Room         *Room
	Participants []*User
}

// Repository stores and retrieves relations between users and rooms.
// TODO: implement as a persistent storage, preferrably Redis cache plus SQL server on background.
type Repository interface {
	CreateUser(ctx context.Context, userName string) (*User, error)
	FindUser(ctx context.Context, userName string) *User

	CreateRoom(ctx context.Context, roomName, creatorUserName string) (*Room, error)
	FindRoom(ctx context.Context, roomName string) *Room
	JoinRoom(ctx context.Context, userName, roomName string) error
	LeaveRoom(ctx context.Context, userName, roomName string) error

	ListRooms(ctx context.Context) ([]*Room, error)
	ListParticipants(ctx context.Context, roomName string) ([]*User, error)
	ListParticipantsForAllRooms(ctx context.Context) ([]*RoomParticipation, error)
}

type InMemoryRepository struct {
	users map[string]*User
	rooms map[string]*Room

	userToRooms map[*User][]*Room
	roomToUsers map[*Room][]*User

	mu sync.RWMutex
}

func NewInMemoryRepository() Repository {
	r := &InMemoryRepository{
		users: make(map[string]*User),
		rooms: make(map[string]*Room),

		userToRooms: make(map[*User][]*Room),
		roomToUsers: make(map[*Room][]*User),
	}

	r.rooms[defaultRoom.Name] = defaultRoom

	return r
}

func (r *InMemoryRepository) CreateUser(ctx context.Context, userName string) (*User, error) {
	if r.FindUser(ctx, userName) != nil {
		return nil, errors.New("user with this name already exists")
	}

	newUser := &User{
		Name: userName,
	}
	r.mu.Lock()
	r.users[userName] = newUser
	r.mu.Unlock()

	return newUser, nil
}

func (r *InMemoryRepository) FindUser(_ context.Context, userName string) *User {
	if user, ok := r.users[userName]; ok {
		return user
	}
	return nil
}

func (r *InMemoryRepository) FindRoom(_ context.Context, roomName string) *Room {
	if room, ok := r.rooms[roomName]; ok {
		return room
	}
	return nil
}

func (r *InMemoryRepository) JoinRoom(ctx context.Context, userName, roomName string) error {
	user := r.FindUser(ctx, userName)
	if user == nil {
		return errors.New("no user registered under this name")
	}
	room := r.FindRoom(ctx, roomName)
	if room == nil {
		return errors.New("no room with this name exists")
	}

	// Update indexes.
	r.mu.Lock()
	if _, ok := r.userToRooms[user]; ok {
		if slices.Index(r.userToRooms[user], room) < 0 {
			r.userToRooms[user] = append(r.userToRooms[user], room)
		}
	} else {
		r.userToRooms[user] = []*Room{room}
	}

	if _, ok := r.roomToUsers[room]; ok {
		if slices.Index(r.roomToUsers[room], user) < 0 {
			r.roomToUsers[room] = append(r.roomToUsers[room], user)
		}
	} else {
		r.roomToUsers[room] = []*User{user}
	}
	r.mu.Unlock()

	return nil
}

func (r *InMemoryRepository) LeaveRoom(ctx context.Context, userName, roomName string) error {
	user := r.FindUser(ctx, userName)
	if user == nil {
		return errors.New("no user registered under this name")
	}
	room := r.FindRoom(ctx, roomName)
	if room == nil {
		return errors.New("no room with this name exists")
	}

	// Update indexes.
	r.mu.Lock()
	if _, ok := r.userToRooms[user]; ok {
		index := slices.Index(r.userToRooms[user], room)
		if index >= 0 {
			r.userToRooms[user] = slices.Delete(r.userToRooms[user], index, index+1)
		}
	}

	if _, ok := r.roomToUsers[room]; ok {
		index := slices.Index(r.roomToUsers[room], user)
		if index >= 0 {
			r.roomToUsers[room] = slices.Delete(r.roomToUsers[room], index, index+1)
		}
	}
	r.mu.Unlock()

	return nil
}

func (r *InMemoryRepository) CreateRoom(ctx context.Context, roomName, creatorUserName string) (*Room, error) {
	creatorUser := r.FindUser(ctx, creatorUserName)
	if creatorUser == nil {
		return nil, errors.New("no user registered under this name")
	}

	if r.FindRoom(ctx, roomName) != nil {
		return nil, errors.New("room with this name already exists")
	}

	newRoom := &Room{
		Name:    roomName,
		Creator: creatorUser,
	}

	r.mu.Lock()
	r.rooms[roomName] = newRoom
	r.mu.Unlock()

	// User who is creating room automatically joins it.
	if err := r.JoinRoom(ctx, creatorUserName, roomName); err != nil {
		return nil, err
	}

	return newRoom, nil
}

func (r *InMemoryRepository) ListRooms(_ context.Context) ([]*Room, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rooms := make([]*Room, len(r.rooms))
	i := 0
	for _, room := range r.rooms {
		rooms[i] = room
		i++
	}
	return rooms, nil
}

func (r *InMemoryRepository) ListParticipants(ctx context.Context, roomName string) ([]*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	room := r.FindRoom(ctx, roomName)
	if room == nil {
		return nil, errors.New("no room with this name exists")
	}
	return r.roomToUsers[room], nil
}

func (r *InMemoryRepository) ListParticipantsForAllRooms(_ context.Context) ([]*RoomParticipation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	roomsParticipation := make([]*RoomParticipation, len(r.rooms))
	i := 0
	for _, room := range r.rooms {
		roomsParticipation[i] = &RoomParticipation{
			Room:         room,
			Participants: r.roomToUsers[room],
		}
		i++
	}
	return roomsParticipation, nil
}
