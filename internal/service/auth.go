package service

import (
	"fmt"
	"github.com/ca11ou5/avito-trainee-assignment/internal/entity"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (s *Service) createJWT(username string) (string, error) {
	expTime := time.Now().Add(24 * time.Hour)

	claims := &entity.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(s.JWTSalt))
	if err != nil {
		return "", fmt.Errorf("sign token: %s", err)
	}

	return signedToken, nil
}

func (s *Service) verifyJWT(token string) (string, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &entity.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}

		return s.JWTSalt, nil
	})
	if err != nil {
		return "", fmt.Errorf("parse jwt: %s", err)
	}

	if claims, ok := jwtToken.Claims.(*entity.Claims); ok && jwtToken.Valid {
		username := claims.Username

		if username == "" {
			return "", fmt.Errorf("username is empty")
		}

		exp := claims.ExpiresAt
		if exp == nil {
			return "", fmt.Errorf("expiration time is empty")
		}

		if time.Now().After(exp.Time) {
			return "", fmt.Errorf("token expired")
		}

		return username, nil
	}

	return "", fmt.Errorf("invalid jwt token")
}

func (s *Service) hashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}

func (s *Service) comparePasswords(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
