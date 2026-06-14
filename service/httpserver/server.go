package httpserver

import (
	"context"
	"go-svr/catalog"
	"go-svr/config"
	"go-svr/health"
	"go-svr/loan"
	"go-svr/membership"
	"go-svr/overdue"
	"go-svr/reservation"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

func NewHandler(
	ctx context.Context,
	logger *slog.Logger,
	healthService health.Service,
	catalogService *catalog.Service,
	loanService *loan.Service,
	overdueService *overdue.Service,
	reservationService *reservation.Service,
	membershipService *membership.Service,
) http.Handler {
	mux := http.NewServeMux()

	addRoutes(
		ctx,
		mux,
		logger,
		healthHandler{s: healthService},
		catalogHandler{s: catalogService},
		loanHandler{s: loanService},
		overdueHandler{s: overdueService},
		reservationHandler{s: reservationService},
		membershipHandler{s: membershipService},
	)

	return mux
}

// RunHTTPServer is responsible for all top-level functionality
// common to all endpoints: routes, middlewares, CORS, auth, logging, etc.
func RunHTTPServer(
	ctx context.Context,
	wg *sync.WaitGroup,
	logger *slog.Logger,
	cfg *config.ServerConfig,
	healthService health.Service,
	catalogService *catalog.Service,
	loanService *loan.Service,
	overdueService *overdue.Service,
	reservationService *reservation.Service,
	membershipService *membership.Service,
) {
	wg.Add(1)

	handler := NewHandler(
		ctx,
		logger,
		healthService,
		catalogService,
		loanService,
		overdueService,
		reservationService,
		membershipService,
	)

	// handler = Recover(handler)
	// handler = loggerMW(handler)
	// handler = authMW(handler)
	// handler = corsMW(handler)

	httpServer := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: handler,
	}

	ready := make(chan bool)

	go func() {
		logger.Info("starting server", "addr", cfg.HTTPAddr)

		ready <- true
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("failed to start server", "err", err)
		}
	}()

	go func() {
		<-ready

		defer wg.Done()
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			time.Duration(cfg.ServerShutdownTime)*time.Second,
		)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error("failed to shutdown server", "err", err)
		}
	}()
}
