package http

import (
	"github.com/ca11ou5/avito-trainee-assignment/internal/models"
	"github.com/ca11ou5/slogging"
	"github.com/gorilla/mux"

	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Service interface {
	// AuthenticateUser POST /api/auth
	AuthenticateUser(ctx context.Context, creds models.Credentials) (string, error)

	// ExtractUserInfo GET /api/info
	ExtractUserInfo(ctx context.Context, token string) (models.EmployeeInfo, error)

	// SendCoin POST /api/sendCoin
	SendCoin(ctx context.Context, token string, trans models.SentTransaction) error

	// BuyItem GET /api/buy/{item}
	BuyItem(ctx context.Context, token string, item string) error
}

type Server struct {
	svc Service
}

func NewServer(svc Service) *Server {
	return &Server{
		svc: svc,
	}
}

func (s *Server) StartListening(port string) error {
	router := s.initRouter()

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           router,
		ReadHeaderTimeout: 3 * time.Second,
	}

	return server.ListenAndServe()
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
