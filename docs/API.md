# API Документация — Taxi

Документация по работе с HTTP API сервиса заказа такси.

## Содержание

1. [Базовый URL](#базовый-url)
2. [Переменные окружения](#переменные-окружения)
3. [Эндпоинты](#эндпоинты)
   - [Health Check](#1-health-check)
   - [Авторизация (OTP)](#2-авторизация-otp)
   - [Отправка OTP (legacy)](#2b-отправка-otp-legacy)
   - [Регистрация пользователя](#3-регистрация-пользователя)
   - [Получить пользователя](#4-получить-пользователя)
   - [Получить пользователя по телефону](#4b-получить-пользователя-по-телефону)
   - [Редактирование пользователя](#5-редактирование-пользователя)
   - [Удаление пользователя](#6-удаление-пользователя)
   - [Регистрация автомобиля](#7-регистрация-автомобиля)
   - [Список авто водителя](#8-список-авто-водителя)
   - [Получить / редактировать / удалить авто](#9-автомобиль)
4. [Формат номера телефона](#формат-номера-телефона)
5. [Коды ответов HTTP](#коды-ответов-http)
6. [Краткая справка](#краткая-справка)
7. [Модель User](#модель-user)
8. [Структура проекта](#структура-проекта-api)
9. [Схема БД](#схема-бд)
10. [CORS](#cors)

---

## Базовый URL

```
http://localhost:8080
```

В продакшене замените на ваш домен.

---

## Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `DATABASE_URL` | DSN PostgreSQL | `postgres://postgres:postgres@localhost:5432/taxi?sslmode=disable` |
| `API_URL` | URL Green API | `https://api.green-api.com` |
| `ID_INSTANCE` | ID инстанса Green API | — |
| `API_TOKEN` | Токен Green API | — |
| `USE_MOCK_WHATSAPP` | `true` — логировать вместо отправки | `false` |
| `PORT` | Порт сервера | `8080` |
| `MIGRATIONS_DIR` | Путь к папке миграций | `migrations` |

---

## Эндпоинты

### 1. Health Check

Проверка работоспособности сервера.

**Запрос:**
```
GET /health
```

**Ответ (200 OK):**
```json
{
  "status": "ok"
}
```

**Пример:**
```bash
curl http://localhost:8080/health
```

---

### 2. Авторизация (OTP)

Вход по номеру телефона. Код генерируется сервером и отправляется в WhatsApp.

**Шаг 1 — Запросить код:**
```
POST /api/auth/send-otp
Content-Type: application/json
```

```json
{ "phone": "77001234567" }
```

**Ответ (200 OK):** `{"status": "sent"}`

**Шаг 2 — Ввести код и получить пользователя:**
```
POST /api/auth/verify
Content-Type: application/json
```

```json
{
  "phone": "77001234567",
  "code": "1234"
}
```

**Ответ при успехе (200 OK):** объект User (если пользователь зарегистрирован)

**Ошибки:**
| Код | Причина |
|-----|---------|
| 401 | Неверный или истёкший код |
| 404 | Пользователь не найден — нужно зарегистрироваться |

Код действителен 5 минут.

**Тестовые пользователи** (код 0000, WhatsApp не вызывается):
| Телефон | Имя | Роль |
|---------|-----|------|
| 77000000000 | Samat | passenger |
| 77000000001 | Nurik | driver |

**Примеры:**
```bash
# 1. Запросить код (придёт в WhatsApp)
curl -X POST http://localhost:8080/api/auth/send-otp \
  -H "Content-Type: application/json" \
  -d '{"phone":"77001234567"}'

# 2. Ввести код из WhatsApp — получить пользователя
curl -X POST http://localhost:8080/api/auth/verify \
  -H "Content-Type: application/json" \
  -d '{"phone":"77001234567","code":"1234"}'
```

---

### 2b. Отправка OTP (legacy)

Ручная отправка кода (для тестов). В продакшене используйте `/api/auth/send-otp`.

```
POST /api/send-otp
Body: { "phone": "77001234567", "code": "1234" }
```

**Форматы номера телефона (принимаются):**
- `77001234567`
- `+77001234567`
- `87001234567`
- ` 77001234567 ` (с пробелами)

**Ответ при успехе (200 OK):**
```json
{
  "status": "sent"
}
```

**Ошибки:**

| Код | Причина |
|-----|---------|
| 400 | Невалидный JSON |
| 400 | Неверный формат телефона (ожидается 7XXXXXXXXXX) |
| 400 | Код должен быть 4 цифры |
| 405 | Метод не POST |
| 500 | Ошибка Green API (инстанс недоступен и т.д.) |

**Примеры:**

```bash
# Успешная отправка
curl -X POST http://localhost:8080/api/send-otp \
  -H "Content-Type: application/json" \
  -d '{"phone":"77001234567","code":"1234"}'

# Неверный формат телефона
curl -X POST http://localhost:8080/api/send-otp \
  -H "Content-Type: application/json" \
  -d '{"phone":"123","code":"1234"}'
# 400 Invalid phone format (expected 7XXXXXXXXXX)

# Код не 4 цифры
curl -X POST http://localhost:8080/api/send-otp \
  -H "Content-Type: application/json" \
  -d '{"phone":"77001234567","code":"12"}'
# 400 Code must be 4 digits
```

---

### 3. Регистрация пользователя

Создание нового пользователя (пассажир или водитель).

**Запрос:**
```
POST /api/users
Content-Type: application/json
```

**Тело запроса:**
```json
{
  "phone": "77001234567",
  "name": "Иван",
  "role": "passenger"
}
```

| Поле | Тип | Обязательное | Описание |
|------|-----|--------------|----------|
| `phone` | string | да | Номер телефона (7XXXXXXXXXX) |
| `name` | string | да | Имя или позывной |
| `role` | string | да | `passenger` или `driver` |

**Ответ при успехе (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "phone": "77001234567",
  "name": "Иван",
  "role": "passenger",
  "is_active": false,
  "created_at": "2026-03-01T12:00:00Z"
}
```

**Ошибки:**

| Код | Сообщение | Причина |
|-----|-----------|---------|
| 400 | Invalid phone format (expected 7XXXXXXXXXX) | Неверный формат номера |
| 400 | Name is required | Поле name пустое |
| 400 | Role must be passenger or driver | role не passenger и не driver |
| 409 | Phone already registered | Номер уже зарегистрирован |

---

### 4. Получить пользователя

**Запрос:**
```
GET /api/users/{id}
```

**Ответ (200 OK):** объект User (как выше)

**Ошибки:**

| Код | Сообщение | Причина |
|-----|-----------|---------|
| 404 | User not found | Пользователь с таким id не найден |

---

### 4b. Получить пользователя по телефону

**Запрос:**
```
GET /api/users/phone/{phone}
```

**Пример:** `GET /api/users/phone/77001234567`

**Ответ (200 OK):** объект User

**Ошибки:** 404 (пользователь не найден)

---

### 5. Редактирование пользователя

Обновление данных пользователя. Все поля опциональны — обновляются только переданные.

**Запрос:**
```
PUT /api/users/{id}
PATCH /api/users/{id}
Content-Type: application/json
```

**Тело запроса:**
```json
{
  "phone": "77009876543",
  "name": "Иван Петров",
  "role": "driver",
  "is_active": true
}
```

| Поле | Тип | Описание |
|------|-----|----------|
| `phone` | string | Новый номер |
| `name` | string | Новое имя |
| `role` | string | `passenger` или `driver` |
| `is_active` | boolean | Статус верификации |

**Ответ (200 OK):** обновлённый объект User

**Ошибки:**

| Код | Сообщение | Причина |
|-----|-----------|---------|
| 400 | Invalid phone format | Неверный формат номера |
| 400 | Name is required | name передан пустым |
| 400 | Role must be passenger or driver | Недопустимая роль |
| 404 | User not found | Пользователь не найден |
| 409 | Phone already registered | Номер уже занят другим пользователем |

---

### 6. Удаление пользователя

**Запрос:**
```
DELETE /api/users/{id}
```

**Ответ (204 No Content):** пустое тело

**Ошибки:**

| Код | Сообщение | Причина |
|-----|-----------|---------|
| 404 | User not found | Пользователь не найден |

**Примеры пользователей:**

```bash
# Регистрация пассажира
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"phone":"77001234567","name":"Алексей","role":"passenger"}'

# Регистрация водителя
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"phone":"77009876543","name":"Сергей","role":"driver"}'

# Получить пользователя (подставьте id из ответа регистрации)
curl http://localhost:8080/api/users/550e8400-e29b-41d4-a716-446655440000

# Редактировать
curl -X PATCH http://localhost:8080/api/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{"name":"Алексей Иванов","is_active":true}'

# Удалить
curl -X DELETE http://localhost:8080/api/users/550e8400-e29b-41d4-a716-446655440000
```

---

### 7. Регистрация автомобиля

Добавление автомобиля для водителя. **Только водители** могут иметь автомобили.

**Запрос:**
```
POST /api/cars
Content-Type: application/json
```

**Тело запроса:**
```json
{
  "driver_id": "550e8400-e29b-41d4-a716-446655440000",
  "model": "Toyota Camry",
  "number": "777 ABC 01",
  "color": "Белый"
}
```

| Поле | Тип | Обязательное | Описание |
|------|-----|--------------|----------|
| `driver_id` | string (UUID) | да | ID водителя |
| `model` | string | да | Марка/модель |
| `number` | string | да | Гос. номер |
| `color` | string | да | Цвет |

**Ответ (201 Created):**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "driver_id": "550e8400-e29b-41d4-a716-446655440000",
  "model": "Toyota Camry",
  "number": "777 ABC 01",
  "color": "Белый"
}
```

**Ошибки:** 400 (не водитель, пустые поля), 404 (водитель не найден)

---

### 8. Список авто водителя

**Запрос:**
```
GET /api/users/{driver_id}/cars
```

**Ответ (200 OK):** массив Car

**Ошибки:** 400 (не водитель), 404 (водитель не найден)

---

### 9. Автомобиль

**Получить:** `GET /api/cars/{id}`  
**Редактировать:** `PUT /api/cars/{id}` или `PATCH /api/cars/{id}`  
**Удалить:** `DELETE /api/cars/{id}`

**Тело для PATCH/PUT (все поля опциональны):**
```json
{
  "model": "Toyota Camry 70",
  "number": "123 ABC 01",
  "color": "Чёрный"
}
```

**Примеры:**
```bash
# Добавить авто (driver_id — ID водителя)
curl -X POST http://localhost:8080/api/cars \
  -H "Content-Type: application/json" \
  -d '{"driver_id":"UUID","model":"Toyota Camry","number":"777 ABC 01","color":"Белый"}'

# Список авто водителя
curl http://localhost:8080/api/users/{driver_id}/cars

# Редактировать
curl -X PATCH http://localhost:8080/api/cars/{car_id} \
  -H "Content-Type: application/json" \
  -d '{"color":"Серый"}'

# Удалить
curl -X DELETE http://localhost:8080/api/cars/{car_id}
```

---

## Формат номера телефона

Принимаются номера Казахстана:
- **11 цифр**, начинаются с `7` (код страны)
- Пример: `77001234567`, `7771234567`

Поддерживаемые варианты ввода:
- `77001234567` — основной формат
- `+77001234567` — с плюсом
- `87001234567` — старый формат (8 заменяется на 7)
- `7001234567` — 10 цифр (добавляется 7 в начало для формата 70X)

---

## Коды ответов HTTP

| Код | Описание |
|-----|----------|
| 200 | Успех |
| 201 | Создано (регистрация) |
| 204 | Успех без тела (удаление) |
| 400 | Ошибка в запросе (валидация) |
| 401 | Не авторизован (неверный/истёкший OTP) |
| 404 | Не найдено |
| 405 | Метод не разрешён |
| 409 | Конфликт (телефон уже занят) |
| 500 | Внутренняя ошибка сервера |

---

## Краткая справка

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/health` | Проверка сервера |
| POST | `/api/auth/send-otp` | Запросить код (генерируется сервером) |
| POST | `/api/auth/verify` | Ввести код, получить пользователя |
| POST | `/api/send-otp` | Отправка OTP вручную (legacy) |
| POST | `/api/users` | Регистрация |
| GET | `/api/users/{id}` | Получить пользователя |
| GET | `/api/users/phone/{phone}` | Получить пользователя по телефону |
| PUT | `/api/users/{id}` | Полное обновление |
| PATCH | `/api/users/{id}` | Частичное обновление |
| DELETE | `/api/users/{id}` | Удаление |
| POST | `/api/cars` | Регистрация авто (только водители) |
| GET | `/api/cars/{id}` | Получить авто |
| GET | `/api/users/{driver_id}/cars` | Список авто водителя |
| PUT | `/api/cars/{id}` | Обновить авто |
| PATCH | `/api/cars/{id}` | Обновить авто |
| DELETE | `/api/cars/{id}` | Удалить авто |

---

## Модель User

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | string (UUID) | Уникальный идентификатор |
| `phone` | string | Номер телефона (7XXXXXXXXXX) |
| `name` | string | Имя или позывной |
| `role` | string | `passenger` или `driver` |
| `is_active` | boolean | Статус верификации (по умолчанию false) |
| `created_at` | string (RFC3339) | Время регистрации |

---

## Модель Car

| Поле | Тип | Описание |
|------|-----|----------|
| `id` | string (UUID) | Уникальный идентификатор |
| `driver_id` | string (UUID) | ID водителя |
| `model` | string | Марка/модель |
| `number` | string | Гос. номер |
| `color` | string | Цвет |

---

## Структура проекта API

```
cmd/api/main.go          — точка входа
internal/handler/        — обработчики HTTP
  health.go              — GET /health
  auth.go                — POST /api/auth/send-otp, /api/auth/verify
  otp.go                 — POST /api/send-otp (legacy)
  user.go                — CRUD /api/users
  car.go                 — CRUD /api/cars
internal/models/         — модели данных
internal/repository/     — работа с БД
internal/service/        — бизнес-логика
pkg/whatsapp/            — Green API (отправка в WhatsApp)
pkg/validator/           — валидация (телефон)
```

---

## Схема БД

### Таблица `users`

| Поле | Тип | Описание |
|------|-----|----------|
| id | UUID | Первичный ключ |
| phone | VARCHAR(20) | Номер (уникальный) |
| name | VARCHAR(255) | Имя или позывной |
| role | VARCHAR(20) | `passenger` или `driver` |
| is_active | BOOLEAN | Статус верификации |
| created_at | TIMESTAMP | Время регистрации |

### Таблица `cars` (для водителей)

| Поле | Тип | Описание |
|------|-----|----------|
| id | UUID | Первичный ключ |
| driver_id | UUID | FK → users.id |
| model | VARCHAR(255) | Марка/модель |
| number | VARCHAR(20) | Гос. номер |
| color | VARCHAR(100) | Цвет |

### Таблица `otp_codes` (коды для авторизации)

| Поле | Тип | Описание |
|------|-----|----------|
| phone | VARCHAR(20) | PK, номер телефона |
| code | VARCHAR(4) | 4-значный код |
| expires_at | TIMESTAMP | Время истечения (5 мин) |

### Таблица `orders` (архив заказов)

| Поле | Тип | Описание |
|------|-----|----------|
| id | UUID | Первичный ключ |
| passenger_id | UUID | FK → users.id |
| driver_id | UUID | FK → users.id (nullable) |
| status | VARCHAR(20) | `created`, `completed`, `cancelled` |
| price | INTEGER | Сумма в тенге |
| created_at | TIMESTAMP | Время создания |

---

## CORS

API поддерживает CORS для запросов с фронтенда. Разрешены origins `https://*`, `http://*`, методы GET, POST, PUT, PATCH, DELETE, OPTIONS.
