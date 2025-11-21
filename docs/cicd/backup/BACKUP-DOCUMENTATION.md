# Backup de Workflows Original - Pre-Migración

**Fecha de Backup:** 2025-11-21
**Sprint:** SPRINT-4 - Workflows Reusables
**Fase:** 1 - Implementación
**Tarea:** 4.5 - Backup de workflows actuales

---

## Archivos Respaldados

| Archivo Original | Backup | Líneas | Tamaño |
|------------------|--------|--------|--------|
| `.github/workflows/pr-to-dev.yml` | `workflows-original/pr-to-dev.yml` | 154 | ~4.8 KB |
| `.github/workflows/pr-to-main.yml` | `workflows-original/pr-to-main.yml` | 250 | ~7.9 KB |
| `.github/workflows/sync-main-to-dev.yml` | `workflows-original/sync-main-to-dev.yml` | 128 | ~4.5 KB |
| **TOTAL** | | **532** | **~17.2 KB** |

---

## Estado Actual de Workflows

### 1. pr-to-dev.yml (154 líneas)

**Trigger:** Pull Request a `dev`

**Jobs:**
1. **unit-tests** (paralelo)
   - Setup Go 1.25
   - Download dependencies
   - Run tests con coverage
   - Check coverage >= 33%
   - Upload coverage artifact

2. **lint** (paralelo)
   - Setup Go 1.25
   - Run golangci-lint v2.4.0

**Características:**
- ✅ Paralelismo: 2 jobs sin `needs:`
- ✅ Coverage threshold: 33%
- ✅ Go version: 1.25
- ✅ Lint timeout: 5m
- ✅ Cache habilitado

---

### 2. pr-to-main.yml (250 líneas)

**Trigger:** Pull Request a `main`

**Jobs:**
1. **unit-tests** (paralelo)
   - Setup Go 1.25
   - Download dependencies
   - Run tests con coverage
   - Check coverage >= 33%
   - Upload coverage artifact

2. **integration-tests** (paralelo)
   - Setup Go 1.25
   - Setup Docker con testcontainers
   - Run integration tests
   - Upload logs como artifact

3. **lint** (paralelo)
   - Setup Go 1.25
   - Run golangci-lint v2.4.0

4. **security-scan** (paralelo)
   - Setup Go 1.25
   - Run Gosec
   - Upload SARIF results

**Características:**
- ✅ Paralelismo: 4 jobs sin `needs:`
- ✅ Coverage threshold: 33%
- ✅ Go version: 1.25
- ✅ Integration tests con Docker
- ✅ Security scan con Gosec
- ✅ Lint timeout: 5m
- ✅ Cache habilitado

---

### 3. sync-main-to-dev.yml (128 líneas)

**Trigger:**
- Push a `main`
- Creación de tags `v*`

**Job:**
- **sync**
  - Checkout con fetch-depth: 0
  - Configure Git
  - Attempt merge main → dev
  - Push si exitoso
  - Create PR si conflictos

**Características:**
- ✅ Auto-merge si no hay conflictos
- ✅ Creación automática de PR en conflictos
- ✅ Manejo de errores robusto

---

## Métricas Antes de Migración

| Métrica | Valor |
|---------|-------|
| **Total de líneas** | 532 |
| **Workflows totales** | 3 |
| **Jobs totales** | 6 (2 + 4 + 1) |
| **Jobs en paralelo PR→dev** | 2 |
| **Jobs en paralelo PR→main** | 4 |
| **Duplicación de código** | Alta (lógica repetida en pr-to-dev y pr-to-main) |
| **Mantenibilidad** | Media (cambios deben replicarse) |

---

## Objetivos de Migración

### Reducción de Código Esperada

| Workflow | Antes (líneas) | Después (estimado) | Reducción |
|----------|----------------|-------------------|-----------|
| `pr-to-dev.yml` | 154 | ~30-40 | ~74-79% |
| `pr-to-main.yml` | 250 | ~50-60 | ~76-80% |
| `sync-main-to-dev.yml` | 128 | ~15-20 | ~84-88% |
| **TOTAL** | **532** | **~95-120** | **~77-82%** |

### Mejoras Esperadas

1. **Mantenibilidad**
   - ✅ Cambios centralizados en infrastructure
   - ✅ Un solo lugar para actualizar versiones
   - ✅ Consistencia garantizada

2. **Reusabilidad**
   - ✅ Workflows reusables para otros proyectos
   - ✅ Patrón validado en api-mobile (PILOTO)
   - ✅ Fácil replicación a api-administracion, worker

3. **Flexibilidad**
   - ✅ Parámetros configurables
   - ✅ Composición modular de workflows
   - ✅ Fácil agregar/quitar validaciones

---

## Estrategia de Rollback

En caso de problemas durante o después de la migración:

### Opción A: Git Revert (Recomendado)
```bash
# Revertir commits de migración
git revert <commit-hash>
git push origin dev
```

### Opción B: Restaurar desde Backup
```bash
# Restaurar workflows originales
cp docs/cicd/backup/workflows-original/pr-to-dev.yml .github/workflows/
cp docs/cicd/backup/workflows-original/pr-to-main.yml .github/workflows/
cp docs/cicd/backup/workflows-original/sync-main-to-dev.yml .github/workflows/

git add .github/workflows/
git commit -m "chore: rollback workflows to pre-migration state"
git push origin dev
```

---

## Validaciones Post-Migración

Después de migrar, validar:

- [ ] pr-to-dev.yml sintaxis correcta (yamllint)
- [ ] pr-to-main.yml sintaxis correcta
- [ ] sync-main-to-dev.yml sintaxis correcta
- [ ] Crear PR de prueba a dev
- [ ] Verificar que todos los jobs corren
- [ ] Verificar paralelismo mantenido
- [ ] Verificar coverage threshold funciona
- [ ] Verificar lint funciona
- [ ] Verificar security scan funciona (PR a main)
- [ ] Verificar sync funciona (push a main)

---

**Generado por:** Claude Code
**Fecha:** 2025-11-21
**Sprint:** SPRINT-4 FASE 1
**Tarea:** 4.5 completada ✅
