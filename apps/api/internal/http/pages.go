package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

// getPage handles GET /v1/pages/{slug}. Only published pages are visible;
// an unpublished or missing slug is a 404 either way.
func (s *Server) getPage(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	page, err := s.q.GetPublishedPageBySlug(r.Context(), slug)
	if errors.Is(err, pgx.ErrNoRows) {
		writeError(w, http.StatusNotFound, "page_not_found", "page not found")
		return
	}
	if err != nil {
		serverError(w, "get published page", err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"slug":   page.Slug,
		"title":  page.Title,
		"seo":    json.RawMessage(page.Seo),
		"blocks": json.RawMessage(page.Blocks),
	})
}