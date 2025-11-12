package sender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"auto-message-sender/internal/models"
)

type messageSender interface {
	SendMessage(ctx context.Context, message models.Message) (models.MessageSenderResponse, error)
}

var _ messageSender = (*WebhookMessageSender)(nil)

type WebhookMessageSender struct {
	client         *http.Client
	webhookSiteURL string
}

func NewWebhookMessageSender(webhookSiteURL string) *WebhookMessageSender {
	return &WebhookMessageSender{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		webhookSiteURL: webhookSiteURL,
	}
}

type webhookMessage struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

type webhookMessageResponse struct {
	Message   string `json:"message"`
	MessageID string `json:"messageId"`
}

func (s *WebhookMessageSender) SendMessage(ctx context.Context, message models.Message) (models.MessageSenderResponse, error) {
	newWebhookMessage := webhookMessage{
		To:      message.PhoneNumber,
		Content: message.MessageContent,
	}
	body := bytes.Buffer{}
	err := json.NewEncoder(&body).Encode(newWebhookMessage)
	if err != nil {
		return models.MessageSenderResponse{}, fmt.Errorf("json.NewEncoder error: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.webhookSiteURL, &body)
	if err != nil {
		return models.MessageSenderResponse{}, fmt.Errorf("http.NewRequestWithContext error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return models.MessageSenderResponse{}, fmt.Errorf("s.client.Do error: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		return models.MessageSenderResponse{}, fmt.Errorf("webhook message sender unexpected response code error: %d", resp.StatusCode)
	}
	var webhookMessageResponseData webhookMessageResponse
	err = json.NewDecoder(resp.Body).Decode(&webhookMessageResponseData)
	if err != nil {
		return models.MessageSenderResponse{}, fmt.Errorf("json.NewDecoder error: %w", err)
	}
	return models.MessageSenderResponse{
		Message:   webhookMessageResponseData.Message,
		MessageID: webhookMessageResponseData.MessageID,
		SentAt:    time.Now().UTC(),
	}, nil
}
