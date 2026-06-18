import { useState, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { authApi } from '../api/auth'
import { useAuth } from '../contexts/AuthContext'
import Logo from '../components/Logo'

export default function RegisterPage() {
  const [form, setForm] = useState({ email: '', password: '', first_name: '', nickname: '' })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const { setAuth } = useAuth()
  const navigate = useNavigate()

  const set = (key: keyof typeof form) => (e: React.ChangeEvent<HTMLInputElement>) =>
    setForm((f) => ({ ...f, [key]: e.target.value }))

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const { token, refresh_token, user } = await authApi.register(
        form.email,
        form.password,
        form.first_name,
        form.nickname || undefined,
      )
      setAuth(token, refresh_token, user)
      navigate('/')
    } catch (err) {
      setError(err instanceof Error ? err.message : "Erreur lors de l'inscription")
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="relative h-[100dvh] flex flex-col items-center justify-center px-6 overflow-hidden pt-[env(safe-area-inset-top)] pb-[env(safe-area-inset-bottom)]">
      {/* Halos décoratifs */}
      <div className="absolute -top-16 -right-16 w-60 h-60 rounded-full bg-green-brand/8 pointer-events-none" />
      <div className="absolute bottom-10 -left-12 w-36 h-36 rounded-full bg-yellow-brand/30 pointer-events-none" />

      {/* Logo centré */}
      <div className="flex flex-col items-center mb-6">
        <Logo size="md" variant="dark" />
        <p className="mt-2 text-green-brand/55 text-sm tracking-wide">
          Créez votre compte
        </p>
      </div>

      {/* Carte blanche centrée */}
      <div className="w-full max-w-sm bg-white/55 backdrop-blur-md border border-white/60 rounded-3xl px-6 pt-7 pb-6 shadow-2xl">
        <h1 className="text-2xl font-bold text-green-brand mb-6">Inscription</h1>

        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <div className="grid grid-cols-2 gap-3">
            <Field label="Prénom *" htmlFor="first_name">
              <input
                id="first_name"
                type="text"
                value={form.first_name}
                onChange={set('first_name')}
                required
                placeholder="Alice"
                className={inputClass}
              />
            </Field>
          </div>

          <Field label="Email *" htmlFor="email">
            <input
              id="email"
              type="email"
              value={form.email}
              onChange={set('email')}
              required
              autoComplete="email"
              placeholder="you@example.com"
              className={inputClass}
            />
          </Field>

          <Field label="Mot de passe *" htmlFor="password">
            <input
              id="password"
              type="password"
              value={form.password}
              onChange={set('password')}
              required
              autoComplete="new-password"
              minLength={8}
              placeholder="8 caractères minimum"
              className={inputClass}
            />
          </Field>

          {error && (
            <p className="text-red-600 text-sm bg-red-50 border border-red-200 rounded-lg px-3 py-2">
              {error}
            </p>
          )}

          <button
            type="submit"
            disabled={loading}
            className="mt-2 w-full bg-yellow-brand text-green-brand font-bold text-base py-3.5 rounded-xl
                       active:scale-[0.98] transition-transform
                       disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? 'Inscription…' : "S'inscrire"}
          </button>
        </form>

        <p className="mt-6 text-center text-sm text-green-brand/55">
          Déjà un compte ?{' '}
          <Link to="/login" className="text-green-brand font-semibold underline underline-offset-2">
            Se connecter
          </Link>
        </p>
      </div>
    </div>
  )
}

function Field({ label, htmlFor, children }: Readonly<{ label: string; htmlFor: string; children: React.ReactNode }>) {
  return (
    <div className="flex flex-col gap-1.5">
      <label htmlFor={htmlFor} className="text-sm font-semibold text-green-brand/75">
        {label}
      </label>
      {children}
    </div>
  )
}

const inputClass =
  'w-full px-4 py-3 rounded-xl border border-green-brand/15 bg-white/55 text-sm ' +
  'focus:outline-none focus:border-green-brand focus:ring-2 focus:ring-green-brand/20 ' +
  'placeholder:text-green-brand/45 transition-colors'
