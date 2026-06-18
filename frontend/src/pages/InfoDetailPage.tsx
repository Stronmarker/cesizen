import { useEffect, useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import AppShell from '../components/AppShell'
import { contentApi } from '../api/content'
import type { Content } from '../types/content'

export default function InfoDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [content, setContent] = useState<Content | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    if (!id) return
    contentApi.getById(Number(id))
      .then(setContent)
      .catch(() => setError('Article introuvable.'))
      .finally(() => setLoading(false))
  }, [id])

  return (
    <AppShell>
      {/* Header */}
      <div className="relative px-6 pt-[calc(2rem_+_env(safe-area-inset-top))] pb-6 overflow-hidden">
        <div className="absolute -top-10 -right-10 w-44 h-44 rounded-full bg-green-brand/5 pointer-events-none" />
        <div className="absolute top-4 right-8 w-16 h-16 rounded-full bg-yellow-brand/10 pointer-events-none" />

        <button
          onClick={() => navigate('/info')}
          className="flex items-center gap-1.5 text-green-brand/45 text-sm mb-4 hover:text-green-brand/65 transition-colors"
        >
          <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" d="M15 19l-7-7 7-7" />
          </svg>
          Retour
        </button>

        {loading && <div className="h-6 w-48 bg-green-brand/15 rounded animate-pulse" />}
        {content && (
          <>
            <h1 className="text-green-brand text-xl font-black leading-snug">{content.title}</h1>
            {content.author && (
              <p className="text-green-brand/45 text-sm mt-1">Par {content.author}</p>
            )}
          </>
        )}
      </div>

      <div className="px-4 pb-4">
        {error && (
          <div className="bg-red-50 border border-red-200 rounded-2xl p-4 text-red-600 text-sm">{error}</div>
        )}

        {content && !loading && (
          <div className="bg-white/55 backdrop-blur-md border border-white/60 rounded-2xl p-5 shadow-sm">
            <p className="text-green-brand/80 text-sm leading-relaxed whitespace-pre-wrap">{content.content}</p>
            <p className="text-green-brand/30 text-xs mt-4 pt-4 border-t border-green-brand/10">
              Publié le {new Date(content.created_at).toLocaleDateString('fr-FR', { day: 'numeric', month: 'long', year: 'numeric' })}
            </p>
          </div>
        )}
      </div>
    </AppShell>
  )
}
