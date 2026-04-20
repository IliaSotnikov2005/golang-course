package v1

import "time"

// RepositoryResponse represents a successful response
type RepositoryResponse struct {
	Name            string    `json:"full_name" example:"go"`
	Description     string    `json:"description" example:"The Go programming language"`
	StargazersCount int       `json:"stars" example:"123456"`
	ForksCount      int       `json:"forks" example:"12345"`
	CreatedAt       time.Time `json:"created_at" example:"2014-08-19T22:33:41Z"`
	HTMLURL         string    `json:"html_url" example:"https://github.com/golang/go"`
}

// HealthResponse represents a services status response
type HealthResponse struct {
	Status   string          `json:"status" example:"ok"`
	Services []ServiceStatus `json:"services"`
}

// ServiceStatus represents a service status
type ServiceStatus struct {
	Name   string `json:"name" example:"processor"`
	Status string `json:"status" example:"up"`
}

// ErrorResponse represents a error response
type ErrorResponse struct {
	Error string `json:"error" example:"Repository not found"`
}
