package cache

import (
	"context"
	"log/slog"

	"auto-message-sender/internal/models"
)

var _ getListCache = (*GetListCacheWithLogger)(nil)

type GetListCacheWithLogger struct {
	logger      *slog.Logger
	baseService getListCache
}

func NewGetListCacheWithLogger(logger *slog.Logger, baseService getListCache) *GetListCacheWithLogger {
	return &GetListCacheWithLogger{
		logger:      logger,
		baseService: baseService,
	}
}

func (g *GetListCacheWithLogger) GetList(ctx context.Context) ([]models.MessageSenderResponse, error) {
	messages, err := g.baseService.GetList(ctx)
	if err != nil {
		g.logger.Error("GetListCacheWithLogger.GetList error:", "error", err)
		return messages, err
	}
	if len(messages) == 0 {
		g.logger.Debug("GetListCacheWithLogger.GetList success but message not found")
		return messages, nil
	}
	g.logger.Debug("GetListCacheWithLogger.GetList success:", "count", len(messages))
	return messages, nil
}
