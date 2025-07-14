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
