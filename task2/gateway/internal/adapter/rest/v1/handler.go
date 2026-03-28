package v1

import (
	"net/http"

	"github.com/IliaSotnikov2005/golang-course/task2/gateway/internal/usecase"
	"github.com/go-chi/chi"
)

type Handler struct {
	getRepoUseCase *usecase.GetRepositoryUseCase
}

func NewHandler(getRepoUseCase *usecase.GetRepositoryUseCase) *Handler {
	return &Handler{
		getRepoUseCase: getRepoUseCase,
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
