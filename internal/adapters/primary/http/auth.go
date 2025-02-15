package http

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

const (
	authHeader      = "Authorization"
	contextTokenKey = "token"
)

var (
	errInvalidToken = errors.New("invalid authorization token format")
	errMissingToken = errors.New("missing bearer token")
)

func extractBearerToken(h http.Header) (string, error) {
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

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		token, err := extractBearerToken(r.Header)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(beatifyError(err))
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextTokenKey, token)))
	})
}
