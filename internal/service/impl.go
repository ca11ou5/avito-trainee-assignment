package service

import (
	"github.com/ca11ou5/avito-trainee-assignment/internal/adapters/secondary/postgres"
	"github.com/ca11ou5/avito-trainee-assignment/internal/models"

	"context"
	"errors"
	"fmt"
)

var (
	ErrWrongPassword      = errors.New("wrong password")
	ErrInvalidToken       = errors.New("invalid token")
	ErrCantSentToYourself = errors.New("can't sent to yourself")
)

func (s *Service) AuthenticateUser(ctx context.Context, creds models.Credentials) (string, error) {
	err := s.merch.IsEmployeeExists(ctx, creds.Username)
	// employee does not exist, create
	if err != nil {
		creds.Password, err = s.auth.HashPassword(creds.Password)
		if err != nil {
			return "", fmt.Errorf("hash password: %s", err)
		}

		err = s.merch.InsertEmployee(ctx, creds)
		if err != nil {
			return "", fmt.Errorf("create employee: %s", err)
		}

		token, err := s.auth.CreateAuthToken(creds.Username)
		if err != nil {
			return "", fmt.Errorf("create jwt: %s", err)
		}

		return token, nil
	}

	hashedPassword, err := s.merch.GetHashedPassword(ctx, creds.Username)
	if err != nil {
		return "", fmt.Errorf("get hashed password: %s", err)
	}

	err = s.auth.ComparePasswords(hashedPassword, creds.Password)
	if err != nil {
		return "", ErrWrongPassword
	}

	token, err := s.auth.CreateAuthToken(creds.Username)
	if err != nil {
		return "", fmt.Errorf("create jwt: %s", err)
	}

	return token, nil
}

func (s *Service) ExtractUserInfo(ctx context.Context, token string) (models.EmployeeInfo, error) {
	username, err := s.auth.VerifyAuthToken(token)
	if err != nil {
		return models.EmployeeInfo{}, fmt.Errorf("%w: %s", ErrInvalidToken, err)
	}

	err = s.merch.IsEmployeeExists(ctx, username)
	if err != nil {
		return models.EmployeeInfo{}, err
	}

	info, err := s.merch.GetEmployeeInfo(ctx, username)
	if err != nil {
		return info, fmt.Errorf("get employee info: %s", err)
	}

	return info, nil
}

func (s *Service) SendCoin(ctx context.Context, token string, trans models.SentTransaction) error {
	username, err := s.auth.VerifyAuthToken(token)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidToken, err)
	}

	if username == trans.ToUser {
		return ErrCantSentToYourself
	}

	for _, name := range []string{username, trans.ToUser} {
		err = s.merch.IsEmployeeExists(ctx, name)
		if err != nil {
			return err
		}
	}

	err = s.merch.SendCoin(ctx, username, trans)
	if err != nil {
		if errors.Is(err, postgres.ErrNotEnoughBalance) {
			return err
		}
		return fmt.Errorf("send coin: %w", err)
	}

	return nil
}

func (s *Service) BuyItem(ctx context.Context, token string, item string) error {
	username, err := s.auth.VerifyAuthToken(token)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidToken, err)
	}

	err = s.merch.IsEmployeeExists(ctx, username)
	if err != nil {
		return err
	}

	err = s.merch.IsMerchExists(ctx, item)
	if err != nil {
		return err
	}

	err = s.merch.InsertEmployeeMerch(ctx, username, item)
	if err != nil {
		if errors.Is(err, postgres.ErrNotEnoughBalance) {
			return err
		}
		return fmt.Errorf("insert employee merch: %s", err)
	}

	return nil
}
