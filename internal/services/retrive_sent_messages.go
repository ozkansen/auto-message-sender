package services

import (
	"context"
	"fmt"

	"auto-message-sender/internal/models"
)

type getListCache interface {
	GetList(ctx context.Context) ([]models.MessageSenderResponse, error)
}

type RetrieveSentMessagesService struct {
	getListCache getListCache
}

func NewRetrieveSentMessagesService(getListCache getListCache) *RetrieveSentMessagesService {
	return &RetrieveSentMessagesService{
		getListCache: getListCache,
	}
}

func (s *RetrieveSentMessagesService) RetrieveSentMessages(ctx context.Context) ([]models.MessageSenderResponse, error) {
	messages, err := s.getListCache.GetList(ctx)
	if err != nil {
		return nil, fmt.Errorf("getListCache.GetList error: %w", err)
	}
	return messages, nil
}
