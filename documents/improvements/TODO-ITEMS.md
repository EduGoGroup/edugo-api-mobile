# 📝 TODOs Pendientes

> **Última revisión:** Diciembre 2024  
> **Total TODOs:** 15+

Este documento cataloga todos los comentarios TODO encontrados en el codebase, organizados por prioridad y área.

---

## 🔴 Prioridad Alta

### TODO-001: Obtener SchoolID del Contexto de Autenticación

**Archivo:** `internal/application/service/material_service.go`  
**Línea:** 63-64

```go
// Crear entidad Material manualmente
// TODO: Obtener schoolID del contexto de autenticación
schoolID := uuid.New() // Temporal
```

#### Problema
Actualmente se genera un UUID aleatorio para `schoolID` en lugar de obtenerlo del JWT del usuario autenticado.

#### Impacto
- Los materiales no se asocian correctamente a la escuela del docente
- Afecta queries de filtrado por escuela
- Puede causar problemas de aislamiento de datos multi-tenant

#### Solución Propuesta

```go
// 1. Agregar schoolID al JWT claims en api-admin
type CustomClaims struct {
    jwt.RegisteredClaims
    UserID   string `json:"user_id"`
    Email    string `json:"email"`
    Role     string `json:"role"`
    SchoolID string `json:"school_id"` // ← Agregar
}

// 2. Extraer en middleware
func RemoteAuthMiddleware(...) gin.HandlerFunc {
    return func(c *gin.Context) {
        // ...validación...
        c.Set("school_id", claims.SchoolID)
    }
}

// 3. Usar en servicio
func (s *materialService) CreateMaterial(...) {
    schoolIDStr := ginmiddleware.MustGetSchoolID(c)
    schoolID, _ := uuid.Parse(schoolIDStr)
    // ...
}
```

#### Dependencias
- Requiere cambio en `api-admin` para incluir `school_id` en JWT
- Requiere actualizar middleware de autenticación

---

### TODO-002: Agregar Middleware de Autorización para Admins

**Archivo:** `internal/infrastructure/http/router/router.go`  
**Línea:** 135

```go
// setupStatsRoutes configura rutas de estadísticas globales del sistema.
func setupStatsRoutes(rg *gin.RouterGroup, c *container.Container) {
    stats := rg.Group("/stats")
    {
        // Estadísticas globales del sistema (Fase 6)
        // TODO: Agregar middleware de autorización para solo admins
        stats.GET("/global", c.Handlers.StatsHandler.GetGlobalStats)
    }
}
```

#### Problema
El endpoint `/v1/stats/global` está accesible para cualquier usuario autenticado, pero debería ser solo para administradores.

#### Impacto
- Exposición de métricas sensibles del sistema
- Cualquier usuario puede ver estadísticas globales

#### Solución Propuesta

```go
// 1. Crear middleware de autorización por rol
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        role := c.GetString("role")
        for _, allowed := range allowedRoles {
            if role == allowed {
                c.Next()
                return
            }
        }
        c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
            "error": "insufficient permissions",
            "code":  "FORBIDDEN",
        })
    }
}

// 2. Aplicar en router
stats.GET("/global",
    middleware.RequireRole("admin", "super_admin"),
    c.Handlers.StatsHandler.GetGlobalStats,
)
```

---

### TODO-003: Verificación de Rol Admin en Progress Handler

**Archivo:** `internal/infrastructure/http/handler/progress_handler.go`  
**Línea:** 109-110

```go
// Autorización: Usuario solo puede actualizar su propio progreso (a menos que sea admin)
// TODO: Agregar verificación de rol admin cuando exista
if req.UserID != authenticatedUserID {
    // ...error 403...
}
```

#### Problema
Un admin debería poder actualizar el progreso de cualquier usuario (para soporte técnico), pero actualmente no hay bypass para admins.

#### Solución Propuesta

```go
// Verificar si es admin
role := c.GetString("role")
isAdmin := role == "admin" || role == "super_admin"

// Autorización: Usuario solo puede actualizar su propio progreso (a menos que sea admin)
if !isAdmin && req.UserID != authenticatedUserID {
    h.logger.Warn("user attempting to update progress of another user",
        "authenticated_user_id", authenticatedUserID,
        "target_user_id", req.UserID,
    )
    c.JSON(http.StatusForbidden, ErrorResponse{
        Error: "you can only update your own progress",
        Code:  "FORBIDDEN",
    })
    return
}
```

---

## 🟡 Prioridad Media

### TODO-004: URL Real de S3 en Material Service

**Archivo:** `internal/application/service/material_service.go`  
**Líneas:** 116-117

```go
payload := rabbitmq.MaterialUploadedPayload{
    MaterialID:    material.ID.String(),
    SchoolID:      material.SchoolID.String(),
    TeacherID:     authorID.String(),
    FileURL:       "s3://edugo/materials/" + material.ID.String() + ".pdf", // TODO: URL real de S3
    FileSizeBytes: 0,  // TODO: obtener tamaño real del archivo
    FileType:      "application/pdf",
}
```

#### Problema
El evento de RabbitMQ se publica con una URL placeholder en lugar de la URL real de S3.

#### Solución Propuesta
El evento debería publicarse DESPUÉS de `NotifyUploadComplete`, no en `CreateMaterial`:

```go
// En NotifyUploadComplete, después de actualizar el material:
payload := rabbitmq.MaterialUploadedPayload{
    MaterialID:    material.ID.String(),
    SchoolID:      material.SchoolID.String(),
    TeacherID:     material.UploadedByTeacherID.String(),
    FileURL:       req.FileURL,           // URL real de S3
    FileSizeBytes: req.FileSizeBytes,     // Tamaño real
    FileType:      req.FileType,          // Tipo real
}
```

---

### TODO-005: Restaurar Publicación de Eventos de Assessment

**Archivo:** `internal/application/service/assessment_service.go`  
**Líneas:** 100-132

```go
// TODO(sprint-00): Restaurar publicación de eventos cuando se defina schema
// para assessment.attempt.recorded en edugo-infrastructure/schemas
/*
    event := messaging.AssessmentAttemptRecordedEvent{
        AttemptID:    attempt.ID,
        UserID:       userID.String(),
        AssessmentID: assessment.MaterialID.String(),
        Score:        score,
        SubmittedAt:  time.Now(),
    }
    // ... publicación comentada ...
*/
```

#### Problema
El código para publicar eventos de intentos de assessment está comentado porque falta definir el schema en `edugo-infrastructure`.

#### Pasos para Resolver
1. Definir schema `assessment.attempt.recorded` en `edugo-infrastructure/schemas`
2. Generar tipos Go desde el schema
3. Descomentar y adaptar el código
4. Agregar tests

---

### TODO-006: Implementar FindByIDWithVersions Completo

**Archivo:** `internal/infrastructure/persistence/postgres/repository/material_repository_impl.go`  
**Línea:** 369

```go
func (r *postgresMaterialRepository) FindByIDWithVersions(ctx context.Context, id valueobject.MaterialID) (*pgentities.Material, []*pgentities.MaterialVersion, error) {
    // Por ahora solo retorna el material sin versiones
    // TODO: Implementar join con material_versions cuando se necesite
    material, err := r.FindByID(ctx, id)
    if err != nil {
        return nil, nil, err
    }
    return material, nil, nil  // ← Sin versiones
}
```

#### Problema
El método no hace join con `material_versions` y siempre retorna un slice vacío de versiones.

#### Solución Propuesta

```go
func (r *postgresMaterialRepository) FindByIDWithVersions(ctx context.Context, id valueobject.MaterialID) (*pgentities.Material, []*pgentities.MaterialVersion, error) {
    // 1. Obtener material
    material, err := r.FindByID(ctx, id)
    if err != nil {
        return nil, nil, err
    }

    // 2. Obtener versiones
    query := `
        SELECT id, material_id, version_number, title, content_url, changed_by, created_at
        FROM material_versions
        WHERE material_id = $1
        ORDER BY version_number DESC
    `

    rows, err := r.db.QueryContext(ctx, query, id.UUID().UUID)
    if err != nil {
        return nil, nil, err
    }
    defer rows.Close()

    var versions []*pgentities.MaterialVersion
    for rows.Next() {
        var v pgentities.MaterialVersion
        if err := rows.Scan(&v.ID, &v.MaterialID, &v.VersionNumber, &v.Title, &v.ContentURL, &v.ChangedBy, &v.CreatedAt); err != nil {
            return nil, nil, err
        }
        versions = append(versions, &v)
    }

    return material, versions, nil
}
```

---

### TODO-007: Publicar Evento material_completed

**Archivo:** `internal/application/service/progress_service.go`  
**Líneas:** 110-118

```go
// Verificar si material fue completado (progress = 100)
isCompleted := updatedProgress.Percentage == 100
if isCompleted {
    s.logger.Info("material completed by user", ...)

    // TODO (Fase futura): Publicar evento "material_completed" a RabbitMQ
    // Example:
    // event := events.MaterialCompleted{
    //     MaterialID: materialID,
    //     UserID: userIDStr,
    //     CompletedAt: updatedProgress.UpdatedAt(),
    // }
    // s.eventPublisher.Publish(ctx, "material.completed", event)
}
```

#### Problema
Cuando un usuario completa un material (100%), no se publica ningún evento para analytics o gamificación.

#### Uso del Evento
- Actualizar dashboard de completados
- Otorgar badges/logros
- Notificar al docente
- Analytics de engagement

---

## 🟢 Prioridad Baja

### TODO-008: Implementar Lógica de Deshabilitación de Recursos

**Archivo:** `internal/bootstrap/config.go`  
**Línea:** 96-97

```go
func WithDisabledResource(resourceName string) BootstrapOption {
    return func(opts *BootstrapOptions) {
        if opts.OptionalResources == nil {
            opts.OptionalResources = make(map[string]bool)
        }
        // Marcar como opcional para que use noop
        opts.OptionalResources[resourceName] = true
        // TODO: Implementar lógica de deshabilitación en el bootstrapper
    }
}
```

#### Problema
La función `WithDisabledResource` solo marca como opcional pero no deshabilita completamente.

---

### TODO-009: Tests de Integración con Testcontainers MongoDB

**Archivos:**
- `internal/infrastructure/persistence/mongodb/repository/assessment_document_repository_test.go:379`
- `internal/infrastructure/persistence/postgres/repository/answer_repository_test.go:305`
- `internal/infrastructure/persistence/postgres/repository/assessment_repository_test.go:201`

```go
// TODO: Claude Local - Tests de integración con testcontainers MongoDB
// func TestMongoAssessmentDocumentRepository_Integration(t *testing.T) {
//     if testing.Short() {
//         t.Skip("Skipping integration test")
//     }
// ...código comentado...
```

#### Problema
Hay tests de integración comentados que deberían activarse o eliminarse.

#### Acción
- Si los tests en `*_integration_test.go` cubren los casos, eliminar el código comentado
- Si no, implementar los tests faltantes

---

## 📊 Resumen por Área

| Área | Cantidad | Prioridad Promedio |
|------|----------|-------------------|
| Autenticación/Autorización | 3 | 🔴 Alta |
| Eventos/Messaging | 3 | 🟡 Media |
| Persistencia | 2 | 🟡 Media |
| Tests | 4 | 🟢 Baja |
| Bootstrap/Config | 2 | 🟢 Baja |

---

## 🗓️ Historial de Resolución

| Fecha | TODO | PR | Descripción |
|-------|------|-----|-------------|
| - | - | - | Ningún TODO resuelto aún |
