import { NavLink, Outlet } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { roleLabels } from '../utils/labels'

export function Layout() {
  const { user, logout } = useAuth()

  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div className="brand">
          <span className="brand-mark">R</span>
          <div>
            <strong>Requests</strong>
            <p>Личный кабинет</p>
          </div>
        </div>

        <nav className="nav">
          <NavLink to="/" end>Обзор</NavLink>
          <NavLink to="/applications">Заявки</NavLink>
          {(user?.role === 'admin' || user?.role === 'manager') && (
            <NavLink to="/users">Пользователи</NavLink>
          )}
          <NavLink to="/profile">Профиль</NavLink>
          <NavLink to="/notifications">Уведомления</NavLink>
        </nav>

        <div className="sidebar-footer">
          <div className="user-card">
            <span className="user-email">{user?.email}</span>
            <span className="user-role">{user ? roleLabels[user.role] : ''}</span>
          </div>
          <button type="button" className="btn btn-ghost" onClick={logout}>
            Выйти
          </button>
        </div>
      </aside>

      <main className="main">
        <Outlet />
      </main>
    </div>
  )
}
