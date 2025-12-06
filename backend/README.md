# Self Boardcast Backend API

## โครงสร้าง Hexagonal Architecture

```
backend/
├── cmd/
│   └── api/
│       └── main.go                 # Entry point สำหรับ API server
├── internal/
│   ├── core/                       # Core Business Logic Layer
│   │   ├── domain/                 # Domain entities
│   │   │   └── water_level.go
│   │   ├── ports/                  # Interfaces (Ports)
│   │   │   ├── input/              # Input ports (use cases)
│   │   │   │   └── water_level_service.go
│   │   │   └── output/             # Output ports (repositories)
│   │   │       └── water_level_repository.go
│   │   └── services/               # Business logic implementation
│   │       └── water_level_service.go
│   └── adapters/                   # Adapters Layer
│       ├── input/                  # Input adapters
│       │   └── http/               # HTTP handlers
│       │       ├── handler.go
│       │       └── router.go
│       └── output/                 # Output adapters
│           └── persistence/        # Database/storage implementations
│               └── water_level_repository.go
├── pkg/                            # Public packages
│   └── config/
│       └── config.go
├── go.mod
├── go.sum
└── main.go                         # Alternative entry point
```

## การติดตั้ง

```bash
cd backend
go mod tidy
```

## การรัน

### วิธีที่ 1: รันจาก cmd/api

```bash
go run cmd/api/main.go
```

### วิธีที่ 2: รันจาก main.go

```bash
go run main.go
```

### วิธีที่ 3: Build และรัน

```bash
go build -o server cmd/api/main.go
./server
```

## API Endpoints

### Health Check

```bash
GET http://localhost:8080/health
```

### Get Latest Water Level

```bash
GET http://localhost:8080/api/v1/water-levels/latest
```

### Get Water Level History

```bash
GET http://localhost:8080/api/v1/water-levels/history?limit=10
```

### Process Image

```bash
POST http://localhost:8080/api/v1/water-levels/process
Content-Type: application/json

{
  "image_url": "http://example.com/image.png"
}
```

## ทดสอบ API

```bash
# Health check
curl http://localhost:8080/health

# Get latest water level
curl http://localhost:8080/api/v1/water-levels/latest

# Get history
curl http://localhost:8080/api/v1/water-levels/history?limit=5

# Process image
curl -X POST http://localhost:8080/api/v1/water-levels/process \
  -H "Content-Type: application/json" \
  -d '{"image_url": "http://example.com/water.jpg"}'
```

## Hexagonal Architecture Layers

### 1. Core Domain Layer

- **Domain entities**: โมเดลหลักของ business (WaterLevel)
- **Pure business logic**: ไม่มี dependencies กับ framework หรือ infrastructure

### 2. Ports Layer

- **Input Ports**: Interfaces ที่กำหนด use cases (WaterLevelService)
- **Output Ports**: Interfaces สำหรับ external dependencies (WaterLevelRepository)

### 3. Services Layer

- **Business logic implementation**: Implement input ports
- **Use output ports**: เรียกใช้ repositories ผ่าน interfaces

### 4. Adapters Layer

- **Input Adapters**: HTTP handlers, CLI, gRPC, etc.
- **Output Adapters**: Database, external APIs, file system, etc.

## ข้อดีของ Hexagonal Architecture

✅ **Testability**: ง่ายต่อการ mock interfaces และเขียน unit tests
✅ **Flexibility**: เปลี่ยน database หรือ framework ได้ง่าย
✅ **Maintainability**: แยก business logic ออกจาก infrastructure
✅ **Scalability**: เพิ่ม adapters ใหม่ได้โดยไม่กระทบ core logic
✅ **Clean Dependencies**: Dependencies ชี้เข้าหา core เสมอ

## ขั้นตอนต่อไป

1. เพิ่ม PostgreSQL adapter แทน in-memory repository
2. เพิ่ม HTTP client เพื่อเรียก Python service
3. เพิ่ม logging และ monitoring
4. เพิ่ม authentication และ authorization
5. เขียน unit tests และ integration tests
