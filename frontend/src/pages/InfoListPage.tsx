import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import AppShell from '../components/AppShell'
import { contentApi } from '../api/content'
import type { Content } from '../types/content'

export default function InfoListPage() {
  const [contents, setContents] = useState<Content[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    contentApi.list()
      .then(setContents)
      .catch(() => setError('Impossible de charger les articles.'))
      .finally(() => setLoading(false))
  }, [])

  return (
    <AppShell>
      {/* Header */}
      <div className="relative px-6 pt-[calc(2rem_+_env(safe-area-inset-top))] pb-6 overflow-hidden">
        <div className="absolute -top-10 -right-10 w-44 h-44 rounded-full bg-green-brand/5 pointer-events-none" />
        <div className="absolute top-4 right-8 w-16 h-16 rounded-full bg-yellow-brand/10 pointer-events-none" />

        <h1 className="text-green-brand text-2xl font-black">Informations</h1>
        <p className="text-green-brand/45 text-sm mt-0.5">Articles et ressources bien-être</p>
      </div>

      <div className="px-4 pb-4">
        {loading && (
          <div className="bg-white/55 backdrop-blur-md border border-white/60 rounded-2xl p-6 text-center text-green-brand/50 text-sm">
            Chargement…
          </div>
        )}

        {error && (
          <div className="bg-red-50 border border-red-200 rounded-2xl p-4 text-red-600 text-sm">{error}</div>
        )}

        {!loading && !error && contents.length === 0 && (
          <div className="bg-white/55 backdrop-blur-md border border-white/60 rounded-2xl p-8 text-center">
            <p className="text-4xl mb-3">📄</p>
            <p className="text-green-brand/70 text-sm font-medium">Aucun article disponible</p>
          </div>
        )}

        {!loading && !error && contents.length > 0 && (
          <div className="flex flex-col gap-3">
            {contents.map((c) => (
              <Link
                key={c.id}
                to={`/info/${c.id}`}
                className="bg-white/55 backdrop-blur-md border border-white/60 rounded-2xl p-5 block active:scale-[0.99] active:opacity-90 transition-all shadow-sm"
              >
                <h2 className="text-green-brand font-bold text-base leading-snug mb-1">{c.title}</h2>
                {c.author && (
                  <p className="text-green-brand/40 text-xs mb-2">Par {c.author}</p>
                )}
                <p className="text-green-brand/60 text-sm line-clamp-2 leading-relaxed">{c.content}</p>
                <span className="inline-flex items-center gap-1 mt-3 text-xs font-semibold text-green-brand">
                  Lire la suite
                  <svg className="w-3 h-3" fill="none" stroke="currentColor" strokeWidth={2.5} viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" d="M9 5l7 7-7 7" />
                  </svg>
                </span>
              </Link>
            ))}
          </div>
        )}
      </div>
    </AppShell>
  )
}
