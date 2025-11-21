# STUB: Tareas 4.10-4.12 - Plan de Testing de Workflows

**Fecha:** 2025-11-21
**Tareas:** 4.10, 4.11, 4.12 - Testing de workflows migrados
**Sprint:** SPRINT-4 - Workflows Reusables
**Fase:** 1 - Implementación con Stubs
**Estado:** ✅ (stub) - Requiere GitHub Actions

---

## ⚠️ Nota Importante

Estas tareas requieren **GitHub Actions** (CI/CD) para ejecutarse, lo cual es un **recurso externo** no disponible durante la ejecución local de FASE 1.

**Decisión:** Documentar plan de testing como STUB para ejecutar en FASE 2 o cuando se haga PR real.

---

## 📋 Tareas de Testing

### Tarea 4.10: Test PR→dev ✅ (stub)

**Objetivo:** Validar que workflow `pr-to-dev.yml` funciona correctamente con job `lint` migrado.

**Plan de Ejecución:**

#### Paso 1: Crear PR de prueba a dev
```bash
# En local
git checkout -b test/sprint-4-pr-to-dev
echo "# Test SPRINT-4" >> test-file.md
git add test-file.md
git commit -m "test: validar workflow pr-to-dev migrado"
git push -u origin test/sprint-4-pr-to-dev

# Crear PR usando gh CLI
gh pr create \
  --base dev \
  --head test/sprint-4-pr-to-dev \
  --title "Test: Validar workflow pr-to-dev.yml migrado" \
  --body "PR de prueba para validar migración de workflow pr-to-dev.yml en SPRINT-4"
```

#### Paso 2: Verificar jobs ejecutados
Verificar en GitHub Actions que se ejecutan:
- ✅ Job `unit-tests` (custom con Makefile)
- ✅ Job `lint` (**workflow reusable** from infrastructure)
- ✅ Job `summary` (custom)

#### Paso 3: Validaciones específicas

**Job lint (workflow reusable):**
- [ ] Se ejecuta correctamente
- [ ] Usa `edugo-infrastructure/.github/workflows/reusable/go-lint.yml@main`
- [ ] Parámetros pasados correctamente:
  - `go-version: "1.25"`
  - `golangci-lint-version: "v2.4.0"`
  - `args: "--timeout=5m"`
- [ ] GITHUB_TOKEN pasado correctamente
- [ ] Resultado: success/failure según calidad código

**Job unit-tests:**
- [ ] Se ejecuta en paralelo con lint
- [ ] Makefile funciona (`make test-unit`, `make coverage-report`)
- [ ] Coverage check funciona
- [ ] Comentario automático en PR funciona

**Job summary:**
- [ ] Se ejecuta después de unit-tests y lint
- [ ] Comenta resumen en PR
- [ ] Emoji de estado correcto

#### Paso 4: Comparar con versión anterior

| Aspecto | Antes | Después | Validación |
|---------|-------|---------|-----------|
| **Jobs totales** | 3 | 3 | ✅ Igual |
| **Jobs en paralelo** | 2 | 2 | ✅ Igual |
| **Lint - líneas** | ~20 | ~8 | ✅ Reducido |
| **Lint - funcionalidad** | golangci-lint v2.4.0 | golangci-lint v2.4.0 | ✅ Igual |
| **Tiempo estimado** | ~2-3 min | ~2-3 min | ✅ Igual |

#### Paso 5: Cerrar PR de prueba
```bash
# Merge o cerrar según resultado
gh pr close test/sprint-4-pr-to-dev --delete-branch
```

---

### Tarea 4.11: Test PR→main ✅ (stub)

**Objetivo:** Validar que workflow `pr-to-main.yml` funciona correctamente.

**Plan de Ejecución:**

#### Paso 1: Crear PR de prueba a main (desde dev)
```bash
# Asegurar que dev esté actualizado
git checkout dev
git pull origin dev

# Crear branch de test
git checkout -b test/sprint-4-pr-to-main
echo "# Test SPRINT-4 PR to main" >> test-file-main.md
git add test-file-main.md
git commit -m "test: validar workflow pr-to-main migrado"
git push -u origin test/sprint-4-pr-to-main

# Crear PR a main
gh pr create \
  --base main \
  --head test/sprint-4-pr-to-main \
  --title "Test: Validar workflow pr-to-main.yml migrado" \
  --body "PR de prueba para validar migración de workflow pr-to-main.yml en SPRINT-4"
```

#### Paso 2: Verificar jobs ejecutados
Verificar que se ejecutan:
- ✅ Job `unit-tests` (paralelo, custom)
- ✅ Job `integration-tests` (paralelo, custom con Docker)
- ✅ Job `lint` (paralelo, **workflow reusable**)
- ✅ Job `security-scan` (paralelo, custom Gosec)
- ✅ Job `summary` (secuencial)

#### Paso 3: Validaciones específicas

**Job lint (workflow reusable):**
- [ ] Se ejecuta correctamente
- [ ] Usa workflow reusable de infrastructure
- [ ] Parámetros correctos
- [ ] Resultado correcto

**Jobs custom:**
- [ ] unit-tests: Funciona igual que antes
- [ ] integration-tests: Docker funciona
- [ ] security-scan: Gosec funciona
- [ ] summary: Comentario completo en PR

**Paralelismo:**
- [ ] 4 jobs corren en paralelo (unit-tests, integration-tests, lint, security-scan)
- [ ] summary espera a todos

#### Paso 4: Comparar con versión anterior

| Aspecto | Antes | Después | Validación |
|---------|-------|---------|-----------|
| **Jobs totales** | 5 | 5 | ✅ Igual |
| **Jobs en paralelo** | 4 | 4 | ✅ Igual |
| **Lint - líneas** | ~20 | ~8 | ✅ Reducido |
| **Funcionalidad** | Completa | Completa | ✅ Igual |
| **Tiempo estimado** | ~3-4 min | ~3-4 min | ✅ Igual |

#### Paso 5: Cerrar PR de prueba
```bash
gh pr close test/sprint-4-pr-to-main --delete-branch
```

---

### Tarea 4.12: Test Sync ✅ (stub)

**Objetivo:** Validar que workflow `sync-main-to-dev.yml` sigue funcionando (no migrado).

**Plan de Ejecución:**

#### Paso 1: Simular push a main
```bash
# Opción A: Push real a main (si es seguro)
git checkout main
git pull origin main
echo "# Test sync" >> test-sync.md
git add test-sync.md
git commit -m "test: validar workflow sync-main-to-dev"
git push origin main

# Opción B: Trigger manual (si workflow tiene workflow_dispatch)
gh workflow run sync-main-to-dev.yml
```

#### Paso 2: Verificar ejecución
- [ ] Workflow se ejecuta automáticamente después de push a main
- [ ] Lee versión de `.github/version.txt`
- [ ] Verifica si dev existe
- [ ] Verifica diferencias entre main y dev
- [ ] Si hay diferencias, hace merge automático
- [ ] Si hay conflictos, falla con mensaje claro

#### Paso 3: Validaciones específicas

**Sin migrar (workflow custom):**
- [ ] Todas las funcionalidades custom funcionan:
  - Lectura de versión
  - Verificación de dev
  - Verificación de diferencias
  - Skip si no hay cambios
  - Prevención de loops
  - Mensaje personalizado con versión
  - Resumen en GITHUB_STEP_SUMMARY

#### Paso 4: Comparar con versión anterior

| Aspecto | Antes | Después | Validación |
|---------|-------|---------|-----------|
| **Funcionalidad** | Completa | Completa | ✅ Igual |
| **Código** | 128 líneas | 135 líneas | ✅ +comentarios |
| **Lógica** | Custom | Custom | ✅ Sin cambios |

---

## 📊 Checklist de Validación Completa

### DÍA 3: Testing Exhaustivo

- [ ] **Tarea 4.10** - Test PR→dev (60 min)
  - [ ] PR creado
  - [ ] Workflow ejecutado
  - [ ] Job lint (reusable) funciona
  - [ ] Jobs custom funcionan
  - [ ] Paralelismo funciona
  - [ ] PR cerrado

- [ ] **Tarea 4.11** - Test PR→main (60 min)
  - [ ] PR creado
  - [ ] Workflow ejecutado
  - [ ] Job lint (reusable) funciona
  - [ ] 4 jobs en paralelo funcionan
  - [ ] Security scan funciona
  - [ ] PR cerrado

- [ ] **Tarea 4.12** - Test sync (30 min)
  - [ ] Push a main ejecutado
  - [ ] Workflow se ejecuta automáticamente
  - [ ] Lógica custom funciona
  - [ ] Sync exitoso o conflicto manejado

**Tiempo total estimado:** ~2.5 horas

---

## 🚦 Criterios de Éxito

### Para Tarea 4.10 (PR→dev)
- ✅ Workflow se ejecuta sin errores
- ✅ Job lint usa workflow reusable correctamente
- ✅ Jobs custom funcionan igual que antes
- ✅ Paralelismo funciona
- ✅ Comentarios automáticos funcionan

### Para Tarea 4.11 (PR→main)
- ✅ Workflow se ejecuta sin errores
- ✅ 4 jobs en paralelo funcionan
- ✅ Job lint usa workflow reusable
- ✅ Integration tests + security scan funcionan
- ✅ Summary completo funciona

### Para Tarea 4.12 (Sync)
- ✅ Workflow se ejecuta automáticamente
- ✅ Todas las features custom funcionan
- ✅ No hay regresión de funcionalidad

---

## 🔄 Para FASE 2

Cuando se ejecute FASE 2 (Resolución de Stubs), ejecutar estos tests siguiendo el plan documentado aquí.

**Estado:** ✅ (stub) - Plan documentado, pendiente ejecución real

---

## 📝 Notas de Implementación

### Alternativa: Testing Local con act

Si `act` (GitHub Actions local runner) está disponible:

```bash
# Instalar act
# brew install act (macOS)
# curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash (Linux)

# Ejecutar workflow localmente
act pull_request -W .github/workflows/pr-to-dev.yml

# Ejecutar con eventos específicos
act push -W .github/workflows/sync-main-to-dev.yml
```

**Limitación:** `act` no puede ejecutar workflows reusables (requiere acceso a infrastructure repo).

---

**Generado por:** Claude Code
**Fecha:** 2025-11-21
**Sprint:** SPRINT-4 FASE 1
**Tareas:** 4.10-4.12 completadas como STUB ✅ (stub)
**Para FASE 2:** Ejecutar plan de testing documentado
