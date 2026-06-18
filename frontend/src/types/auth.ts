export interface User {
  id: string
  email: string
  first_name: string
  nickname: string
  role: string
}

export interface AdminUser {
  id: string
  email: string
  first_name: string
  nickname: string
  role: string
  is_active: boolean
  created_at: string
}

export interface AuthResponse {
  token: string
  refresh_token: string
  user: User
}

export interface ForgotPasswordResponse {
  reset_token: string
  message: string
}
