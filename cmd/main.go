package main

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"auto-message-sender/internal/handlers"
	"auto-message-sender/internal/services"
)

func main() {
	autoMessageSenderServices := services.NewAutoMessageSender(nil, nil, nil, 2)
	messagesService := services.NewRetrieveSentMessagesService(nil)

	messagesHandler := handlers.NewMessagesHandler(messagesService)
	autoSenderStartStopHandler := handlers.NewAutoSenderStartStopHandler(autoMessageSenderServices)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /messages", messagesHandler.RetrieveSentMessagesHandler)
	mux.HandleFunc("POST /start", autoSenderStartStopHandler.Start)
	mux.HandleFunc("POST /stop", autoSenderStartStopHandler.Stop)

	wg := sync.WaitGroup{}
	server := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	wg.Go(func() {
		err := server.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			panic(err)
		}
	})
	wg.Go(func() {
		err := server.Shutdown(context.TODO())
		if err != nil {
			panic(err)
		}
	})
	wg.Go(func() {
		err := autoMessageSenderServices.Run(context.TODO())
		if err != nil {
			panic(err)
		}
	})
	wg.Wait()
}
