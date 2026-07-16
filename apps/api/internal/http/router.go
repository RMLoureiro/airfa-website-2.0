package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/RMLoureiro/airfa-website-2.0/apps/api/internal/db"
)

// Server holds the dependencies shared by the handlers.
type Server struct {
	q *db.Queries
}

// Router builds the full HTTP surface: /healthz plus the public /v1 reads.
func Router(q *db.Queries) http.Handler {
	s := &Server{q: q}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/v1", func(r chi.Router) {
		r.Get("/site", s.getSite)
		r.Get("/pages/{slug}", s.getPage)
		r.Get("/collections/{name}", s.getCollection)
	})

	r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		writeError(w, http.StatusNotFound, "not_found", "resource not found")
	})

	return r
}