package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
)

// menuZones are always present in the response, empty when unset, so the
// frontend can read menus.main without a nil check.
var menuZones = []string{"main", "secondary", "utility", "footer"}

// getSite handles GET /v1/site.
func (s *Server) getSite(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	settings := json.RawMessage(`{}`)
	data, err := s.q.GetSettings(ctx)
	switch {
	case err == nil:
		settings = json.RawMessage(data)
	case errors.Is(err, pgx.ErrNoRows):
		// No settings row seeded yet — an empty object is a valid answer.
	default:
		serverError(w, "get settings", err)
		return
	}

	rows, err := s.q.ListMenus(ctx)
	if err != nil {
		serverError(w, "list menus", err)
		return
	}

	menus := make(map[string]json.RawMessage, len(menuZones))
	for _, zone := range menuZones {
		menus[zone] = json.RawMessage(`[]`)
	}
	for _, m := range rows {
		menus[m.Zone] = json.RawMessage(m.Items)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"settings": settings,
		"menus":    menus,
	})
}