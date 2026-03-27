package v1

import "time"

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
