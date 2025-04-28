package models

import "time"

type Message struct {
	ID        string    `json:"id"`
	Sender    string    `json:"sender"`
	Content   string    `json:"content"`
	RoomId    string    `json:"roomId"`
	Timestamp time.Time `json:"timestamp"`
}
