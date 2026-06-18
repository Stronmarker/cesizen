package domain

import "time"

type Entry struct {
	ID           int       `json:"id"`
	UserID       string    `json:"user_id"`
	EmotionID    int       `json:"emotion_id"`
	EmotionLabel string    `json:"emotion_label,omitempty"`
	PrimaryLabel string    `json:"primary_label,omitempty"`
	Intensity    int       `json:"intensity"`
	Comment      string    `json:"comment"`
	EntryDate    time.Time `json:"entry_date"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateEntryInput struct {
	EmotionID int       `json:"emotion_id"`
	Intensity int       `json:"intensity"`
	Comment   string    `json:"comment"`
	EntryDate time.Time `json:"entry_date"`
}

type UpdateEntryInput struct {
	EmotionID int       `json:"emotion_id"`
	Intensity int       `json:"intensity"`
	Comment   string    `json:"comment"`
	EntryDate time.Time `json:"entry_date"`
}

type EmotionStat struct {
	EmotionID    int    `json:"emotion_id"`
	EmotionLabel string `json:"emotion_label"`
	PrimaryLabel string `json:"primary_label"`
	Count        int    `json:"count"`
}
