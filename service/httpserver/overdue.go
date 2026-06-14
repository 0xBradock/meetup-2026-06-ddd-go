package httpserver

import (
	"context"
	"go-svr/overdue"
	"log/slog"
	"net/http"
)

type overdueHandler struct {
	s *overdue.Service
}

func (h *overdueHandler) listOpenCases(ctx context.Context, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		out, err := h.s.ListOpenCases(ctx)
		if err != nil {
			logger.Error("list open cases failed", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := httpResponse(w, http.StatusOK, out); err != nil {
			logger.Error("failed to encode response", "error", err)
		}
	}
}

func (h *overdueHandler) applyFine(ctx context.Context, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			DaysOverdue int `json:"daysOverdue"`
		}

		caseID := overdue.OverdueCaseID(r.PathValue("id"))
		req, err := decode[request](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		out, err := h.s.ApplyFine(ctx, overdue.ApplyFineInput{
			CaseID:      caseID,
			DaysOverdue: req.DaysOverdue,
		})
		if err != nil {
			logger.Error("apply fine failed", "error", err)
			w.WriteHeader(http.StatusConflict)
			return
		}

		if err := httpResponse(w, http.StatusOK, out); err != nil {
			logger.Error("failed to encode response", "error", err)
		}
	}
}

func (h *overdueHandler) waiveFine(ctx context.Context, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Reason string `json:"reason"`
		}

		caseID := overdue.OverdueCaseID(r.PathValue("id"))
		req, err := decode[request](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := h.s.WaiveFine(ctx, overdue.WaiveFineInput{
			CaseID: caseID,
			Reason: req.Reason,
		}); err != nil {
			logger.Error("waive fine failed", "error", err)
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *overdueHandler) resolveCase(ctx context.Context, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		caseID := overdue.OverdueCaseID(r.PathValue("id"))

		if err := h.s.ResolveCase(ctx, overdue.ResolveCaseInput{CaseID: caseID}); err != nil {
			logger.Error("resolve case failed", "error", err)
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
