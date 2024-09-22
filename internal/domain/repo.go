package domain

import (
	"errors"
	"slices"
)

type RoomParticipation struct {
	Room         *Room
	Participants []*User
}

type Repository interface {
	CreateUser(userName string) (*User, error)
	FindUser(userName string) *User

	CreateRoom(roomName, creatorUserName string) (*Room, error)
	FindRoom(roomName string) *Room
	JoinRoom(userName, roomName string) error
	LeaveRoom(userName, roomName string) error

	ListRooms() ([]*Room, error)
	ListParticipants(roomName string) ([]*User, error)
	ListParticipantsForAllRooms() ([]*RoomParticipation, error)
}

type InMemoryRepository struct {
	users map[string]*User
	rooms map[string]*Room

	userToRooms map[*User][]*Room
	roomToUsers map[*Room][]*User
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

func (r *InMemoryRepository) CreateUser(userName string) (*User, error) {
	if r.FindUser(userName) != nil {
		return nil, errors.New("user with this name already exists")
	}

	newUser := &User{
		Name: userName,
	}
	r.users[userName] = newUser

	return newUser, nil
}

func (r *InMemoryRepository) UserExists(userName string) bool {
	return r.FindUser(userName) != nil
}

func (r *InMemoryRepository) JoinRoom(userName, roomName string) error {
	user := r.FindUser(userName)
	if user == nil {
		return errors.New("no user registered under this name")
	}
	room := r.FindRoom(roomName)
	if room == nil {
		return errors.New("no room with this name exists")
	}

	// Update indexes.
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

	return nil
}

func (r *InMemoryRepository) LeaveRoom(userName, roomName string) error {
	user := r.FindUser(userName)
	if user == nil {
		return errors.New("no user registered under this name")
	}
	room := r.FindRoom(roomName)
	if room == nil {
		return errors.New("no room with this name exists")
	}

	// Update indexes.
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

	return nil
}

func (r *InMemoryRepository) CreateRoom(roomName, creatorUserName string) (*Room, error) {
	creatorUser := r.FindUser(creatorUserName)
	if creatorUser == nil {
		return nil, errors.New("no user registered under this name")
	}

	if r.FindRoom(roomName) != nil {
		return nil, errors.New("room with this name already exists")
	}

	newRoom := &Room{
		Name:    roomName,
		Creator: creatorUser,
	}
	r.rooms[roomName] = newRoom

	// User who is creating room automatically joins it.
	if err := r.JoinRoom(creatorUserName, roomName); err != nil {
		return nil, err
	}

	return newRoom, nil
}

func (r *InMemoryRepository) ListRooms() ([]*Room, error) {
	rooms := make([]*Room, len(r.rooms))
	i := 0
	for _, room := range r.rooms {
		rooms[i] = room
		i++
	}
	return rooms, nil
}

func (r *InMemoryRepository) ListParticipants(roomName string) ([]*User, error) {
	room := r.FindRoom(roomName)
	if room == nil {
		return nil, errors.New("no room with this name exists")
	}
	return r.roomToUsers[room], nil
}

func (r *InMemoryRepository) ListParticipantsForAllRooms() ([]*RoomParticipation, error) {
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

func (r *InMemoryRepository) FindUser(userName string) *User {
	if user, ok := r.users[userName]; ok {
		return user
	}
	return nil
}

func (r *InMemoryRepository) FindRoom(roomName string) *Room {
	if room, ok := r.rooms[roomName]; ok {
		return room
	}
	return nil
}
