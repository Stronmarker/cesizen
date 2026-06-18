import type { User, AdminUser } from '../types/auth'
import { authFetch } from './fetchAuth'

export const userApi = {
  getMe: () => authFetch<User>('/users/me', { method: 'GET' }),

  updateMe: (first_name: string, nickname: string) =>
    authFetch<User>('/users/me', {
      method: 'PUT',
      body: JSON.stringify({ first_name, nickname }),
    }),

  deleteMe: () => authFetch<void>('/users/me', { method: 'DELETE' }),

  adminList: () => authFetch<AdminUser[]>('/admin/users', { method: 'GET' }),

  adminUpdate: (id: string, role: string, is_active: boolean) =>
    authFetch<AdminUser>(`/admin/users/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ role, is_active }),
    }),
}
