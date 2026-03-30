package v1

import (
	"log/slog"
	"net/http"

	_ "github.com/IliaSotnikov2005/golang-course/task3/repo-stat/api/docs"
	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/api/internal/usecase"
)

type Handler struct {
	log                  *slog.Logger
	getRepositoryUseCase *usecase.GetRepositoryUseCase
	pingUseCase          *usecase.PingUseCase
}

func NewHandler(l *slog.Logger, g *usecase.GetRepositoryUseCase, p *usecase.PingUseCase) *Handler {
	return &Handler{
		log:                  l,
		getRepositoryUseCase: g,
		pingUseCase:          p,
	}
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
	url := r.URL.Query().Get("url")
	if url == "" {
		h.respondJSON(w, http.StatusBadRequest, ErrorResponse{Message: "url parameter is required"})
		return
	}

	repository, err := h.getRepositoryUseCase.Execute(r.Context(), url)
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

func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	res, isOk := h.pingUseCase.Execute(r.Context())

	status := http.StatusOK
	if !isOk {
		status = http.StatusServiceUnavailable
	}

	h.respondJSON(w, status, res)
}
