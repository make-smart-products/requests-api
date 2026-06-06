# Frontend — Личный кабинет

React + Vite + TypeScript UI для Requests API.

## Запуск

### Docker (рекомендуется)

Из корня репозитория:

```bash
docker compose up --build
```

UI: http://localhost:3000

### Локально (dev)

```bash
# Терминал 1 — backend
cd ..
go run ./cmd/api

# Терминал 2 — frontend
cd frontend
npm install
npm run dev
```

UI: http://localhost:5173

## Страницы

- `/login` — вход
- `/register` — регистрация клиента
- `/` — обзор
- `/applications` — заявки (фильтры, поиск)
- `/users` — пользователи (admin, manager)
- `/profile` — профиль
- `/notifications` — уведомления

## Деплой

- **Vercel:** root = `frontend`, env `VITE_API_URL=https://api.example.com`
- **Netlify:** base = `frontend`, env `VITE_API_URL=https://api.example.com`

Подробнее: [../docs/DEPLOY.md](../docs/DEPLOY.md)

## Тестовый вход

- `admin@requests.local` / `admin12345`
