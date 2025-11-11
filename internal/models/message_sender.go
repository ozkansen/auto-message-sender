package models

import (
	"time"
)

type MessageSenderResponse struct {
	Message   string    `json:"message"`
	MessageID string    `json:"message_id"`
	SentAt    time.Time `json:"sent_at"`
}
