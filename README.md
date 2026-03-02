# Taxi — Сервис заказа такси

Проект на Go с Clean Architecture. Регистрация пассажиров и водителей, отправка OTP через WhatsApp (Green API).

## Возможности

- Авторизация по OTP (код в WhatsApp, возврат пользователя если зарегистрирован)
- Регистрация, редактирование и удаление пользователей (пассажиры и водители)
- Регистрация и редактирование автомобилей (только для водителей)
- Отправка OTP-кодов через WhatsApp (Green API)
- Валидация казахстанских номеров (7XXXXXXXXXX)

## Структура проекта

```
taxi/
├── cmd/api/main.go          # Точка входа
├── internal/
│   ├── models/              # User, Order, Car
│   ├── repository/          # Работа с БД, миграции
│   ├── service/             # Бизнес-логика
│   └── handler/             # HTTP-эндпоинты
├── pkg/
│   ├── whatsapp/            # Green API (Messenger interface)
│   └── validator/           # Валидация (телефон)
├── migrations/               # SQL-миграции
├── scripts/create_db.go      # Создание БД taxi
├── docs/API.md               # Документация API
└── docker-compose.yml        # PostgreSQL для разработки
```

## Запуск

1. Запустить PostgreSQL:
   ```bash
   docker compose up -d
   ```
   Если база `taxi` не существует: `go run scripts/create_db.go`

2. Скопировать `.env.example` в `.env` и настроить переменные.

3. Запустить API:
   ```bash
   go run cmd/api/main.go
   ```

4. Проверить: `curl http://localhost:8080/health`

## Переменные окружения

| Переменная | Описание |
|------------|----------|
| `DATABASE_URL` | DSN PostgreSQL |
| `API_URL` | URL Green API |
| `ID_INSTANCE` | ID инстанса Green API |
| `API_TOKEN` | Токен Green API |
| `USE_MOCK_WHATSAPP` | `true` — логировать вместо отправки |
| `PORT` | Порт сервера (по умолчанию 8080) |
| `MIGRATIONS_DIR` | Путь к миграциям (по умолчанию `migrations`) |

## Документация API

Подробное описание эндпоинтов и примеры запросов: [docs/API.md](docs/API.md)

## Тесты

```bash
go test ./...
```

Покрыты: validator, repository (User, OTP), service (User, Auth), handler (Health, Auth).
