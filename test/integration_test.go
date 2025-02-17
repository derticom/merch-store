package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	baseURL = "http://localhost:8080"
)

func TestPurchaseItem(t *testing.T) {
	// Перед тестом выполнить запуск контейнеров
	// docker compose up --build

	username := "testuser"
	password := "testpassword"
	userID, token := registerUser(t, username, password)
	assert.NotEmpty(t, userID, "User ID should not be empty")
	assert.NotEmpty(t, token, "Token should not be empty")

	// Покупка товара.
	item := "book"
	buyItem(t, token, item)

	// Проверка баланса.
	coins := getUserInfo(t, token)
	assert.Equal(t, 950, coins, "User should have 950 coins after buying a book")
}

// Регистрация пользователя
func registerUser(t *testing.T, username, password string) (string, string) {
	url := fmt.Sprintf("%s/api/auth/register", baseURL)
	payload := map[string]string{
		"username": username,
		"password": password,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err, "Failed to register user")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Failed to register user")

	var result struct {
		Token string `json:"token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err, "Failed to decode response")

	userID := "testuser"
	return userID, result.Token
}

// Покупка товара
func buyItem(t *testing.T, token, item string) {
	url := fmt.Sprintf("%s/api/buy/%s", baseURL, item)
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err, "Failed to create request")

	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err, "Failed to buy item")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Failed to buy item")
}

// Получение информации о пользователе
func getUserInfo(t *testing.T, token string) int {
	url := fmt.Sprintf("%s/api/info", baseURL)
	req, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err, "Failed to create request")

	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err, "Failed to get user info")
	defer resp.Body.Close()

	var result struct {
		Coins int `json:"coins"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err, "Failed to decode response")

	return result.Coins
}
