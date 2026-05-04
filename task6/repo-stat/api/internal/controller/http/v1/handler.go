package v1

import (
	"log/slog"
	"net/http"
	"strings"

	_ "github.com/IliaSotnikov2005/golang-course/task6/repo-stat/api/docs"
	"github.com/IliaSotnikov2005/golang-course/task6/repo-stat/api/internal/usecase"
	"github.com/go-chi/chi"
)

type Handler struct {
	log                      *slog.Logger
	getRepositoryUseCase     *usecase.GetRepositoryUseCase
	subscribeUseCase         *usecase.SubscribeUseCase
	unsubscribeUseCase       *usecase.UnsubscribeUseCase
	listSubscriptionsUseCase *usecase.ListSubscriptionsUseCase
	getSubsInfoUseCase       *usecase.GetSubscriptionsInfoUseCase
	pingUseCase              *usecase.PingUseCase
}

func NewHandler(
	log *slog.Logger,
	g *usecase.GetRepositoryUseCase,
	s *usecase.SubscribeUseCase,
	u *usecase.UnsubscribeUseCase,
	l *usecase.ListSubscriptionsUseCase,
	gs *usecase.GetSubscriptionsInfoUseCase,
	p *usecase.PingUseCase) *Handler {
	return &Handler{
		log:                      log,
		getRepositoryUseCase:     g,
		subscribeUseCase:         s,
		unsubscribeUseCase:       u,
		listSubscriptionsUseCase: l,
		getSubsInfoUseCase:       gs,
		pingUseCase:              p,
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

	repository, err := h.getRepositoryUseCase.Execute(r.Context(), rawURL)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response := RepositoryResponse{
		FullName:        repository.FullName,
		Description:     repository.Description,
		StargazersCount: repository.Stargazers,
		ForksCount:      repository.Forks,
		CreatedAt:       repository.CreatedAt,
		HTMLURL:         repository.HTMLURL,
	}

	h.respondJSON(w, http.StatusOK, response)
}

// subscribe godoc
// @Summary      Subscribe to a repository
// @Tags         subscriptions
// @Param        owner  path  string  true  "Repository Owner"
// @Param        repo   path  string  true  "Repository Name"
// @Success      201    {object}  map[string]string
// @Router       /v1/subscriptions/{owner}/{repo} [post]
func (h *Handler) subscribe(w http.ResponseWriter, r *http.Request) {
	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")

	if err := h.subscribeUseCase.Execute(r.Context(), owner, repo); err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusCreated, map[string]string{"status": "subscribed"})

}

// unsubscribe godoc
// @Summary      Unsubscribe from a repository
// @Tags         subscriptions
// @Param        owner  path  string  true  "Repository Owner"
// @Param        repo   path  string  true  "Repository Name"
// @Success      200    {object}  map[string]string
// @Router       /v1/subscriptions/{owner}/{repo} [delete]
func (h *Handler) unsubscribe(w http.ResponseWriter, r *http.Request) {
	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")

	if err := h.unsubscribeUseCase.Execute(r.Context(), owner, repo); err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"status": "unsubscribed"})
}

// listSubscriptions godoc
// @Summary      List all subscriptions
// @Tags         subscriptions
// @Success      200  {array}  v1.SubscriptionResponse
// @Router       /v1/subscriptions [get]
func (h *Handler) listSubscriptions(w http.ResponseWriter, r *http.Request) {
	subs, err := h.listSubscriptionsUseCase.Execute(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	response := make([]SubscriptionResponse, 0, len(subs))
	for _, s := range subs {
		response = append(response, SubscriptionResponse{
			Owner: s.Owner,
			Repo:  s.Repo,
		})
	}

	h.respondJSON(w, http.StatusOK, response)
}

// getSubscriptionsInfo godoc
// @Summary      Get detailed info about all subscribed repositories
// @Description  Triggers a chain API -> Processor -> Collector (which fetches list from Subscriber)
// @Tags         subscriptions
// @Produce      json
// @Success      200  {array}   v1.RepositoryResponse
// @Failure      500  {object}  v1.ErrorResponse
// @Router       /v1/subscriptions/info [get]
func (h *Handler) getSubscriptionsInfo(w http.ResponseWriter, r *http.Request) {
	repositories, err := h.getSubsInfoUseCase.Execute(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	response := make([]RepositoryResponse, 0, len(repositories))
	for _, repo := range repositories {
		response = append(response, RepositoryResponse{
			FullName:        repo.FullName,
			Description:     repo.Description,
			StargazersCount: repo.Stargazers,
			ForksCount:      repo.Forks,
			CreatedAt:       repo.CreatedAt,
			HTMLURL:         repo.HTMLURL,
		})
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
