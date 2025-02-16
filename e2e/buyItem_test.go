package e2e

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestE2E_BuyItem(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e tests")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Предположим, что мы покупаем предмет "cup".
	url := "http://localhost:8080/api/buy/cup"

	// Создаем GET‑запрос.
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)

	// Добавляем заголовок авторизации с валидным токеном.
	token, _ := getAuthToken(t)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Отправляем запрос.
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Проверяем, что код ответа 200.
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.NoError(t, err)
}
