package e2e

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomString(n int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

var userTestCreds = map[string]string{}

// Махинации чтобы тесты проводились с одним пользователем
func getUsername() string {
	if v, ok := userTestCreds["username"]; ok {
		return v
	}

	username := randomString(16)
	userTestCreds["username"] = username
	return username
}

// Махинации чтобы тесты проводились с одним пользователем
func getPassword() string {
	if v, ok := userTestCreds["password"]; ok {
		return v
	}

	password := randomString(12)
	userTestCreds["password"] = password
	return password
}

func getAuthToken(t *testing.T) (string, string) {
	client := &http.Client{Timeout: 10 * time.Second}
	url := "http://localhost:8080/api/auth"

	body := map[string]string{
		"username": getUsername(),
		"password": getPassword(),
	}
	bb, err := json.Marshal(body)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bb))
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respbb, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	// Ответ выглядит как "token" : "example"
	var authResp struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(respbb, &authResp)
	assert.NoError(t, err)

	return authResp.Token, body["username"]
}

func TestE2E_Auth(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e tests")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	url := "http://localhost:8080/api/auth"

	body := map[string]string{
		"username": "ca11ou5",
		"password": "12345678",
	}
	bb, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bb))
	assert.NoError(t, err)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Проверяем, что код ответа 200.
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respbb, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Contains(t, string(respbb), "token")
	t.Logf(string(respbb))
}
