package http

import (
	"context"
	"fmt"
	"github.com/ca11ou5/avito-trainee-assignment/internal/entity"
	"github.com/ca11ou5/avito-trainee-assignment/internal/payload"
	"github.com/ca11ou5/slogging"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
)

type Service interface {
	AuthenticateUser(ctx context.Context, req payload.AuthRequest) (string, error)
	ExtractUserInfo(ctx context.Context, token string) (entity.EmployeeInfo, error)
}

type Server struct {
	svc Service

	Port string
}

func NewServer(svc Service, port string) *Server {
	return &Server{
		svc: svc,

		Port: port,
	}
}

func (s *Server) StartListening() error {
	router := s.initRouter()

	return http.ListenAndServe(fmt.Sprintf(":%s", s.Port), router)
}

func (s *Server) initRouter() *mux.Router {
	router := mux.NewRouter()

	traceMiddleware := slogging.MuxHTTPTraceMiddleware(slog.Default())

	api := router.PathPrefix("/api").Subrouter()
	{
		api.HandleFunc("/auth", s.auth).Methods(http.MethodPost)

		protectedAPI := api.PathPrefix("/").Subrouter()
		{
			protectedAPI.HandleFunc("/info", s.getInfo).Methods(http.MethodGet)
			protectedAPI.HandleFunc("/sendCoin", s.sendCoin).Methods(http.MethodPost)
			protectedAPI.HandleFunc("/buy/{item}", s.buyItem).Methods(http.MethodGet)
		}
		protectedAPI.Use(authMiddleware)
	}
	api.Use(traceMiddleware)

	return router
}
