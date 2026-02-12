# 🌊 Water Level Monitor

Real-time water level monitoring system for Bangkok and Pathum Thani — automated data collection, interactive map visualization, and LINE push notifications.

![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)
![React](https://img.shields.io/badge/React-19-61DAFB?logo=react&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white)

## ✨ Features

- 📊 **Real-time Dashboard** — Water level charts with bank-level reference lines
- 🗺️ **Interactive Map** — Leaflet map with GeoJSON province boundaries (Bangkok, Pathum Thani)
- ⏰ **Automated Data Collection** — Cron job fetches data every 20 minutes from external APIs
- 🔔 **LINE Notifications** — Automatic alerts when water levels reach DANGER/WATCH status
- 🔐 **JWT Authentication** — Register / Login / Refresh Token / Logout
- 🐳 **Dockerized** — Multi-stage build with Docker Compose

## 🏗️ Architecture

```
┌─────────────┐     ┌──────────────────────────────────────────┐
│   Frontend   │────▶│              API Server (Echo)            │
│  React+Vite  │     │  handlers → services → repositories      │
└─────────────┘     └───────┬──────────────────┬───────────────┘
                            │                  │
                    ┌───────▼──────┐    ┌──────▼──────┐
                    │  PostgreSQL  │    │    Redis     │
                    │  (Data)      │    │  (Queue)     │
                    └──────────────┘    └──────┬──────┘
                                               │
┌──────────────┐                       ┌───────▼──────┐
│  Cron Service │──── enqueue ────────▶│   Worker     │
│  (Scheduler) │                       │  (asynq)     │
└──────────────┘                       └──────┬──────┘
                                               │
                                       ┌───────▼──────┐
                                       │ LINE Notify  │
                                       └──────────────┘
```

## 📁 Project Structure

```
self-boardcast/
├── backend/           # Go backend (Echo + sqlx + asynq)
│   ├── cmd/
│   │   ├── api/       # REST API server
│   │   ├── cron/      # Scheduled data fetching
│   │   └── worker/    # Async notification worker
│   ├── internal/      # Business logic (layered architecture)
│   └── Dockerfile     # Multi-stage build
├── frontend/          # React + Vite + TypeScript + Tailwind
│   ├── src/
│   │   ├── components/  # Map, Chart, Header, Detail views
│   │   ├── pages/       # Login, Register
│   │   ├── services/    # API client
│   │   ├── stores/      # Zustand auth store
│   │   └── types/       # TypeScript types
│   └── package.json
└── docker-compose.yml
```

## 🚀 Quick Start

### Docker Compose (Recommended)

```bash
cd backend
docker compose up -d
```

### Manual Setup

```bash
# Backend (3 terminals)
cd backend/cmd/api && go run main.go
cd backend/cmd/cron && go run main.go
cd backend/cmd/worker && go run main.go

# Frontend
cd frontend && npm install && npm run dev
```

See [backend/README.md](backend/README.md) and [frontend/README.md](frontend/README.md) for detailed setup instructions.

## 🧪 Testing

```bash
cd backend && go test ./... -cover
```

## Tech Stack

| Layer | Technology |
|---|---|
| **Frontend** | React 19, Vite, TypeScript, Tailwind CSS, Leaflet, Recharts |
| **Backend** | Go 1.24, Echo v4, sqlx, JWT |
| **Queue** | Redis + asynq |
| **Database** | PostgreSQL 16 |
| **Notification** | LINE Messaging API |
| **Infrastructure** | Docker, Docker Compose |
