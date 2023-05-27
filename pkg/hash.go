package pkg

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func MustHashPassword(password string) string {
	hash, err := HashPassword(password)
	if err != nil {
		log.Fatalf("error while hashing password, %v", err)
	}
	return hash
}

func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
