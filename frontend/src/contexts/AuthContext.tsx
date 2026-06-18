import { createContext, useContext, useState, useCallback, type ReactNode } from 'react'
import type { User } from '../types/auth'
import { authApi } from '../api/auth'

interface AuthContextType {
  user: User | null
  token: string | null
  setAuth: (token: string, refreshToken: string, user: User) => void
  logout: () => void
  refreshAccessToken: () => Promise<string | null>
}

const AuthContext = createContext<AuthContextType | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(localStorage.getItem('token'))
  const [user, setUser] = useState<User | null>(() => {
    const stored = localStorage.getItem('user')
    return stored ? JSON.parse(stored) : null
  })

  const setAuth = useCallback((t: string, refreshToken: string, u: User) => {
    localStorage.setItem('token', t)
    localStorage.setItem('refresh_token', refreshToken)
    localStorage.setItem('user', JSON.stringify(u))
    setToken(t)
    setUser(u)
  }, [])

  const logout = useCallback(() => {
    const rt = localStorage.getItem('refresh_token')
    if (rt) authApi.logout(rt).catch(() => {})
    localStorage.removeItem('token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user')
    setToken(null)
    setUser(null)
  }, [])

  const refreshAccessToken = useCallback(async (): Promise<string | null> => {
    const rt = localStorage.getItem('refresh_token')
    if (!rt) return null
    try {
      const resp = await authApi.refresh(rt)
      setAuth(resp.token, resp.refresh_token, resp.user)
      return resp.token
    } catch {
      logout()
      return null
    }
  }, [setAuth, logout])

  return (
    <AuthContext.Provider value={{ user, token, setAuth, logout, refreshAccessToken }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
