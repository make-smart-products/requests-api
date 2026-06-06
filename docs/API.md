# API Documentation

Base URL: `http://localhost:8080`

Все ответы — JSON. Ошибки:

```json
{ "error": "описание ошибки" }
```

## Auth

### POST /api/v1/auth/register

Регистрация клиента.

```json
{
  "email": "client@example.com",
  "password": "password123",
  "first_name": "Иван",
  "last_name": "Иванов",
  "phone": "+79001234567",
  "company": "ООО Пример"
}
```

Ответ `201`:

```json
{
  "token": "jwt-token",
  "user": { "id": 1, "email": "client@example.com", "role": "client" }
}
```

### POST /api/v1/auth/login

```json
{
  "email": "client@example.com",
  "password": "password123"
}
```

Ответ `200` — как у register.

## Profile

Требуется заголовок: `Authorization: Bearer <token>`

### GET /api/v1/profile

Текущий профиль пользователя.

### PUT /api/v1/profile

```json
{
  "first_name": "Иван",
  "last_name": "Иванов",
  "phone": "+79001234567",
  "company": "ООО Пример"
}
```

## Users

Доступ: `manager`, `admin`

### GET /api/v1/users?role=client

Список пользователей, опциональный фильтр по роли.

## Applications

### GET /api/v1/applications?status=submitted

- `client` — только свои заявки
- `manager` — назначенные заявки
- `admin` — все заявки

### POST /api/v1/applications

```json
{
  "title": "Подключение услуги",
  "description": "Нужна консультация",
  "status": "draft"
}
```

### GET /api/v1/applications/{id}

Детали заявки.

### PATCH /api/v1/applications/{id}

```json
{
  "status": "in_review",
  "assigned_manager_id": 2
}
```

### DELETE /api/v1/applications/{id}

Только клиент, только статус `draft`.

## Notifications

### GET /api/v1/notifications?unread=true

Список уведомлений текущего пользователя.

### PATCH /api/v1/notifications/{id}/read

Пометить уведомление прочитанным.

## Health

### GET /health

```json
{ "status": "ok" }
```

## Коды ответов

| Код | Значение |
|-----|----------|
| 200 | Успех |
| 201 | Создано |
| 204 | Удалено |
| 400 | Невалидные данные |
| 401 | Не авторизован |
| 403 | Нет прав |
| 404 | Не найдено |
| 409 | Конфликт (например, email занят) |
| 500 | Внутренняя ошибка |

## Тестовый админ

При первом запуске создаётся:

- email: `admin@requests.local`
- password: `admin12345`
