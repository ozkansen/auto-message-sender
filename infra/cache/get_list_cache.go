package cache

import (
	"context"
	"time"

	"auto-message-sender/internal/models"

	"github.com/redis/go-redis/v9"
)

type getListCache interface {
	GetList(ctx context.Context) ([]models.MessageSenderResponse, error)
}

var _ getListCache = (*GetListCache)(nil)

type GetListCache struct {
	client *redis.Client
}

func NewGetListCache(client *redis.Client) *GetListCache {
	return &GetListCache{
		client: client,
	}
}

type responseData struct {
	Message   string    `redis:"message"`
	MessageID string    `redis:"message_id"`
	SentAt    time.Time `redis:"sent_at"`
}

func (c *GetListCache) GetList(ctx context.Context) ([]models.MessageSenderResponse, error) {
	scanKeys := c.client.Keys(ctx, "sent_message*")
	if scanKeys.Err() != nil {
		return nil, scanKeys.Err()
	}
	resultKeys := scanKeys.Val()
	messages := make([]models.MessageSenderResponse, 0, 10)
	for _, key := range resultKeys {
		scan := c.client.HGetAll(ctx, key)
		if scan.Err() != nil {
			return nil, scan.Err()
		}
		var data responseData
		err2 := scan.Scan(&data)
		if err2 != nil {
			return nil, err2
		}
		messages = append(messages, models.MessageSenderResponse{
			Message:   data.Message,
			MessageID: data.MessageID,
			SentAt:    data.SentAt.UTC(),
		})
	}
	return messages, nil
}
