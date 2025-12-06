# 💳 Deuda Técnica

> **Última revisión:** Diciembre 2024  
> **Nivel de deuda:** 🟡 Medio

Este documento cataloga la deuda técnica acumulada que debe abordarse para mantener la salud del codebase.

---

## DEBT-001: Duplicación de Carpetas ValueObject

### Severidad: 🔴 Alta

### Descripción

Existen dos carpetas para value objects con contenido diferente:

```
internal/domain/
├── valueobject/     # 3 archivos
│   ├── ids.go
│   ├── email.go
│   └── ...
└── valueobjects/    # 8 archivos
    ├── mongo_document_id.go
    ├── question_id.go
    ├── score.go
    ├── time_spent.go
    └── ...
```

### Problema

1. **Confusión:** No está claro cuál usar para nuevos value objects
2. **Inconsistencia:** Singular vs plural en nombres
3. **Imports dispersos:** Diferentes partes del código importan de diferentes ubicaciones

### Solución Propuesta

1. **Consolidar en `valueobject/` (singular)** - Es el estándar en Go
2. **Mover todos los archivos de `valueobjects/` a `valueobject/`**
3. **Actualizar todos los imports**
4. **Eliminar carpeta `valueobjects/`**

### Script de Migración

```bash
#!/bin/bash
# migrate-valueobjects.sh

# 1. Mover archivos
mv internal/domain/valueobjects/*.go internal/domain/valueobject/

# 2. Actualizar imports en todo el proyecto
find . -name "*.go" -exec sed -i '' \
  's|domain/valueobjects|domain/valueobject|g' {} \;

# 3. Eliminar carpeta vacía
rmdir internal/domain/valueobjects

# 4. Verificar compilación
go build ./...

# 5. Ejecutar tests
go test ./...
```

### Esfuerzo Estimado
2-3 horas (incluyendo testing)

---

## DEBT-002: Código Comentado en Assessment Service

### Severidad: 🟡 Media

### Descripción

Hay bloques grandes de código comentado que deberían eliminarse o restaurarse:

**Archivo:** `internal/application/service/assessment_service.go:104-132`

```go
// TODO(sprint-00): Restaurar publicación de eventos cuando se defina schema
/*
    event := messaging.AssessmentAttemptRecordedEvent{
        AttemptID:    attempt.ID,
        UserID:       userID.String(),
        AssessmentID: assessment.MaterialID.String(),
        Score:        score,
        SubmittedAt:  time.Now(),
    }

    eventJSON, err := event.ToJSON()
    if err != nil {
        s.logger.Warn("failed to serialize assessment attempt recorded event",
            zap.String("attempt_id", attempt.ID),
            zap.Error(err),
        )
    } else {
        // ... más código comentado ...
    }
*/
```

### Problema

1. **Ruido visual:** Dificulta leer el código activo
2. **Mantenimiento:** El código comentado se desactualiza
3. **Git history:** El código está en el historial si se necesita

### Solución Propuesta

1. **Crear issue/ticket** para implementar la funcionalidad
2. **Eliminar el código comentado** del archivo
3. **Documentar en este archivo** qué funcionalidad falta
4. **Cuando se implemente:** Escribir código nuevo, no descomentar

### Archivos con Código Comentado

| Archivo | Líneas | Descripción |
|---------|--------|-------------|
| `assessment_service.go` | 104-132 | Publicación de eventos |
| `answer_repository_test.go` | 305+ | Tests de integración |
| `assessment_repository_test.go` | 201+ | Tests de integración |
| `assessment_document_repository_test.go` | 379+ | Tests de integración |

---

## DEBT-003: SchoolID Hardcodeado

### Severidad: 🔴 Alta

### Descripción

En `MaterialService.CreateMaterial`, el `schoolID` se genera como UUID aleatorio en lugar de obtenerlo del contexto de autenticación.

**Archivo:** `internal/application/service/material_service.go:63-64`

```go
// TODO: Obtener schoolID del contexto de autenticación
schoolID := uuid.New() // Temporal
```

### Impacto

1. **Multi-tenancy roto:** Los materiales no se asocian a la escuela correcta
2. **Queries incorrectas:** Filtrar por school_id no funciona
3. **Seguridad:** Potencial fuga de datos entre escuelas

### Dependencias para Resolver

1. `api-admin` debe incluir `school_id` en el JWT
2. Middleware de auth debe extraer y guardar en contexto
3. Handlers deben pasar al service

### Workaround Temporal

Si se necesita antes de la solución completa:

```go
// Extraer school_id del primer material del usuario (si existe)
existingMaterials, _ := s.materialRepo.FindByAuthor(ctx, authorID)
if len(existingMaterials) > 0 {
    schoolID = existingMaterials[0].SchoolID
} else {
    // Fallback: usar school_id de una tabla de users
    user, _ := s.userRepo.FindByID(ctx, authorID)
    if user != nil {
        schoolID = user.SchoolID
    }
}
```

---

## DEBT-004: Coexistencia de Sistemas de Assessment

### Severidad: 🟡 Media

### Descripción

Existen dos sistemas de assessment funcionando en paralelo:

```
Sistema Legacy (MongoDB)              Sistema Nuevo (PostgreSQL)
├── AssessmentRepository              ├── AssessmentRepoV2
├── assessment_service.go             ├── assessment_attempt_service.go
├── SubmitAssessment handler          ├── CreateMaterialAttempt handler
└── /assessments/:id/submit           └── /materials/:id/assessment/attempts
```

### Problema

1. **Confusión:** ¿Cuál usar para nuevas features?
2. **Mantenimiento doble:** Bugs deben arreglarse en ambos
3. **Datos dispersos:** Resultados en MongoDB y PostgreSQL
4. **Complejidad:** Más código que mantener

### Plan de Consolidación

```
┌─────────────────────────────────────────────────────────────────┐
│                    PLAN DE CONSOLIDACIÓN                         │
└─────────────────────────────────────────────────────────────────┘

Fase 1: Migración de Clientes (2 semanas)
├── Identificar todos los clientes del sistema legacy
├── Crear guía de migración
└── Comunicar timeline

Fase 2: Migración de Datos (1 semana)
├── Script para migrar datos de MongoDB a PostgreSQL
├── Validar integridad de datos migrados
└── Backup de datos originales

Fase 3: Deprecación (2 semanas)
├── Agregar warnings a endpoints legacy
├── Monitorear uso
└── Comunicar fecha de eliminación

Fase 4: Eliminación (1 día)
├── Eliminar código legacy
├── Eliminar colecciones MongoDB (assessment_results, assessment_attempts)
└── Actualizar documentación
```

---

## DEBT-005: Tests Unitarios con TODOs

### Severidad: 🟢 Baja

### Descripción

Varios archivos de tests tienen comentarios indicando que necesitan actualización:

```go
// TODO: Estos tests unitarios requieren actualización para usar mocks reales (sqlmock)
// Los tests de integración en *_integration_test.go
// validan el funcionamiento real con testcontainers
```

**Archivos afectados:**
- `answer_repository_test.go`
- `assessment_repository_test.go`
- `assessment_document_repository_test.go`

### Problema

1. **Cobertura incompleta:** Tests unitarios no validando correctamente
2. **Falsa seguridad:** Tests pasan pero no prueban el código real
3. **Mantenimiento:** Tests de integración son más lentos

### Solución

Opciones:
1. **Implementar mocks reales** con `sqlmock` y `mongomock`
2. **Eliminar tests unitarios** si los de integración son suficientes
3. **Documentar decisión** de solo usar tests de integración

---

## DEBT-006: Logger Inconsistente

### Severidad: 🟢 Baja

### Descripción

Hay inconsistencia en cómo se usa el logger:

```go
// Usando zap.String
s.logger.Info("material created",
    zap.String("material_id", material.ID.String()),
)

// Usando key-value pairs
s.logger.Info("attempt recorded", "material_id", materialID, "score", score)
```

### Problema

1. **Inconsistencia:** Diferentes estilos en diferentes archivos
2. **Refactoring difícil:** No hay un estándar claro

### Solución Propuesta

Estandarizar en key-value pairs (más simple y compatible con múltiples backends):

```go
// Estándar recomendado
s.logger.Info("message here",
    "key1", value1,
    "key2", value2,
)
```

---

## 📊 Resumen de Deuda Técnica

| ID | Descripción | Severidad | Esfuerzo | Estado |
|----|-------------|-----------|----------|--------|
| DEBT-001 | Duplicación valueobject/valueobjects | 🔴 Alta | 2-3h | ✅ Completado |
| DEBT-002 | Código comentado | 🟡 Media | 1h | ✅ Completado |
| DEBT-003 | SchoolID hardcodeado | 🔴 Alta | 4-8h | Pendiente |
| DEBT-004 | Sistemas assessment duplicados | 🟡 Media | 2-3 semanas | Pendiente |
| DEBT-005 | Tests con TODOs | 🟢 Baja | 4-8h | Pendiente |
| DEBT-006 | Logger inconsistente | 🟢 Baja | 2-4h | Pendiente |

---

## 📈 Métricas de Deuda

### Puntos de Deuda por Severidad

| Severidad | Puntos | Items |
|-----------|--------|-------|
| 🔴 Alta | 10 pts c/u | 1 |
| 🟡 Media | 5 pts c/u | 1 (era 2) |
| 🟢 Baja | 2 pts c/u | 2 |

**Total:** 19 puntos de deuda técnica (era 24)

### Objetivo

Reducir a menos de 15 puntos antes del Q2 2025.

---

## 🗓️ Historial de Pago de Deuda

| Fecha | Item | Puntos Pagados | PR |
|-------|------|----------------|-----|
| Dic 2024 | DEBT-001 | 10 pts | - |
| Dic 2024 | DEBT-002 | 5 pts | - |

**Puntos pagados este trimestre:** 15
