// Package auth ...
// TODO
package auth

import (
	"errors"
	"os"
)

var (
	// ErrMissingAPIKey is the error message for missing API key.
	ErrMissingAPIKey = errors.New("missing API key in the request header or parameters")
)

// LookupAndSetAuth sets the env variable API_KEY in the supplied request.
func LookupAndSetAuth() (string, error) {
	k, ok := os.LookupEnv("API_KEY")
	if !ok {
		return "", ErrMissingAPIKey
	}

	return k, nil
}
