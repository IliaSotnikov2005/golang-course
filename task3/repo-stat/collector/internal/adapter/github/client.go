package github

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/IliaSotnikov2005/golang-course/task3/repo-stat/collector/internal/domain"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	userAgent  string
	log        *slog.Logger
}

func NewClient(httpClient *http.Client, baseURL string, userAgent string, log *slog.Logger) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL:    baseURL,
		userAgent:  userAgent,
		log:        log,
	}
}

func (c *Client) GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	const operation = "github.Client.GetRepository"

	log := c.log.With(slog.String("operation", operation))

	url := fmt.Sprintf("%s/%s/%s", c.baseURL, owner, repo)

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("request creation error: %w", err)
	}

	request.Header.Set("Accept", "application/vnd.github.v3+json")
	request.Header.Set("User-Agent", c.userAgent)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request execution error: %w", err)
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Warn("failed to close response body: %v\n", slog.String("error", err.Error()))
		}
	}()

	switch response.StatusCode {
	case http.StatusOK:

	case http.StatusMovedPermanently:
		location := response.Header.Get("Location")
		return nil, fmt.Errorf("%w: %s", domain.ErrMovedPermanently, location)

	case http.StatusUnauthorized:
		return nil, fmt.Errorf("%w: invalid or missing token", domain.ErrUnauthorized)

	case http.StatusForbidden:
		return nil, fmt.Errorf("%w: access denied", domain.ErrForbidden)

	case http.StatusNotFound:
		return nil, fmt.Errorf("%w: repository %s/%s not found", domain.ErrNotFound, owner, repo)
	default:
		return nil, fmt.Errorf("%w: unexpected status code %d", domain.ErrInternal, response.StatusCode)
	}

	var ghResponse githubResponse
	if err := json.NewDecoder(response.Body).Decode(&ghResponse); err != nil {
		return nil, fmt.Errorf("json decoding error: %w", err)
	}

	return &domain.Repository{
		FullName:    ghResponse.FullName,
		Description: ghResponse.Description,
		Stargazers:  ghResponse.Stargazers,
		Forks:       ghResponse.Forks,
		CreatedAt:   ghResponse.CreatedAt,
		HTMLURL:     ghResponse.HTMLURL,
	}, nil
}

type githubResponse struct {
	FullName    string    `json:"full_name"`
	Description string    `json:"description"`
	Stargazers  int       `json:"stargazers_count"`
	Forks       int       `json:"forks_count"`
	CreatedAt   time.Time `json:"created_at"`
	HTMLURL     string    `json:"html_url"`
}
