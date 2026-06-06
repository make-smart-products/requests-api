import { Navigate, Outlet } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import type { Role } from '../types'

export function RoleRoute({ roles }: { roles: Role[] }) {
  const { user, isAuthenticated } = useAuth()

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  if (!user || !roles.includes(user.role)) {
    return <Navigate to="/" replace />
  }

  return <Outlet />
}
