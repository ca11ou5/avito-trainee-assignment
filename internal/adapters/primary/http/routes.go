package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ca11ou5/avito-trainee-assignment/internal/payload"
	"github.com/ca11ou5/avito-trainee-assignment/internal/service"
	"io"
	"net/http"
)

var (
	errEmptyBody  = errors.New("empty request body")
	errFailedAuth = errors.New("failed to authenticate employee")
)

func (s *Server) getInfo(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(contextTokenKey).(string)

	info, err := s.svc.ExtractUserInfo(r.Context(), token)
	if err != nil {
		// TODO:
	}

}

func (s *Server) sendCoin(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) buyItem(w http.ResponseWriter, r *http.Request) {

	//item := mux.Vars(r)["item"]
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

	var req payload.AuthRequest
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

	token, err := s.svc.AuthenticateUser(r.Context(), req)
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
