package http

import (
	"github.com/ca11ou5/slogging"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

type Service interface{}

type Server struct {
	svc Service
}

func NewServer(svc Service) *Server {
	return &Server{
		svc: svc,
	}
}

func (s *Server) StartListening() error {

}

func (s *Server) initRouter() *mux.Router {
	router := mux.NewRouter()

	tracemw := slogging.MuxHTTPTraceMiddleware(slog.Default())

	api := router.PathPrefix("/api").Subrouter()
	{
		api.HandleFunc("/auth").Methods(http.MethodPost)
		api.HandleFunc("/sendCoin").Methods(http.MethodPost)
		api.HandleFunc("/info").Methods(http.MethodGet)
		api.HandleFunc("/buy/{item}").Methods(http.MethodGet)
	}
	api.Use(tracemw)

	return router
}
