package domain

import "time"

type Repository struct {
	Name        string
	Description string
	Stargazers  int
	Forks       int
	CreatedAt   time.Time
	HTMLURL     string
}

type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
