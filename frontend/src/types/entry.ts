export interface Entry {
  id: number
  user_id: string
  emotion_id: number
  emotion_label: string
  primary_label: string
  intensity: number
  comment: string
  entry_date: string
  created_at: string
}

export interface EmotionStat {
  emotion_id: number
  emotion_label: string
  primary_label: string
  count: number
}
