# 🔍 Análisis de Modernización - Sprint-00

**Proyecto:** edugo-api-mobile  
**Fecha:** 16 de Noviembre, 2025  
**Objetivo:** Actualizar proyecto para usar edugo-infrastructure v0.5.0 y deprecar código obsoleto

---

## 📊 Estado Actual vs. Estado Deseado

### Versiones Actuales en go.mod

```go
// ❌ DESACTUALIZADO - Módulos de edugo-shared
github.com/EduGoGroup/edugo-shared/auth v0.3.3              // Latest: v0.7.0
github.com/EduGoGroup/edugo-shared/bootstrap v0.5.0         // ✅ OK
github.com/EduGoGroup/edugo-shared/common v0.5.0            // ✅ OK
github.com/EduGoGroup/edugo-shared/lifecycle v0.5.0         // ✅ OK
github.com/EduGoGroup/edugo-shared/logger v0.5.0            // ✅ OK
github.com/EduGoGroup/edugo-shared/middleware/gin v0.3.3    // Latest: v0.7.0
github.com/EduGoGroup/edugo-shared/testing v0.6.2           // ✅ OK

// ❌ FALTANTE - Módulos de edugo-infrastructure
github.com/EduGoGroup/edugo-infrastructure/postgres         // NECESARIO: v0.5.0
github.com/EduGoGroup/edugo-infrastructure/mongodb          // NECESARIO: v0.5.0
github.com/EduGoGroup/edugo-infrastructure/messaging        // NECESARIO: v0.5.0
github.com/EduGoGroup/edugo-infrastructure/database         // NECESARIO: v0.1.1

// ❌ FALTANTE - Módulos adicionales de edugo-shared
github.com/EduGoGroup/edugo-shared/database/postgres        // NECESARIO: v0.7.0
github.com/EduGoGroup/edugo-shared/database/mongodb         // NECESARIO: v0.7.0
github.com/EduGoGroup/edugo-shared/config                   // NECESARIO: v0.7.0
```

### Releases Disponibles de edugo-infrastructure

```bash
# Módulos por separado (USAR ESTOS)
postgres/v0.5.0          # ✅ Migraciones PostgreSQL
mongodb/v0.5.0           # ✅ Migraciones MongoDB
messaging/v0.5.0         # ✅ Validación de eventos
database/v0.1.1          # ✅ Utilities de database

# Releases globales (legacy)
v0.1.0, v0.1.1, v0.5.0   # No usar, usar módulos específicos
```

---

## 🗑️ CÓDIGO DEPRECATED A ELIMINAR

### 1. Scripts SQL Locales (DUPLICADOS)

**Ubicación:** `scripts/postgresql/`

```bash
scripts/postgresql/
├── 01_create_schema.sql           # ❌ DUPLICADO - Existe en infrastructure
├── 02_seed_data.sql               # ❌ DUPLICADO - Existe en infrastructure
├── 03_refresh_tokens.sql          # ❌ DUPLICADO - Migración 001 en infrastructure
├── 04_login_attempts.sql          # ⚠️  REVISAR - Puede ser específico del proyecto
├── 04_material_versions.sql       # ❌ DUPLICADO - Migración 005 en infrastructure
├── 05_indexes_materials.sql       # ❌ DUPLICADO - Migración 005 en infrastructure
├── 05_user_progress_upsert.sql    # ⚠️  REVISAR - Puede ser específico del proyecto
```

**Acción:**
- ✅ **ELIMINAR completamente** si las migraciones están en infrastructure
- ⚠️  **MIGRAR a infrastructure** si son específicas del proyecto pero deberían ser compartidas
- 📋 **MANTENER** solo si son 100% específicas de api-mobile y no compartibles

**Migraciones en infrastructure/postgres:**
```
001_create_users.up.sql           # Tabla users
002_create_schools.up.sql         # Tabla schools
003_create_academic_units.up.sql  # Cursos, clases
004_create_memberships.up.sql     # Relación user-school-course
005_create_materials.up.sql       # Materiales educativos
006_create_assessments.up.sql     # Quizzes (ref MongoDB)
```

**Recomendación:**
- Eliminar `01_create_schema.sql`, `02_seed_data.sql`, `03_refresh_tokens.sql`, `04_material_versions.sql`, `05_indexes_materials.sql`
- Revisar si `04_login_attempts.sql` y `05_user_progress_upsert.sql` deben ir a infrastructure

---

### 2. Implementación Local de Conexión a BD (REEMPLAZAR)

**Archivo:** `internal/infrastructure/database/postgres.go`

```go
// ❌ CÓDIGO ACTUAL (custom)
func InitPostgreSQL(ctx context.Context, cfg *config.Config, log logger.Logger) (*sql.DB, error) {
    connStr := cfg.Database.Postgres.GetConnectionString()
    db, err := sql.Open("postgres", connStr)
    // ... configuración manual del pool
    // ... ping manual
}
```

**Problema:**
- Duplica lógica que ya existe en `edugo-shared/database/postgres`
- No usa el connector estándar del ecosistema
- Configuración manual propensa a errores

**Solución:** Usar módulo de shared

```go
// ✅ CÓDIGO NUEVO (usando shared)
import "github.com/EduGoGroup/edugo-shared/database/postgres"

func InitPostgreSQL(ctx context.Context, cfg *config.Config, log logger.Logger) (*sql.DB, error) {
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

**Acción:**
- ❌ **ELIMINAR** `internal/infrastructure/database/postgres.go`
- ❌ **ELIMINAR** `internal/infrastructure/database/postgres_test.go`
- ✅ **USAR** `edugo-shared/database/postgres` directamente

---

### 3. Implementación Local de MongoDB (REEMPLAZAR)

**Archivo:** `internal/infrastructure/database/mongodb.go`

Similar a PostgreSQL, eliminar implementación custom y usar shared.

```go
// ✅ NUEVO
import "github.com/EduGoGroup/edugo-shared/database/mongodb"

func InitMongoDB(ctx context.Context, cfg *config.Config, log logger.Logger) (*mongo.Client, error) {
    return mongodb.NewConnection(ctx, mongodb.Config{
        URI:      cfg.Database.MongoDB.URI,
        Database: cfg.Database.MongoDB.Database,
        Timeout:  cfg.Database.MongoDB.Timeout,
    })
}
```

**Acción:**
- ❌ **ELIMINAR** `internal/infrastructure/database/mongodb.go`
- ❌ **ELIMINAR** `internal/infrastructure/database/mongodb_test.go`
- ✅ **USAR** `edugo-shared/database/mongodb`

---

### 4. Configuración Manual de RabbitMQ (MODERNIZAR)

**Ubicación:** `internal/infrastructure/messaging/`

**Problema:**
- Configuración manual de conexiones, exchanges, queues
- No usa templates de infrastructure
- Validación de eventos inexistente

**Solución:**
```go
// ✅ NUEVO
import "github.com/EduGoGroup/edugo-infrastructure/messaging"

// Validar eventos antes de publicar
validator := messaging.NewValidator()
if err := validator.Validate(event, "material-uploaded-v1"); err != nil {
    return fmt.Errorf("evento inválido: %w", err)
}
```

**Acción:**
- ✅ **AGREGAR** validación de eventos con schemas de infrastructure
- ✅ **MODERNIZAR** código de messaging para usar patterns de infrastructure
- 📋 **MANTENER** lógica específica del negocio

---

## 🚀 OPORTUNIDADES DE MEJORA

### 1. Usar Módulo Config de Shared

**Actual:** Configuración custom en `internal/config/`

**Mejora:**
```go
import "github.com/EduGoGroup/edugo-shared/config"

// Usar loader estándar
cfg, err := config.Load(
    config.WithEnvPrefix("EDUGO_API_MOBILE"),
    config.WithConfigPaths("./config"),
)
```

**Beneficios:**
- ✅ Validación automática de configuración
- ✅ Hot-reload de configuración
- ✅ Soporte multi-ambiente consistente
- ✅ Variables de entorno estandarizadas

---

### 2. Migrar a Bootstrap de Shared

**Actual:** Inicialización manual de servicios en `cmd/main.go`

**Mejora:**
```go
import "github.com/EduGoGroup/edugo-shared/bootstrap"

app := bootstrap.New(
    bootstrap.WithLogger(),
    bootstrap.WithPostgres(),
    bootstrap.WithMongoDB(),
    bootstrap.WithRabbitMQ(),
    bootstrap.WithGracefulShutdown(),
)

if err := app.Run(ctx); err != nil {
    log.Fatal(err)
}
```

**Beneficios:**
- ✅ Inicialización declarativa
- ✅ Graceful shutdown automático
- ✅ Health checks integrados
- ✅ Menos boilerplate

---

### 3. Ejecutar Migraciones desde Código

**Actual:** Migraciones manuales ejecutadas fuera del proyecto

**Mejora:**
```go
import "github.com/EduGoGroup/edugo-infrastructure/postgres"

// En inicialización de la app
migrator := postgres.NewMigrator(db)
if err := migrator.Up(ctx); err != nil {
    log.Fatal("error ejecutando migraciones:", err)
}
```

**Beneficios:**
- ✅ Migraciones automáticas en startup (opcional)
- ✅ Rollback programático
- ✅ Estado de migraciones en health check
- ✅ Tests con migraciones incluidas

---

### 4. Testing con Infrastructure

**Actual:** Tests con setup manual de BD

**Mejora:**
```go
import "github.com/EduGoGroup/edugo-infrastructure/database"

func TestUserRepository(t *testing.T) {
    // Testcontainer con migraciones incluidas
    db := database.NewTestPostgres(t,
        database.WithMigrations(), // Ejecuta migraciones de infrastructure
    )

    repo := NewUserRepository(db)
    // ... tests
}
```

**Beneficios:**
- ✅ Tests con schema real (no mocks)
- ✅ Migraciones aplicadas automáticamente
- ✅ Cleanup automático después de tests
- ✅ Paralelización segura

---

### 5. Validación de Eventos en CI/CD

**Actual:** Sin validación de eventos

**Mejora:**
```yaml
# .github/workflows/validate-events.yml
- name: Validar eventos
  run: |
    cd messaging
    go run validator.go validate ../internal/events/*.json
```

**Beneficios:**
- ✅ Detectar eventos inválidos en PR
- ✅ Breaking changes visibles antes de merge
- ✅ Contratos de eventos documentados

---

## 📋 RESUMEN DE CAMBIOS NECESARIOS

### Eliminar (Deprecated)

| Archivo/Carpeta | Razón | Reemplazo |
|-----------------|-------|-----------|
| `scripts/postgresql/*.sql` | Duplicado en infrastructure | `edugo-infrastructure/postgres/migrations/` |
| `internal/infrastructure/database/postgres.go` | Custom connector | `edugo-shared/database/postgres` |
| `internal/infrastructure/database/mongodb.go` | Custom connector | `edugo-shared/database/mongodb` |
| `internal/infrastructure/database/*_test.go` | Tests de connectors custom | Tests con `database.NewTestPostgres()` |

### Actualizar (Versiones)

| Módulo | Versión Actual | Versión Nueva | Razón |
|--------|----------------|---------------|-------|
| `edugo-shared/auth` | v0.3.3 | v0.7.0 | JWT improvements, security patches |
| `edugo-shared/middleware/gin` | v0.3.3 | v0.7.0 | CORS, rate limiting, nuevos middlewares |

### Agregar (Nuevos Módulos)

| Módulo | Versión | Propósito |
|--------|---------|-----------|
| `edugo-infrastructure/postgres` | v0.5.0 | Migraciones centralizadas |
| `edugo-infrastructure/mongodb` | v0.5.0 | Migraciones MongoDB |
| `edugo-infrastructure/messaging` | v0.5.0 | Validación de eventos |
| `edugo-infrastructure/database` | v0.1.1 | Utilities de testing |
| `edugo-shared/database/postgres` | v0.7.0 | Connector estándar PostgreSQL |
| `edugo-shared/database/mongodb` | v0.7.0 | Connector estándar MongoDB |
| `edugo-shared/config` | v0.7.0 | Loader de configuración |

### Modernizar (Refactorizar)

| Componente | Acción | Impacto |
|------------|--------|---------|
| `cmd/main.go` | Migrar a bootstrap | Menos código, más robusto |
| `internal/config/` | Usar config de shared | Validación automática |
| `internal/infrastructure/messaging/` | Agregar validación de eventos | Contratos explícitos |
| Tests de integración | Usar testcontainers con migraciones | Tests más reales |

---

## 🎯 IMPACTO ESTIMADO

### Código a Eliminar
- **~500 líneas** de SQL (scripts locales)
- **~200 líneas** de código Go (connectors custom)
- **~300 líneas** de tests (connectors)
- **Total:** ~1,000 líneas eliminadas

### Código a Agregar
- **~100 líneas** de configuración (imports nuevos)
- **~50 líneas** de validación de eventos
- **~150 líneas** de tests modernizados
- **Total:** ~300 líneas agregadas

### Resultado Neto
- **-700 líneas** de código
- **+5 módulos** de dependencias (más mantenibles)
- **100%** de migraciones centralizadas
- **100%** de eventos validados

---

## ⚠️ RIESGOS Y MITIGACIONES

### Riesgo 1: Breaking Changes en Shared v0.7.0

**Probabilidad:** Media  
**Impacto:** Alto

**Mitigación:**
1. Leer CHANGELOG de edugo-shared v0.7.0
2. Ejecutar tests después de actualizar cada módulo
3. Hacer commits atómicos por módulo actualizado
4. Rollback fácil si algo falla

---

### Riesgo 2: Migraciones de Infrastructure No Cubren Casos Específicos

**Probabilidad:** Baja  
**Impacto:** Medio

**Mitigación:**
1. Revisar diff entre migraciones locales vs. infrastructure
2. Crear PRs en infrastructure para agregar migraciones faltantes
3. Documentar cualquier tabla específica de api-mobile

---

### Riesgo 3: Performance de Validación de Eventos

**Probabilidad:** Baja  
**Impacto:** Bajo

**Mitigación:**
1. Cachear validator de schemas
2. Solo validar en dev/qa (opcional en prod)
3. Benchmark de validación en tests

---

## 📝 PRÓXIMOS PASOS

Ver: `TASKS_ACTUALIZADO.md`

---

**Generado por:** Claude Code  
**Aprobado por:** [Pendiente]  
**Fecha de ejecución:** [Pendiente]

---

## 🔗 Referencias

- [edugo-infrastructure README](https://github.com/EduGoGroup/edugo-infrastructure/blob/main/README.md)
- [edugo-infrastructure CHANGELOG](https://github.com/EduGoGroup/edugo-infrastructure/blob/main/CHANGELOG.md)
- [edugo-shared Releases](https://github.com/EduGoGroup/edugo-shared/releases)
- [Documentación de Migraciones](https://github.com/EduGoGroup/edugo-infrastructure/blob/main/docs/TABLE_OWNERSHIP.md)
