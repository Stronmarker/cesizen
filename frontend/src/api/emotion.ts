import type { Emotion, PrimaryEmotion } from '../types/emotion'
import { authFetch } from './fetchAuth'

const API_URL = import.meta.env.VITE_API_URL

async function get<T>(path: string): Promise<T> {
  const res = await fetch(`${API_URL}${path}`)
  const data = await res.json()
  if (!res.ok) throw new Error(data.error ?? 'Request failed')
  return data as T
}

export const emotionApi = {
  listPrimary: () => get<PrimaryEmotion[]>('/primary-emotions'),
  list: () => get<Emotion[]>('/emotions'),

  adminCreatePrimary: (label: string) =>
    authFetch<PrimaryEmotion>('/admin/primary-emotions', {
      method: 'POST',
      body: JSON.stringify({ label }),
    }),
  adminUpdatePrimary: (id: number, label: string, is_active: boolean) =>
    authFetch<PrimaryEmotion>(`/admin/primary-emotions/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ label, is_active }),
    }),
  adminDeletePrimary: (id: number) =>
    authFetch<void>(`/admin/primary-emotions/${id}`, { method: 'DELETE' }),

  adminCreateEmotion: (label: string, primary_emotion_id: number) =>
    authFetch<Emotion>('/admin/emotions', {
      method: 'POST',
      body: JSON.stringify({ label, primary_emotion_id }),
    }),
  adminUpdateEmotion: (id: number, label: string, primary_emotion_id: number, is_active: boolean) =>
    authFetch<Emotion>(`/admin/emotions/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ label, primary_emotion_id, is_active }),
    }),
  adminDeleteEmotion: (id: number) =>
    authFetch<void>(`/admin/emotions/${id}`, { method: 'DELETE' }),
}
