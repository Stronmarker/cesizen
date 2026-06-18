import type { Entry, EmotionStat } from '../types/entry'
import { authFetch } from './fetchAuth'

export const trackerApi = {
  list: () => authFetch<Entry[]>('/tracker/entries', { method: 'GET' }),
  create: (emotion_id: number, intensity: number, comment: string, entry_date: string) =>
    authFetch<Entry>('/tracker/entries', {
      method: 'POST',
      body: JSON.stringify({ emotion_id, intensity, comment, entry_date }),
    }),
  update: (id: number, emotion_id: number, intensity: number, comment: string, entry_date: string) =>
    authFetch<Entry>(`/tracker/entries/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ emotion_id, intensity, comment, entry_date }),
    }),
  delete: (id: number) =>
    authFetch<void>(`/tracker/entries/${id}`, { method: 'DELETE' }),
  stats: (period: string) =>
    authFetch<EmotionStat[]>(`/tracker/stats?period=${period}`, { method: 'GET' }),
}
