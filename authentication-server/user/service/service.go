package service

import (
	"github.com/tanopwan/oauth-farm/authentication-server/user/repository"
)

// Service contains business logic to the domain
type Service struct {
	usersByIDs      map[int]*Model
	usersByUsername map[string]*Model
	repo            repository.Repository
}

// NewService return new object
func NewService(repo repository.Repository) *Service {
	user := Model{
		ID:       1,
		Username: "tanopwan@gmail.com",
	}
	user2 := Model{
		ID:       2,
		Username: "undermatthew@gmail.com",
	}
	return &Service{
		usersByIDs:      map[int]*Model{1: &user, 2: &user2},
		usersByUsername: map[string]*Model{"tanopwan@gmail.com": &user, "undermatthew@gmail.com": &user2},
		repo:            repo,
	}
}

// GetActiveUserByEmail function, find user with active status
func (s *Service) GetActiveUserByEmail(email string) (*Model, error) {
	return s.usersByUsername[email], nil
}

// GetUserByID function, find user
func (s *Service) GetUserByID(userID int) (*Model, error) {
	return s.usersByIDs[userID], nil
}
