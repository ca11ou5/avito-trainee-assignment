package http

import (
	"github.com/stretchr/testify/assert"

	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExtractBearerToken_MissingHeader(t *testing.T) {
	headers := http.Header{}
	token, err := extractBearerToken(headers)
	assert.Equal(t, "", token)
	assert.Equal(t, errMissingToken, err)
}

func TestExtractBearerToken_InvalidFormat(t *testing.T) {
	headers := http.Header{}
	headers.Set(authHeader, "InvalidTokenFormat")
	token, err := extractBearerToken(headers)
	assert.Equal(t, "", token)
	assert.Equal(t, errInvalidToken, err)
}

func TestExtractBearerToken_MissingTokenValue(t *testing.T) {
	headers := http.Header{}
	headers.Set(authHeader, "Bearer ")
	token, err := extractBearerToken(headers)
	assert.Equal(t, "", token)
	assert.Equal(t, errMissingTokenValue, err)
}

func TestExtractBearerToken_Valid(t *testing.T) {
	headers := http.Header{}
	expectedToken := "abc123"
	headers.Set(authHeader, "Bearer "+expectedToken)
	token, err := extractBearerToken(headers)
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
}

func TestContextToken(t *testing.T) {
	// Создаем контекст с токеном
	ctx := context.WithValue(context.Background(), tokenKey{}, "mytoken")
	token := contextToken(ctx)
	assert.Equal(t, "mytoken", token)
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	// Создаем next handler, который проверяет, что токен корректно установлен в контексте
	nextCalled := false
	var extractedToken string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		extractedToken = contextToken(r.Context())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	req := httptest.NewRequest("GET", "http://example.com", nil)
	expectedToken := "validtoken"
	req.Header.Set(authHeader, "Bearer "+expectedToken)
	rr := httptest.NewRecorder()

	handler := authMiddleware(next)
	handler.ServeHTTP(rr, req)

	assert.True(t, nextCalled, "Next handler should be called")
	assert.Equal(t, expectedToken, extractedToken)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("next handler should not be called")
	})

	req := httptest.NewRequest("GET", "http://example.com", nil)
	rr := httptest.NewRecorder()

	handler := authMiddleware(next)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	assert.Contains(t, rr.Body.String(), errMissingToken.Error())
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("next handler should not be called")
	})

	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set(authHeader, "InvalidTokenFormat")
	rr := httptest.NewRecorder()

	handler := authMiddleware(next)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), errInvalidToken.Error())
}

func TestAuthMiddleware_EmptyTokenValue(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("next handler should not be called")
	})

	req := httptest.NewRequest("GET", "http://example.com", nil)
	req.Header.Set(authHeader, "Bearer ")
	rr := httptest.NewRecorder()

	handler := authMiddleware(next)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), errMissingTokenValue.Error())
}
