import { FormEvent, useEffect, useState } from 'react'
import { api } from '../api/client'
import type { Profile } from '../types'

const emptyProfile: Profile = {
  user_id: 0,
  first_name: '',
  last_name: '',
  phone: '',
  company: '',
}

export function ProfilePage() {
  const [profile, setProfile] = useState<Profile>(emptyProfile)
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    api.getProfile()
      .then(setProfile)
      .catch((err: Error) => setError(err.message))
      .finally(() => setLoading(false))
  }, [])

  async function handleSubmit(event: FormEvent) {
    event.preventDefault()
    setMessage('')
    setError('')
    setSaving(true)
    try {
      const updated = await api.updateProfile(profile)
      setProfile(updated)
      setMessage('Профиль сохранён')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка сохранения')
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return <div className="page"><p className="muted">Загрузка…</p></div>
  }

  return (
    <div className="page">
      <header className="page-header">
        <h1>Профиль</h1>
      </header>

      {error && <div className="alert alert-error">{error}</div>}
      {message && <div className="alert alert-success">{message}</div>}

      <form className="card form-card" onSubmit={handleSubmit}>
        <div className="grid-2">
          <label>
            Имя
            <input
              value={profile.first_name}
              onChange={(e) => setProfile({ ...profile, first_name: e.target.value })}
            />
          </label>
          <label>
            Фамилия
            <input
              value={profile.last_name}
              onChange={(e) => setProfile({ ...profile, last_name: e.target.value })}
            />
          </label>
        </div>

        <label>
          Телефон
          <input
            value={profile.phone}
            onChange={(e) => setProfile({ ...profile, phone: e.target.value })}
          />
        </label>

        <label>
          Компания
          <input
            value={profile.company}
            onChange={(e) => setProfile({ ...profile, company: e.target.value })}
          />
        </label>

        <button type="submit" className="btn btn-primary" disabled={saving}>
          {saving ? 'Сохраняем…' : 'Сохранить'}
        </button>
      </form>
    </div>
  )
}
