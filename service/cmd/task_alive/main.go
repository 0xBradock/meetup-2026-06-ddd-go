package main

import (
	"context"
	"fmt"
	"go-svr/config"
	"go-svr/db"
	"io"
	"log/slog"
	"os"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Args, os.Getenv, os.Stdout, os.Stderr); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "task error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, _ []string, getenv func(string) string, stdout io.Writer, _ io.Writer) error {
	taskID := getenv("TASK_ID")
	if taskID == "" {
		return fmt.Errorf("missing required env var: TASK_ID")
	}

	cfg, err := config.LoadServerConfig(getenv)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(stdout, &slog.HandlerOptions{Level: slog.LevelInfo})).With(
		"service", cfg.ServiceName,
		"environment", cfg.Environment,
		"task", "alive",
		"task_id", taskID,
	)

	pgRepos, closePG, err := db.NewPG(ctx, cfg.PGDriver, getenv)
	if err != nil {
		return fmt.Errorf("init repository: %w", err)
	}
	defer closePG()

	logger.Info("Running task")

	ping, err := pgRepos.Health.Ping(ctx, taskID)
	if err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	logger.Info("Task completed",
		"ping_message", ping.Message,
		"ping_received_at", ping.ReceivedAtUnix,
	)

	return nil
}
