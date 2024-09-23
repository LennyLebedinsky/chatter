package message

import (
	"context"
	"reflect"
	"testing"
	"time"
)

const testRoomName = "tower"

var testMessage1 = &Message{
	User:           "ultron",
	Room:           testRoomName,
	IsNotification: false,
	Value:          "Bow to me, minion!",
	ServerTime:     time.Time{},
}

var testMessage2 = &Message{
	User:           "jarvis",
	Room:           testRoomName,
	IsNotification: false,
	Value:          "Never!!!",
	ServerTime:     time.Time{},
}

func TestInMemoryStore_GetMessages(t *testing.T) {
	type fields struct {
		history map[string][]*Message
	}
	type args struct {
		roomName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*Message
		wantErr bool
	}{
		{
			name: "Getting all messages by the room name should succeed",
			fields: fields{
				history: map[string][]*Message{
					testRoomName: {testMessage1, testMessage2},
				},
			},
			args: args{
				roomName: testRoomName,
			},
			want: []*Message{
				testMessage1,
				testMessage2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &InMemoryStore{
				history: tt.fields.history,
			}
			got, err := s.GetMessages(context.Background(), tt.args.roomName)
			if (err != nil) != tt.wantErr {
				t.Errorf("InMemoryStore.GetMessages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InMemoryStore.GetMessages() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemoryStore_SaveMessage(t *testing.T) {
	type fields struct {
		history map[string][]*Message
	}
	type args struct {
		roomName string
		msg      *Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Saving message to the room should succeed",
			fields: fields{
				history: map[string][]*Message{},
			},
			args: args{
				roomName: testRoomName,
				msg:      testMessage1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &InMemoryStore{
				history: tt.fields.history,
			}
			if err := s.SaveMessage(context.Background(), tt.args.roomName, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("InMemoryStore.SaveMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
