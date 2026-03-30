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
	case errors.Is(err, domain.ErrNotFound), errors.Is(err, domain.ErrMovedPermanently):
		statusCode = http.StatusNotFound
		message = "Repository not found"

	case errors.Is(err, domain.ErrInvalidInput):
		statusCode = http.StatusBadRequest
		message = err.Error()

	case errors.Is(err, domain.ErrForbidden):
		statusCode = http.StatusForbidden
		message = "Access forbidden"

	case errors.Is(err, domain.ErrUnauthorized):
		statusCode = http.StatusUnauthorized
		message = "Authentication required"

	case errors.Is(err, domain.ErrRateLimit):
		statusCode = http.StatusTooManyRequests
		message = "Rate limit exceeded, please try again later"

	case errors.Is(err, domain.ErrTimeout):
		statusCode = http.StatusGatewayTimeout
		message = "Request timeout"

	default:
		statusCode = http.StatusInternalServerError
		message = "Internal server error"
	}

	h.respondJSON(w, statusCode, ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	})
}

func (h *Handler) respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
