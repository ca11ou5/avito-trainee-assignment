package mock

import "github.com/stretchr/testify/mock"

type AuthRepository struct {
	mock.Mock
}

func (a *AuthRepository) HashPassword(password string) (hashedPassword string, err error) {
	args := a.Called(password)
	return args.Get(0).(string), args.Error(1)
}

func (a *AuthRepository) ComparePasswords(hashedPassword string, password string) error {
	args := a.Called(hashedPassword, password)
	return args.Error(0)
}

func (a *AuthRepository) CreateAuthToken(username string) (token string, err error) {
	args := a.Called(username)
	return args.Get(0).(string), args.Error(1)
}

func (a *AuthRepository) VerifyAuthToken(token string) (username string, err error) {
	args := a.Called(token)
	return args.Get(0).(string), args.Error(1)
}
