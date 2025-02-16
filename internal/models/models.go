package models

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type EmployeeInfo struct {
	Coins       int         `json:"coins"`
	Inventory   []Item      `json:"inventory"`
	CoinHistory CoinHistory `json:"coinHistory"`
}

type Item struct {
	Type     string `json:"type" db:"merch_name"`
	Quantity int    `json:"quantity" db:"count"`
}

type CoinHistory struct {
	Received []ReceivedTransaction `json:"received"`
	Sent     []SentTransaction     `json:"sent"`
}

type ReceivedTransaction struct {
	FromUser string `json:"fromUser" db:"sender_username"`
	Amount   int    `json:"amount" db:"amount"`
}

type SentTransaction struct {
	ToUser string `json:"toUser" db:"receiver_username"`
	Amount int    `json:"amount" db:"amount"`
}

func (st *SentTransaction) Validate() error {
	if st.ToUser == "" {
		return errors.New("'toUser' field is required")
	}

	if st.Amount == 0 {
		return errors.New("'amount' field is required")
	}

	if st.Amount < 0 {
		return errors.New("'amount' field cannot be less than zero")
	}

	return nil
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	errRequiredUsername = errors.New("username is required")
	errRequiredPassword = errors.New("password is required")
	errPasswordTooLong  = errors.New("password is too long")
	errPasswordTooShort = errors.New("password is too short")
)

func (r *Credentials) Validate() error {
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
