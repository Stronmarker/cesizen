import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import AppShell from '../components/AppShell'
import { trackerApi } from '../api/tracker'
import type { EmotionStat } from '../types/entry'

type Period = 'week' | 'month' | 'quarter' | 'year'

const PERIODS: { key: Period; label: string }[] = [
  { key: 'week', label: '7j' },
  { key: 'month', label: '1m' },
  { key: 'quarter', label: '3m' },
  { key: 'year', label: '1an' },
]

const EC: Record<string, string> = {
  Joie: 'bg-amber-300', Colère: 'bg-rose-300', Peur: 'bg-violet-300',
  Tristesse: 'bg-sky-300', Surprise: 'bg-orange-300', Dégoût: 'bg-teal-300',
}
function barColor(p: string) { return EC[p] ?? 'bg-green-brand/40' }

export default function StatsPage() {
  const navigate = useNavigate()
  const [period, setPeriod] = useState<Period>('month')
  const [stats, setStats] = useState<EmotionStat[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    setLoading(true); setError('')
    trackerApi.stats(period)
      .then(setStats)
      .catch(() => setError('Erreur de chargement'))
      .finally(() => setLoading(false))
  }, [period])

  const maxCount = stats.length > 0 ? Math.max(...stats.map((s) => s.count)) : 1
  const grouped = stats.reduce<Record<string, EmotionStat[]>>((acc, s) => {
    if (!acc[s.primary_label]) acc[s.primary_label] = []
    acc[s.primary_label].push(s)
    return acc
  }, {})
  const totalEntries = stats.reduce((sum, s) => sum + s.count, 0)

  return (
    <AppShell>
      {/* Header */}
      <div className="relative px-6 pt-[calc(2rem_+_env(safe-area-inset-top))] pb-6 overflow-hidden">
        <div className="absolute -top-10 -right-10 w-44 h-44 rounded-full bg-green-brand/5 pointer-events-none" />
        <div className="absolute top-4 right-16 w-14 h-14 rounded-full bg-yellow-brand/10 pointer-events-none" />

        <button
          onClick={() => navigate('/tracker')}
          className="flex items-center gap-1.5 text-green-brand/45 text-sm mb-4 hover:text-green-brand/65 transition-colors"
        >
          <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" d="M15 19l-7-7 7-7" />
          </svg>
          Retour
        </button>
        <h1 className="text-green-brand text-2xl font-black">Statistiques</h1>
        <p className="text-green-brand/45 text-sm mt-0.5">Analyse de vos émotions</p>
      </div>

      <div className="px-4 pb-4">
        {/* Period selector */}
        <div className="bg-white/55 backdrop-blur-md border border-white/60 rounded-2xl p-1 mb-4 flex gap-1">
          {PERIODS.map((p) => (
            <button
              key={p.key}
              onClick={() => setPeriod(p.key)}
              className={`flex-1 py-2 rounded-xl text-sm font-bold transition-colors ${
                period === p.key ? 'bg-green-brand text-white' : 'text-green-brand/50 hover:text-green-brand/70'
              }`}
            >
              {p.label}
            </button>
          ))}
        </div>

        {loading && (
          <div className="bg-white/55 backdrop-blur-md border border-white/60 rounded-2xl p-6 text-center text-green-brand/50 text-sm">
            Chargement…
          </div>
        )}

        {error && (
          <div className="bg-red-50 border border-red-200 rounded-2xl p-4 text-red-600 text-sm">{error}</div>
        )}

        {!loading && !error && stats.length === 0 && (
          <div className="bg-white/55 backdrop-blur-md border border-white/60 rounded-2xl p-8 text-center">
            <p className="text-4xl mb-3">📊</p>
            <p className="text-green-brand/70 text-sm font-medium">Aucune donnée sur cette période</p>
            <p className="text-green-brand/40 text-xs mt-1">Enregistrez vos émotions pour voir vos stats !</p>
          </div>
        )}

        {!loading && !error && stats.length > 0 && (
          <>
            <div className="bg-white/55 backdrop-blur-md border border-white/60 rounded-2xl p-4 mb-4 flex items-center justify-between shadow-sm">
              <div>
                <p className="text-green-brand/40 text-xs">Total entrées</p>
                <p className="text-2xl font-black text-green-brand">{totalEntries}</p>
              </div>
              <div className="w-px h-8 bg-green-brand/15" />
              <div className="text-right">
                <p className="text-green-brand/40 text-xs">Émotions distinctes</p>
                <p className="text-2xl font-black text-green-brand">{stats.length}</p>
              </div>
            </div>

            {Object.entries(grouped).map(([primary, list]) => (
              <div key={primary} className="bg-white/55 backdrop-blur-md border border-white/60 rounded-2xl p-4 mb-3 shadow-sm">
                <div className="flex items-center gap-2 mb-3">
                  <div className={`w-2.5 h-2.5 rounded-full ${barColor(primary)}`} />
                  <h3 className="text-green-brand font-bold text-sm">{primary}</h3>
                </div>
                <div className="flex flex-col gap-2.5">
                  {list.map((s) => (
                    <div key={s.emotion_id}>
                      <div className="flex justify-between text-xs mb-1">
                        <span className="text-green-brand/70 font-medium">{s.emotion_label}</span>
                        <span className="text-green-brand/40">{s.count} fois</span>
                      </div>
                      <div className="h-2 bg-green-brand/15 rounded-full overflow-hidden">
                        <div
                          className={`h-full ${barColor(primary)} rounded-full transition-all`}
                          style={{ width: `${(s.count / maxCount) * 100}%` }}
                        />
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </>
        )}
      </div>
    </AppShell>
  )
}
