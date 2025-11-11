package models

import (
	"time"
)

type Message struct {
	MessageID      string
	PhoneNumber    string
	MessageContent string
	SendingStatus  string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
