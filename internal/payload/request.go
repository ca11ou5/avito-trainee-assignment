package payload

import "errors"

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	errRequiredUsername = errors.New("username is required")
	errRequiredPassword = errors.New("password is required")
	errPasswordTooLong  = errors.New("password is too long")
	errPasswordTooShort = errors.New("password is too short")
)

func (r *AuthRequest) Validate() error {
	if r.Username == "" {
		return errRequiredUsername
	}

	if r.Password == "" {
		return errRequiredPassword
	}

	if len(r.Password) > 72 {
		return errPasswordTooLong
	}

	if len(r.Password) < 8 {
		return errPasswordTooShort
	}

	return nil
}
