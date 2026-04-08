package domain

type PingStatus string

const (
	PingStatusUp   PingStatus = "up"
	PingStatusDown PingStatus = "down"
)

type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type PingResponse struct {
	Status   string          `json:"status"`
	Services []ServiceStatus `json:"services"`
}
