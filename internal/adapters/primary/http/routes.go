package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ca11ou5/avito-trainee-assignment/internal/adapters/secondary/postgres"
	"github.com/ca11ou5/avito-trainee-assignment/internal/models"
	"github.com/ca11ou5/avito-trainee-assignment/internal/service"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"time"
)

var (
	errEmptyBody            = errors.New("empty request body")
	errFailedAuth           = errors.New("failed to authenticate employee")
	errMissingPathParameter = errors.New("missing path 'item' parameter")
)

const defaultRequestTimeout = 1 * time.Second

func (s *Server) getInfo(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), defaultRequestTimeout)
	defer cancel()

	token := contextToken(ctx)

	info, err := s.svc.ExtractUserInfo(ctx, token)
	switch {
	case errors.Is(err, service.ErrInvalidToken):
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(beatifyError(err))
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(beatifyError(err))
		return
	default:
		respbb, _ := json.Marshal(info)
		w.WriteHeader(http.StatusOK)
		w.Write(respbb)
		return
	}
}

func (s *Server) sendCoin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bb, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(beatifyError(err))
		return
	}
	defer r.Body.Close()

	var req models.SentTransaction
	err = json.Unmarshal(bb, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(beatifyError(fmt.Errorf("%s: %s", errEmptyBody, err)))
		return
	}

	err = req.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(beatifyError(err))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), defaultRequestTimeout)
	defer cancel()

	token := contextToken(ctx)

	err = s.svc.SendCoin(ctx, token, req)
	switch {
	case errors.Is(err, service.ErrInvalidToken):
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(beatifyError(err))
		return
	case errors.Is(err, service.ErrCantSentToYourself) || errors.Is(err, postgres.ErrEmployeeNotExists) || errors.Is(err, postgres.ErrNotEnoughBalance):
		w.WriteHeader(http.StatusBadRequest)
		w.Write(beatifyError(err))
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(beatifyError(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) buyItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	item := mux.Vars(r)["item"]
	if item == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(beatifyError(errMissingPathParameter))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), defaultRequestTimeout)
	defer cancel()

	token := contextToken(ctx)

	err := s.svc.BuyItem(ctx, token, item)
	switch {
	case errors.Is(err, service.ErrInvalidToken) || errors.Is(err, postgres.ErrEmployeeNotExists):
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(beatifyError(err))
		return
	case errors.Is(err, postgres.ErrMerchNotExists) || errors.Is(err, postgres.ErrNotEnoughBalance):
		w.WriteHeader(http.StatusBadRequest)
		w.Write(beatifyError(err))
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(beatifyError(err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) auth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bb, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(beatifyError(err))
		return
	}
	defer r.Body.Close()

	var creds models.Credentials
	err = json.Unmarshal(bb, &creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(beatifyError(fmt.Errorf("%s: %s", errEmptyBody, err)))
		return
	}

	err = creds.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(beatifyError(err))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), defaultRequestTimeout)
	defer cancel()

	token, err := s.svc.AuthenticateUser(ctx, creds)
	switch {
	case errors.Is(err, service.ErrWrongPassword):
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(beatifyError(err))
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(beatifyError(fmt.Errorf("%s: %s", errFailedAuth, err)))
		return
	default:
		w.WriteHeader(http.StatusOK)
		w.Write(beatifyToken(token))
		return
	}
}
