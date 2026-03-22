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

	_ "github.com/IliaSotnikov2005/golang-course/task2/gateway/docs"
	httpSwagger "github.com/swaggo/http-swagger"
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
		r.Get("/health", h.healthCheck)
	})

	return r
}

// GetRepository godoc
// @Summary      Get repository information
// @Description  Returns information about a GitHub repository
// @Tags         repositories
// @Accept       json
// @Produce      json
// @Param        owner   path      string  true  "Repository owner (user or organization)"
// @Param        repo    path      string  true  "Repository name"
// @Success      200     {object}  RepositoryResponse
// @Failure      400     {object}  ErrorResponse  "Invalid input"
// @Failure      403     {object}  ErrorResponse  "Access forbidden"
// @Failure      404     {object}  ErrorResponse  "Repository not found"
// @Failure      429     {object}  ErrorResponse  "Rate limit exceeded"
// @Failure      500     {object}  ErrorResponse  "Internal server error"
// @Failure      504     {object}  ErrorResponse  "Request timeout"
// @Router       /api/v1/repos/{owner}/{repo} [get]
func (h *Handler) getRepository(w http.ResponseWriter, r *http.Request) {
	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")

	repository, err := h.getRepoUseCase.Execute(r.Context(), owner, repo)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response := RepositoryResponse{
		Name:            repository.Name,
		Description:     repository.Description,
		StargazersCount: repository.Stargazers,
		ForksCount:      repository.Forks,
		CreatedAt:       repository.CreatedAt,
		HTMLURL:         repository.HTMLURL,
	}

	h.respondJSON(w, http.StatusOK, response)
}

// HealthCheck godoc
// @Summary      Health check
// @Description  Returns service health status
// @Tags         system
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Router       /api/v1/health [get]
func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	h.respondJSON(w, http.StatusOK, HealthResponse{Status: "ok"})
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	var statusCode int
	var message string
	var errorCode string

	switch {
	case errors.Is(err, domain.ErrNotFound):
		statusCode = http.StatusNotFound
		errorCode = "NOT_FOUND"
		message = "Repository not found"

	case errors.Is(err, domain.ErrMovedPermanently):
		statusCode = http.StatusNotFound
		errorCode = "NOT_FOUND"
		message = "Repository not found or moved"

	case errors.Is(err, domain.ErrInvalidInput):
		statusCode = http.StatusBadRequest
		errorCode = "INVALID_INPUT"
		message = err.Error()

	case errors.Is(err, domain.ErrForbidden):
		statusCode = http.StatusForbidden
		errorCode = "FORBIDDEN"
		message = "Access forbidden"

	case errors.Is(err, domain.ErrUnauthorized):
		statusCode = http.StatusUnauthorized
		errorCode = "UNAUTHORIZED"
		message = "Authentication required"

	case errors.Is(err, domain.ErrRateLimit):
		statusCode = http.StatusTooManyRequests
		errorCode = "RATE_LIMIT_EXCEEDED"
		message = "Rate limit exceeded, please try again later"

	case errors.Is(err, domain.ErrTimeout):
		statusCode = http.StatusGatewayTimeout
		errorCode = "TIMEOUT"
		message = "Request timeout"

	default:
		statusCode = http.StatusInternalServerError
		errorCode = "INTERNAL_ERROR"
		message = "Internal server error"
	}

	h.respondJSON(w, statusCode, ErrorResponse{
		Error:   errorCode,
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

func (h *Handler) serveSwagger(w http.ResponseWriter, r *http.Request) {
	httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	).ServeHTTP(w, r)
}

// RepositoryResponse represents the API response for repository info
type RepositoryResponse struct {
	Name            string    `json:"name" example:"go"`
	Description     string    `json:"description" example:"The Go programming language"`
	StargazersCount int       `json:"stargazers_count" example:"123456"`
	ForksCount      int       `json:"forks_count" example:"12345"`
	CreatedAt       time.Time `json:"created_at" example:"2014-08-19T22:33:41Z"`
	HTMLURL         string    `json:"html_url" example:"https://github.com/golang/go"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Not Found"`
	Message string `json:"message" example:"Repository not found"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status string `json:"status" example:"ok"`
}
