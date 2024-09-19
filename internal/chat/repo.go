package chat

import (
	"errors"
	"slices"
)

type Repository interface {
	Register(userName string) (*User, error)
	Unregister(userName string) error

	Join(userName, roomName string) error
	Leave(userName, roomName string) error

	CreateRoom(roomName, creatorUserName string) (*Room, error)
	ListRooms() ([]*Room, error)
}

type InMemoryRepository struct {
	users map[string]*User
	rooms map[string]*Room

	userInRooms   map[*User][]*Room
	roomWithUsers map[*Room][]*User
}

func NewInMemoryRepository() Repository {
	r := &InMemoryRepository{
		users: make(map[string]*User),
		rooms: make(map[string]*Room),

		userInRooms:   make(map[*User][]*Room),
		roomWithUsers: make(map[*Room][]*User),
	}

	r.rooms[defaultRoom.Name] = defaultRoom

	return r
}

func (r *InMemoryRepository) Register(userName string) (*User, error) {
	if _, ok := r.users[userName]; ok {
		return nil, errors.New("user with this name already exists")
	}

	newUser := &User{
		Name: userName,
	}
	r.users[userName] = newUser

	return newUser, nil
}

func (r *InMemoryRepository) Unregister(userName string) error {
	var user *User
	var ok bool
	if user, ok = r.users[userName]; !ok {
		return errors.New("no user registered under this name")
	}

	// Make user leave all joined rooms.
	if joinedRooms, ok := r.userInRooms[user]; ok {
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
	var user *User
	var room *Room
	var ok bool
	if user, ok = r.users[userName]; !ok {
		return errors.New("no user registered under this name")
	}

	if room, ok = r.rooms[roomName]; !ok {
		return errors.New("no room with this name exists")
	}

	// Update indexes.
	if _, ok := r.userInRooms[user]; ok {
		if slices.Index(r.userInRooms[user], room) < 0 {
			r.userInRooms[user] = append(r.userInRooms[user], room)
		}
	} else {
		r.userInRooms[user] = []*Room{room}
	}

	if _, ok := r.roomWithUsers[room]; ok {
		if slices.Index(r.roomWithUsers[room], user) < 0 {
			r.roomWithUsers[room] = append(r.roomWithUsers[room], user)
		}
	} else {
		r.roomWithUsers[room] = []*User{user}
	}

	return nil
}

func (r *InMemoryRepository) Leave(userName, roomName string) error {
	var user *User
	var room *Room
	var ok bool
	if user, ok = r.users[userName]; !ok {
		return errors.New("no user registered under this name")
	}

	if room, ok = r.rooms[roomName]; !ok {
		return errors.New("no room with this name exists")
	}

	// Update indexes.
	if _, ok := r.userInRooms[user]; ok {
		index := slices.Index(r.userInRooms[user], room)
		if index >= 0 {
			r.userInRooms[user] = slices.Delete(r.userInRooms[user], index, index+1)
		}
	}

	if _, ok := r.roomWithUsers[room]; ok {
		index := slices.Index(r.roomWithUsers[room], user)
		if index >= 0 {
			r.roomWithUsers[room] = slices.Delete(r.roomWithUsers[room], index, index+1)
		}
	}

	return nil
}

func (r *InMemoryRepository) CreateRoom(roomName, creatorUserName string) (*Room, error) {
	var user *User
	var ok bool
	if user, ok = r.users[creatorUserName]; !ok {
		return nil, errors.New("no user registered under this name")
	}

	if _, ok = r.rooms[roomName]; ok {
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
