# Decisión: Tarea 4.6 - Migración Híbrida de pr-to-dev.yml

**Fecha:** 2025-11-21
**Tarea:** 4.6 - Migrar pr-to-dev.yml
**Sprint:** SPRINT-4 - Workflows Reusables
**Fase:** 1 - Implementación

---

## Contexto

Al intentar migrar `pr-to-dev.yml` a usar workflows reusables de infrastructure, se identificaron características personalizadas del proyecto:

### Características Actuales de pr-to-dev.yml

1. **Usa Makefile:**
   - `make test-unit`
   - `make coverage-report`

2. **Scripts Personalizados:**
   - `./scripts/check-coverage.sh`

3. **Funcionalidades Específicas:**
   - Comentarios automáticos de coverage en PR
   - Job de summary con resumen de checks
   - Label `skip-coverage` para bypass temporal

4. **Configuración:**
   - GO_VERSION: 1.25
   - COVERAGE_THRESHOLD: 33
   - Timeouts específicos

### Workflows Reusables en Infrastructure

Los workflows reusables (`go-test.yml`, `go-lint.yml`) usan:
- Comandos estándar de Go (`go test`, `golangci-lint`)
- No usan Makefile
- No tienen lógica de comentarios en PR (está en composite action)

---

## Análisis de Compatibilidad

| Feature | pr-to-dev.yml Actual | Workflow Reusable | Compatible |
|---------|---------------------|-------------------|-----------|
| **Tests** | `make test-unit` | `go test` | ❌ No directamente |
| **Coverage** | `make coverage-report` + script custom | `go tool cover` + composite action | ⚠️ Parcial |
| **Lint** | `golangci-lint v2.4.0` | `golangci-lint` configurable | ✅ Sí |
| **Comments** | github-script custom | No incluido | ❌ No |
| **Summary** | Job separado custom | No incluido | ❌ No |

---

## Opciones Evaluadas

### Opción A: Migración Completa (Descartada)
Reemplazar todo pr-to-dev.yml con llamadas a workflows reusables.

**Pros:**
- Máxima reducción de código
- Centralización total

**Contras:**
- ❌ Pierde funcionalidades (comentarios, summary)
- ❌ Requiere eliminar Makefile (cambio disruptivo)
- ❌ Requiere eliminar scripts custom
- ❌ Regresión de UX (sin comentarios automáticos)

**Decisión:** ❌ Rechazada

---

### Opción B: Migración Híbrida (SELECCIONADA)
Migrar solo componentes compatibles, mantener lógica personalizada.

**Estrategia:**
1. ✅ Migrar job `lint` → usar `go-lint.yml` reusable
2. ⚠️ Mantener job `unit-tests` custom (por Makefile)
3. ✅ Mantener job `summary` custom

**Pros:**
- ✅ Reduce algo de código (job lint)
- ✅ Mantiene funcionalidades personalizadas
- ✅ Sin cambios disruptivos
- ✅ Incremento gradual de reusabilidad

**Contras:**
- ⚠️ Solo ~20-30% de reducción (vs 74% esperado)
- ⚠️ Migración parcial

**Decisión:** ✅ Seleccionada (para FASE 1)

---

### Opción C: Adaptar Proyecto + Migración Completa (FASE 2)
Modificar proyecto para usar comandos estándar, luego migrar completamente.

**Tareas requeridas:**
1. Eliminar/adaptar Makefile para usar `go test` directo
2. Eliminar `./scripts/check-coverage.sh` (usar composite action de infrastructure)
3. Adaptar lógica de comentarios a composite action reutilizable

**Pros:**
- ✅ Migración completa posible
- ✅ Proyecto más estándar
- ✅ Máxima reducción de código

**Contras:**
- ⚠️ Requiere cambios en proyecto (no solo workflows)
- ⚠️ Testing extensivo necesario
- ⚠️ Fuera del alcance de FASE 1

**Decisión:** ⏳ Pospuesto para FASE 2 o sprint futuro

---

## Implementación de Opción B (Migración Híbrida)

### Archivo Migrado: pr-to-dev.yml

```yaml
name: PR to Dev - Unit Tests

on:
  pull_request:
    branches: [dev]
    types: [opened, synchronize, reopened]

env:
  GO_VERSION: "1.25"
  COVERAGE_THRESHOLD: 33

jobs:
  # =====================================================
  # MIGRADO: Job lint usando workflow reusable
  # =====================================================
  lint:
    name: Lint & Format Check
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/go-lint.yml@main
    with:
      go-version: "1.25"
      golangci-lint-version: "v2.4.0"
      args: "--timeout=5m"
    secrets:
      GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

  # =====================================================
  # MANTENIDO: Job tests custom (usa Makefile)
  # =====================================================
  unit-tests:
    name: Unit Tests
    runs-on: ubuntu-latest
    timeout-minutes: 10

    steps:
      - name: 📥 Checkout código
        uses: actions/checkout@v4

      - name: 🔧 Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true

      - name: 🔐 Configurar acceso a repos privados
        run: |
          git config --global url."https://${{ secrets.GITHUB_TOKEN }}@github.com/".insteadOf "https://github.com/"
        env:
          GOPRIVATE: github.com/EduGoGroup/*

      - name: 📦 Descargar dependencias
        run: go mod download

      - name: 🧪 Ejecutar tests unitarios
        run: make test-unit
        timeout-minutes: 5

      - name: 📊 Generar reporte de cobertura
        run: make coverage-report
        timeout-minutes: 5

      - name: ✅ Verificar umbral de cobertura
        if: |
          !contains(github.event.pull_request.labels.*.name, 'skip-coverage')
        run: |
          ./scripts/check-coverage.sh coverage/coverage-filtered.out ${{ env.COVERAGE_THRESHOLD }} || {
            echo "::warning::Cobertura por debajo del umbral de ${COVERAGE_THRESHOLD}%"
            echo "💡 Tip: Agrega label 'skip-coverage' al PR si es temporal"
            exit 1
          }
        continue-on-error: false

      - name: 📤 Subir reporte de cobertura
        uses: actions/upload-artifact@v4
        if: always()
        with:
          name: coverage-report-unit
          path: coverage/
          retention-days: 7

      - name: 📈 Comentar cobertura en PR
        uses: actions/github-script@v7
        if: always()
        with:
          script: |
            const fs = require('fs');
            const coverage = fs.readFileSync('coverage/coverage-filtered.out', 'utf8');
            const lines = coverage.split('\n');
            const totalLine = lines[lines.length - 2];
            const match = totalLine.match(/(\d+\.\d+)%/);
            const coveragePercent = match ? match[1] : 'N/A';

            const comment = `## 📊 Cobertura de Tests Unitarios

            **Cobertura Total**: ${coveragePercent}%
            **Umbral Mínimo**: ${process.env.COVERAGE_THRESHOLD}%

            ${parseFloat(coveragePercent) >= parseFloat(process.env.COVERAGE_THRESHOLD) ? '✅ Cobertura cumple con el umbral' : '⚠️ Cobertura por debajo del umbral'}

            📄 [Ver reporte completo](https://github.com/${context.repo.owner}/${context.repo.repo}/actions/runs/${context.runId})
            `;

            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: comment
            });

  # =====================================================
  # MANTENIDO: Job summary custom
  # =====================================================
  summary:
    name: PR Summary
    runs-on: ubuntu-latest
    needs: [unit-tests, lint]
    if: always()

    steps:
      - name: 📋 Generar resumen
        uses: actions/github-script@v7
        with:
          script: |
            const unitTests = '${{ needs.unit-tests.result }}';
            const lint = '${{ needs.lint.result }}';

            const statusEmoji = (status) => {
              switch(status) {
                case 'success': return '✅';
                case 'failure': return '❌';
                case 'cancelled': return '⏸️';
                default: return '⚠️';
              }
            };

            const summary = `## 🔍 Resumen de Checks - PR a Dev

            | Check | Estado |
            |-------|--------|
            | Tests Unitarios | ${statusEmoji(unitTests)} ${unitTests} |
            | Lint & Format | ${statusEmoji(lint)} ${lint} |

            ${unitTests === 'success' && lint === 'success' ? '✅ **Todos los checks pasaron** - PR listo para review' : '⚠️ **Algunos checks fallaron** - Por favor revisa los errores'}
            `;

            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: summary
            });
```

---

## Métricas de Migración Híbrida

| Métrica | Antes | Después | Reducción |
|---------|-------|---------|-----------|
| **Líneas totales** | 154 | ~140 | ~9% |
| **Jobs migrados** | 1/3 (lint) | - | 33% |
| **Jobs custom** | 2/3 (unit-tests, summary) | - | 67% |

**Nota:** Reducción menor a esperada (74%) pero mantiene funcionalidades.

---

## Plan para FASE 2 o Sprint Futuro

### Tareas para Migración Completa

1. **Eliminar dependencia de Makefile:**
   - Mover lógica de `make test-unit` a comandos Go directos
   - Documentar comandos equivalentes

2. **Estandarizar coverage check:**
   - Usar composite action de infrastructure
   - Eliminar `./scripts/check-coverage.sh`

3. **Crear composite action para comentarios:**
   - Mover lógica de comentarios a action reutilizable
   - Usar en todos los proyectos

4. **Migrar completamente:**
   - Usar `go-test.yml` para unit-tests
   - Usar `go-lint.yml` para lint
   - Reducción esperada: ~70-80%

---

## Conclusiones

✅ **Migración híbrida es el enfoque correcto para FASE 1**
✅ **Mantiene funcionalidades personalizadas**
✅ **Reduce algo de código (job lint)**
⏳ **Migración completa pospuesta a FASE 2**

**Próximos pasos:**
1. Implementar pr-to-dev.yml híbrido
2. Documentar en SPRINT-STATUS.md
3. Continuar con pr-to-main.yml (similar estrategia)

---

**Generado por:** Claude Code
**Fecha:** 2025-11-21
**Sprint:** SPRINT-4 FASE 1
**Tarea:** 4.6 en progreso
