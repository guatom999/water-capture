# 🌊 Water Level Monitor — Backend

Backend system for real-time water level monitoring, built with Go (Echo framework). Split into 3 services: **API**, **Cron**, and **Worker**.

## Architecture

Uses **Layered Architecture** with interface-based Dependency Injection at every layer.

```
backend/
├── cmd/
│   ├── api/          # HTTP API server (Echo)
│   ├── cron/         # Scheduled jobs (robfig/cron)
│   └── worker/       # Async task worker (asynq + Redis)
├── internal/
│   ├── config/       # App configuration (.env)
│   ├── database/     # Database connection + migrations
│   │   └── migrate/  # SQL migration scripts
│   ├── entities/     # Domain entities (Province, Location, WaterLevel)
│   ├── models/       # Request/Response DTOs
│   ├── handlers/     # HTTP handlers (controller layer)
│   ├── services/     # Business logic layer
│   ├── repositories/ # Data access layer (sqlx + PostgreSQL)
│   ├── middleware/    # Auth & security middleware
│   ├── jobs/         # Cron job definitions
│   ├── tasks/        # Asynq task types & handlers
│   ├── notifiers/    # Notification channels (LINE, Webhook)
│   └── utils/        # Shared utility functions
├── Dockerfile        # Multi-stage build (api/cron/worker)
├── docker-compose.yml
└── go.mod
```

### Dependency Flow

```
Handlers → Services → Repositories → PostgreSQL
                ↘ Tasks → Redis Queue → Worker → Notifiers (LINE)
```

## Tech Stack

| Technology | Usage |
|---|---|
| **Go 1.24** | Core language |
| **Echo v4** | HTTP framework |
| **sqlx + lib/pq** | PostgreSQL driver & query builder |
| **golang-jwt/jwt** | JWT authentication |
| **asynq + Redis** | Async task queue (notifications) |
| **robfig/cron** | Scheduled job runner |
| **testify + sqlmock** | Unit testing |

## API Endpoints

### Health Check

```
GET /heath
```

### Water Level

```
GET  /markers                              # Get all stations with latest water level
GET  /markers/detail?station_id={id}       # Get detailed water level data for a station
POST /add_station_location                 # Add new station from external API
```

### Authentication

```
POST /auth/register    # Register new user
POST /auth/login       # Login (returns access + refresh token)
POST /auth/refresh     # Refresh access token
POST /auth/logout      # Logout (revoke refresh token)
```

### Image Serving

```
GET /images/:filename    # Serve water level images
GET /images/health       # Image service health check
```

## Services

### 1. API Server (`cmd/api`)

Main HTTP server serving the REST API.

```bash
cd cmd/api && go run main.go
```

### 2. Cron Service (`cmd/cron`)

Fetches water level data from external APIs every 20 minutes. Enqueues LINE alerts when levels reach DANGER/WATCH status.

```bash
cd cmd/cron && go run main.go
```

### 3. Worker Service (`cmd/worker`)

Consumes tasks from Redis queue (asynq) and sends notifications via LINE Messaging API.

```bash
cd cmd/worker && go run main.go
```

## Setup & Running

### Prerequisites

- Go 1.24+
- PostgreSQL 16+
- Redis 7+

### Option 1: Docker Compose (Recommended)

```bash
docker compose up -d
```

This starts all services: PostgreSQL, Redis, API, Cron, and Worker.

### Option 2: Local Development

1. Create `.env` from `.env.example`

2. Start Database & Redis

```bash
docker compose -f docker-compose-db.yaml up -d
```

3. Run migrations

```bash
psql -h localhost -p 5503 -U postgres -d self_boardcast-db -f internal/database/migrate/migrate.sql
psql -h localhost -p 5503 -U postgres -d self_boardcast-db -f internal/database/migrate/auth.sql
```

4. Start services (each in a separate terminal)

```bash
cd cmd/api && go run main.go
cd cmd/cron && go run main.go
cd cmd/worker && go run main.go
```

## Environment Variables

| Variable | Description | Default |
|---|---|---|
| `APP_PORT` | API server port | `8080` |
| `DB_HOST` | PostgreSQL host | - |
| `DB_PORT` | PostgreSQL port | - |
| `DB_USERNAME` | PostgreSQL username | - |
| `DB_PASSWORD` | PostgreSQL password | - |
| `DB_NAME` | Database name | - |
| `DB_SSLMODE` | SSL mode | - |
| `JWT_SECRET` | JWT signing secret | - |
| `ACCESS_TOKEN_EXPIRY` | Access token expiry (minutes) | `15` |
| `REFRESH_TOKEN_EXPIRY` | Refresh token expiry (days) | `1` |
| `UPLOAD_DIR` | Image upload directory | - |
| `IMAGE_PROCESSING_DIR` | Image processing directory | - |
| `REDIS_ADDR` | Redis address | `localhost:6379` |
| `LINE_CHANNEL_ACCESS_TOKEN` | LINE Messaging API token | - |

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test ./internal/repositories/ ./internal/services/ -cover

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```
