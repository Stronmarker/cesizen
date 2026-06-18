import { useEffect, useState, type FormEvent } from 'react'
import AdminHeader from '../../components/AdminHeader'
import { contentApi } from '../../api/content'
import type { Content } from '../../types/content'

type FormState = { title: string; content: string; author: string; is_published: boolean }

const emptyForm: FormState = { title: '', content: '', author: '', is_published: true }

export default function AdminContentsPage() {
  const [contents, setContents] = useState<Content[]>([])
  const [loading, setLoading] = useState(true)
  const [editing, setEditing] = useState<Content | null>(null)
  const [creating, setCreating] = useState(false)
  const [form, setForm] = useState<FormState>(emptyForm)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')

  function loadContents() {
    contentApi.adminList()
      .then(setContents)
      .catch(() => setError('Erreur de chargement'))
      .finally(() => setLoading(false))
  }

  useEffect(() => { loadContents() }, [])

  function openCreate() {
    setEditing(null)
    setForm(emptyForm)
    setError('')
    setCreating(true)
  }

  function openEdit(c: Content) {
    setCreating(false)
    setForm({ title: c.title, content: c.content, author: c.author, is_published: c.is_published })
    setError('')
    setEditing(c)
  }

  function closeForm() {
    setCreating(false)
    setEditing(null)
    setError('')
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setSaving(true)
    setError('')
    try {
      if (editing) {
        await contentApi.adminUpdate(editing.id, form)
      } else {
        await contentApi.adminCreate({ title: form.title, content: form.content, author: form.author })
      }
      closeForm()
      loadContents()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erreur')
    } finally {
      setSaving(false)
    }
  }

  async function handleDelete(id: number) {
    if (!confirm('Supprimer cet article ?')) return
    try {
      await contentApi.adminDelete(id)
      setContents((prev) => prev.filter((c) => c.id !== id))
    } catch {
      setError('Erreur lors de la suppression')
    }
  }

  const set = (key: keyof FormState) => (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => setForm((f) => ({ ...f, [key]: e.target.value }))

  return (
    <div className="min-h-screen bg-transparent">
      <AdminHeader />

      <div className="max-w-2xl mx-auto px-4 py-6">
        <div className="flex items-center justify-between mb-4">
          <h1 className="text-green-brand font-black text-xl">Articles</h1>
          <button
            onClick={openCreate}
            className="bg-yellow-brand text-green-brand text-sm font-bold px-4 py-2 rounded-xl active:scale-[0.98] transition-transform"
          >
            + Nouveau
          </button>
        </div>

        {error && !creating && !editing && (
          <div className="mb-3 bg-red-50 border border-red-200 rounded-xl p-3 text-red-600 text-sm">{error}</div>
        )}

        {/* Form */}
        {(creating || editing) && (
          <div className="bg-white/55 backdrop-blur-md rounded-2xl shadow-sm border border-white/60 p-5 mb-4">
            <h2 className="text-green-brand font-bold text-base mb-4">
              {editing ? 'Modifier l\'article' : 'Nouvel article'}
            </h2>
            <form onSubmit={handleSubmit} className="flex flex-col gap-3">
              <input
                type="text"
                placeholder="Titre *"
                value={form.title}
                onChange={set('title')}
                required
                className={inputClass}
              />
              <input
                type="text"
                placeholder="Auteur"
                value={form.author}
                onChange={set('author')}
                className={inputClass}
              />
              <textarea
                placeholder="Contenu *"
                value={form.content}
                onChange={set('content')}
                required
                rows={6}
                className={`${inputClass} resize-none`}
              />
              {editing && (
                <label className="flex items-center gap-2 text-sm text-green-brand/75 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={form.is_published}
                    onChange={(e) => setForm((f) => ({ ...f, is_published: e.target.checked }))}
                    className="w-4 h-4 accent-green-brand"
                  />
                  Publié
                </label>
              )}
              {error && (
                <p className="text-red-600 text-sm bg-red-50 border border-red-200 rounded-lg px-3 py-2">{error}</p>
              )}
              <div className="flex gap-2 pt-1">
                <button
                  type="submit"
                  disabled={saving}
                  className="flex-1 bg-yellow-brand text-green-brand font-bold text-sm py-3 rounded-xl disabled:opacity-50 active:scale-[0.98] transition-transform"
                >
                  {saving ? 'Sauvegarde…' : editing ? 'Mettre à jour' : 'Créer'}
                </button>
                <button
                  type="button"
                  onClick={closeForm}
                  className="px-4 py-3 rounded-xl border border-green-brand/15 text-green-brand/65 text-sm font-medium"
                >
                  Annuler
                </button>
              </div>
            </form>
          </div>
        )}

        {/* List */}
        {loading ? (
          <p className="text-center text-green-brand/45 text-sm py-8">Chargement…</p>
        ) : contents.length === 0 ? (
          <p className="text-center text-green-brand/45 text-sm py-8">Aucun article.</p>
        ) : (
          <div className="flex flex-col gap-3">
            {contents.map((c) => (
              <div key={c.id} className="bg-white/55 backdrop-blur-md rounded-2xl shadow-sm border border-white/60 p-4">
                <div className="flex items-start justify-between gap-3">
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-0.5">
                      <h3 className="text-green-brand font-bold text-sm truncate">{c.title}</h3>
                      <span className={`shrink-0 text-xs px-2 py-0.5 rounded-full font-medium ${
                        c.is_published
                          ? 'bg-green-light text-green-brand'
                          : 'bg-green-brand/8 text-green-brand/55'
                      }`}>
                        {c.is_published ? 'Publié' : 'Brouillon'}
                      </span>
                    </div>
                    {c.author && <p className="text-xs text-green-brand/45">Par {c.author}</p>}
                    <p className="text-green-brand/65 text-xs mt-1 line-clamp-1">{c.content}</p>
                  </div>
                  <div className="flex gap-1.5 shrink-0">
                    <button
                      onClick={() => openEdit(c)}
                      className="p-2 rounded-lg bg-white/55 hover:bg-green-brand/8 text-green-brand/65 transition-colors"
                    >
                      <PencilIcon />
                    </button>
                    <button
                      onClick={() => handleDelete(c.id)}
                      className="p-2 rounded-lg bg-red-50 hover:bg-red-100 text-red-500 transition-colors"
                    >
                      <TrashIcon />
                    </button>
                  </div>
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
