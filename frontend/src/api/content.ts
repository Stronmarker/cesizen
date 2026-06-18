import type { Content } from '../types/content'

const API_URL = import.meta.env.VITE_API_URL

async function get<T>(path: string): Promise<T> {
  const res = await fetch(`${API_URL}${path}`)
  const data = await res.json()
  if (!res.ok) throw new Error(data.error ?? 'Request failed')
  return data as T
}

async function authRequest<T>(path: string, options: RequestInit): Promise<T> {
  const token = localStorage.getItem('token')
  const res = await fetch(`${API_URL}${path}`, {
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    ...options,
  })
  if (res.status === 204) return undefined as T
  const data = await res.json()
  if (!res.ok) throw new Error(data.error ?? 'Request failed')
  return data as T
}

export const contentApi = {
  list: () => get<Content[]>('/contents'),
  getById: (id: number) => get<Content>(`/contents/${id}`),

  adminList: () => authRequest<Content[]>('/admin/contents', { method: 'GET' }),
  adminCreate: (input: { title: string; content: string; author: string }) =>
    authRequest<Content>('/admin/contents', { method: 'POST', body: JSON.stringify(input) }),
  adminUpdate: (id: number, input: { title: string; content: string; author: string; is_published: boolean }) =>
    authRequest<Content>(`/admin/contents/${id}`, { method: 'PUT', body: JSON.stringify(input) }),
  adminDelete: (id: number) =>
    authRequest<void>(`/admin/contents/${id}`, { method: 'DELETE' }),
}
