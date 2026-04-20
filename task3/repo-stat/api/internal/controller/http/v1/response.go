package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/api/internal/domain"
)

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	var statusCode int
	var message string

	switch {
	case errors.Is(err, domain.ErrNotFound):
		statusCode = http.StatusNotFound
		message = http.StatusText(http.StatusNotFound)
	case errors.Is(err, domain.ErrMovedPermanently):
		statusCode = http.StatusMovedPermanently
		message = http.StatusText(http.StatusMovedPermanently)
	case errors.Is(err, domain.ErrInvalidInput):
		statusCode = http.StatusBadRequest
		message = err.Error()
	case errors.Is(err, domain.ErrForbidden):
		statusCode = http.StatusForbidden
		message = http.StatusText(http.StatusForbidden)

	case errors.Is(err, domain.ErrUnauthorized):
		statusCode = http.StatusUnauthorized
		message = http.StatusText(http.StatusUnauthorized)

	case errors.Is(err, domain.ErrRateLimit):
		statusCode = http.StatusTooManyRequests
		message = http.StatusText(http.StatusTooManyRequests)

	case errors.Is(err, domain.ErrTimeout):
		statusCode = http.StatusGatewayTimeout
		message = http.StatusText(http.StatusGatewayTimeout)
	default:
		statusCode = http.StatusInternalServerError
		message = http.StatusText(http.StatusInternalServerError)
	}

	h.respondJSON(w, statusCode, ErrorResponse{Message: message})
}

func (h *Handler) respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
