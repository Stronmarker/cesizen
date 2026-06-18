package repository

import (
	"context"
	"errors"

	"cesizen/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContentRepository struct {
	db *pgxpool.Pool
}

func NewContentRepository(db *pgxpool.Pool) *ContentRepository {
	return &ContentRepository{db: db}
}

func (r *ContentRepository) FindAllPublished(ctx context.Context) ([]domain.Content, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, title, content, author, is_published, created_at, updated_at
		FROM contents WHERE is_published = true ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanContents(rows)
}

func (r *ContentRepository) FindAll(ctx context.Context) ([]domain.Content, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, title, content, author, is_published, created_at, updated_at
		FROM contents ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanContents(rows)
}

func (r *ContentRepository) FindByID(ctx context.Context, id int) (*domain.Content, error) {
	c := &domain.Content{}
	err := r.db.QueryRow(ctx, `
		SELECT id, title, content, author, is_published, created_at, updated_at
		FROM contents WHERE id = $1
	`, id).Scan(&c.ID, &c.Title, &c.Content, &c.Author, &c.IsPublished, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return c, err
}

func (r *ContentRepository) Create(ctx context.Context, c *domain.Content) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO contents (title, content, author)
		VALUES ($1, $2, $3)
		RETURNING id, is_published, created_at, updated_at
	`, c.Title, c.Content, c.Author).Scan(&c.ID, &c.IsPublished, &c.CreatedAt, &c.UpdatedAt)
}

func (r *ContentRepository) Update(ctx context.Context, c *domain.Content) error {
	_, err := r.db.Exec(ctx, `
		UPDATE contents SET title = $1, content = $2, author = $3, is_published = $4, updated_at = now()
		WHERE id = $5
	`, c.Title, c.Content, c.Author, c.IsPublished, c.ID)
	return err
}

func (r *ContentRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM contents WHERE id = $1`, id)
	return err
}

func scanContents(rows interface {
	Next() bool
	Scan(...any) error
	Err() error
}) ([]domain.Content, error) {
	var contents []domain.Content
	for rows.Next() {
		var c domain.Content
		if err := rows.Scan(&c.ID, &c.Title, &c.Content, &c.Author, &c.IsPublished, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		contents = append(contents, c)
	}
	return contents, rows.Err()
}
