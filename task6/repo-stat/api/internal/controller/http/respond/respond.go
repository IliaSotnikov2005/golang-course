package respond

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/IliaSotnikov2005/golang-course/task6/repo-stat/api/internal/domain"
)

type ErrorResponse struct {
	Error string `json:"error" example:"Repository not found"`
}

func Error(w http.ResponseWriter, err error) {
	var statusCode int
	var errorMessage string

	switch {
	case errors.Is(err, domain.ErrNotFound):
		statusCode = http.StatusNotFound
		errorMessage = http.StatusText(http.StatusNotFound)
	case errors.Is(err, domain.ErrMovedPermanently):
		statusCode = http.StatusMovedPermanently
		errorMessage = http.StatusText(http.StatusMovedPermanently)
	case errors.Is(err, domain.ErrInvalidInput):
		statusCode = http.StatusBadRequest
		errorMessage = err.Error()
	case errors.Is(err, domain.ErrForbidden):
		statusCode = http.StatusForbidden
		errorMessage = http.StatusText(http.StatusForbidden)

	case errors.Is(err, domain.ErrUnauthorized):
		statusCode = http.StatusUnauthorized
		errorMessage = http.StatusText(http.StatusUnauthorized)

	case errors.Is(err, domain.ErrRateLimit):
		statusCode = http.StatusTooManyRequests
		errorMessage = http.StatusText(http.StatusTooManyRequests)

	case errors.Is(err, domain.ErrTimeout):
		statusCode = http.StatusGatewayTimeout
		errorMessage = http.StatusText(http.StatusGatewayTimeout)
	case errors.Is(err, domain.ErrAccepted):
		statusCode = http.StatusAccepted
		errorMessage = "Data is being collected. Please try again in a few seconds."
	default:
		statusCode = http.StatusInternalServerError
		errorMessage = http.StatusText(http.StatusInternalServerError)
	}

	JSON(w, statusCode, ErrorResponse{Error: errorMessage})
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
