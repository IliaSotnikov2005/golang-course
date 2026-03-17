package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/IliaSotnikov2005/golang-course/task2/gateway/internal/domain"
	"github.com/IliaSotnikov2005/golang-course/task2/gateway/internal/usecase"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

type Handler struct {
	getRepoUseCase *usecase.GetRepositoryUseCase
}

func NewHandler(getRepoUseCase *usecase.GetRepositoryUseCase) *Handler {
	return &Handler{
		getRepoUseCase: getRepoUseCase,
	}
}

func (h *Handler) Router() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET"},
		AllowedHeaders: []string{"Accept"},
		MaxAge:         300,
	}))

	r.Get("/swagger/*", h.serveSwagger)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/repos/{owner}/{repo}", h.getRepository)
	})

	return r
}

func (h *Handler) getRepository(w http.ResponseWriter, r *http.Request) {
	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")

	repository, err := h.getRepoUseCase.Execute(r.Context(), owner, repo)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response := map[string]any{
		"name": repository.Name,

		"description":      repository.Description,
		"stargazers_count": repository.Stargazers,
		"forks_count":      repository.Forks,
		"created_at":       repository.CreatedAt,
		"html_url":         repository.HTMLURL,
	}

	h.respondJSON(w, http.StatusOK, response)
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	var statusCode int
	var message string

	switch {
	case errors.Is(err, domain.ErrNotFound):
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
		message = "Unauthorized"
	case errors.Is(err, domain.ErrRateLimit):
		statusCode = http.StatusTooManyRequests
		message = "Rate limit exceeded"
	case errors.Is(err, domain.ErrTimeout):
		statusCode = http.StatusGatewayTimeout
		message = "Request timeout"
	default:
		statusCode = http.StatusInternalServerError
		message = "Internal server error"
	}

	h.respondJSON(w, statusCode, map[string]string{
		"error":   http.StatusText(statusCode),
		"message": message,
	})
}

func (h *Handler) respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) serveSwagger(w http.ResponseWriter, r *http.Request) {
	// TODO
}
