package repository

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/knadh/koanf/v2"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Repository struct {
	Ad AdRepoInterface
}

func NewRepository(ctx context.Context, cfg *koanf.Koanf) *Repository {
	pgUrl := cfg.String("pg_url")
	// подключение к postgres
	conn, err := pgx.Connect(context.Background(), pgUrl)
	if err != nil {

		slog.Error("Unable to connect to database",
			"err", err.Error())
	}
	
	// запуск миграций
	m, err := migrate.New(cfg.String("migrations_dir"), pgUrl)
	if err != nil {
		slog.Error("new migrations",
			"err", err.Error())
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("migrations up",
			"err", err.Error())
	}

	var repository = &Repository{
		Ad: NewAdRepo(ctx, conn),
	}

	return repository
}
