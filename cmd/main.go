package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"auto-message-sender/infra/repository"
	"auto-message-sender/internal/handlers"
	"auto-message-sender/internal/services"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()
	connectionString := getPostgresqlDSNFromEnv()
	conn, err := newPostgresqlDBConn(ctx, connectionString)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	messageRepository := repository.NewMessagePostgresqlRepository(conn)
	autoMessageSenderServices := services.NewAutoMessageSender(messageRepository, nil, nil, 2)
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

func newPostgresqlDBConn(ctx context.Context, dbConnectionDsn string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, dbConnectionDsn)
	if err != nil {
		return nil, fmt.Errorf("pgx.Connect error: %w", err)
	}
	err = conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("pgx.Ping error: %w", err)
	}
	return conn, nil
}

// getPostgresqlDSNFromEnv return back example connection string like
// "postgres://dbuser:dbpassword@postgresdb:5432/automessagesenderdb?sslmode=disable"
func getPostgresqlDSNFromEnv() string {
	return os.Getenv("POSTGRESQL_DSN")
}
