package service

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"github.com/pkg/errors"
	"io"
)

// Service contains business logic to the domain
type Service struct {
	sessions map[string]int
}

// NewService return new object
func NewService() *Service {
	return &Service{
		sessions: make(map[string]int),
	}
}

// CreateSession function
func (s *Service) CreateSession(userID int) (*Model, error) {
	b := make([]byte, 64)
	_, err := io.ReadFull(crand.Reader, b)
	if err != nil {
		return nil, errors.Wrap(err, "createsession: read crypto random err")
	}

	h := sha256.New()
	_, err = h.Write(b)
	if err != nil {
		return nil, errors.Wrap(err, "createsession: hash sha256 err")
	}

	hash := hex.EncodeToString(h.Sum(nil))
	s.sessions[hash] = userID
	return &Model{
		UserID: userID,
		Hash:   hash,
	}, nil
}

// ValidateSession function
func (s *Service) ValidateSession(hash string) (*Model, error) {
	return &Model{
		UserID: s.sessions[hash],
		Hash:   hash,
	}, nil
}

// RemoveSession function
func (s *Service) RemoveSession(hash string) {
	delete(s.sessions, hash)
}
