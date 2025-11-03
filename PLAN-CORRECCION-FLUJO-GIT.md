# 🔧 PLAN DE CORRECCIÓN Y PREVENCIÓN DEL FLUJO GIT

**Fecha**: 2025-11-02
**Proyecto**: edugo-api-mobile (aplicable a todos los proyectos hermanos)
**Objetivo**: Prevenir commits directos y asegurar sincronización automática

---

## 🎯 PROBLEMAS IDENTIFICADOS

### 1. Workflow sync-main-to-dev NO hace merge directo

**Problema actual** (líneas 88-110):
```yaml
# El workflow CREA UN PR pero NO mergea automáticamente
gh pr create --base dev --head main ...
gh pr merge "$PR_NUMBER" --auto --squash  # ← ESTO FALLA
```

**Por qué falla el auto-merge**:
- ❌ Requiere permisos especiales (`pull-requests: write` no es suficiente)
- ❌ Si hay branch protection con required approvals, no puede auto-mergear
- ❌ El `--auto` merge requiere que el repo tenga "Allow auto-merge" habilitado

**Resultado**: PR queda abierto esperando merge manual → usuario hace merge con estrategia incorrecta

---

### 2. Commits directos a main y dev están permitidos

**Problema**: Nada impide hacer:
```bash
git checkout main
# hacer cambios
git commit -m "fix: algo"
git push origin main   # ← ESTO DEBERÍA ESTAR BLOQUEADO
```

**Resultado**: Commits que no pasan por PR → no hay review → no hay CI/CD previo

---

### 3. No hay estrategia de merge forzada

**Problema**: Al mergear PRs, GitHub permite elegir entre:
- Merge commit
- Squash and merge
- Rebase and merge

**Resultado**: Inconsistencia en el historial de commits

---

## ✅ SOLUCIONES PROPUESTAS

### Solución 1: Mejorar el Workflow sync-main-to-dev

#### Opción A: Merge Directo (Recomendado)

Cambiar el workflow para que haga merge directo sin crear PR:

```yaml
- name: Merge main to dev directamente
  if: steps.check_diff.outputs.has_diff == 'true'
  run: |
    VERSION="${{ steps.version.outputs.version }}"

    # Configurar git
    git config user.name "github-actions[bot]"
    git config user.email "github-actions[bot]@users.noreply.github.com"

    # Checkout dev y merge main
    git checkout dev
    git merge origin/main --no-ff -m "chore: sync main v$VERSION to dev

    Sincronización automática de main a dev después de release.

    🤖 Generated with [Claude Code](https://claude.com/claude-code)

    Co-Authored-By: Claude <noreply@anthropic.com>"

    # Push a dev
    git push origin dev

    echo "✅ main sincronizado a dev correctamente"
```

**Ventajas**:
- ✅ Sincronización inmediata
- ✅ No requiere aprobación manual
- ✅ Funciona aunque haya branch protection

**Desventajas**:
- ⚠️ Si hay conflictos, el workflow falla (pero esto es bueno, indica un problema real)

---

#### Opción B: PR con Auto-Merge Mejorado

Si prefieres mantener el PR para tener trazabilidad:

```yaml
- name: Habilitar auto-merge en el repo
  env:
    GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  run: |
    # Crear PR
    PR_NUMBER=$(gh pr create \
      --base dev \
      --head main \
      --title "chore: sync main v$VERSION to dev" \
      --body "$BODY" \
      --label "sync" \
      --label "automated" \
      --json number \
      --jq '.number')

    # Habilitar auto-merge (requiere que el repo tenga la feature habilitada)
    gh pr merge "$PR_NUMBER" --merge --auto

    # Si falla, mergear manualmente (solo si no hay conflictos)
    if ! gh pr checks "$PR_NUMBER" --watch; then
      echo "⚠️ Checks fallaron, revisar manualmente"
      exit 1
    fi

    # Aprobar el PR (requiere PAT o GitHub App)
    gh pr review "$PR_NUMBER" --approve --body "Auto-aprobado: sincronización automática"

    # Mergear
    gh pr merge "$PR_NUMBER" --merge
```

**Desventaja**: Requiere configuración adicional (PAT con permisos o GitHub App)

---

**RECOMENDACIÓN**: Usar **Opción A (merge directo)** por simplicidad y confiabilidad.

---

### Solución 2: Configurar Branch Protection Rules

Ir a **GitHub → Settings → Branches → Add rule** para **main** y **dev**:

#### Para `main`:

```yaml
Branch name pattern: main

☑️ Require a pull request before merging
  ☑️ Require approvals: 1
  ☑️ Dismiss stale pull request approvals when new commits are pushed
  ☑️ Require review from Code Owners (opcional)

☑️ Require status checks to pass before merging
  ☑️ Require branches to be up to date before merging
  Status checks:
    - CI
    - Tests

☑️ Require conversation resolution before merging

☑️ Require linear history (fuerza squash o rebase)

☑️ Do not allow bypassing the above settings
  Exceptions:
    - github-actions[bot] (para permitir workflow sync)

☐ Allow force pushes
  - NO marcar (fuerza push bloqueado)

☐ Allow deletions
  - NO marcar (eliminación bloqueada)
```

#### Para `dev`:

```yaml
Branch name pattern: dev

☑️ Require a pull request before merging
  ☑️ Require approvals: 1 (opcional, puede ser 0 para desarrollo rápido)

☑️ Require status checks to pass before merging
  Status checks:
    - CI
    - Tests

☑️ Do not allow bypassing the above settings
  Exceptions:
    - github-actions[bot] (para permitir workflow sync)

☐ Allow force pushes
  - NO marcar

☐ Allow deletions
  - NO marcar
```

**Resultado**:
- ✅ Commits directos a main/dev: **BLOQUEADOS**
- ✅ Todo debe pasar por PR
- ✅ CI/CD se ejecuta antes de mergear
- ✅ github-actions[bot] puede sincronizar

---

### Solución 3: Forzar Estrategia de Merge

En **GitHub → Settings → General → Pull Requests**:

```
☐ Allow merge commits
☑️ Allow squash merging (RECOMENDADO)
☐ Allow rebase merging

Default to squash merging
```

**Resultado**: Todos los PRs se mergean con squash → historial limpio

---

### Solución 4: Configurar CODEOWNERS (Opcional)

Crear `.github/CODEOWNERS`:

```
# Requiere aprobación de owners para archivos críticos

# Workflows CI/CD
/.github/workflows/  @jhoanmedina

# Configuración del proyecto
/config/  @jhoanmedina
/.env.example  @jhoanmedina

# Scripts de base de datos
/scripts/postgresql/  @jhoanmedina

# Módulos críticos
/internal/domain/  @jhoanmedina
```

---

## 📋 FLUJO DE TRABAJO OBLIGATORIO

### Desarrollo Normal

```
1. Crear feature branch desde dev:
   git checkout dev
   git pull origin dev
   git checkout -b feature/mi-feature

2. Desarrollar y commitear:
   git add .
   git commit -m "feat: mi nueva feature"

3. Push de feature branch:
   git push origin feature/mi-feature

4. Crear PR: feature/mi-feature → dev
   - CI/CD se ejecuta
   - Requiere 1 aprobación (si está configurado)
   - Squash merge

5. Mergear a dev:
   - PR se mergea con squash
   - dev se actualiza

6. RELEASE: Cuando dev esté listo para producción
   - Crear PR: dev → main
   - CI/CD se ejecuta
   - Requiere 1 aprobación
   - Squash merge

7. Mergear a main:
   - PR se mergea
   - Workflow sync-main-to-dev SE EJECUTA AUTOMÁTICAMENTE
   - main se sincroniza a dev sin intervención manual
```

### Hotfix (Emergencia en producción)

```
1. Crear hotfix branch desde main:
   git checkout main
   git pull origin main
   git checkout -b hotfix/fix-critico

2. Desarrollar fix:
   git add .
   git commit -m "fix: corregir bug crítico en producción"

3. Push de hotfix branch:
   git push origin hotfix/fix-critico

4. Crear PR: hotfix/fix-critico → main
   - CI/CD se ejecuta
   - Requiere 1 aprobación
   - Squash merge

5. Mergear a main:
   - PR se mergea
   - Workflow sync-main-to-dev sincroniza a dev automáticamente
```

---

## ⚠️ LO QUE NUNCA DEBES HACER

❌ **Commits directos a main o dev**:
```bash
git checkout main
git commit -m "algo"  # ← BLOQUEADO por branch protection
git push origin main  # ← FALLA
```

❌ **Merge manual sin PR**:
```bash
git checkout dev
git merge main  # ← EVITAR (no hay trazabilidad ni CI/CD)
git push origin dev
```

❌ **Estrategias de merge inconsistentes**:
- Usar "merge commit" en un PR y "squash" en otro
- Esto genera el problema de "25 commits fantasma"

❌ **Force push a main o dev**:
```bash
git push --force origin main  # ← BLOQUEADO
```

---

## 🔧 PASOS DE IMPLEMENTACIÓN

### Paso 1: Actualizar Workflow sync-main-to-dev.yml

- [ ] Reemplazar sección "Create PR" con merge directo (Opción A)
- [ ] Commit y push del cambio
- [ ] Verificar que el workflow se actualiza en GitHub Actions

### Paso 2: Configurar Branch Protection en GitHub

- [ ] Ir a Settings → Branches
- [ ] Configurar rule para `main` (según especificaciones arriba)
- [ ] Configurar rule para `dev` (según especificaciones arriba)
- [ ] Agregar exception para `github-actions[bot]`

### Paso 3: Configurar Merge Strategy

- [ ] Ir a Settings → General → Pull Requests
- [ ] Desmarcar "Allow merge commits"
- [ ] Marcar "Allow squash merging"
- [ ] Set default to "Squash merging"

### Paso 4: Crear CODEOWNERS (Opcional)

- [ ] Crear `.github/CODEOWNERS`
- [ ] Agregar owners para archivos críticos
- [ ] Commit y push

### Paso 5: Documentar en README

- [ ] Agregar sección "Flujo de Trabajo Git"
- [ ] Documentar proceso de PR
- [ ] Documentar proceso de release

### Paso 6: Probar el Flujo Completo

- [ ] Crear feature branch de prueba
- [ ] Crear PR a dev → verificar CI/CD
- [ ] Mergear a dev → verificar squash
- [ ] Crear PR de dev a main → verificar CI/CD
- [ ] Mergear a main → **VERIFICAR QUE SYNC AUTOMÁTICO FUNCIONA**
- [ ] Verificar que dev se actualizó automáticamente

---

## 📊 CHECKLIST DE VERIFICACIÓN POST-IMPLEMENTACIÓN

Después de implementar, verificar:

- [ ] ✅ No puedo hacer commit directo a main
- [ ] ✅ No puedo hacer commit directo a dev
- [ ] ✅ Los PRs requieren CI/CD passing
- [ ] ✅ Los PRs requieren aprobación (si configurado)
- [ ] ✅ Los merges usan squash automáticamente
- [ ] ✅ Workflow sync-main-to-dev mergea automáticamente
- [ ] ✅ Después de mergear a main, dev se actualiza solo
- [ ] ✅ `git rev-list --count origin/main...origin/dev` da `0  0`

---

## 🎯 BENEFICIOS ESPERADOS

Después de implementar estas correcciones:

1. ✅ **Historial limpio**: Squash merge elimina commits intermedios
2. ✅ **Sincronización automática**: main → dev sin intervención manual
3. ✅ **Trazabilidad**: Todo cambio pasa por PR
4. ✅ **CI/CD garantizado**: Nada llega a main/dev sin pasar tests
5. ✅ **Protección**: Commits accidentales bloqueados
6. ✅ **Consistencia**: Misma estrategia de merge en todos los PRs

---

## 📚 REFERENCIAS

- [GitHub Branch Protection](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches)
- [GitHub Auto-merge](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/incorporating-changes-from-a-pull-request/automatically-merging-a-pull-request)
- [GitHub CODEOWNERS](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners)

---

**Próximos pasos**: Implementar estas correcciones en edugo-api-mobile y luego replicar en proyectos hermanos.
