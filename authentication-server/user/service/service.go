package service

// Service contains business logic to the domain
type Service struct {
	usersByIDs      map[int]*Model
	usersByUsername map[string]*Model
}

// NewService return new object
func NewService() *Service {
	user := Model{
		ID:       1,
		Username: "tanopwan@gmail.com",
	}
	return &Service{
		usersByIDs:      map[int]*Model{1: &user},
		usersByUsername: map[string]*Model{"tanopwan@gmail.com": &user},
	}
}

// GetActiveUser function, find user with active status
func (s *Service) GetActiveUser(username string) (*Model, error) {
	return s.usersByUsername[username], nil
}

// GetUserByID function, find user
func (s *Service) GetUserByID(userID int) (*Model, error) {
	return s.usersByIDs[userID], nil
}
