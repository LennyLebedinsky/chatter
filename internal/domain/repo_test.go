package domain

import (
	"context"
	"reflect"
	"testing"
)

var testUserA = &User{
	Name: "jarvis",
}

var testUserB = &User{
	Name: "ultron",
}

var testRoomA = &Room{
	Name:    "tower",
	Creator: testUserA,
}

var testRoomB = &Room{
	Name: "sokovia",
}

func TestInMemoryRepository_CreateUser(t *testing.T) {
	type fields struct {
		users       map[string]*User
		rooms       map[string]*Room
		userToRooms map[*User][]*Room
		roomToUsers map[*Room][]*User
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
				users:       map[string]*User{},
				rooms:       map[string]*Room{},
				userToRooms: map[*User][]*Room{},
				roomToUsers: map[*Room][]*User{},
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
				rooms:       map[string]*Room{},
				userToRooms: map[*User][]*Room{},
				roomToUsers: map[*Room][]*User{},
			},
			args: args{
				userName: testUserA.Name,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InMemoryRepository{
				users:       tt.fields.users,
				rooms:       tt.fields.rooms,
				userToRooms: tt.fields.userToRooms,
				roomToUsers: tt.fields.roomToUsers,
			}
			got, err := r.CreateUser(context.Background(), tt.args.userName)
			if (err != nil) != tt.wantErr {
				t.Errorf("InMemoryRepository.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InMemoryRepository.CreateUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemoryRepository_FindUser(t *testing.T) {
	type fields struct {
		users       map[string]*User
		rooms       map[string]*Room
		userToRooms map[*User][]*Room
		roomToUsers map[*Room][]*User
	}
	type args struct {
		userName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *User
	}{
		{
			name: "Searching for existing user should succeed",
			fields: fields{
				users: map[string]*User{
					testUserA.Name: testUserA,
				},
				rooms:       map[string]*Room{},
				userToRooms: map[*User][]*Room{},
				roomToUsers: map[*Room][]*User{},
			},
			args: args{
				userName: testUserA.Name,
			},
			want: testUserA,
		},
		{
			name: "Searching for non-existing user should return nil",
			fields: fields{
				users:       map[string]*User{},
				rooms:       map[string]*Room{},
				userToRooms: map[*User][]*Room{},
				roomToUsers: map[*Room][]*User{},
			},
			args: args{
				userName: testUserA.Name,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InMemoryRepository{
				users:       tt.fields.users,
				rooms:       tt.fields.rooms,
				userToRooms: tt.fields.userToRooms,
				roomToUsers: tt.fields.roomToUsers,
			}
			if got := r.FindUser(context.Background(), tt.args.userName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InMemoryRepository.FindUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemoryRepository_FindRoom(t *testing.T) {
	type fields struct {
		users       map[string]*User
		rooms       map[string]*Room
		userToRooms map[*User][]*Room
		roomToUsers map[*Room][]*User
	}
	type args struct {
		roomName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Room
	}{
		{
			name: "Searching for existing room should succeed",
			fields: fields{
				users: map[string]*User{},
				rooms: map[string]*Room{
					testRoomB.Name: testRoomB,
				},
				userToRooms: map[*User][]*Room{},
				roomToUsers: map[*Room][]*User{},
			},
			args: args{
				roomName: testRoomB.Name,
			},
			want: testRoomB,
		},
		{
			name: "Searching for non-existing room should return nil",
			fields: fields{
				users:       map[string]*User{},
				rooms:       map[string]*Room{},
				userToRooms: map[*User][]*Room{},
				roomToUsers: map[*Room][]*User{},
			},
			args: args{
				roomName: testRoomB.Name,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InMemoryRepository{
				users:       tt.fields.users,
				rooms:       tt.fields.rooms,
				userToRooms: tt.fields.userToRooms,
				roomToUsers: tt.fields.roomToUsers,
			}
			if got := r.FindRoom(context.Background(), tt.args.roomName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InMemoryRepository.FindRoom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemoryRepository_JoinRoom(t *testing.T) {
	type fields struct {
		users       map[string]*User
		rooms       map[string]*Room
		userToRooms map[*User][]*Room
		roomToUsers map[*Room][]*User
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
		{
			name: "Existing user joining existing room should succeed",
			fields: fields{
				users: map[string]*User{
					testUserA.Name: testUserA,
				},
				rooms: map[string]*Room{
					testRoomB.Name: testRoomB,
				},
				userToRooms: map[*User][]*Room{},
				roomToUsers: map[*Room][]*User{},
			},
			args: args{
				userName: testUserA.Name,
				roomName: testRoomB.Name,
			},
			wantErr: false,
		},
		{
			name: "Existing user joining existing room which they previously joined should succeed (no-op)",
			fields: fields{
				users: map[string]*User{
					testUserA.Name: testUserA,
				},
				rooms: map[string]*Room{
					testRoomB.Name: testRoomB,
				},
				userToRooms: map[*User][]*Room{
					testUserA: {testRoomB},
				},
				roomToUsers: map[*Room][]*User{
					testRoomB: {testUserA},
				},
			},
			args: args{
				userName: testUserA.Name,
				roomName: testRoomB.Name,
			},
			wantErr: false,
		},
		{
			name: "Non-existing user joining existing room should fail",
			fields: fields{
				users: map[string]*User{},
				rooms: map[string]*Room{
					testRoomB.Name: testRoomB,
				},
				userToRooms: map[*User][]*Room{},
				roomToUsers: map[*Room][]*User{},
			},
			args: args{
				userName: testUserA.Name,
				roomName: testRoomB.Name,
			},
			wantErr: true,
		},
		{
			name: "Existing user joining non-existing room should fail",
			fields: fields{
				users: map[string]*User{
					testUserA.Name: testUserA,
				},
				rooms:       map[string]*Room{},
				userToRooms: map[*User][]*Room{},
				roomToUsers: map[*Room][]*User{},
			},
			args: args{
				userName: testUserA.Name,
				roomName: testRoomB.Name,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InMemoryRepository{
				users:       tt.fields.users,
				rooms:       tt.fields.rooms,
				userToRooms: tt.fields.userToRooms,
				roomToUsers: tt.fields.roomToUsers,
			}
			if err := r.JoinRoom(context.Background(), tt.args.userName, tt.args.roomName); (err != nil) != tt.wantErr {
				t.Errorf("InMemoryRepository.JoinRoom() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInMemoryRepository_LeaveRoom(t *testing.T) {
	type fields struct {
		users       map[string]*User
		rooms       map[string]*Room
		userToRooms map[*User][]*Room
		roomToUsers map[*Room][]*User
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
		{
			name: "Existing user leaving existing room they previously joined should succeed",
			fields: fields{
				users: map[string]*User{
					testUserA.Name: testUserA,
				},
				rooms: map[string]*Room{
					testRoomB.Name: testRoomB,
				},
				userToRooms: map[*User][]*Room{
					testUserA: {testRoomB},
				},
				roomToUsers: map[*Room][]*User{
					testRoomB: {testUserA},
				},
			},
			args: args{
				userName: testUserA.Name,
				roomName: testRoomB.Name,
			},
			wantErr: false,
		},

		{
			name: "Existing user 'leaving' existing room which they never joined should succeed (no-op)",
			fields: fields{
				users: map[string]*User{
					testUserA.Name: testUserA,
				},
				rooms: map[string]*Room{
					testRoomB.Name: testRoomB,
				},
				userToRooms: map[*User][]*Room{},
				roomToUsers: map[*Room][]*User{},
			},
			args: args{
				userName: testUserA.Name,
				roomName: testRoomB.Name,
			},
			wantErr: false,
		},
		{
			name: "Non-existing user leaving existing room should fail",
			fields: fields{
				users: map[string]*User{},
				rooms: map[string]*Room{
					testRoomB.Name: testRoomB,
				},
				userToRooms: map[*User][]*Room{},
				roomToUsers: map[*Room][]*User{},
			},
			args: args{
				userName: testUserA.Name,
				roomName: testRoomB.Name,
			},
			wantErr: true,
		},
		{
			name: "Existing user leaving non-existing room should fail",
			fields: fields{
				users: map[string]*User{
					testUserA.Name: testUserA,
				},
				rooms:       map[string]*Room{},
				userToRooms: map[*User][]*Room{},
				roomToUsers: map[*Room][]*User{},
			},
			args: args{
				userName: testUserA.Name,
				roomName: testRoomB.Name,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InMemoryRepository{
				users:       tt.fields.users,
				rooms:       tt.fields.rooms,
				userToRooms: tt.fields.userToRooms,
				roomToUsers: tt.fields.roomToUsers,
			}
			if err := r.LeaveRoom(context.Background(), tt.args.userName, tt.args.roomName); (err != nil) != tt.wantErr {
				t.Errorf("InMemoryRepository.LeaveRoom() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInMemoryRepository_CreateRoom(t *testing.T) {
	type fields struct {
		users       map[string]*User
		rooms       map[string]*Room
		userToRooms map[*User][]*Room
		roomToUsers map[*Room][]*User
	}
	type args struct {
		roomName        string
		creatorUserName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Room
		wantErr bool
	}{
		{
			name: "Existing user creating a new room should succeed",
			fields: fields{
				users: map[string]*User{
					testUserA.Name: testUserA,
				},
				rooms:       map[string]*Room{},
				userToRooms: map[*User][]*Room{},
				roomToUsers: map[*Room][]*User{},
			},
			args: args{
				roomName:        testRoomA.Name,
				creatorUserName: testUserA.Name,
			},
			want:    testRoomA,
			wantErr: false,
		},
		{
			name: "Existing user trying to create a room with the same name should fail",
			fields: fields{
				users: map[string]*User{
					testUserA.Name: testUserA,
				},
				rooms: map[string]*Room{
					testRoomA.Name: testRoomA,
				},
				userToRooms: map[*User][]*Room{
					testUserA: {testRoomA},
				},
				roomToUsers: map[*Room][]*User{
					testRoomA: {testUserA},
				},
			},
			args: args{
				roomName:        testRoomA.Name,
				creatorUserName: testUserA.Name,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Non-existing user creating room should fail",
			fields: fields{
				users:       map[string]*User{},
				rooms:       map[string]*Room{},
				userToRooms: map[*User][]*Room{},
				roomToUsers: map[*Room][]*User{},
			},
			args: args{
				roomName:        testRoomA.Name,
				creatorUserName: testUserA.Name,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InMemoryRepository{
				users:       tt.fields.users,
				rooms:       tt.fields.rooms,
				userToRooms: tt.fields.userToRooms,
				roomToUsers: tt.fields.roomToUsers,
			}
			got, err := r.CreateRoom(context.Background(), tt.args.roomName, tt.args.creatorUserName)
			if (err != nil) != tt.wantErr {
				t.Errorf("InMemoryRepository.CreateRoom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InMemoryRepository.CreateRoom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemoryRepository_ListRooms(t *testing.T) {
	type fields struct {
		users       map[string]*User
		rooms       map[string]*Room
		userToRooms map[*User][]*Room
		roomToUsers map[*Room][]*User
	}
	type args struct{}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Room
		wantErr bool
	}{
		{
			name: "Successful call should return list of all rooms",
			fields: fields{
				users: map[string]*User{},
				rooms: map[string]*Room{
					testRoomA.Name: testRoomA,
					testRoomB.Name: testRoomB,
				},
				userToRooms: map[*User][]*Room{},
				roomToUsers: map[*Room][]*User{},
			},
			args: args{},
			want: []*Room{
				testRoomA,
				testRoomB,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InMemoryRepository{
				users:       tt.fields.users,
				rooms:       tt.fields.rooms,
				userToRooms: tt.fields.userToRooms,
				roomToUsers: tt.fields.roomToUsers,
			}
			got, err := r.ListRooms(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("InMemoryRepository.ListRooms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InMemoryRepository.ListRooms() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemoryRepository_ListParticipants(t *testing.T) {
	type fields struct {
		users       map[string]*User
		rooms       map[string]*Room
		userToRooms map[*User][]*Room
		roomToUsers map[*Room][]*User
	}
	type args struct {
		roomName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*User
		wantErr bool
	}{
		{
			name: "Successful call should return list of users who have joined the room",
			fields: fields{
				users: map[string]*User{
					testUserA.Name: testUserA,
					testUserB.Name: testUserB,
				},
				rooms: map[string]*Room{
					testRoomA.Name: testRoomA,
				},
				userToRooms: map[*User][]*Room{
					testUserA: {testRoomA},
					testUserB: {testRoomA},
				},
				roomToUsers: map[*Room][]*User{
					testRoomA: {testUserA, testUserB},
				},
			},
			args: args{
				roomName: testRoomA.Name,
			},
			want: []*User{
				testUserA,
				testUserB,
			},
			wantErr: false,
		},
		{
			name: "Call for non-existent room should fail",
			fields: fields{
				users: map[string]*User{
					testUserA.Name: testUserA,
					testUserB.Name: testUserB,
				},
				rooms:       map[string]*Room{},
				userToRooms: map[*User][]*Room{},
				roomToUsers: map[*Room][]*User{},
			},
			args: args{
				roomName: testRoomA.Name,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InMemoryRepository{
				users:       tt.fields.users,
				rooms:       tt.fields.rooms,
				userToRooms: tt.fields.userToRooms,
				roomToUsers: tt.fields.roomToUsers,
			}
			got, err := r.ListParticipants(context.Background(), tt.args.roomName)
			if (err != nil) != tt.wantErr {
				t.Errorf("InMemoryRepository.ListParticipants() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InMemoryRepository.ListParticipants() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemoryRepository_ListParticipantsForAllRooms(t *testing.T) {
	type fields struct {
		users       map[string]*User
		rooms       map[string]*Room
		userToRooms map[*User][]*Room
		roomToUsers map[*Room][]*User
	}
	type args struct{}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*RoomParticipation
		wantErr bool
	}{
		{
			name: "Successful call should return list of all rooms and all users joinied each room",
			fields: fields{
				users: map[string]*User{
					testUserA.Name: testUserA,
					testUserB.Name: testUserB,
				},
				rooms: map[string]*Room{
					testRoomA.Name: testRoomA,
					testRoomB.Name: testRoomB,
				},
				userToRooms: map[*User][]*Room{
					testUserA: {testRoomA},
					testUserB: {testRoomA, testRoomB},
				},
				roomToUsers: map[*Room][]*User{
					testRoomA: {testUserA, testUserB},
					testRoomB: {testUserB},
				},
			},
			args: args{},
			want: []*RoomParticipation{
				{
					Room:         testRoomA,
					Participants: []*User{testUserA, testUserB},
				},
				{
					Room:         testRoomB,
					Participants: []*User{testUserB},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &InMemoryRepository{
				users:       tt.fields.users,
				rooms:       tt.fields.rooms,
				userToRooms: tt.fields.userToRooms,
				roomToUsers: tt.fields.roomToUsers,
			}
			got, err := r.ListParticipantsForAllRooms(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("InMemoryRepository.ListParticipantsForAllRooms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InMemoryRepository.ListParticipantsForAllRooms() = %v, want %v", got, tt.want)
			}
		})
	}
}
