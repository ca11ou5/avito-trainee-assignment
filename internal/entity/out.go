package entity

// INFO RESPONSE

type infoResponse struct {
	Coins       int         `json:"coins"`
	Inventory   []Item      `json:"inventory"`
	CoinHistory CoinHistory `json:"coinHistory"`
}

type Item struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type CoinHistory struct {
	Received []FromTransaction `json:"received"`
	Sent     []ToTransaction   `json:"sent"`
}

type FromTransaction struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type ToTransaction struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

// INFO RESPONSE

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
