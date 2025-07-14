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
					ID:    1,
					Name:  "John Doe",
					Email: "john@example.com",
				}
				service.On("CreateExample", mock.Anything, "John Doe", "john@example.com").
					Return(expected, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"id": 1,
				"name": "John Doe",
				"email": "john@example.com",
				"created_at": expected.CreatedAt,
				"updated_at": expected.UpdatedAt,
			},
		},
		{
			name: "empty name",
			requestBody: CreateExampleRequest{
				Name:  "",
				Email: "john@example.com",
			},
			mockSetup: func(service *MockExampleService) {
				service.On("CreateExample", mock.Anything, "", "john@example.com").
					Return(nil, errors.New("name is required"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid request body",
			},
		},
		{
			name: "empty email",
			requestBody: CreateExampleRequest{
				Name:  "John Doe",
				Email: "",
			},
			mockSetup: func(service *MockExampleService) {
				service.On("CreateExample", mock.Anything, "John Doe", "").
					Return(nil, errors.New("email is required"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid request body",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockExampleService)
			tt.mockSetup(mockService)
			
			handler := NewExampleHandler(mockService)
			router := setupTestRouter(handler)
			
			requestBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/examples/", bytes.NewReader(requestBody))
			req.Header.Set("Content-Type", "application/json")
			
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			
			assert.Equal(t, tt.expectedStatus, resp.Code)
			assert.Equal(t, tt.expectedBody, resp.Body.Body)
			
			mockService.AssertExpectations(t)
		})
	}
}

func TestExampleHandler_GetExampleByID(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockSetup      func(*MockExampleService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful retrieval",
			id:   "1",
			mockSetup: func(service *MockExampleService) {
				expected := &models.ExampleModel{
					ID:    1,
					Name:  "John Doe",
					Email: "john@example.com",
				}
				service.On("GetExampleByID", mock.Anything, int64(1)).
					Return(expected, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id": 1,
				"name": "John Doe",
				"email": "john@example.com",
				"created_at": expected.CreatedAt,
				"updated_at": expected.UpdatedAt,
			},
		},
		{
			name: "invalid ID",
			id:   "invalid",
			mockSetup: func(service *MockExampleService) {
				// No mock setup needed for invalid ID as it will return 400 directly
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid ID parameter",
			},
		},
		{
			name: "not found",
			id:   "999",
			mockSetup: func(service *MockExampleService) {
				service.On("GetExampleByID", mock.Anything, int64(999)).
					Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "Example not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockExampleService)
			tt.mockSetup(mockService)
			
			handler := NewExampleHandler(mockService)
			router := setupTestRouter(handler)
			
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/examples/"+tt.id, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			
			assert.Equal(t, tt.expectedStatus, resp.Code)
			assert.Equal(t, tt.expectedBody, resp.Body.Body)
			
			mockService.AssertExpectations(t)
		})
	}
}

func TestExampleHandler_GetAllExamples(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockExampleService)
		expectedStatus int
		expectedBody   []map[string]interface{}
	}{
		{
			name: "successful retrieval",
			mockSetup: func(service *MockExampleService) {
				examples := []*models.ExampleModel{
					{ID: 1, Name: "John Doe", Email: "john@example.com", CreatedAt: "2023-10-27T10:00:00Z", UpdatedAt: "2023-10-27T10:00:00Z"},
					{ID: 2, Name: "Jane Doe", Email: "jane@example.com", CreatedAt: "2023-10-27T11:00:00Z", UpdatedAt: "2023-10-27T11:00:00Z"},
				}
				service.On("GetAllExamples", mock.Anything).
					Return(examples, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: []map[string]interface{}{
				{
					"id": 1,
					"name": "John Doe",
					"email": "john@example.com",
					"created_at": "2023-10-27T10:00:00Z",
					"updated_at": "2023-10-27T10:00:00Z",
				},
				{
					"id": 2,
					"name": "Jane Doe",
					"email": "jane@example.com",
					"created_at": "2023-10-27T11:00:00Z",
					"updated_at": "2023-10-27T11:00:00Z",
				},
			},
		},
		{
			name: "no examples",
			mockSetup: func(service *MockExampleService) {
				service.On("GetAllExamples", mock.Anything).
					Return([]*models.ExampleModel{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: []map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockExampleService)
			tt.mockSetup(mockService)
			
			handler := NewExampleHandler(mockService)
			router := setupTestRouter(handler)
			
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/examples/", nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			
			assert.Equal(t, tt.expectedStatus, resp.Code)
			assert.Equal(t, tt.expectedBody, resp.Body.Body)
			
			mockService.AssertExpectations(t)
		})
	}
}

func TestExampleHandler_UpdateExample(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		requestBody    UpdateExampleRequest
		mockSetup      func(*MockExampleService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful update",
			id:   "1",
			requestBody: UpdateExampleRequest{
				Name:  "John Doe Updated",
				Email: "john.updated@example.com",
			},
			mockSetup: func(service *MockExampleService) {
				existing := &models.ExampleModel{
					ID:    1,
					Name:  "John Doe",
					Email: "john@example.com",
				}
				service.On("GetExampleByID", mock.Anything, int64(1)).
					Return(existing, nil)
				service.On("UpdateExample", mock.Anything, int64(1), "John Doe Updated", "john.updated@example.com").
					Return(existing, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id": 1,
				"name": "John Doe Updated",
				"email": "john.updated@example.com",
				"created_at": existing.CreatedAt,
				"updated_at": existing.UpdatedAt,
			},
		},
		{
			name: "invalid ID",
			id:   "invalid",
			requestBody: UpdateExampleRequest{
				Name:  "John Doe Updated",
				Email: "john.updated@example.com",
			},
			mockSetup: func(service *MockExampleService) {
				// No mock setup needed for invalid ID as it will return 400 directly
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid ID parameter",
			},
		},
		{
			name: "not found",
			id:   "999",
			requestBody: UpdateExampleRequest{
				Name:  "John Doe Updated",
				Email: "john.updated@example.com",
			},
			mockSetup: func(service *MockExampleService) {
				service.On("GetExampleByID", mock.Anything, int64(999)).
					Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "Example not found",
			},
		},
		{
			name: "empty name",
			id:   "1",
			requestBody: UpdateExampleRequest{
				Name:  "",
				Email: "john.updated@example.com",
			},
			mockSetup: func(service *MockExampleService) {
				existing := &models.ExampleModel{
					ID:    1,
					Name:  "John Doe",
					Email: "john@example.com",
				}
				service.On("GetExampleByID", mock.Anything, int64(1)).
					Return(existing, nil)
				service.On("UpdateExample", mock.Anything, int64(1), "", "john.updated@example.com").
					Return(nil, errors.New("name is required"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid request body",
			},
		},
		{
			name: "empty email",
			id:   "1",
			requestBody: UpdateExampleRequest{
				Name:  "John Doe Updated",
				Email: "",
			},
			mockSetup: func(service *MockExampleService) {
				existing := &models.ExampleModel{
					ID:    1,
					Name:  "John Doe",
					Email: "john@example.com",
				}
				service.On("GetExampleByID", mock.Anything, int64(1)).
					Return(existing, nil)
				service.On("UpdateExample", mock.Anything, int64(1), "John Doe Updated", "").
					Return(nil, errors.New("email is required"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid request body",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockExampleService)
			tt.mockSetup(mockService)
			
			handler := NewExampleHandler(mockService)
			router := setupTestRouter(handler)
			
			reqBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPut, "/api/v1/examples/"+tt.id, bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			
			assert.Equal(t, tt.expectedStatus, resp.Code)
			assert.Equal(t, tt.expectedBody, resp.Body.Body)
			
			mockService.AssertExpectations(t)
		})
	}
}

func TestExampleHandler_DeleteExample(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockSetup      func(*MockExampleService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful deletion",
			id:   "1",
			mockSetup: func(service *MockExampleService) {
				service.On("DeleteExample", mock.Anything, int64(1)).
					Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Example deleted successfully",
			},
		},
		{
			name: "invalid ID",
			id:   "invalid",
			mockSetup: func(service *MockExampleService) {
				// No mock setup needed for invalid ID as it will return 400 directly
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid ID parameter",
			},
		},
		{
			name: "not found",
			id:   "999",
			mockSetup: func(service *MockExampleService) {
				service.On("DeleteExample", mock.Anything, int64(999)).
					Return(errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "Example not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockExampleService)
			tt.mockSetup(mockService)
			
			handler := NewExampleHandler(mockService)
			router := setupTestRouter(handler)
			
			req, _ := http.NewRequest(http.MethodDelete, "/api/v1/examples/"+tt.id, nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)
			
			assert.Equal(t, tt.expectedStatus, resp.Code)
			assert.Equal(t, tt.expectedBody, resp.Body.Body)
			
			mockService.AssertExpectations(t)
		})
	}
}
```

```
