package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vokinneberg/ya-praktukum-go-testing-workshop/config"
	"github.com/vokinneberg/ya-praktukum-go-testing-workshop/internal/http/handler"
	"github.com/vokinneberg/ya-praktukum-go-testing-workshop/internal/http/middleware"
	"github.com/vokinneberg/ya-praktukum-go-testing-workshop/internal/url"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	// Initialize the logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})).With("app", "shortener")

	// Initialize the context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize the database
	dbpool, err := pgxpool.New(ctx, cfg.DatabaseDSN)
	if err != nil {
		logger.Error("Unable to create connection pool", "error", err)
		panic(err)
	}
	defer dbpool.Close()

	// Check if the database is accessible
	if err := dbpool.Ping(ctx); err != nil {
		logger.Error("Unable to access database", "error", err)
		panic(err)
	}

	// Initialize the repository
	repository := url.NewRepository(dbpool, cfg.SyncWrite)
	if cfg.Restore {
		if err := repository.Restore(ctx); err != nil {
			logger.Error("Error restoring repository", "error", err)
			panic(err)
		}
	}
	defer func() {
		logger.Info("Flushing repository before shutdown")
		if err := repository.Flush(ctx); err != nil {
			logger.Error("Error flushing repository", "error", err)
		}
	}()

	// Initialize the service
	service := url.NewService(repository, cfg)

	//Initialize the handler
	handler := handler.New(service, logger)

	// Initialize the mux
	r := http.NewServeMux()
	r.HandleFunc("POST /api/shorten", handler.ShortURLHandler)
	r.HandleFunc("GET /{id}", handler.RedirectURLHandler)
	r.HandleFunc("GET /ping", handler.PingHandler)

	// Initialize the logger middleware
	loggingMiddleware := middleware.Logging(logger)

	go func() {
		if err := http.ListenAndServe(":8080", loggingMiddleware(r)); err != nil {
			logger.Error("Error starting server", "error", err)
			panic(err)
		}
	}()

	logger.Info("Server started", "port", 8080)

	// Implement the graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logger.Info("Shutting down server")
}
