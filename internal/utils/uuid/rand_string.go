package uuid

import (
	"crypto/rand"
	"errors"
	"fmt"
)

func GenerateRandomStringId(n uint) (string, error) {
	if n > 128 {
		return "", errors.New("max than 128")
	}

	randomBytes := make([]byte, n) // 1 byte = 8 bits = 2 characters
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	res := fmt.Sprintf("%x", randomBytes[0:n])[0:n]
	return res, nil
}
