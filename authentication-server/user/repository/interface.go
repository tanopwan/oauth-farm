package repository

// Model is exported user model
type Model struct {
	ID       int
	Username string
	IsActive bool
}

// Repository with basic operations to user domain
type Repository interface {
	GetByID(id int) (*Model, error)
	GetByUsername(username string) (*Model, error)
	GetByEmail(email string) (*Model, error)
	Create(m Model) error
	SetIsActive(id int, isActive bool) error
}
