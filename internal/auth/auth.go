// Package auth contains functions and helpers for authentication.
package auth

import (
	"errors"
	"os"
)

var (
	// ErrMissingAPIKey is the error message for missing API key.
	ErrMissingAPIKey = errors.New("missing API key in the environment")
)

// LookupAPIAuthKey lookups the env variable API_KEY,
// returns an ErrMissingAPIKey if not found.
func LookupAPIAuthKey() (string, error) {
	k, ok := os.LookupEnv("API_KEY")
	if !ok {
		return "", ErrMissingAPIKey
	}
	return k, nil
}
