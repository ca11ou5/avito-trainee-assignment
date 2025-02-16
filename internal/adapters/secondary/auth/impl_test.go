package auth

import (
	"github.com/ca11ou5/avito-trainee-assignment/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func newTestAdapter() *Adapter {
	return &Adapter{
		JWTSalt: "testsecret",
	}
}

// TestCreateAndVerifyAuthToken проверяет, что созданный токен можно верифицировать и извлечь из него username.
func TestCreateAndVerifyAuthToken(t *testing.T) {
	adpt := newTestAdapter()
	username := "testuser"

	// Создаём токен
	token, err := adpt.CreateAuthToken(username)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Верифицируем токен
	extractedUsername, err := adpt.VerifyAuthToken(token)
	assert.NoError(t, err)
	assert.Equal(t, username, extractedUsername)
}

func TestVerifyAuthToken_InvalidToken(t *testing.T) {
	adpt := newTestAdapter()

	_, err := adpt.VerifyAuthToken("invalid.token.value")
	assert.Error(t, err)
}

// TestVerifyAuthToken_WrongSigningMethod проверяет, что при изменении алгоритма подписи в токене верификация падает.
func TestVerifyAuthToken_WrongSigningMethod(t *testing.T) {
	adpt := newTestAdapter()
	username := "testuser"

	token, err := adpt.CreateAuthToken(username)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	adpt.JWTSalt = "wrongsecret"
	_, err = adpt.VerifyAuthToken(token)
	assert.Error(t, err)
}

func TestHashAndComparePassword(t *testing.T) {
	adpt := newTestAdapter()
	password := "mysecretpassword"

	hashed, err := adpt.HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hashed)
	assert.NotEqual(t, password, hashed)

	err = adpt.ComparePasswords(hashed, password)
	assert.NoError(t, err)

	err = adpt.ComparePasswords(hashed, "wrongpassword")
	assert.Error(t, err)
}

// TestVerifyAuthToken_ExpiredToken проверяет сценарий истечения срока действия токена.
func TestVerifyAuthToken_ExpiredToken(t *testing.T) {
	adpt := newTestAdapter()
	username := "testuser"

	expTime := time.Now().Add(-1 * time.Minute) // токен уже истёк
	claims := &models.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Minute)),
		},
	}
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := tokenObj.SignedString([]byte(adpt.JWTSalt))
	assert.NoError(t, err)

	_, err = adpt.VerifyAuthToken(tokenStr)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is expired")
}
