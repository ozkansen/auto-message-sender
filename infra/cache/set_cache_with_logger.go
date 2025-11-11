package cache

import (
	"context"
	"log/slog"

	"auto-message-sender/internal/models"
)

var _ setCache = (*SetCacheWithLogger)(nil)

type SetCacheWithLogger struct {
	logger      *slog.Logger
	baseService setCache
}

func NewSetCacheWithLogger(logger *slog.Logger, baseService setCache) *SetCacheWithLogger {
	return &SetCacheWithLogger{
		logger:      logger,
		baseService: baseService,
	}
}

func (s *SetCacheWithLogger) Set(ctx context.Context, message models.MessageSenderResponse) error {
	err := s.baseService.Set(ctx, message)
	if err != nil {
		s.logger.Error("SetCacheWithLogger.Set error:", "error", err)
		return err
	}
	s.logger.Debug("SetCacheWithLogger.Set success:", "messageID", message.MessageID)
	return nil
}
