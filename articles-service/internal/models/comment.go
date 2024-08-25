package models

import (
	"encoding/json"
	"time"
)

// представляет сущность комментария
type Comment struct {
	Id        int       `json:"id,omitempty" db:"id"`
	UserId    int       `json:"user_id,omitempty" db:"user_id"`
	ArticleId int       `json:"article_id,omitempty" db:"article_id" validate:"required"`
	Text      string    `json:"text,omitempty" db:"text" validate:"required,max=500"`
	Timestamp time.Time `json:"timestamp,omitempty" db:"timestamp"`
}

func (s *Comment) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Comment) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}
