# Decisión: Tarea 4.8 - NO Migrar sync-main-to-dev.yml

**Fecha:** 2025-11-21
**Tarea:** 4.8 - Migrar sync-main-to-dev.yml
**Sprint:** SPRINT-4 - Workflows Reusables
**Fase:** 1 - Implementación

---

## Contexto

Al intentar migrar `sync-main-to-dev.yml` al workflow reusable `sync-branches.yml` de infrastructure, se identificaron diferencias significativas en funcionalidad.

---

## Análisis de Compatibilidad

### Workflow Actual (sync-main-to-dev.yml)

**Funcionalidades:**
1. ✅ Lectura de versión desde archivo (`.github/version.txt`)
2. ✅ Verificación si rama dev existe
3. ✅ Creación automática de dev si no existe
4. ✅ Verificación de commits ahead (diferencias)
5. ✅ Skip si no hay diferencias
6. ✅ Prevención de loops (no ejecutar en commits de sync)
7. ✅ Mensaje de commit personalizado con versión
8. ✅ Resumen detallado en GITHUB_STEP_SUMMARY
9. ✅ Manejo de conflictos con abort

**Líneas:** 128

### Workflow Reusable (sync-branches.yml)

**Funcionalidades:**
1. ✅ Sincronización básica main → dev
2. ✅ Creación de PR en conflictos
3. ✅ Auto-merge si no hay conflictos
4. ⚠️ **NO** lee versión
5. ⚠️ **NO** verifica diferencias (siempre intenta merge)
6. ⚠️ **NO** previene loops
7. ⚠️ **NO** mensajes personalizados
8. ⚠️ **NO** resumen en GITHUB_STEP_SUMMARY

**Líneas:** ~60 (estimado si se usa)

---

## Matriz de Comparación

| Feature | sync-main-to-dev.yml | sync-branches.yml | Compatible |
|---------|---------------------|-------------------|-----------|
| **Lectura de versión** | ✅ `.github/version.txt` | ❌ No soporta | ❌ No |
| **Verificar diferencias** | ✅ Skip si iguales | ❌ Siempre intenta | ❌ No |
| **Prevención loops** | ✅ Skip commits 'sync' | ❌ No soporta | ❌ No |
| **Mensaje custom** | ✅ Con versión | ⚠️ Genérico | ❌ No |
| **Resumen detallado** | ✅ STEP_SUMMARY | ⚠️ Básico | ❌ No |
| **Merge directo** | ✅ Push a dev | ✅ Soporta | ✅ Sí |
| **Manejo conflictos** | ✅ Abort + error | ✅ Crea PR | ⚠️ Diferente |

**Compatibilidad:** 1/7 features (14%) - **NO COMPATIBLE**

---

## Opciones Evaluadas

### Opción A: Migración Completa (Descartada)

Reemplazar todo con workflow reusable.

**Pros:**
- Reduce código (~50%)
- Centralización

**Contras:**
- ❌ Pierde lectura de versión
- ❌ Pierde verificación de diferencias (ejecuta siempre)
- ❌ Pierde prevención de loops
- ❌ Pierde mensaje personalizado con versión
- ❌ Pierde resumen detallado
- ❌ Regresión de funcionalidad

**Decisión:** ❌ Rechazada

---

### Opción B: Adaptar Workflow Reusable (Fuera de Alcance FASE 1)

Extender `sync-branches.yml` para soportar features de api-mobile.

**Tareas requeridas:**
1. Agregar parámetro `version-file` al workflow reusable
2. Agregar lógica de verificación de diferencias
3. Agregar prevención de loops
4. Agregar templates de mensaje de commit
5. Mejorar resumen en STEP_SUMMARY

**Pros:**
- ✅ Mantiene funcionalidades
- ✅ Workflow reusable más robusto
- ✅ Beneficia a otros proyectos

**Contras:**
- ⚠️ Requiere modificar infrastructure
- ⚠️ Testing extensivo necesario
- ⚠️ Fuera del alcance de FASE 1 de api-mobile
- ⚠️ Requiere coordinación entre proyectos

**Decisión:** ⏳ Pospuesto para sprint futuro en infrastructure

---

### Opción C: NO Migrar (SELECCIONADA)

Mantener `sync-main-to-dev.yml` custom en api-mobile.

**Justificación:**
1. ✅ Mantiene todas las funcionalidades actuales
2. ✅ Sin regresión
3. ✅ Workflow específico del proyecto (versión, etc.)
4. ✅ Migración completa no aporta valor en FASE 1

**Pros:**
- ✅ Sin riesgo de romper funcionalidad
- ✅ Mantiene lógica de versionado
- ✅ Prevención de loops funciona
- ✅ Mensajes informativos mantienen contexto

**Contras:**
- ⚠️ No reduce código de este workflow
- ⚠️ No centraliza (pero es específico del proyecto)

**Decisión:** ✅ Seleccionada para FASE 1

---

## Implementación de Opción C

### Archivo Actual: sync-main-to-dev.yml

**Estado:** MANTENIDO SIN CAMBIOS

**Razón:** Lógica específica del proyecto que no tiene equivalente en workflow reusable.

**Marcado como:** "No migrado (custom logic)" en documentación.

---

## Métricas

| Métrica | Valor |
|---------|-------|
| **Líneas actuales** | 128 |
| **Reducción de código** | 0% (no migrado) |
| **Funcionalidades mantenidas** | 100% |
| **Compatibilidad con reusable** | 14% |

---

## Plan para Migración Futura (Opcional)

### FASE 2 o Sprint en Infrastructure

**Estrategia:**
1. Extender workflow reusable `sync-branches.yml` en infrastructure
2. Agregar soporte para:
   - Version files
   - Skip si no hay diferencias
   - Prevención de loops
   - Mensajes personalizados
   - Resúmenes detallados
3. Migrar api-mobile una vez disponible

**Beneficio:**
- Todos los proyectos podrían usar sync con estas features
- Centralización real con funcionalidades completas

**Esfuerzo estimado:** 4-6 horas

---

## Documentación en Workflows

Agregar comentario al inicio de `sync-main-to-dev.yml`:

```yaml
# =====================================================
# NO MIGRADO a workflow reusable en SPRINT-4
# Razón: Lógica específica del proyecto (versión, loops, etc.)
# Ver: docs/cicd/tracking/decisions/TASK-4.8-NO-MIGRATION.md
# Migración futura: Extender sync-branches.yml en infrastructure
# =====================================================
```

---

## Conclusiones

✅ **NO migrar sync-main-to-dev.yml es la decisión correcta para FASE 1**
✅ **Mantiene todas las funcionalidades actuales**
✅ **Sin regresión de features**
⏳ **Migración posible en futuro si se extiende workflow reusable**

**Impacto en objetivos del sprint:**
- Workflows migrados: 2/3 (pr-to-dev, pr-to-main)
- Workflows no migrados: 1/3 (sync-main-to-dev)
- Reducción total de código: ~2-3% (vs ~75% esperado inicialmente)

**Razón principal:** Características personalizadas del proyecto incompatibles con workflows reusables actuales.

---

**Generado por:** Claude Code
**Fecha:** 2025-11-21
**Sprint:** SPRINT-4 FASE 1
**Tarea:** 4.8 completada (con decisión de no migración)
