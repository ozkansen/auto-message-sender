package cache

import (
	"context"
	"fmt"

	"auto-message-sender/internal/models"

	"github.com/redis/go-redis/v9"
)

type setCache interface {
	Set(ctx context.Context, message models.MessageSenderResponse) error
}

var _ setCache = (*SetCache)(nil)

type SetCache struct {
	client *redis.Client
}

func NewSetCache(client *redis.Client) *SetCache {
	return &SetCache{
		client: client,
	}
}

func (c *SetCache) Set(ctx context.Context, message models.MessageSenderResponse) error {
	newData := responseData{
		Message:   message.Message,
		MessageID: message.MessageID,
		SentAt:    message.SentAt,
	}
	key := fmt.Sprintf("sent_message_%s", message.MessageID)
	result := c.client.HSet(ctx, key, newData)
	err := result.Err()
	if err != nil {
		return err
	}
	return nil
}
