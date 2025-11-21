# 📊 Comparación de Migraciones: Local vs Infrastructure

**Fecha:** 16 de Noviembre, 2025  
**Decisión:** Eliminar TODAS las migraciones locales  
**Razón:** Centralización en edugo-infrastructure

---

## 🎯 Decisión Arquitectónica

**TODAS las migraciones locales están DEPRECATED y se eliminarán.**

### Principio de Centralización

- ✅ **Infrastructure es la fuente de verdad**: Todas las migraciones PostgreSQL vienen de `edugo-infrastructure/postgres`
- ✅ **Eliminar duplicación**: No mantener scripts SQL en cada proyecto
- ✅ **Responsabilidades claras**: Infrastructure maneja schema, proyectos consumen

---

## 📋 Scripts Locales a ELIMINAR

### Archivo: `01_create_schema.sql` (297 líneas)
**Estado:** ELIMINAR COMPLETAMENTE  
**Razón:** Duplica migraciones de infrastructure (001-008)

**Tablas en local:**
- `users` → Existe en infrastructure `001_create_users.up.sql`
- `materials` → Existe en infrastructure `005_create_materials.up.sql`
- `material_progress` → **NO existe en infrastructure** (específico de api-mobile)

**Diferencias clave:**
```diff
Local (api-mobile):
- materials.author_id (directo al teacher)
- materials.subject_id (string simple)
- Solo 3 tablas básicas

Infrastructure:
+ materials.school_id (REQUIRED)
+ materials.uploaded_by_teacher_id (más descriptivo)
+ materials.academic_unit_id (estructura académica)
+ materials.file_size_bytes, file_type (metadata completa)
+ Incluye: schools, academic_units, memberships
```

**Acción:** Eliminar archivo. Usar migraciones de infrastructure.

---

### Archivo: `02_seed_data.sql` (424 líneas)
**Estado:** ELIMINAR COMPLETAMENTE  
**Razón:** Datos de prueba locales. Infrastructure no maneja seeds (correcto).

**Contenido:**
- Inserts de usuarios de prueba
- Inserts de materiales de ejemplo
- Datos para desarrollo local

**Acción:** Eliminar archivo. Los seeds se manejan vía testcontainers en tests o scripts de desarrollo separados (no en migrations).

---

### Archivo: `03_refresh_tokens.sql` (133 líneas)
**Estado:** ELIMINAR COMPLETAMENTE  
**Razón:** Funcionalidad movida a edugo-shared/auth

**Contenido:**
- Tabla `refresh_tokens`
- Índices, vistas, funciones de limpieza

**Acción:** Eliminar archivo. Si refresh tokens es necesario, debe estar en:
- Infrastructure (tabla global) O
- Shared/auth (manejo de autenticación)

---

### Archivo: `04_login_attempts.sql` (185 líneas)
**Estado:** ELIMINAR COMPLETAMENTE  
**Razón:** Seguridad debe estar centralizada en infrastructure o shared

**Contenido:**
- Tabla `login_attempts` (rate limiting)
- Rate limiting functions
- Security triggers

**Acción:** Eliminar archivo. Si se necesita rate limiting:
- Infrastructure (tabla global de auditoría) O
- Shared/auth (lógica de seguridad) O
- Middleware de Gin (shared/middleware/gin)

---

### Archivo: `04_material_versions.sql` (72 líneas)
**Estado:** ELIMINAR COMPLETAMENTE  
**Razón:** Versionado de materiales - funcionalidad no prioritaria

**Contenido:**
- Tabla `material_versions` (historial de cambios)
- Triggers de versionado automático

**Acción:** Eliminar archivo. Si se necesita versionado en el futuro:
- Agregar a infrastructure como feature completa
- No mantener en proyectos individuales

---

### Archivo: `05_indexes_materials.sql` (33 líneas)
**Estado:** ELIMINAR COMPLETAMENTE  
**Razón:** Índices ya están en infrastructure migrations

**Contenido:**
- `idx_materials_updated_at`
- Otros índices de performance

**Acción:** Eliminar archivo. Infrastructure `005_create_materials.up.sql` ya incluye todos los índices necesarios.

---

### Archivo: `05_user_progress_upsert.sql` (113 líneas)
**Estado:** ELIMINAR COMPLETAMENTE  
**Razón:** Material progress debe estar en infrastructure

**Contenido:**
- Función `upsert_user_progress()`
- Lógica de actualización de progreso

**Acción:**
1. Eliminar archivo local
2. **VERIFICAR** si `material_progress` existe en infrastructure
3. Si NO existe, crear PR en infrastructure para agregarlo (es tabla compartida)

---

## 📊 Resumen de Eliminación

| Archivo | Líneas | Acción | Razón |
|---------|--------|--------|-------|
| `01_create_schema.sql` | 297 | ❌ ELIMINAR | Duplica infrastructure |
| `02_seed_data.sql` | 424 | ❌ ELIMINAR | Seeds locales (no en migrations) |
| `03_refresh_tokens.sql` | 133 | ❌ ELIMINAR | Debe estar en shared/auth |
| `04_login_attempts.sql` | 185 | ❌ ELIMINAR | Debe estar en infrastructure/shared |
| `04_material_versions.sql` | 72 | ❌ ELIMINAR | Feature no prioritaria |
| `05_indexes_materials.sql` | 33 | ❌ ELIMINAR | Ya en infrastructure |
| `05_user_progress_upsert.sql` | 113 | ❌ ELIMINAR | Verificar en infrastructure |
| **TOTAL** | **1257** | **❌ ELIMINAR TODO** | Centralización |

---

## 🚨 Tablas/Features Faltantes en Infrastructure

### 1. `material_progress`
**Descripción:** Progreso de lectura de estudiantes en materiales  
**Columnas clave:**
- `material_id`, `user_id` (PK compuesta)
- `percentage`, `last_page`, `status`
- `last_accessed_at`

**Acción requerida:**
- ⚠️ **VERIFICAR** si existe en infrastructure
- Si NO existe → Crear issue/PR en edugo-infrastructure
- Esta tabla es **crítica** para funcionalidad de api-mobile

### 2. `refresh_tokens`
**Descripción:** Tokens de renovación JWT  
**Decisión:** Debe estar en edugo-shared/auth (no en DB si es posible) o en infrastructure como tabla global

### 3. `login_attempts`
**Descripción:** Rate limiting y seguridad  
**Decisión:** Debe estar en infrastructure (auditoría global) o manejarse vía middleware

---

## ✅ Próximos Pasos (TASK-005)

1. ✅ Análisis completado (este documento)
2. ⏭️ Eliminar carpeta completa `scripts/postgresql/`
3. ⏭️ Verificar que no haya referencias a estos scripts en código
4. ⏭️ Actualizar documentación sobre migraciones
5. ⏭️ Commit: "refactor(sprint-00): eliminar migraciones locales deprecated"

---

## 📝 Notas Importantes

### ¿Por qué eliminar en vez de migrar?

**Razón 1: Evitar duplicación**
- Infrastructure ya tiene las tablas core (users, materials, schools, etc.)
- Mantener scripts locales genera inconsistencias

**Razón 2: Responsabilidad clara**
- Infrastructure: Schema global
- API Mobile: Consumir schema via migraciones de infrastructure
- Shared: Lógica de negocio reutilizable

**Razón 3: Mantenibilidad**
- Un solo lugar para cambios de schema
- Versionado claro con tags de infrastructure
- Rollbacks centralizados

### ¿Qué pasa con datos de prueba?

**Opción 1: Testcontainers** (Recomendado)
```go
// En tests, usar infrastructure migrations + seeds en código
func setupTestDB(t *testing.T) *sql.DB {
    // 1. Crear container con postgres
    // 2. Ejecutar migrations de infrastructure
    // 3. Insertar datos de prueba via código Go
}
```

**Opción 2: Scripts de desarrollo** (No en repo principal)
```bash
# scripts/dev/seed-local.sh (gitignored)
# Solo para desarrollo local, no para producción
```

---

**Documento generado:** TASK-004 - Análisis de Migraciones  
**Responsable:** Claude Code  
**Próximo paso:** TASK-005 - Eliminar scripts SQL
