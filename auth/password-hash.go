package auth

import (
	"crypto/sha256"
	"encoding/hex"
)

//HashPassword hashes the password
func HashPassword(password string) string {
	hashedPassword := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hashedPassword[:])
}
