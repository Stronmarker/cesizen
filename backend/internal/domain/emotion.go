package domain

type PrimaryEmotion struct {
	ID       int    `json:"id"`
	Label    string `json:"label"`
	IsActive bool   `json:"is_active"`
}

type Emotion struct {
	ID               int    `json:"id"`
	Label            string `json:"label"`
	PrimaryEmotionID int    `json:"primary_emotion_id"`
	PrimaryLabel     string `json:"primary_label"`
	IsActive         bool   `json:"is_active"`
}

type CreatePrimaryEmotionInput struct {
	Label string `json:"label"`
}

type UpdatePrimaryEmotionInput struct {
	Label    string `json:"label"`
	IsActive bool   `json:"is_active"`
}

type CreateEmotionInput struct {
	Label            string `json:"label"`
	PrimaryEmotionID int    `json:"primary_emotion_id"`
}

type UpdateEmotionInput struct {
	Label            string `json:"label"`
	PrimaryEmotionID int    `json:"primary_emotion_id"`
	IsActive         bool   `json:"is_active"`
}
