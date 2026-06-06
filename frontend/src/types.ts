export type Role = 'client' | 'manager' | 'admin'

export type ApplicationStatus =
  | 'draft'
  | 'submitted'
  | 'in_review'
  | 'approved'
  | 'rejected'

export interface User {
  id: number
  email: string
  role: Role
  created_at?: string
  updated_at?: string
}

export interface Profile {
  user_id: number
  first_name: string
  last_name: string
  phone: string
  company: string
  created_at?: string
  updated_at?: string
}

export interface Application {
  id: number
  user_id: number
  assigned_manager_id?: number
  title: string
  description: string
  status: ApplicationStatus
  created_at: string
  updated_at: string
}

export interface Notification {
  id: number
  user_id: number
  channel: 'email' | 'sms' | 'in_app'
  type: string
  title: string
  body: string
  is_read: boolean
  sent_at?: string
  created_at: string
}

export interface AuthResponse {
  token: string
  user: User
}
