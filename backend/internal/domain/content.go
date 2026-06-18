package domain

import "time"

type Content struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Author      string    `json:"author"`
	IsPublished bool      `json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateContentInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

type UpdateContentInput struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	Author      string `json:"author"`
	IsPublished bool   `json:"is_published"`
}
