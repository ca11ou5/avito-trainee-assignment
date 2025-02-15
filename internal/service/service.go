package service

import (
	"context"
	"github.com/ca11ou5/avito-trainee-assignment/internal/payload"
)

type Repository interface {
	// TODO: decompose
	IsEmployeeExists(ctx context.Context, username string) (exists bool, err error)
	InsertEmployee(ctx context.Context, req payload.AuthRequest) (err error)
	GetHashedPassword(ctx context.Context, username string) (hashedPassword string, err error)
	GetCoinsTransactions(ctx context.Context, username string)
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
