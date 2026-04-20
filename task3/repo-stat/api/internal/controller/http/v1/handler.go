package v1

import (
	"log/slog"
	"net/http"
	"strings"

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
// @Summary      Gets information about GitHub repository
// @Description  Returns data about stars, forks, and the repository description by its URL
// @Tags         repositories
// @Accept       json
// @Produce      json
// @Param        url   query     string  true  "GitHub URL (e.g. https://github.com/google/go-github)"
// @Success      200   {object}  v1.RepositoryResponse
// @Failure      400   {object}  v1.ErrorResponse
// @Failure      404   {object}  v1.ErrorResponse
// @Failure      500   {object}  v1.ErrorResponse
// @Router       /v1/repositories/info [get]
func (h *Handler) getRepository(w http.ResponseWriter, r *http.Request) {
	rawURL := strings.TrimSpace(r.URL.Query().Get("url"))
	if rawURL == "" {
		h.respondJSON(w, http.StatusBadRequest, ErrorResponse{Message: "url parameter is required"})
		return
	}

	repository, err := h.getRepositoryUseCase.Execute(r.Context(), rawURL)
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

// healthCheck godoc
// @Summary      Checking the service status (Healthcheck)
// @Description  Returns API status and availability of dependent microservices (Collector, Subscriber)
// @Tags         system
// @Produce      json
// @Success      200  {object}  HealthResponse  "All systems are working fine"
// @Failure      503  {object}  HealthResponse  "One or more services are unavailable"
// @Router       /ping [get]
func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	res, isOk := h.pingUseCase.Execute(r.Context())

	status := http.StatusOK
	if !isOk {
		status = http.StatusServiceUnavailable
	}

	h.respondJSON(w, status, res)
}
