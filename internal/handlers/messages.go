package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"auto-message-sender/internal/models"
)

type retrieveSentMessagesService interface {
	RetrieveSentMessages(ctx context.Context) ([]models.MessageSenderResponse, error)
}
type MessagesHandler struct {
	retrieveSentMessagesService retrieveSentMessagesService
}

func NewMessagesHandler(retrieveSentMessagesService retrieveSentMessagesService) *MessagesHandler {
	return &MessagesHandler{
		retrieveSentMessagesService: retrieveSentMessagesService,
	}
}

func (h *MessagesHandler) RetrieveSentMessagesHandler(w http.ResponseWriter, r *http.Request) {
	messages, err := h.retrieveSentMessagesService.RetrieveSentMessages(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(messages)
}
