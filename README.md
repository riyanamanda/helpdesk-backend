# Helpdesk REST API

A REST API for IT Helpdesk management built with Go and Echo v5, following Clean Architecture (Handler → Service → Repository). Designed for internal employee needs, this project might seem over-engineered for a helpdesk system — and it is. This was done intentionally to explore and master scalable software architecture.

**ERD:** [View on Eraser](https://app.eraser.io/workspace/MCKUzCCls92JCU5rpuew?origin=share)
**API:** [View on Postman](https://ranstack.postman.co/workspace/Personal~a161c18a-46cb-43a7-9866-eca9e1d1d19d/collection/10256898-3098bfb3-d3fd-4e97-8939-23523608d0ea?action=share&source=copy-link&creator=10256898)

## Stack

| Concern | Technology |
|---|---|
| Language | Go 1.26+ |
| Web framework | Echo v5 |
| Database | PostgreSQL + sqlx |
| Migrations | Goose |
| Cache | Redis |
| Object storage | MinIO |
| Message queue | RabbitMQ |
| Push notifications | Firebase (FCM) |
| Auth | JWT + Google OAuth |
| Logging | slog |
| Hot reload | Air |

## Prerequisites

- Go 1.26+
- PostgreSQL
- Redis
- MinIO
- RabbitMQ
- [Air](https://github.com/air-verse/air) — `go install github.com/air-verse/air@latest`
- [Goose](https://github.com/pressly/goose) — `go install github.com/pressly/goose/v3/cmd/goose@latest`
- [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports) — `go install golang.org/x/tools/cmd/goimports@latest`

## Getting Started

**1. Copy and fill in environment variables:**

```bash
cp .env.example .env
```

See [Environment Variables](#environment-variables) for details.

**2. Run migrations:**

```bash
make migrate-up
```

**3. Seed initial data:**

```bash
make seed
```

**4. Start the API server and worker:**

```bash
make run
```

This starts both the HTTP server and the background worker with hot reload via Air.

## Environment Variables

```env
# App
APP_NAME="Helpdesk App"
APP_HOST=localhost
APP_PORT=8080
CORS_ORIGINS=http://localhost:3000

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=helpdesk_db
DB_USERNAME=postgres
DB_PASSWORD=postgres
DB_SSLMODE=disable

# Goose (migrations)
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://postgres:postgres@localhost:5432/helpdesk_db
GOOSE_MIGRATION_DIR=./migrations

# Auth
JWT_SECRET=your-secret
JWT_EXP=24h
FIREBASE_PROJECT_ID=your-firebase-project-id

# MinIO (object storage)
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=helpdesk-dev
MINIO_USE_SSL=false

# Mail
MAIL_HOST=sandbox.smtp.mailtrap.io
MAIL_PORT=587
MAIL_USERNAME=
MAIL_PASSWORD=
MAIL_FROM=noreply@example.com
MAIL_SSL=false

# RabbitMQ
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
```

## Make Commands

| Command | Description |
|---|---|
| `make run` | Start API server + worker with hot reload |
| `make build` | Compile all packages |
| `make seed` | Seed the database with initial data |
| `make migrate-up` | Apply all pending migrations |
| `make migrate-down` | Roll back the last migration |
| `make migrate-create name=<name>` | Create a new migration file |
| `make test` | Run all tests |
| `make format` | Format code with goimports |
| `make tidy` | Tidy go.mod and go.sum |
| `make clean` | Remove build artifacts |

## Project Structure

```
.
├── cmd/
│   ├── api/        # HTTP server entry point
│   ├── seed/       # Database seeder
│   └── worker/     # Background job worker entry point
├── internal/
│   ├── platform/   # Infrastructure (DB, cache, storage, middleware, etc.)
│   ├── shared/     # Cross-domain utilities (errors, response, validation, etc.)
│   ├── auth/
│   ├── category/
│   ├── dashboard/
│   ├── division/
│   ├── feedback/
│   ├── mailer/
│   ├── notification/
│   ├── profile/
│   ├── ticket/
│   ├── user/
│   └── user_device/
└── migrations/     # SQL migration files
```
