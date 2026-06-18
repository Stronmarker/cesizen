import { useEffect, useState, type FormEvent } from 'react'
import AdminHeader from '../../components/AdminHeader'
import { userApi } from '../../api/user'
import type { AdminUser } from '../../types/auth'

export default function AdminUsersPage() {
  const [users, setUsers] = useState<AdminUser[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const [editing, setEditing] = useState<AdminUser | null>(null)
  const [form, setForm] = useState({ role: 'user', is_active: true })
  const [saving, setSaving] = useState(false)

  function load() {
    userApi.adminList()
      .then(setUsers)
      .catch(() => setError('Erreur de chargement'))
      .finally(() => setLoading(false))
  }

  useEffect(() => { load() }, [])

  function openEdit(u: AdminUser) {
    setForm({ role: u.role, is_active: u.is_active })
    setError('')
    setEditing(u)
  }

  function closeEdit() {
    setEditing(null)
    setError('')
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    if (!editing) return
    setSaving(true)
    setError('')
    try {
      await userApi.adminUpdate(editing.id, form.role, form.is_active)
      closeEdit()
      load()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Erreur')
    } finally {
      setSaving(false)
    }
  }

  return (
    <div className="min-h-screen bg-transparent">
      <AdminHeader />

      <div className="max-w-2xl mx-auto px-4 py-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-green-brand font-black text-lg">Utilisateurs</h2>
          <span className="text-xs text-green-brand/45">{users.length} compte{users.length !== 1 ? 's' : ''}</span>
        </div>

        {error && !editing && (
          <div className="mb-4 bg-red-50 border border-red-200 rounded-xl p-3 text-red-600 text-sm">{error}</div>
        )}

        {editing && (
          <div className="bg-white/55 backdrop-blur-md rounded-2xl shadow-sm border border-white/60 p-5 mb-4">
            <h3 className="text-green-brand font-bold text-base mb-1">Modifier l'utilisateur</h3>
            <p className="text-green-brand/45 text-xs mb-4">{editing.email}</p>
            <form onSubmit={handleSubmit} className="flex flex-col gap-3">
              <div>
                <label className="text-xs font-medium text-green-brand/55 mb-1 block">Rôle</label>
                <select
                  value={form.role}
                  onChange={(e) => setForm((f) => ({ ...f, role: e.target.value }))}
                  className={inputClass}
                >
                  <option value="user">Utilisateur</option>
                  <option value="admin">Administrateur</option>
                </select>
              </div>
              <label className="flex items-center gap-2 text-sm text-green-brand/75 cursor-pointer">
                <input
                  type="checkbox"
                  checked={form.is_active}
                  onChange={(e) => setForm((f) => ({ ...f, is_active: e.target.checked }))}
                  className="w-4 h-4 accent-green-brand"
                />
                Compte actif
              </label>
              {error && (
                <p className="text-red-600 text-sm bg-red-50 border border-red-200 rounded-lg px-3 py-2">{error}</p>
              )}
              <div className="flex gap-2">
                <button
                  type="submit"
                  disabled={saving}
                  className={`flex-1 ${btnYellow} py-3 disabled:opacity-50`}
                >
                  {saving ? 'Sauvegarde…' : 'Mettre à jour'}
                </button>
                <button
                  type="button"
                  onClick={closeEdit}
                  className="px-4 py-3 rounded-xl border border-green-brand/15 text-green-brand/65 text-sm font-medium"
                >
                  Annuler
                </button>
              </div>
            </form>
          </div>
        )}

        {loading ? (
          <p className="text-center text-green-brand/45 text-sm py-4">Chargement…</p>
        ) : (
          <div className="flex flex-col gap-2">
            {users.map((u) => (
              <div key={u.id} className="bg-white/55 backdrop-blur-md rounded-2xl shadow-sm border border-white/60 p-4">
                <div className="flex items-center justify-between gap-3">
                  <div className="min-w-0">
                    <div className="flex items-center gap-2 flex-wrap">
                      <span className="text-green-brand/80 font-semibold text-sm truncate">
                        {u.first_name || '(sans prénom)'}
                        {u.nickname ? ` · ${u.nickname}` : ''}
                      </span>
                      <span className={`text-xs px-2 py-0.5 rounded-full font-medium shrink-0 ${
                        u.role === 'admin'
                          ? 'bg-yellow-brand/20 text-yellow-600'
                          : 'bg-green-brand/8 text-green-brand/55'
                      }`}>
                        {u.role}
                      </span>
                      {!u.is_active && (
                        <span className="text-xs px-2 py-0.5 rounded-full bg-red-50 text-red-400 font-medium shrink-0">
                          Inactif
                        </span>
                      )}
                    </div>
                    <p className="text-green-brand/45 text-xs mt-0.5 truncate">{u.email}</p>
                    <p className="text-green-brand/30 text-xs">
                      Inscrit le {new Date(u.created_at).toLocaleDateString('fr-FR')}
                    </p>
                  </div>
                  <button
                    onClick={() => openEdit(u)}
                    className="p-2 rounded-lg bg-white/55 hover:bg-green-brand/8 text-green-brand/65 transition-colors shrink-0"
                  >
                    <PencilIcon />
                  </button>
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

const inputClass =
  'w-full px-4 py-3 rounded-xl border border-green-brand/15 bg-white/55 text-sm ' +
  'focus:outline-none focus:border-green-brand focus:ring-2 focus:ring-green-brand/20 transition-colors'

const btnYellow = 'bg-yellow-brand text-green-brand text-sm font-bold px-4 py-2 rounded-xl active:scale-[0.98] transition-transform'
