import { FormEvent, useEffect, useState } from 'react'
import { api } from '../api/client'
import { useAuth } from '../context/AuthContext'
import type { Application, ApplicationStatus } from '../types'
import { statusClass, statusLabels } from '../utils/labels'

const clientStatuses: ApplicationStatus[] = ['draft', 'submitted']
const managerStatuses: ApplicationStatus[] = ['in_review', 'approved', 'rejected']

export function ApplicationsPage() {
  const { user } = useAuth()
  const [applications, setApplications] = useState<Application[]>([])
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)
  const [showForm, setShowForm] = useState(false)
  const [form, setForm] = useState({ title: '', description: '', status: 'draft' as ApplicationStatus })

  function load() {
    setLoading(true)
    api.listApplications()
      .then(setApplications)
      .catch((err: Error) => setError(err.message))
      .finally(() => setLoading(false))
  }

  useEffect(() => {
    load()
  }, [])

  async function handleCreate(event: FormEvent) {
    event.preventDefault()
    setError('')
    try {
      await api.createApplication(form)
      setForm({ title: '', description: '', status: 'draft' })
      setShowForm(false)
      load()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка создания')
    }
  }

  async function updateStatus(app: Application, status: ApplicationStatus) {
    setError('')
    try {
      await api.updateApplication(app.id, { status })
      load()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка обновления')
    }
  }

  async function handleDelete(id: number) {
    setError('')
    try {
      await api.deleteApplication(id)
      load()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка удаления')
    }
  }

  const canCreate = user?.role === 'client' || user?.role === 'admin'

  return (
    <div className="page">
      <header className="page-header">
        <div>
          <h1>Заявки</h1>
          <p className="muted">Управление заявками</p>
        </div>
        {canCreate && (
          <button type="button" className="btn btn-primary" onClick={() => setShowForm(!showForm)}>
            {showForm ? 'Отмена' : 'Новая заявка'}
          </button>
        )}
      </header>

      {error && <div className="alert alert-error">{error}</div>}

      {showForm && (
        <form className="card form-card" onSubmit={handleCreate}>
          <h2>Новая заявка</h2>
          <label>
            Заголовок
            <input
              value={form.title}
              onChange={(e) => setForm({ ...form, title: e.target.value })}
              required
            />
          </label>
          <label>
            Описание
            <textarea
              rows={4}
              value={form.description}
              onChange={(e) => setForm({ ...form, description: e.target.value })}
            />
          </label>
          <label>
            Статус
            <select
              value={form.status}
              onChange={(e) => setForm({ ...form, status: e.target.value as ApplicationStatus })}
            >
              <option value="draft">Черновик</option>
              <option value="submitted">Отправить</option>
            </select>
          </label>
          <button type="submit" className="btn btn-primary">Создать</button>
        </form>
      )}

      {loading ? (
        <p className="muted">Загрузка…</p>
      ) : applications.length === 0 ? (
        <div className="card empty-state">Заявок пока нет</div>
      ) : (
        <div className="table-card">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>Заголовок</th>
                <th>Статус</th>
                <th>Дата</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              {applications.map((app) => (
                <tr key={app.id}>
                  <td>#{app.id}</td>
                  <td>
                    <strong>{app.title}</strong>
                    <p className="muted">{app.description}</p>
                  </td>
                  <td>
                    <span className={`badge ${statusClass[app.status]}`}>
                      {statusLabels[app.status]}
                    </span>
                  </td>
                  <td>{new Date(app.created_at).toLocaleDateString('ru-RU')}</td>
                  <td className="actions">
                    {user?.role === 'client' && app.status === 'draft' && (
                      <>
                        <button
                          type="button"
                          className="btn btn-small"
                          onClick={() => updateStatus(app, 'submitted')}
                        >
                          Отправить
                        </button>
                        <button
                          type="button"
                          className="btn btn-small btn-danger"
                          onClick={() => handleDelete(app.id)}
                        >
                          Удалить
                        </button>
                      </>
                    )}
                    {user?.role === 'manager' && (
                      <>
                        {app.status === 'submitted' && (
                          <button
                            type="button"
                            className="btn btn-small"
                            onClick={() => updateStatus(app, 'in_review')}
                          >
                            В работу
                          </button>
                        )}
                        {app.status === 'in_review' && managerStatuses.slice(1).map((status) => (
                          <button
                            key={status}
                            type="button"
                            className="btn btn-small"
                            onClick={() => updateStatus(app, status)}
                          >
                            {statusLabels[status]}
                          </button>
                        ))}
                      </>
                    )}
                    {user?.role === 'admin' && clientStatuses
                      .concat(managerStatuses)
                      .filter((status) => status !== app.status)
                      .map((status) => (
                        <button
                          key={status}
                          type="button"
                          className="btn btn-small"
                          onClick={() => updateStatus(app, status)}
                        >
                          → {statusLabels[status]}
                        </button>
                      ))}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
