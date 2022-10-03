package avp

import (
	"crypto/sha256"
	"fmt"
)

const (
	SALT = "AVPFOREVER"
)

// HashPassword produces a sha256 hash from a given password after salting it.
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(SALT + password))
	return fmt.Sprintf("%x", hash)
}

// CheckPasswordHash checks a challenge password against a hash.
func CheckPasswordHash(password, hash string) bool {
	return HashPassword(password) == hash
}
