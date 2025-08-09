# Core Service

A Go microservice with health/readiness/version endpoints, structured logging, and container-first CI.

## Features

- **Health Endpoints**: `/healthz`, `/readyz`, `/version`
- **Graceful Shutdown**: Proper signal handling and graceful server shutdown
- **Environment Configuration**: Configuration via environment variables
- **Structured Logging**: Built-in logging with structured output
- **Container Ready**: Dockerfile and docker-compose.yml included
- **CI/CD**: GitHub Actions workflow for testing and building
- **Make Commands**: Convenient Makefile for common tasks

## Quick Start

### Local Development

```bash
# Install dependencies
make tidy

# Run the service
make run

# Test the endpoints
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
curl http://localhost:8080/version
```

### Docker

```bash
# Build the Docker image
make docker-build

# Run with Docker
make docker-run

# Or use docker-compose
docker-compose up
```

## Configuration

The service is configured via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Port to listen on |

## Endpoints

- `GET /healthz` - Health check endpoint
- `GET /readyz` - Readiness check endpoint  
- `GET /version` - Version information

## Development

```bash
# Format code
make fmt

# Run tests
make test

# Run all checks (format, vet, test)
make check

# Clean build artifacts
make clean
```

## Building

```bash
# Build binary
make build

# Build Docker image
make docker-build
```

## Deployment

The service includes:

- Multi-stage Dockerfile for optimized production builds
- GitHub Actions CI/CD pipeline
- Health checks for container orchestration
- Graceful shutdown handling

## License

MIT License - see LICENSE file for details.
