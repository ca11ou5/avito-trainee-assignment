package mock

import (
	"github.com/ca11ou5/avito-trainee-assignment/internal/models"
	"github.com/stretchr/testify/mock"

	"context"
)

type MerchRepository struct {
	mock.Mock
}

func (m *MerchRepository) IsEmployeeExists(ctx context.Context, username string) error {
	args := m.Called(ctx, username)
	return args.Error(0)
}

func (m *MerchRepository) InsertEmployee(ctx context.Context, creds models.Credentials) error {
	args := m.Called(ctx, creds)
	return args.Error(0)
}

func (m *MerchRepository) GetHashedPassword(ctx context.Context, username string) (string, error) {
	args := m.Called(ctx, username)
	return args.String(0), args.Error(1)
}

func (m *MerchRepository) GetEmployeeInfo(ctx context.Context, username string) (models.EmployeeInfo, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(models.EmployeeInfo), args.Error(1)
}

func (m *MerchRepository) SendCoin(ctx context.Context, username string, trans models.SentTransaction) error {
	args := m.Called(ctx, username, trans)
	return args.Error(0)
}

func (m *MerchRepository) IsMerchExists(ctx context.Context, itemName string) error {
	args := m.Called(ctx, itemName)
	return args.Error(0)
}

func (m *MerchRepository) InsertEmployeeMerch(ctx context.Context, username string, merch string) error {
	args := m.Called(ctx, username, merch)
	return args.Error(0)
}
