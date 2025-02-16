package postgres

import (
	"github.com/ca11ou5/slogging"
	"github.com/jmoiron/sqlx"

	// postgres driver, need sslmode=disable
	_ "github.com/lib/pq"

	"log/slog"
	"os"
)

type Adapter struct {
	db *sqlx.DB
}

func NewAdapter(postgresURL string) *Adapter {
	return &Adapter{
		db: connect(postgresURL),
	}
}

func connect(connstring string) *sqlx.DB {
	db, err := sqlx.Connect("postgres", connstring)
	if err != nil {
		slog.Error("failed connect to postgres",
			slogging.ErrAttr(err))
		os.Exit(1)
	}

	err = db.Ping()
	if err != nil {
		slog.Error("failed ping postgres",
			slogging.ErrAttr(err))
		os.Exit(1)
	}

	return db
}
