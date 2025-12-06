# Self Boardcast Backend API

## Simple Clean Architecture

```
backend/
├── internal/
│   ├── models/              # Data models
│   │   └── water_level.go
│   ├── repositories/        # Data access layer
│   │   ├── water_level_repository.go
│   │   └── inmemory_repository.go
│   ├── services/            # Business logic
│   │   └── water_level_service.go
│   ├── handlers/            # HTTP handlers
│   │   └── water_level_handler.go
│   ├── server/              # Server setup
│   │   └── server.go
│   ├── utils/               # Utilities
│   │   ├── errors.go
│   │   └── time.go
│   └── db/                  # Database config
│       └── config.go
├── main.go
└── go.mod
```

## Architecture Layers

### 1. Models

- Data structures
- Business entities

### 2. Repositories

- Data access interfaces
- Database implementations

### 3. Services

- Business logic
- Use repository interfaces

### 4. Handlers

- HTTP request/response
- Call services

### 5. Server

- Route setup
- Middleware
- Server configuration

### 6. Utils & DB

- Helper functions
- Database configuration

## ติดตั้งและรัน

```bash
cd backend
go mod tidy
go run main.go
```

## API Endpoints

```bash
# Health check
GET http://localhost:8080/health

# Get latest water level
GET http://localhost:8080/api/v1/water-levels/latest

# Get history
GET http://localhost:8080/api/v1/water-levels/history?limit=10

# Get by ID
GET http://localhost:8080/api/v1/water-levels/:id

# Process image
POST http://localhost:8080/api/v1/water-levels/process
{
  "image_url": "http://example.com/image.png"
}
```

## ทดสอบ

```bash
# Health
curl http://localhost:8080/health

# Process image
curl -X POST http://localhost:8080/api/v1/water-levels/process \
  -H "Content-Type: application/json" \
  -d '{"image_url": "http://example.com/water.jpg"}'

# Get latest
curl http://localhost:8080/api/v1/water-levels/latest

# Get history
curl http://localhost:8080/api/v1/water-levels/history?limit=5
```

## Dependency Flow

```
Handler → Service → Repository
   ↓         ↓          ↓
 HTTP    Business    Database
Layer     Logic      Access
```

## ข้อดี

✅ **ง่ายต่อการเข้าใจ** - โครงสร้างชัดเจน ไม่ซับซ้อน
✅ **แยก concerns** - แต่ละ layer ทำหน้าที่ของตัวเอง
✅ **ทดสอบได้** - Mock repository ได้ง่าย
✅ **ขยายได้** - เพิ่ม features ใหม่ได้สะดวก
✅ **Maintainable** - แก้ไขง่าย ไม่กระทบส่วนอื่น

## ขั้นตอนต่อไป

- เพิ่ม PostgreSQL repository
- เพิ่ม Python service client
- เพิ่ม authentication
- เพิ่ม logging
- เขียน tests
