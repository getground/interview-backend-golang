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
```

Now let me create the test files:

```go:internal/app/example/service_test.go
package example

import (
	"context"
	"testing"

	"github.com/getground/interview-backend-golang/models"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockExampleRepository struct {
	mock.Mock
}

func (m *MockExampleRepository) Create(ctx context.Context, example *models.ExampleModel) error {
	args := m.Called(ctx, example)
	return args.Error(0)
}

func (m *MockExampleRepository) GetByID(ctx context.Context, id int64) (*models.ExampleModel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ExampleModel), args.Error(1)
}

func (m *MockExampleRepository) GetAll(ctx context.Context) ([]*models.ExampleModel, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.ExampleModel), args.Error(1)
}

func (m *MockExampleRepository) Update(ctx context.Context, example *models.ExampleModel) error {
	args := m.Called(ctx, example)
	return args.Error(0)
}

func (m *MockExampleRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestService_CreateExample(t *testing.T) {
	tests := []struct {
		name          string
		inputName     string
		inputEmail    string
		mockSetup     func(*MockExampleRepository)
		expectedError bool
	}{
		{
			name:       "successful creation",
			inputName:  "John Doe",
			inputEmail: "john@example.com",
			mockSetup: func(repo *MockExampleRepository) {
				repo.On("Create", mock.Anything, mock.AnythingOfType("*models.ExampleModel")).
					Return(nil)
			},
			expectedError: false,
		},
		{
			name:          "empty name",
			inputName:     "",
			inputEmail:    "john@example.com",
			mockSetup:     func(repo *MockExampleRepository) {},
			expectedError: true,
		},
		{
			name:          "empty email",
			inputName:     "John Doe",
			inputEmail:    "",
			mockSetup:     func(repo *MockExampleRepository) {},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockExampleRepository)
			tt.mockSetup(mockRepo)
			
			service := NewService(mockRepo)
			
			result, err := service.CreateExample(context.Background(), tt.inputName, tt.inputEmail)
			
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.inputName, result.Name)
				assert.Equal(t, tt.inputEmail, result.Email)
			}
			
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetExampleByID(t *testing.T) {
	tests := []struct {
		name          string
		inputID       int64
		mockSetup     func(*MockExampleRepository)
		expectedError bool
	}{
		{
			name:    "successful retrieval",
			inputID: 1,
			mockSetup: func(repo *MockExampleRepository) {
				expected := &models.ExampleModel{
					ID:    1,
					Name:  "John Doe",
					Email: "john@example.com",
				}
				repo.On("GetByID", mock.Anything, int64(1)).
					Return(expected, nil)
			},
			expectedError: false,
		},
		{
			name:    "not found",
			inputID: 999,
			mockSetup: func(repo *MockExampleRepository) {
				repo.On("GetByID", mock.Anything, int64(999)).
					Return(nil, errors.New("not found"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockExampleRepository)
			tt.mockSetup(mockRepo)
			
			service := NewService(mockRepo)
			
			result, err := service.GetExampleByID(context.Background(), tt.inputID)
			
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.inputID, result.ID)
			}
			
			mockRepo.AssertExpectations(t)
		})
	}
}
```

```go:handlers/example_handler_test.go
package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getground/interview-backend-golang/internal/app/example"
	"github.com/getground/interview-backend-golang/models"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockExampleService struct {
	mock.Mock
}

func (m *MockExampleService) CreateExample(ctx context.Context, name, email string) (*models.ExampleModel, error) {
	args := m.Called(ctx, name, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ExampleModel), args.Error(1)
}

func (m *MockExampleService) GetExampleByID(ctx context.Context, id int64) (*models.ExampleModel, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ExampleModel), args.Error(1)
}

func (m *MockExampleService) GetAllExamples(ctx context.Context) ([]*models.ExampleModel, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.ExampleModel), args.Error(1)
}

func (m *MockExampleService) UpdateExample(ctx context.Context, id int64, name, email string) (*models.ExampleModel, error) {
	args := m.Called(ctx, id, name, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ExampleModel), args.Error(1)
}

func (m *MockExampleService) DeleteExample(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupTestRouter(handler *ExampleHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	api := router.Group("/api/v1")
	{
		examples := api.Group("/examples")
		{
			examples.POST("/", handler.CreateExample)
			examples.GET("/", handler.GetAllExamples)
			examples.GET("/:id", handler.GetExampleByID)
			examples.PUT("/:id", handler.UpdateExample)
			examples.DELETE("/:id", handler.DeleteExample)
		}
	}
	
	return router
}

func TestExampleHandler_CreateExample(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    CreateExampleRequest
		mockSetup      func(*MockExampleService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful creation",
			requestBody: CreateExampleRequest{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			mockSetup: func(service *MockExampleService) {
				expected := &models.ExampleModel{