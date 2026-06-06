# Requests API

Backend + frontend личного кабинета клиентов.

- **Backend:** Go + SQLite
- **Frontend:** React + Vite + TypeScript (`frontend/`)

## Возможности

- REST API: авторизация, профиль, заявки, уведомления
- Web UI: логин, регистрация, заявки, профиль, уведомления
- Роли: `client`, `manager`, `admin`
- Email/SMS/in-app уведомления
- JWT-аутентификация
- Документация: [docs/API.md](docs/API.md), [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)

## Быстрый старт (Docker — одна команда)

Требуется [Docker Desktop](https://www.docker.com/products/docker-desktop/).

```bash
git clone https://github.com/make-smart-products/requests-api.git
cd requests-api
docker compose up --build
```

Откройте **http://localhost:3000**

Логин: `admin@requests.local` / `admin12345`

Остановка: `Ctrl+C` или `docker compose down`

Опционально — свой JWT-секрет:

```bash
cp .env.docker.example .env
docker compose up --build
```

## Локальная разработка

### Backend

```bash
cp .env.example .env
go mod tidy
go run ./cmd/api
```

API: `http://localhost:8080`

### Frontend

```bash
cd frontend
npm install
npm run dev
```

UI: `http://localhost:5173`

Подробнее: [frontend/README.md](frontend/README.md)

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
frontend/                — React UI
docker-compose.yml       — запуск API + UI
docs/                    — документация
```

## Лицензия

MIT
