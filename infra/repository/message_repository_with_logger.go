package repository

import (
	"context"
	"log/slog"

	"auto-message-sender/internal/models"
)

var _ messageRepository = (*MessageRepositoryWithLogger)(nil)

type MessageRepositoryWithLogger struct {
	logger      *slog.Logger
	baseService messageRepository
}

func NewMessageRepositoryWithLogger(logger *slog.Logger, baseService messageRepository) *MessageRepositoryWithLogger {
	return &MessageRepositoryWithLogger{
		logger:      logger,
		baseService: baseService,
	}
}

func (m *MessageRepositoryWithLogger) GetUnsentMessages(ctx context.Context, limit int) ([]models.Message, error) {
	messages, err := m.baseService.GetUnsentMessages(ctx, limit)
	if err != nil {
		m.logger.Error("MessageRepositoryWithLogger.GetUnsentMessages error:", "error", err)
		return messages, err
	}
	if len(messages) == 0 {
		m.logger.Debug("MessageRepositoryWithLogger.GetUnsentMessages success but new message not found")
		return messages, nil
	}
	m.logger.Debug("MessageRepositoryWithLogger.GetUnsentMessages success:", "count", len(messages))
	return messages, nil
}

func (m *MessageRepositoryWithLogger) UpdateMessageStatus(ctx context.Context, messageID, sendingStatus string) error {
	err := m.baseService.UpdateMessageStatus(ctx, messageID, sendingStatus)
	if err != nil {
		m.logger.Error("UpdateMessageStatus error:", "error", err)
		return err
	}
	m.logger.Debug("UpdateMessageStatus success:", "messageID", messageID, "sendingStatus", sendingStatus)
	return nil
}
