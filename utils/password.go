package utils

import (
	"github.com/alexedwards/argon2id"
)

var params = &argon2id.Params{
	Memory:      128 * 1024, 
	Iterations:  3,
	Parallelism: 4,
	SaltLength:  16,
	KeyLength:   32,
}

// HashPassword creates a password hash using Argon2
func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, params)
	if err != nil {
		return "", err
	}
	return hash, nil
}

// ComparePasswords compares a password with a hash
func ComparePasswords(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
