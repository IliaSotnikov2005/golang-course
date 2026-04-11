package github

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
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

func (c *Client) Exists(ctx context.Context, owner, repo string) (bool, error) {
	url, err := url.JoinPath(c.baseURL, owner, repo)
	if err != nil {
		return false, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return false, err
	}

	request.Header.Set("Accept", "application/vnd.github.v3+json")
	request.Header.Set("User-Agent", c.userAgent)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return false, err
	}

	defer func() {
		io.Copy(io.Discard, response.Body)
		if err := response.Body.Close(); err != nil {
			c.log.Warn("failed to close response body", slog.String("error", err.Error()))
		}
	}()

	switch response.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}
}
