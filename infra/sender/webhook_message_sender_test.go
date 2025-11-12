package sender

import (
	"context"
	"testing"
	"time"

	"auto-message-sender/internal/models"
)

func TestWebhookMessageSender(t *testing.T) {
	// webhook.site test
	webhookMessageSender := NewWebhookMessageSender("https://webhook.site/264d7ada-f7a7-40e9-8f30-eb0bde016436")
	_, err := webhookMessageSender.SendMessage(context.Background(), models.Message{
		MessageID:      "123",
		PhoneNumber:    "+901231231",
		MessageContent: "asdadasdsa",
		SendingStatus:  "pending",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		t.Fatal(err)
	}
}
