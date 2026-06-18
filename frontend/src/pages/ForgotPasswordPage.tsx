import { useState, type FormEvent } from 'react'
import { Link } from 'react-router-dom'
import { authApi } from '../api/auth'

export default function ForgotPasswordPage() {
  const [email, setEmail] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [resetToken, setResetToken] = useState('')

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setLoading(true); setError('')
    try {
      const resp = await authApi.forgotPassword(email)
      setResetToken(resp.reset_token)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erreur')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="relative h-[100dvh] flex flex-col justify-center px-6 overflow-hidden pt-[env(safe-area-inset-top)] pb-[env(safe-area-inset-bottom)]">
      {/* Halos décoratifs */}
      <div className="absolute -top-16 -right-16 w-60 h-60 rounded-full bg-green-brand/8 pointer-events-none" />
      <div className="absolute bottom-10 -left-12 w-36 h-36 rounded-full bg-yellow-brand/30 pointer-events-none" />

      <div className="relative w-full max-w-sm mx-auto">
        <Link
          to="/login"
          className="inline-flex items-center gap-1.5 text-green-brand/50 text-sm mb-6 hover:text-green-brand/70 transition-colors"
        >
          <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" d="M15 19l-7-7 7-7" />
          </svg>
          Retour à la connexion
        </Link>

        <div className="bg-white/55 backdrop-blur-md border border-white/60 rounded-3xl px-6 pt-8 pb-8 shadow-2xl">
          <h1 className="text-2xl font-bold text-green-brand mb-2">Mot de passe oublié</h1>
          <p className="text-green-brand/45 text-sm mb-6">
            Entrez votre email pour recevoir un token de réinitialisation.
          </p>

          {resetToken ? (
            <div className="flex flex-col gap-4">
              <div className="bg-green-light/30 border border-green-brand/20 rounded-2xl p-4">
                <p className="text-green-brand font-semibold text-sm mb-2">Token généré :</p>
                <p className="font-mono text-xs text-green-brand/75 break-all bg-white rounded-lg p-3 border border-green-brand/12">
                  {resetToken}
                </p>
                <p className="text-green-brand/45 text-xs mt-2">
                  Copiez ce token et utilisez-le sur la page de réinitialisation. Valide 1 heure.
                </p>
              </div>
              <Link
                to="/reset-password"
                className="w-full bg-yellow-brand text-green-brand font-bold text-sm py-3 rounded-xl text-center block"
              >
                Réinitialiser le mot de passe →
              </Link>
            </div>
          ) : (
            <form onSubmit={handleSubmit} className="flex flex-col gap-4">
              <div className="flex flex-col gap-1.5">
                <label htmlFor="email" className="text-sm font-semibold text-green-brand/75">Email</label>
                <input
                  id="email" type="email" value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required placeholder="you@example.com" className={inputClass}
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
                {loading ? 'Envoi…' : 'Générer un token'}
              </button>
            </form>
          )}
        </div>
      </div>
    </div>
  )
}

const inputClass =
  'w-full px-4 py-3 rounded-xl border border-green-brand/15 bg-white/55 text-sm ' +
  'focus:outline-none focus:border-green-brand focus:ring-2 focus:ring-green-brand/20 ' +
  'placeholder:text-green-brand/45 transition-colors'
