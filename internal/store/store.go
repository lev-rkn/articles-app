package store

import (
	"ads-service/internal/store/postgres"
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/knadh/koanf/v2"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Store struct {
	Ad AdRepo
}

func NewRepository(logger *slog.Logger, cfg *koanf.Koanf) *Store {
	// connect to postgres
	conn, err := pgx.Connect(context.Background(), cfg.MustString("pg_url"))
	if err != nil {

		logger.Error("Unable to connect to database",
			"err", err.Error())
	}

	// load migrations
	m, err := migrate.New(cfg.MustString("migrations_dir"), cfg.MustString("pg_url"))
	if err != nil {
		logger.Error("new migrations",
			"err", err.Error())
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Error("migrations up",
			"err", err.Error())
	}

	var store = &Store{
		Ad: postgres.NewAdPgRepo(conn),
	}

	return store
}
