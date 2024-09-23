package message

import "time"

// Message represents main object of exchange between users which is published in the rooms.
// Some of messages can be marked as notifications for housekeeping and letting users know
// on what's going on with other users or the system.
// Notifications can be treated separately from messages, e.g. not published in the rooms.
type Message struct {
	User           string    `json:"user"`
	Room           string    `json:"room"`
	IsNotification bool      `json:"isNotification"`
	Value          string    `json:"value"`
	ServerTime     time.Time `json:"serverTime"`
}
