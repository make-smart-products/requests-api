# Деплой: API на VPS, UI на Vercel/Netlify

## Архитектура

```text
Пользователь
    ↓
Vercel / Netlify  →  React UI (статика)
    ↓ HTTPS
VPS (nginx)       →  Go API :8080
    ↓
SQLite volume
```

---

## 1. API на VPS (Ubuntu)

### Требования

- Ubuntu 22.04+
- Docker **или** Go 1.22+
- Домен, например `api.example.com`

### Вариант A — Docker (рекомендуется)

```bash
ssh user@your-vps
git clone https://github.com/make-smart-products/requests-api.git
cd requests-api
cp .env.example .env
```

Отредактируйте `.env`:

```env
APP_ENV=production
HTTP_ADDR=:8080
JWT_SECRET=длинный-случайный-секрет
DB_PATH=/app/data/requests.db
CORS_ALLOWED_ORIGINS=https://your-app.vercel.app,https://your-app.netlify.app
```

Запуск только API:

```bash
docker build -t requests-api .
docker run -d \
  --name requests-api \
  --restart unless-stopped \
  -p 127.0.0.1:8080:8080 \
  -v requests-api-data:/app/data \
  --env-file .env \
  requests-api
```

### Вариант B — systemd + бинарник

```bash
go build -o /opt/requests-api/api ./cmd/api
sudo cp deploy/systemd/requests-api.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now requests-api
```

### Nginx + HTTPS

```bash
sudo apt install nginx certbot python3-certbot-nginx
sudo cp deploy/nginx/api.conf.example /etc/nginx/sites-available/requests-api
# замените api.example.com на ваш домен
sudo ln -s /etc/nginx/sites-available/requests-api /etc/nginx/sites-enabled/
sudo nginx -t
sudo certbot --nginx -d api.example.com
sudo systemctl reload nginx
```

Проверка:

```bash
curl https://api.example.com/health
```

---

## 2. UI на Vercel

### Шаги

1. https://vercel.com → **Add New Project**
2. Import репозитория `make-smart-products/requests-api`
3. **Root Directory:** `frontend`
4. **Framework Preset:** Vite
5. **Environment Variable:**

| Name | Value |
|------|-------|
| `VITE_API_URL` | `https://api.example.com` |

6. Deploy

`frontend/vercel.json` уже настроен для SPA-роутинга.

### Обновите CORS на VPS

Добавьте URL Vercel в `.env`:

```env
CORS_ALLOWED_ORIGINS=https://your-app.vercel.app
```

Перезапустите API.

---

## 3. UI на Netlify

1. https://app.netlify.com → **Add new site** → Git
2. Репозиторий `requests-api`
3. **Base directory:** `frontend`
4. **Build command:** `npm run build`
5. **Publish directory:** `frontend/dist`
6. **Environment variable:**

| Name | Value |
|------|-------|
| `VITE_API_URL` | `https://api.example.com` |

`frontend/netlify.toml` уже содержит SPA redirect.

Добавьте Netlify URL в `CORS_ALLOWED_ORIGINS` на VPS.

---

## 4. Чеклист после деплоя

- [ ] `https://api.example.com/health` → `{"status":"ok"}`
- [ ] UI открывается на Vercel/Netlify
- [ ] Логин `admin@requests.local` / `admin12345` работает
- [ ] Создание заявки работает
- [ ] CORS не блокирует запросы (проверить DevTools → Network)

---

## 5. Безопасность (production)

- Смените пароль админа после первого входа
- Используйте длинный `JWT_SECRET` (32+ символов)
- Включите HTTPS на API и UI
- Не коммитьте `.env` в git
- Регулярно обновляйте VPS и зависимости
