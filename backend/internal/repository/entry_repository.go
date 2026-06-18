package repository

import (
	"context"
	"errors"
	"time"

	"cesizen/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EntryRepository struct {
	db *pgxpool.Pool
}

func NewEntryRepository(db *pgxpool.Pool) *EntryRepository {
	return &EntryRepository{db: db}
}

func (r *EntryRepository) FindByUser(ctx context.Context, userID string, from, to time.Time) ([]domain.Entry, error) {
	rows, err := r.db.Query(ctx, `
		SELECT e.id, e.user_id, e.emotion_id, em.label, p.label,
		       e.intensity, COALESCE(e.comment, ''), e.entry_date, e.created_at
		FROM entries e
		JOIN emotions em ON em.id = e.emotion_id
		JOIN primary_emotions p ON p.id = em.primary_emotion_id
		WHERE e.user_id = $1 AND e.entry_date >= $2 AND e.entry_date <= $3
		ORDER BY e.entry_date DESC
	`, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Entry
	for rows.Next() {
		var e domain.Entry
		if err := rows.Scan(&e.ID, &e.UserID, &e.EmotionID, &e.EmotionLabel, &e.PrimaryLabel,
			&e.Intensity, &e.Comment, &e.EntryDate, &e.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (r *EntryRepository) FindByID(ctx context.Context, id int) (*domain.Entry, error) {
	e := &domain.Entry{}
	err := r.db.QueryRow(ctx, `
		SELECT e.id, e.user_id, e.emotion_id, em.label, p.label,
		       e.intensity, COALESCE(e.comment, ''), e.entry_date, e.created_at
		FROM entries e
		JOIN emotions em ON em.id = e.emotion_id
		JOIN primary_emotions p ON p.id = em.primary_emotion_id
		WHERE e.id = $1
	`, id).Scan(&e.ID, &e.UserID, &e.EmotionID, &e.EmotionLabel, &e.PrimaryLabel,
		&e.Intensity, &e.Comment, &e.EntryDate, &e.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return e, err
}

func (r *EntryRepository) Create(ctx context.Context, e *domain.Entry) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO entries (user_id, emotion_id, intensity, comment, entry_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`, e.UserID, e.EmotionID, e.Intensity, e.Comment, e.EntryDate).
		Scan(&e.ID, &e.CreatedAt)
}

func (r *EntryRepository) Update(ctx context.Context, e *domain.Entry) error {
	_, err := r.db.Exec(ctx, `
		UPDATE entries SET emotion_id = $1, intensity = $2, comment = $3, entry_date = $4
		WHERE id = $5
	`, e.EmotionID, e.Intensity, e.Comment, e.EntryDate, e.ID)
	return err
}

func (r *EntryRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM entries WHERE id = $1`, id)
	return err
}

func (r *EntryRepository) StatsForUser(ctx context.Context, userID string, from, to time.Time) ([]domain.EmotionStat, error) {
	rows, err := r.db.Query(ctx, `
		SELECT e.emotion_id, em.label, p.label, COUNT(*) AS count
		FROM entries e
		JOIN emotions em ON em.id = e.emotion_id
		JOIN primary_emotions p ON p.id = em.primary_emotion_id
		WHERE e.user_id = $1 AND e.entry_date >= $2 AND e.entry_date <= $3
		GROUP BY e.emotion_id, em.label, p.label
		ORDER BY count DESC
	`, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.EmotionStat
	for rows.Next() {
		var s domain.EmotionStat
		if err := rows.Scan(&s.EmotionID, &s.EmotionLabel, &s.PrimaryLabel, &s.Count); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}
