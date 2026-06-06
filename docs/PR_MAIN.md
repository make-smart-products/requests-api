# Pull Request: feature/auth → main

## Summary

- Go backend личного кабинета: auth, profile, applications, notifications
- Роли: client, manager, admin
- SQLite + JWT + email/SMS (dev log mode)
- Документация: `docs/API.md`, `docs/ARCHITECTURE.md`
- GitHub Actions CI: тесты и сборка

## Test plan

- [ ] `go test ./...` проходит локально
- [ ] `go run ./cmd/api` поднимает сервер на `:8080`
- [ ] `GET /health` возвращает `{"status":"ok"}`
- [ ] `POST /api/v1/auth/login` с admin@requests.local работает
- [ ] `POST /api/v1/auth/register` создаёт клиента
- [ ] `POST /api/v1/applications` создаёт заявку с токеном
- [ ] CI workflow зелёный после merge
