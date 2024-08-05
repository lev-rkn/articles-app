package types

import "errors"

var (
	ErrInvalidPageNumber = errors.New("невалидный номер страницы")
	ErrInvalidUserId     = errors.New("невалидный идентификатор пользователя")
	ErrInvalidPriceSort  = errors.New("невалидный параметр сортировки по цене")
	ErrInvalidDateSort   = errors.New("невалидный параметр сортировки по дате")
	ErrInvalidId         = errors.New("невалидный идентификатор (id) объявления")
	ErrInvalidToken      = errors.New("invalid token")
	ErrArticleNotFound   = errors.New("статья не найдена")
)

var (
	KeyUser  = "user"
	KeyError = "error"
)
