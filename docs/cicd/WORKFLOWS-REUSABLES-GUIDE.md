# Guía de Workflows Reusables - edugo-api-mobile

**Proyecto:** edugo-api-mobile
**Sprint:** SPRINT-4 - Workflows Reusables
**Fecha:** 2025-11-21
**Versión:** 1.0

---

## 📋 Tabla de Contenidos

1. [Resumen Ejecutivo](#resumen-ejecutivo)
2. [Workflows Migrados](#workflows-migrados)
3. [Workflows Reusables Disponibles](#workflows-reusables-disponibles)
4. [Cómo Usar Workflows Reusables](#cómo-usar-workflows-reusables)
5. [Migración Híbrida Explicada](#migración-híbrida-explicada)
6. [Mantenimiento](#mantenimiento)
7. [Troubleshooting](#troubleshooting)
8. [Próximos Pasos](#próximos-pasos)

---

## 📊 Resumen Ejecutivo

### ¿Qué se hizo?

En **SPRINT-4** se migró parcialmente `edugo-api-mobile` para usar **workflows reusables centralizados** en `edugo-infrastructure`.

**Resultado:**
- ✅ 2/3 workflows migrados parcialmente (pr-to-dev, pr-to-main)
- ⚠️ 1/3 workflow no migrado (sync-main-to-dev)
- ✅ Job `lint` migrado a workflow reusable en ambos workflows
- ✅ Jobs custom mantenidos (tests, summary, security)

### Beneficios

1. **Centralización parcial** - Job lint centralizado
2. **Mantenibilidad mejorada** - Actualizar lint en un solo lugar
3. **Consistencia** - Mismo linting en todos los proyectos que usen el workflow reusable
4. **Sin regresión** - Funcionalidades personalizadas mantenidas

### Limitaciones

- **Reducción de código limitada** (~3-5% vs 75% esperado)
- **Migración parcial** - Solo jobs compatibles migrados
- **Dependencia de Makefile** - Impide migración completa de tests

---

## 🔄 Workflows Migrados

### 1. pr-to-dev.yml ✅ (Migración Híbrida)

**Estado:** Parcialmente migrado

#### Antes (154 líneas)
```yaml
jobs:
  unit-tests:  # Custom con Makefile
  lint:        # Custom con golangci-lint-action
  summary:     # Custom con github-script
```

#### Después (147 líneas)
```yaml
jobs:
  unit-tests:  # ⚠️ CUSTOM - Mantenido (usa Makefile)
  lint:        # ✅ MIGRADO - Workflow reusable
  summary:     # ⚠️ CUSTOM - Mantenido (comentarios personalizados)
```

**Reducción:** 4.5% (~7 líneas)

---

### 2. pr-to-main.yml ✅ (Migración Híbrida)

**Estado:** Parcialmente migrado

#### Antes (250 líneas)
```yaml
jobs:
  unit-tests:       # Custom con Makefile
  integration-tests:# Custom con Docker + Makefile
  lint:             # Custom con golangci-lint-action
  security-scan:    # Custom con Gosec
  summary:          # Custom con github-script
```

#### Después (242 líneas)
```yaml
jobs:
  unit-tests:       # ⚠️ CUSTOM - Mantenido
  integration-tests:# ⚠️ CUSTOM - Mantenido
  lint:             # ✅ MIGRADO - Workflow reusable
  security-scan:    # ⚠️ CUSTOM - Mantenido
  summary:          # ⚠️ CUSTOM - Mantenido
```

**Reducción:** 3.2% (~8 líneas)

---

### 3. sync-main-to-dev.yml ❌ (No Migrado)

**Estado:** NO migrado

**Razón:** Lógica específica del proyecto incompatible con workflow reusable:
- Lectura de versión desde archivo
- Verificación de diferencias
- Prevención de loops
- Mensajes personalizados
- Resumen detallado

**Ver decisión completa:** `docs/cicd/tracking/decisions/TASK-4.8-NO-MIGRATION.md`

---

## 🔧 Workflows Reusables Disponibles

Ubicación: `edugo-infrastructure/.github/workflows/reusable/`

### 1. go-lint.yml ✅ (En uso)

**Propósito:** Linting con golangci-lint

**Uso:**
```yaml
lint:
  name: Lint & Format Check
  uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-lint.yml@main
  with:
    go-version: "1.25"
    golangci-lint-version: "v2.4.0"
    args: "--timeout=5m"
  secrets:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

**Parámetros:**
- `go-version` - Versión de Go (default: 1.25)
- `golangci-lint-version` - Versión del linter (default: v1.64.7)
- `args` - Argumentos adicionales (default: --timeout=5m)
- `working-directory` - Directorio de trabajo (default: .)
- `skip-cache` - Saltar cache (default: false)

---

### 2. go-test.yml ⚠️ (Disponible, no usado)

**Propósito:** Tests unitarios/integración con coverage

**Por qué NO se usa:**
- api-mobile usa `make test-unit` y `make coverage-report`
- Workflow reusable usa comandos Go estándar
- Incompatible sin modificar proyecto

**Uso potencial (si se elimina Makefile):**
```yaml
test:
  name: Run Tests
  uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-test.yml@main
  with:
    go-version: "1.25"
    coverage-threshold: 33
    run-race: true
```

---

### 3. sync-branches.yml ⚠️ (Disponible, no usado)

**Propósito:** Sincronización básica main→dev

**Por qué NO se usa:**
- api-mobile necesita lógica específica (versión, loops, etc.)
- Workflow reusable es más simple
- Incompatible sin perder features

**Uso potencial:**
```yaml
sync:
  uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/sync-branches.yml@main
  with:
    source-branch: main
    target-branch: dev
```

---

### 4. docker-build.yml ⚠️ (Disponible, no usado)

**Propósito:** Build de imágenes Docker

**Nota:** api-mobile tiene workflows de release separados, no evaluado en SPRINT-4.

---

## 📚 Cómo Usar Workflows Reusables

### Estructura Básica

```yaml
jobs:
  job-name:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/WORKFLOW.yml@REFERENCIA
    with:
      parametro1: valor1
      parametro2: valor2
    secrets:
      SECRET_NAME: ${{ secrets.SECRET_NAME }}
```

### Componentes

1. **uses** - Ruta al workflow reusable
   - Formato: `org/repo/.github/workflows/file.yml@ref`
   - Ejemplo: `EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-lint.yml@main`

2. **with** - Parámetros de entrada
   - Definidos en `workflow_call.inputs` del workflow reusable
   - Todos son opcionales (tienen defaults)

3. **secrets** - Secrets a pasar
   - Necesario pasar `GITHUB_TOKEN` para workflows que lo requieran
   - Formato: `GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}`

### Ejemplo Completo

```yaml
name: PR to Dev

on:
  pull_request:
    branches: [dev]

jobs:
  lint:
    name: Lint Code
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-lint.yml@main
    with:
      go-version: "1.25"
      golangci-lint-version: "v2.4.0"
      args: "--timeout=5m"
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

---

## 🔀 Migración Híbrida Explicada

### ¿Qué es Migración Híbrida?

**Definición:** Migrar solo los componentes compatibles a workflows reusables, manteniendo lógica personalizada donde sea necesario.

### Ventajas

✅ **Mantiene funcionalidades** - No se pierde lógica personalizada
✅ **Sin cambios disruptivos** - Makefile, scripts, etc. siguen funcionando
✅ **Incremento gradual** - Se puede migrar más en el futuro
✅ **Reduce riesgo** - Cambios pequeños y controlados

### Desventajas

⚠️ **Menor reducción de código** - 3-5% vs 75% esperado
⚠️ **Migración incompleta** - Algunos jobs siguen siendo custom
⚠️ **Beneficios parciales** - Centralización solo en jobs migrados

### ¿Por qué es Necesaria?

**Razones en api-mobile:**

1. **Uso de Makefile**
   - `make test-unit`
   - `make coverage-report`
   - Workflows reusables usan `go test` directo

2. **Scripts personalizados**
   - `./scripts/check-coverage.sh`
   - Lógica específica del proyecto

3. **Comentarios automáticos en PRs**
   - github-script personalizado
   - Resúmenes detallados
   - No incluido en workflows reusables

4. **Lógica de negocio**
   - Lectura de versión
   - Prevención de loops
   - Manejo específico de errores

---

## 🔧 Mantenimiento

### Actualizar Versión de Go

**En api-mobile:**
```yaml
# .github/workflows/pr-to-dev.yml
env:
  GO_VERSION: "1.25"  # Actualizar aquí

# ...

lint:
  uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-lint.yml@main
  with:
    go-version: "1.25"  # Y aquí
```

**Mejor práctica:** Usar variable de entorno en todo el workflow.

### Actualizar Versión de golangci-lint

**Opción A:** Actualizar solo en api-mobile
```yaml
lint:
  uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-lint.yml@main
  with:
    golangci-lint-version: "v2.5.0"  # Nueva versión
```

**Opción B:** Actualizar default en infrastructure
- Editar `edugo-infrastructure/.github/workflows/reusable/go-lint.yml`
- Cambiar default de `golangci-lint-version`
- Afecta a TODOS los proyectos que usen el workflow

### Actualizar Workflow Reusable

**Cuando se actualiza workflow en infrastructure:**
1. Cambios se reflejan automáticamente (usa `@main`)
2. Para versión específica, usar tag: `@v1.0.0`
3. Probar en api-mobile antes de actualizar otros proyectos

---

## 🚨 Troubleshooting

### Problema: Job lint falla con error de permisos

**Síntoma:**
```
Error: Resource not accessible by integration
```

**Solución:**
Verificar que `GITHUB_TOKEN` se pasa correctamente:
```yaml
lint:
  uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-lint.yml@main
  secrets:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}  # ← Necesario
```

---

### Problema: Workflow reusable no encontrado

**Síntoma:**
```
Error: Unable to resolve action `EduGoGroup/edugo-infrastructure/...`
```

**Causas posibles:**
1. Workflow reusable no existe en la ruta especificada
2. Rama incorrecta (`@main` vs `@dev`)
3. Permisos de repositorio

**Solución:**
```yaml
# Verificar ruta correcta
uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-lint.yml@main
#     ^^^^ org        ^^^^ repo          ^^^^ ruta completa          ^^^^ ref
```

---

### Problema: Parámetros no funcionan

**Síntoma:**
Workflow reusable usa valores default en lugar de los especificados.

**Causa:**
Parámetros mal formateados o tipo incorrecto.

**Solución:**
```yaml
lint:
  uses: ...
  with:
    go-version: "1.25"              # String entrecomillado
    golangci-lint-version: "v2.4.0" # String entrecomillado
    args: "--timeout=5m"            # String entrecomillado
```

---

## 🚀 Próximos Pasos

### Migración Completa (FASE 2)

**Para lograr 70-80% de reducción de código:**

#### 1. Eliminar dependencia de Makefile
- Mover lógica de `make test-unit` a comandos Go directos
- Mover lógica de `make coverage-report` a scripts
- Documentar comandos equivalentes

#### 2. Estandarizar coverage check
- Usar composite action de infrastructure
- Eliminar `./scripts/check-coverage.sh`

#### 3. Crear composite action para comentarios
- Mover lógica de comentarios a action reutilizable
- Usar en todos los proyectos

#### 4. Migrar completamente
- Usar `go-test.yml` para unit-tests
- Usar `go-lint.yml` para lint (ya migrado)
- Reducción esperada: ~70-80%

---

### Extender Workflows Reusables (Sprint en Infrastructure)

**Para mejorar workflows reusables:**

#### 1. Agregar soporte para Makefile
```yaml
# En go-test.yml
inputs:
  use-makefile:
    type: boolean
    default: false
  makefile-target:
    type: string
    default: test-unit
```

#### 2. Agregar templates de comentarios
```yaml
# En go-test.yml
inputs:
  comment-template:
    type: string
    default: standard
```

#### 3. Extender sync-branches.yml
- Agregar lectura de versión
- Agregar verificación de diferencias
- Agregar prevención de loops
- Agregar templates de mensaje

---

### Replicar a Otros Proyectos

**Una vez validado en api-mobile:**

1. **edugo-api-administracion**
   - Misma estructura que api-mobile
   - Migración híbrida similar
   - Tiempo estimado: 4-6 horas

2. **edugo-worker**
   - Estructura diferente
   - Evaluar compatibilidad
   - Tiempo estimado: 6-8 horas

3. **Otros proyectos**
   - Evaluar caso por caso
   - Usar api-mobile como referencia

---

## 📚 Referencias

### Documentación del Sprint

| Documento | Propósito |
|-----------|-----------|
| `TASK-4.1-DISCOVERY.md` | Hallazgo workflows pre-existentes |
| `WORKFLOWS-REUSABLES-VALIDATION.md` | Validación completa workflows |
| `TASK-4.6-HYBRID-MIGRATION.md` | Decisión migración híbrida |
| `TASK-4.8-NO-MIGRATION.md` | Por qué sync-main-to-dev no se migró |
| `WORKFLOWS-SYNTAX-VALIDATION.md` | Validación sintaxis YAML |
| `TASKS-4.10-4.12-TESTING-STUB.md` | Plan de testing |
| `SPRINT-4-FASE-1-PROGRESS.md` | Reporte de progreso |

### Archivos Modificados

- `.github/workflows/pr-to-dev.yml` (147 líneas, -7 líneas)
- `.github/workflows/pr-to-main.yml` (242 líneas, -8 líneas)
- `.github/workflows/sync-main-to-dev.yml` (135 líneas, +7 líneas comentarios)

### Workflows Reusables

- `edugo-infrastructure/.github/workflows/reusable/go-lint.yml`
- `edugo-infrastructure/.github/workflows/reusable/go-test.yml`
- `edugo-infrastructure/.github/workflows/reusable/sync-branches.yml`
- `edugo-infrastructure/.github/workflows/reusable/docker-build.yml`

---

**Generado por:** Claude Code
**Fecha:** 2025-11-21
**Sprint:** SPRINT-4 completado
**Versión:** 1.0
