# Tareas del Sprint 06 - CI/CD

## Objetivo
Implementar pipeline completo de CI/CD con GitHub Actions para garantizar calidad automática.

---

## Tareas

### TASK-06-001: GitHub Actions Workflow
**Prioridad:** HIGH  
**Estimación:** 3h

#### Implementación

Archivo: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/.github/workflows/ci.yml`

```yaml
name: CI Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

jobs:
  test:
    name: Tests
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:15-alpine
        env:
          POSTGRES_USER: postgres
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: edugo_test
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
      
      mongodb:
        image: mongo:7-alpine
        ports:
          - 27017:27017
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Install dependencies
        run: go mod download
      
      - name: Run migrations
        run: |
          psql -h localhost -U postgres -d edugo_test < scripts/postgresql/06_assessments.sql
        env:
          PGPASSWORD: postgres
      
      - name: Run tests
        run: go test ./... -v -cover -coverprofile=coverage.out
        env:
          DB_HOST: localhost
          DB_PORT: 5432
          DB_USER: postgres
          DB_PASSWORD: postgres
          DB_NAME: edugo_test
          MONGO_URI: mongodb://localhost:27017
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./coverage.out
  
  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest
  
  build:
    name: Build
    runs-on: ubuntu-latest
    needs: [test, lint]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Build
        run: go build -o bin/api-mobile ./cmd/api
```

#### Criterios
- [ ] Tests automáticos en cada push
- [ ] PostgreSQL y MongoDB en services
- [ ] Coverage reportado a Codecov
- [ ] Linting automático
- [ ] Build exitoso

---

### TASK-06-002: Dockerfile
**Prioridad:** MEDIUM  
**Estimación:** 2h

#### Implementación

Archivo: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/Dockerfile`

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin/api-mobile ./cmd/api

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/bin/api-mobile .
COPY --from=builder /app/config ./config

EXPOSE 8080

CMD ["./api-mobile"]
```

#### Comandos
```bash
# Build
docker build -t edugo-api-mobile:latest .

# Run
docker run -p 8080:8080 edugo-api-mobile:latest
```

---

### TASK-06-003: Documentación README
**Prioridad:** MEDIUM  
**Estimación:** 2h

#### Actualizar README con:
- Instalación
- Configuración
- Migraciones
- Tests
- Deployment

---

## Resumen

**Tareas:** 5  
**Estimación:** 12 horas  
**Entregables:** Pipeline CI/CD completo

---

**Sprint:** 06/06
