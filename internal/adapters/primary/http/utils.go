package http

import (
	"encoding/json"
)

func beatifyError(err error) []byte {
	bb, _ := json.Marshal(map[string]string{
		"errors": err.Error(),
	})

	return bb
}

func beatifyToken(token string) []byte {
	bb, _ := json.Marshal(map[string]string{
		"token": token,
	})

	return bb
}
