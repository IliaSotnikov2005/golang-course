# Distributed Repository Viewer

A microservices-based system for retrieving GitHub repository information. Consists of two services: Collector (gRPC) and Gateway (HTTP REST).

## Architecture
```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Client    │---->│   Gateway   │---->│  Collector  │---->│  GitHub API │
│ (HTTP/REST) │     │  (REST API) │     │   (gRPC)    │     │             │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
      │                   │                   │
      v                   v                   v
 HTTP 200-500        gRPC codes         HTTP 200-500
 (from gateway)    (from collector)    (from github API)
```

### Services

1. **Collector** (gRPC Server)
   - Fetches data from GitHub API
   - Returns repository information
   - Port (default): 50051

2. **Gateway** (HTTP Server)
   - Provides REST API for clients
   - Proxies requests to Collector
   - Generates Swagger documentation
   - Port (default): 8080

## Project Structure
```
.
├── collector/ # Collector service
│ ├── cmd/
│ │ └── main.go # Entry point
│ ├── internal/
│ │ ├── adapter/ # External system adapters
│ │ │ ├── github/ # GitHub API client
│ │ │ └── grpc/ # gRPC handler
│ │ ├── api/ # API contracts
│ │ │ └── proto/ # Protocol Buffers
│ │ ├── app/ # Application
│ │ │ └── grpc/ # gRPC server
│ │ ├── config/ # Configuration
│ │ ├── domain/ # Domain models and errors
│ │ └── usecase/ # Business logic
│ ├── config/ # Configuration files
│ ├── Dockerfile
│ └── go.mod
│
├── gateway/ # Gateway service
│ ├── cmd/
│ │ └── main.go # Entry point
│ ├── internal/
│ │ ├── adapter/ # Adapters
│ │ │ ├── grpc/ # gRPC client
│ │ │ └── rest/ # REST handler
│ │ ├── api/ # API contracts
│ │ │ └── proto/ # Protocol Buffers
│ │ ├── app/ # Application
│ │ │ └── http/ # HTTP server
│ │ ├── config/ # Configuration
│ │ ├── domain/ # Domain models
│ │ └── usecase/ # Business logic
│ ├── docs/ # Swagger documentation
│ ├── config/ # Configuration files
│ ├── Dockerfile
│ └── go.mod
│
├── docker-compose.yaml # Service orchestration
└── README.md
```

## Installation and Running

### Requirements

- Docker and Docker Compose v2
- Go 1.26+

### Run with Docker Compose

1. Clone the repository:
```bash
git clone git@github.com:IliaSotnikov2005/golang-course.git
cd task2
```

2. Configure services:
```bash
collector/config/<your-config.yaml>
gateway/config/<your-config.yaml>
```

3. Start services:
```bash
docker-compose up -d
```

4. Check status:
```bash
docker-compose ps
curl http://localhost:8080/api/v1/health
```

### Usage Examples
Once both services are running:
#### curl requests
```bash
# Get Go repository information
curl http://localhost:8080/api/v1/repos/golang/go

# Get Kubernetes repository
curl http://localhost:8080/api/v1/repos/kubernetes/kubernetes

# Health check
curl http://localhost:8080/api/v1/health

# Error handling - non-existent repository
curl http://localhost:8080/api/v1/repos/golang/nonexistent
```

### Error Mapping

| GitHub API | Domain Error | gRPC Code | HTTP Code |
|------------|--------------|-----------|-----------|
| 404 Not Found | `ErrNotFound` | `NotFound` | 404 |
| 301 Moved Permanently | `ErrMovedPermanently` | `NotFound` | 404 |
| 403 Forbidden | `ErrForbidden` | `PermissionDenied` | 403 |
| 401 Unauthorized | `ErrUnauthorized` | `Unauthenticated` | 401 |
| 429 Rate Limit | `ErrRateLimit` | `ResourceExhausted` | 429 |
| 400 Bad Request | `ErrInvalidInput` | `InvalidArgument` | 400 |
| Timeout | `ErrTimeout` | `DeadlineExceeded` | 504 |
| Other | `ErrInternal` | `Internal` | 500 |
