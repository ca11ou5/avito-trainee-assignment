package e2e

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestE2E_Info(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e tests")
	}

	t.Run("Auth", TestE2E_Auth)
	t.Run("Buy Item", TestE2E_BuyItem)
	t.Run("Send Coin", TestE2E_SendCoin)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Наш токен
	token, _ := getAuthToken(t)

	url := "http://localhost:8080/api/info"

	// Создаем GET‑запрос.
	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.NoError(t, err)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Отправляем запрос.
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Тело ответа
	bb, _ := io.ReadAll(resp.Body)
	t.Log(string(bb))

	// Проверяем, что код ответа 200.
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.NoError(t, err)
}
