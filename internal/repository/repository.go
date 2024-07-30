package repository

import (
	"ads-service/internal/config"
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Repository struct {
	Ad AdRepoInterface
}

func NewRepository(ctx context.Context, cfg *config.Config) *Repository {
	// подключение к postgres
	conn, err := pgx.Connect(context.Background(), cfg.PgUrl)
	if err != nil {

		slog.Error("Unable to connect to database",
			"err", err.Error())
	}
	
	// запуск миграций
	m, err := migrate.New(cfg.MigrationsDir, cfg.PgUrl)
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
