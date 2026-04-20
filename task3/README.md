# Repo Stat Service

A GitHub repository statistics monitoring system built with a microservices architecture in Go.
## Architecture

The project follows a clean, decoupled microservices approach:
- API Gateway: The entry point for clients. It provides a REST API, serves Swagger documentation, and orchestrates requests to internal services.

- Collector Service: Responsible for fetching raw data from the GitHub API.

- Processor Service: Processes and transforms data received from the Collector.

- Subscriber Service: subscriber service.

## Getting Started
Prerequisites
- Docker and Docker Compose installed.
- Go 1.26+ (if running locally).

### Quick Start with Docker

From the `task3` directory, run:
```bash
docker compose up --build
```

The services will be available at:
- REST API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html
- Collector (gRPC): localhost:8081
- Processor (gRPC): localhost:8082
- Subscriber (gRPC): localhost:8083

## Testing

The project includes integration tests that run in an isolated Docker container.

Run tests via Docker Compose:
```bash
docker compose up --exit-code-from tests
```

Run tests locally:
```bash
go test -v ./tests/...
```

## Project Structure
```plaintext
.
├── api/             # REST Gateway (Handlers, Usecases)
├── collector/       # GitHub Data Collection Service
├── processor/       # Data Processing Service
├── subscriber/      # Subscription Management Service
├── proto/           # Protobuf definitions & generated code
├── platform/        # Shared packages (logger, grpcserver, etc.)
└── compose.yaml     # Docker orchestration
```
