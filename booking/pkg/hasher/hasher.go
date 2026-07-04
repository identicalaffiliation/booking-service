package hasher

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Hasher struct{}

func NewHasher() *Hasher { return &Hasher{} }

func (h *Hasher) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("generate hash: %w", err)
	}

	return string(hash), nil
}

func (h *Hasher) ComparePassword(password, reqPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(reqPassword))
}

func (h *Hasher) Hash(val string) string {
	sum := sha256.Sum256([]byte(val))
	return hex.EncodeToString(sum[:])
}

func (h *Hasher) CompareHash(hash, val string) bool {
	expected := h.Hash(val)
	return subtle.ConstantTimeCompare(
		[]byte(hash),
		[]byte(expected),
	) == 1
}
