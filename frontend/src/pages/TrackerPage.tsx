import { useEffect, useState, type FormEvent } from 'react'
import { Link } from 'react-router-dom'
import AppShell from '../components/AppShell'
import { trackerApi } from '../api/tracker'
import { emotionApi } from '../api/emotion'
import type { Entry } from '../types/entry'
import type { Emotion } from '../types/emotion'

type Form = { emotionId: number; intensity: number; comment: string; entryDate: string }

const today = () => new Date().toISOString().slice(0, 10)
const emptyForm = (): Form => ({ emotionId: 0, intensity: 5, comment: '', entryDate: today() })

// Teintes douces et translucides — un repère par émotion, harmonisé au reste
const EC: Record<string, { border: string; pill: string; bar: string }> = {
  Joie:      { border: 'border-l-amber-300',  pill: 'bg-amber-100/70 text-amber-700',   bar: 'bg-amber-300' },
  Colère:    { border: 'border-l-rose-300',   pill: 'bg-rose-100/70 text-rose-600',     bar: 'bg-rose-300' },
  Peur:      { border: 'border-l-violet-300', pill: 'bg-violet-100/70 text-violet-600', bar: 'bg-violet-300' },
  Tristesse: { border: 'border-l-sky-300',    pill: 'bg-sky-100/70 text-sky-600',       bar: 'bg-sky-300' },
  Surprise:  { border: 'border-l-orange-300', pill: 'bg-orange-100/70 text-orange-600', bar: 'bg-orange-300' },
  Dégoût:    { border: 'border-l-teal-300',   pill: 'bg-teal-100/70 text-teal-600',     bar: 'bg-teal-300' },
}
const DEC = { border: 'border-l-green-brand/30', pill: 'bg-green-brand/10 text-green-brand/60', bar: 'bg-green-brand/40' }
function ec(p: string) { return EC[p] ?? DEC }

export default function TrackerPage() {
  const [entries, setEntries] = useState<Entry[]>([])
  const [emotions, setEmotions] = useState<Emotion[]>([])
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [editing, setEditing] = useState<Entry | null>(null)
  const [form, setForm] = useState<Form>(emptyForm())
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')

  function load() {
    Promise.all([trackerApi.list(), emotionApi.list()])
      .then(([e, em]) => { setEntries(e); setEmotions(em) })
      .catch(() => setError('Erreur de chargement'))
      .finally(() => setLoading(false))
  }

  useEffect(() => { load() }, [])

  function openCreate() { setEditing(null); setForm(emptyForm()); setError(''); setShowForm(true) }
  function openEdit(e: Entry) {
    setEditing(e)
    setForm({ emotionId: e.emotion_id, intensity: e.intensity, comment: e.comment, entryDate: e.entry_date.slice(0, 10) })
    setError('')
    setShowForm(true)
  }
  function closeForm() { setShowForm(false); setEditing(null); setError('') }

  async function handleSubmit(ev: FormEvent) {
    ev.preventDefault()
    setSaving(true); setError('')
    try {
      if (editing) {
        await trackerApi.update(editing.id, form.emotionId, form.intensity, form.comment, form.entryDate + 'T12:00:00Z')
      } else {
        await trackerApi.create(form.emotionId, form.intensity, form.comment, form.entryDate + 'T12:00:00Z')
      }
      closeForm(); load()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erreur')
    } finally {
      setSaving(false)
    }
  }

  async function handleDelete(id: number) {
    if (!confirm('Supprimer cette entrée ?')) return
    try {
      await trackerApi.delete(id)
      setEntries((prev) => prev.filter((e) => e.id !== id))
    } catch {
      setError('Erreur lors de la suppression')
    }
  }

  const primaryGroups = emotions.reduce<Record<string, Emotion[]>>((acc, e) => {
    if (!acc[e.primary_label]) acc[e.primary_label] = []
    acc[e.primary_label].push(e)
    return acc
  }, {})

  let submitLabel = editing ? 'Mettre à jour' : 'Enregistrer'
  if (saving) submitLabel = 'Sauvegarde…'

  return (
    <AppShell>
      {/* Header */}
      <div className="relative px-6 pt-[calc(2rem_+_env(safe-area-inset-top))] pb-6 overflow-hidden">
        <div className="absolute -top-10 -right-10 w-44 h-44 rounded-full bg-green-brand/5 pointer-events-none" />
        <div className="absolute top-4 right-16 w-14 h-14 rounded-full bg-yellow-brand/10 pointer-events-none" />

        <div className="flex items-end justify-between">
          <div>
            <h1 className="text-green-brand text-2xl font-black">Mes émotions</h1>
            <p className="text-green-brand/45 text-sm mt-0.5">Journal émotionnel</p>
          </div>
          <Link
            to="/tracker/stats"
            className="text-green-brand text-xs font-semibold bg-green-brand/10 border border-green-brand/20 rounded-xl px-3.5 py-2"
          >
            Stats →
          </Link>
        </div>
      </div>

      <div className="px-4 pb-4">
        {/* Add button */}
        {!showForm && (
          <button
            onClick={openCreate}
            className="w-full bg-yellow-brand text-green-brand font-bold text-sm py-3 rounded-2xl mb-4 active:scale-[0.98] transition-transform"
          >
            + Nouvelle entrée
          </button>
        )}

        {error && !showForm && (
          <div className="mb-3 bg-red-50 border border-red-200 rounded-xl p-3 text-red-600 text-sm">{error}</div>
        )}

        {/* Form */}
        {showForm && (
          <div className="bg-white/55 backdrop-blur-md border border-white/60 rounded-2xl p-5 mb-4 shadow-sm">
            <h2 className="text-green-brand font-bold text-base mb-4">
              {editing ? "Modifier l'entrée" : 'Nouvelle entrée'}
            </h2>
            <form onSubmit={handleSubmit} className="flex flex-col gap-3">
              <div>
                <label htmlFor="form-emotion" className="text-green-brand/70 text-xs font-semibold mb-1.5 block">Émotion *</label>
                <select
                  id="form-emotion"
                  required
                  value={form.emotionId}
                  onChange={(e) => setForm((f) => ({ ...f, emotionId: Number(e.target.value) }))}
                  className={inputClass}
                >
                  <option value={0} disabled>Choisir une émotion…</option>
                  {Object.entries(primaryGroups).map(([primary, list]) => (
                    <optgroup key={primary} label={primary}>
                      {list.map((em) => (
                        <option key={em.id} value={em.id}>{em.label}</option>
                      ))}
                    </optgroup>
                  ))}
                </select>
              </div>

              <div>
                <label htmlFor="form-intensity" className="text-green-brand/70 text-xs font-semibold mb-1.5 block">
                  Intensité : <span className="text-green-brand font-bold">{form.intensity}/10</span>
                </label>
                <input
                  id="form-intensity"
                  type="range" min={1} max={10} value={form.intensity}
                  onChange={(e) => setForm((f) => ({ ...f, intensity: Number(e.target.value) }))}
                  className="w-full accent-yellow-brand"
                />
              </div>

              <div>
                <label htmlFor="form-date" className="text-green-brand/70 text-xs font-semibold mb-1.5 block">Date</label>
                <input
                  id="form-date"
                  type="date" value={form.entryDate} max={today()}
                  onChange={(e) => setForm((f) => ({ ...f, entryDate: e.target.value }))}
                  className={inputClass}
                />
              </div>

              <textarea
                placeholder="Commentaire (optionnel)"
                value={form.comment}
                onChange={(e) => setForm((f) => ({ ...f, comment: e.target.value }))}
                rows={3} className={`${inputClass} resize-none`}
              />

              {error && (
                <p className="text-red-600 text-sm bg-red-50 border border-red-200 rounded-lg px-3 py-2">{error}</p>
              )}

              <div className="flex gap-2 pt-1">
                <button
                  type="submit"
                  disabled={saving || form.emotionId === 0}
                  className="flex-1 bg-yellow-brand text-green-brand font-bold text-sm py-3 rounded-xl disabled:opacity-50 active:scale-[0.98] transition-transform"
                >
                  {submitLabel}
                </button>
                <button
                  type="button" onClick={closeForm}
                  className="px-4 py-3 rounded-xl bg-white/55 backdrop-blur-md border border-white/60 text-green-brand/70 text-sm font-medium"
                >
                  Annuler
                </button>
              </div>
            </form>
          </div>
        )}

        {/* Journal */}
        {loading && (
          <div className="bg-white/55 backdrop-blur-md border border-white/60 rounded-2xl p-6 text-center text-green-brand/50 text-sm">
            Chargement…
          </div>
        )}
        {!loading && entries.length === 0 ? (
          <div className="bg-white/55 backdrop-blur-md border border-white/60 rounded-2xl p-8 text-center">
            <p className="text-4xl mb-3">💭</p>
            <p className="text-green-brand/70 text-sm font-medium">Aucune entrée pour le moment</p>
            <p className="text-green-brand/40 text-xs mt-1">Commencez à tracker vos émotions !</p>
          </div>
        ) : null}
        {!loading && entries.length > 0 && (
          <div className="flex flex-col gap-2.5">
            {entries.map((e) => {
              const c = ec(e.primary_label)
              return (
                <div key={e.id} className={`bg-white/55 backdrop-blur-md border border-white/60 rounded-2xl border-l-4 ${c.border} p-4 shadow-sm`}>
                  <div className="flex items-start justify-between gap-3">
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2 mb-1">
                        <span className={`text-xs font-semibold ${c.pill} px-2 py-0.5 rounded-full shrink-0`}>
                          {e.primary_label}
                        </span>
                        <span className="text-sm font-bold text-green-brand truncate">{e.emotion_label}</span>
                      </div>
                      <div className="flex items-center gap-3">
                        <IntensityBar value={e.intensity} barClass={c.bar} />
                        <span className="text-xs text-green-brand/40">
                          {new Date(e.entry_date).toLocaleDateString('fr-FR', { day: 'numeric', month: 'short' })}
                        </span>
                      </div>
                      {e.comment && (
                        <p className="text-green-brand/50 text-xs mt-1.5 line-clamp-1">{e.comment}</p>
                      )}
                    </div>
                    <div className="flex gap-1.5 shrink-0">
                      <button
                        onClick={() => openEdit(e)}
                        className="p-2 rounded-lg bg-green-brand/10 hover:bg-green-brand/20 text-green-brand/50 transition-colors"
                      >
                        <PencilIcon />
                      </button>
                      <button
                        onClick={() => handleDelete(e.id)}
                        className="p-2 rounded-lg bg-red-100 hover:bg-red-200 text-red-500 transition-colors"
                      >
                        <TrashIcon />
                      </button>
                    </div>
                  </div>
                </div>
              )
            })}
          </div>
        )}
      </div>
    </AppShell>
  )
}

function IntensityBar({ value, barClass }: Readonly<{ value: number; barClass: string }>) {
  return (
    <div className="flex items-center gap-1.5">
      <div className="w-20 h-1.5 bg-green-brand/15 rounded-full overflow-hidden">
        <div className={`h-full ${barClass} rounded-full`} style={{ width: `${(value / 10) * 100}%` }} />
      </div>
      <span className="text-xs text-green-brand/50 font-medium">{value}/10</span>
    </div>
  )
}

function PencilIcon() {
  return (
    <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
    </svg>
  )
}

function TrashIcon() {
  return (
    <svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth={2} viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
    </svg>
  )
}

const inputClass =
  'w-full px-4 py-3 rounded-xl bg-white/70 border border-green-brand/20 text-green-brand text-sm ' +
  'placeholder:text-green-brand/30 focus:outline-none focus:border-green-brand/40 focus:ring-2 focus:ring-green-brand/10 transition-colors'
