package dtokafka

import "time"

type RepoRequest struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

type RepoResponse struct {
	FullName    string    `json:"full_name"`
	Description string    `json:"description"`
	Stargazers  int       `json:"stargazers"`
	Forks       int       `json:"forks"`
	CreatedAt   time.Time `json:"created_at"`
	HTMLURL     string    `json:"html_url"`

	Error string `json:"error,omitempty"`
}
