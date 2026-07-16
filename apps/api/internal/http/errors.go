package httpapi

import (
	"encoding/json"
	"log"
	"net/http"
)

type apiError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

type errorEnvelope struct {
	Error apiError `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write json response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorEnvelope{Error: apiError{Code: code, Message: message}})
}

// serverError logs the real cause and returns an opaque 500 to the client.
func serverError(w http.ResponseWriter, what string, err error) {
	log.Printf("error: %s: %v", what, err)
	writeError(w, http.StatusInternalServerError, "internal", "internal server error")
}