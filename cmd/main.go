package main

import (
	"github.com/ca11ou5/avito-trainee-assignment/configs"
	"github.com/ca11ou5/avito-trainee-assignment/internal/adapters/primary/http"
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

	// TODO:
	_ = slogging.NewLogger(
		slogging.SetLevel(cfg.LogLevel),
		slogging.WithSource(true),
		slogging.SetDefault(true),
	)

	pg := postgres.NewAdapter(cfg.PostgresURL)

	svc := service.New(pg)

	srv := http.NewServer(svc)

	// TODO:
	srv.StartListening()
}
