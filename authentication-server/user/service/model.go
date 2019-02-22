package service

// Status for user status
type Status int

const (
	_ Status = iota
	// Active user can login
	Active
	// Inactive user cannot login
	Inactive
)

func (s Status) String() string {
	return [...]string{"Active", "Inactive"}[s]
}

// Model for exported session
type Model struct {
	ID       int
	Username string
	Status   Status
	IsActive bool
}
