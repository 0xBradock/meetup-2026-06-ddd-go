package httpserver

import (
	"context"
	"go-svr/reservation"
	"log/slog"
	"net/http"
)

type reservationHandler struct {
	s *reservation.Service
}

func (h *reservationHandler) placeReservation(ctx context.Context, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			MemberID string `json:"memberID"`
			TitleID  string `json:"titleID"`
		}

		req, err := decode[request](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		out, err := h.s.PlaceReservation(ctx, reservation.PlaceReservationInput{
			MemberID: reservation.MemberID(req.MemberID),
			TitleID:  reservation.TitleID(req.TitleID),
		})
		if err != nil {
			logger.Error("place reservation failed", "error", err)
			w.WriteHeader(http.StatusConflict)
			return
		}

		if err := httpResponse(w, http.StatusCreated, out); err != nil {
			logger.Error("failed to encode response", "error", err)
		}
	}
}

func (h *reservationHandler) cancelReservation(ctx context.Context, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := reservation.ReservationID(r.PathValue("id"))

		if err := h.s.CancelReservation(ctx, reservation.CancelReservationInput{ReservationID: id}); err != nil {
			logger.Error("cancel reservation failed", "error", err)
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *reservationHandler) getQueue(ctx context.Context, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		titleID := reservation.TitleID(r.PathValue("titleID"))

		out, err := h.s.GetQueue(ctx, titleID)
		if err != nil {
			logger.Error("get queue failed", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := httpResponse(w, http.StatusOK, out); err != nil {
			logger.Error("failed to encode response", "error", err)
		}
	}
}
