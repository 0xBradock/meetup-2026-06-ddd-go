package httpserver

import (
	"context"
	"go-svr/loan"
	"log/slog"
	"net/http"
)

type loanHandler struct {
	s *loan.Service
}

func (h *loanHandler) borrowBook(ctx context.Context, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			MemberID string `json:"memberID"`
			CopyID   string `json:"copyID"`
		}

		req, err := decode[request](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		out, err := h.s.BorrowBook(ctx, loan.BorrowBookInput{
			MemberID: loan.MemberID(req.MemberID),
			CopyID:   loan.CopyID(req.CopyID),
		})
		if err != nil {
			logger.Error("borrow book failed", "error", err)
			w.WriteHeader(http.StatusConflict)
			return
		}

		if err := httpResponse(w, http.StatusCreated, out); err != nil {
			logger.Error("failed to encode response", "error", err)
		}
	}
}

func (h *loanHandler) returnCopy(ctx context.Context, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loanID := loan.LoanID(r.PathValue("id"))

		out, err := h.s.ReturnCopy(ctx, loan.ReturnCopyInput{LoanID: loanID})
		if err != nil {
			logger.Error("return copy failed", "error", err)
			w.WriteHeader(http.StatusConflict)
			return
		}

		if err := httpResponse(w, http.StatusOK, out); err != nil {
			logger.Error("failed to encode response", "error", err)
		}
	}
}

func (h *loanHandler) extendLoan(ctx context.Context, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			ExtendDays int `json:"extendDays"`
		}

		loanID := loan.LoanID(r.PathValue("id"))
		req, err := decode[request](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		out, err := h.s.ExtendLoan(ctx, loan.ExtendLoanInput{
			LoanID:     loanID,
			ExtendDays: req.ExtendDays,
		})
		if err != nil {
			logger.Error("extend loan failed", "error", err)
			w.WriteHeader(http.StatusConflict)
			return
		}

		if err := httpResponse(w, http.StatusOK, out); err != nil {
			logger.Error("failed to encode response", "error", err)
		}
	}
}

func (h *loanHandler) getLoan(ctx context.Context, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		loanID := loan.LoanID(r.PathValue("id"))

		out, err := h.s.GetLoan(ctx, loanID)
		if err != nil {
			logger.Error("get loan failed", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := httpResponse(w, http.StatusOK, out); err != nil {
			logger.Error("failed to encode response", "error", err)
		}
	}
}
