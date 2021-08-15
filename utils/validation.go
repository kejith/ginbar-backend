package utils

import (
	"net"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// CreatePasswordHash creates a Hash from a String with bcrypt
func CreatePasswordHash(password []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), err
}

// IsEmailValid checks if the email provided passes the required structure
// and length test. It also checks the domain has a valid MX record.
func IsEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}

	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !emailRegex.MatchString(e) {
		return false
	}
	parts := strings.Split(e, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return false
	}
	return true
}
