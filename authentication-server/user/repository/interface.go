package repository

import (
	"time"
)

// Model is exported user model
type Model struct {
	ID       int
	Username string
	Created  time.Time
	Updated  time.Time
	IsActive bool
}

// Repository with basic operations to user domain
type Repository interface {
	GetByID(id uint64) (*Model, error)
	GetByUsername(username string) (*Model, error)
	Create(m Model) error
	SetIsActive(id uint64, isActive bool) error
}
