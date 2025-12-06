# 📋 TASKS Sprint-00: Integrar Infrastructure (ACTUALIZADO)

**Fecha:** 16 de Noviembre, 2025  
**Versión:** 2.0.0  
**Duración Estimada:** 3-4 horas  
**Prioridad:** CRÍTICA

---

## 🎯 Objetivo

Modernizar edugo-api-mobile para usar:
- ✅ `edugo-infrastructure` v0.5.0 (módulos por separado)
- ✅ `edugo-shared` v0.7.0 (últimas versiones)
- ✅ Eliminar código deprecated
- ✅ Aprovechar nuevas funcionalidades

---

## 📦 FASE 1: Actualizar Dependencias (30 min)

### TASK-001: Agregar Módulos de Infrastructure

```bash
# Módulos de infrastructure (nuevos)
go get github.com/EduGoGroup/edugo-infrastructure/postgres@v0.5.0
go get github.com/EduGoGroup/edugo-infrastructure/mongodb@v0.5.0
go get github.com/EduGoGroup/edugo-infrastructure/messaging@v0.5.0
go get github.com/EduGoGroup/edugo-infrastructure/database@v0.1.1
```

**Validación:**
```bash
grep "edugo-infrastructure" go.mod | wc -l
# Debe mostrar: 4
```

---

### TASK-002: Actualizar Módulos de Shared

```bash
# Actualizar módulos existentes
go get github.com/EduGoGroup/edugo-shared/auth@v0.7.0
go get github.com/EduGoGroup/edugo-shared/middleware/gin@v0.7.0

# Agregar módulos nuevos
go get github.com/EduGoGroup/edugo-shared/database/postgres@v0.7.0
go get github.com/EduGoGroup/edugo-shared/database/mongodb@v0.7.0
go get github.com/EduGoGroup/edugo-shared/config@v0.7.0
```

**Validación:**
```bash
go list -m all | grep "edugo-shared"
# Verificar versiones v0.7.0 o superiores
```

---

### TASK-003: Limpiar Dependencias

```bash
go mod tidy
go mod verify
```

**Validación:**
```bash
go build ./...
# Debe compilar sin errores
```

**Duración:** 30 minutos  
**Output:** `go.mod` y `go.sum` actualizados

---

## 🗑️ FASE 2: Eliminar Código Deprecated (1 hora)

### TASK-004: Analizar Migraciones Locales vs Infrastructure

```bash
# Ver migraciones en infrastructure
cd /path/to/edugo-infrastructure/postgres/migrations
ls -la *.up.sql

# Comparar con locales
cd /path/to/edugo-api-mobile/scripts/postgresql
ls -la *.sql
```

**Acción:**
1. Crear tabla comparativa (archivo, líneas, existe en infrastructure)
2. Identificar scripts 100% duplicados
3. Identificar scripts específicos del proyecto

**Output:** `MIGRACIONES_COMPARACION.md`

---

### TASK-005: Eliminar Scripts SQL Duplicados

Basado en análisis de TASK-004:

```bash
cd scripts/postgresql

# Eliminar SOLO si están en infrastructure
rm -f 01_create_schema.sql          # ✅ Existe en 001_create_users.up.sql
rm -f 02_seed_data.sql               # ✅ Existe en seeds/
rm -f 03_refresh_tokens.sql          # ✅ Parte de 001_create_users.up.sql
rm -f 04_material_versions.sql       # ✅ Parte de 005_create_materials.up.sql
rm -f 05_indexes_materials.sql       # ✅ Parte de 005_create_materials.up.sql

# MANTENER si son específicos de api-mobile
# 04_login_attempts.sql               # ⚠️ REVISAR
# 05_user_progress_upsert.sql         # ⚠️ REVISAR
```

**Validación:**
```bash
ls scripts/postgresql/*.sql
# Solo deben quedar scripts específicos (si los hay)
```

**Duración:** 20 minutos  
**Output:** Scripts eliminados, carpeta limpia

---

### TASK-006: Eliminar Connectors Custom de Database

```bash
cd internal/infrastructure/database

# Respaldar (por si acaso)
git mv postgres.go postgres.go.deprecated
git mv postgres_test.go postgres_test.go.deprecated
git mv mongodb.go mongodb.go.deprecated
git mv mongodb_test.go mongodb_test.go.deprecated

# Commit
git commit -m "chore: deprecar connectors custom de database

- Usar edugo-shared/database/postgres
- Usar edugo-shared/database/mongodb
- Eliminar implementación duplicada"
```

**Validación:**
```bash
ls internal/infrastructure/database/
# Solo deben quedar archivos .deprecated
```

**Duración:** 10 minutos  
**Output:** Archivos deprecated

---

### TASK-007: Actualizar Importaciones en el Proyecto

Buscar y reemplazar imports:

```bash
# Buscar usos del connector custom
grep -r "internal/infrastructure/database" internal/ --include="*.go"

# Reemplazar con:
# import "github.com/EduGoGroup/edugo-shared/database/postgres"
# import "github.com/EduGoGroup/edugo-shared/database/mongodb"
```

**Archivos a actualizar:**
- `internal/container/container.go` (DI)
- `cmd/main.go` (inicialización)
- Tests de integración que usen DB

**Ejemplo de cambio:**

```go
// ❌ ANTES
import "github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/database"

func InitDB(ctx context.Context, cfg *config.Config, log logger.Logger) (*sql.DB, error) {
    return database.InitPostgreSQL(ctx, cfg, log)
}

// ✅ DESPUÉS
import "github.com/EduGoGroup/edugo-shared/database/postgres"

func InitDB(ctx context.Context, cfg *config.Config, log logger.Logger) (*sql.DB, error) {
    return postgres.NewConnection(ctx, postgres.Config{
        Host:         cfg.Database.Postgres.Host,
        Port:         cfg.Database.Postgres.Port,
        Database:     cfg.Database.Postgres.Database,
        User:         cfg.Database.Postgres.User,
        Password:     cfg.Database.Postgres.Password,
        SSLMode:      cfg.Database.Postgres.SSLMode,
        MaxOpenConns: cfg.Database.Postgres.MaxConnections,
    })
}
```

**Validación:**
```bash
go build ./...
# Debe compilar sin errores
```

**Duración:** 30 minutos  
**Output:** Imports actualizados, código compilando

---

## ✨ FASE 3: Integrar Nuevas Funcionalidades (1.5 horas)

### TASK-008: Integrar Validador de Eventos

Crear archivo: `internal/infrastructure/messaging/validator.go`

```go
package messaging

import (
    "fmt"

    "github.com/EduGoGroup/edugo-infrastructure/messaging"
)

var validator *messaging.Validator

// InitValidator inicializa el validador de eventos una sola vez
func InitValidator() error {
    var err error
    validator, err = messaging.NewValidator()
    if err != nil {
        return fmt.Errorf("error inicializando validador: %w", err)
    }
    return nil
}

// ValidateEvent valida un evento contra su schema
func ValidateEvent(event interface{}, schemaName string) error {
    if validator == nil {
        return fmt.Errorf("validador no inicializado")
    }
    return validator.Validate(event, schemaName)
}
```

**Usar en publishers:**

```go
// En internal/infrastructure/messaging/publisher.go

func (p *Publisher) PublishMaterialUploaded(ctx context.Context, event MaterialUploadedEvent) error {
    // Validar antes de publicar
    if err := ValidateEvent(event, "material-uploaded-v1"); err != nil {
        return fmt.Errorf("evento inválido: %w", err)
    }

    // Publicar...
    return p.channel.PublishWithContext(ctx, ...)
}
```

**Validación:**
```bash
# Test con evento inválido
go test ./internal/infrastructure/messaging -run TestValidateEvent
```

**Duración:** 30 minutos  
**Output:** Validación de eventos funcionando

---

### TASK-009: Configurar Migraciones de Infrastructure

Actualizar `README.md` del proyecto:

```markdown
## 🗄️ Setup Base de Datos

### Opción 1: Usar infrastructure (RECOMENDADO)

```bash
# Clonar infrastructure si no existe
cd /path/to/edugo-infrastructure

# Levantar servicios
make dev-up-messaging  # PostgreSQL + MongoDB + RabbitMQ

# Ejecutar migraciones
cd postgres && make migrate-up
cd ../mongodb && make migrate-up
```

### Opción 2: Docker Compose local (legacy)

Solo si NO tienes infrastructure clonado:

```bash
docker-compose up -d postgres mongodb rabbitmq
```

**NOTA:** Las migraciones están en `edugo-infrastructure`, NO en este proyecto.
```

**Validación:**
```bash
# Verificar que README menciona infrastructure
grep "edugo-infrastructure" README.md
```

**Duración:** 15 minutos  
**Output:** README actualizado

---

### TASK-010: Actualizar Tests con Infrastructure Database

Actualizar tests de integración para usar testcontainers con migraciones:

```go
// ❌ ANTES
func TestUserRepository_Integration(t *testing.T) {
    db := setupTestDB(t) // Custom setup
    defer db.Close()

    // ... tests
}

// ✅ DESPUÉS
import "github.com/EduGoGroup/edugo-infrastructure/database"

func TestUserRepository_Integration(t *testing.T) {
    db := database.NewTestPostgres(t,
        database.WithMigrations(), // Aplica migraciones de infrastructure
    )

    // ... tests
    // Cleanup automático
}
```

**Archivos a actualizar:**
- `internal/infrastructure/persistence/*_test.go`
- `internal/application/services/*_integration_test.go`

**Validación:**
```bash
go test ./internal/infrastructure/persistence/... -tags=integration
# Todos los tests deben pasar con migraciones reales
```

**Duración:** 45 minutos  
**Output:** Tests modernizados

---

## ✅ FASE 4: Validación y Documentación (30 min)

### TASK-011: Ejecutar Tests Completos

```bash
# Tests unitarios
go test ./... -short

# Tests de integración
go test ./... -tags=integration

# Coverage
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

**Validación:**
- ✅ Todos los tests pasan
- ✅ Coverage >= 80%
- ✅ Sin warnings de deprecated

**Duración:** 10 minutos

---

### TASK-012: Verificar Build y Lint

```bash
# Build
go build ./...

# Lint (si existe golangci-lint)
golangci-lint run

# Vet
go vet ./...
```

**Validación:**
- ✅ Build exitoso
- ✅ 0 errores de lint
- ✅ 0 errores de vet

**Duración:** 10 minutos

---

### TASK-013: Actualizar Documentación del Sprint

Crear archivo: `EXECUTION_REPORT.md`

```markdown
# Execution Report - Sprint-00

## Cambios Realizados

### Dependencias Agregadas
- edugo-infrastructure/postgres v0.5.0
- edugo-infrastructure/mongodb v0.5.0
- edugo-infrastructure/messaging v0.5.0
- edugo-infrastructure/database v0.1.1
- edugo-shared/database/postgres v0.7.0
- edugo-shared/database/mongodb v0.7.0
- edugo-shared/config v0.7.0

### Dependencias Actualizadas
- edugo-shared/auth: v0.3.3 → v0.7.0
- edugo-shared/middleware/gin: v0.3.3 → v0.7.0

### Código Eliminado
- scripts/postgresql/*.sql (5 archivos, ~300 líneas)
- internal/infrastructure/database/postgres.go (~100 líneas)
- internal/infrastructure/database/mongodb.go (~100 líneas)
- internal/infrastructure/database/*_test.go (~300 líneas)

Total: ~800 líneas eliminadas

### Código Agregado
- internal/infrastructure/messaging/validator.go (~50 líneas)
- Tests modernizados (~150 líneas)

Total: ~200 líneas agregadas

### Resultado Neto
- **-600 líneas** de código
- **100%** de migraciones centralizadas
- **100%** de eventos validados
- **+7 módulos** de dependencias

## Métricas

- Tests: PASS (100%)
- Coverage: XX%
- Build: OK
- Lint: 0 errores

## Próximos Pasos

- Sprint-01: Implementar Schema de Evaluaciones
```

**Duración:** 10 minutos  
**Output:** Reporte de ejecución

---

## 📊 CHECKLIST DE COMPLETACIÓN

### Dependencias
- [ ] `go.mod` contiene `edugo-infrastructure/postgres@v0.5.0`
- [ ] `go.mod` contiene `edugo-infrastructure/mongodb@v0.5.0`
- [ ] `go.mod` contiene `edugo-infrastructure/messaging@v0.5.0`
- [ ] `go.mod` contiene `edugo-infrastructure/database@v0.1.1`
- [ ] `go.mod` contiene `edugo-shared/auth@v0.7.0`
- [ ] `go.mod` contiene `edugo-shared/middleware/gin@v0.7.0`
- [ ] `go.mod` contiene `edugo-shared/database/postgres@v0.7.0`
- [ ] `go.mod` contiene `edugo-shared/database/mongodb@v0.7.0`
- [ ] `go.mod` contiene `edugo-shared/config@v0.7.0`

### Código Eliminado
- [ ] `scripts/postgresql/` contiene solo scripts específicos (o está vacío)
- [ ] `internal/infrastructure/database/postgres.go` deprecated o eliminado
- [ ] `internal/infrastructure/database/mongodb.go` deprecated o eliminado
- [ ] Imports de connectors custom reemplazados

### Código Agregado
- [ ] Validador de eventos integrado
- [ ] Tests usan `database.NewTestPostgres()`
- [ ] README actualizado con instrucciones de infrastructure

### Validaciones
- [ ] `go build ./...` compila sin errores
- [ ] `go test ./... -short` todos los tests pasan
- [ ] `go test ./... -tags=integration` todos los tests pasan
- [ ] `golangci-lint run` sin errores (si aplica)
- [ ] Coverage >= 80%

### Documentación
- [ ] `EXECUTION_REPORT.md` creado
- [ ] `ANALISIS_MODERNIZACION.md` revisado
- [ ] `README.md` actualizado

---

## ⏱️ TIEMPO TOTAL ESTIMADO

| Fase | Duración |
|------|----------|
| Fase 1: Actualizar Dependencias | 30 min |
| Fase 2: Eliminar Código Deprecated | 1 hora |
| Fase 3: Integrar Nuevas Funcionalidades | 1.5 horas |
| Fase 4: Validación y Documentación | 30 min |
| **TOTAL** | **3-4 horas** |

---

## 🚨 PUNTOS DE ATENCIÓN

1. **Backup antes de eliminar:** Crear branch `backup/sprint-00` por si necesitas rollback
2. **Tests primero:** Asegurar que tests pasan ANTES de eliminar código
3. **Commits atómicos:** Un commit por tarea completada
4. **Validar cada fase:** No pasar a siguiente fase si hay errores

---

## 🎯 CRITERIOS DE ÉXITO

Sprint-00 está COMPLETO cuando:

- ✅ Proyecto compila sin errores
- ✅ Todos los tests pasan (unit + integration)
- ✅ 0 código deprecated en uso
- ✅ Migraciones de infrastructure funcionan
- ✅ Eventos se validan contra schemas
- ✅ Coverage >= 80%
- ✅ Documentación actualizada

---

**Siguiente Sprint:** Sprint-01 - Schema de Base de Datos (Evaluaciones)

---

**Última actualización:** 16 de Noviembre, 2025  
**Versión:** 2.0.0  
**Mantenedor:** EduGo Team
