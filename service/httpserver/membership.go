package httpserver

import (
	"context"
	"go-svr/membership"
	"log/slog"
	"net/http"
)

type membershipHandler struct {
	s *membership.Service
}

func (h *membershipHandler) registerMember(ctx context.Context, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Name  string `json:"name"`
			Email string `json:"email"`
			Phone string `json:"phone"`
		}

		req, err := decode[request](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		out, err := h.s.RegisterMember(ctx, membership.RegisterMemberInput{
			Name:  req.Name,
			Email: req.Email,
			Phone: req.Phone,
		})
		if err != nil {
			logger.Error("register member failed", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := httpResponse(w, http.StatusCreated, out); err != nil {
			logger.Error("failed to encode response", "error", err)
		}
	}
}

func (h *membershipHandler) getMember(ctx context.Context, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := membership.MemberID(r.PathValue("id"))

		out, err := h.s.GetMember(ctx, id)
		if err != nil {
			logger.Error("get member failed", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := httpResponse(w, http.StatusOK, out); err != nil {
			logger.Error("failed to encode response", "error", err)
		}
	}
}

func (h *membershipHandler) blockMember(ctx context.Context, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := membership.MemberID(r.PathValue("id"))

		if err := h.s.BlockMember(ctx, membership.BlockMemberInput{MemberID: id}); err != nil {
			logger.Error("block member failed", "error", err)
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *membershipHandler) unblockMember(ctx context.Context, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := membership.MemberID(r.PathValue("id"))

		if err := h.s.UnblockMember(ctx, membership.UnblockMemberInput{MemberID: id}); err != nil {
			logger.Error("unblock member failed", "error", err)
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
