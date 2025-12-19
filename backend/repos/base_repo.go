package repos

import (
	"github.com/google/uuid"
)

type BaseRepository[T any] interface {
	FindByID(id uuid.UUID) (*T, error)
	FindAll() ([]T, error)
	Create(entity *T) error
	Update(entity *T) error
	Delete(id uuid.UUID) error
}
