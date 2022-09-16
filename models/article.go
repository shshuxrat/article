package models

import "time"

type Content struct {
	Title string `json:"title"`
	Body  string `json:"b"`
}

type Article struct {
	ID        int        `json:"id"`
	Content              // Promoted fields
	Author    Person     `json:"p"` // Nested structs
	CreatedAt *time.Time `json:"created_at"`
}

type ArticleCreateModel struct {
	Content      // Promoted fields
	AuthorID int `json:"author_id"`
}

type ArticleUpdateModel struct {
	ID       int `json:"id"`
	Content      // Promoted fields
	AuthorID int `json:"author_id"`
}
