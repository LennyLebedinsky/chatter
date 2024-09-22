package message

import "time"

type Message struct {
	User           string    `json:"user"`
	Room           string    `json:"room"`
	IsNotification bool      `json:"isNotification"`
	Value          string    `json:"value"`
	ServerTime     time.Time `json:"serverTime"`
}
