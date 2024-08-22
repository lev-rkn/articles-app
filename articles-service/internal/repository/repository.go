package repository

import (
	"articles-service/internal/config"
	"articles-service/internal/lib/utils"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Repository struct {
	Article ArticleRepoInterface
	Comment CommentRepoInterface
}

func NewRepository(ctx context.Context) *Repository {
	pool, err := pgxpool.New(ctx, config.Cfg.Postgres.PgUrl)
	if err != nil {
		utils.ErrorLog("Unable to connect to database", err)
	}
	pool.Config().MaxConns = config.Cfg.Postgres.MaxConnections
	// запуск миграций
	m, err := migrate.New(config.Cfg.Postgres.MigrationsDir, config.Cfg.Postgres.PgUrl)
	if err != nil {
		utils.ErrorLog("new migrations", err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		utils.ErrorLog("migrations up", err)
	}

	rdb := redis.NewClient(&redis.Options{
        Addr:     config.Cfg.Redis.Address,
        Password: config.Cfg.Redis.Password,
        DB:       config.Cfg.Redis.DB,
    })

	var repository = &Repository{
		Article: NewArticleRepo(ctx, pool, rdb),
		Comment: NewCommentRepo(ctx, pool, rdb),
	}

	return repository
}
