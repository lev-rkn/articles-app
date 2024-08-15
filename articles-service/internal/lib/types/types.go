package types

import "errors"

var (
	ErrInvalidPageNumber = errors.New("невалидный номер страницы")
	ErrInvalidUserId     = errors.New("невалидный идентификатор пользователя")
	ErrInvalidDateSort   = errors.New("невалидный параметр сортировки по дате")
	ErrInvalidArticleId  = errors.New("невалидный идентификатор (id) статьи")
	ErrInvalidToken      = errors.New("невалидный токен доступа")
	ErrArticleNotFound   = errors.New("статья не найдена")
	ErrNoComments        = errors.New("комментарии отсутствуют")
)

var (
	KeyUser  = "user"
	KeyError = "error"
)
