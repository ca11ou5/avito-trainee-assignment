package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ca11ou5/avito-trainee-assignment/internal/entity"
	"github.com/ca11ou5/avito-trainee-assignment/internal/payload"
)

var ErrWrongPassword = errors.New("wrong password")

func (s *Service) AuthenticateUser(ctx context.Context, req payload.AuthRequest) (string, error) {
	exists, err := s.repo.IsEmployeeExists(ctx, req.Username)
	if err != nil {
		return "", fmt.Errorf("is employee exists: %s", err)
	}

	incomingPassword := req.Password

	req.Password = s.hashPassword(req.Password)
	if !exists {
		err = s.repo.InsertEmployee(ctx, req)
		if err != nil {
			return "", fmt.Errorf("create employee: %s", err)
		}

		token, err := s.createJWT(req.Username)
		if err != nil {
			return "", fmt.Errorf("create jwt: %s", err)
		}

		return token, nil
	}

	hashedPassword, err := s.repo.GetHashedPassword(ctx, req.Username)
	if err != nil {
		return "", fmt.Errorf("get hashed password: %s", err)
	}

	err = s.comparePasswords(hashedPassword, incomingPassword)
	if err != nil {
		return "", ErrWrongPassword
	}

	token, err := s.createJWT(req.Username)
	if err != nil {
		return "", fmt.Errorf("create jwt: %s", err)
	}

	return token, nil
}

func (s *Service) ExtractUserInfo(ctx context.Context, token string) (entity.EmployeeInfo, error) {
	username, err := s.verifyJWT(token)
	if err != nil {
		return entity.EmployeeInfo{}, fmt.Errorf("verify jwt: %s", err)
	}

	s.gatherEmployeeInfo(ctx, username)

}

func (s *Service) gatherEmployeeInfo(ctx context.Context, employeeUsername string) (entity.EmployeeInfo, error) {
	s.repo.
}
