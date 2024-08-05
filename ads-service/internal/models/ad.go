package models

import "time"

type Ad struct {
	Id          int       `json:"id,omitempty" swaggerignore:"true"`
	Title       string    `json:"title,omitempty" validate:"required,max=200"`
	Description string    `json:"description,omitempty" validate:"max=1000"`
	Photos      []string  `json:"photos,omitempty" validate:"required,max=3"`
	Price       uint       `json:"price,omitempty" validate:"required"`
	UserId      int       `json:"user_id,omitempty"`
	Timestamp   time.Time `json:"timestamp,omitempty" swaggerignore:"true"`
}
