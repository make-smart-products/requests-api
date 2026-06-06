import type { ApplicationStatus, Role } from '../types'

export const roleLabels: Record<Role, string> = {
  client: 'Клиент',
  manager: 'Менеджер',
  admin: 'Администратор',
}

export const statusLabels: Record<ApplicationStatus, string> = {
  draft: 'Черновик',
  submitted: 'Отправлена',
  in_review: 'На рассмотрении',
  approved: 'Одобрена',
  rejected: 'Отклонена',
}

export const statusClass: Record<ApplicationStatus, string> = {
  draft: 'badge-neutral',
  submitted: 'badge-info',
  in_review: 'badge-warning',
  approved: 'badge-success',
  rejected: 'badge-danger',
}
