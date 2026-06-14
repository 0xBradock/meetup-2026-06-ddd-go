package httpserver

import (
	"context"
	"errors"
	"go-svr/health"
	"log/slog"
	"net/http"
)

type healthHandler struct {
	s health.Service
}

func (h *healthHandler) getHealth(ctx context.Context, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		out, err := h.s.GetHealth(ctx)
		if err != nil {
			logger.Error("health check failed", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := httpResponse(w, http.StatusOK, out); err != nil {
			logger.Error("failed to encode health response", "error", err)
		}
	}
}

func (h *healthHandler) getPing(ctx context.Context, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Message string `json:"message"`
		}

		req, err := decode[request](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		out, err := h.s.Ping(ctx, health.PingInput{Message: req.Message})
		if err != nil {
			if errors.Is(err, health.ErrEmptyPingMessage) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := httpResponse(w, http.StatusOK, out); err != nil {
			logger.Error("failed to encode ping response", "error", err)
		}
	}
}
