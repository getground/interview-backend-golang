package models

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type ExampleRepositoryImpl struct {
	data   map[int64]*ExampleModel
	mu     sync.RWMutex
	nextID int64
}

func NewExampleRepository() ExampleRepository {
	return &ExampleRepositoryImpl{
		data:   make(map[int64]*ExampleModel),
		nextID: 1,
	}
}

func (r *ExampleRepositoryImpl) Create(ctx context.Context, example *ExampleModel) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if example.Name == "" {
		return errors.New("name is required")
	}
	if example.Email == "" {
		return errors.New("email is required")
	}
	for _, existing := range r.data {
		if existing.Email == example.Email {
			return errors.New("email already exists")
		}
	}
	example.ID = r.nextID
	now := time.Now().Format(time.RFC3339)
	example.CreatedAt = now
	example.UpdatedAt = now
	r.data[example.ID] = example
	r.nextID++
	return nil
}

func (r *ExampleRepositoryImpl) GetByID(ctx context.Context, id int64) (*ExampleModel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	example, exists := r.data[id]
	if !exists {
		return nil, errors.Wrapf(errors.New("not found"), "example not found with id: %d", id)
	}
	return &ExampleModel{
		ID:        example.ID,
		Name:      example.Name,
		Email:     example.Email,
		CreatedAt: example.CreatedAt,
		UpdatedAt: example.UpdatedAt,
	}, nil
}

func (r *ExampleRepositoryImpl) GetAll(ctx context.Context) ([]*ExampleModel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	examples := make([]*ExampleModel, 0, len(r.data))
	for _, example := range r.data {
		examples = append(examples, &ExampleModel{
			ID:        example.ID,
			Name:      example.Name,
			Email:     example.Email,
			CreatedAt: example.CreatedAt,
			UpdatedAt: example.UpdatedAt,
		})
	}
	return examples, nil
}

func (r *ExampleRepositoryImpl) Update(ctx context.Context, example *ExampleModel) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if example.Name == "" {
		return errors.New("name is required")
	}
	if example.Email == "" {
		return errors.New("email is required")
	}
	existing, exists := r.data[example.ID]
	if !exists {
		return errors.Wrapf(errors.New("not found"), "example not found with id: %d", example.ID)
	}
	for id, other := range r.data {
		if id != example.ID && other.Email == example.Email {
			return errors.New("email already exists")
		}
	}
	example.CreatedAt = existing.CreatedAt
	example.UpdatedAt = time.Now().Format(time.RFC3339)
	r.data[example.ID] = example
	return nil
}

func (r *ExampleRepositoryImpl) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.data[id]; !exists {
		return errors.Wrapf(errors.New("not found"), "example not found with id: %d", id)
	}
	delete(r.data, id)
	return nil
}
