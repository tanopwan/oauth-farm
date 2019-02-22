package postgres

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/tanopwan/oauth-farm/authentication-server/user/repository"
	"log"
)

// Repository implements with in-memory for testing purpose
type Repository struct {
	dbx *sqlx.DB
}

// NewRepository return new object
func NewRepository(db *sql.DB) *Repository {
	dbx := sqlx.NewDb(db, "postgres")
	return &Repository{
		dbx: dbx,
	}
}

// GetByID get user by ID
func (r *Repository) GetByID(id int) (*repository.Model, error) {
	q := `select * from app_user where id = :id`
	e := user{}
	stmt, err := r.dbx.PrepareNamed(q)
	if err != nil {
		log.Printf("[GetByID] failed to prepare named with reason: %s\n", err.Error())
		return nil, err
	}

	err = stmt.Get(&e, id)
	if err != nil {
		log.Printf("[GetByID] failed to get with reason: %s\n", err.Error())
		return nil, err
	}

	return e.toModel(), nil
}

// GetByUsername get user by username
func (r *Repository) GetByUsername(username string) (*repository.Model, error) {
	return nil, nil
}

// GetByEmail get user by email
func (r *Repository) GetByEmail(email string) (*repository.Model, error) {
	return nil, nil
}

// Create create user
func (r *Repository) Create(m repository.Model) error {
	return nil
}

// SetIsActive set user active status
func (r *Repository) SetIsActive(id int, isActive bool) error {
	return nil
}
