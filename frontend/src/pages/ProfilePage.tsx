import { useState, useEffect, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { Capacitor } from '@capacitor/core'
import { useAuth } from '../contexts/AuthContext'
import { userApi } from '../api/user'
import AppShell from '../components/AppShell'

export default function ProfilePage() {
  const { user, setAuth, token, logout } = useAuth()
  const navigate = useNavigate()
  const [form, setForm] = useState({ first_name: user?.first_name ?? '', nickname: user?.nickname ?? '' })
  const [status, setStatus] = useState<'idle' | 'saving' | 'saved' | 'error'>('idle')
  const [error, setError] = useState('')
  const [deleting, setDeleting] = useState(false)

  useEffect(() => {
    if (user) setForm({ first_name: user.first_name, nickname: user.nickname })
  }, [user])

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setStatus('saving'); setError('')
    try {
      const updated = await userApi.updateMe(form.first_name, form.nickname)
      const rt = localStorage.getItem('refresh_token') ?? ''
      setAuth(token!, rt, updated)
      setStatus('saved')
      setTimeout(() => setStatus('idle'), 2000)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erreur lors de la sauvegarde')
      setStatus('error')
    }
  }

  async function handleDeleteAccount() {
    if (!confirm('Supprimer définitivement votre compte ? Cette action est irréversible.')) return
    setDeleting(true)
    try {
      await userApi.deleteMe()
      logout(); navigate('/login')
    } catch {
      setError('Erreur lors de la suppression du compte')
      setDeleting(false)
    }
  }

  function handleLogout() { logout(); navigate('/login') }

  const set = (key: keyof typeof form) => (e: React.ChangeEvent<HTMLInputElement>) =>
    setForm((f) => ({ ...f, [key]: e.target.value }))

  const initial = user?.first_name?.[0]?.toUpperCase() ?? '?'

  let saveLabel = 'Enregistrer'
  if (status === 'saving') saveLabel = 'Sauvegarde…'
  if (status === 'saved') saveLabel = '✓ Sauvegardé'

  return (
    <AppShell>
      {/* Hero */}
      <div className="relative px-6 pt-[calc(2.5rem_+_env(safe-area-inset-top))] pb-8 overflow-hidden">
        <div className="absolute -top-12 -right-12 w-52 h-52 rounded-full bg-green-brand/5 pointer-events-none" />
        <div className="absolute top-8 right-8 w-20 h-20 rounded-full bg-yellow-brand/10 pointer-events-none" />

        <div className="flex flex-col items-center relative z-10">
          <div className="w-20 h-20 rounded-full bg-green-brand flex items-center justify-center mb-3 shadow-lg">
            <span className="text-3xl font-black text-yellow-brand">{initial}</span>
          </div>
          <h1 className="text-green-brand text-xl font-bold">{user?.first_name}</h1>
          <p className="text-green-brand/45 text-sm">{user?.email}</p>
          <span className="mt-2 px-3 py-0.5 rounded-full bg-green-brand/10 text-green-brand text-xs font-semibold uppercase tracking-wide">
            {user?.role}
          </span>
          <button
            onClick={handleLogout}
            className="mt-4 text-green-brand/45 text-xs font-medium hover:text-green-brand/65 transition-colors"
          >
            Se déconnecter
          </button>
        </div>
      </div>

      <div className="px-4">
        {/* Accès back-office — admins, web uniquement */}
        {user?.role === 'admin' && !Capacitor.isNativePlatform() && (
          <Link
            to="/admin"
            className="bg-green-brand text-white rounded-2xl p-4 mb-4 flex items-center gap-4 shadow-sm
                       active:scale-[0.99] transition-transform"
          >
            <div className="w-11 h-11 bg-white/15 rounded-xl flex items-center justify-center text-xl shrink-0">
              🛠️
            </div>
            <div className="flex-1 min-w-0">
              <p className="font-bold text-sm">Espace administration</p>
              <p className="text-green-light/70 text-xs mt-0.5">Gérer l'application (back-office)</p>
            </div>
            <svg className="w-4 h-4 text-white/60 shrink-0" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" d="M9 5l7 7-7 7" />
            </svg>
          </Link>
        )}

        {/* Form */}
        <div className="bg-white/55 backdrop-blur-md border border-white/60 rounded-2xl p-5 mb-4 shadow-sm">
          <h2 className="text-green-brand font-bold text-base mb-4">Mes informations</h2>

          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            <Field label="Prénom *" htmlFor="first_name">
              <input
                id="first_name" type="text" value={form.first_name}
                onChange={set('first_name')} required className={inputClass}
              />
            </Field>

            <Field label="Email" htmlFor="email">
              <input
                id="email" type="email" value={user?.email ?? ''}
                disabled className={`${inputClass} opacity-40 cursor-not-allowed`}
              />
              <p className="text-green-brand/30 text-xs mt-1">L'email ne peut pas être modifié.</p>
            </Field>

            {error && (
              <p className="text-red-600 text-sm bg-red-50 border border-red-200 rounded-lg px-3 py-2">{error}</p>
            )}

            <button
              type="submit" disabled={status === 'saving'}
              className="w-full bg-yellow-brand text-green-brand font-bold text-sm py-3 rounded-xl
                         active:scale-[0.98] transition-transform disabled:opacity-50"
            >
              {saveLabel}
            </button>
          </form>
        </div>

        {/* Danger zone */}
        <div className="mb-4 bg-red-50 rounded-2xl border border-red-100 p-5">
          <h2 className="text-red-600 font-bold text-base mb-1">Supprimer mon compte</h2>
          <p className="text-green-brand/45 text-xs mb-4">Cette action est définitive et supprime toutes vos données.</p>
          <button
            type="button" onClick={handleDeleteAccount} disabled={deleting}
            className="w-full bg-red-100 text-red-600 font-bold text-sm py-3 rounded-xl border border-red-200
                       hover:bg-red-200 active:scale-[0.98] transition-all disabled:opacity-50"
          >
            {deleting ? 'Suppression…' : 'Supprimer mon compte'}
          </button>
        </div>
      </div>
    </AppShell>
  )
}

function Field({ label, htmlFor, children }: Readonly<{ label: string; htmlFor: string; children: React.ReactNode }>) {
  return (
    <div className="flex flex-col gap-1.5">
      <label htmlFor={htmlFor} className="text-green-brand/70 text-sm font-semibold">{label}</label>
      {children}
    </div>
  )
}

const inputClass =
  'w-full px-4 py-3 rounded-xl bg-white/70 border border-green-brand/20 text-green-brand text-sm ' +
  'placeholder:text-green-brand/30 focus:outline-none focus:border-green-brand/40 focus:ring-2 focus:ring-green-brand/10 transition-colors'
