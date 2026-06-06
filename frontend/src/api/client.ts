import type {
  Application,
  ApplicationStatus,
  AuthResponse,
  Notification,
  Profile,
  User,
} from '../types'

const API_BASE = import.meta.env.VITE_API_URL ?? ''

class ApiError extends Error {
  status: number

  constructor(status: number, message: string) {
    super(message)
    this.status = status
  }
}

function getToken(): string | null {
  return localStorage.getItem('token')
}

async function request<T>(
  path: string,
  options: RequestInit = {},
  auth = true,
): Promise<T> {
  const headers = new Headers(options.headers)
  headers.set('Content-Type', 'application/json')

  if (auth) {
    const token = getToken()
    if (token) {
      headers.set('Authorization', `Bearer ${token}`)
    }
  }

  const response = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers,
  })

  if (response.status === 204) {
    return undefined as T
  }

  const payload = await response.json().catch(() => ({}))
  if (!response.ok) {
    throw new ApiError(response.status, payload.error ?? 'Ошибка запроса')
  }

  return payload as T
}

export const api = {
  login(email: string, password: string) {
    return request<AuthResponse>('/api/v1/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    }, false)
  },

  register(data: {
    email: string
    password: string
    first_name: string
    last_name: string
    phone: string
    company?: string
  }) {
    return request<AuthResponse>('/api/v1/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
    }, false)
  },

  getProfile() {
    return request<Profile>('/api/v1/profile')
  },

  updateProfile(data: Partial<Profile>) {
    return request<Profile>('/api/v1/profile', {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  },

  listApplications(status?: ApplicationStatus) {
    const query = status ? `?status=${status}` : ''
    return request<Application[]>(`/api/v1/applications${query}`)
  },

  getApplication(id: number) {
    return request<Application>(`/api/v1/applications/${id}`)
  },

  createApplication(data: {
    title: string
    description: string
    status?: ApplicationStatus
  }) {
    return request<Application>('/api/v1/applications', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  updateApplication(
    id: number,
    data: {
      title?: string
      description?: string
      status?: ApplicationStatus
      assigned_manager_id?: number
    },
  ) {
    return request<Application>(`/api/v1/applications/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    })
  },

  deleteApplication(id: number) {
    return request<void>(`/api/v1/applications/${id}`, {
      method: 'DELETE',
    })
  },

  listNotifications(unreadOnly = false) {
    const query = unreadOnly ? '?unread=true' : ''
    return request<Notification[]>(`/api/v1/notifications${query}`)
  },

  markNotificationRead(id: number) {
    return request<Notification>(`/api/v1/notifications/${id}/read`, {
      method: 'PATCH',
    })
  },

  listUsers(role?: string) {
    const query = role ? `?role=${role}` : ''
    return request<User[]>(`/api/v1/users${query}`)
  },
}

export { ApiError }
