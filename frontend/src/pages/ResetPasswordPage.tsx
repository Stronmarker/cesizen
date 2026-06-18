import { useState, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { authApi } from '../api/auth'

export default function ResetPasswordPage() {
  const navigate = useNavigate()
  const [form, setForm] = useState({ token: '', password: '', confirm: '' })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    if (form.password !== form.confirm) {
      setError('Les mots de passe ne correspondent pas')
      return
    }
    setLoading(true); setError('')
    try {
      await authApi.resetPassword(form.token, form.password)
      navigate('/login', { state: { message: 'Mot de passe réinitialisé. Connectez-vous.' } })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erreur')
    } finally {
      setLoading(false)
    }
  }

  const set = (key: keyof typeof form) => (e: React.ChangeEvent<HTMLInputElement>) =>
    setForm((f) => ({ ...f, [key]: e.target.value }))

  return (
    <div className="relative h-[100dvh] flex flex-col justify-center px-6 overflow-hidden pt-[env(safe-area-inset-top)] pb-[env(safe-area-inset-bottom)]">
      {/* Halos décoratifs */}
      <div className="absolute -top-16 -right-16 w-60 h-60 rounded-full bg-green-brand/8 pointer-events-none" />
      <div className="absolute bottom-10 -left-12 w-36 h-36 rounded-full bg-yellow-brand/30 pointer-events-none" />

      <div className="relative w-full max-w-sm mx-auto">
        <Link
          to="/forgot-password"
          className="inline-flex items-center gap-1.5 text-green-brand/50 text-sm mb-6 hover:text-green-brand/70 transition-colors"
        >
          <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" d="M15 19l-7-7 7-7" />
          </svg>
          Retour
        </Link>

        <div className="bg-white/55 backdrop-blur-md border border-white/60 rounded-3xl px-6 pt-8 pb-8 shadow-2xl">
          <h1 className="text-2xl font-bold text-green-brand mb-2">Nouveau mot de passe</h1>
          <p className="text-green-brand/45 text-sm mb-6">
            Collez le token reçu et choisissez un nouveau mot de passe.
          </p>

          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            <div className="flex flex-col gap-1.5">
              <label htmlFor="token" className="text-sm font-semibold text-green-brand/75">Token de réinitialisation</label>
              <input
                id="token" type="text" value={form.token}
                onChange={set('token')} required
                placeholder="Collez votre token ici" className={inputClass}
              />
            </div>

            <div className="flex flex-col gap-1.5">
              <label htmlFor="password" className="text-sm font-semibold text-green-brand/75">Nouveau mot de passe</label>
              <input
                id="password" type="password" value={form.password}
                onChange={set('password')} required minLength={6}
                placeholder="••••••••" className={inputClass}
              />
            </div>

            <div className="flex flex-col gap-1.5">
              <label htmlFor="confirm" className="text-sm font-semibold text-green-brand/75">Confirmer</label>
              <input
                id="confirm" type="password" value={form.confirm}
                onChange={set('confirm')} required
                placeholder="••••••••" className={inputClass}
              />
            </div>

            {error && (
              <p className="text-red-600 text-sm bg-red-50 border border-red-200 rounded-lg px-3 py-2">{error}</p>
            )}

            <button
              type="submit" disabled={loading}
              className="mt-2 w-full bg-yellow-brand text-green-brand font-bold text-base py-3.5 rounded-xl
                         active:scale-[0.98] transition-transform disabled:opacity-50"
            >
              {loading ? 'Réinitialisation…' : 'Changer le mot de passe'}
            </button>
          </form>

          <p className="mt-6 text-center text-sm text-green-brand/55">
            <Link to="/login" className="text-green-brand font-semibold underline underline-offset-2">
              Retour à la connexion
            </Link>
          </p>
        </div>
      </div>
    </div>
  )
}

const inputClass =
  'w-full px-4 py-3 rounded-xl border border-green-brand/15 bg-white/55 text-sm ' +
  'focus:outline-none focus:border-green-brand focus:ring-2 focus:ring-green-brand/20 ' +
  'placeholder:text-green-brand/45 transition-colors'
