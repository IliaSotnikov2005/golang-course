package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/IliaSotnikov2005/golang-course/task4/repo-stat/collector/internal/domain"
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
	url, err := url.JoinPath(c.baseURL, owner, repo)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("request creation error: %w", err)
	}

	request.Header.Set("Accept", "application/vnd.github.v3+json")
	request.Header.Set("User-Agent", c.userAgent)

	response, err := c.httpClient.Do(request)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("%w: %v", domain.ErrTimeout, err)
		}
		if errors.Is(err, context.Canceled) {
			return nil, fmt.Errorf("%w: request canceled", err)
		}

		return nil, fmt.Errorf("request execution error: %w", err)
	}

	defer func() {
		if _, err := io.Copy(io.Discard, response.Body); err != nil {
			c.log.Warn("failed to drain response body", "error", err)
		}
		if err := response.Body.Close(); err != nil {
			c.log.Warn("failed to close response body", slog.String("error", err.Error()))
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
