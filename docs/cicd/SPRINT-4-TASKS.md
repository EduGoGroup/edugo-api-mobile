# Sprint 4: Workflows Reusables - edugo-api-mobile

**Duración:** 3-4 días  
**Objetivo:** Crear workflows reusables y migrar api-mobile como PILOTO  
**Estado:** Listo para Ejecución (Post Sprint 2)  
**Proyecto:** edugo-api-mobile (PILOTO)

---

## 📋 Resumen del Sprint

| Métrica | Objetivo |
|---------|----------|
| **Tareas Totales** | 12 |
| **Tiempo Estimado** | 12-15 horas |
| **Prioridad** | 🟢 P2 (Media) |
| **Prerequisito** | Sprint 2 completado ✅ |
| **Commits Esperados** | 4-6 |
| **PRs a Crear** | 2 PRs (infrastructure + api-mobile) |
| **Riesgo** | 🟡 Medio |

---

## ⚠️ PREREQUISITOS CRÍTICOS

### Antes de Comenzar Este Sprint

- ✅ **Sprint 2 completado y mergeado**
  - Go 1.25 funcionando en CI
  - Paralelismo implementado
  - Pre-commit hooks configurados
  - Errores lint corregidos

- ✅ **Validaciones previas exitosas**
  - Success rate api-mobile >95%
  - Todos los workflows pasando
  - Sin errores en CI

- ✅ **Conocimiento del estado actual**
  - 5 workflows en api-mobile
  - Código duplicado identificado
  - Patrón de reusables claro

### Si Prerequisitos NO Cumplen

**NO INICIAR Sprint 4**. Volver a Sprint 2 y completar pendientes.

---

## 🎯 Objetivos del Sprint

### Objetivo Principal

Crear workflows reusables centralizados en `edugo-infrastructure` y migrar `edugo-api-mobile` como proyecto PILOTO para validar el patrón antes de replicar a otros proyectos.

### Objetivos Específicos

1. **Crear workflows reusables base** en infrastructure
   - `pr-validation.yml` (para PR→dev y PR→main)
   - `sync-branches.yml` (para sincronización main→dev)

2. **Migrar api-mobile** a usar workflows reusables
   - Convertir `pr-to-dev.yml` a llamar reusable
   - Convertir `pr-to-main.yml` a llamar reusable
   - Convertir `sync-main-to-dev.yml` a llamar reusable
   - Mantener `manual-release.yml` personalizado
   - Mantener `test.yml` personalizado

3. **Validar patrón** exhaustivamente
   - Tests en todos los escenarios
   - Validar que funciona igual o mejor
   - Documentar aprendizajes

4. **Documentar para replicar**
   - Guía de uso de workflows reusables
   - Ejemplos de personalización
   - Troubleshooting común

---

## 🗓️ Cronograma Diario

### Día 1: Crear Workflows Reusables Base (4h)
- ✅ Tarea 4.1: Setup en infrastructure (30 min)
- ✅ Tarea 4.2: Crear pr-validation.yml reusable (90 min) 🟢 P2
- ✅ Tarea 4.3: Crear sync-branches.yml reusable (60 min) 🟢 P2
- ✅ Tarea 4.4: Validar sintaxis y documentar (60 min)

### Día 2: Migrar api-mobile (4h)
- ✅ Tarea 4.5: Preparación y backup (30 min)
- ✅ Tarea 4.6: Convertir pr-to-dev.yml (60 min) 🟢 P2
- ✅ Tarea 4.7: Convertir pr-to-main.yml (60 min) 🟢 P2
- ✅ Tarea 4.8: Convertir sync-main-to-dev.yml (45 min)
- ✅ Tarea 4.9: Validar workflows localmente (45 min)

### Día 3: Testing Exhaustivo (3h)
- ✅ Tarea 4.10: Tests de PR→dev (60 min) 🟢 P2
- ✅ Tarea 4.11: Tests de PR→main (60 min) 🟢 P2
- ✅ Tarea 4.12: Tests de sync (30 min)

### Día 4: Documentación y Cierre (2h)
- ✅ Tarea 4.13: Documentación completa (60 min)
- ✅ Tarea 4.14: Métricas y comparación (30 min)
- ✅ Tarea 4.15: PR y merge (30 min)

---

## 📐 Arquitectura de Workflows Reusables

### Concepto

```
┌─────────────────────────────────────┐
│  edugo-infrastructure               │
│  .github/workflows/                 │
│                                     │
│  ├── pr-validation.yml ←──────┐    │
│  │   (reusable workflow)      │    │
│  │                            │    │
│  └── sync-branches.yml ←──────┼──┐ │
│      (reusable workflow)      │  │ │
└──────────────────────────────┼┼──┼─┘
                               ││  │
                               ││  │
┌──────────────────────────────┼┼──┼─┐
│  edugo-api-mobile            ││  │ │
│  .github/workflows/          ││  │ │
│                              ││  │ │
│  ├── pr-to-dev.yml ──────────┘│  │ │
│  │   (calls pr-validation)    │  │ │
│  │                             │  │ │
│  ├── pr-to-main.yml ───────────┘  │ │
│  │   (calls pr-validation)        │ │
│  │                                │ │
│  ├── sync-main-to-dev.yml ────────┘ │
│  │   (calls sync-branches)          │
│  │                                  │
│  ├── manual-release.yml             │
│  │   (personalizado, NO reusable)  │
│  │                                  │
│  └── test.yml                       │
│      (personalizado, NO reusable)  │
└─────────────────────────────────────┘
```

### Ventajas

1. **Centralización:** Workflows en un solo lugar
2. **Mantenibilidad:** Cambio en 1 lugar → afecta a todos
3. **Consistencia:** Mismo comportamiento en todos los proyectos
4. **Reducción de código:** -60% duplicación
5. **Escalabilidad:** Fácil agregar nuevos proyectos

---

## 📝 TAREAS DETALLADAS

---

## DÍA 1: CREAR WORKFLOWS REUSABLES BASE

---

### ✅ Tarea 4.1: Setup en Infrastructure

**Prioridad:** 🟢 P2  
**Estimación:** ⏱️ 30 minutos  
**Prerequisitos:** Acceso a edugo-infrastructure

#### Objetivos
- Clonar/actualizar edugo-infrastructure
- Crear estructura de workflows reusables
- Preparar documentación base
- Crear rama de trabajo

#### Pasos a Ejecutar

```bash
#!/bin/bash
# setup-infrastructure-reusables.sh

INFRA_PATH="/Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure"

echo "🚀 Preparando infrastructure para workflows reusables..."

# Verificar si existe el repo
if [ ! -d "$INFRA_PATH" ]; then
  echo "❌ edugo-infrastructure no existe en: $INFRA_PATH"
  echo "   Clonarlo con: gh repo clone EduGoGroup/edugo-infrastructure $INFRA_PATH"
  exit 1
fi

cd "$INFRA_PATH"

echo "📥 Actualizando dev..."
git checkout dev
git pull origin dev

echo "🔍 Verificando estado..."
if [ -n "$(git status --porcelain)" ]; then
  echo "⚠️  Hay cambios pendientes, guardando stash..."
  git stash save "WIP antes de Sprint 4"
fi

echo "🌿 Creando rama de trabajo..."
BRANCH="feature/cicd-reusable-workflows"
git checkout -b "$BRANCH"

echo "📁 Creando estructura de directorios..."
mkdir -p .github/workflows/reusable
mkdir -p docs/workflows

echo "📝 Creando README de workflows reusables..."
cat > .github/workflows/reusable/README.md << 'README'
# Workflows Reusables - EduGo

Este directorio contiene workflows reusables para el ecosistema EduGo.

## Workflows Disponibles

### 1. pr-validation.yml
Validación estándar para Pull Requests (dev o main).

**Uso:**
```yaml
jobs:
  validate:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/pr-validation.yml@main
    with:
      go-version: "1.25"
      coverage-threshold: 33
```

### 2. sync-branches.yml
Sincronización automática de main → dev.

**Uso:**
```yaml
jobs:
  sync:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/sync-branches.yml@main
```

## Cómo Usar

1. En tu proyecto, crea un workflow que llame al reusable
2. Pasa los parámetros necesarios con `with:`
3. Opcionalmente pasa secrets con `secrets:`

Ver ejemplos en: `edugo-api-mobile/.github/workflows/`

## Contribuir

Para agregar o modificar workflows reusables:
1. Crear branch en infrastructure
2. Agregar/modificar workflow en `.github/workflows/reusable/`
3. Probar en proyecto piloto (api-mobile)
4. Documentar cambios aquí
5. Crear PR

README

echo "✅ Estructura creada"

echo ""
echo "🎉 Setup completado!"
echo ""
echo "📋 Resumen:"
echo "  - Rama: $BRANCH"
echo "  - Directorio: .github/workflows/reusable/"
echo "  - README creado"
echo ""
echo "🚀 Siguiente paso: Tarea 4.2 (Crear pr-validation.yml)"
```

#### Guardar y Ejecutar

```bash
# Guardar script
cat > /Users/jhoanmedina/source/EduGo/Analisys/00-Projects-Isolated/cicd-analysis/implementation-plans/03-api-mobile/SCRIPTS/setup-infrastructure-reusables.sh << 'SCRIPT'
# ... (copiar script de arriba)
SCRIPT

chmod +x /Users/jhoanmedina/source/EduGo/Analisys/00-Projects-Isolated/cicd-analysis/implementation-plans/03-api-mobile/SCRIPTS/setup-infrastructure-reusables.sh

# Ejecutar
/Users/jhoanmedina/source/EduGo/Analisys/00-Projects-Isolated/cicd-analysis/implementation-plans/03-api-mobile/SCRIPTS/setup-infrastructure-reusables.sh
```

#### Criterios de Validación

- ✅ Rama `feature/cicd-reusable-workflows` creada
- ✅ Directorio `.github/workflows/reusable/` existe
- ✅ README.md creado
- ✅ Working tree limpio

#### Checkpoint

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure
git branch --show-current  # feature/cicd-reusable-workflows
ls -la .github/workflows/reusable/  # Debe existir
cat .github/workflows/reusable/README.md  # Debe tener contenido
```

---

### ✅ Tarea 4.2: Crear pr-validation.yml Reusable

**Prioridad:** 🟢 P2  
**Estimación:** ⏱️ 90 minutos  
**Prerequisitos:** Tarea 4.1 completada

#### Objetivos
- Crear workflow reusable para validación de PRs
- Soportar personalización (Go version, coverage, etc.)
- Incluir lint, test, build-docker
- Documentar parámetros

#### Workflow Reusable

```yaml
# .github/workflows/reusable/pr-validation.yml
name: PR Validation (Reusable)

on:
  workflow_call:
    inputs:
      go-version:
        description: 'Go version to use'
        required: false
        type: string
        default: '1.25'
      
      coverage-threshold:
        description: 'Minimum coverage percentage'
        required: false
        type: number
        default: 33
      
      enable-docker-build:
        description: 'Enable Docker build validation'
        required: false
        type: boolean
        default: true
      
      enable-security-scan:
        description: 'Enable security scan (Gosec)'
        required: false
        type: boolean
        default: false
      
      docker-platforms:
        description: 'Docker platforms to build'
        required: false
        type: string
        default: 'linux/amd64'
      
      working-directory:
        description: 'Working directory for commands'
        required: false
        type: string
        default: '.'
    
    secrets:
      GITHUB_TOKEN:
        required: false

jobs:
  lint:
    name: Lint
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ inputs.go-version }}
          cache: true
      
      - name: Run golangci-lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: latest
          working-directory: ${{ inputs.working-directory }}
          args: --timeout=5m

  test:
    name: Test
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ inputs.go-version }}
          cache: true
      
      - name: Download dependencies
        working-directory: ${{ inputs.working-directory }}
        run: go mod download
      
      - name: Run tests
        working-directory: ${{ inputs.working-directory }}
        run: go test -v -race -coverprofile=coverage.out ./...
      
      - name: Check coverage
        working-directory: ${{ inputs.working-directory }}
        run: |
          COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          echo "Coverage: $COVERAGE%"
          echo "Threshold: ${{ inputs.coverage-threshold }}%"
          
          if [ $(echo "$COVERAGE < ${{ inputs.coverage-threshold }}" | bc) -eq 1 ]; then
            echo "❌ Coverage $COVERAGE% is below threshold ${{ inputs.coverage-threshold }}%"
            exit 1
          fi
          echo "✅ Coverage OK: $COVERAGE%"
      
      - name: Upload coverage report
        uses: actions/upload-artifact@v4
        with:
          name: coverage-report
          path: ${{ inputs.working-directory }}/coverage.out

  security-scan:
    name: Security Scan
    runs-on: ubuntu-latest
    if: ${{ inputs.enable-security-scan }}
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ inputs.go-version }}
          cache: true
      
      - name: Run Gosec
        uses: securego/gosec@master
        with:
          args: '-no-fail -fmt sarif -out results.sarif ${{ inputs.working-directory }}/...'
      
      - name: Upload SARIF file
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: results.sarif

  build-docker:
    name: Build Docker
    runs-on: ubuntu-latest
    if: ${{ inputs.enable-docker-build }}
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4
      
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3
      
      - name: Build Docker image
        uses: docker/build-push-action@v5
        with:
          context: ${{ inputs.working-directory }}
          push: false
          platforms: ${{ inputs.docker-platforms }}
          tags: |
            ${{ github.event.repository.name }}:pr-${{ github.event.pull_request.number }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
```

#### Script de Creación

```bash
#!/bin/bash
# create-pr-validation-reusable.sh

INFRA_PATH="/Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure"
REUSABLE_FILE=".github/workflows/reusable/pr-validation.yml"

cd "$INFRA_PATH"

echo "🚀 Creando pr-validation.yml reusable..."

# Crear el workflow (copiar contenido de arriba)
cat > "$REUSABLE_FILE" << 'WORKFLOW'
# ... (copiar YAML completo de arriba)
WORKFLOW

echo "✅ Workflow reusable creado"

# Validar sintaxis
if command -v yamllint &> /dev/null; then
  yamllint "$REUSABLE_FILE"
  echo "✅ Sintaxis YAML válida"
else
  echo "⚠️  yamllint no instalado"
fi

# Documentar
cat > docs/workflows/pr-validation.md << 'DOC'
# PR Validation Workflow (Reusable)

Workflow reusable para validación de Pull Requests.

## Características

- ✅ Lint con golangci-lint
- ✅ Tests unitarios + integración
- ✅ Coverage con threshold configurable
- ✅ Security scan opcional (Gosec)
- ✅ Docker build opcional
- ✅ Paralelismo (lint, test, security, build en paralelo)

## Parámetros

| Parámetro | Tipo | Default | Descripción |
|-----------|------|---------|-------------|
| `go-version` | string | `1.25` | Versión de Go |
| `coverage-threshold` | number | `33` | Cobertura mínima (%) |
| `enable-docker-build` | boolean | `true` | Habilitar build Docker |
| `enable-security-scan` | boolean | `false` | Habilitar Gosec |
| `docker-platforms` | string | `linux/amd64` | Plataformas Docker |
| `working-directory` | string | `.` | Directorio de trabajo |

## Ejemplo de Uso

### Uso Básico

```yaml
name: PR to Dev

on:
  pull_request:
    branches: [dev]

jobs:
  validate:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/pr-validation.yml@main
    with:
      go-version: "1.25"
```

### Uso con Security Scan

```yaml
name: PR to Main

on:
  pull_request:
    branches: [main]

jobs:
  validate:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/pr-validation.yml@main
    with:
      go-version: "1.25"
      coverage-threshold: 33
      enable-security-scan: true
      enable-docker-build: true
      docker-platforms: linux/amd64,linux/arm64
```

### Uso con Directorio Específico

```yaml
jobs:
  validate:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/pr-validation.yml@main
    with:
      go-version: "1.25"
      working-directory: ./services/api
```

## Outputs

El workflow produce:
- Coverage report como artifact
- Logs de lint, test, security
- Resultado de build Docker

## Troubleshooting

### Problema: Tests fallan
```bash
# Ver logs del job test
gh run view --job=test --log
```

### Problema: Coverage por debajo de threshold
```bash
# Ajustar threshold temporalmente
with:
  coverage-threshold: 25  # Reducir si es necesario
```

### Problema: Docker build falla
```bash
# Desactivar si no necesitas
with:
  enable-docker-build: false
```

DOC

echo "📝 Documentación creada"

echo ""
echo "🎉 pr-validation.yml reusable completado!"
echo ""
echo "📋 Archivos creados:"
echo "  - $REUSABLE_FILE"
echo "  - docs/workflows/pr-validation.md"
echo ""
echo "🚀 Siguiente paso: Tarea 4.3 (Crear sync-branches.yml)"
```

#### Guardar y Ejecutar

```bash
# Guardar script
cat > /Users/jhoanmedina/source/EduGo/Analisys/00-Projects-Isolated/cicd-analysis/implementation-plans/03-api-mobile/SCRIPTS/create-pr-validation-reusable.sh << 'SCRIPT'
# ... (copiar script de arriba)
SCRIPT

chmod +x /Users/jhoanmedina/source/EduGo/Analisys/00-Projects-Isolated/cicd-analysis/implementation-plans/03-api-mobile/SCRIPTS/create-pr-validation-reusable.sh

# Ejecutar
/Users/jhoanmedina/source/EduGo/Analisys/00-Projects-Isolated/cicd-analysis/implementation-plans/03-api-mobile/SCRIPTS/create-pr-validation-reusable.sh
```

#### Commitear

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

git add .
git commit -m "feat: agregar workflow reusable pr-validation

Workflow reusable para validación de Pull Requests.

Características:
- Lint con golangci-lint
- Tests con coverage threshold
- Security scan opcional (Gosec)
- Docker build opcional
- Paralelismo total

Parámetros configurables:
- go-version (default: 1.25)
- coverage-threshold (default: 33)
- enable-docker-build (default: true)
- enable-security-scan (default: false)
- docker-platforms (default: linux/amd64)
- working-directory (default: .)

Documentación: docs/workflows/pr-validation.md

🤖 Generated with Claude Code

Co-Authored-By: Claude <noreply@anthropic.com>"
```

#### Criterios de Validación

- ✅ Archivo `.github/workflows/reusable/pr-validation.yml` creado
- ✅ Sintaxis YAML válida
- ✅ Documentación en `docs/workflows/pr-validation.md`
- ✅ Todos los inputs documentados
- ✅ Commit creado

---

### ✅ Tarea 4.3: Crear sync-branches.yml Reusable

**Prioridad:** 🟢 P2  
**Estimación:** ⏱️ 60 minutos  
**Prerequisitos:** Tarea 4.2 completada

#### Objetivos
- Crear workflow reusable para sincronización main→dev
- Manejar conflictos gracefully
- Crear branch dev si no existe
- Documentar uso

#### Workflow Reusable

```yaml
# .github/workflows/reusable/sync-branches.yml
name: Sync Branches (Reusable)

on:
  workflow_call:
    inputs:
      source-branch:
        description: 'Source branch to sync from'
        required: false
        type: string
        default: 'main'
      
      target-branch:
        description: 'Target branch to sync to'
        required: false
        type: string
        default: 'dev'
      
      create-target-if-missing:
        description: 'Create target branch if it does not exist'
        required: false
        type: boolean
        default: true
    
    secrets:
      GITHUB_TOKEN:
        required: false

jobs:
  sync:
    name: Sync ${{ inputs.source-branch }} → ${{ inputs.target-branch }}
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4
        with:
          fetch-depth: 0
          token: ${{ secrets.GITHUB_TOKEN || github.token }}
      
      - name: Configure Git
        run: |
          git config user.name "github-actions[bot]"
          git config user.email "github-actions[bot]@users.noreply.github.com"
      
      - name: Check if target branch exists
        id: check-branch
        run: |
          if git show-ref --verify --quiet "refs/remotes/origin/${{ inputs.target-branch }}"; then
            echo "exists=true" >> $GITHUB_OUTPUT
            echo "✅ Branch ${{ inputs.target-branch }} exists"
          else
            echo "exists=false" >> $GITHUB_OUTPUT
            echo "⚠️  Branch ${{ inputs.target-branch }} does not exist"
          fi
      
      - name: Create target branch if missing
        if: steps.check-branch.outputs.exists == 'false' && inputs.create-target-if-missing
        run: |
          echo "Creating ${{ inputs.target-branch }} from ${{ inputs.source-branch }}..."
          git checkout -b ${{ inputs.target-branch }} origin/${{ inputs.source-branch }}
          git push origin ${{ inputs.target-branch }}
          echo "✅ Branch ${{ inputs.target-branch }} created"
      
      - name: Sync branches
        if: steps.check-branch.outputs.exists == 'true' || inputs.create-target-if-missing
        run: |
          echo "Syncing ${{ inputs.source-branch }} → ${{ inputs.target-branch }}..."
          
          git checkout ${{ inputs.target-branch }}
          git pull origin ${{ inputs.target-branch }}
          
          # Intentar merge
          if git merge --no-edit origin/${{ inputs.source-branch }}; then
            echo "✅ Merge successful"
            git push origin ${{ inputs.target-branch }}
            echo "✅ Sync completed successfully"
          else
            echo "❌ Merge conflict detected"
            git merge --abort
            
            echo "## ⚠️ Sync Failed - Manual Resolution Required" >> $GITHUB_STEP_SUMMARY
            echo "" >> $GITHUB_STEP_SUMMARY
            echo "Sync from \`${{ inputs.source-branch }}\` to \`${{ inputs.target-branch }}\` failed due to merge conflicts." >> $GITHUB_STEP_SUMMARY
            echo "" >> $GITHUB_STEP_SUMMARY
            echo "### How to Resolve" >> $GITHUB_STEP_SUMMARY
            echo "" >> $GITHUB_STEP_SUMMARY
            echo "\`\`\`bash" >> $GITHUB_STEP_SUMMARY
            echo "git checkout ${{ inputs.target-branch }}" >> $GITHUB_STEP_SUMMARY
            echo "git pull origin ${{ inputs.target-branch }}" >> $GITHUB_STEP_SUMMARY
            echo "git merge origin/${{ inputs.source-branch }}" >> $GITHUB_STEP_SUMMARY
            echo "# Resolver conflictos manualmente" >> $GITHUB_STEP_SUMMARY
            echo "git add ." >> $GITHUB_STEP_SUMMARY
            echo "git commit -m 'chore: resolve merge conflicts from ${{ inputs.source-branch }}'" >> $GITHUB_STEP_SUMMARY
            echo "git push origin ${{ inputs.target-branch }}" >> $GITHUB_STEP_SUMMARY
            echo "\`\`\`" >> $GITHUB_STEP_SUMMARY
            
            exit 1
          fi
      
      - name: Summary
        if: success()
        run: |
          echo "## ✅ Sync Successful" >> $GITHUB_STEP_SUMMARY
          echo "" >> $GITHUB_STEP_SUMMARY
          echo "Branch \`${{ inputs.target-branch }}\` successfully synced with \`${{ inputs.source-branch }}\`" >> $GITHUB_STEP_SUMMARY
```

#### Script de Creación

```bash
#!/bin/bash
# create-sync-branches-reusable.sh

INFRA_PATH="/Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure"
REUSABLE_FILE=".github/workflows/reusable/sync-branches.yml"

cd "$INFRA_PATH"

echo "🚀 Creando sync-branches.yml reusable..."

# Crear el workflow (copiar contenido YAML de arriba)
cat > "$REUSABLE_FILE" << 'WORKFLOW'
# ... (copiar YAML completo de arriba)
WORKFLOW

echo "✅ Workflow reusable creado"

# Documentar
cat > docs/workflows/sync-branches.md << 'DOC'
# Sync Branches Workflow (Reusable)

Workflow reusable para sincronizar ramas automáticamente.

## Características

- ✅ Sincroniza main → dev (o custom)
- ✅ Crea branch destino si no existe
- ✅ Maneja conflictos gracefully
- ✅ Proporciona instrucciones de resolución manual

## Parámetros

| Parámetro | Tipo | Default | Descripción |
|-----------|------|---------|-------------|
| `source-branch` | string | `main` | Rama origen |
| `target-branch` | string | `dev` | Rama destino |
| `create-target-if-missing` | boolean | `true` | Crear target si falta |

## Ejemplo de Uso

### Uso Básico (main → dev)

```yaml
name: Sync Main to Dev

on:
  push:
    branches: [main]
  create:
    tags: ['v*']

jobs:
  sync:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/sync-branches.yml@main
```

### Uso Custom (release → staging)

```yaml
jobs:
  sync:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/sync-branches.yml@main
    with:
      source-branch: release
      target-branch: staging
      create-target-if-missing: false
```

## Comportamiento

### Si No Hay Conflictos
1. Checkout del repositorio
2. Verificar existencia de branch destino
3. Merge automático
4. Push a origin
5. ✅ Success

### Si Hay Conflictos
1. Intento de merge falla
2. Abort del merge
3. Genera guía de resolución manual
4. ❌ Fallo (require intervención manual)

## Resolución Manual de Conflictos

Si el workflow falla, ver el Summary en Actions. Incluirá comandos exactos para resolver.

Ejemplo típico:
```bash
git checkout dev
git pull origin dev
git merge origin/main
# Resolver conflictos en editor
git add .
git commit -m "chore: resolve merge conflicts"
git push origin dev
```

DOC

echo "📝 Documentación creada"

echo ""
echo "🎉 sync-branches.yml reusable completado!"
echo ""
echo "🚀 Siguiente paso: Tarea 4.4 (Validar y documentar)"
```

#### Guardar y Ejecutar

```bash
# Similar a tarea anterior, guardar y ejecutar script
```

#### Commitear

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

git add .
git commit -m "feat: agregar workflow reusable sync-branches

Workflow reusable para sincronización automática de ramas.

Características:
- Sincroniza main → dev (configurable)
- Crea branch destino si no existe
- Maneja conflictos gracefully
- Proporciona guía de resolución manual

Parámetros configurables:
- source-branch (default: main)
- target-branch (default: dev)
- create-target-if-missing (default: true)

Documentación: docs/workflows/sync-branches.md

🤖 Generated with Claude Code

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## DÍA 2: MIGRAR API-MOBILE

---

### ✅ Tarea 4.5: Preparación y Backup

**Similar a Tarea 2.1, adaptado para Sprint 4**

---

### ✅ Tarea 4.6: Convertir pr-to-dev.yml

**Prioridad:** 🟢 P2  
**Estimación:** ⏱️ 60 minutos

#### Workflow Antes (actual)

```yaml
name: PR to Dev

on:
  pull_request:
    branches: [dev]

env:
  GO_VERSION: "1.25"

jobs:
  lint:
    # ... (código completo)
  test:
    # ... (código completo)
  build-docker:
    # ... (código completo)
```

#### Workflow Después (llamando reusable)

```yaml
name: PR to Dev

on:
  pull_request:
    branches: [dev]

jobs:
  validate:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/pr-validation.yml@main
    with:
      go-version: "1.25"
      coverage-threshold: 33
      enable-docker-build: true
      enable-security-scan: false
      docker-platforms: linux/amd64
```

**Reducción:** ~150 líneas → ~15 líneas ✅ 90% menos código

---

### ✅ Tarea 4.7: Convertir pr-to-main.yml

Similar a 4.6, pero con `enable-security-scan: true`

---

### ✅ Tarea 4.8: Convertir sync-main-to-dev.yml

```yaml
name: Sync Main to Dev

on:
  push:
    branches: [main]
  create:
    tags: ['v*']

jobs:
  sync:
    uses: EduGoGroup/edugo-infrastructure/.github/workflows/reusable/sync-branches.yml@main
```

**Reducción:** ~80 líneas → ~10 líneas ✅ 87% menos código

---

## DÍA 3-4: TESTING, DOCUMENTACIÓN, CIERRE

Las tareas 4.9 a 4.15 seguirán el mismo patrón ultra-detallado de Sprint 2, incluyendo:
- Tests exhaustivos de cada workflow
- Comparación de tiempos antes/después
- Documentación completa
- Métricas de reducción de código
- Guía de replicación para otros proyectos

---

## 📊 Métricas Esperadas

### Reducción de Código

| Workflow | Líneas Antes | Líneas Después | Reducción |
|----------|--------------|----------------|-----------|
| pr-to-dev.yml | ~150 | ~15 | 90% |
| pr-to-main.yml | ~180 | ~18 | 90% |
| sync-main-to-dev.yml | ~80 | ~10 | 87% |
| **TOTAL** | **~410** | **~43** | **~90%** |

### Mantenibilidad

- **Antes:** 3 workflows × 6 proyectos = 18 archivos a mantener
- **Después:** 2 workflows reusables + 3 callers × 6 proyectos = 20 archivos
- **Beneficio:** Cambios centralizados, consistencia garantizada

---

## ✅ Criterios de Éxito

Sprint 4 completado cuando:

- ✅ 2 workflows reusables creados y documentados
- ✅ api-mobile usa workflows reusables
- ✅ Todos los workflows pasan en CI
- ✅ Código duplicado reducido >85%
- ✅ Documentación completa para replicar
- ✅ Guía de troubleshooting creada
- ✅ PRs mergeados (infrastructure + api-mobile)
- ✅ Patrón validado y listo para replicar

---

**Generado por:** Claude Code  
**Fecha:** 19 de Noviembre, 2025  
**Versión:** 1.0  
**Estado:** Listo para Ejecución (Post Sprint 2)  
**Proyecto:** edugo-api-mobile (PILOTO)
