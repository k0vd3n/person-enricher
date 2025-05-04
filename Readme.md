# Person Enricher
[RU](#русская-версия) | [EN](#english-version)

## <a name="english-version"></a> English version

The Person Enricher application implements a RESTful HTTP server to manage person records. It enriches incoming requests with probable age, gender, and nationality by querying public APIs, stores enriched data in PostgreSQL, and exposes metrics for monitoring.


## Technology stack 

- **Programming Language:**  Go
- **Web Framework:**  Gorilla Mux
- **Database Management System:**  PostgreSQL
- **ORM:**  GORM
- **HTTP Client:**  net/http
- **Environment Variables:**  godotenv
- **Metrics Collection:**  Prometheus client_golang
- **Dashboard:**  Grafana (uses `person_enricher_dashboard.json`)
- **UUID Generation:**  (handled by GORM/model)
- **Swagger / OpenAPI:**  Swagger YAML/JSON in `docs/` (generated via swag or similar)
- **Testing & Mocks:** 
  - Unit testing: Go’s testing package
  - Assertions: testify/assert
  - Mock generation: gomock / mockgen

## Database 

A single `people` table is created via GORM auto-migrations. The schema corresponds to the `models.Person` struct and includes:

```sql
CREATE TABLE people (
  id           UUID PRIMARY KEY,
  name         VARCHAR NOT NULL,
  surname      VARCHAR NOT NULL,
  patronymic   VARCHAR,
  age          INT,
  gender       VARCHAR,
  nationality  VARCHAR
);
```

## HTTP API 

We expose the following endpoints under `/people`:
 
- **GET /people** 
List people with optional `filter`, `page` and `size` query parameters.
 
- **GET /people/{id}** 
Retrieve a single person by ID.
 
- **POST /people** 
Create a new person:

```json
{
  "name": "Dmitriy",
  "surname": "Ushakov",
  "patronymic": "Vasilevich"  // optional
}
```
 
- **PUT /people/{id}** 

Update an existing person (all fields required).
 
- **DELETE /people/{id}** 

Remove a person by ID.

All responses are JSON; errors return a standardized `{ "error": "message" }` format.

## External Data Enrichment 

On creation, we call three public APIs to enrich the record:

- **Age:**  [https://api.agify.io/?name={name}]()
- **Gender:**  [https://api.genderize.io/?name={name}]()
- **Nationality:**  [https://api.nationalize.io/?name={name}]()

Enriched data (age, gender, nationality) is merged into the `Person` model before saving.

## Logging 

Debug and info logs throughout the application (using Go’s `log.Printf`), covering:

- Startup and shutdown sequences
- Incoming requests and parsed parameters
- External API calls and responses
- Database operations

## Configuration 

All runtime configuration (DB credentials, HTTP port, metrics port, SSL mode) is loaded from a `.env` file via `github.com/joho/godotenv`, or directly from environment variables.

## Swagger / OpenAPI 

Swagger definitions generated into `docs/swagger.yaml` and `docs/swagger.json`. The `docs/docs.go` file exposes these via an embedded handler.

## Metrics & Monitoring 

Prometheus metrics are collected for:

- HTTP request durations (`http_request_duration_seconds{method,path}`)
- Service method durations (`service_method_duration_seconds{method}`)
- Repository method durations (`repo_method_duration_seconds{method}`)
- External API call durations (`enricher_request_duration_seconds{type}`)

Metrics are exposed on the `/metrics` endpoint (configurable port).

A Grafana dashboard is provided in the root:

```bash
person_enricher_dashboard.json
```

## Important files 

- **`cmd/main.go`** 
Application entry point: loads config, initializes DB, metrics, services, router, and starts HTTP & metrics servers with graceful shutdown.
- **`internal/repository/db.go`** 
Returns a GORM `*DB`, sets up connection pool, and auto-migrates the `Person` model.
- **`internal/handlers/router.go`** 
Configures Gorilla Mux routes and middleware for HTTP metrics.
- **`internal/externalapi/data_enricher.go`** 
Implements `EnrichPersonalData`: calls agify, genderize, nationalize.
- **`internal/metrics/metrics.go`** 
Defines all Prometheus metrics vectors.
- **`docs/swagger.yaml` / `swagger.json`** 
OpenAPI definitions for the REST API.
- **`Dockerfile`** 
Multi-stage build for a minimal production container image.

## Dockerfile 

```dockerfile
FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o person-enricher ./cmd

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/person-enricher .
EXPOSE 8080 8081
ENTRYPOINT ["./person-enricher"]
```

Build and run:

```bash
docker build -t person-enricher:latest .
docker run -d \
  -p 8080:8080 -p 8081:8081 \
  --env-file .env \
  person-enricher:latest
```

## Security 

- All DB operations use GORM’s parameterized queries.
- Secrets (DB credentials) are kept out of version control in `.env` or orchestration secrets.
- HTTPS termination and secret management recommended in production (e.g., Kubernetes Secrets).

## Examples 

Creating a person:

```bash
curl -X POST http://localhost:8080/v1/people \
  -H "Content-Type: application/json" \
  -d '{"name":"Dmitry","surname":"Ushakov"}'
```

Response:

```json
{
  "id":"...",
  "name":"Dmitry",
  "surname":"Ushakov",
  "patronymic":"",
  "age":35,
  "gender":"male",
  "nationality":"RU"
}
```

Listing with filter & pagination:

```bash
curl "http://localhost:8080/v1/people?filter=ushakov&page=1&size=5"
```

Updating a person:

```bash
curl -X PUT http://localhost:8080/v1/people/{id} \
  -H "Content-Type: application/json" \
  -d '{"name":"Dmitriy","surname":"Ushakov","patronymic":"","age":36,"gender":"male","nationality":"RU"}'
```

Deleting a person:


```bash
curl -X DELETE http://localhost:8080/v1/people/{id}
```

## <a name="русская-версия"></a> Русская Версия

Jump to [English version](#english-version)

Приложение Person Enricher реализует RESTful HTTP‑сервер для управления информацией о людях. При создании записи оно обогащает данные наиболее вероятным возрастом, полом и национальностью, сохраняет результат в PostgreSQL и экспортирует метрики Prometheus для мониторинга.

## Технологический стек 
 
- **Язык программирования:**  Go
- **Веб-фреймворк:**  Gorilla Mux
- **СУБД:**  PostgreSQL
- **ORM:**  GORM
- **HTTP‑клиент:**  net/http
- **Переменные окружения:**  godotenv
- **Метрики:**  Prometheus client_golang
- **Дашборд:**  Grafana (использует `person_enricher_dashboard.json`)
- **UUID Генератор:**  (обрабатывается GORM/model)
- **Swagger / OpenAPI:**  `docs/swagger.yaml`, `docs/swagger.json`
- **Тестирование и моки:** 
	- Unit тесты: Go testing пакет
	- Assertions: testify/assert
	- Mock generation: gomock / mockgen
## База данных 

Единая таблица `people` создаётся через авто‑миграцию GORM. Схема:

```sql
CREATE TABLE people (
  id           UUID PRIMARY KEY,
  name         VARCHAR NOT NULL,
  surname      VARCHAR NOT NULL,
  patronymic   VARCHAR,
  age          INT,
  gender       VARCHAR,
  nationality  VARCHAR
);
```

## HTTP API 

Маршруты:
 
- **GET /people**  — список с фильтром `filter`, пагинацией `page` и `size`.
- **GET /people/{id}**  — получить по ID.
- **POST /people**  — создать новую запись.

```json
{
  "name": "Dmitriy",
  "surname": "Ushakov",
  "patronymic": "Vasilevich"  // optional
}
```

- **PUT /people/{id}**  — обновить запись.
- **DELETE /people/{id}**  — удалить запись.

Ошибки возвращают JSON `{ "error": "сообщение" }`.

## Обогащение данных 


При создании записи вызываются публичные API:

2. Age: `https://api.agify.io/?name={name}`
3. Gender: `https://api.genderize.io/?name={name}`
4. Nationality: `https://api.nationalize.io/?name={name}`

Результаты сохраняются в таблице `people`.

## Логирование 

Используются `log.Printf` для debug и info:

- Параметры запуска и остановки.
- Параметры HTTP‑запросов.
- Результаты внешних вызовов и операций с БД.

## Конфигурация 
Все параметры (DB_HOST, DB_PORT, HTTP_PORT, METRICS_PORT и т.д.) загружаются из `.env` либо из переменных окружения.

## Swagger / OpenAPI 
Сгенерированные спецификации находятся в `docs/swagger.yaml` и `docs/swagger.json`. Встраиваются через `docs/docs.go`.
## Метрики и мониторинг 


Собираются метрики:
- Время HTTP‑запросов: `http_request_duration_seconds{method,path}`
- Время методов сервиса: `service_method_duration_seconds{method}`
- Время методов репозитория: `repo_method_duration_seconds{method}`
 - Время вызовов внешних API: `enricher_request_duration_seconds{type}`

Экспонируются на `/metrics` (порт настраивается). Есть готовый дашборд Grafana: `person_enricher_dashboard.json`.
## Важные файлы 

- **`cmd/main.go`**  — точка входа, конфигурация, запуск HTTP & метрик серверов.
- **`internal/repository/db.go`**  — настройка GORM и авто‑миграция.
- **`internal/handlers/router.go`**  — маршрутизация Gorilla Mux и middleware.
- **`internal/externalapi/data_enricher.go`**  — вызовы Agify, Genderize, Nationalize.
- **`internal/metrics/metrics.go`**  — определения метрик Prometheus.
- **`docs/swagger.yaml` / `swagger.json`**  — спецификация API.
- **`Dockerfile`**  — multi‑stage сборка образа.
## Dockerfile 

```dockerfile
FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o person-enricher ./cmd

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/person-enricher .
EXPOSE 8080 8081
ENTRYPOINT ["./person-enricher"]
```
## Безопасность 

 
- Параметризованные запросы через GORM (защита от SQL‑инъекций).
 
- Конфиденциальные данные в `.env` или оркестрационных секретах.

- HTTPS и управление секретами рекомендуются в продакшене.
## Примеры 
Создание записи:

```bash
curl -X POST http://localhost:8080/v1/people \
  -H "Content-Type: application/json" \
  -d '{"name":"Dmitry","surname":"Ushakov"}'
```

Ответ:

```json
{
  "id":"...",
  "name":"Dmitry",
  "surname":"Ushakov",
  "patronymic":"",
  "age":35,
  "gender":"male",
  "nationality":"RU"
}
```

Вывод всех записей с фильтрами и пагинацией:

```bash
curl "http://localhost:8080/v1/people?filter=ushakov&page=1&size=5"
```

Получение записи по id:

```bash
curl "http://localhost:8080/v1/people/{id}"
```

Обновление:

```bash
curl -X PUT http://localhost:8080/v1/people/{id} \
  -H "Content-Type: application/json" \
  -d '{"name":"Dmitriy","surname":"Ushakov","patronymic":"","age":36,"gender":"male","nationality":"RU"}'
```

Удаление:

```bash
curl -X DELETE http://localhost:8080/v1/people/{id}
```