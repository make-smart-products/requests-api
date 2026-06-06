# Requests API

Backend личного кабинета клиентов на **Go + SQLite**.

## Возможности

- REST API: авторизация, профиль, заявки, уведомления
- Роли: `client`, `manager`, `admin`
- Email/SMS/in-app уведомления
- JWT-аутентификация
- Документация: [docs/API.md](docs/API.md), [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)

## Быстрый старт

```bash
git clone https://github.com/make-smart-products/requests-api.git
cd requests-api
cp .env.example .env
go mod tidy
go run ./cmd/api
```

Сервер: `http://localhost:8080`

## Примеры запросов

```bash
# Health
curl http://localhost:8080/health

# Регистрация
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"client@example.com","password":"password123","first_name":"Иван","last_name":"Иванов","phone":"+79001234567"}'

# Логин админа
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@requests.local","password":"admin12345"}'

# Создать заявку
curl -X POST http://localhost:8080/api/v1/applications \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Новая заявка","description":"Описание","status":"submitted"}'
```

## Конфигурация

См. `.env.example`:

| Переменная | Описание |
|------------|----------|
| `HTTP_ADDR` | Адрес сервера (`:8080`) |
| `JWT_SECRET` | Секрет для JWT |
| `DB_PATH` | Путь к SQLite (`./data/requests.db`) |
| `SMTP_*` | Настройки email |
| `SMS_PROVIDER` | `log` для dev, иначе внешний провайдер |

## Структура проекта

```text
cmd/api/                 — main
internal/config/         — конфигурация
internal/handler/        — HTTP handlers
internal/service/        — бизнес-логика
internal/repository/     — SQLite
internal/notification/   — email/SMS
docs/                    — документация
```

## Push в GitHub

```bash
git init
git add .
git commit -m "feat: initial Go backend for client cabinet"
git branch -M main
git remote add origin https://github.com/make-smart-products/requests-api.git
git push -u origin main
```

## Лицензия

MIT
