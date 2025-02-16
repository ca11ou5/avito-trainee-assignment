package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestE2E_SendCoin(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e tests")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Наш токен
	token, _ := getAuthToken(t)

	// Юзернейм другого пользователя, которому мы будем отправлять монеты
	body := map[string]interface{}{
		"toUser": "ca11ou5",
		"amount": 200,
	}
	bb, err := json.Marshal(body)

	url := "http://localhost:8080/api/sendCoin"

	// Создаем GET‑запрос.
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bb))
	assert.NoError(t, err)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Отправляем запрос.
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Проверяем, что код ответа 200.
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.NoError(t, err)
}
