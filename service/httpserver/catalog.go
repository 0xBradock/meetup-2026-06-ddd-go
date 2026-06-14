package httpserver

import (
	"context"
	"go-svr/catalog"
	"log/slog"
	"net/http"
)

type catalogHandler struct {
	s *catalog.Service
}

func (h *catalogHandler) registerTitle(ctx context.Context, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Name      string `json:"name"`
			ISBN      string `json:"isbn"`
			Author    string `json:"author"`
			Publisher string `json:"publisher"`
			Category  string `json:"category"`
		}

		req, err := decode[request](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		out, err := h.s.RegisterTitle(ctx, catalog.RegisterTitleInput{
			Name:      req.Name,
			ISBN:      req.ISBN,
			Author:    req.Author,
			Publisher: req.Publisher,
			Category:  req.Category,
		})
		if err != nil {
			logger.Error("register title failed", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := httpResponse(w, http.StatusCreated, out); err != nil {
			logger.Error("failed to encode response", "error", err)
		}
	}
}

func (h *catalogHandler) getTitle(ctx context.Context, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := catalog.TitleID(r.PathValue("id"))

		out, err := h.s.GetTitle(ctx, id)
		if err != nil {
			logger.Error("get title failed", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := httpResponse(w, http.StatusOK, out); err != nil {
			logger.Error("failed to encode response", "error", err)
		}
	}
}

func (h *catalogHandler) addCopy(ctx context.Context, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			Condition    string `json:"condition"`
			ShelfSection string `json:"shelfSection"`
			ShelfRow     string `json:"shelfRow"`
		}

		titleID := catalog.TitleID(r.PathValue("id"))
		req, err := decode[request](r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		out, err := h.s.AddCopy(ctx, catalog.AddCopyInput{
			TitleID:      titleID,
			Condition:    req.Condition,
			ShelfSection: req.ShelfSection,
			ShelfRow:     req.ShelfRow,
		})
		if err != nil {
			logger.Error("add copy failed", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := httpResponse(w, http.StatusCreated, out); err != nil {
			logger.Error("failed to encode response", "error", err)
		}
	}
}
