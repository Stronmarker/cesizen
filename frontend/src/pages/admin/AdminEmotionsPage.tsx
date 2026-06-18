import { useEffect, useState, type FormEvent } from 'react'
import AdminHeader from '../../components/AdminHeader'
import { emotionApi } from '../../api/emotion'
import type { Emotion, PrimaryEmotion } from '../../types/emotion'

export default function AdminEmotionsPage() {
  const [primaries, setPrimaries] = useState<PrimaryEmotion[]>([])
  const [emotions, setEmotions] = useState<Emotion[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  // Primary form state
  const [editingPrimary, setEditingPrimary] = useState<PrimaryEmotion | null>(null)
  const [creatingPrimary, setCreatingPrimary] = useState(false)
  const [primaryForm, setPrimaryForm] = useState({ label: '', is_active: true })

  // Emotion form state
  const [editingEmotion, setEditingEmotion] = useState<Emotion | null>(null)
  const [creatingEmotion, setCreatingEmotion] = useState(false)
  const [emotionForm, setEmotionForm] = useState({ label: '', primary_emotion_id: 0, is_active: true })

  const [saving, setSaving] = useState(false)

  function load() {
    Promise.all([emotionApi.listPrimary(), emotionApi.list()])
      .then(([p, e]) => { setPrimaries(p); setEmotions(e) })
      .catch(() => setError('Erreur de chargement'))
      .finally(() => setLoading(false))
  }

  useEffect(() => { load() }, [])

  // ── Primary emotion handlers ─────────────────────────────────

  function openCreatePrimary() {
    setEditingPrimary(null)
    setPrimaryForm({ label: '', is_active: true })
    setError('')
    setCreatingPrimary(true)
    setCreatingEmotion(false)
    setEditingEmotion(null)
  }

  function openEditPrimary(p: PrimaryEmotion) {
    setCreatingPrimary(false)
    setPrimaryForm({ label: p.label, is_active: p.is_active })
    setError('')
    setEditingPrimary(p)
  }

  function closePrimaryForm() {
    setCreatingPrimary(false)
    setEditingPrimary(null)
    setError('')
  }

  async function handlePrimarySubmit(e: FormEvent) {
    e.preventDefault()
    setSaving(true)
    setError('')
    try {
      if (editingPrimary) {
        await emotionApi.adminUpdatePrimary(editingPrimary.id, primaryForm.label, primaryForm.is_active)
      } else {
        await emotionApi.adminCreatePrimary(primaryForm.label)
      }
      closePrimaryForm()
      load()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erreur')
    } finally {
      setSaving(false)
    }
  }

  async function handleDeletePrimary(id: number) {
    if (!confirm('Supprimer cette émotion primaire ?')) return
    try {
      await emotionApi.adminDeletePrimary(id)
      load()
    } catch {
      setError('Erreur lors de la suppression')
    }
  }

  // ── Secondary emotion handlers ───────────────────────────────

  function openCreateEmotion() {
    setEditingEmotion(null)
    setEmotionForm({ label: '', primary_emotion_id: primaries[0]?.id ?? 0, is_active: true })
    setError('')
    setCreatingEmotion(true)
    setCreatingPrimary(false)
    setEditingPrimary(null)
  }

  function openEditEmotion(e: Emotion) {
    setCreatingEmotion(false)
    setEmotionForm({ label: e.label, primary_emotion_id: e.primary_emotion_id, is_active: e.is_active })
    setError('')
    setEditingEmotion(e)
  }

  function closeEmotionForm() {
    setCreatingEmotion(false)
    setEditingEmotion(null)
    setError('')
  }

  async function handleEmotionSubmit(e: FormEvent) {
    e.preventDefault()
    setSaving(true)
    setError('')
    try {
      if (editingEmotion) {
        await emotionApi.adminUpdateEmotion(editingEmotion.id, emotionForm.label, emotionForm.primary_emotion_id, emotionForm.is_active)
      } else {
        await emotionApi.adminCreateEmotion(emotionForm.label, emotionForm.primary_emotion_id)
      }
      closeEmotionForm()
      load()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erreur')
    } finally {
      setSaving(false)
    }
  }

  async function handleDeleteEmotion(id: number) {
    if (!confirm('Supprimer cette émotion ?')) return
    try {
      await emotionApi.adminDeleteEmotion(id)
      load()
    } catch {
      setError('Erreur lors de la suppression')
    }
  }

  // Group emotions by primary
  const byPrimary = emotions.reduce<Record<number, Emotion[]>>((acc, e) => {
    if (!acc[e.primary_emotion_id]) acc[e.primary_emotion_id] = []
    acc[e.primary_emotion_id].push(e)
    return acc
  }, {})

  return (
    <div className="min-h-screen bg-transparent">
      <AdminHeader />

      <div className="max-w-2xl mx-auto px-4 py-6">
        {error && !creatingPrimary && !editingPrimary && !creatingEmotion && !editingEmotion && (
          <div className="mb-4 bg-red-50 border border-red-200 rounded-xl p-3 text-red-600 text-sm">{error}</div>
        )}

        {/* ── Primary emotions ─────────────────────── */}
        <div className="flex items-center justify-between mb-3">
          <h2 className="text-green-brand font-black text-lg">Émotions primaires</h2>
          <button onClick={openCreatePrimary} className={btnYellow}>+ Nouveau</button>
        </div>

        {(creatingPrimary || editingPrimary) && (
          <div className="bg-white/55 backdrop-blur-md rounded-2xl shadow-sm border border-white/60 p-5 mb-4">
            <h3 className="text-green-brand font-bold text-base mb-4">
              {editingPrimary ? 'Modifier' : 'Nouvelle émotion primaire'}
            </h3>
            <form onSubmit={handlePrimarySubmit} className="flex flex-col gap-3">
              <input
                type="text" placeholder="Label *" required
                value={primaryForm.label}
                onChange={(e) => setPrimaryForm((f) => ({ ...f, label: e.target.value }))}
                className={inputClass}
              />
              {editingPrimary && (
                <label className="flex items-center gap-2 text-sm text-green-brand/75 cursor-pointer">
                  <input type="checkbox" checked={primaryForm.is_active}
                    onChange={(e) => setPrimaryForm((f) => ({ ...f, is_active: e.target.checked }))}
                    className="w-4 h-4 accent-green-brand" />
                  Active
                </label>
              )}
              {error && <p className="text-red-600 text-sm bg-red-50 border border-red-200 rounded-lg px-3 py-2">{error}</p>}
              <div className="flex gap-2">
                <button type="submit" disabled={saving} className={`flex-1 ${btnYellow} py-3 disabled:opacity-50`}>
                  {saving ? 'Sauvegarde…' : editingPrimary ? 'Mettre à jour' : 'Créer'}
                </button>
                <button type="button" onClick={closePrimaryForm} className="px-4 py-3 rounded-xl border border-green-brand/15 text-green-brand/65 text-sm font-medium">
                  Annuler
                </button>
              </div>
            </form>
          </div>
        )}

        {loading ? (
          <p className="text-center text-green-brand/45 text-sm py-4">Chargement…</p>
        ) : (
          <div className="flex flex-col gap-2 mb-8">
            {primaries.map((p) => (
              <div key={p.id} className="bg-white/55 backdrop-blur-md rounded-2xl shadow-sm border border-white/60 p-4 flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <span className="text-green-brand/80 font-semibold text-sm">{p.label}</span>
                  <span className={`text-xs px-2 py-0.5 rounded-full font-medium ${p.is_active ? 'bg-green-light text-green-brand' : 'bg-green-brand/8 text-green-brand/45'}`}>
                    {p.is_active ? 'Active' : 'Inactive'}
                  </span>
                </div>
                <div className="flex gap-1.5">
                  <button onClick={() => openEditPrimary(p)} className={iconBtn}><PencilIcon /></button>
                  <button onClick={() => handleDeletePrimary(p.id)} className={iconBtnRed}><TrashIcon /></button>
                </div>
              </div>
            ))}
          </div>
        )}

        {/* ── Secondary emotions ─────────────────────── */}
        <div className="flex items-center justify-between mb-3">
          <h2 className="text-green-brand font-black text-lg">Émotions secondaires</h2>
          <button onClick={openCreateEmotion} className={btnYellow}>+ Nouveau</button>
        </div>

        {(creatingEmotion || editingEmotion) && (
          <div className="bg-white/55 backdrop-blur-md rounded-2xl shadow-sm border border-white/60 p-5 mb-4">
            <h3 className="text-green-brand font-bold text-base mb-4">
              {editingEmotion ? 'Modifier' : 'Nouvelle émotion'}
            </h3>
            <form onSubmit={handleEmotionSubmit} className="flex flex-col gap-3">
              <input
                type="text" placeholder="Label *" required
                value={emotionForm.label}
                onChange={(e) => setEmotionForm((f) => ({ ...f, label: e.target.value }))}
                className={inputClass}
              />
              <select
                required
                value={emotionForm.primary_emotion_id}
                onChange={(e) => setEmotionForm((f) => ({ ...f, primary_emotion_id: Number(e.target.value) }))}
                className={inputClass}
              >
                <option value={0} disabled>Émotion primaire *</option>
                {primaries.map((p) => (
                  <option key={p.id} value={p.id}>{p.label}</option>
                ))}
              </select>
              {editingEmotion && (
                <label className="flex items-center gap-2 text-sm text-green-brand/75 cursor-pointer">
                  <input type="checkbox" checked={emotionForm.is_active}
                    onChange={(e) => setEmotionForm((f) => ({ ...f, is_active: e.target.checked }))}
                    className="w-4 h-4 accent-green-brand" />
                  Active
                </label>
              )}
              {error && <p className="text-red-600 text-sm bg-red-50 border border-red-200 rounded-lg px-3 py-2">{error}</p>}
              <div className="flex gap-2">
                <button type="submit" disabled={saving} className={`flex-1 ${btnYellow} py-3 disabled:opacity-50`}>
                  {saving ? 'Sauvegarde…' : editingEmotion ? 'Mettre à jour' : 'Créer'}
                </button>
                <button type="button" onClick={closeEmotionForm} className="px-4 py-3 rounded-xl border border-green-brand/15 text-green-brand/65 text-sm font-medium">
                  Annuler
                </button>
              </div>
            </form>
          </div>
        )}

        {!loading && (
          <div className="flex flex-col gap-3">
            {primaries.map((p) => (
              <div key={p.id}>
                <p className="text-xs font-bold text-green-brand/45 uppercase tracking-wider mb-1.5 px-1">{p.label}</p>
                <div className="flex flex-col gap-1.5">
                  {(byPrimary[p.id] ?? []).map((e) => (
                    <div key={e.id} className="bg-white/55 backdrop-blur-md rounded-xl shadow-sm border border-white/60 px-4 py-2.5 flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <span className="text-green-brand/80 text-sm">{e.label}</span>
                        {!e.is_active && (
                          <span className="text-xs px-2 py-0.5 rounded-full bg-green-brand/8 text-green-brand/45">Inactive</span>
                        )}
                      </div>
                      <div className="flex gap-1.5">
                        <button onClick={() => openEditEmotion(e)} className={iconBtn}><PencilIcon /></button>
                        <button onClick={() => handleDeleteEmotion(e.id)} className={iconBtnRed}><TrashIcon /></button>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
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
  'w-full px-4 py-3 rounded-xl border border-green-brand/15 bg-white/55 text-sm ' +
  'focus:outline-none focus:border-green-brand focus:ring-2 focus:ring-green-brand/20 transition-colors'

const btnYellow = 'bg-yellow-brand text-green-brand text-sm font-bold px-4 py-2 rounded-xl active:scale-[0.98] transition-transform'
const iconBtn = 'p-2 rounded-lg bg-white/55 hover:bg-green-brand/8 text-green-brand/65 transition-colors'
const iconBtnRed = 'p-2 rounded-lg bg-red-50 hover:bg-red-100 text-red-500 transition-colors'
