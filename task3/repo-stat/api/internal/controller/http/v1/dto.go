package v1

import "time"

type RepositoryResponse struct {
	Name            string    `json:"name" example:"go"`
	Description     string    `json:"description" example:"The Go programming language"`
	StargazersCount int       `json:"stargazers_count" example:"123456"`
	ForksCount      int       `json:"forks_count" example:"12345"`
	CreatedAt       time.Time `json:"created_at" example:"2014-08-19T22:33:41Z"`
	HTMLURL         string    `json:"html_url" example:"https://github.com/golang/go"`
}

type HealthResponse struct {
	Status   string          `json:"status" example:"ok"`
	Services []ServiceStatus `json:"services"`
}

type ServiceStatus struct {
	Name   string `json:"name" example:"processor"`
	Status string `json:"status" example:"up"`
}

type ErrorResponse struct {
	Error   string `json:"error" example:"Not Found"`
	Message string `json:"message" example:"Repository not found"`
}
