# 🚨 INFORME CRÍTICO: Problemas de Versionado en Proyectos EduGo

**Fecha**: 2025-11-01
**Autor**: Claude Code + Jhoan Medina
**Prioridad**: 🔴 ALTA - Requiere decisión inmediata

---

## 📊 Situación Actual

### ❌ Problema Detectado: Versionado Incorrecto en Proyectos de Desarrollo

Todos los proyectos EduGo están usando versiones **v1.x.x** y **v2.x.x** cuando deberían estar en **v0.x.x** porque:

1. ❌ NO han salido a producción
2. ❌ NO están cerca de producción
3. ❌ Están en fase de desarrollo activo
4. ❌ Pueden tener breaking changes frecuentes

---

## 🔍 Estado Actual de Cada Proyecto

### 1. **edugo-shared** (Librería)

#### Tags Actuales
```
v0.1.0          ← v0 (correcto para desarrollo)
v1.0.0          ← v1 (❌ incorrecto - fue error mío)
v2.0.0          ← v2 (❌ incorrecto - por migración a módulos)
v2.0.1
v2.0.5
v2.0.6          ← Recién creado hoy
```

#### Módulos en go.mod
```go
// ❌ INCORRECTO para versiones v2.x.x
module github.com/EduGoGroup/edugo-shared/common
module github.com/EduGoGroup/edugo-shared/auth
module github.com/EduGoGroup/edugo-shared/logger

// ✅ DEBERÍA SER (para v2.x.x)
module github.com/EduGoGroup/edugo-shared/common/v2
module github.com/EduGoGroup/edugo-shared/auth/v2

// ✅ O MEJOR AÚN (para desarrollo)
module github.com/EduGoGroup/edugo-shared/common
tag: v0.3.0 (no v2.0.6)
```

#### Problema Go Modules
```bash
$ go list -m github.com/EduGoGroup/edugo-shared/auth@v2.0.5
ERROR: invalid version: module path must match major version ("github.com/EduGoGroup/edugo-shared/auth/v2")

# Go detecta que:
# - Tag es v2.0.5 (versión mayor = 2)
# - Pero module path no tiene /v2 al final
# - Esto viola las reglas de Go modules
```

---

### 2. **edugo-api-mobile** (API)

#### Tags Actuales
```
v1.0.0          ← v1 (❌ incorrecto para desarrollo)
v1.0.1
v1.0.2
```

#### Módulo en go.mod
```go
module github.com/EduGoGroup/edugo-api-mobile  // Sin /v2, pero usa v1.x.x
```

#### Dependencias de edugo-shared
```go
github.com/EduGoGroup/edugo-shared/auth v0.0.2                              // ✅ v0 (correcto)
github.com/EduGoGroup/edugo-shared/common v0.0.0-20251031204120-ecc6528... // Pseudo-version
github.com/EduGoGroup/edugo-shared/logger v0.0.0-20251031204214-949cb60... // Pseudo-version
github.com/EduGoGroup/edugo-shared/middleware/gin v0.0.1                   // ✅ v0 (correcto)
```

**Observación**: Las dependencias usan v0.x.x o pseudo-versions, NO v2.0.x (porque v2.0.x es incompatible)

---

### 3. **edugo-api-administracion** y **edugo-worker**
(Asumo misma situación - requiere verificación)

---

## 🎯 ¿Por Qué Es Un Problema?

### 1. **Violación de Semantic Versioning**

Según [Semantic Versioning 2.0.0](https://semver.org/):

```
v0.x.x = En desarrollo, puede tener breaking changes
v1.0.0 = Primera versión estable para producción
v2.0.0 = Breaking change desde v1.x.x
```

**Problema**:
- Estamos en **v2.0.6** pero NO hemos llegado ni a QA
- Implica que hubo **DOS releases estables a producción** (v1.0.0 y v2.0.0)
- Mentira técnica sobre el estado del proyecto

### 2. **Incompatibilidad con Go Modules**

Para versiones v2+, Go requiere:

```go
// Si el tag es v2.0.5:
module github.com/EduGoGroup/edugo-shared/v2  // ← DEBE tener /v2
```

**Problema Actual**:
```go
// Tenemos:
module github.com/EduGoGroup/edugo-shared/auth
tag: v2.0.5

// Go espera:
module github.com/EduGoGroup/edugo-shared/auth/v2
tag: v2.0.5
```

**Consecuencia**: Los proyectos consumidores NO pueden usar `go get` con tags v2.x.x

### 3. **Proyectos Consumidores Usan Pseudo-Versions**

```go
// En go.mod de api-mobile:
github.com/EduGoGroup/edugo-shared/common v0.0.0-20251031204120-ecc6528ef4b6
                                           ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
                                           Esto es una PSEUDO-VERSION (commit hash)
```

**Qué significa**:
- NO está usando tags
- Está usando commits directos
- Menos predecible, menos estable
- Dificulta rollback

---

## 📋 SOLUCIONES PROPUESTAS

### ✅ **OPCIÓN 1: Resetear a v0.x.x** (RECOMENDADO)

#### Para edugo-shared:

1. **Crear nuevo tag v0.3.0** (ignora v1.x.x y v2.x.x anteriores)
   ```bash
   cd edugo-shared
   git checkout main
   git tag v0.3.0
   git push origin v0.3.0
   ```

2. **Deprecar tags v1.x.x y v2.x.x** (dejarlos pero no usarlos)
   ```bash
   # Agregar nota en GitHub Releases:
   gh release edit v2.0.6 --notes "⚠️ DEPRECADO: Usar v0.3.0 en su lugar"
   ```

3. **Actualizar CHANGELOG.md** indicando el cambio de versionado

#### Para edugo-api-mobile, api-administracion, worker:

1. **Crear nuevos tags v0.x.x**
   ```bash
   # edugo-api-mobile
   git tag v0.1.0  # Ignorar v1.0.2

   # edugo-api-administracion
   git tag v0.1.0

   # edugo-worker
   git tag v0.1.0
   ```

2. **Actualizar dependencias a v0.x.x de shared**
   ```bash
   go get github.com/EduGoGroup/edugo-shared/auth@v0.3.0
   go get github.com/EduGoGroup/edugo-shared/common@v0.3.0
   go get github.com/EduGoGroup/edugo-shared/logger@v0.3.0
   ```

#### Ventajas:
- ✅ Semánticamente correcto (v0.x.x = desarrollo)
- ✅ Compatible con Go modules (no requiere /v2)
- ✅ Permite breaking changes sin culpa
- ✅ Cuando salgas a producción → v1.0.0 tendrá significado real

#### Desventajas:
- ⚠️ Los tags v1.x.x y v2.x.x quedan como "errores históricos"
- ⚠️ Puede confundir si alguien ve tags viejos

---

### ⚠️ **OPCIÓN 2: Migrar a v2 con /v2** (Complejo)

#### Para edugo-shared:

1. **Agregar /v2 a TODOS los go.mod**:
   ```go
   // common/go.mod
   module github.com/EduGoGroup/edugo-shared/common/v2

   // auth/go.mod
   module github.com/EduGoGroup/edugo-shared/auth/v2

   // ... (7 módulos)
   ```

2. **Actualizar imports internos**:
   ```go
   // En auth/jwt.go que importa common:
   import "github.com/EduGoGroup/edugo-shared/common/v2/errors"
   ```

3. **Actualizar TODOS los proyectos consumidores**:
   ```go
   // edugo-api-mobile
   import "github.com/EduGoGroup/edugo-shared/auth/v2"
   import "github.com/EduGoGroup/edugo-shared/logger/v2"
   ```

#### Ventajas:
- ✅ Correcto según Go modules v2+ spec
- ✅ Permite coexistencia de v1 y v2

#### Desventajas:
- ❌ Cambio masivo en 7 módulos de shared
- ❌ Cambio masivo en 3+ proyectos consumidores
- ❌ Breaking change gigante
- ❌ Muchísimo trabajo
- ❌ Sigue siendo v2 cuando debería ser v0

---

### 🔄 **OPCIÓN 3: Mantener Como Está** (No Recomendado)

Seguir usando v2.x.x con `+incompatible`:

```bash
# Los consumidores pueden usar:
go get github.com/EduGoGroup/edugo-shared/auth@v2.0.6

# Go.mod mostrará:
github.com/EduGoGroup/edugo-shared/auth v2.0.6+incompatible
```

#### Ventajas:
- ✅ No requiere cambios

#### Desventajas:
- ❌ "+incompatible" en todos los go.mod
- ❌ Violación de semantic versioning
- ❌ Violación de Go modules spec
- ❌ Confusión sobre estado del proyecto
- ❌ Dificultad para explicar a nuevos desarrolladores

---

## 🎯 MI RECOMENDACIÓN FUERTE: OPCIÓN 1

### Plan de Acción Inmediato

#### **Paso 1: Resetear edugo-shared a v0.3.0**

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared

# 1. Crear tag v0.3.0 (nuevo esquema de versionado)
git checkout main
git pull
git tag v0.3.0
git push origin v0.3.0

# 2. Actualizar CHANGELOG.md indicando el cambio
cat >> CHANGELOG.md << 'EOF'

## [0.3.0] - 2025-11-01

### ⚠️ BREAKING: Cambio de Esquema de Versionado

Este proyecto vuelve a versionado **v0.x.x** para reflejar correctamente su estado de desarrollo.

**Razones**:
- El proyecto NO ha salido a producción
- Permite breaking changes sin violar semantic versioning
- Compatible con Go modules (no requiere /v2 en module path)

**Migración**:
- Tags v1.x.x y v2.x.x quedan deprecados
- Usar v0.3.0+ en adelante
- Cuando salga a producción → v1.0.0

### Added
- Copilot custom instructions
- Workflows CI/CD optimizados (matrix strategy)
- Workflow sync-main-to-dev

### Changed
- Módulo middleware/gin agregado a workflows
- Go version 1.23 → 1.25
- codecov-action v3 → v4
EOF

git add CHANGELOG.md
git commit -m "docs: cambio de esquema de versionado v2.x.x → v0.x.x

El proyecto está en desarrollo, no ha salido a producción.
Versiones v1.x.x y v2.x.x fueron error de versionado.

Migración:
- Tags v2.x.x quedan deprecados
- Nuevo esquema: v0.3.0+
- Al salir a producción: v1.0.0

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"

git push origin main

# 3. Deprecar releases anteriores
gh release edit v2.0.6 --notes "⚠️ DEPRECADO: Este proyecto cambió a versionado v0.x.x. Usar v0.3.0 en su lugar."
gh release edit v2.0.5 --notes "⚠️ DEPRECADO: Usar v0.3.0 en su lugar."
```

#### **Paso 2: Actualizar edugo-api-mobile**

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# 1. Crear branch
git checkout dev
git pull
git checkout -b fix/corregir-versionado

# 2. Actualizar dependencias de shared a v0.3.0
go get github.com/EduGoGroup/edugo-shared/auth@v0.3.0
go get github.com/EduGoGroup/edugo-shared/common@v0.3.0
go get github.com/EduGoGroup/edugo-shared/logger@v0.3.0
go get github.com/EduGoGroup/edugo-shared/middleware/gin@v0.3.0
go mod tidy

# 3. Eliminar tags v1.x.x y crear v0.x.x
git tag -d v1.0.0 v1.0.1 v1.0.2
git push origin :refs/tags/v1.0.0 :refs/tags/v1.0.1 :refs/tags/v1.0.2

# 4. Actualizar archivo de versión
echo "0.1.0" > .github/version.txt

# 5. Commit y PR
git add go.mod go.sum .github/version.txt
git commit -m "fix: corregir versionado a v0.x.x (proyecto en desarrollo)"
git push origin fix/corregir-versionado
gh pr create --base dev
```

#### **Paso 3: Repetir para api-administracion y worker**

```bash
# Similar al paso 2, para cada proyecto
```

---

## 📚 Explicación Técnica Detallada

### ¿Qué Son Tags vs Releases?

```
┌─────────────────────────────────────────────────────────────┐
│  GIT TAG (Técnico)                                          │
│  ────────────────                                           │
│  - Puntero a un commit específico en Git                    │
│  - Vive en el repositorio Git                               │
│  - Usado por: go get, git checkout, herramientas CLI        │
│  - Ejemplo: v0.3.0, auth/v0.1.0                             │
│                                                              │
│  Comando: git tag v0.3.0                                    │
│           git push origin v0.3.0                            │
└─────────────────────────────────────────────────────────────┘
                         │
                         │ (se basa en)
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  GITHUB RELEASE (Visual/UI)                                 │
│  ────────────────────────                                   │
│  - Feature de interfaz web de GitHub                        │
│  - Incluye: changelog, notas, archivos adjuntos             │
│  - Usado por: Humanos navegando en github.com               │
│  - NO usado por: go get (solo lee tags)                     │
│                                                              │
│  Se crea: Automáticamente (release.yml) o manual en UI      │
└─────────────────────────────────────────────────────────────┘
```

### ¿Cuál Usa `go get`?

```bash
# go get SOLO lee TAGS (no releases)
go get github.com/EduGoGroup/edugo-shared/auth@v0.3.0
                                                 ^^^^^^
                                                 TAG (no release)

# Proceso:
# 1. go get busca el tag v0.3.0 en Git
# 2. Descarga el código de ese commit
# 3. Ignora completamente los GitHub Releases
```

### ¿Para Qué Sirven los Releases Entonces?

```
GitHub Releases son para:
✅ Documentación visual del cambio
✅ Notas para humanos (qué cambió, cómo migrar)
✅ Adjuntar binarios compilados (no aplica para librerías Go)
✅ Marketing/comunicación del release

NO para:
❌ go get (usa tags)
❌ Versionado técnico (usa tags)
❌ Instalación de dependencias (usa tags)
```

---

## 🔢 Explicación de Versionado Semántico

### Para Proyectos en Desarrollo (v0.x.x)

```
v0.1.0 → Primera versión usable
v0.2.0 → Nueva feature
v0.2.1 → Bugfix
v0.3.0 → Otra feature
v0.10.0 → Décima feature
v0.99.0 → Feature 99 (aún en desarrollo)

Breaking changes: PERMITIDOS en cualquier momento
Semántica: "Esto está en desarrollo, puede cambiar"
```

### Para Proyectos Estables (v1.x.x)

```
v1.0.0 → PRIMERA VERSIÓN EN PRODUCCIÓN (hito importante)
v1.1.0 → Nueva feature (retrocompatible)
v1.1.1 → Bugfix
v1.2.0 → Otra feature (retrocompatible)

Breaking changes: PROHIBIDOS (requiere v2.0.0)
Semántica: "Esto está en producción, estable, confiable"
```

### Para Breaking Changes de Producción (v2.x.x)

```
v2.0.0 → BREAKING CHANGE desde v1.x.x
         (Requiere module path /v2 en Go)
v2.1.0 → Feature nueva (retrocompatible con v2.0.0)

Semántica: "Nueva versión mayor, incompatible con v1.x.x"
```

---

## 📊 Comparación: Estado Actual vs Estado Correcto

### Estado ACTUAL (Incorrecto)

```
edugo-shared:
  Tags: v2.0.6 ❌
  Module path: github.com/EduGoGroup/edugo-shared/auth ❌
  Implicación: "2 releases estables en producción" (falso)

edugo-api-mobile:
  Tags: v1.0.2 ❌
  Module path: github.com/EduGoGroup/edugo-api-mobile ✅
  Implicación: "Release estable en producción" (falso)

Dependencias:
  api-mobile usa: v0.0.0-20251031... (pseudo-versions) ⚠️
  Problema: No usa tags, usa commits directos
```

### Estado CORRECTO (Recomendado)

```
edugo-shared:
  Tags: v0.3.0 ✅
  Module path: github.com/EduGoGroup/edugo-shared/auth ✅
  Implicación: "En desarrollo, pre-producción" (verdadero)

edugo-api-mobile:
  Tags: v0.1.0 ✅
  Module path: github.com/EduGoGroup/edugo-api-mobile ✅
  Implicación: "En desarrollo" (verdadero)

Dependencias:
  api-mobile usa: auth@v0.3.0 (tags limpios) ✅
  Beneficio: Versionado claro, fácil rollback
```

---

## 🚨 IMPACTO DE NO CORREGIR

### Corto Plazo
- ⚠️ Confusión en el equipo sobre estado del proyecto
- ⚠️ `go get` no funciona con tags v2.x.x (requiere workarounds)
- ⚠️ Pseudo-versions en lugar de tags limpios

### Mediano Plazo
- ❌ Al llegar a producción real, ¿usar v3.0.0?
- ❌ Explicar a stakeholders por qué v2.0.6 no está en producción
- ❌ Deuda técnica acumulada

### Largo Plazo
- ❌ Historial de versiones engañoso
- ❌ Dificultad para auditorías
- ❌ Violación de mejores prácticas de Go

---

## 📝 PLAN DE ACCIÓN RECOMENDADO

### Fase 1: Documentar Decisión (HOY)
1. ✅ Leer este informe
2. ✅ Decidir: ¿Opción 1 (v0.x.x) o mantener v2.x.x?
3. ✅ Documentar decisión en este archivo

### Fase 2: Corregir edugo-shared (1-2 días)
1. Crear tag v0.3.0
2. Actualizar CHANGELOG.md
3. Deprecar releases v1.x.x y v2.x.x
4. Comunicar cambio a equipo

### Fase 3: Actualizar Proyectos Consumidores (2-3 días)
1. edugo-api-mobile: v1.0.2 → v0.1.0
2. edugo-api-administracion: revisar y ajustar
3. edugo-worker: revisar y ajustar
4. Actualizar dependencias de shared a v0.3.0

### Fase 4: Estandarizar (1 día)
1. Actualizar documentación de todos los proyectos
2. Establecer política de versionado
3. Configurar branch protection para prevenir tags incorrectos

---

## 🎯 DECISIÓN REQUERIDA

**Jhoan, necesito que decidas:**

### Pregunta 1: ¿Resetear a v0.x.x?
- [ ] **Sí** → Proceder con Opción 1 (resetear a v0.3.0)
- [ ] **No** → Mantener v2.x.x y agregar /v2 a module paths (Opción 2)
- [ ] **Más tarde** → Documentar deuda técnica y resolver después

### Pregunta 2: ¿Cuándo ejecutar la corrección?
- [ ] **Ahora** → Empezar inmediatamente
- [ ] **Después de completar FASE 3 del sprint** → Posponer
- [ ] **Antes de salir a QA** → Incluir en checklist pre-QA

### Pregunta 3: ¿Qué hacer con tags históricos?
- [ ] **Eliminar tags v1.x.x y v2.x.x** del repositorio
- [ ] **Deprecar** (dejar con nota de deprecación)
- [ ] **Dejar como están** (solo no usar más)

---

## 📋 Sobre el Error de test.yml

### Problema Reportado
```
Runs 19002030386, 19001997913, 19001961065
→ Aparecen en ROJO en lista de actions
→ Duración: 0s
→ Workflow: .github/workflows/test.yml
```

### Explicación Técnica

**NO es un error real, es comportamiento esperado**:

1. **Qué pasa**:
   - GitHub Actions detecta cambios en `.github/workflows/test.yml`
   - Intenta ejecutar el workflow
   - El workflow tiene `on: pull_request` (NO tiene `on: push`)
   - GitHub falla inmediatamente (0s) porque el trigger no coincide

2. **Por qué aparece en rojo**:
   - GitHub registra el intento como "failed"
   - Pero no es un fallo real de código
   - Es un "no-op" (no operation)

3. **Por qué NO afecta los PRs**:
   - Los PRs solo muestran checks triggered por `pull_request`
   - Estos "errores" son triggered por `push`
   - Dos contextos diferentes en GitHub Actions

### Solución Aplicada

✅ **Agregado comentario explicativo** en test.yml (PR #3 ya mergeado a dev)

```yaml
# IMPORTANTE: Este workflow NO se ejecuta en push (solo PRs y manual)
# Los "errores" en push son esperados - GitHub intenta ejecutar el workflow
# pero falla inmediatamente (0s) porque no tiene trigger para push.
# Esto es comportamiento normal y no afecta el flujo de trabajo.
```

### ¿Se Puede Eliminar Estos Errores?

**Opción A**: Agregar trigger push (no recomendado)
```yaml
on:
  workflow_dispatch:
  pull_request:
    branches: [ main, dev ]
  push:  # ← Esto eliminaría los errores
    branches: [ main, dev ]
```

**Problema**: Ejecutaría tests en CADA push a dev/main, desperdiciando minutos de GitHub Actions.

**Opción B**: Dejar como está (recomendado)
- Los "errores" no afectan nada
- Son solo visuales en la lista de actions
- Los PRs siguen mostrando todo verde

---

## 📊 Resumen Ejecutivo

### ✅ Lo Que Funciona Bien
1. CI/CD completo implementado
2. Workflows ejecutándose correctamente en PRs
3. Release automático funcionando
4. Sync main → dev funcionando

### ❌ Problemas Críticos a Resolver
1. **Versionado incorrecto**: v2.0.6 debería ser v0.3.0
2. **Module path incompatible**: v2.x.x requiere /v2 (pero no lo tiene)
3. **Semantic versioning violado**: v1.0.0 y v2.0.0 implican producción
4. **Pseudo-versions en consumidores**: En lugar de tags limpios

### ⚠️ Problemas Menores (Cosméticos)
1. test.yml muestra errores en push (0s) - Esperado, no crítico

---

## 📞 ACCIÓN REQUERIDA DE JHOAN

Por favor, responde:

1. **¿Procedo a resetear todo a v0.x.x?** (Mi recomendación fuerte: SÍ)
2. **¿Elimino los tags v1.x.x y v2.x.x o solo los depreco?**
3. **¿Quieres que haga esto ahora o después del sprint actual?**

**Nota**: Esto afectará los 4 proyectos (shared, api-mobile, api-administracion, worker) pero es la decisión técnicamente correcta.

---

**Creado**: 2025-11-01 20:30
**Requiere decisión**: URGENTE
**Impacto**: Alto (4 proyectos)
**Esfuerzo estimado**: 1-2 días
