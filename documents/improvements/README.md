# 🔧 Mejoras y Refactorizaciones Pendientes

> **Propósito:** Este directorio documenta código que debe ser mejorado, eliminado o refactorizado para mantener la calidad del codebase.

## 📋 Índice de Mejoras

| Documento | Prioridad | Descripción |
|-----------|-----------|-------------|
| [DEPRECATED-CODE.md](./DEPRECATED-CODE.md) | 🔴 Alta | Código marcado como deprecado para eliminar |
| [TODO-ITEMS.md](./TODO-ITEMS.md) | 🟡 Media | TODOs pendientes en el código |
| [LEGACY-ENDPOINTS.md](./LEGACY-ENDPOINTS.md) | 🟡 Media | Endpoints legacy a migrar/eliminar |
| [REFACTORING-OPPORTUNITIES.md](./REFACTORING-OPPORTUNITIES.md) | 🟢 Baja | Oportunidades de mejora de código |
| [TECHNICAL-DEBT.md](./TECHNICAL-DEBT.md) | 🔴 Alta | Deuda técnica acumulada |

---

## 📊 Resumen de Estado

### Código Deprecado
- **5 funciones** en `bootstrap/bootstrap.go` marcadas como DEPRECATED
- **1 repositorio** legacy (`AssessmentRepository`) coexiste con versión nueva
- **2 endpoints** legacy que deberían migrarse

### TODOs Pendientes
- **15+ TODOs** identificados en el codebase
- Principalmente relacionados con:
  - Obtener `schoolID` del contexto de autenticación
  - Implementar verificación de rol admin
  - Publicación de eventos RabbitMQ pendientes
  - Tests de integración incompletos

### Deuda Técnica
- Duplicación entre `valueobject/` y `valueobjects/`
- Coexistencia de sistema de assessments legacy y nuevo
- Código comentado que debería eliminarse

---

## 🎯 Plan de Acción Recomendado

### Fase 1: Limpieza Inmediata (1-2 días)
1. Eliminar funciones `WithInjected*` deprecadas
2. Eliminar código comentado
3. Consolidar carpetas `valueobject/` y `valueobjects/`

### Fase 2: Migración de Legacy (1 semana)
1. Migrar clientes del endpoint `PATCH /materials/:id/progress` a `PUT /progress`
2. Migrar clientes del endpoint `POST /assessments/:id/submit` al nuevo sistema
3. Eliminar repositorio legacy de assessments

### Fase 3: Completar Funcionalidad (2 semanas)
1. Implementar obtención de `schoolID` desde JWT
2. Agregar middleware de autorización para admins
3. Implementar eventos pendientes de RabbitMQ

---

## 🔍 Cómo Usar Esta Documentación

1. **Antes de trabajar en una área:** Revisar si hay mejoras pendientes
2. **Al encontrar código problemático:** Agregar a esta documentación
3. **Al completar una mejora:** Marcarla como completada con fecha y PR
4. **En code reviews:** Verificar que no se agregue más deuda técnica

---

## 📝 Plantilla para Nuevas Mejoras

```markdown
## [TIPO-XXX] Título descriptivo

**Archivo(s):** `ruta/al/archivo.go`
**Prioridad:** 🔴 Alta | 🟡 Media | 🟢 Baja
**Esfuerzo estimado:** X horas/días
**Impacto:** Descripción del impacto

### Descripción
Explicación detallada del problema.

### Código Actual
```go
// Código problemático
```

### Solución Propuesta
```go
// Código mejorado
```

### Pasos de Migración
1. Paso 1
2. Paso 2

### Riesgos
- Riesgo 1
- Riesgo 2
```
