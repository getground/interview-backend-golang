package example

import (
	"context"

	"github.com/getground/interview-backend-golang/models"
	"github.com/pkg/errors"
)

type Service interface {
	CreateExample(ctx context.Context, name, email string) (*models.ExampleModel, error)
	GetExampleByID(ctx context.Context, id int64) (*models.ExampleModel, error)
	GetAllExamples(ctx context.Context) ([]*models.ExampleModel, error)
	UpdateExample(ctx context.Context, id int64, name, email string) (*models.ExampleModel, error)
	DeleteExample(ctx context.Context, id int64) error
}

type service struct {
	repo models.ExampleRepository
}

func NewService(repo models.ExampleRepository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateExample(ctx context.Context, name, email string) (*models.ExampleModel, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	example := &models.ExampleModel{
		Name:  name,
		Email: email,
	}
	err := s.repo.Create(ctx, example)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create example")
	}
	return example, nil
}

func (s *service) GetExampleByID(ctx context.Context, id int64) (*models.ExampleModel, error) {
	example, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get example with id: %d", id)
	}
	return example, nil
}

func (s *service) GetAllExamples(ctx context.Context) ([]*models.ExampleModel, error) {
	examples, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all examples")
	}
	return examples, nil
}

func (s *service) UpdateExample(ctx context.Context, id int64, name, email string) (*models.ExampleModel, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	example := &models.ExampleModel{
		ID:    id,
		Name:  name,
		Email: email,
	}
	err := s.repo.Update(ctx, example)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to update example with id: %d", id)
	}
	return example, nil
}

func (s *service) DeleteExample(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "failed to delete example with id: %d", id)
	}
	return nil
}
