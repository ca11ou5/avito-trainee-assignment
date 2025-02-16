package service

import (
	"github.com/ca11ou5/avito-trainee-assignment/internal/models"

	"context"
)

type MerchRepository interface {
	// POST /api/auth
	// IsEmployeeExists returns nil if exists, if not ErrEmployeeNotExists
	IsEmployeeExists(ctx context.Context, username string) error
	InsertEmployee(ctx context.Context, creds models.Credentials) (err error)
	GetHashedPassword(ctx context.Context, username string) (hashedPassword string, err error)

	// GET /api/info
	GetEmployeeInfo(ctx context.Context, username string) (models.EmployeeInfo, error)

	// POST /api/sendCoin
	SendCoin(ctx context.Context, username string, trans models.SentTransaction) error

	// GET /api/buy/{item}
	IsMerchExists(ctx context.Context, itemName string) error
	InsertEmployeeMerch(ctx context.Context, username string, merch string) error
}

type AuthRepository interface {
	HashPassword(password string) (hashedPassword string, err error)
	ComparePasswords(hashedPassword string, password string) error
	CreateAuthToken(username string) (token string, err error)
	VerifyAuthToken(token string) (username string, err error)
}

type Service struct {
	merch MerchRepository
	auth  AuthRepository
}

func New(merch MerchRepository, auth AuthRepository) *Service {
	return &Service{
		merch: merch,
		auth:  auth,
	}
}
