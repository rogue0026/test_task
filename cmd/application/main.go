package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/rogue0026/test_/internal/storage/users/postgres"
	"github.com/rogue0026/test_/internal/transport/http/handlers"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	repo, err := postgres.New(context.Background(), "postgres://user:password@localhost:5432/test")
	if err != nil {
		logger.Error("error", err.Error())
		return
	}

	appRouter := chi.NewRouter()
	appRouter.Method(http.MethodPost, "/api/test-service/users", handlers.RegisterUser(logger, repo))
	appRouter.Method(http.MethodPost, "/api/test-service/users/verify", handlers.VerifyUser(logger, repo))

	s := http.Server{
		Addr:    ":8080",
		Handler: appRouter,
	}
	logger.Debug("starting listening at", slog.String("address", s.Addr))
	err = s.ListenAndServe()
	if err != nil {
		slog.Error("error", err.Error())
		return
	}
}
