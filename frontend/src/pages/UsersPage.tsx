import { useEffect, useMemo, useState } from 'react'
import { api } from '../api/client'
import { FiltersBar } from '../components/FiltersBar'
import type { Role, User } from '../types'
import { roleLabels } from '../utils/labels'
import { matchesSearch } from '../utils/search'

export function UsersPage() {
  const [users, setUsers] = useState<User[]>([])
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')
  const [roleFilter, setRoleFilter] = useState<Role | ''>('')

  useEffect(() => {
    setLoading(true)
    api.listUsers(roleFilter || undefined)
      .then(setUsers)
      .catch((err: Error) => setError(err.message))
      .finally(() => setLoading(false))
  }, [roleFilter])

  const filteredUsers = useMemo(
    () => users.filter((user) => matchesSearch(search, [user.email, user.id, roleLabels[user.role]])),
    [users, search],
  )

  return (
    <div className="page">
      <header className="page-header">
        <div>
          <h1>Пользователи</h1>
          <p className="muted">Список клиентов, менеджеров и администраторов</p>
        </div>
      </header>

      {error && <div className="alert alert-error">{error}</div>}

      <FiltersBar
        search={search}
        onSearchChange={setSearch}
        searchPlaceholder="Email, ID или роль"
      >
        <label>
          Роль
          <select
            value={roleFilter}
            onChange={(e) => setRoleFilter(e.target.value as Role | '')}
          >
            <option value="">Все роли</option>
            <option value="client">Клиент</option>
            <option value="manager">Менеджер</option>
            <option value="admin">Администратор</option>
          </select>
        </label>
      </FiltersBar>

      {loading ? (
        <p className="muted">Загрузка…</p>
      ) : filteredUsers.length === 0 ? (
        <div className="card empty-state">Пользователи не найдены</div>
      ) : (
        <div className="table-card">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>Email</th>
                <th>Роль</th>
                <th>Создан</th>
              </tr>
            </thead>
            <tbody>
              {filteredUsers.map((user) => (
                <tr key={user.id}>
                  <td>#{user.id}</td>
                  <td>{user.email}</td>
                  <td>
                    <span className="badge badge-neutral">{roleLabels[user.role]}</span>
                  </td>
                  <td>
                    {user.created_at
                      ? new Date(user.created_at).toLocaleDateString('ru-RU')
                      : '—'}
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
