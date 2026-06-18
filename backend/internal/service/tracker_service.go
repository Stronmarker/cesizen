package service

import (
	"context"
	"errors"
	"time"

	"cesizen/internal/domain"
)

type EntryRepo interface {
	FindByUser(ctx context.Context, userID string, from, to time.Time) ([]domain.Entry, error)
	FindByID(ctx context.Context, id int) (*domain.Entry, error)
	Create(ctx context.Context, e *domain.Entry) error
	Update(ctx context.Context, e *domain.Entry) error
	Delete(ctx context.Context, id int) error
	StatsForUser(ctx context.Context, userID string, from, to time.Time) ([]domain.EmotionStat, error)
}

type TrackerService struct {
	entries EntryRepo
}

func NewTrackerService(entries EntryRepo) *TrackerService {
	return &TrackerService{entries: entries}
}

func (s *TrackerService) ListEntries(ctx context.Context, userID string, from, to time.Time) ([]domain.Entry, error) {
	return s.entries.FindByUser(ctx, userID, from, to)
}

func (s *TrackerService) CreateEntry(ctx context.Context, userID string, input domain.CreateEntryInput) (*domain.Entry, error) {
	entryDate := input.EntryDate
	if entryDate.IsZero() {
		entryDate = time.Now()
	}
	e := &domain.Entry{
		UserID:    userID,
		EmotionID: input.EmotionID,
		Intensity: input.Intensity,
		Comment:   input.Comment,
		EntryDate: entryDate,
	}
	if err := s.entries.Create(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *TrackerService) UpdateEntry(ctx context.Context, userID string, id int, input domain.UpdateEntryInput) (*domain.Entry, error) {
	e, err := s.entries.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if e.UserID != userID {
		return nil, domain.ErrNotFound
	}
	e.EmotionID = input.EmotionID
	e.Intensity = input.Intensity
	e.Comment = input.Comment
	e.EntryDate = input.EntryDate
	if err := s.entries.Update(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *TrackerService) DeleteEntry(ctx context.Context, userID string, id int) error {
	e, err := s.entries.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.ErrNotFound
		}
		return err
	}
	if e.UserID != userID {
		return domain.ErrNotFound
	}
	return s.entries.Delete(ctx, id)
}

func (s *TrackerService) Stats(ctx context.Context, userID string, period string) ([]domain.EmotionStat, error) {
	from, to := periodRange(period)
	return s.entries.StatsForUser(ctx, userID, from, to)
}

func periodRange(period string) (time.Time, time.Time) {
	now := time.Now()
	// Borne haute = fin de journée : les entrées du jour sont horodatées à midi
	// (T12:00:00Z), donc un simple time.Now() les exclurait avant midi.
	to := endOfDay(now)
	switch period {
	case "week":
		return now.AddDate(0, 0, -7), to
	case "quarter":
		return now.AddDate(0, -3, 0), to
	case "year":
		return now.AddDate(-1, 0, 0), to
	default: // "month" and fallback
		return now.AddDate(0, -1, 0), to
	}
}

// endOfDay renvoie le dernier instant du jour de t, pour que la borne haute d'un
// intervalle inclue toutes les entrées datées de ce jour.
func endOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
}
