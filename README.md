# Go Clean Architecture Backend

A Go backend application following clean architecture principles with dependency injection.

## Features

- Clean architecture with clear separation of concerns
- Dependency injection using fx framework
- HTTP handlers with Gin router
- In-memory data storage (JSON-based)
- Comprehensive testing infrastructure

## Project Structure

```
.
├── cmd/                    # Application entry points
├── internal/              # Private application code
│   ├── app/              # Application services
│   └── pkg/              # Shared packages
├── handlers/              # HTTP handlers
├── models/                # Data models and storage
└── server.go             # Server setup
```

## Getting Started

### Prerequisites

- Go 1.23+

### Running Locally

1. Clone the repository
2. Install dependencies: `go mod download`
3. Run the application: `go run server.go`
4. The server will start on `http://localhost:3001`

### API Endpoints

- `GET /health` - Health check
- `POST /api/v1/examples/` - Create example
- `GET /api/v1/examples/` - Get all examples
- `GET /api/v1/examples/:id` - Get example by ID
- `PUT /api/v1/examples/:id` - Update example
- `DELETE /api/v1/examples/:id` - Delete example

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/app/example
```

## Architecture

This application follows clean architecture principles:

1. **Handlers** - Handle HTTP requests and responses
2. **Services** - Contain business logic
3. **Models** - Handle data storage and retrieval
4. **Dependency Injection** - Manages dependencies using fx

## Dependencies

- `go.uber.org/fx` - Dependency injection
- `github.com/gin-gonic/gin` - HTTP router
- `github.com/samber/lo` - Utility functions
- `github.com/stretchr/testify` - Testing utilities
- `github.com/pkg/errors` - Error handling
- `github.com/spf13/viper` - Configuration management
