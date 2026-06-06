import { useEffect, useState } from 'react'
import { api } from '../api/client'
import type { Notification } from '../types'

export function NotificationsPage() {
  const [notifications, setNotifications] = useState<Notification[]>([])
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)

  function load() {
    setLoading(true)
    api.listNotifications()
      .then(setNotifications)
      .catch((err: Error) => setError(err.message))
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    load()
  }, [])

  async function markRead(id: number) {
    try {
      await api.markNotificationRead(id)
      load()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка')
    }
  }

  return (
    <div className="page">
      <header className="page-header">
        <h1>Уведомления</h1>
      </header>

      {error && <div className="alert alert-error">{error}</div>}

      {loading ? (
        <p className="muted">Загрузка…</p>
      ) : notifications.length === 0 ? (
        <div className="card empty-state">Уведомлений нет</div>
      ) : (
        <div className="notifications-list">
          {notifications.map((note) => (
            <article key={note.id} className={`card notification ${note.is_read ? 'read' : ''}`}>
              <div className="notification-header">
                <div>
                  <span className="badge badge-neutral">{note.channel}</span>
                  <h3>{note.title}</h3>
                </div>
                {!note.is_read && (
                  <button
                    type="button"
                    className="btn btn-small"
                    onClick={() => markRead(note.id)}
                  >
                    Прочитано
                  </button>
                )}
              </div>
              <p>{note.body}</p>
              <time className="muted">
                {new Date(note.created_at).toLocaleString('ru-RU')}
              </time>
            </article>
          ))}
        </div>
      )}
    </div>
  )
}
