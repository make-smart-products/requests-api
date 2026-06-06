import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { api } from '../api/client'
import { useAuth } from '../context/AuthContext'
import type { Application, Notification } from '../types'
import { roleLabels, statusLabels } from '../utils/labels'

export function DashboardPage() {
  const { user } = useAuth()
  const [applications, setApplications] = useState<Application[]>([])
  const [notifications, setNotifications] = useState<Notification[]>([])
  const [error, setError] = useState('')

  useEffect(() => {
    Promise.all([
      api.listApplications(),
      api.listNotifications(true),
    ])
      .then(([apps, notes]) => {
        setApplications(apps.slice(0, 5))
        setNotifications(notes.slice(0, 5))
      })
      .catch((err: Error) => setError(err.message))
  }, [])

  return (
    <div className="page">
      <header className="page-header">
        <div>
          <h1>Обзор</h1>
          <p className="muted">
            {user ? roleLabels[user.role] : ''} · {user?.email}
          </p>
        </div>
      </header>

      {error && <div className="alert alert-error">{error}</div>}

      <div className="stats-grid">
        <div className="stat-card">
          <span>Заявок</span>
          <strong>{applications.length}</strong>
        </div>
        <div className="stat-card">
          <span>Непрочитанных</span>
          <strong>{notifications.length}</strong>
        </div>
      </div>

      <div className="grid-2 stretch">
        <section className="card">
          <div className="card-header">
            <h2>Последние заявки</h2>
            <Link to="/applications">Все</Link>
          </div>
          {applications.length === 0 ? (
            <p className="muted">Заявок пока нет</p>
          ) : (
            <ul className="list">
              {applications.map((app) => (
                <li key={app.id}>
                  <div>
                    <strong>{app.title}</strong>
                    <span className="muted">#{app.id}</span>
                  </div>
                  <span>{statusLabels[app.status]}</span>
                </li>
              ))}
            </ul>
          )}
        </section>

        <section className="card">
          <div className="card-header">
            <h2>Уведомления</h2>
            <Link to="/notifications">Все</Link>
          </div>
          {notifications.length === 0 ? (
            <p className="muted">Новых уведомлений нет</p>
          ) : (
            <ul className="list">
              {notifications.map((note) => (
                <li key={note.id}>
                  <div>
                    <strong>{note.title}</strong>
                    <p className="muted">{note.body}</p>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </section>
      </div>
    </div>
  )
}
