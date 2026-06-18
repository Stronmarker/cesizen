import { useState, type FormEvent } from 'react'
import { Link, useNavigate, useLocation } from 'react-router-dom'
import { authApi } from '../api/auth'
import { useAuth } from '../contexts/AuthContext'
import Logo from '../components/Logo'

export default function LoginPage() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const { setAuth } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()
  const successMessage = (location.state as { message?: string } | null)?.message

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      const { token, refresh_token, user } = await authApi.login(email, password)
      setAuth(token, refresh_token, user)
      navigate('/')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erreur de connexion')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div
      className="fixed inset-0 flex flex-col overflow-hidden"
      style={{ paddingTop: 'env(safe-area-inset-top)', paddingBottom: 'env(safe-area-inset-bottom)' }}
    >
      {/* Halos décoratifs */}
      <div className="absolute -top-16 -right-16 w-60 h-60 rounded-full bg-green-brand/8 pointer-events-none" />
      <div className="absolute top-28 -left-12 w-36 h-36 rounded-full bg-yellow-brand/30 pointer-events-none" />

      {/* Logo centré */}
      <div className="flex flex-col items-center px-6 pt-8 pb-6 shrink-0">
        <Logo size="lg" variant="dark" />
        <p className="mt-2 text-green-brand/55 text-sm tracking-wide">
          Votre bien-être mental au quotidien
        </p>
      </div>

      {/* Zone scrollable interne — la page ne scroll jamais, seule la carte si nécessaire */}
      <div className="flex-1 overflow-y-auto px-6 pb-6">
      <div className="w-full max-w-sm mx-auto bg-white/55 backdrop-blur-md border border-white/60 rounded-3xl px-6 pt-8 pb-8 shadow-2xl">
        <h1 className="text-2xl font-bold text-green-brand mb-6">Connexion</h1>

        {successMessage && (
          <p className="mb-4 text-green-brand text-sm bg-green-light/30 border border-green-brand/20 rounded-lg px-3 py-2">
            {successMessage}
          </p>
        )}

        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <Field label="Email" htmlFor="email">
            <input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              autoComplete="email"
              placeholder="you@example.com"
              className={inputClass}
            />
          </Field>

          <Field label="Mot de passe" htmlFor="password">
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              autoComplete="current-password"
              placeholder="••••••••"
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
            {loading ? 'Connexion…' : 'Se connecter'}
          </button>
        </form>

        <p className="mt-4 text-center text-sm text-green-brand/55">
          <Link to="/forgot-password" className="text-green-brand underline underline-offset-2">
            Mot de passe oublié ?
          </Link>
        </p>

        <p className="mt-3 text-center text-sm text-green-brand/55">
          Pas encore de compte ?{' '}
          <Link to="/register" className="text-green-brand font-semibold underline underline-offset-2">
            S'inscrire
          </Link>
        </p>
      </div>
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
