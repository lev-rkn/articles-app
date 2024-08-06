package repository

import (
	"articles-service/internal/config"
	"articles-service/internal/lib/utils"
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Repository struct {
	Article ArticleRepoInterface
	Comment CommentRepoInterface
}

func NewRepository(ctx context.Context) *Repository {
	// подключение к postgres
	conn, err := pgx.Connect(context.Background(), config.Cfg.PgUrl)
	if err != nil {
		utils.ErrorLog("Unable to connect to database", err)
	}

	// запуск миграций
	m, err := migrate.New(config.Cfg.MigrationsDir, config.Cfg.PgUrl)
	if err != nil {
		utils.ErrorLog("new migrations", err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		utils.ErrorLog("migrations up", err)
	}

	var repository = &Repository{
		Article: NewArticleRepo(ctx, conn),
		Comment: NewCommentRepo(ctx, conn),
	}

	return repository
}
