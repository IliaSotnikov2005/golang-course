# Repo Stat Service

A GitHub repository statistics monitoring system built with a Go microservices architecture, gRPC, REST, Kafka, and PostgreSQL.
## Architecture

The project follows Clean Architecture principles and consists of the following components:

- API Gateway (Port 8080): The entry point for clients. Provides a REST API and Swagger UI, translating HTTP requests into gRPC calls.

- Subscriber Service (Port 8083): Manages repository subscriptions. Persists data in its own PostgreSQL instance and publishes "subscription created" events to Kafka.

- Processor Service (Port 8082): Maintains a "Smart Cache" in its own PostgreSQL instance. It consumes results from Kafka to store repository statistics and serves them via gRPC.

- Collector Service (Port 8081): A worker service that listens to Kafka for new tasks, fetches data from the GitHub API, and publishes the results back to Kafka.

- PostgreSQL: The relational database used by the Subscriber Service to persist data.
- Kafka: Provides decoupled, reliable data transfer between the services.

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
│   └── migrations/  # SQL Migration files
├── subscriber/      # Subscription Management Service
│   └── migrations/  # SQL Migration files
├── proto/           # Protobuf definitions & generated code
├── platform/        # Shared packages
├── .env.example     # Template for environment variables
└── compose.yaml     # Docker orchestration
```
