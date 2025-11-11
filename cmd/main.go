package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"auto-message-sender/infra/cache"
	"auto-message-sender/infra/repository"
	"auto-message-sender/internal/handlers"
	"auto-message-sender/internal/services"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func main() {
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

	setCache := cache.NewSetCache(client)
	messageRepository := repository.NewMessagePostgresqlRepository(conn)
	autoMessageSenderServices := services.NewAutoMessageSender(messageRepository, nil, setCache, 2)

	getListCache := cache.NewGetListCache(client)
	messagesService := services.NewRetrieveSentMessagesService(getListCache)

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
		<-ctx.Done()
		err := server.Shutdown(ctx)
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
