import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { Capacitor } from '@capacitor/core'
import { AuthProvider, useAuth } from './contexts/AuthContext'
import type { ReactNode } from 'react'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'
import ForgotPasswordPage from './pages/ForgotPasswordPage'
import ResetPasswordPage from './pages/ResetPasswordPage'
import HomePage from './pages/HomePage'
import ProfilePage from './pages/ProfilePage'
import InfoListPage from './pages/InfoListPage'
import InfoDetailPage from './pages/InfoDetailPage'
import AdminDashboardPage from './pages/admin/AdminDashboardPage'
import AdminContentsPage from './pages/admin/AdminContentsPage'
import AdminEmotionsPage from './pages/admin/AdminEmotionsPage'
import AdminUsersPage from './pages/admin/AdminUsersPage'
import TrackerPage from './pages/TrackerPage'
import StatsPage from './pages/StatsPage'

function ProtectedRoute({ children }: Readonly<{ children: ReactNode }>) {
  const { token } = useAuth()
  return token ? <>{children}</> : <Navigate to="/login" replace />
}

function AdminRoute({ children }: Readonly<{ children: ReactNode }>) {
  const { token, user } = useAuth()
  // Back-office réservé au web : indisponible sur l'app mobile native (Capacitor).
  if (Capacitor.isNativePlatform()) return <Navigate to="/" replace />
  if (!token) return <Navigate to="/login" replace />
  if (user?.role !== 'admin') return <Navigate to="/" replace />
  return <>{children}</>
}

function AppRoutes() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/register" element={<RegisterPage />} />
      <Route path="/forgot-password" element={<ForgotPasswordPage />} />
      <Route path="/reset-password" element={<ResetPasswordPage />} />
      <Route path="/" element={<ProtectedRoute><HomePage /></ProtectedRoute>} />
      <Route path="/profile" element={<ProtectedRoute><ProfilePage /></ProtectedRoute>} />
      <Route path="/info" element={<ProtectedRoute><InfoListPage /></ProtectedRoute>} />
      <Route path="/info/:id" element={<ProtectedRoute><InfoDetailPage /></ProtectedRoute>} />
      <Route path="/tracker" element={<ProtectedRoute><TrackerPage /></ProtectedRoute>} />
      <Route path="/tracker/stats" element={<ProtectedRoute><StatsPage /></ProtectedRoute>} />
      <Route path="/admin" element={<AdminRoute><AdminDashboardPage /></AdminRoute>} />
      <Route path="/admin/contents" element={<AdminRoute><AdminContentsPage /></AdminRoute>} />
      <Route path="/admin/emotions" element={<AdminRoute><AdminEmotionsPage /></AdminRoute>} />
      <Route path="/admin/users" element={<AdminRoute><AdminUsersPage /></AdminRoute>} />
    </Routes>
  )
}

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <AppRoutes />
      </BrowserRouter>
    </AuthProvider>
  )
}
