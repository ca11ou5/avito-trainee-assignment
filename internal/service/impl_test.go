package service

import (
	"context"
	"github.com/ca11ou5/avito-trainee-assignment/internal/adapters/secondary/postgres"
	"github.com/ca11ou5/avito-trainee-assignment/internal/mock"
	"github.com/ca11ou5/avito-trainee-assignment/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func prepareTestDependency() (context.Context, *mock.MerchRepository, *mock.AuthRepository, *Service) {
	ctx := context.Background()

	merchRepo := new(mock.MerchRepository)
	authRepo := new(mock.AuthRepository)

	svc := &Service{
		merch: merchRepo,
		auth:  authRepo,
	}

	return ctx, merchRepo, authRepo, svc
}

func TestAuthenticateUser_NewUser(t *testing.T) {
	ctx, merchRepo, authRepo, svc := prepareTestDependency()
	creds := models.Credentials{
		Username: "mocktesting",
		Password: "mocktesting",
	}

	// setup dependency
	merchRepo.On("IsEmployeeExists", ctx, creds.Username).Return(postgres.ErrEmployeeNotExists)

	hashedPassword := "examplehashedpassword"
	authRepo.On("HashPassword", creds.Password).Return(hashedPassword, nil)

	merchRepo.On("InsertEmployee", ctx, models.Credentials{
		Username: creds.Username,
		Password: hashedPassword,
	}).Return(nil)

	expectedToken := "exampletoken"
	authRepo.On("CreateAuthToken", creds.Username).Return(expectedToken, nil)

	actualToken, err := svc.AuthenticateUser(ctx, creds)

	// assertion
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, actualToken)

	merchRepo.AssertExpectations(t)
	authRepo.AssertExpectations(t)
}

func TestAuthenticateUser_IncorrectCredentialsToExistingUser(t *testing.T) {
	ctx, merchRepo, authRepo, svc := prepareTestDependency()
	creds := models.Credentials{
		Username: "mocktesting",
		Password: "mocktesting",
	}

	// setup dependency
	merchRepo.On("IsEmployeeExists", ctx, creds.Username).Return(nil)

	hashedPassword := "examplehashedpassword"
	merchRepo.On("GetHashedPassword", ctx, creds.Username).Return(hashedPassword, nil)

	authRepo.On("ComparePasswords", hashedPassword, creds.Password).Return(ErrWrongPassword)

	token, err := svc.AuthenticateUser(ctx, creds)

	// assertion
	assert.ErrorIs(t, err, ErrWrongPassword)
	assert.Equal(t, "", token)

	merchRepo.AssertExpectations(t)
	authRepo.AssertExpectations(t)
}

func TestAuthenticateUser_СorrectCredentialsToExistingUser(t *testing.T) {
	ctx, merchRepo, authRepo, svc := prepareTestDependency()
	creds := models.Credentials{
		Username: "mocktesting",
		Password: "mocktesting",
	}

	// setup dependency
	merchRepo.On("IsEmployeeExists", ctx, creds.Username).Return(nil)

	hashedPassword := "examplehashedpassword"
	merchRepo.On("GetHashedPassword", ctx, creds.Username).Return(hashedPassword, nil)

	authRepo.On("ComparePasswords", hashedPassword, creds.Password).Return(nil)

	expectedToken := "exampletoken"
	authRepo.On("CreateAuthToken", creds.Username).Return(expectedToken, nil)

	actualToken, err := svc.AuthenticateUser(ctx, creds)

	// assertion
	assert.NoError(t, err, nil)
	assert.Equal(t, expectedToken, actualToken)

	merchRepo.AssertExpectations(t)
	authRepo.AssertExpectations(t)
}

func TestExtractUserInfo_EmployeeNotExists(t *testing.T) {
	ctx, merchRepo, authRepo, svc := prepareTestDependency()
	token := "exampletoken"

	// setup dependency
	username := "exampleusername"
	authRepo.On("VerifyAuthToken", token).Return(username, nil)

	merchRepo.On("IsEmployeeExists", ctx, username).Return(postgres.ErrEmployeeNotExists)

	actualInfo, err := svc.ExtractUserInfo(ctx, token)

	// assertion
	assert.ErrorIs(t, err, postgres.ErrEmployeeNotExists)
	assert.Equal(t, models.EmployeeInfo{}, actualInfo)

	merchRepo.AssertExpectations(t)
	authRepo.AssertExpectations(t)
}

func TestExtractUserInfo_Successful(t *testing.T) {
	ctx, merchRepo, authRepo, svc := prepareTestDependency()
	token := "exampletoken"

	// setup dependencies
	username := "exampleusername"
	authRepo.On("VerifyAuthToken", token).Return(username, nil)

	merchRepo.On("IsEmployeeExists", ctx, username).Return(nil)

	expected := models.EmployeeInfo{
		Coins:     1000,
		Inventory: []models.Item{},
		CoinHistory: models.CoinHistory{
			Received: []models.ReceivedTransaction{},
			Sent:     []models.SentTransaction{},
		},
	}
	merchRepo.On("GetEmployeeInfo", ctx, username).Return(expected, nil)

	actualInfo, err := svc.ExtractUserInfo(ctx, token)

	// assertion
	assert.NoError(t, err)
	assert.Equal(t, expected, actualInfo)

	merchRepo.AssertExpectations(t)
	authRepo.AssertExpectations(t)
}

func TestSendCoin_Successful(t *testing.T) {
	ctx, merchRepo, authRepo, svc := prepareTestDependency()
	token := "exampletoken"
	trans := models.SentTransaction{
		ToUser: "examplereceiver",
		Amount: 500,
	}

	// setup dependencies
	username := "exampleusername"
	authRepo.On("VerifyAuthToken", token).Return(username, nil)

	merchRepo.On("IsEmployeeExists", ctx, username).Return(nil)
	merchRepo.On("IsEmployeeExists", ctx, trans.ToUser).Return(nil)

	merchRepo.On("SendCoin", ctx, username, trans).Return(nil)

	err := svc.SendCoin(ctx, token, trans)

	// assertion
	assert.NoError(t, err)
	assert.Equal(t, nil, err)

	merchRepo.AssertExpectations(t)
	authRepo.AssertExpectations(t)
}

func TestSendCoin_InvalidToken(t *testing.T) {
	ctx, merchRepo, authRepo, svc := prepareTestDependency()
	token := "exampletoken"
	trans := models.SentTransaction{
		ToUser: "examplereceiver",
		Amount: 500,
	}

	// setup dependencies
	authRepo.On("VerifyAuthToken", token).Return("", ErrInvalidToken)

	err := svc.SendCoin(ctx, token, trans)

	// assertion
	assert.ErrorIs(t, err, ErrInvalidToken)

	merchRepo.AssertExpectations(t)
	authRepo.AssertExpectations(t)
}

func TestBuyItem_Successful(t *testing.T) {
	ctx, merchRepo, authRepo, svc := prepareTestDependency()
	token := "exampletoken"
	item := "pink-hoody"

	// setup dependencies
	username := "exampleusername"
	authRepo.On("VerifyAuthToken", token).Return(username, nil)

	merchRepo.On("IsEmployeeExists", ctx, username).Return(nil)

	merchRepo.On("IsMerchExists", ctx, item).Return(nil)

	merchRepo.On("InsertEmployeeMerch", ctx, username, item).Return(nil)

	err := svc.BuyItem(ctx, token, item)

	// assertion
	assert.NoError(t, err)

	merchRepo.AssertExpectations(t)
	authRepo.AssertExpectations(t)
}

func TestBuyItem_NotEnoughBalance(t *testing.T) {
	ctx, merchRepo, authRepo, svc := prepareTestDependency()
	token := "exampletoken"
	item := "pink-hoody"

	// setup dependencies
	username := "exampleusername"
	authRepo.On("VerifyAuthToken", token).Return(username, nil)

	merchRepo.On("IsEmployeeExists", ctx, username).Return(nil)

	merchRepo.On("IsMerchExists", ctx, item).Return(nil)

	merchRepo.On("InsertEmployeeMerch", ctx, username, item).Return(postgres.ErrNotEnoughBalance)

	err := svc.BuyItem(ctx, token, item)

	// assertion
	assert.ErrorIs(t, postgres.ErrNotEnoughBalance, err)

	merchRepo.AssertExpectations(t)
	authRepo.AssertExpectations(t)
}
