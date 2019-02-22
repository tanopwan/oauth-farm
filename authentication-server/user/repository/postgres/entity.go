package postgres

import (
	"github.com/tanopwan/oauth-farm/authentication-server/user/repository"
	"time"
)

type user struct {
	ID        int       `db:"id"`
	Username  string    `db:"username"`
	Firstname string    `db:"firstname"`
	Lastname  string    `db:"lastname"`
	Email     string    `db:"email"`
	Status    int       `db:"status"`
	Created   time.Time `db:"created"`
	Updated   time.Time `db:"updated"`
}

func (u user) toModel() *repository.Model {
	return &repository.Model{
		ID:       u.ID,
		Username: u.Username,
		IsActive: u.Status == 1,
	}
}
