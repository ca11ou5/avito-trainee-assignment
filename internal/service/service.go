package service

import (
	"context"
	"github.com/ca11ou5/avito-trainee-assignment/internal/models"
)

type Repository interface {
	// POST /api/auth
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

type Service struct {
	repo Repository

	JWTSalt string
}

func New(repo Repository, jwtSalt string) *Service {
	return &Service{
		repo:    repo,
		JWTSalt: jwtSalt,
	}
}
