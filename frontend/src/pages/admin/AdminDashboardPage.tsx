import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import AdminHeader from '../../components/AdminHeader'
import { userApi } from '../../api/user'
import { contentApi } from '../../api/content'
import { emotionApi } from '../../api/emotion'

type Stats = {
  users: number
  activeUsers: number
  admins: number
  contents: number
  published: number
  primaries: number
  emotions: number
}

export default function AdminDashboardPage() {
  const [stats, setStats] = useState<Stats | null>(null)
  const [error, setError] = useState('')

  useEffect(() => {
    Promise.all([
      userApi.adminList(),
      contentApi.adminList(),
      emotionApi.listPrimary(),
      emotionApi.list(),
    ])
      .then(([users, contents, primaries, emotions]) => {
        setStats({
          users: users.length,
          activeUsers: users.filter((u) => u.is_active).length,
          admins: users.filter((u) => u.role === 'admin').length,
          contents: contents.length,
          published: contents.filter((c) => c.is_published).length,
          primaries: primaries.length,
          emotions: emotions.length,
        })
      })
      .catch(() => setError('Erreur de chargement des statistiques'))
  }, [])

  return (
    <div className="min-h-screen bg-transparent">
      <AdminHeader />

      <div className="max-w-2xl mx-auto px-4 py-6">
        <h2 className="text-green-brand font-black text-lg mb-0.5">Tableau de bord</h2>
        <p className="text-green-brand/45 text-sm mb-5">Vue d'ensemble et gestion de l'application</p>

        {error && (
          <div className="mb-4 bg-red-50 border border-red-200 rounded-xl p-3 text-red-600 text-sm">{error}</div>
        )}

        {/* Indicateurs */}
        <div className="grid grid-cols-2 gap-3 mb-7">
          <StatCard label="Utilisateurs" value={stats?.users}
            hint={stats ? `${stats.activeUsers} actifs · ${stats.admins} admin` : undefined} />
          <StatCard label="Contenus" value={stats?.contents}
            hint={stats ? `${stats.published} publiés` : undefined} />
          <StatCard label="Émotions de base" value={stats?.primaries} hint="niveau 1" />
          <StatCard label="Émotions" value={stats?.emotions} hint="niveau 2" />
        </div>

        {/* Accès gestion */}
        <p className="text-green-brand/55 text-xs font-semibold uppercase tracking-wide mb-2">Gérer</p>
        <div className="flex flex-col gap-3">
          <ManageLink to="/admin/users" emoji="👥" title="Utilisateurs"
            desc="Comptes, rôles, activation / désactivation" />
          <ManageLink to="/admin/contents" emoji="📰" title="Contenus"
            desc="Pages d'information : créer, éditer, publier" />
          <ManageLink to="/admin/emotions" emoji="🎭" title="Référentiel d'émotions"
            desc="Émotions niveau 1 et niveau 2 du tracker" />
        </div>
      </div>
    </div>
  )
}

function StatCard({ label, value, hint }: Readonly<{ label: string; value?: number; hint?: string }>) {
  return (
    <div className="bg-white/55 backdrop-blur-md rounded-2xl shadow-sm border border-white/60 p-4">
      <p className="text-green-brand/55 text-xs font-medium">{label}</p>
      <p className="text-green-brand text-3xl font-black mt-1 leading-none">
        {value ?? '—'}
      </p>
      {hint && <p className="text-green-brand/40 text-xs mt-1.5">{hint}</p>}
    </div>
  )
}

function ManageLink({ to, emoji, title, desc }: Readonly<{ to: string; emoji: string; title: string; desc: string }>) {
  return (
    <Link
      to={to}
      className="bg-white/55 backdrop-blur-md border border-white/60 rounded-2xl p-4 flex items-center gap-4
                 active:scale-[0.99] active:opacity-90 transition-all shadow-sm"
    >
      <div className="w-11 h-11 bg-green-brand/8 rounded-xl flex items-center justify-center text-xl shrink-0">
        {emoji}
      </div>
      <div className="flex-1 min-w-0">
        <p className="font-semibold text-green-brand text-sm">{title}</p>
        <p className="text-green-brand/50 text-xs mt-0.5">{desc}</p>
      </div>
      <svg className="w-4 h-4 text-green-brand/30 shrink-0" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" d="M9 5l7 7-7 7" />
      </svg>
    </Link>
  )
}
