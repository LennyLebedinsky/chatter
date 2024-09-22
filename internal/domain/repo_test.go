package domain

import (
	"reflect"
	"testing"
)

var testUserA = &User{
	Name: "Jarvis",
}

var testUserB = &User{
	Name: "Ultron",
}

var testRoomA = &Room{
	Name: "Sokovia",
}

var testRoomB = &Room{
	Name: "AvengersHub",
}

func TestInMemoryRepository_Register(t *testing.T) {
	type fields struct {
		users         map[string]*User
		rooms         map[string]*Room
		userInRooms   map[*User][]*Room
		roomWithUsers map[*Room][]*User
	}
	type args struct {
		userName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *User
		wantErr bool
	}{
		{
			name: "Registering new user should succeed",
			fields: fields{
				users:         map[string]*User{},
				rooms:         map[string]*Room{},
				userInRooms:   map[*User][]*Room{},
				roomWithUsers: map[*Room][]*User{},
			},
			args: args{
				userName: testUserA.Name,
			},
			want:    testUserA,
			wantErr: false,
		},
		{
			name: "Registering user with the same name as already registered should fail",
			fields: fields{
				users: map[string]*User{
					testUserA.Name: testUserA,
				},
				rooms:         map[string]*Room{},
				userInRooms:   map[*User][]*Room{},
				roomWithUsers: map[*Room][]*User{},
			},
			args: args{
				userName: testUserA.Name,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Registering another new user should succeed",
			fields: fields{
				users: map[string]*User{
					testUserA.Name: testUserA,
				},
				rooms:         map[string]*Room{},
				userInRooms:   map[*User][]*Room{},
				roomWithUsers: map[*Room][]*User{},
			},
			args: args{
				userName: testUserB.Name,
			},
			want:    testUserB,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InMemoryRepository{
				users:       tt.fields.users,
				rooms:       tt.fields.rooms,
				userToRooms: tt.fields.userInRooms,
				roomToUsers: tt.fields.roomWithUsers,
			}
			got, err := r.CreateUser(tt.args.userName)
			if (err != nil) != tt.wantErr {
				t.Errorf("InMemoryRepository.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InMemoryRepository.Register() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemoryRepository_Join(t *testing.T) {
	type fields struct {
		users         map[string]*User
		rooms         map[string]*Room
		userInRooms   map[*User][]*Room
		roomWithUsers map[*Room][]*User
	}
	type args struct {
		userName string
		roomName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InMemoryRepository{
				users:       tt.fields.users,
				rooms:       tt.fields.rooms,
				userToRooms: tt.fields.userInRooms,
				roomToUsers: tt.fields.roomWithUsers,
			}
			if err := r.JoinRoom(tt.args.userName, tt.args.roomName); (err != nil) != tt.wantErr {
				t.Errorf("InMemoryRepository.Join() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
