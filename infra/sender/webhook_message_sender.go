package sender

import (
	"context"
	"time"

	"auto-message-sender/internal/models"
)

type messageSender interface {
	SendMessage(ctx context.Context, message models.Message) (models.MessageSenderResponse, error)
}

var _ messageSender = (*WebhookMessageSender)(nil)

type WebhookMessageSender struct{}

func NewWebhookMessageSender() *WebhookMessageSender {
	return &WebhookMessageSender{}
}

func (s *WebhookMessageSender) SendMessage(ctx context.Context, message models.Message) (models.MessageSenderResponse, error) {
	return models.MessageSenderResponse{
		Message:   "Accepted",
		MessageID: message.MessageID,
		SentAt:    time.Now().UTC(),
	}, nil
}
