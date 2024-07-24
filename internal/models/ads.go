package models

type Ad struct {
	Id          int      `json:"id,omitempty"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Price       int      `json:"price,omitempty"`
	Photos      []string `json:"photos,omitempty"`
	Timestamp   string   `json:"timestamp,omitempty"`
}
