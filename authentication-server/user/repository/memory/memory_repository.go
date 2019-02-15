package memory

// Repository implements with in-memory for testing purpose
type Repository struct {
}

// NewRepository return new object
func NewRepository() *Repository {
	return &Repository{}
}
