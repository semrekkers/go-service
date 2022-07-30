// Package apih contains helpers for HTTP API services.
package apih

import (
	"encoding/json"
	"net/http"
	"strings"
)

type (
	HandlerFunc func(w http.ResponseWriter, r *http.Request) error
	Middleware  func(next http.Handler) http.Handler
)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if f == nil {
		ServeError(w, &Error{
			StatusCode: http.StatusNotImplemented,
			Message:    "Handler is not implemented",
		})
		return
	}

	if err := f(w, r); err != nil {
		ServeError(w, err)
	}
}

func (f HandlerFunc) HttpHandlerFunc() http.HandlerFunc {
	return f.ServeHTTP
}

func ServeJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(v) // TODO: handle error
}

func DecodeJSON(r *http.Request, v any) error {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		return &Error{
			StatusCode: http.StatusUnsupportedMediaType,
			Message:    "Unsupported request media type, expected JSON content",
		}
	}
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return &Error{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid request, unable to parse JSON content",
			Inner:      err,
		}
	}
	return nil
}
