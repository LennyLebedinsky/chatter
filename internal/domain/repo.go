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
	Register(userName string) (*User, error)
	Unregister(userName string) error

	Join(userName, roomName string) error
	Leave(userName, roomName string) error

	CreateRoom(roomName, creatorUserName string) (*Room, error)

	ListRooms() ([]*Room, error)
	ListRoomsParticipation() ([]*RoomParticipation, error)
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

func (r *InMemoryRepository) Register(userName string) (*User, error) {
	if r.findUser(userName) != nil {
		return nil, errors.New("user with this name already exists")
	}

	newUser := &User{
		Name: userName,
	}
	r.users[userName] = newUser

	return newUser, nil
}

func (r *InMemoryRepository) Unregister(userName string) error {
	user := r.findUser(userName)
	if user == nil {
		return errors.New("no user registered under this name")
	}

	// Make user leave all joined rooms.
	if joinedRooms, ok := r.userToRooms[user]; ok {
		for _, room := range joinedRooms {
			if err := r.Leave(userName, room.Name); err != nil {
				return err
			}
		}
	}
	delete(r.users, userName)

	return nil
}

func (r *InMemoryRepository) Join(userName, roomName string) error {
	user := r.findUser(userName)
	if user == nil {
		return errors.New("no user registered under this name")
	}
	room := r.findRoom(roomName)
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

func (r *InMemoryRepository) Leave(userName, roomName string) error {
	user := r.findUser(userName)
	if user == nil {
		return errors.New("no user registered under this name")
	}
	room := r.findRoom(roomName)
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
	user := r.findUser(creatorUserName)
	if user == nil {
		return nil, errors.New("no user registered under this name")
	}

	if r.findRoom(roomName) != nil {
		return nil, errors.New("room with this name already exists")
	}

	newRoom := &Room{
		Name:    roomName,
		Creator: user,
	}
	r.rooms[roomName] = newRoom

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

func (r *InMemoryRepository) ListRoomsParticipation() ([]*RoomParticipation, error) {
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

func (r *InMemoryRepository) findUser(userName string) *User {
	if user, ok := r.users[userName]; ok {
		return user
	}
	return nil
}

func (r *InMemoryRepository) findRoom(roomName string) *Room {
	if room, ok := r.rooms[roomName]; ok {
		return room
	}
	return nil
}
