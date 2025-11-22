# Decisión: Tarea 2.3 Bloqueada - Validar Compilación Local

**Fecha:** 2025-11-21
**Tarea:** 2.3 - Validar compilación local exhaustiva
**Fase:** FASE 1
**Sprint:** SPRINT-2
**Razón del Bloqueo:** Go no disponible

---

## Comandos para FASE 2

Cuando Go 1.25 esté disponible:

```bash
cd /home/user/edugo-api-mobile

# 1. Limpiar cache
go clean -cache
go clean -testcache
go clean -modcache
go mod download

# 2. Verificar módulos
go mod verify

# 3. Compilar
go build -v ./...

# 4. Tests unitarios
go test -v -short ./...

# 5. Tests completos (con Docker)
go test -v ./...

# 6. Race detector
go test -race -short ./...

# 7. Linter
golangci-lint run --timeout=5m

# 8. Cobertura
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

## Criterios de Aceptación

- [ ] go mod verify exitoso
- [ ] Compilación exitosa
- [ ] Tests unitarios pasan
- [ ] Tests integración pasan
- [ ] Race detector pasa
- [ ] Cobertura >= 33%

---

**Estado:** STUB completado - Validación pendiente para FASE 2
