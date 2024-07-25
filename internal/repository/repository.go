package repository

import (
	"ads-service/internal/repository/postgres"
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/knadh/koanf/v2"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Repository struct {
	Ad AdRepo
}

func NewRepository(logger *slog.Logger, cfg *koanf.Koanf) *Repository {
	// подключение к postgres
	conn, err := pgx.Connect(context.Background(), cfg.String("pg_url"))
	if err != nil {

		logger.Error("Unable to connect to database",
			"err", err.Error())
	}

	// запуск миграций
	m, err := migrate.New(cfg.String("migrations_dir"), cfg.String("pg_url"))
	if err != nil {
		logger.Error("new migrations",
			"err", err.Error())
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Error("migrations up",
			"err", err.Error())
	}

	var repository = &Repository{
		Ad: postgres.NewAdPgRepo(conn),
	}

	return repository
}
