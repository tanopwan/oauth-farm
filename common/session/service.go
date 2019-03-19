package session

import (
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"github.com/pkg/errors"
	"github.com/tanopwan/oauth-farm/common"
	"io"
)

// Service contains business logic to the domain
type Service struct {
	// sessions  map[string]int
	DataStore common.Cache
}

// NewService return new object
func NewService(cache common.Cache) *Service {
	return &Service{
		DataStore: cache,
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
	s.DataStore.Set(hash, userID)
	return &Model{
		UserID: userID,
		Hash:   hash,
	}, nil
}

// ValidateSession function
func (s *Service) ValidateSession(hash string) (*Model, error) {
	userID, err := s.DataStore.Get(hash)
	if err != nil {
		return nil, errors.Wrap(err, "validatesession: get err")
	}
	return &Model{
		UserID: userID.(int),
		Hash:   hash,
	}, nil
}

// RemoveSession function
func (s *Service) RemoveSession(hash string) {
	s.DataStore.Del(hash)
}
