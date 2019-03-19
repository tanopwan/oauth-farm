package session

// Model for exported session
// Empty user(value) means non-loggedin session
type Model struct {
	Hash  string
	Value string
}
