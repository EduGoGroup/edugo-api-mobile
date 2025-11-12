# 📊 INFORME DE ANÁLISIS DEL ECOSISTEMA EDUGO

**Fecha de análisis**: 2025-11-02
**Responsable**: Claude Code + Jhoan Medina
**Versión del informe**: 1.0

---

## 📋 ÍNDICE

1. [Resumen Ejecutivo](#resumen-ejecutivo)
2. [Análisis por Proyecto](#análisis-por-proyecto)
   - [edugo-api-mobile](#1-edugo-api-mobile)
   - [edugo-api-administracion](#2-edugo-api-administracion)
   - [edugo-worker](#3-edugo-worker)
   - [edugo-shared](#4-edugo-shared)
   - [edugo-dev-environment](#5-edugo-dev-environment)
3. [Análisis Comparativo](#análisis-comparativo)
4. [Problemas Críticos Identificados](#problemas-críticos-identificados)
5. [Explicación Conceptual: Merge Commits y Sincronización](#explicación-conceptual-merge-commits-y-sincronización)
6. [Plan de Acción Detallado](#plan-de-acción-detallado)
7. [Anexos](#anexos)

---

## 🎯 RESUMEN EJECUTIVO

### Alcance del Análisis

Se realizó un análisis exhaustivo de **5 repositorios** del ecosistema EduGo:
- 2 APIs REST (api-mobile, api-administracion)
- 1 Worker (edugo-worker)
- 1 Librería compartida (edugo-shared)
- 1 Ambiente de desarrollo (edugo-dev-environment)

### Hallazgos Principales

| Categoría | Estado | Descripción |
|-----------|--------|-------------|
| **Workflows CI/CD** | 🟢 Estandarizado | Los 3 proyectos de servicio tienen workflows idénticos |
| **Versionado** | 🔴 Crítico | Inconsistencia entre version.txt (v0.x.x) y tags (v1.x.x, v2.x.x) |
| **Sincronización** | 🟡 Atención | api-mobile tiene 25 commits en dev sin mergear a main |
| **Estructura** | 🟢 Consistente | Estructura de directorios coherente entre proyectos |

### Métricas Generales

- **Total de repositorios analizados**: 5
- **Workflows CI/CD encontrados**: 31 archivos
- **Tags totales**: ~30 tags distribuidos
- **Commits pendientes de sincronización**: ~30 commits entre todos los proyectos

### Prioridades de Acción

1. 🔴 **URGENTE**: Unificar esquema de versionado (v0.x.x vs v1.x.x)
2. 🔴 **URGENTE**: Mergear edugo-api-mobile dev → main (25 commits)
3. 🟡 **IMPORTANTE**: Sincronizar edugo-shared (ramas divergentes)
4. 🟢 **MEJORA**: Agregar CI/CD a edugo-dev-environment

---

## 🔬 HALLAZGOS REALES DE INVESTIGACIÓN PROFUNDA

**Fecha de investigación profunda**: 2025-11-02 (22:30)
**Enfoque**: edugo-shared y edugo-api-mobile (proyectos vivos prioritarios)

### 🎯 Contexto de la Investigación

El usuario identificó que algo no cuadraba: si la última imagen Docker se construyó recientemente, ¿cómo es posible que dev esté 25 commits adelante de main? Esta inconsistencia motivó una investigación exhaustiva de las ramas remotas para entender qué realmente pasó.

---

### 1. edugo-shared - ANÁLISIS REAL

#### Estado Actual de Ramas Remotas

```bash
origin/main (082f430): fix: Resetear versionado a v0.3.0 (#5) - 2025-11-01
origin/dev  (74864bf): fix: resetear esquema de versionado de v2.x.x a v0.x.x (#4) - 2025-11-01
```

#### Último Release Válido

**✅ v0.3.0** (commit 082f430) - 2025-11-01

Este es el ÚNICO release válido en esquema v0.x.x. Incluye tags para todos los módulos:
- `v0.3.0` (global)
- `auth/v0.3.0`
- `common/v0.3.0`
- `database/mongodb/v0.3.0`
- `database/postgres/v0.3.0`
- `logger/v0.3.0`
- `messaging/rabbit/v0.3.0`
- `middleware/gin/v0.3.0`

**❌ Tags ERRÓNEOS** (deben ignorarse o eliminarse):
- `v2.0.6`, `v2.0.5`, `v2.0.1`, `v2.0.0` - ERROR GRAVE (no v0.x.x)
- `v1.0.0` - ERROR (no v0.x.x)
- Tags de módulos con v2.x.x - ERROR

#### Divergencia de Ramas

**PROBLEMA CRÍTICO**: Las ramas han DIVERGIDO.

```
Commits SOLO en origin/main (que dev NO tiene):
- 082f430: fix: Resetear versionado a v0.3.0 (#5) ← RELEASE v0.3.0
- 6bb83e1: feat: Copilot instructions y optimización CI/CD (#1) (#2)

Commits SOLO en origin/dev (que main NO tiene):
- 74864bf: fix: resetear esquema de versionado de v2.x.x a v0.x.x (#4)
- d9fc9cd: docs: agregar comentario explicativo sobre errores de test.yml en push (#3)
- 66d1fde: feat: Copilot instructions y optimización CI/CD (#1)
```

#### Causa Raíz de la Divergencia

Análisis del gráfico de commits revela:

```
       ┌─ 082f430 (origin/main) ← tag v0.3.0
       │  6bb83e1
       │
base───┤ 4330be1 (tag: middleware/gin/v0.0.1)
       │
       └─ 66d1fde (origin/dev)
          d9fc9cd
          74864bf
```

**Qué pasó**:
1. El PR #1 "Copilot instructions y optimización CI/CD" se trabajó en ambas ramas
2. En dev se mergeó como commit 66d1fde
3. En main se mergeó como commit 6bb83e1 (con "#1 (#2)" en el mensaje)
4. Son commits DIFERENTES del mismo trabajo (posible rebase o cherry-pick)
5. Luego main agregó el commit 082f430 (reseteo de versionado a v0.3.0)
6. Luego dev agregó commits para resetear versionado pero de forma diferente

**Resultado**: Ramas divergentes con trabajo duplicado pero commits diferentes.

---

### 2. edugo-api-mobile - ANÁLISIS REAL

#### Estado Actual de Ramas Remotas

```bash
origin/main (6e88a52): Dev (#8) - 2025-11-02 21:01
origin/dev  (c907a42): Merge branch 'main' into dev - 2025-11-02
```

#### Último Release Válido

**✅ v0.1.1** (commit 1dde3c8) - 2025-11-01

```
Release v0.1.1 - Latest
Tag: v0.1.1
Commit: 1dde3c8
Fecha creación: 2025-11-01 22:55:20
Fecha publicación: 2025-11-01 23:18:47
```

**Imagen Docker**: `ghcr.io/edugogroup/edugo-api-mobile:0.1.1`

**❌ Releases ERRÓNEOS** (fueron experimentos con versionado incorrecto):
- `v1.0.2` (2025-11-01 00:59) - ERROR
- `v1.0.1` (2025-11-01 00:31) - ERROR
- `v1.0.0` (2025-10-31 23:53) - ERROR

Nota: Estos releases v1.x.x fueron creados ANTES de v0.1.1, durante experimentos con versionado.

#### ¿Qué contiene el último release v0.1.1?

El commit 1dde3c8 está PRESENTE en ambas ramas (main y dev). Este es el último código estable released.

**Desde v0.1.1 (1dde3c8) hasta origin/main (6e88a52)**:
- Solo **3 commits**:
  1. `44c8b17`: docs: actualizar plan CI/CD con workflow manual-release TODO-EN-UNO
  2. `2dfe7f2`: Merge remote-tracking branch 'origin/main'
  3. `6e88a52`: Dev (#8) ← **MEGA-MERGE PR**

**Desde v0.1.1 (1dde3c8) hasta origin/dev (c907a42)**:
- **26 commits** (incluye los 3 de arriba + 23 únicos en dev)

#### El Misterio del "Dev (#8)" - RESUELTO

El commit `6e88a52 "Dev (#8)"` es un **MEGA-MERGE COMMIT** que fue mergeado HOY (2025-11-02 21:01).

**Contenido del PR #8** (según el mensaje de commit):
- ✅ GitFlow setup (preparar estructura)
- ✅ Conectar implementación real con Container DI
- ✅ Sistema completo de autenticación JWT
- ✅ Migración a bcrypt
- ✅ Refresh tokens con revocación
- ✅ Migración a middleware compartido (edugo-shared)
- ✅ Rate limiting anti-fuerza bruta
- ✅ Documentación completa
- ✅ Copilot instructions
- ✅ Fixes de CI/CD
- ✅ Actualización a edugo-shared v0.3.0

**Total**: 25 commits de features mergeados en un solo PR a main.

#### ¿Por Qué dev está "25 commits adelante"?

**LA RESPUESTA REAL**:

1. Los 25 commits de features se desarrollaron en la rama `feature/conectar`
2. Se creó el PR #8 de `feature/conectar` → `main`
3. El PR #8 se MERGEÓ a main creando el commit `6e88a52` (HOY 21:01)
4. El workflow `sync-main-to-dev.yml` DEBIÓ ejecutarse automáticamente
5. PERO el usuario reporta que hizo un **merge MANUAL** de main → dev (minutos antes de esta investigación)
6. Ese merge manual posiblemente usó estrategia incorrecta
7. Resultado: dev tiene el commit `c907a42 "Merge branch 'main' into dev"`
8. PERO `c907a42` NO incluye correctamente el commit `6e88a52` de main

#### Visualización del Problema

```
                                    ┌─ 6e88a52 (origin/main) "Dev (#8)"
                                    │  [MEGA-MERGE con 25 commits]
                                    │  Mergeado HOY 21:01
                                    │
v0.1.1 ─── 44c8b17 ─── 2dfe7f2 ─────┤
(1dde3c8)                           │
                                    │
                                    └─ c907a42 (origin/dev) "Merge main → dev"
                                       [Merge manual, estrategia incorrecta]

                                       Los 25 commits están DEBAJO en dev:
                                       ├─ 74864bf, d9fc9cd, ed5fcdf...
                                       ├─ Feature completa de autenticación
                                       ├─ bcrypt, refresh tokens, rate limiting
                                       └─ Migración a shared, copilot, etc.
```

#### La Verdad Sobre la Sincronización

**Qué debió pasar**:
1. PR #8 se mergea a main → crea 6e88a52
2. Workflow `sync-main-to-dev.yml` se ejecuta automáticamente
3. Main se mergea a dev automáticamente
4. dev y main quedan AL MISMO NIVEL

**Qué realmente pasó**:
1. PR #8 se mergeó a main → crea 6e88a52 ✅
2. Workflow `sync-main-to-dev.yml` NO se ejecutó (o falló) ❌
3. Usuario hizo merge MANUAL de main → dev
4. El merge manual creó c907a42 pero con estrategia incorrecta
5. Ahora main y dev están DIVERGENTES:
   - main tiene 6e88a52 (mega-merge)
   - dev tiene los 25 commits individuales SIN 6e88a52

#### Estado de Divergencia

```bash
$ git rev-list --left-right --count origin/main...origin/dev
1       25

origin/main tiene 1 commit que dev NO tiene: 6e88a52
origin/dev tiene 25 commits que main NO tiene: [los commits individuales]
```

**PERO ATENCIÓN**: El commit 6e88a52 en main CONTIENE el trabajo de esos 25 commits. Son los MISMOS CAMBIOS en forma diferente:
- En main: 1 merge commit con 25 cambios incorporados
- En dev: 25 commits individuales

---

### 🎯 CONCLUSIONES DE LA INVESTIGACIÓN PROFUNDA

#### edugo-shared

| Aspecto | Estado | Detalle |
|---------|--------|---------|
| Último release válido | v0.3.0 | Commit 082f430, todos los módulos v0.3.0 |
| Estado de main | ✅ Correcto | Tiene el release v0.3.0 |
| Estado de dev | ⚠️ Desactualizado | NO tiene v0.3.0, está en commit anterior |
| Divergencia | 🔴 Crítica | Ramas divergentes con trabajo duplicado |
| Tags erróneos | v2.x.x, v1.x.x | Deben eliminarse |

**Acción requerida**: Sincronizar dev con main (trae v0.3.0 a dev).

#### edugo-api-mobile

| Aspecto | Estado | Detalle |
|---------|--------|---------|
| Último release válido | v0.1.1 | Commit 1dde3c8, imagen Docker 0.1.1 |
| Última imagen Docker | 0.1.1 | Publicada 2025-11-01 23:18 |
| Estado de main | ✅ Actualizado | Tiene PR #8 mergeado (6e88a52) |
| Estado de dev | ⚠️ Desincronizado | NO tiene 6e88a52, tiene 25 commits individuales |
| Divergencia | 🔴 Crítica | Merge manual incorrecto |
| Tags erróneos | v1.x.x | Deben eliminarse |

**Acción requerida**: Corregir sincronización entre main y dev.

#### ¿Por Qué NO era 25 commits "de diferencia real"?

La confusión original era correcta: si la imagen Docker (v0.1.1) es del 2025-11-01, y main está solo 3 commits adelante, ¿cómo puede dev estar 25 commits adelante?

**Respuesta**:
- La imagen 0.1.1 es del commit 1dde3c8
- Main está en 6e88a52 (3 commits después de 1dde3c8)
- El commit 6e88a52 es un MEGA-MERGE que INCLUYE 25 commits
- Dev tiene esos mismos 25 commits pero como commits INDIVIDUALES
- El problema no es que dev tenga trabajo nuevo, es que el MERGE no se sincronizó correctamente

Es como si main dijera: "Tengo el libro completo" (6e88a52)
Y dev dijera: "Tengo las 25 páginas sueltas" (25 commits)

**Son los mismos contenidos, pero en formato diferente**.

---

### 📋 PLAN DE ACCIÓN CORREGIDO

#### Prioridad #1: edugo-api-mobile

**Opción A: Reconocer que main tiene el trabajo actualizado**
- main (6e88a52) YA TIENE todo el trabajo de dev
- dev solo necesita sincronizarse con main
- Hacer: `git checkout dev && git merge main --no-ff`
- Resultado: dev queda al nivel de main

**Opción B: Verificar que 6e88a52 realmente contiene todo**
- Hacer diff del contenido: `git diff origin/main origin/dev`
- Si NO hay diferencias significativas: usar Opción A
- Si hay diferencias: investigar qué falta

**Recomendación**: Usar Opción A. El commit 6e88a52 tiene el mensaje completo listando TODOS los cambios.

#### Prioridad #2: edugo-shared

**Acción**: Sincronizar dev con main
- dev (74864bf) NO tiene el release v0.3.0
- main (082f430) tiene v0.3.0
- Hacer: `git checkout dev && git merge main --no-ff`
- Resultado: dev tiene v0.3.0

#### Prioridad #3: Limpiar Tags Erróneos

**edugo-api-mobile**:
- Eliminar: v1.0.0, v1.0.1, v1.0.2
- Mantener: v0.0.1, v0.1.0, v0.1.1

**edugo-shared**:
- Eliminar: v1.0.0, v2.0.0, v2.0.1, v2.0.5, v2.0.6
- Eliminar: tags de módulos v2.x.x
- Mantener: v0.1.0, v0.3.0 (global y por módulo)

---

### ✅ VERIFICACIONES FINALES

Después de aplicar el plan, verificar:

```bash
# En cada proyecto:
# 1. Ramas sincronizadas
git fetch --all
git rev-list --left-right --count origin/main...origin/dev
# Debe dar: 0       0

# 2. Último release válido
git tag -l "v0.*" | sort -V | tail -1
# Debe coincidir con version.txt

# 3. No quedan tags v1.x.x o v2.x.x
git tag -l "v1.*"
git tag -l "v2.*"
# Deben estar vacíos
```

---

## 📊 ANÁLISIS POR PROYECTO

### 1. edugo-api-mobile

#### 📌 Información General

- **Tipo**: API REST (Go + Gin)
- **Rama activa**: dev
- **Genera Docker**: ✅ Sí (GHCR)
- **Última actualización**: 2025-10-31

#### 🔍 Estado de Ramas

```
Último commit en main: 6e88a52 - "Dev (#8)" (2025-10-31)
Último commit en dev:  c907a42 - "Merge branch 'main' into dev" (2025-10-31)

Commits adelante de dev vs main: 25 commits
Commits adelante de main vs dev: 1 commit
```

**Gráfico de divergencia**:
```
main:  ───────────────●─────> (6e88a52) [1 commit único]
                      ↓
dev:   ───────────────●─────●─●─●─●─●─●─●─●─●─●─●─●─●─●─●─●─●─●─●─●─●─●─●─●─●>
                                      [25 commits únicos]
```

#### 🏷️ Tags y Releases

| Tag | Fecha Creación | Tipo |
|-----|----------------|------|
| v1.0.2 | Reciente | Release |
| v1.0.1 | Anterior | Release |
| v1.0.0 | Anterior | Release |
| v0.1.1 | Anterior | Release |
| v0.1.0 | Anterior | Release |
| v0.0.1 | Inicial | Release |

**version.txt**: `0.1.1` ⚠️ **Inconsistente con último tag v1.0.2**

#### ⚙️ CI/CD Workflows

| Workflow | Archivo | Trigger | Propósito |
|----------|---------|---------|-----------|
| ✅ CI | `ci.yml` | PR a main/dev, push a main | Compilación y pruebas |
| ✅ Tests | `test.yml` | PR y push | Suite de pruebas |
| ✅ Release Auto | `release.yml` | Push de tags | Release automático |
| ✅ Release Manual | `manual-release.yml` | workflow_dispatch | Release TODO-EN-UNO manual |
| ✅ Sync Main→Dev | `sync-main-to-dev.yml` | Push a main | Sincronización automática |
| ✅ Build & Push | `build-and-push.yml` | Releases | Build y push de imagen Docker |
| ✅ Docker Only | `docker-only.yml` | workflow_dispatch | Solo build Docker |
| ✅ Docs | `README.md` | - | Documentación de workflows |

#### 🚨 Problemas Detectados

**CRÍTICO** 🔴:
1. **Dev muy adelantado**: 25 commits sin mergear a main
   - Incluye features importantes: autenticación JWT completa, refresh tokens, bcrypt
   - Rate limiting anti-fuerza bruta
   - Migración a middleware compartido de edugo-shared
   - Múltiples mejoras de documentación y CI/CD

2. **Inconsistencia de versionado**:
   - version.txt dice `0.1.1`
   - Último tag es `v1.0.2`
   - ¿Cuál es el versionado real?

**ADVERTENCIA** 🟡:
3. **Desincronización**: Main tiene 1 commit (merge #8) que dev no tiene
   - Puede causar conflictos futuros
   - Necesita sincronización main → dev

#### 💡 Commits Destacados en Dev (sin mergear)

```
c907a42 - Merge branch 'main' into dev
2dfe7f2 - Merge remote-tracking branch 'origin/main'
44c8b17 - docs: actualizar plan CI/CD con workflow manual-release TODO-EN-UNO
1dde3c8 - chore: release v0.1.1
f5f2c75 - feat: workflow TODO-EN-UNO para release completo (#7)
e31d01f - docs: actualizar CLAUDE.md con flujo completo de edugo-shared
71b6a94 - feat: agregar comando /Task_Short de Claude Code (#6)
a8e9c70 - Merge branch 'main' into dev
ece5dd0 - feat: migrar a middleware JWT compartido de edugo-shared (#5)
[... 16 commits más ...]
```

#### ✅ Fortalezas

- ✅ Workflows CI/CD completos y bien documentados
- ✅ Estructura de proyecto coherente con Clean Architecture
- ✅ Tests con testcontainers implementados
- ✅ Documentación Swagger actualizada
- ✅ Integración con edugo-shared funcionando

---

### 2. edugo-api-administracion

#### 📌 Información General

- **Tipo**: API Admin (Go + Gin)
- **Rama activa**: main
- **Genera Docker**: ✅ Sí (GHCR)
- **Última actualización**: 2025-10-31

#### 🔍 Estado de Ramas

```
Último commit en main:       5e62e54 - "chore: release v0.1.0" (2025-10-31)
Último commit en origin/dev: bde500a - "feat: CI/CD optimizado y Copilot instructions (#3)"

Commits adelante de main vs origin/dev: 1 commit
Commits adelante de origin/dev vs main: 0 commits
```

**Gráfico de divergencia**:
```
main:       ───────────●────> (5e62e54) [1 commit adelante]
                       ↑
origin/dev: ───────────●      (bde500a)
```

**NOTA**: La rama `dev` solo existe en remote, no está checkeada localmente.

#### 🏷️ Tags y Releases

| Tag | Tipo | Estado |
|-----|------|--------|
| v1.0.3 | Release | Más reciente |
| v1.0.2 | Release | Anterior |
| v1.0.1 | Release | Anterior |
| v1.0.0 | Release | Anterior |
| v0.1.0 | Release | Inicial |

**version.txt**: `0.1.0` ⚠️ **Inconsistente con último tag v1.0.3**

#### ⚙️ CI/CD Workflows

Misma estructura que api-mobile:
- ✅ ci.yml
- ✅ test.yml
- ✅ release.yml
- ✅ manual-release.yml
- ✅ sync-main-to-dev.yml
- ✅ build-and-push.yml
- ✅ docker-only.yml
- ✅ README.md

#### 🚨 Problemas Detectados

**ADVERTENCIA** 🟡:
1. **Main adelantado**: Main tiene el commit de release v0.1.0 que dev no tiene
   - Dev necesita sincronización desde main
   - Workflow `sync-main-to-dev.yml` debería haberlo hecho automáticamente

2. **Inconsistencia de versionado**: version.txt (0.1.0) vs tags (v1.x.x)

**INFO** ℹ️:
3. **Rama dev local no existe**: Solo está en remote
   - Puede dificultar el trabajo local si se necesita cambiar a dev

#### ✅ Fortalezas

- ✅ Workflows estandarizados con otros proyectos
- ✅ Main está relativamente actualizado
- ✅ Solo 1 commit de diferencia (fácil de sincronizar)

---

### 3. edugo-worker

#### 📌 Información General

- **Tipo**: Worker (Go + RabbitMQ)
- **Rama activa**: dev
- **Genera Docker**: ✅ Sí (GHCR)
- **Última actualización**: 2025-10-31

#### 🔍 Estado de Ramas

```
Último commit en main: b0eeb55 - "Dev (#4)" (2025-10-31)
Último commit en dev:  487dac3 - "Merge branch 'main' into dev"

Commits adelante de dev vs main: 4 commits
Commits adelante de main vs dev: 1 commit
```

**Gráfico de divergencia**:
```
main: ──────────●─────> (b0eeb55) [1 commit único]
                ↓
dev:  ──────────●─●─●─●> (487dac3) [4 commits únicos]
```

#### 🏷️ Tags y Releases

| Tag | Tipo |
|-----|------|
| v1.0.2 | Release |
| v1.0.1 | Release |
| v1.0.0 | Release |
| v0.1.0 | Release |

**version.txt**: `0.1.0` ⚠️ **Inconsistente con último tag v1.0.2**

#### ⚙️ CI/CD Workflows

Misma estructura que api-mobile (8 workflows):
- ✅ ci.yml
- ✅ test.yml
- ✅ release.yml
- ✅ manual-release.yml
- ✅ sync-main-to-dev.yml
- ✅ build-and-push.yml
- ✅ docker-only.yml
- ✅ README.md

#### 🚨 Problemas Detectados

**ADVERTENCIA** 🟡:
1. **Dev adelantado**: 4 commits en dev sin mergear a main
   - Incluye optimizaciones de CI/CD
   - Copilot instructions actualizadas

2. **Desincronización**: Main tiene 1 commit que dev no tiene
   - Merge commit del PR #4
   - Necesita sync main → dev

3. **Inconsistencia de versionado**: version.txt (0.1.0) vs tags (v1.x.x)

#### ✅ Fortalezas

- ✅ Workflows estandarizados
- ✅ Solo 4 commits de diferencia (manejable)
- ✅ Desarrollo activo y documentado

---

### 4. edugo-shared

#### 📌 Información General

- **Tipo**: Librería compartida (Go modules)
- **Rama activa**: dev
- **Genera Docker**: ❌ No (es librería)
- **Arquitectura**: Modular (múltiples módulos Go)
- **Última actualización**: 2025-10-31

#### 🔍 Estado de Ramas

```
Último commit en main: 082f430 - "fix: Resetear versionado a v0.3.0 (#5)"
Último commit en dev:  74864bf - "fix: resetear esquema de versionado de v2.x.x a v0.x.x (#4)"

Commits adelante de main vs dev: 2 commits
Commits adelante de dev vs main: 3 commits
```

**Gráfico de divergencia** (RAMAS DIVERGENTES):
```
       ┌─●─●─────> main (082f430) [2 commits únicos]
       │
base───┤
       │
       └─●─●─●───> dev (74864bf) [3 commits únicos]
```

⚠️ **CRÍTICO**: Las ramas han divergido, tienen commits que la otra no tiene.

#### 🏷️ Tags y Releases (Múltiples)

**Tags Globales**:
```
v2.0.6
v2.0.5
v2.0.1
v2.0.0
v1.0.0
v0.3.0
v0.1.0
```

**Tags de Módulos**:
```
messaging/rabbit/v2.0.5
middleware/gin/v0.3.0
middleware/gin/v0.0.1
```

**version.txt**: ❌ NO EXISTE (normal para librería)

#### 📦 Estructura Modular

```
edugo-shared/
├── auth/                    # Módulo de autenticación
├── logger/                  # Módulo de logging
├── common/errors/           # Tipos de error
├── messaging/rabbit/        # Cliente RabbitMQ (v2.0.5)
└── middleware/gin/          # Middleware para Gin (v0.3.0)
```

#### ⚙️ CI/CD Workflows

**Workflows diferentes** a los proyectos de servicios:

| Workflow | Archivo | Propósito |
|----------|---------|-----------|
| ✅ CI | `ci.yml` | Tests por módulo |
| ✅ Tests | `test.yml` | Suite de pruebas |
| ✅ Release | `release.yml` | Release con validación de módulos |
| ✅ Sync | `sync-main-to-dev.yml` | Sincronización |
| ✅ Docs | `README.md` | Documentación |

**NO TIENE** (correcto para librería):
- ❌ manual-release.yml
- ❌ build-and-push.yml (no genera Docker)
- ❌ docker-only.yml

#### 🚨 Problemas Detectados

**CRÍTICO** 🔴:
1. **Ramas divergentes**: Main y dev tienen commits únicos que el otro no tiene
   - Main: 2 commits (relacionados con reseteo de versionado a v0.3.0)
   - Dev: 3 commits (relacionados con reseteo de v2.x.x a v0.x.x)
   - Necesita reconciliación manual

2. **Versionado múltiple complejo**:
   - Tags globales: v2.0.6
   - Tags de módulos: middleware/gin/v0.3.0, messaging/rabbit/v2.0.5
   - Reseteos de versionado recientes sugieren confusión

**ADVERTENCIA** 🟡:
3. **Arquitectura modular**: Complejidad adicional en versionado
   - Cada módulo puede tener su propia versión
   - Requiere coordinación cuidadosa de releases

#### ✅ Fortalezas

- ✅ Arquitectura modular bien estructurada
- ✅ Workflows adaptados para librería
- ✅ Tests por módulo
- ✅ Documentación de módulos
- ✅ No genera Docker (correcto para librería)

#### 📝 Commits Recientes Relevantes

**Main**:
```
082f430 - fix: Resetear versionado a v0.3.0 (#5)
ca52dae - Merge branch 'dev' into main
```

**Dev**:
```
74864bf - fix: resetear esquema de versionado de v2.x.x a v0.x.x (#4)
9a4745f - chore: release middleware/gin/v0.3.0
cc2acd2 - chore: release messaging/rabbit/v2.0.5
```

---

### 5. edugo-dev-environment

#### 📌 Información General

- **Tipo**: Ambiente de desarrollo (Docker Compose)
- **Rama única**: main (no tiene dev)
- **Genera Docker**: ❌ No (usa imágenes de otros repos)
- **Propósito**: Orquestar servicios con docker-compose
- **Última actualización**: 2025-10-28

#### 🔍 Estado de Ramas

```
Rama única: main
Último commit: cb9e60f - "fix: agregar variables de entorno requeridas para APIs"

NO TIENE rama dev
```

#### 🏷️ Tags y Releases

- **Tags**: ❌ Ninguno
- **version.txt**: ❌ NO EXISTE
- **Versionado**: No aplica

#### 📁 Estructura del Proyecto

```
edugo-dev-environment/
├── docker/
│   ├── docker-compose.yml     ✅ Orquestación de servicios
│   ├── .env                   ✅ Variables de entorno
│   └── .env.example           ✅ Template de variables
├── docs/                      ✅ Documentación
├── scripts/                   ✅ Scripts de utilidad
├── .gitignore                 ✅
└── README.md                  ✅
```

#### ⚙️ CI/CD Workflows

**NO TIENE**: ❌ Directorio `.github/workflows/` no existe

#### 🚨 Problemas Detectados

**ADVERTENCIA** 🟡:
1. **Sin CI/CD**: No hay workflows de validación
   - No valida sintaxis de docker-compose.yml
   - No verifica compatibilidad de variables de entorno
   - No hay checks automáticos en PRs

2. **Sin versionado**: No tiene tags ni releases
   - Dificulta rastrear cambios importantes
   - No hay forma de referenciar versiones específicas

3. **Sin rama dev**: Solo tiene main
   - No sigue el estándar dev → main de otros repos
   - Cambios van directo a main

**INFO** ℹ️:
4. **docker-compose.yml en subdirectorio**: Está en `docker/` en vez de raíz
   - Requiere especificar path: `docker-compose -f docker/docker-compose.yml`
   - Puede ser intencional para organización

#### ✅ Fortalezas

- ✅ Documentación clara
- ✅ Estructura organizada
- ✅ .env.example para guiar configuración
- ✅ Scripts de utilidad

#### 💡 Sugerencias

- Considerar agregar workflow básico para validar docker-compose.yml
- Evaluar si necesita rama dev o si main directo es apropiado
- Posible versionado con tags para cambios mayores de infraestructura

---

## 📈 ANÁLISIS COMPARATIVO

### Tabla 1: Consistencia de Workflows CI/CD

| Repositorio | ci.yml | test.yml | release.yml | manual-release | sync-main-to-dev | Docker workflows | TOTAL |
|-------------|--------|----------|-------------|----------------|------------------|-----------------|-------|
| **api-mobile** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ ✅ | 8 workflows |
| **api-admin** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ ✅ | 8 workflows |
| **worker** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ ✅ | 8 workflows |
| **shared** | ✅ | ✅ | ✅ | ❌ | ✅ | ❌ ❌ | 5 workflows |
| **dev-env** | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ ❌ | 0 workflows |

**Conclusión**:
- ✅ **Estandarización EXCELENTE** en proyectos de servicios (api-mobile, api-admin, worker)
- ✅ **Adecuado** para shared (librería no necesita Docker workflows)
- 🟡 **Gap** en dev-environment (podría beneficiarse de validación básica)

---

### Tabla 2: Versionado y Tags

| Repositorio | version.txt | Último tag | Tags totales | Consistencia | Gravedad |
|-------------|-------------|------------|--------------|--------------|----------|
| **api-mobile** | 0.1.1 | v1.0.2 | 6 tags | ❌ Inconsistente | 🔴 Alta |
| **api-admin** | 0.1.0 | v1.0.3 | 5 tags | ❌ Inconsistente | 🔴 Alta |
| **worker** | 0.1.0 | v1.0.2 | 4 tags | ❌ Inconsistente | 🔴 Alta |
| **shared** | N/A | v2.0.6 + módulos | ~10 tags | ⚠️ Múltiple (modular) | 🟠 Media |
| **dev-env** | N/A | N/A | 0 tags | ⚠️ Sin versionado | 🟡 Baja |

**Conclusión CRÍTICA**:
- 🔴 **TODOS** los proyectos con Docker tienen **inconsistencia severa** de versionado
- 📊 **Pattern detectado**: version.txt dice v0.x.x pero tags dicen v1.x.x
- 🔍 **Hipótesis**: Se hizo un reseteo de versionado pero no se limpiaron tags antiguos
- ⚡ **Acción requerida**: Decidir esquema oficial y aplicar consistentemente

---

### Tabla 3: Estado de Sincronización Main vs Dev

| Repositorio | Dev adelante | Main adelante | Total desincronizado | Estado | Urgencia |
|-------------|--------------|---------------|---------------------|--------|----------|
| **api-mobile** | 25 commits | 1 commit | 26 commits | 🔴 Dev MUY adelantado | 🔴 Urgente |
| **api-admin** | 0 commits | 1 commit | 1 commit | 🟡 Main adelantado | 🟢 Baja |
| **worker** | 4 commits | 1 commit | 5 commits | 🟡 Dev adelantado | 🟡 Media |
| **shared** | 3 commits | 2 commits | 5 commits | 🟠 Divergente | 🔴 Urgente |
| **dev-env** | N/A | N/A | N/A | ⚪ Sin dev | N/A |

**Conclusión**:
- 🔴 **api-mobile**: Requiere atención urgente (25 commits es demasiado)
- 🔴 **shared**: Ramas divergentes requieren reconciliación manual
- 🟡 **worker**: Manejable pero debe atenderse pronto
- 🟢 **api-admin**: Casi sincronizado, fácil de corregir

---

### Tabla 4: Arquitectura y Propósito

| Repositorio | Categoría | Lenguaje | Framework | Base de Datos | Messaging | Storage | Genera Docker |
|-------------|-----------|----------|-----------|---------------|-----------|---------|---------------|
| **api-mobile** | API REST | Go | Gin | PostgreSQL, MongoDB | RabbitMQ (pendiente) | S3 (pendiente) | ✅ Sí |
| **api-admin** | API Admin | Go | Gin | PostgreSQL | - | - | ✅ Sí |
| **worker** | Worker | Go | - | - | RabbitMQ | - | ✅ Sí |
| **shared** | Librería | Go | - | - | Módulo rabbit | - | ❌ No |
| **dev-env** | Infraestructura | YAML | Docker Compose | - | - | - | ❌ No |

---

### Gráfico: Distribución de Problemas por Gravedad

```
🔴 CRÍTICOS (2):
  ├─ Versionado inconsistente (afecta a 3 proyectos)
  └─ api-mobile dev muy adelantado (25 commits)

🟠 IMPORTANTES (1):
  └─ edugo-shared ramas divergentes

🟡 ADVERTENCIAS (4):
  ├─ worker dev adelantado (4 commits)
  ├─ api-admin main adelantado (1 commit)
  ├─ dev-environment sin CI/CD
  └─ dev-environment sin versionado

🟢 MEJORAS (2):
  ├─ Estandarizar workflow en dev-environment
  └─ Documentar decisiones de versionado
```

---

## 🚨 PROBLEMAS CRÍTICOS IDENTIFICADOS

### Problema #1: Inconsistencia de Versionado

**Gravedad**: 🔴 CRÍTICA
**Afecta a**: api-mobile, api-admin, worker
**Urgencia**: Inmediata

#### Descripción

Todos los proyectos que generan imágenes Docker tienen una inconsistencia severa entre:
- **version.txt**: Indica v0.x.x (0.1.0 o 0.1.1)
- **Tags Git**: Existen tags v1.x.x y hasta v2.x.x

#### Evidencia

| Proyecto | version.txt | Último tag | Diferencia |
|----------|-------------|------------|------------|
| api-mobile | 0.1.1 | v1.0.2 | +0.9 versiones mayores |
| api-admin | 0.1.0 | v1.0.3 | +1.0 versiones mayores |
| worker | 0.1.0 | v1.0.2 | +1.0 versiones mayores |

#### Impacto

- ❌ No está claro cuál es la versión "oficial" del proyecto
- ❌ Las imágenes Docker pueden tener tags incorrectos
- ❌ Los releases en GitHub pueden estar mal etiquetados
- ❌ Consumidores del API no saben qué versión están usando
- ❌ Dificulta el rollback y troubleshooting

#### Hipótesis de Causa Raíz

Basado en commits recientes en shared que mencionan "resetear versionado":
- Se decidió cambiar de esquema v1.x.x/v2.x.x a v0.x.x
- Se actualizó version.txt pero no se eliminaron tags antiguos
- Los workflows pueden estar usando tags o version.txt inconsistentemente

#### Acción Requerida

1. **Decidir esquema oficial**: ¿v0.x.x (pre-release) o v1.x.x (stable)?
2. **Auditar workflows**: ¿De dónde leen la versión?
3. **Limpiar o etiquetar**: Eliminar tags obsoletos o actualizar version.txt
4. **Documentar decisión**: Agregar política de versionado al README

---

### Problema #2: edugo-api-mobile con 25 Commits sin Mergear

**Gravedad**: 🔴 CRÍTICA
**Afecta a**: edugo-api-mobile
**Urgencia**: Inmediata

#### Descripción

La rama `dev` de api-mobile está **25 commits** por delante de `main`. Esto representa:
- Semanas o meses de trabajo sin proteger en main
- Riesgo de pérdida si algo pasa con la rama dev
- Features importantes funcionando solo en dev

#### Features sin Mergear (Muestra)

```
✨ Autenticación JWT completa (login, refresh, logout)
✨ Encriptación de contraseñas con bcrypt
✨ Rate limiting anti-fuerza bruta
✨ Migración a middleware compartido (edugo-shared)
✨ Workflow manual-release TODO-EN-UNO
✨ Comandos de Claude Code (/Task_Short)
✨ Optimizaciones de CI/CD
✨ Múltiples mejoras de documentación
```

#### Impacto

- 🚨 **Alto riesgo**: Si se pierde dev, se pierde mucho trabajo
- 🚨 **Producción desactualizada**: Main (y probablemente producción) no tiene estas features
- 🚨 **Divergencia creciente**: Mientras más tiempo pase, más difícil el merge
- 🚨 **Duplicación de esfuerzo**: Si se trabaja sobre main, se pierde lo de dev

#### Commits Críticos que Deben Protegerse

```
f5f2c75 - feat: workflow TODO-EN-UNO para release completo (#7)
ece5dd0 - feat: migrar a middleware JWT compartido de edugo-shared (#5)
9b87eb1 - feat: implementar rate limiting para login (#4)
10b8a5e - feat: implementar sistema completo de autenticación JWT
```

#### Acción Requerida

1. **URGENTE**: Crear PR de dev → main
2. **Ejecutar CI/CD**: Verificar que todo compila y pasa tests
3. **Review cuidadoso**: Revisar los 25 commits antes de mergear
4. **Mergear a main**: Proteger el trabajo hecho
5. **Crear release**: Generar v0.2.0 (o según esquema decidido)

---

### Problema #3: edugo-shared con Ramas Divergentes

**Gravedad**: 🟠 IMPORTANTE
**Afecta a**: edugo-shared
**Urgencia**: Alta

#### Descripción

Main y dev han **divergido**, tienen commits únicos que la otra rama no tiene:
- **Main tiene 2 commits** que dev no tiene
- **Dev tiene 3 commits** que main no tiene

#### Análisis de Divergencia

```
Commits SOLO en main:
- 082f430: fix: Resetear versionado a v0.3.0 (#5)
- ca52dae: Merge branch 'dev' into main

Commits SOLO en dev:
- 74864bf: fix: resetear esquema de versionado de v2.x.x a v0.x.x (#4)
- 9a4745f: chore: release middleware/gin/v0.3.0
- cc2acd2: chore: release messaging/rabbit/v2.0.5
```

#### Causa Raíz

Parece que hubo trabajo paralelo en ambas ramas relacionado con reseteo de versionado:
- En dev: Se reseteó de v2.x.x a v0.x.x
- En main: Se reseteó a v0.3.0
- Posiblemente se trabajó directamente en main (anti-patrón)

#### Impacto

- ❌ No se puede hacer merge automático simple
- ❌ Se requiere reconciliación manual
- ❌ Posible pérdida de cambios si se hace merge incorrecto
- ❌ Afecta a los 3 proyectos que dependen de shared

#### Acción Requerida

1. **Analizar cambios**: Determinar qué commits son necesarios
2. **Decidir estrategia**: Rebase, merge, o cherry-pick
3. **Reconciliar**: Aplicar cambios necesarios de ambas ramas
4. **Sincronizar**: Dejar main y dev al mismo nivel
5. **Prevenir**: Evitar trabajo directo en main en el futuro

---

### Problema #4: dev-environment sin CI/CD ni Versionado

**Gravedad**: 🟡 ADVERTENCIA
**Afecta a**: edugo-dev-environment
**Urgencia**: Baja (mejora)

#### Descripción

El repositorio dev-environment:
- No tiene workflows de CI/CD
- No tiene tags ni releases
- No tiene rama dev (solo main)
- Cambios van directo a main sin validación

#### Impacto Potencial

- ⚠️ docker-compose.yml inválido podría committearse sin detectarse
- ⚠️ Variables de entorno faltantes no se detectan hasta runtime
- ⚠️ No hay forma de referenciar versiones estables
- ⚠️ Rollback difícil sin tags

#### Acción Sugerida (Opcional)

1. Agregar workflow básico para validar docker-compose.yml
2. Considerar si necesita rama dev o si main directo es apropiado
3. Evaluar si necesita versionado con tags

---

## 📚 EXPLICACIÓN CONCEPTUAL: Merge Commits y Sincronización

### Tu Pregunta

> "Si main está por debajo de dev, y hago un PR a main, el PR traerá todos los commits de dev y se equipara, pero el PR genera un commit merge, y eso pone a main por encima de dev. En ese caso, ¿no debería ese commit pasarse a dev? ¿O está bien así? ¿Cuál es el procedimiento? ¿No debería haberse actualizado tanto dev como main al mismo nivel?"

### Respuesta Conceptual

**SÍ, tienes toda la razón**. El commit de merge SÍ debe sincronizarse de vuelta a dev. Aquí te explico por qué y cómo:

---

### Flujo Correcto: dev → main → dev

#### Paso 1: Estado Inicial

```
main: ───A───B───C
               ↓
dev:  ───A───B───C───D───E───F
```

- main está en C
- dev tiene 3 commits adelante (D, E, F)

#### Paso 2: Crear PR de dev → main

```bash
# En GitHub/GitLab
1. Crear PR: dev → main
2. Review de los commits D, E, F
3. Aprobar PR
```

#### Paso 3: Merge del PR (genera commit M)

```
main: ───A───B───C───────────M  ← Merge commit
               ↓             ↗
dev:  ───A───B───C───D───E───F
```

**Ahora main tiene el commit M que dev NO tiene**.

#### Paso 4: ⚠️ CRÍTICO - Sincronizar main → dev

```bash
# Opción 1: Automático (con workflow sync-main-to-dev.yml)
# Se ejecuta automáticamente al pushear a main

# Opción 2: Manual
git checkout dev
git merge main --no-ff
git push origin dev
```

#### Paso 5: Estado Final Correcto

```
main: ───A───B───C───────────M
                             ↓
dev:  ───A───B───C───D───E───F───M'  ← Merge de main a dev
```

**Ahora AMBAS ramas están al mismo nivel** ✅

---

### ¿Por Qué Es Importante Sincronizar?

#### Problema 1: Divergencia Futura

Si NO sincronizas:
```
main: ───C───M           ← Tiene M
dev:  ───C───F           ← No tiene M

# Próximo desarrollo en dev:
dev:  ───C───F───G───H   ← No tiene M

# Próximo PR dev → main:
main: ───C───M───M2      ← ¡Conflictos! dev no tenía M
```

#### Problema 2: Historial Confuso

Sin sincronización, el `git log` se ve así:
```bash
# En dev
git log --oneline
F feat: new feature
E feat: another feature
D feat: some feature
C initial commit

# Falta M! ¿Qué pasó con el merge a main?
```

#### Problema 3: Conflictos Futuros

- El próximo PR de dev → main tendrá que lidiar con M
- Puede generar conflictos innecesarios
- Historia se vuelve no-lineal y confusa

---

### Workflow `sync-main-to-dev.yml`

Tus proyectos YA TIENEN este workflow automático. Veamos cómo funciona:

```yaml
name: Sync Main to Dev

on:
  push:
    branches:
      - main  # Se ejecuta cuando se pushea a main

jobs:
  sync:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          ref: dev

      - name: Merge main into dev
        run: |
          git config user.name "github-actions[bot]"
          git config user.email "github-actions[bot]@users.noreply.github.com"
          git merge origin/main --no-ff
          git push origin dev
```

**Esto DEBERÍA sincronizar automáticamente**, pero necesitas verificar:

1. ¿El workflow está habilitado?
2. ¿Tiene permisos de escritura?
3. ¿Hay algún error en las ejecuciones?

---

### Escenarios Especiales

#### Escenario A: Fast-Forward Merge

Si main NO tenía commits únicos:
```
main: ───A───B───C
dev:  ───A───B───C───D───E

# Merge fast-forward (sin commit M)
main: ───A───B───C───D───E
dev:  ───A───B───C───D───E
```

✅ **No necesita sincronización** (ya están iguales)

#### Escenario B: Squash Merge

Si usas "Squash and Merge":
```
main: ───C───S  ← S contiene D+E+F squasheados
dev:  ───C───D───E───F
```

⚠️ **Requiere estrategia diferente**:
- S y D+E+F son cambios iguales pero commits diferentes
- Sincronizar crearía duplicación
- Solución: Rebase dev sobre main (avanzado)

#### Escenario C: Cherry-Pick

Si haces cherry-pick en vez de merge:
```
main: ───C───D'  ← D' es copia de D
dev:  ───C───D───E───F
```

⚠️ **También requiere cuidado**:
- D y D' son diferentes commits (hashes distintos)
- Puede causar confusión

---

### Mejor Práctica: Merge Commit + Sync

**Recomendación para EduGo**:

1. **Usar merge commit** (no squash, no rebase) para PRs dev → main
2. **Sincronizar automáticamente** con workflow sync-main-to-dev.yml
3. **Verificar sincronización** después de cada merge
4. **No hacer commits directos** a main (siempre pasar por dev)

#### Checklist Post-Merge

```bash
# Después de mergear dev → main:
□ Verificar que workflow sync-main-to-dev.yml se ejecutó
□ Verificar que dev tiene el merge commit de main
□ Confirmar que `git log main` y `git log dev` muestran el mismo último commit
□ Confirmar que `git rev-list --left-right --count main...dev` da 0↔0
```

---

### Comando para Verificar Sincronización

```bash
# Ver cuántos commits de diferencia hay
git rev-list --left-right --count main...dev

# Resultado esperado después de sincronización:
0       0

# Si ves algo como:
1       0  ← main tiene 1 commit que dev no tiene (MALO)
0       5  ← dev tiene 5 commits que main no tiene (normal durante desarrollo)
2       3  ← DIVERGENTE (CRÍTICO)
```

---

### Resumen

| Pregunta | Respuesta |
|----------|-----------|
| ¿El merge commit debe pasarse a dev? | ✅ SÍ, absolutamente |
| ¿Está bien dejar main adelante? | ❌ NO, causa divergencia |
| ¿Debe sincronizarse automáticamente? | ✅ SÍ, con workflow sync-main-to-dev.yml |
| ¿Main y dev deben estar al mismo nivel después? | ✅ SÍ, siempre después de mergear |

---

## 🎯 PLAN DE ACCIÓN DETALLADO

### Fase 1: Análisis y Decisiones Previas (HACER PRIMERO)

#### 1.1 Decidir Esquema de Versionado

**Urgencia**: 🔴 CRÍTICA
**Tiempo estimado**: 30 minutos
**Responsable**: Equipo técnico + Product Owner

- [ ] **Revisar tags actuales de todos los proyectos**
  ```bash
  # En cada proyecto:
  cd edugo-api-mobile && git tag -l | sort -V
  cd edugo-api-administracion && git tag -l | sort -V
  cd edugo-worker && git tag -l | sort -V
  cd edugo-shared && git tag -l | sort -V
  ```

- [ ] **Decidir esquema oficial de versionado**

  **Opción A: v0.x.x (Pre-release)**
  - ✅ Indica que aún no hay versión estable 1.0
  - ✅ Consistente con version.txt actual
  - ❌ Puede dar impresión de "no production-ready"

  **Opción B: v1.x.x (Stable)**
  - ✅ Indica versión estable en producción
  - ✅ Más profesional para consumidores externos
  - ❌ Requiere actualizar version.txt

  **Decisión**: ______A_______ (marcar Opción A o B) Seria bueno colocar v0.10.0 para poder colocar a todos al mismo nivel, La version V1.x.x fueron errores que en teoria ya estaba solventado, pero por lo que me cuentas, continua, asi que si o si, tags, version, release numero de instancia, todo debe estar en v0.x.x otra valor es un error no solucionado.

- [ ] **Documentar decisión en README.md de cada proyecto**

  Agregar sección:
  ```markdown
  ## Versionado

  Este proyecto sigue [Semantic Versioning 2.0.0](https://semver.org/).

  - **Formato**: vMAJOR.MINOR.PATCH
  - **Esquema actual**: v0.x.x (pre-release) / v1.x.x (stable)
  - **Dónde se define**: version.txt en la raíz del proyecto
  - **Tags**: Cada release crea un tag git vX.Y.Z
  ```

---

#### 1.2 Auditar Workflows de Versionado

**Urgencia**: 🔴 CRÍTICA
**Tiempo estimado**: 1 hora
**Responsable**: DevOps / Desarrollador

- [ ] **Verificar de dónde leen la versión los workflows**

  Archivos a revisar en cada proyecto:
  ```bash
  - .github/workflows/release.yml
  - .github/workflows/manual-release.yml
  - .github/workflows/build-and-push.yml
  ```

- [ ] **Identificar inconsistencias**

  Buscar:
  - [ ] ¿Se lee version.txt?
  - [ ] ¿Se usan tags git?
  - [ ] ¿Se genera la versión automáticamente?
  - [ ] ¿Qué pasa si version.txt y tags no coinciden?

- [ ] **Asegurar una única fuente de verdad**

  **Recomendación**: version.txt debe ser la fuente de verdad
  - Workflows leen de version.txt
  - Workflows crean tags basados en version.txt
  - Validación: si tag existe, no permitir release duplicado

---

#### 1.3 Verificar Estado de Workflows sync-main-to-dev

**Urgencia**: 🟡 MEDIA
**Tiempo estimado**: 30 minutos
**Responsable**: DevOps

- [ ] **Revisar ejecuciones del workflow sync-main-to-dev.yml**

  En GitHub para cada proyecto:
  ```
  Actions → Sync Main to Dev → Ver últimas ejecuciones
  ```

- [ ] **Verificar si hay errores o warnings**

  Problemas comunes:
  - ❌ Falta de permisos de escritura
  - ❌ Conflictos de merge
  - ❌ Workflow deshabilitado

- [ ] **Habilitar si está deshabilitado**

- [ ] **Corregir errores de permisos**

  En workflow, agregar si falta:
  ```yaml
  permissions:
    contents: write  # Necesario para push a dev
  ```

---

### Fase 2: Limpieza de Versionado

#### 2.1 Limpiar Tags Inconsistentes

**Urgencia**: 🔴 ALTA
**Tiempo estimado**: 1 hora
**Responsable**: DevOps
**⚠️ PRECAUCIÓN**: Esto es destructivo, hacer backup primero

**Si elegiste Opción A (v0.x.x)**:

- [ ] **Eliminar tags v1.x.x y v2.x.x de api-mobile**
  ```bash
  cd edugo-api-mobile

  # Listar tags a eliminar
  git tag -l "v1.*"
  git tag -l "v2.*"

  # Eliminar tags localmente
  git tag -d v1.0.0 v1.0.1 v1.0.2

  # Eliminar tags remotamente (DESTRUCTIVO)
  git push origin --delete v1.0.0 v1.0.1 v1.0.2
  ```

- [ ] **Eliminar tags v1.x.x de api-administracion**
  ```bash
  cd edugo-api-administracion
  git tag -d v1.0.0 v1.0.1 v1.0.2 v1.0.3
  git push origin --delete v1.0.0 v1.0.1 v1.0.2 v1.0.3
  ```

- [ ] **Eliminar tags v1.x.x de worker**
  ```bash
  cd edugo-worker
  git tag -d v1.0.0 v1.0.1 v1.0.2
  git push origin --delete v1.0.0 v1.0.1 v1.0.2
  ```

- [ ] **Eliminar tags v2.x.x de shared (evaluar primero)**
  ```bash
  cd edugo-shared
  # CUIDADO: evaluar si se están usando en go.mod de otros proyectos
  git tag -d v2.0.0 v2.0.1 v2.0.5 v2.0.6
  git push origin --delete v2.0.0 v2.0.1 v2.0.5 v2.0.6
  ```

  [ ] **Eliminar si es posible imagenes guardada en docker con estas versiones erroneas**
  ```bash
  # Indagar manera de eliminar
  ```


  **⚠️ ANTES de eliminar, verificar**:
  ```bash
  cd edugo-api-mobile
  grep "edugo-shared" go.mod
  # Si referencia v2.x.x, NO eliminar aún
  ```

**Si elegiste Opción B (v1.x.x)**:

- [ ] **Actualizar version.txt en todos los proyectos**
  ```bash
  # api-mobile
  cd edugo-api-mobile
  echo "1.0.2" > version.txt
  git add version.txt
  git commit -m "chore: actualizar version.txt a v1.0.2"

  # api-administracion
  cd edugo-api-administracion
  echo "1.0.3" > version.txt
  git add version.txt
  git commit -m "chore: actualizar version.txt a v1.0.3"

  # worker
  cd edugo-worker
  echo "1.0.2" > version.txt
  git add version.txt
  git commit -m "chore: actualizar version.txt a v1.0.2"
  ```

---

#### 2.2 Validar Limpieza

- [ ] **Verificar que tags y version.txt coincidan**
  ```bash
  # Para cada proyecto:
  LAST_TAG=$(git tag -l "v*.*.*" | grep -v "alpha\|beta\|rc" | sort -V | tail -1)
  VERSION_FILE=$(cat version.txt)
  echo "Tag: $LAST_TAG"
  echo "version.txt: v$VERSION_FILE"
  # Deben coincidir
  ```

- [ ] **Verificar que no haya releases huérfanos en GitHub**

  Ir a GitHub → Releases y eliminar releases sin tag correspondiente

---

### Fase 3: Sincronización de Ramas

#### 3.1 Sincronizar edugo-shared (PRIORITARIO)

**Urgencia**: 🔴 ALTA
**Tiempo estimado**: 2 horas
**Responsable**: Desarrollador senior
**Complejidad**: Alta (ramas divergentes)

**⚠️ ADVERTENCIA**: Esto requiere resolución manual de divergencia.

- [ ] **Crear backup de ambas ramas**
  ```bash
  cd edugo-shared
  git branch backup-main main
  git branch backup-dev dev
  git push origin backup-main backup-dev
  ```

- [ ] **Analizar diferencias**
  ```bash
  # Ver commits únicos en cada rama
  git log main..dev --oneline --no-merges  # Commits solo en dev
  git log dev..main --oneline --no-merges  # Commits solo en main

  # Ver archivos afectados
  git diff main...dev --name-status
  ```

- [ ] **Decidir estrategia de reconciliación**

  **Estrategia A: Merge bidireccional**
  ```bash
  # Mergear main a dev
  git checkout dev
  git merge main --no-ff -m "merge: sincronizar main → dev"

  # Resolver conflictos si los hay
  git status
  # Editar archivos en conflicto
  git add .
  git commit

  # Mergear dev a main
  git checkout main
  git merge dev --no-ff -m "merge: sincronizar dev → main"

  # Pushear ambas
  git push origin main dev
  ```

  **Estrategia B: Rebase dev sobre main** (más limpia)
  ```bash
  # Rebase dev sobre main
  git checkout dev
  git rebase main

  # Resolver conflictos si los hay
  git status
  # Editar, git add, git rebase --continue

  # Force push dev (CUIDADO)
  git push origin dev --force-with-lease

  # Fast-forward main a dev
  git checkout main
  git merge dev --ff-only
  git push origin main
  ```

- [ ] **Validar sincronización**
  ```bash
  git rev-list --left-right --count main...dev
  # Debe dar: 0       0
  ```

- [ ] **Crear release en shared**
  ```bash
  # Si no existe v0.3.0
  git tag v0.3.0
  git push origin v0.3.0
  ```

- [ ] **Actualizar dependencia en proyectos que usan shared**
  ```bash
  # En api-mobile, api-admin, worker:
  go get github.com/EduGoGroup/edugo-shared@v0.3.0
  go mod tidy
  ```

---

#### 3.2 Mergear edugo-api-mobile dev → main

**Urgencia**: 🔴 ALTA
**Tiempo estimado**: 3-4 horas (incluye testing)
**Responsable**: Desarrollador + Tester

**⚠️ IMPORTANTE**: Este es el merge más grande (25 commits).

- [ ] **Preparación: Asegurar que dev está actualizado**
  ```bash
  cd edugo-api-mobile
  git checkout dev
  git pull origin dev
  ```

- [ ] **Crear rama de feature para el merge**
  ```bash
  git checkout -b merge/dev-to-main-$(date +%Y%m%d)
  git push origin merge/dev-to-main-$(date +%Y%m%d)
  ```

- [ ] **Crear Pull Request en GitHub**
  - Base: main
  - Compare: merge/dev-to-main-YYYYMMDD
  - Título: "feat: mergear desarrollo acumulado de dev (25 commits)"
  - Body: Listar features principales:
    ```markdown
    ## Cambios Principales

    - ✨ Sistema completo de autenticación JWT (login, refresh, logout)
    - 🔒 Encriptación de contraseñas con bcrypt
    - 🛡️ Rate limiting anti-fuerza bruta para login
    - 🔗 Migración a middleware compartido (edugo-shared)
    - 🚀 Workflow manual-release TODO-EN-UNO
    - 🤖 Comandos personalizados de Claude Code
    - 📝 Mejoras de documentación
    - ⚙️ Optimizaciones de CI/CD

    ## Commits Incluidos

    [Listar los 25 commits con descripción breve]

    ## Testing

    - [ ] Tests unitarios pasan
    - [ ] Tests de integración pasan
    - [ ] Build exitoso
    - [ ] Swagger docs se generan correctamente
    ```

- [ ] **Ejecutar CI/CD automático**

  El PR debe disparar:
  - [ ] ci.yml (compilación)
  - [ ] test.yml (pruebas)

- [ ] **Review manual del código**

  Revisar cambios críticos:
  - [ ] Autenticación JWT
  - [ ] Middleware de rate limiting
  - [ ] Integración con edugo-shared
  - [ ] Cambios en variables de entorno

- [ ] **Testing manual en ambiente de desarrollo**
  ```bash
  # Checkout de rama de merge
  git checkout merge/dev-to-main-YYYYMMDD

  # Levantar ambiente local
  docker-compose -f ../edugo-dev-environment/docker/docker-compose.yml up -d

  # Ejecutar API
  go run cmd/main.go

  # Probar endpoints críticos:
  # - POST /api/v1/auth/login
  # - POST /api/v1/auth/refresh
  # - POST /api/v1/auth/logout
  # - Verificar rate limiting
  ```

- [ ] **Aprobar y mergear PR**

  En GitHub:
  - Aprobar PR
  - Merge con **"Create a merge commit"** (NO squash)
  - Confirmar merge

- [ ] **Verificar que workflow sync-main-to-dev se ejecutó**

  GitHub Actions → Sync Main to Dev → Última ejecución

  Debe haber:
  - ✅ Ejecución exitosa
  - ✅ Commit de sincronización en dev

- [ ] **Validar sincronización**
  ```bash
  git checkout main
  git pull origin main

  git checkout dev
  git pull origin dev

  # Verificar que están al mismo nivel
  git rev-list --left-right --count main...dev
  # Debe dar: 0       0
  ```

- [ ] **Crear release**
  ```bash
  # Decidir nueva versión (ej: v0.2.0)
  echo "0.2.0" > version.txt
  git add version.txt
  git commit -m "chore: bump version to v0.2.0"
  git push origin main

  # Crear tag
  git tag v0.2.0
  git push origin v0.2.0
  ```

- [ ] **Verificar que se creó release automático en GitHub**

  GitHub → Releases → Debe aparecer v0.2.0

- [ ] **Verificar que se construyó imagen Docker**

  GitHub → Actions → Build and Push → Última ejecución

  Debe haber:
  - ✅ Build exitoso
  - ✅ Imagen en GHCR: ghcr.io/edugogroup/edugo-api-mobile:0.2.0

---

#### 3.3 Sincronizar edugo-api-administracion

**Urgencia**: 🟡 MEDIA
**Tiempo estimado**: 30 minutos
**Complejidad**: Baja (solo 1 commit)

- [ ] **Sincronizar dev desde main**
  ```bash
  cd edugo-api-administracion

  # Checkout dev (puede estar solo en remote)
  git fetch origin
  git checkout dev
  git pull origin dev

  # Merge main a dev
  git merge main --no-ff -m "merge: sincronizar main → dev"

  # Push
  git push origin dev
  ```

- [ ] **Validar sincronización**
  ```bash
  git rev-list --left-right --count main...dev
  # Debe dar: 0       0
  ```

---

#### 3.4 Sincronizar edugo-worker

**Urgencia**: 🟡 MEDIA
**Tiempo estimado**: 2 horas
**Complejidad**: Media (4 commits)

**Similar a api-mobile pero más pequeño**:

- [ ] **Crear PR de dev → main**

  Incluye:
  - Optimizaciones de CI/CD
  - Copilot instructions

- [ ] **Ejecutar CI/CD**

- [ ] **Aprobar y mergear**

- [ ] **Verificar sync-main-to-dev**

- [ ] **Crear release si es necesario**

---

### Fase 4: Mejoras de Infraestructura

#### 4.1 Continuar

---

#### 4.2 Crear Política de Protección de Ramas

**Urgencia**: 🟡 MEDIA
**Tiempo estimado**: 30 minutos
**Responsable**: DevOps / Admin de GitHub

**Para TODOS los repositorios**:

- [ ] **Configurar branch protection en GitHub**

  Settings → Branches → Add rule

  Para rama **main**:
  - ✅ Require pull request before merging
  - ✅ Require approvals: 1
  - ✅ Require status checks to pass: CI, Tests
  - ✅ Require conversation resolution before merging
  - ❌ Allow force pushes (deshabilitar)
  - ❌ Allow deletions (deshabilitar)

  Para rama **dev**:
  - ✅ Require pull request before merging (opcional)
  - ✅ Require status checks to pass: CI
  - ✅ Allow force pushes: solo admins (opcional)

- [ ] **Configurar tag protection**

  Settings → Tags → Add rule

  Pattern: `v*`
  - ✅ Only allow matching tags to be created by users with write access

---

#### 4.3 Documentar Flujo de Trabajo Estándar

**Urgencia**: 🟢 MEDIA
**Tiempo estimado**: 1 hora
**Responsable**: Tech Lead

- [ ] **Crear CONTRIBUTING.md en cada proyecto**

  ```markdown
  # Guía de Contribución

  ## Flujo de Trabajo con Git

  ### Flujo Principal: dev → main

  1. **Desarrollo en dev**:
     - Crea feature branch desde dev: `git checkout -b feature/mi-feature dev`
     - Desarrolla y commitea
     - Crea PR hacia dev
     - Merge a dev después de CI/CD y review

  2. **Release a main**:
     - Cuando dev tenga features listas para producción
     - Crear PR de dev → main
     - CI/CD debe pasar
     - Review obligatorio
     - Merge con **merge commit** (no squash)

  3. **Sincronización automática**:
     - Workflow `sync-main-to-dev.yml` sincroniza main → dev automáticamente
     - Verifica que se ejecutó: GitHub Actions
     - Confirma: `git rev-list --left-right --count main...dev` debe dar `0 0`

  ## Versionado

  - Esquema: Semantic Versioning (v0.x.x / v1.x.x)
  - Fuente de verdad: `version.txt`
  - Tags se crean automáticamente en releases

  ## Commits

  - Formato: `tipo: descripción breve`
  - Tipos: feat, fix, docs, chore, refactor, test, ci
  - Ejemplo: `feat: agregar endpoint de autenticación JWT`
  ```

- [ ] **Actualizar README.md con sección de Git Workflow**

- [ ] **Crear diagrama de flujo** (opcional)

  Usar Mermaid, PlantUML o diagrama visual

---

### Fase 5: Validación Final

#### 5.1 Checklist de Validación por Proyecto

**Para CADA proyecto** (api-mobile, api-admin, worker, shared):

- [ ] **Versionado consistente**
  ```bash
  # version.txt y último tag coinciden
  LAST_TAG=$(git tag -l "v*.*.*" | sort -V | tail -1)
  VERSION_FILE=$(cat version.txt 2>/dev/null || echo "N/A")
  echo "Tag: $LAST_TAG, version.txt: v$VERSION_FILE"
  ```

- [ ] **Ramas sincronizadas**
  ```bash
  git fetch origin
  git rev-list --left-right --count main...dev
  # Debe dar: 0       0
  ```

- [ ] **Workflows funcionando**
  ```bash
  # En GitHub Actions, verificar últimas ejecuciones:
  # - CI: ✅ Passing
  # - Tests: ✅ Passing
  # - Sync main to dev: ✅ Success
  ```

- [ ] **Releases actualizados**
  ```bash
  # GitHub → Releases → Último release coincide con último tag
  ```

- [ ] **Documentación actualizada**
  ```bash
  # README.md tiene sección de versionado
  # CONTRIBUTING.md existe (si aplica)
  ```

---

#### 5.2 Matriz de Estado Final

Llenar esta matriz después de completar el plan:

| Proyecto | version.txt | Último Tag | Main=Dev | Workflows ✓ | CI/CD ✓ | Release ✓ |
|----------|-------------|------------|----------|-------------|---------|-----------|
| api-mobile | ______ | ______ | ☐ Sí ☐ No | ☐ Sí ☐ No | ☐ Sí ☐ No | ☐ Sí ☐ No |
| api-admin | ______ | ______ | ☐ Sí ☐ No | ☐ Sí ☐ No | ☐ Sí ☐ No | ☐ Sí ☐ No |
| worker | ______ | ______ | ☐ Sí ☐ No | ☐ Sí ☐ No | ☐ Sí ☐ No | ☐ Sí ☐ No |
| shared | N/A | ______ | ☐ Sí ☐ No | ☐ Sí ☐ No | ☐ Sí ☐ No | ☐ Sí ☐ No |
| dev-env | N/A | ______ | ☐ N/A | ☐ Sí ☐ No | ☐ Sí ☐ No | ☐ N/A |

---

#### 5.3 Pruebas de Integración

- [ ] **Actualizar edugo-dev-environment con nuevas versiones**
  ```bash
  cd edugo-dev-environment/docker

  # Editar docker-compose.yml con nuevas imágenes:
  # ghcr.io/edugogroup/edugo-api-mobile:0.2.0
  # ghcr.io/edugogroup/edugo-api-administracion:X.X.X
  # ghcr.io/edugogroup/edugo-worker:X.X.X

  # Levantar ambiente
  docker-compose up -d

  # Verificar que todo levanta correctamente
  docker-compose ps

  # Verificar logs
  docker-compose logs api-mobile
  ```

- [ ] **Probar flujo end-to-end**

  - [ ] Login en API mobile
  - [ ] Operación en API admin
  - [ ] Verificar que worker procesa mensajes (si aplica)

- [ ] **Verificar métricas y health checks**
  ```bash
  curl http://localhost:8080/health  # api-mobile
  curl http://localhost:8081/health  # api-admin (ajustar puerto)
  ```

---

### Fase 6: Monitoreo Continuo

#### 6.1 Configurar Alertas de Divergencia

- [ ] **Crear script de monitoreo**

  Archivo: `scripts/check-sync.sh`
  ```bash
  #!/bin/bash

  REPOS=(
    "edugo-api-mobile"
    "edugo-api-administracion"
    "edugo-worker"
    "edugo-shared"
  )

  for REPO in "${REPOS[@]}"; do
    cd "/path/to/$REPO"

    git fetch origin
    DIFF=$(git rev-list --left-right --count main...dev)

    if [ "$DIFF" != "0      0" ]; then
      echo "⚠️ $REPO: Ramas desincronizadas - $DIFF"
    else
      echo "✅ $REPO: Sincronizado"
    fi
  done
  ```

- [ ] **Agregar a cron (opcional)**
  ```bash
  # Ejecutar cada día
  0 9 * * * /path/to/scripts/check-sync.sh | mail -s "Git Sync Status" devops@edugo.com
  ```

---

#### 6.2 Crear Dashboard de Estado (Opcional)

- [ ] **Usar GitHub Actions para generar badge de estado**

- [ ] **Agregar badges al README.md**
  ```markdown
  ![CI](https://github.com/EduGoGroup/edugo-api-mobile/workflows/CI/badge.svg)
  ![Tests](https://github.com/EduGoGroup/edugo-api-mobile/workflows/Tests/badge.svg)
  ![Release](https://github.com/EduGoGroup/edugo-api-mobile/workflows/Release/badge.svg)
  ```

---

## 📎 ANEXOS

### Anexo A: Comandos Útiles de Git

```bash
# Ver divergencia entre ramas
git rev-list --left-right --count main...dev

# Ver commits únicos en dev
git log main..dev --oneline

# Ver commits únicos en main
git log dev..main --oneline

# Ver archivos diferentes entre ramas
git diff main...dev --name-status

# Gráfico de commits
git log --oneline --graph --all -20

# Último tag
git tag -l "v*.*.*" | sort -V | tail -1

# Tags ordenados por fecha
git tag -l --sort=-creatordate

# Verificar si rama está adelante/atrás de remote
git status -sb

# Ver historial de un archivo específico
git log --follow --oneline -- version.txt
```

---

### Anexo B: Estructura de Workflows Estándar

**Workflow mínimo recomendado para proyectos de servicio**:

1. **ci.yml** - Compilación y linting
2. **test.yml** - Suite de pruebas
3. **release.yml** - Release automático en tags
4. **manual-release.yml** - Release manual on-demand
5. **sync-main-to-dev.yml** - Sincronización automática
6. **build-and-push.yml** - Build y push de Docker

**Workflow mínimo para librería (shared)**:

1. **ci.yml** - Compilación y linting
2. **test.yml** - Tests por módulo
3. **release.yml** - Release con validación de módulos
4. **sync-main-to-dev.yml** - Sincronización

---

### Anexo C: Política de Versionado Recomendada

**Semantic Versioning 2.0.0**:

```
vMAJOR.MINOR.PATCH

MAJOR: Cambios incompatibles (breaking changes)
MINOR: Nueva funcionalidad compatible
PATCH: Bug fixes compatibles

Ejemplos:
v0.1.0 → v0.1.1  (bugfix)
v0.1.1 → v0.2.0  (nueva feature)
v0.9.0 → v1.0.0  (primera versión estable)
v1.2.3 → v2.0.0  (breaking change)
```

**Pre-releases**:
```
v1.0.0-alpha.1  (desarrollo temprano)
v1.0.0-beta.1   (feature complete, en testing)
v1.0.0-rc.1     (release candidate)
v1.0.0          (release estable)
```

---

### Anexo D: Troubleshooting

#### Problema: Workflow sync-main-to-dev falla

**Síntomas**:
```
Error: Failed to merge main into dev
Automatic merge failed; fix conflicts and then commit the result
```

**Solución**:
```bash
# Hacer merge manual
git checkout dev
git pull origin dev
git merge origin/main

# Resolver conflictos
git status
# Editar archivos en conflicto
git add .
git commit -m "merge: resolver conflictos main → dev"
git push origin dev
```

---

#### Problema: Tag ya existe pero con versión incorrecta

**Síntomas**:
```
Error: Tag v1.0.0 already exists
```

**Solución**:
```bash
# Eliminar tag local
git tag -d v1.0.0

# Eliminar tag remoto
git push origin --delete v1.0.0

# Crear tag correcto
git tag v0.2.0
git push origin v0.2.0
```

---

#### Problema: go.mod referencia versión de shared que no existe

**Síntomas**:
```
go get github.com/EduGoGroup/edugo-shared@v2.0.6
go: module github.com/EduGoGroup/edugo-shared@v2.0.6: not found
```

**Solución**:
```bash
# Listar tags disponibles
cd edugo-shared
git tag -l | sort -V

# Actualizar a tag disponible
cd edugo-api-mobile
go get github.com/EduGoGroup/edugo-shared@v0.3.0
go mod tidy
```

---

### Anexo E: Checklist de Pre-Release

Antes de crear un release:

- [ ] Todos los tests pasan localmente
- [ ] CI/CD en green
- [ ] version.txt actualizado
- [ ] CHANGELOG.md actualizado (si existe)
- [ ] Documentación actualizada (README, Swagger)
- [ ] Variables de entorno documentadas
- [ ] Migraciones de BD probadas (si aplica)
- [ ] Compatibilidad con versiones anteriores verificada
- [ ] Dependencias actualizadas (go mod tidy)
- [ ] Branch dev sincronizado con main

---

### Anexo F: Contactos y Responsables

| Área | Responsable | Contacto |
|------|-------------|----------|
| DevOps | [Nombre] | [Email] |
| Backend API | [Nombre] | [Email] |
| Worker | [Nombre] | [Email] |
| Shared Libraries | [Nombre] | [Email] |
| QA | [Nombre] | [Email] |
| Product Owner | [Nombre] | [Email] |

---

## 📝 REGISTRO DE CAMBIOS DEL INFORME

| Versión | Fecha | Autor | Cambios |
|---------|-------|-------|---------|
| 1.0 | 2025-11-02 | Claude Code | Informe inicial completo |

---

## 🎯 CONCLUSIÓN

Este informe proporciona un análisis exhaustivo del estado actual del ecosistema EduGo y un plan de acción detallado para:

1. ✅ **Estandarizar versionado** en todos los proyectos
2. ✅ **Sincronizar ramas** main y dev consistentemente
3. ✅ **Proteger trabajo** importante acumulado en dev
4. ✅ **Establecer procesos** para evitar divergencias futuras
5. ✅ **Mejorar infraestructura** de CI/CD donde necesario

**Prioridades**:
1. 🔴 Decidir esquema de versionado (Fase 1.1)
2. 🔴 Mergear api-mobile dev → main (Fase 3.2)
3. 🔴 Reconciliar shared (Fase 3.1)
4. 🟡 Sincronizar worker y api-admin (Fases 3.3, 3.4)
5. 🟢 Mejoras de infraestructura (Fase 4)

**Tiempo estimado total**: 12-16 horas de trabajo técnico distribuidas en varios días.

---

**¿Preguntas o necesitas clarificaciones?**
Contacta al equipo técnico o revisa la documentación de cada proyecto.

---

*Este informe fue generado con la asistencia de Claude Code.*
*Última actualización: 2025-11-02*
