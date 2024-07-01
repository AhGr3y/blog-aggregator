package auth

import (
	"errors"
	"net/http"
	"strings"
)

// ExtractApiKeyFromRequest -
func ExtractApiKeyFromRequest(r *http.Request) (string, error) {

	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("missing or invalid authorization header")
	}

	parts := strings.Split(header, " ")
	if len(parts) < 2 || parts[0] != "ApiKey" {
		return "", errors.New("invalid authorization header format: " + header)
	}

	apiKey := parts[1]

	return apiKey, nil
}
