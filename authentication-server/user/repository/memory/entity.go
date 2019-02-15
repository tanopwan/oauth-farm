package memory

import (
	"github.com/tanopwan/oauth-farm/authentication-server/user/repository"
	"time"
)

type user struct {
	ID       int
	Username string
	Created  time.Time
	Updated  time.Time
	IsActive bool
}

func (e user) toModel() *repository.Model {
	return &repository.Model{
		ID:       e.ID,
		Username: e.Username,
		Created:  e.Created,
		Updated:  e.Updated,
		IsActive: e.IsActive,
	}
}

type provider struct {
	ID              int
	UserID          int
	ProviderName    string
	ProviderID      string
	ProviderProfile map[string]interface{}
	Firstname       string
	Lastname        string
	Email           string
	Picture         string
	Created         time.Time
	Updated         time.Time
	IsActive        bool
}
