package entity

import "github.com/golang-jwt/jwt/v5"

// ERROR RESPONSE

type errorResponse struct {
	Errors string `json:"errors"`
}

// ERROR RESPONSE

// AUTH REQUEST

type AuthRequest struct {
	// required
	Username string `json:"username"`

	// required
	Password string `json:"password"`
}

// AUTH REQUEST

// AUTH RESPONSE

type authResponse struct {
	Token string `json:"token"`
}

// AUTH RESPONSE

// SEND COIN REQUEST

type sendCoinRequest struct {
	// required
	ToUser string `json:"toUser"`

	// required
	Amount int `json:"amount"`
}

// SEND COIN REQUEST

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
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type CoinHistory struct {
	Received []ReceivedTransaction `json:"received"`
	Sent     []SentTransaction     `json:"sent"`
}

type ReceivedTransaction struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type SentTransaction struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}
