package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	CreateArticleRequest = promauto.NewCounter(prometheus.CounterOpts{
		Name: "create_article_request_count",
		Help: "Количество запросов на создание Статьи",
	})
	CreateArticleOK = promauto.NewCounter(prometheus.CounterOpts{
		Name: "create_article_ok_count",
		Help: "Количество успешно созданных Статей",
	})
	CreateArticleError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "create_article_error_count",
		Help: "Количестов ошибок при создании Статьи",
	})

	GetArticleRequest = promauto.NewCounter(prometheus.CounterOpts{
		Name: "get_article_request_count",
		Help: "Количество запросов на получение Статьи",
	})
	GetArticleOK = promauto.NewCounter(prometheus.CounterOpts{
		Name: "get_article_ok_count",
		Help: "Количество успешно полученных Статей",
	})
	GetArticleError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "get_article_error_count",
		Help: "Количестов ошибок при получении Статьи",
	})

	GetArticlesRequest = promauto.NewCounter(prometheus.CounterOpts{
		Name: "get_articles_request_count",
		Help: "Количество запросов на получение Статей",
	})
	GetArticlesOK = promauto.NewCounter(prometheus.CounterOpts{
		Name: "get_articles_ok_count",
		Help: "Количество успешно полученных Статей",
	})
	GetArticlesError = promauto.NewCounter(prometheus.CounterOpts{
		Name: "get_articles_error_count",
		Help: "Количестов ошибок при получении Статей",
	})
)

