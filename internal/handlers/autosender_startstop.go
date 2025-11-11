package handlers

import (
	"fmt"
	"net/http"
)

type autoSenderStartStopService interface {
	Start()
	Stop()
}

type AutoSenderStartStopHandler struct {
	autoMessageSender autoSenderStartStopService
}

func NewAutoSenderStartStopHandler(autoMessageSender autoSenderStartStopService) *AutoSenderStartStopHandler {
	return &AutoSenderStartStopHandler{
		autoMessageSender: autoMessageSender,
	}
}

func (h *AutoSenderStartStopHandler) Start(w http.ResponseWriter, r *http.Request) {
	h.autoMessageSender.Start()
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, "OK")
}

func (h *AutoSenderStartStopHandler) Stop(w http.ResponseWriter, r *http.Request) {
	h.autoMessageSender.Stop()
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, "OK")
}
