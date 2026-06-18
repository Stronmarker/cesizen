import { NavLink, useNavigate } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'

const TABS = [
  { to: '/admin', label: 'Tableau de bord', end: true },
  { to: '/admin/users', label: 'Utilisateurs', end: false },
  { to: '/admin/contents', label: 'Contenus', end: false },
  { to: '/admin/emotions', label: 'Émotions', end: false },
]

/** En-tête commun du back-office : identité, retour à l'app, déconnexion et navigation par onglets. */
export default function AdminHeader() {
  const { user, logout } = useAuth()
  const navigate = useNavigate()

  return (
    <header className="bg-gradient-to-r from-green-deep via-green-brand to-green-mid shadow-md">
      <div className="max-w-2xl mx-auto px-4 pt-3">
        <div className="flex items-center justify-between gap-3">
          <div className="min-w-0">
            <p className="text-white font-bold text-sm">
              Back-office <span className="text-green-light/50 font-normal">· CESIZen</span>
            </p>
            <p className="text-green-light/60 text-xs truncate">{user?.email}</p>
          </div>
          <div className="flex items-center gap-1 shrink-0">
            <button
              onClick={() => navigate('/')}
              className="text-green-light/80 text-sm px-3 py-1.5 rounded-lg hover:bg-white/10 transition-colors"
            >
              ↩ App
            </button>
            <button
              onClick={() => { logout(); navigate('/login') }}
              className="text-green-light/80 text-sm px-3 py-1.5 rounded-lg hover:bg-white/10 transition-colors"
            >
              Déconnexion
            </button>
          </div>
        </div>

        <nav className="flex gap-1 mt-3 overflow-x-auto">
          {TABS.map((t) => (
            <NavLink
              key={t.to}
              to={t.to}
              end={t.end}
              className={({ isActive }) =>
                `px-3.5 py-2 text-sm font-semibold rounded-t-lg whitespace-nowrap transition-colors ${
                  isActive ? 'bg-cream text-green-brand' : 'text-green-light/75 hover:bg-white/10'
                }`
              }
            >
              {t.label}
            </NavLink>
          ))}
        </nav>
      </div>
    </header>
  )
}
