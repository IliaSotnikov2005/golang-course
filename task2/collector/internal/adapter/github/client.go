package github

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/IliaSotnikov2005/golang-course/task2/collector/internal/domain"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	logger     log.Logger
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    "https://api.github.com/repos",
		logger:     *log.Default(),
	}
}

type gitHubResponse struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Stargazers  int       `json:"stargazers_count"`
	Forks       int       `json:"forks_count"`
	CreatedAt   time.Time `json:"created_at"`
	HTMLURL     string    `json:"html_url"`
}

func (c *Client) GetRepository(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	apiURL := fmt.Sprintf("%s/%s/%s", c.baseURL, owner, repo)

	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("request creation error: %w", err)
	}

	request.Header.Set("Accept", "application/vnd.github.v3+json")
	request.Header.Set("User-Agent", "Go-Client")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request execution error: %w", err)
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			c.logger.Printf("Warning: failed to close response body: %v\n", err)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", response.Status)
	}

	var ghResponse gitHubResponse
	if err := json.NewDecoder(response.Body).Decode(&ghResponse); err != nil {
		return nil, fmt.Errorf("json decoding error: %w", err)
	}

	return &domain.Repository{
		Name:        ghResponse.Name,
		Description: ghResponse.Description,
		Stargazers:  ghResponse.Stargazers,
		Forks:       ghResponse.Forks,
		CreatedAt:   ghResponse.CreatedAt,
		HTMLURL:     ghResponse.HTMLURL,
	}, nil
}
