package http

import (
	"errors"
	"net/http"
	"strings"
)

const authHeader = "Authorization"

var (
	errInvalidToken = errors.New("invalid authorization token format")
	errMissingToken = errors.New("missing bearer token")
)

func extractBearerToken(h *http.Header) (string, error) {
	bearerToken := h.Get(authHeader)

	tokenParts := strings.Split(bearerToken, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", errInvalidToken
	}

	if tokenParts[1] == "" {
		return "", errMissingToken
	}

	return tokenParts[1], nil
}
