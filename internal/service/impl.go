package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ca11ou5/avito-trainee-assignment/internal/adapters/secondary/postgres"
	"github.com/ca11ou5/avito-trainee-assignment/internal/models"
)

var (
	ErrWrongPassword      = errors.New("wrong password")
	ErrInvalidToken       = errors.New("invalid token")
	ErrCantSentToYourself = errors.New("can't sent to yourself")
)

func (s *Service) AuthenticateUser(ctx context.Context, creds models.Credentials) (string, error) {
	err := s.repo.IsEmployeeExists(ctx, creds.Username)

	incomingPassword := creds.Password
	creds.Password = s.hashPassword(creds.Password)

	// employee does not exist, create
	if err != nil {
		err = s.repo.InsertEmployee(ctx, creds)
		if err != nil {
			return "", fmt.Errorf("create employee: %s", err)
		}

		token, err := s.createJWT(creds.Username)
		if err != nil {
			return "", fmt.Errorf("create jwt: %s", err)
		}

		return token, nil
	}

	hashedPassword, err := s.repo.GetHashedPassword(ctx, creds.Username)
	if err != nil {
		return "", fmt.Errorf("get hashed password: %s", err)
	}

	err = s.comparePasswords(hashedPassword, incomingPassword)
	if err != nil {
		return "", ErrWrongPassword
	}

	token, err := s.createJWT(creds.Username)
	if err != nil {
		return "", fmt.Errorf("create jwt: %s", err)
	}

	return token, nil
}

func (s *Service) ExtractUserInfo(ctx context.Context, token string) (models.EmployeeInfo, error) {
	username, err := s.verifyJWT(token)
	if err != nil {
		return models.EmployeeInfo{}, fmt.Errorf("%w: %s", ErrInvalidToken, err)
	}

	err = s.repo.IsEmployeeExists(ctx, username)
	if err != nil {
		return models.EmployeeInfo{}, err
	}

	info, err := s.repo.GetEmployeeInfo(ctx, username)
	if err != nil {
		return info, fmt.Errorf("get employee info: %s", err)
	}

	return info, nil
}

func (s *Service) SendCoin(ctx context.Context, token string, trans models.SentTransaction) error {
	username, err := s.verifyJWT(token)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidToken, err)
	}

	if username == trans.ToUser {
		return ErrCantSentToYourself
	}

	for _, name := range []string{username, trans.ToUser} {
		err = s.repo.IsEmployeeExists(ctx, name)
		if err != nil {
			return err
		}
	}

	err = s.repo.SendCoin(ctx, username, trans)
	if err != nil {
		if errors.Is(err, postgres.ErrNotEnoughBalance) {
			return err
		}
		return fmt.Errorf("send coin: %w", err)
	}

	return nil
}

func (s *Service) BuyItem(ctx context.Context, token string, item string) error {
	username, err := s.verifyJWT(token)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidToken, err)
	}

	err = s.repo.IsEmployeeExists(ctx, username)
	if err != nil {
		return err
	}

	err = s.repo.IsMerchExists(ctx, item)
	if err != nil {
		return err
	}

	err = s.repo.InsertEmployeeMerch(ctx, username, item)
	if err != nil {
		if errors.Is(err, postgres.ErrNotEnoughBalance) {
			return err
		}
		return fmt.Errorf("insert employee merch: %s", err)
	}

	return nil
}
