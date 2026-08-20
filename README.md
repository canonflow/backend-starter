# Backend Starter

A production-minded Go backend starter built on [Fiber v3](https://gofiber.io/), [GORM](https://gorm.io/), and Redis. It ships with a clean, layered architecture (handler → usecase → repository), JWT cookie-based authentication, request rate limiting, structured logging, and a database migration toolchain — so you can start building features instead of wiring boilerplate.

> This repository is an enhanced fork of [canonflow/canonflow-go-backend-template](https://github.com/canonflow/canonflow-go-backend-template). See [Enhancements Over the Template](#enhancements-over-the-template) for what changed.

## Features

- **Layered architecture** — clear separation between delivery (routes), handlers, usecases, repositories, and models.
- **Pluggable persistence** — repository factory + driver pattern makes it straightforward to add databases beyond MySQL.
- **JWT authentication** — tokens issued as HTTP-only cookies, verified by auth middleware.
- **Rate limiting** — per-route sliding-window throttling backed by Redis.
- **Structured logging** — request-scoped logging with [zap](https://github.com/uber-go/zap).
- **Config from env** — type-safe config access with sensible defaults.
- **Database migrations** — versioned SQL migrations via [golang-migrate](https://github.com/golang-migrate/migrate) with `make` targets.
- **Consistent API envelope** — standard success / error / pagination response shapes.
- **High-performance JSON** — [sonic](https://github.com/bytedance/sonic) wired in as Fiber's JSON encoder and decoder.
- **Live reload** — [Air](https://github.com/air-verse/air) config included for local development.
- **Graceful shutdown** — clean termination on `SIGINT` / `SIGTERM`, with Fiber prefork enabled.

## Tech Stack

| Concern        | Technology                          |
| -------------- | ----------------------------------- |
| Language       | Go 1.26                             |
| Web framework  | Fiber v3                            |
| ORM            | GORM (MySQL driver)                |
| Cache / store  | Redis (go-redis)                   |
| Auth           | JWT (golang-jwt v5)                |
| Validation     | go-playground/validator v10        |
| Logging        | Uber zap                           |
| Migrations     | golang-migrate v4                  |
| JSON codec     | bytedance/sonic (Fiber encoder/decoder) |

## Project Structure

```
.
├── cmd
│   ├── main.go              # Application entrypoint
│   ├── migration/           # Migration CLI (up/down/version/force/drop)
│   └── scripts/             # Utility scripts (APP_KEY generator)
├── internal
│   ├── app/                 # Bootstrap, Fiber, GORM, Redis, validator wiring
│   │   └── scope/           # Reusable GORM query scopes (pagination, soft-delete)
│   ├── config/              # Type-safe env config loader
│   ├── core/                # Logger and Fiber context helpers
│   ├── dto/                 # Request/response data transfer objects
│   ├── http
│   │   ├── delivery/        # Route registration
│   │   ├── handler/         # HTTP handlers
│   │   └── middleware/      # Auth, rate limiter, request logger
│   ├── model/               # GORM models
│   ├── repository/          # Data access (driver + factory pattern)
│   └── usecase/             # Business logic
├── migrations/mysql/        # Versioned SQL migrations
├── pkg/                     # Shared packages (jwt, password, response, helpers)
├── .air.toml                # Live-reload config
├── .env.example             # Environment variable template
└── makefile                 # Migration & tooling shortcuts
```

## Getting Started

### Prerequisites

- Go 1.26+
- MySQL
- Redis
- [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) (only needed to create new migration files via `make migrate-create`)
- [Air](https://github.com/air-verse/air) (optional, for live reload)

### 1. Clone and install dependencies

```bash
git clone <your-repo-url>
cd backend-starter
go mod download
```

### 2. Configure environment

Copy the example env file and fill in your values:

```bash
cp .env.example .env
```

Then generate a secure `APP_KEY` (used to sign JWTs):

```bash
make key-generate
```

### 3. Run migrations

```bash
make migrate-up
```

### 4. Start the server

For development with live reload:

```bash
air
```

Or run directly:

```bash
go run cmd/main.go
```

The server listens on `APP_HOST:APP_PORT` (default `localhost:8000`). In `development`, all registered routes are printed on startup.

## Configuration

All configuration is loaded from `.env`. Key variables:

| Variable                  | Description                                  | Default              |
| ------------------------- | -------------------------------------------- | -------------------- |
| `APP_ENV`                 | `development` or `production`                | `development`        |
| `APP_HOST` / `APP_PORT`   | Bind host and port                           | `localhost` / `8000` |
| `APP_KEY`                 | Secret used to sign JWTs (`make key-generate`) | —                  |
| `BCRYPT_COST`             | bcrypt hashing cost                          | `10`                 |
| `DB_DRIVER`               | Database driver                              | `mysql`              |
| `DB_HOST` / `DB_PORT`     | Database host and port                       | `localhost` / `3306` |
| `DB_NAME`                 | Database name                                | —                    |
| `DB_USERNAME` / `DB_PASSWORD` | Database credentials                     | —                    |
| `DB_IDLE` / `DB_MAX` / `DB_LIFETIME` | Connection pool tuning            | `10` / `100` / `300` |
| `LOG_LEVEL`               | zap log level                                | `debug`              |
| `REDIS_HOST` / `REDIS_PORT` | Redis connection                           | `localhost`          |
| `REDIS_USERNAME` / `REDIST_PASSWORD` | Redis credentials                 | —                    |
| `CORS_ALLOW_*`            | CORS origins, headers, methods, credentials  | see `.env.example`   |
| `JWT_DURATION_IN_MINUTES` | Token lifetime in minutes                    | `3600`               |
| `JWT_SAME_SITE`           | Cookie SameSite policy                       | `Lax`                |
| `JWT_PATH` / `JWT_DOMAIN` | Cookie path and domain                       | `/` / ``             |
| `JWT_SECURE`              | Send cookie over HTTPS only                  | `false`              |
| `JWT_HTTP_ONLY`           | Mark auth cookie HTTP-only                   | `true`               |

## API

All routes are prefixed with `/api`. The user module is versioned under `/v1/user`.

| Method | Endpoint                | Auth | Description                          | Rate limit      |
| ------ | ----------------------- | ---- | ------------------------------------ | --------------- |
| POST   | `/api/v1/user/signup`   | No   | Register a new user                  | 5 / min         |
| POST   | `/api/v1/user/signin`   | No   | Authenticate and receive a JWT cookie | 5 / min        |
| GET    | `/api/v1/user/me`       | Yes  | Return the current authenticated user | 5 / min        |
| POST   | `/api/v1/user/signout`  | Yes  | Clear the auth cookie                | 2 / min         |

Authenticated routes read the JWT from an HTTP-only cookie. On sign-in the token is set automatically; subsequent requests are validated by the auth middleware.

### Request examples

Sign up:

```bash
curl -X POST http://localhost:8000/api/v1/user/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"secret"}'
```

Sign in (stores the auth cookie in `cookies.txt`):

```bash
curl -X POST http://localhost:8000/api/v1/user/signin \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{"email":"user@example.com","password":"secret"}'
```

Fetch the current user (reuses the cookie):

```bash
curl http://localhost:8000/api/v1/user/me -b cookies.txt
```

### Response format

Every response uses a consistent envelope.

Success:

```json
{
  "success": true,
  "data": { "id": 1, "email": "user@example.com", "created_at": "..." },
  "meta": null,
  "errors": null
}
```

Error:

```json
{
  "success": false,
  "data": null,
  "meta": null,
  "errors": [
    { "code": "BAD_REQUEST", "message": "Email is already taken", "field": "email" }
  ]
}
```

## Migrations

Migrations are versioned SQL files under `migrations/<driver>/` and run through the migration CLI, exposed via `make`:

| Command                                   | Description                                    |
| ----------------------------------------- | ---------------------------------------------- |
| `make migrate-create NAME=<name>`         | Create a new up/down migration pair            |
| `make migrate-up`                         | Apply all pending migrations                   |
| `make migrate-down`                       | Revert applied migrations                      |
| `make migrate-version`                    | Show the current migration version             |
| `make migrate-force VERSION=<version>`    | Force set a version (recover from dirty state) |
| `make migrate-drop`                       | Drop everything (irreversible, prompts twice)  |

`migrate-force` and `migrate-drop` are destructive and require interactive confirmation.

## Architecture Notes

- **Bootstrap** (`internal/app/bootstrap.go`) wires repositories, usecases, handlers, middleware, and routes. Add new modules here.
- **Repository pattern** — each domain exposes an interface plus per-driver implementations selected through a factory. Swap or add a database by implementing the interface and registering it in the factory/driver switch.
- **Config** — access values with `config.Get[T](key)` or `config.GetOrDefault[T](key, def)` for type-safe reads.
- **Rate limiting** — the `Throttle` middleware uses a Redis-backed sliding window, keyed per route prefix and client IP.
- **JSON codec** — Fiber is configured with `sonic.Marshal` / `sonic.Unmarshal` (`internal/app/fiber.go`), replacing the standard `encoding/json` for faster request/response (de)serialization.

## Enhancements Over the Template

This project builds on [canonflow/canonflow-go-backend-template](https://github.com/canonflow/canonflow-go-backend-template) with the following changes and additions:

- **Fiber instead of Gin** — the HTTP layer is built on [Fiber v3](https://gofiber.io/) rather than the template's [Gin](https://gin-gonic.com/) framework.
- **Sonic JSON codec** — swapped the default `encoding/json` for [bytedance/sonic](https://github.com/bytedance/sonic) as Fiber's `JSONEncoder` and `JSONDecoder` for faster serialization.
- **Redis-backed rate limiting** — per-route sliding-window throttling (`Throttle` middleware) using Redis as shared storage, applied to the auth endpoints.
- **JWT cookie authentication** — tokens issued as HTTP-only cookies with configurable SameSite, path, domain, secure, and duration settings, plus auth middleware that surfaces user context to handlers.
- **Migration CLI + make targets** — versioned migrations driven by golang-migrate with `up`, `down [steps]`, `version`, `force`, and guarded `drop` commands.
- **Type-safe config layer** — generic `config.Get[T]` / `config.GetOrDefault[T]` accessors over env variables with sane defaults.
- **Structured, request-scoped logging** — zap logger propagated through context via a request logger middleware.
- **Consistent response envelope** — standardized success/error/pagination shapes and a normalized error handler.
- **Graceful shutdown & prefork** — signal-aware shutdown with a timeout, and Fiber prefork enabled for multi-process serving.

> Note: feature specifics reflect this repository's current state; refer to the upstream template for its baseline.
