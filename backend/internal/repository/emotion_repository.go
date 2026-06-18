package repository

import (
	"context"
	"errors"

	"cesizen/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EmotionRepository struct {
	db *pgxpool.Pool
}

func NewEmotionRepository(db *pgxpool.Pool) *EmotionRepository {
	return &EmotionRepository{db: db}
}

func (r *EmotionRepository) FindAllPrimary(ctx context.Context) ([]domain.PrimaryEmotion, error) {
	rows, err := r.db.Query(ctx, `SELECT id, label, is_active FROM primary_emotions ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.PrimaryEmotion
	for rows.Next() {
		var p domain.PrimaryEmotion
		if err := rows.Scan(&p.ID, &p.Label, &p.IsActive); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *EmotionRepository) FindPrimaryByID(ctx context.Context, id int) (*domain.PrimaryEmotion, error) {
	p := &domain.PrimaryEmotion{}
	err := r.db.QueryRow(ctx, `SELECT id, label, is_active FROM primary_emotions WHERE id = $1`, id).
		Scan(&p.ID, &p.Label, &p.IsActive)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return p, err
}

func (r *EmotionRepository) CreatePrimary(ctx context.Context, p *domain.PrimaryEmotion) error {
	return r.db.QueryRow(ctx,
		`INSERT INTO primary_emotions (label) VALUES ($1) RETURNING id, is_active`,
		p.Label,
	).Scan(&p.ID, &p.IsActive)
}

func (r *EmotionRepository) UpdatePrimary(ctx context.Context, p *domain.PrimaryEmotion) error {
	_, err := r.db.Exec(ctx,
		`UPDATE primary_emotions SET label = $1, is_active = $2 WHERE id = $3`,
		p.Label, p.IsActive, p.ID,
	)
	return err
}

func (r *EmotionRepository) DeletePrimary(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM primary_emotions WHERE id = $1`, id)
	return err
}

func (r *EmotionRepository) FindAll(ctx context.Context) ([]domain.Emotion, error) {
	rows, err := r.db.Query(ctx, `
		SELECT e.id, e.label, e.primary_emotion_id, p.label, e.is_active
		FROM emotions e
		JOIN primary_emotions p ON p.id = e.primary_emotion_id
		ORDER BY p.id, e.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanEmotions(rows)
}

func (r *EmotionRepository) FindByID(ctx context.Context, id int) (*domain.Emotion, error) {
	e := &domain.Emotion{}
	err := r.db.QueryRow(ctx, `
		SELECT e.id, e.label, e.primary_emotion_id, p.label, e.is_active
		FROM emotions e
		JOIN primary_emotions p ON p.id = e.primary_emotion_id
		WHERE e.id = $1
	`, id).Scan(&e.ID, &e.Label, &e.PrimaryEmotionID, &e.PrimaryLabel, &e.IsActive)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return e, err
}

func (r *EmotionRepository) CreateEmotion(ctx context.Context, e *domain.Emotion) error {
	return r.db.QueryRow(ctx,
		`INSERT INTO emotions (label, primary_emotion_id) VALUES ($1, $2) RETURNING id, is_active`,
		e.Label, e.PrimaryEmotionID,
	).Scan(&e.ID, &e.IsActive)
}

func (r *EmotionRepository) UpdateEmotion(ctx context.Context, e *domain.Emotion) error {
	_, err := r.db.Exec(ctx,
		`UPDATE emotions SET label = $1, primary_emotion_id = $2, is_active = $3 WHERE id = $4`,
		e.Label, e.PrimaryEmotionID, e.IsActive, e.ID,
	)
	return err
}

func (r *EmotionRepository) DeleteEmotion(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM emotions WHERE id = $1`, id)
	return err
}

func scanEmotions(rows interface {
	Next() bool
	Scan(...any) error
	Err() error
}) ([]domain.Emotion, error) {
	var out []domain.Emotion
	for rows.Next() {
		var e domain.Emotion
		if err := rows.Scan(&e.ID, &e.Label, &e.PrimaryEmotionID, &e.PrimaryLabel, &e.IsActive); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
