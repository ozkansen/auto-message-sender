package services

import (
	"context"
	"fmt"
	"time"

	"auto-message-sender/internal/models"
)

type messageRepository interface {
	GetUnsentMessages(ctx context.Context, limit int) ([]models.Message, error)
	UpdateMessageStatus(ctx context.Context, messageID, sendingStatus string) error
}

type messageSender interface {
	SendMessage(ctx context.Context, message models.Message) (models.MessageSenderResponse, error)
}

type setCache interface {
	Set(ctx context.Context, message models.MessageSenderResponse) error
}

type AutoMessageSender struct {
	messageRepository        messageRepository
	messageSender            messageSender
	cache                    setCache
	maxGetUnsentMessageLimit int
	stopSignal               chan struct{}
	startSignal              chan struct{}
}

func NewAutoMessageSender(
	messageRepository messageRepository,
	messageSender messageSender,
	cache setCache,
	maxGetUnsentMessageLimit int,
) *AutoMessageSender {
	return &AutoMessageSender{
		messageRepository:        messageRepository,
		messageSender:            messageSender,
		cache:                    cache,
		maxGetUnsentMessageLimit: maxGetUnsentMessageLimit,
		stopSignal:               make(chan struct{}),
		startSignal:              make(chan struct{}),
	}
}

func (s *AutoMessageSender) Run(ctx context.Context) error {
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return context.Cause(ctx)
		case <-ticker.C:
			err := s.sendMessages(ctx)
			if err != nil {
				return fmt.Errorf("sendMessages error: %w", err)
			}
		case <-s.stopSignal:
			ticker.Stop()
		case <-s.startSignal:
			ticker.Reset(2 * time.Minute)
		}
	}
}

func (s *AutoMessageSender) Start() {
	s.startSignal <- struct{}{}
}

func (s *AutoMessageSender) Stop() {
	s.stopSignal <- struct{}{}
}

func (s *AutoMessageSender) sendMessages(ctx context.Context) error {
	messages, err := s.messageRepository.GetUnsentMessages(ctx, s.maxGetUnsentMessageLimit)
	if err != nil {
		return fmt.Errorf("messageRepository.GetUnsentMessages error: %w", err)
	}
	for _, message := range messages {
		sendMessageResponse, err2 := s.messageSender.SendMessage(ctx, message)
		if err2 != nil {
			return fmt.Errorf("messageSender.SendMessage error: %w", err2)
		}
		err2 = s.cache.Set(ctx, sendMessageResponse)
		if err2 != nil {
			return fmt.Errorf("cache.Set error: %w", err2)
		}
		err2 = s.messageRepository.UpdateMessageStatus(ctx, message.MessageID, "sent")
		if err2 != nil {
			return fmt.Errorf("messageRepository.UpdateMessageStatus error: %w", err2)
		}
	}
	return nil
}
