package http

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/ca11ou5/avito-trainee-assignment/internal/adapters/secondary/postgres"
	"github.com/ca11ou5/avito-trainee-assignment/internal/models"
	"github.com/ca11ou5/avito-trainee-assignment/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func withToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenKey{}, token)
}

type Svc struct {
	mock.Mock
}

func (s *Svc) ExtractUserInfo(ctx context.Context, token string) (models.EmployeeInfo, error) {
	args := s.Called(ctx, token)
	return args.Get(0).(models.EmployeeInfo), args.Error(1)
}

func (s *Svc) SendCoin(ctx context.Context, token string, req models.SentTransaction) error {
	args := s.Called(ctx, token, req)
	return args.Error(0)
}

func (s *Svc) AuthenticateUser(ctx context.Context, creds models.Credentials) (string, error) {
	return "", nil
}

func (s *Svc) BuyItem(ctx context.Context, token string, item string) error {
	return nil
}

func TestGetInfo_Success(t *testing.T) {
	mockSvc := new(Svc)
	expectedInfo := models.EmployeeInfo{
		Coins: 1000,
	}
	testToken := "validtoken"
	mockSvc.
		On("ExtractUserInfo", mock.Anything, testToken).
		Return(expectedInfo, nil)

	// Создаем сервер с этим сервисом.
	srv := &Server{
		svc: mockSvc,
	}

	// Создаем тестовый запрос и устанавливаем токен в контекст.
	req := httptest.NewRequest("GET", "http://example.com/info", nil)
	req = req.WithContext(withToken(req.Context(), testToken))
	rr := httptest.NewRecorder()

	srv.getInfo(rr, req)

	// Проверяем, что статус ответа 200 и тело содержит ожидаемые данные.
	assert.Equal(t, http.StatusOK, rr.Code)

	var actualInfo models.EmployeeInfo
	err := json.Unmarshal(rr.Body.Bytes(), &actualInfo)
	assert.NoError(t, err)
	assert.Equal(t, expectedInfo, actualInfo)

	mockSvc.AssertExpectations(t)
}

func TestGetInfo_InvalidToken(t *testing.T) {
	mockSvc := new(Svc)
	testToken := "badtoken"

	mockSvc.
		On("ExtractUserInfo", mock.Anything, testToken).
		Return(models.EmployeeInfo{}, service.ErrInvalidToken)

	srv := &Server{
		svc: mockSvc,
	}

	req := httptest.NewRequest("GET", "http://example.com/info", nil)
	req = req.WithContext(withToken(req.Context(), testToken))
	rr := httptest.NewRecorder()

	srv.getInfo(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), service.ErrInvalidToken.Error())

	mockSvc.AssertExpectations(t)
}

func TestSendCoin_Success(t *testing.T) {
	mockSvc := new(Svc)
	testToken := "validtoken"
	reqData := models.SentTransaction{
		Amount: 50,
		ToUser: "receiver",
	}
	// SendCoin должен вернуть nil.
	mockSvc.
		On("SendCoin", mock.Anything, testToken, reqData).
		Return(nil)

	srv := &Server{
		svc: mockSvc,
	}

	// Подготавливаем тело запроса.
	bodyBytes, _ := json.Marshal(reqData)
	req := httptest.NewRequest("POST", "http://example.com/sendCoin", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(withToken(req.Context(), testToken))
	rr := httptest.NewRecorder()

	srv.sendCoin(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockSvc.AssertExpectations(t)
}

func TestSendCoin_InvalidJSON(t *testing.T) {
	srv := &Server{}

	req := httptest.NewRequest("POST", "http://example.com/sendCoin", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	req = req.WithContext(withToken(req.Context(), "sometoken"))
	rr := httptest.NewRecorder()

	srv.sendCoin(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), errEmptyBody.Error())
}

func TestSendCoin_ValidationError(t *testing.T) {
	// Тестируем случай когда метод Validate() возвращает ошибку.
	reqData := models.SentTransaction{
		Amount: 0,
		ToUser: "receiver",
	}
	bodyBytes, _ := json.Marshal(reqData)
	srv := &Server{}

	req := httptest.NewRequest("POST", "http://example.com/sendCoin", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(withToken(req.Context(), "sometoken"))
	rr := httptest.NewRecorder()

	srv.sendCoin(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSendCoin_ServiceError(t *testing.T) {
	mockSvc := new(Svc)
	testToken := "validtoken"
	reqData := models.SentTransaction{
		Amount: 50,
		ToUser: "receiver",
	}
	bodyBytes, _ := json.Marshal(reqData)
	req := httptest.NewRequest("POST", "http://example.com/sendCoin", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(withToken(req.Context(), testToken))
	rr := httptest.NewRecorder()

	// Пример: если сервис возвращает ErrNotEnoughBalance (или ErrCantSentToYourself, или ErrEmployeeNotExists), мы ожидаем 400
	mockSvc.
		On("SendCoin", mock.Anything, testToken, reqData).
		Return(postgres.ErrNotEnoughBalance)

	srv := &Server{
		svc: mockSvc,
	}

	srv.sendCoin(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), postgres.ErrNotEnoughBalance.Error())

	mockSvc.AssertExpectations(t)
}
