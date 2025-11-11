package sender

import (
	"context"
	"log/slog"

	"auto-message-sender/internal/models"
)

var _ messageSender = (*WebhookMessageSenderWithLogger)(nil)

type WebhookMessageSenderWithLogger struct {
	logger      *slog.Logger
	baseService messageSender
}

func NewWebhookMessageSenderWithLogger(logger *slog.Logger, baseService messageSender) *WebhookMessageSenderWithLogger {
	return &WebhookMessageSenderWithLogger{
		logger:      logger,
		baseService: baseService,
	}
}

func (s *WebhookMessageSenderWithLogger) SendMessage(ctx context.Context, message models.Message) (models.MessageSenderResponse, error) {
	s.logger.Info("message sent", "message", message)
	return s.baseService.SendMessage(ctx, message)
}
