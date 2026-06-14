package main

import (
	"context"
	"fmt"
	"go-svr/catalog"
	"go-svr/config"
	"go-svr/db"
	"go-svr/health"
	"go-svr/httpserver"
	"go-svr/loan"
	"go-svr/membership"
	"go-svr/overdue"
	"go-svr/reservation"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()
	if err := Run(ctx, os.Args, os.Getenv, os.Stdout, os.Stderr); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}

var listenFn = net.Listen

func Run(ctx context.Context, _ []string, getenv func(string) string, stdout io.Writer, _ io.Writer) error {
	cfg, err := config.LoadServerConfig(getenv)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(stdout, &slog.HandlerOptions{Level: slog.LevelInfo})).With(
		"service", cfg.ServiceName,
		"environment", cfg.Environment,
	)

	logger.Info("Starting service")

	// Driven
	logger.Info("Initializing repository", "pg_driver", cfg.PGDriver, "mssql_driver", cfg.MSSQLDriver)

	pgRepos, closePG, err := db.NewPG(ctx, cfg.PGDriver, getenv)
	if err != nil {
		return fmt.Errorf("init pg repository: %w", err)
	}
	defer closePG()

	mssqlRepos, closeMSSQL, err := db.NewMSSQL(ctx, cfg.MSSQLDriver, getenv)
	if err != nil {
		return fmt.Errorf("init mssql repository: %w", err)
	}
	defer closeMSSQL()

	healthRepo := pgRepos.Health
	if cfg.MSSQLDriver == "mssql" {
		healthRepo = mssqlRepos.Health
	}

	// Sub-Domains
	logger.Info("Creating service handlers")
	hs := health.NewService(healthRepo, logger.With("domain", "health"))

	// Bounded contexts — repositories are nil in the demo; services contain pseudo-code.
	catalogSvc := catalog.NewService(nil, logger.With("domain", "catalog"))
	loanSvc := loan.NewService(nil, logger.With("domain", "loan"))
	overdueSvc := overdue.NewService(nil, logger.With("domain", "overdue"))
	reservationSvc := reservation.NewService(nil, logger.With("domain", "reservation"))
	membershipSvc := membership.NewService(nil, logger.With("domain", "membership"))

	// Drivers
	logger.Info("Starting HTTP listener", "addr", cfg.HTTPAddr)
	httpLis, err := listenFn("tcp", cfg.HTTPAddr)
	if err != nil {
		return fmt.Errorf("listen http: %w", err)
	}
	defer func() {
		err = httpLis.Close()
		if err != nil {
			logger.Error("Failed to close HTTP listener", "error", err)
		}
	}()

	httpHandler := httpserver.NewHandler(ctx, logger, *hs, catalogSvc, loanSvc, overdueSvc, reservationSvc, membershipSvc)

	svr := &http.Server{
		Handler:      httpHandler,
		ReadTimeout:  cfg.HTTPReadTimeout,
		WriteTimeout: cfg.HTTPWriteTimeout,
		IdleTimeout:  cfg.HTTPIdleTimeout,
	}

	httpErr := make(chan error, 1)
	go func() {
		logger.Info("HTTP server listening", "addr", cfg.HTTPAddr)
		httpErr <- svr.Serve(httpLis)
	}()

	logger.Info("Service started successfully")

	shutdownCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case <-shutdownCtx.Done():
		logger.Info("Shutdown signal received")
		_ = svr.Shutdown(context.Background())
		return nil
	case err := <-httpErr:
		logger.Error("HTTP server error", "error", err)
		return fmt.Errorf("serve http: %w", err)
	}
}
