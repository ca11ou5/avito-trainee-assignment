package http

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) getInfo(w http.ResponseWriter, r *http.Request) {}

func (s *Server) sendCoin(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) buyItem(w http.ResponseWriter, r *http.Request) {
	item := mux.Vars(r)["item"]
}

func (s *Server) auth(w http.ResponseWriter, r *http.Request) {}
