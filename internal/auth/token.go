package auth

import (
	"net/http"
	"strings"
	"errors"
)

func GetBearerToken(headers http.Header) (string, error) {
	token := headers.Get("Authorization")

	if token == "" {
		return "", errors.New("missing Auth token")
	}

	if !strings.HasPrefix(token, "Bearer ") {
		return "", errors.New("missing Bearer prefix")
	}

	bearerToken := strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))

	if bearerToken == "" {
		return "", errors.New("missing token")
	}

	return bearerToken, nil
}