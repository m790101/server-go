package auth

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	data, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// CheckPasswordHash -
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GenerateRefreshToken() (string, error) {
	c := 10
	token := make([]byte, c)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	encodedStr := hex.EncodeToString(token)

	return encodedStr, nil
}
