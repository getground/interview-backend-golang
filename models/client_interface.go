package models

import (
	"context"
)

type ClientInterface interface {
	Close() error
}

type ExampleModel struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ExampleRepository interface {
	Create(ctx context.Context, example *ExampleModel) error
	GetByID(ctx context.Context, id int64) (*ExampleModel, error)
	GetAll(ctx context.Context) ([]*ExampleModel, error)
	Update(ctx context.Context, example *ExampleModel) error
	Delete(ctx context.Context, id int64) error
}
