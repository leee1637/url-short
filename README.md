# URL Shortener

Сервис сокращения ссылок на Go. Одна команда — и работает: Postgres поднимается в Docker, миграции применяются автоматически, сервис слушает `:8080`.

## Стек

- **Go 1.26** + `log/slog`
- HTTP-фреймворк: **Gin** (`gin-gonic/gin`)
- Postgres: **pgx/v5** (`pgxpool`)
- Миграции: **golang-migrate**
- Конфиг: **cleanenv** (yaml + env)
- Архитектура: три слоя `handler → service → storage`

## Быстрый старт

```bash
git clone github.com/leee1637/url-short
cd url-short

cp .env.example .env
docker compose up
```

Сервис доступен на `http://localhost:8888`.

### Только локальная разработка (без контейнера бэкенда)

```bash
docker compose up db migrate   # поднять БД и применить миграции
cp .env.example .env
go run ./cmd/url-short
```

## Конфигурация

Секреты живут только в `.env` (в коде их нет). Значения в `.env.example`, `docker-compose.yml` и `config/local.yaml` совпадают.

| Переменная | Значение по умолчанию | Описание |
|---|---|---|
| `APP_ENV` | `local` | режим окружения (`local` / `production`) |
| `HTTP_ADDR` | `:8888` | адрес HTTP-сервера |
| `DB_HOST` | `localhost` | хост Postgres |
| `DB_PORT` | `5432` | порт Postgres |
| `DB_USER` | `shorter` | пользователь БД |
| `DB_PASS` | `shorter123` | пароль БД |
| `DB_NAME` | `urlshort` | имя БД |
| `AUTH_USER` | `admin` | логин для Basic Auth |
| `AUTH_PASS` | `admin123` | пароль для Basic Auth |

## API

Ручки `/urls*` защищены Basic Auth. Публичная ручка — только редирект.

| Метод | Путь | Описание | Успех |
|---|---|---|---|
| `POST` | `/urls` | создание ссылки. Тело: `{"url":"...","alias":"..."}` — alias опционален, сгенерируется автоматически | `201` |
| `GET` | `/:alias` | редирект на исходный URL (публичная, без auth) | `302` |
| `GET` | `/urls/:alias` | информация о ссылке | `200` |
| `PATCH` | `/urls/:alias` | обновить URL ссылки. Тело: `{"url":"..."}` | `200` |
| `DELETE` | `/urls/:alias` | удалить ссылку | `200` |

### Примеры запросов

```bash
# создать ссылку (алиас сгенерируется)
curl -u admin:admin123 -X POST http://localhost:8080/urls \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com"}'

# создать с собственным алиасом
curl -u admin:admin123 -X POST http://localhost:8080/urls \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com","alias":"example"}'

# редирект
curl -L http://localhost:8080/example

# обновить
curl -u admin:admin123 -X PATCH http://localhost:8080/urls/example \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://new.example.com"}'

# удалить
curl -u admin:admin123 -X DELETE http://localhost:8080/urls/example
```

### Коды статусов

| Ситуация | Статус |
|---|---|
| успешное создание | `201 Created` |
| редирект | `302 Found` |
| битый JSON | `400 Bad Request` |
| невалидный URL или alias | `400 Bad Request` |
| alias уже занят | `409 Conflict` |
| alias не найден | `404 Not Found` |
| нет или неверная авторизация | `401 Unauthorized` |
| внутренняя ошибка БД | `500 Internal Server Error` |

## Структура проекта

```
cmd/url-short/       точка входа: конфиг → логгер → БД → wire → запуск
internal/
  config/            загрузка конфига (cleanenv: yaml + env)
  domain/            сущности, доменные ошибки, интерфейс репозитория
  storage/           PostgreSQL, чистый SQL (pgx)
  service/           бизнес-логика: валидация, генерация алиаса
  handler/           HTTP-роутер, хендлеры, Basic Auth middleware
migrations/          SQL-миграции (golang-migrate)
tests/               юнит-тесты хендлеров
config/local.yaml    конфиг по умолчанию
```

## Тесты

```bash
go test ./...
```

Юнит-тесты покрывают хендлеры через фейковый репозиторий: успешные сценарии, битый JSON, конфликт алиаса, «не найден», авторизацию и внутренние ошибки БД. Basic Auth в тестах берёт креды из того же `.env`, что и сервис.