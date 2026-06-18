export interface PrimaryEmotion {
  id: number
  label: string
  is_active: boolean
}

export interface Emotion {
  id: number
  label: string
  primary_emotion_id: number
  primary_label: string
  is_active: boolean
}
