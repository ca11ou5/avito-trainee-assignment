package main

import (
	"github.com/ca11ou5/avito-trainee-assignment/configs"
	"github.com/ca11ou5/avito-trainee-assignment/internal/adapters/primary/http"
	"github.com/ca11ou5/avito-trainee-assignment/internal/adapters/secondary/auth"
	"github.com/ca11ou5/avito-trainee-assignment/internal/adapters/secondary/postgres"
	"github.com/ca11ou5/avito-trainee-assignment/internal/service"
	"github.com/ca11ou5/slogging"
	"github.com/ilyakaznacheev/cleanenv"

	"log/slog"
	"os"
)

func main() {
	var cfg configs.Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		slog.Error("failed to read env",
			slogging.ErrAttr(err))
		os.Exit(1)
	}

	slog.Debug("config",
		slogging.AnyAttr("data", cfg))

	log := slogging.NewLogger(
		slogging.SetLevel(cfg.LogLevel),
		slogging.WithSource(true),
		slogging.SetDefault(true),
	)

	pg := postgres.NewAdapter(cfg.PostgresURL)
	authProvider := auth.NewAdapter(cfg.JWTSalt)

	svc := service.New(pg, authProvider)

	srv := http.NewServer(svc)

	log.Info("starting http server")
	err = srv.StartListening(cfg.Port)
	if err != nil {
		log.Error("failed to start http server",
			slogging.ErrAttr(err))
		os.Exit(1)
	}
}
