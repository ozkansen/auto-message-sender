package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"auto-message-sender/infra/cache"
	"auto-message-sender/infra/repository"
	"auto-message-sender/infra/sender"
	"auto-message-sender/internal/handlers"
	"auto-message-sender/internal/services"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func main() {
	logger := newSlogLogger()
	logger.Info("starting application")
	defer logger.Info("application stopped")

	ctx := context.Background()

	connectionString := getPostgresqlDSNFromEnv()
	conn, err := newPostgresqlDBConn(ctx, connectionString)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	redisAddr := getRedisAddrFromEnv()
	client, err := newRedisClient(ctx, redisAddr)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	webhookMessageSender := sender.NewWebhookMessageSender()
	webhookMessageSenderWithLogger := sender.NewWebhookMessageSenderWithLogger(logger, webhookMessageSender)
	messageRepository := repository.NewMessagePostgresqlRepository(conn)
	messageRepositoryWithLogger := repository.NewMessageRepositoryWithLogger(logger, messageRepository)
	setCache := cache.NewSetCache(client)
	setCacheWithLogger := cache.NewSetCacheWithLogger(logger, setCache)
	autoMessageSenderServices := services.NewAutoMessageSender(messageRepositoryWithLogger, webhookMessageSenderWithLogger, setCacheWithLogger, 2)

	getListCache := cache.NewGetListCache(client)
	getListCacheWithLogger := cache.NewGetListCacheWithLogger(logger, getListCache)
	messagesService := services.NewRetrieveSentMessagesService(getListCacheWithLogger)

	messagesHandler := handlers.NewMessagesHandler(messagesService)
	autoSenderStartStopHandler := handlers.NewAutoSenderStartStopHandler(autoMessageSenderServices)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /messages", messagesHandler.RetrieveSentMessagesHandler)
	mux.HandleFunc("POST /start", autoSenderStartStopHandler.Start)
	mux.HandleFunc("POST /stop", autoSenderStartStopHandler.Stop)

	server := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	wg := sync.WaitGroup{}
	wg.Go(func() {
		logger.Info("starting http server")
		err2 := server.ListenAndServe()
		if err2 != nil {
			if errors.Is(err2, http.ErrServerClosed) {
				logger.Warn("http server stopped")
				return
			}
			logger.Error("http server error", "error", err2)
			panic(err2)
		}
	})
	wg.Go(func() {
		<-ctx.Done()
		logger.Info("shutting down http server")
		err2 := server.Shutdown(ctx)
		if err2 != nil {
			logger.Error("http server shutdown error", "error", err2)
			panic(err2)
		}
	})
	wg.Go(func() {
		logger.Info("starting auto message sender")
		err2 := autoMessageSenderServices.Run(ctx)
		if err2 != nil {
			logger.Error("auto message sender error", "error", err2)
			panic(err2)
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

func newRedisClient(ctx context.Context, redisAddr string) (*redis.Client, error) {
	opt, err := redis.ParseURL(redisAddr)
	if err != nil {
		return nil, fmt.Errorf("redis.ParseURL error: %w", err)
	}
	client := redis.NewClient(opt)
	status := client.Ping(ctx)
	if status.Err() != nil {
		return nil, fmt.Errorf("redis.Ping error: %w", err)
	}
	return client, nil
}

func getRedisAddrFromEnv() string {
	return os.Getenv("REDIS_ADDR")
}

func newSlogLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}
