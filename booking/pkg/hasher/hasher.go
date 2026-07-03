package hasher

import (
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

func (h *Hasher) Compare(password, reqPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(reqPassword))
}
