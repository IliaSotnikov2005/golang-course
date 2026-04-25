#file:repo-stat взаимодействие между #file:processor и #file:collector должно быть через kafka. Но запрос по /v1/subscriptions/info всегда возвращает "error": "Data is being collected. Please try again in a few seconds."

# Repo Stat Service

A GitHub repository statistics monitoring system built with a Go microservices architecture, gRPC, REST, and PostgreSQL.
## Architecture

The project follows Clean Architecture principles and consists of the following components:

- API Gateway (Port 8080): Acts as the entry point. It serves the REST API, provides Swagger UI, and maps internal gRPC errors to standard HTTP statuses.

- Subscriber Service (Port 8083): Manages repository subscriptions in PostgreSQL. It validates repository existence via the GitHub API.

- Processor Service (Port 8082): Aggregates repository data by orchestrating requests to the Collector.

- Collector Service (Port 8081): A dedicated layer for fetching raw data from the external GitHub API.

- PostgreSQL: The relational database used by the Subscriber Service to persist data.

## Getting Started
Prerequisites
- Docker and Docker Compose installed.
- Go 1.26+ (if running locally).

## Configuration
Create a .env file in the root directory based on the provided example:
```bash
cp .env.example .env
```

### Quick Start with Docker

From the `task4` directory, run:
```bash
docker compose up --build
```

The services will be available at:
- REST API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html
- Collector (gRPC): localhost:8081
- Processor (gRPC): localhost:8082
- Subscriber (gRPC): localhost:8083

## Database & Migrations
Migrations are applied automatically when the subscriber container starts.
Migration files: Located in ./subscriber/migrations

To completely reset the database (including all data):
```bash
docker-compose down -v
```

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
├── api/             # REST Gateway (Handlers, gRPC Clients)
├── collector/       # GitHub Data Collection Service
├── processor/       # Data Processing Service
├── subscriber/      # Subscription Management Service
│   └── migrations/  # SQL Migration files
├── proto/           # Protobuf definitions & generated code
├── platform/        # Shared packages
├── .env.example     # Template for environment variables
└── compose.yaml     # Docker orchestration
```
