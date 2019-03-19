package common

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
)

// GenerateRandomString ... return base 64 string
func GenerateRandomString(n int) string {
	b := make([]byte, n)
	io.ReadFull(rand.Reader, b)
	return base64.RawURLEncoding.EncodeToString(b)
}

// HashSHA256 ... return hash string base64 of input string
func HashSHA256(value string) string {
	h := sha256.New()
	h.Write([]byte(value))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
