package auth

import (
	"net/http"
	"errors"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	apikey := headers.Get("Authorization")

	if apikey == "" {
		return "", errors.New("missing Auth token")
	}

	if !strings.HasPrefix(apikey, "ApiKey ") {
		return "", errors.New("missing ApiKey prefix")
	}

	bearerKey := strings.TrimSpace(strings.TrimPrefix(apikey, "ApiKey "))

	if bearerKey == "" {
		return "", errors.New("missing token")
	}

	return bearerKey, nil
}