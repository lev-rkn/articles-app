package models

import "time"

// представляет сущность статьи
type Article struct {
	Id          int       `json:"id,omitempty" db:"id"`
	UserId      int       `json:"user_id,omitempty" db:"user_id"`
	Title       string    `json:"title,omitempty" db:"title" validate:"required,max=140"`
	Description string    `json:"description,omitempty" db:"description" validate:"required,max=1000"`
	Photos      []string  `json:"photos,omitempty" db:"photos" validate:"required,max=3"`
	Timestamp   time.Time `json:"timestamp,omitempty" db:"timestamp"`
}
