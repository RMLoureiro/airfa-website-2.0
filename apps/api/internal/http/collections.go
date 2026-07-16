package httpapi

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/RMLoureiro/airfa-website-2.0/apps/api/internal/db"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

// getCollection handles GET /v1/collections/{name}?limit=&offset=.
func (s *Server) getCollection(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit, offset, err := pagination(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_pagination", err.Error())
		return
	}

	var (
		items any
		total int64
	)

	switch name := chi.URLParam(r, "name"); name {
	case "events":
		rows, err := s.q.ListPublishedEvents(ctx, db.ListPublishedEventsParams{Limit: limit, Offset: offset})
		if err != nil {
			serverError(w, "list events", err)
			return
		}
		if rows == nil {
			rows = []db.ListPublishedEventsRow{}
		}
		if total, err = s.q.CountPublishedEvents(ctx); err != nil {
			serverError(w, "count events", err)
			return
		}
		items = rows

	case "posts":
		rows, err := s.q.ListPublishedPosts(ctx, db.ListPublishedPostsParams{Limit: limit, Offset: offset})
		if err != nil {
			serverError(w, "list posts", err)
			return
		}
		if rows == nil {
			rows = []db.ListPublishedPostsRow{}
		}
		if total, err = s.q.CountPublishedPosts(ctx); err != nil {
			serverError(w, "count posts", err)
			return
		}
		items = rows

	case "partners":
		rows, err := s.q.ListPartners(ctx, db.ListPartnersParams{Limit: limit, Offset: offset})
		if err != nil {
			serverError(w, "list partners", err)
			return
		}
		// Note: []db.Partner, not []db.ListPartnersRow — ListPartners selects every
		// column of the table, so sqlc reuses the Partner model instead of
		// generating a row type. There is no ListPartnersRow.
		if rows == nil {
			rows = []db.Partner{}
		}
		if total, err = s.q.CountPartners(ctx); err != nil {
			serverError(w, "count partners", err)
			return
		}
		items = rows

	case "activities":
		rows, err := s.q.ListPublishedActivities(ctx, db.ListPublishedActivitiesParams{Limit: limit, Offset: offset})
		if err != nil {
			serverError(w, "list activities", err)
			return
		}
		if rows == nil {
			rows = []db.ListPublishedActivitiesRow{}
		}
		if total, err = s.q.CountPublishedActivities(ctx); err != nil {
			serverError(w, "count activities", err)
			return
		}
		items = rows

	default:
		writeError(w, http.StatusNotFound, "unknown_collection",
			fmt.Sprintf("unknown collection: %q", name))
		return
	}

	w.Header().Set("X-Total-Count", strconv.FormatInt(total, 10))
	writeJSON(w, http.StatusOK, items)
}

// pagination reads ?limit=&offset= with defaults and an upper bound.
func pagination(r *http.Request) (limit, offset int32, err error) {
	limit, offset = defaultLimit, 0

	if v := r.URL.Query().Get("limit"); v != "" {
		n, convErr := strconv.Atoi(v)
		if convErr != nil || n < 1 {
			return 0, 0, fmt.Errorf("invalid limit parameter: %q", v)
		}
		if n > maxLimit {
			n = maxLimit
		}
		limit = int32(n)
	}

	if v := r.URL.Query().Get("offset"); v != "" {
		n, convErr := strconv.Atoi(v)
		if convErr != nil || n < 0 {
			return 0, 0, fmt.Errorf("invalid offset parameter: %q", v)
		}
		offset = int32(n)
	}

	return limit, offset, nil
}