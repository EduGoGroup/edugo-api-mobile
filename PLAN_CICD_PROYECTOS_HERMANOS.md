# 🚀 Plan de Trabajo: Replicar CI/CD y Copilot en Proyectos Hermanos

**Fecha de creación:** 2025-11-01
**Proyecto origen:** edugo-api-mobile
**Commit de referencia:** f0f8ba5 (feat: agregar Copilot custom instructions y optimizar workflows CI/CD)

---

## 🎯 Objetivo General

Replicar las mejoras de CI/CD y configuración de GitHub Copilot implementadas en `edugo-api-mobile` a los proyectos hermanos del ecosistema EduGo, adaptando cada configuración según las necesidades específicas de cada proyecto.

---

## 📦 Proyectos a Actualizar

| # | Proyecto | Ruta | Tipo | Prioridad | Docker |
|---|----------|------|------|-----------|--------|
| 1 | **edugo-api-administracion** | `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion` | API REST (Go) | Alta | ✅ Sí |
| 2 | **edugo-worker** | `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-worker` | Worker (Go) | Alta | ✅ Sí |
| 3 | **edugo-shared** | `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared` | Librería (Go) | Alta | ❌ No (solo releases) |
| 4 | **edugo-dev-environment** | `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-dev-environment` | Docker Compose | Baja | ❌ No (orquestador) |

---

## 🧩 Componentes a Replicar

### ✅ Componentes Obligatorios (Todos los Proyectos)

- [ ] `.github/copilot-instructions.md` - Instrucciones personalizadas de Copilot en español
- [ ] `.github/workflows/README.md` - Documentación de estrategia CI/CD
- [ ] Optimización de workflows existentes
- [ ] Tabla de estrategia por branch

### ⚙️ Workflows Según Tipo de Proyecto

#### **APIs (api-mobile, api-administracion)**
- [ ] `ci.yml` - Pipeline de integración continua
- [ ] `test.yml` - Tests con cobertura
- [ ] `manual-release.yml` - ⭐ **Workflow TODO-EN-UNO** manual para crear releases (reemplaza auto-version)
- [ ] `docker-only.yml` - Build manual de Docker (opcional, manual-release ya incluye Docker)
- [ ] `release.yml` - Disparado automáticamente por tags (puede ejecutarse manual también)
- [ ] `sync-main-to-dev.yml` - Sincronización de branches
- [ ] ❌ **NO auto-version.yml** - Eliminado (inestable, reemplazado por manual-release)

#### **Worker (edugo-worker)**
- [ ] `ci.yml` - Pipeline de integración continua
- [ ] `test.yml` - Tests con cobertura
- [ ] `manual-release.yml` - ⭐ **Workflow TODO-EN-UNO** manual para crear releases
- [ ] `docker-only.yml` - Build manual de Docker (opcional)
- [ ] `release.yml` - Disparado por tags (opcional)
- [ ] ❌ **NO auto-version.yml** - Versionado manual on-demand
- [ ] ❌ **NO sync-main-to-dev** (flujo más simple para workers)

#### **Librería Compartida (edugo-shared)**
- [ ] `ci.yml` - Tests y validación
- [ ] `test.yml` - Tests con cobertura
- [ ] `release.yml` - **CRÍTICO**: Crear GitHub Release con tag
- [ ] ❌ **NO Docker workflows** (no genera imágenes)
- [ ] ❌ **NO auto-version** (versionado manual con tags)

#### **Dev Environment (edugo-dev-environment)**
- [ ] `validation.yml` - Validar docker-compose.yml
- [ ] ❌ **NO workflows complejos** (solo validación)
- [ ] Documentación de uso

---

## 📋 Plan de Ejecución por Proyecto

---

## 🔷 FASE 1: edugo-api-administracion (API Admin)

**Branch de trabajo:** `feature/cicd-copilot-setup`
**PR destino:** `dev`

### 📊 Análisis Previo

- [ ] **Paso 1.1:** Explorar estructura del proyecto
  ```bash
  cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion
  tree -L 2 -I 'node_modules|vendor'
  git log --oneline -10
  git branch -a
  ```

- [ ] **Paso 1.2:** Verificar workflows existentes
  ```bash
  ls -la .github/workflows/
  cat .github/workflows/*.yml  # Revisar cada uno
  ```

- [ ] **Paso 1.3:** Identificar diferencias con api-mobile
  - Versión de Go
  - Dependencias específicas
  - Estructura de directorios
  - Naming conventions

### 🛠️ Implementación

- [ ] **Paso 1.4:** Crear branch de trabajo
  ```bash
  git checkout -b feature/cicd-copilot-setup
  ```

- [ ] **Paso 1.5:** Crear `.github/copilot-instructions.md`
  - Copiar desde api-mobile
  - Adaptar arquitectura específica
  - Actualizar TODOs y deuda técnica
  - Mantener configuración de español

- [ ] **Paso 1.6:** Crear/Actualizar workflows
  - [ ] `ci.yml` - Adaptar versión de Go y dependencias
  - [ ] `test.yml` - Configurar cobertura
  - [ ] `manual-release.yml` - ⭐ **Copiar desde api-mobile** (TODO-EN-UNO: version + Docker + release)
  - [ ] `docker-only.yml` - Adaptar nombre de imagen (opcional, manual-release ya incluye Docker)
  - [ ] `release.yml` - Adaptar tags y nombres + agregar workflow_dispatch
  - [ ] `sync-main-to-dev.yml` - Mantener igual
  - [ ] ❌ **NO auto-version.yml** - No copiar (fue eliminado de api-mobile)

- [ ] **Paso 1.7:** Actualizar `.github/workflows/README.md`
  - Tabla de estrategia por branch
  - Sección de Copilot
  - Documentación específica del proyecto

- [ ] **Paso 1.8:** Verificar configuración Docker
  ```bash
  # Verificar que existe Dockerfile
  cat Dockerfile

  # Verificar GHCR registry en workflows
  grep -r "ghcr.io" .github/workflows/
  ```

### ✅ Validación y PR

- [ ] **Paso 1.9:** Commit de cambios
  ```bash
  git add .github/
  git commit -m "feat: agregar Copilot instructions y optimizar workflows CI/CD

  Adaptado desde edugo-api-mobile (commit f0f8ba5)
  - Copilot custom instructions en español
  - Workflows optimizados por branch
  - Documentación de estrategia CI/CD
  "
  ```

- [ ] **Paso 1.10:** Push y crear PR
  ```bash
  git push origin feature/cicd-copilot-setup
  gh pr create --base dev --head feature/cicd-copilot-setup \
    --title "feat: Copilot instructions y optimización CI/CD" \
    --body "Ver PLAN_CICD_PROYECTOS_HERMANOS.md para detalles"
  ```

- [ ] **Paso 1.11:** Verificar ejecución de workflows
  ```bash
  # Esperar a que se ejecuten ci.yml y test.yml
  gh run watch

  # Verificar que Copilot haga review en español
  gh pr view --web
  ```

- [ ] **Paso 1.12:** Aprobar y mergear PR
  ```bash
  # Revisar checks
  gh pr checks

  # Mergear cuando esté listo
  gh pr merge --squash
  ```

### 🎯 Checklist de Validación

- [ ] ✅ Workflows se ejecutan en PR a dev
- [ ] ✅ Copilot comenta en español
- [ ] ✅ CI pasa correctamente
- [ ] ✅ Tests con cobertura generan reporte
- [ ] ✅ Documentación es clara
- [ ] ❌ NO se ejecutan workflows en feature branches

---

## 🟦 FASE 2: edugo-worker (Worker)

**Branch de trabajo:** `feature/cicd-copilot-setup`
**PR destino:** `dev`

### 📊 Análisis Previo

- [ ] **Paso 2.1:** Explorar estructura del proyecto
  ```bash
  cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-worker
  tree -L 2 -I 'node_modules|vendor'
  git log --oneline -10
  git branch -a
  ```

- [ ] **Paso 2.2:** Identificar características únicas
  - ¿Es un worker con cron jobs?
  - ¿Es un consumer de RabbitMQ?
  - ¿Qué procesos ejecuta?
  - ¿Cómo se diferencia de las APIs?

- [ ] **Paso 2.3:** Verificar dependencias
  ```bash
  cat go.mod | grep edugo-shared
  cat go.mod | grep rabbitmq
  ```

### 🛠️ Implementación

- [ ] **Paso 2.4:** Crear branch de trabajo
  ```bash
  git checkout -b feature/cicd-copilot-setup
  ```

- [ ] **Paso 2.5:** Crear `.github/copilot-instructions.md`
  - Copiar base desde api-mobile
  - **ADAPTAR:** Indicar que es un Worker, no una API
  - **AGREGAR:** Patrones de workers (cron, consumers, processors)
  - **AGREGAR:** Manejo de jobs async y reintentos
  - **AGREGAR:** Logging específico de workers
  - Mantener español

- [ ] **Paso 2.6:** Crear workflows específicos de Worker
  - [ ] `ci.yml` - Tests y validación (sin handlers HTTP)
  - [ ] `test.yml` - Cobertura enfocada en processors
  - [ ] `manual-release.yml` - ⭐ **Copiar desde api-mobile** (TODO-EN-UNO)
  - [ ] `docker-only.yml` - Build manual (opcional)
  - [ ] `release.yml` - Con workflow_dispatch (opcional)
  - [ ] ❌ **NO auto-version.yml** - Versionado manual on-demand
  - [ ] ❌ **NO sync-main-to-dev.yml** (flujo más simple)

- [ ] **Paso 2.7:** Adaptar documentación
  - Explicar que es un worker
  - Documentar patrón de ejecución
  - Estrategia de deployment diferente a APIs

### ✅ Validación y PR

- [ ] **Paso 2.8:** Commit y push
  ```bash
  git add .github/
  git commit -m "feat: agregar Copilot instructions y workflows CI/CD para Worker"
  git push origin feature/cicd-copilot-setup
  ```

- [ ] **Paso 2.9:** Crear PR a dev
  ```bash
  gh pr create --base dev --head feature/cicd-copilot-setup \
    --title "feat: Copilot instructions y CI/CD para Worker" \
    --body "Adaptación específica para proyecto tipo Worker"
  ```

- [ ] **Paso 2.10:** Validar ejecución y mergear
  ```bash
  gh run watch
  gh pr merge --squash
  ```

### 🎯 Checklist de Validación

- [ ] ✅ Copilot entiende que es un Worker (no API)
- [ ] ✅ Tests ejecutan jobs correctamente
- [ ] ✅ Docker build genera imagen funcional
- [ ] ✅ Documentación refleja naturaleza de worker

---

## 🟩 FASE 3: edugo-shared (Librería Compartida)

**Branch de trabajo:** `feature/cicd-copilot-setup`
**PR destino:** `dev`

### 📊 Análisis Previo

- [ ] **Paso 3.1:** Explorar estructura
  ```bash
  cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared
  tree -L 2
  git log --oneline -10
  git tag -l | sort -V | tail -10  # Ver últimos tags
  ```

- [ ] **Paso 3.2:** Verificar estructura de paquetes
  ```bash
  ls -la */
  # Verificar: auth/, logger/, common/errors/, etc.
  ```

- [ ] **Paso 3.3:** Revisar flujo de versionado actual
  ```bash
  git log --tags --simplify-by-decoration --pretty="format:%ai %d"
  ```

### 🛠️ Implementación

- [ ] **Paso 3.4:** Crear branch de trabajo
  ```bash
  git checkout -b feature/cicd-copilot-setup
  ```

- [ ] **Paso 3.5:** Crear `.github/copilot-instructions.md` ESPECIAL
  - **FOCO:** Librería compartida, no aplicación
  - **IMPORTANTE:** Retrocompatibilidad (breaking changes)
  - **REGLAS:**
    - Semantic Versioning estricto
    - Documentar breaking changes
    - Tests exhaustivos (otros proyectos dependen)
    - Ejemplos de uso de cada paquete
  - **PAQUETES:** Documentar auth, logger, errors, middleware
  - Mantener español

- [ ] **Paso 3.6:** Crear workflows específicos de Librería
  - [ ] `ci.yml` - Tests exhaustivos (sin Docker)
  - [ ] `test.yml` - Cobertura alta (>80%)
  - [ ] `release.yml` - **CRÍTICO**:
    - Crear GitHub Release
    - Generar changelog
    - Publicar tag
    - **NO construir Docker**
  - [ ] ❌ **NO docker workflows**
  - [ ] ❌ **NO auto-version** (tags manuales)

- [ ] **Paso 3.7:** Crear workflow de validación de breaking changes
  ```yaml
  # breaking-changes.yml
  # Detectar cambios en interfaces públicas
  # Alertar si hay breaking changes sin bump de major version
  ```

- [ ] **Paso 3.8:** Documentar flujo de actualización
  ```markdown
  ## Cómo Crear un Release de edugo-shared

  1. Hacer cambios en feature branch
  2. Crear PR a dev
  3. Mergear a dev → main
  4. Crear tag: git tag v0.3.0
  5. Push tag: git push origin v0.3.0
  6. Release workflow genera GitHub Release automático
  7. Otros proyectos actualizan: go get github.com/EduGoGroup/edugo-shared@v0.3.0
  ```

### ✅ Validación y PR

- [ ] **Paso 3.9:** Commit y push
  ```bash
  git add .github/
  git commit -m "feat: agregar Copilot instructions y workflows CI/CD para librería compartida"
  git push origin feature/cicd-copilot-setup
  ```

- [ ] **Paso 3.10:** Crear PR a dev
  ```bash
  gh pr create --base dev --head feature/cicd-copilot-setup \
    --title "feat: Copilot instructions y CI/CD para librería compartida" \
    --body "Configuración específica para módulo Go compartido"
  ```

- [ ] **Paso 3.11:** Validar y mergear
  ```bash
  gh run watch
  gh pr merge --squash
  ```

- [ ] **Paso 3.12:** Probar flujo de release
  ```bash
  # Después del merge, crear tag de prueba
  git checkout main
  git pull
  git tag v0.0.0-test-20251101
  git push origin v0.0.0-test-20251101

  # Verificar que release.yml crea GitHub Release
  gh release view v0.0.0-test-20251101
  ```

### 🎯 Checklist de Validación

- [ ] ✅ Copilot enfocado en retrocompatibilidad
- [ ] ✅ Tests de alta cobertura pasan
- [ ] ✅ Release workflow crea GitHub Release
- [ ] ✅ Tag genera release visible en GitHub
- [ ] ✅ Otros proyectos pueden hacer go get del tag
- [ ] ❌ NO intenta construir Docker

---

## 🟪 FASE 4: edugo-dev-environment (Docker Compose)

**Branch de trabajo:** `feature/cicd-validation`
**PR destino:** `dev`

### 📊 Análisis Previo

- [ ] **Paso 4.1:** Explorar estructura
  ```bash
  cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-dev-environment
  ls -la
  cat docker-compose.yml
  ```

- [ ] **Paso 4.2:** Identificar servicios orquestados
  ```bash
  # Ver qué imágenes usa
  grep "image:" docker-compose.yml

  # Verificar si usa imágenes de GHCR
  grep "ghcr.io" docker-compose.yml
  ```

### 🛠️ Implementación

- [ ] **Paso 4.3:** Crear branch de trabajo
  ```bash
  git checkout -b feature/cicd-validation
  ```

- [ ] **Paso 4.4:** Crear `.github/copilot-instructions.md` MINIMALISTA
  - **FOCO:** Orquestación de servicios
  - **REGLAS:**
    - Validar sintaxis de docker-compose.yml
    - Documentar servicios y sus puertos
    - Variables de entorno requeridas
  - Mantener español

- [ ] **Paso 4.5:** Crear workflow simple de validación
  ```yaml
  # validation.yml
  name: Validate Docker Compose

  on:
    pull_request:
      branches: [main, dev]
    push:
      branches: [main]

  jobs:
    validate:
      runs-on: ubuntu-latest
      steps:
        - uses: actions/checkout@v4
        - name: Validate docker-compose.yml
          run: docker-compose config
  ```

- [ ] **Paso 4.6:** Actualizar README.md
  ```markdown
  ## Cómo usar este entorno

  1. Descargar imágenes: docker-compose pull
  2. Levantar servicios: docker-compose up -d
  3. Ver logs: docker-compose logs -f
  4. Detener: docker-compose down
  ```

- [ ] **Paso 4.7:** NO agregar workflows complejos
  - ❌ NO auto-version
  - ❌ NO Docker build (solo orquesta)
  - ❌ NO release workflow
  - ✅ Solo validación de sintaxis

### ✅ Validación y PR

- [ ] **Paso 4.8:** Commit y push
  ```bash
  git add .github/ README.md
  git commit -m "feat: agregar validación CI/CD para docker-compose"
  git push origin feature/cicd-validation
  ```

- [ ] **Paso 4.9:** Crear PR a dev
  ```bash
  gh pr create --base dev --head feature/cicd-validation \
    --title "feat: Validación CI/CD para entorno de desarrollo" \
    --body "Validación simple de docker-compose.yml"
  ```

- [ ] **Paso 4.10:** Validar y mergear
  ```bash
  gh run watch
  gh pr merge --squash
  ```

### 🎯 Checklist de Validación

- [ ] ✅ Workflow valida sintaxis de docker-compose.yml
- [ ] ✅ README actualizado con instrucciones
- [ ] ✅ Copilot ayuda con configuración de servicios
- [ ] ❌ NO workflows innecesarios

---

## 🔄 FASE 5: Sincronización y Validación Final

### 🧪 Validación Integral

- [ ] **Paso 5.1:** Verificar todos los PRs mergeados
  ```bash
  # edugo-api-administracion
  cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion
  git log --oneline -1  # Verificar commit de CI/CD

  # edugo-worker
  cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-worker
  git log --oneline -1

  # edugo-shared
  cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared
  git log --oneline -1

  # edugo-dev-environment
  cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-dev-environment
  git log --oneline -1
  ```

- [ ] **Paso 5.2:** Probar flujo completo end-to-end

  #### 5.2.1: Cambio en edugo-shared
  ```bash
  cd edugo-shared

  # Hacer cambio menor
  echo "// Test comment" >> logger/logger.go
  git add .
  git commit -m "test: validar flujo CI/CD"
  git push

  # Crear tag
  git tag v0.0.0-cicd-test
  git push origin v0.0.0-cicd-test

  # Verificar que se crea Release
  gh release view v0.0.0-cicd-test
  ```

  #### 5.2.2: Actualizar dependencia en api-mobile
  ```bash
  cd edugo-api-mobile

  # Actualizar a nuevo tag
  go get github.com/EduGoGroup/edugo-shared@v0.0.0-cicd-test
  go mod tidy

  # Crear PR
  git checkout -b test/update-shared
  git add go.mod go.sum
  git commit -m "test: actualizar edugo-shared a tag de prueba"
  git push origin test/update-shared
  gh pr create --base dev

  # Verificar que CI/CD se ejecuta
  gh run watch
  ```

  #### 5.2.3: Build de imágenes Docker
  ```bash
  # api-mobile
  cd edugo-api-mobile
  gh workflow run docker-only.yml -f tag=test-cicd

  # api-administracion
  cd edugo-api-administracion
  gh workflow run docker-only.yml -f tag=test-cicd

  # worker
  cd edugo-worker
  gh workflow run docker-only.yml -f tag=test-cicd

  # Verificar que se publican en GHCR
  # https://github.com/orgs/EduGoGroup/packages
  ```

  #### 5.2.4: Actualizar dev-environment
  ```bash
  cd edugo-dev-environment

  # Actualizar docker-compose.yml con nuevos tags
  sed -i 's/:latest/:test-cicd/g' docker-compose.yml

  git add docker-compose.yml
  git commit -m "test: usar tags de test CI/CD"
  git push

  # Verificar que validation workflow pasa
  gh run list --limit 1
  ```

- [ ] **Paso 5.3:** Validar Copilot en todos los proyectos
  - [ ] Crear PR en cada proyecto
  - [ ] Verificar que Copilot comenta en español
  - [ ] Verificar sugerencias contextuales
  - [ ] Verificar detección de anti-patrones

### 📊 Dashboard de Estado Final

| Proyecto | Branch | Copilot | CI/CD | Docker | Release | Estado |
|----------|--------|---------|-------|--------|---------|--------|
| **api-mobile** | ✅ main | ✅ | ✅ | ✅ | ✅ | ✅ Completado |
| **api-administracion** | ⏳ dev | ⏳ | ⏳ | ⏳ | ⏳ | ⏳ En progreso |
| **worker** | ⏳ dev | ⏳ | ⏳ | ⏳ | ⏳ | ⏳ Pendiente |
| **shared** | ⏳ dev | ⏳ | ⏳ | N/A | ⏳ | ⏳ Pendiente |
| **dev-environment** | ⏳ dev | ⏳ | ⏳ | N/A | N/A | ⏳ Pendiente |

---

## 📚 Documentación de Referencia

### Archivos de Origen (edugo-api-mobile)

- ✅ `.github/copilot-instructions.md` (621 líneas)
- ✅ `.github/workflows/ci.yml`
- ✅ `.github/workflows/test.yml`
- ✅ `.github/workflows/manual-release.yml` ⭐ **NUEVO - TODO-EN-UNO**
- ✅ `.github/workflows/docker-only.yml` (opcional)
- ✅ `.github/workflows/release.yml` (con workflow_dispatch)
- ✅ `.github/workflows/sync-main-to-dev.yml`
- ✅ `.github/workflows/README.md`
- ❌ ~~`.github/workflows/auto-version.yml`~~ (eliminado - inestable)

### Adaptaciones por Tipo de Proyecto

#### APIs (api-mobile, api-administracion)
```
Workflows completos + Docker + Manual Release (TODO-EN-UNO) + Sync
- auto-version.yml eliminado (inestable)
+ manual-release.yml (control total, on-demand)
```

#### Worker (edugo-worker)
```
Workflows completos + Docker + Manual Release (TODO-EN-UNO) - Sync
- auto-version.yml NO incluir
+ manual-release.yml (control total, on-demand)
+ Lógica específica de workers en copilot-instructions.md
```

#### Librería (edugo-shared)
```
Workflows básicos + Release con tags - Docker
+ Énfasis en retrocompatibilidad y semantic versioning
```

#### Orquestador (edugo-dev-environment)
```
Solo validación + Documentación
+ Copilot ayuda con docker-compose.yml
```

---

## 🎯 Métricas de Éxito

### Por Proyecto

- [ ] ✅ Copilot responde en español
- [ ] ✅ CI/CD se ejecuta solo en PRs (no en feature branches)
- [ ] ✅ Tests pasan correctamente
- [ ] ✅ Cobertura >70% (APIs y Worker) / >80% (Shared)
- [ ] ✅ Docker images se publican correctamente (donde aplique)
- [ ] ✅ Releases automáticos funcionan (donde aplique)
- [ ] ✅ Documentación clara y actualizada

### Ecosistema Completo

- [ ] ✅ Flujo: Cambio en shared → Tag → Release → Update en APIs
- [ ] ✅ Flujo: PR en API → CI/CD → Docker → Disponible en dev-environment
- [ ] ✅ Flujo: Merge a main → Auto-version → Release → Deployment
- [ ] ✅ Copilot consistente en todos los proyectos
- [ ] ✅ Estrategia CI/CD documentada y entendida

---

## ⚠️ Notas Importantes

### Errores Comunes a Evitar

1. **NO copiar workflows a ciegas** - Cada proyecto necesita adaptación
2. **NO olvidar actualizar nombres de imágenes Docker**
3. **NO usar auto-version en shared** - Versionado manual con tags
4. **NO agregar Docker a shared** - Es una librería, no un servicio
5. **NO sobre-automatizar dev-environment** - Es solo orquestación

### Orden de Ejecución Recomendado

1. ✅ **edugo-shared PRIMERO** - Todos dependen de esta
2. ✅ **api-administracion** - API similar a api-mobile
3. ✅ **worker** - Diferente pero usa shared
4. ✅ **dev-environment** - Orquesta todos los anteriores

### Rollback Plan

Si algo sale mal en un proyecto:
```bash
# Revertir commit
git revert HEAD

# O eliminar branch y workflows
git checkout main
git branch -D feature/cicd-copilot-setup
git push origin --delete feature/cicd-copilot-setup

# Restaurar workflows originales desde backup
```

---

## 📝 Checklist General de Proyecto

Para cada proyecto, verificar:

- [ ] ✅ Branch de trabajo creado
- [ ] ✅ Copilot instructions adaptado
- [ ] ✅ Workflows necesarios creados/actualizados
- [ ] ✅ README de workflows documentado
- [ ] ✅ Commit descriptivo creado
- [ ] ✅ Push realizado
- [ ] ✅ PR creado a dev
- [ ] ✅ CI/CD ejecutándose
- [ ] ✅ Copilot revisando en español
- [ ] ✅ Checks pasando
- [ ] ✅ PR mergeado
- [ ] ✅ Validación post-merge
- [ ] ✅ Documentado en este plan

---

## 🏁 Estado del Plan

**Última actualización:** 2025-11-01 23:20
**Completado:** 2/5 proyectos (40%)
**Nota importante:** auto-version.yml eliminado, usar manual-release.yml en su lugar

| Proyecto | Estado | Fecha Inicio | Fecha Fin | Notas |
|----------|--------|--------------|-----------|-------|
| edugo-api-mobile | ✅ Completado | 2025-11-01 | 2025-11-01 | Workflows optimizados + manual-release.yml TODO-EN-UNO + v0.1.1 |
| edugo-shared | ✅ Completado | 2025-11-01 | 2025-11-01 | CI/CD + versionado v0.3.0 + manual (sin Docker) |
| edugo-api-administracion | ⏸️ Pendiente | - | - | Usar manual-release.yml + v0.x.x |
| edugo-worker | ⏸️ Pendiente | - | - | Usar manual-release.yml + v0.x.x |
| edugo-dev-environment | ⏸️ Pendiente | - | - | - |

---

## 🛠️ Herramientas de Validación Pre-Commit

### ⚠️ IMPORTANTE: Validar Workflows ANTES de Push

Durante la implementación en `edugo-api-mobile` encontramos **errores críticos de sintaxis YAML** que causaron fallos en GitHub Actions. Para evitar esto en los proyectos hermanos, **SIEMPRE** validar workflows localmente.

### 🔧 Instalar actionlint

```bash
# macOS
brew install actionlint

# Linux
wget https://github.com/rhysd/actionlint/releases/latest/download/actionlint_linux_amd64.tar.gz
tar -xzf actionlint_linux_amd64.tar.gz
sudo mv actionlint /usr/local/bin/

# Verificar instalación
actionlint --version
```

### ✅ Flujo de Validación Recomendado

```bash
# 1. Modificar workflow
vim .github/workflows/mi-workflow.yml

# 2. VALIDAR antes de commit
actionlint .github/workflows/mi-workflow.yml

# 3. Si pasa validación → commit
git add .github/workflows/mi-workflow.yml
git commit -m "feat: agregar workflow"

# 4. Push con confianza
git push origin feature/mi-branch
```

### 🚨 Errores Comunes Encontrados y Sus Soluciones

#### **Error 1: Commit Messages Multilinea**

**Problema:**
```yaml
# ❌ INCORRECTO
git commit -m "mensaje línea 1

línea 2"  # Error de parsing YAML
```

**Solución:**
```yaml
# ✅ CORRECTO
git commit -m "mensaje línea 1" -m "" -m "línea 2"
```

#### **Error 2: Backticks en Strings de Bash**

**Problema:**
```yaml
# ❌ INCORRECTO
--body "Este es un \`código\` con backticks"  # Causa command substitution
```

**Solución A (Concatenación):**
```yaml
# ✅ CORRECTO
BODY="Este es un \`código\` con backticks"
--body "$BODY"
```

**Solución B (Archivo temporal):**
```yaml
# ✅ CORRECTO
echo "Este es un \`código\` con backticks" > /tmp/body.txt
--body-file /tmp/body.txt
```

#### **Error 3: Heredocs Dentro de Workflows**

**Problema:**
```yaml
# ❌ INCORRECTO - actionlint falla con heredocs complejos
run: |
  cat <<EOF
  Texto con \`backticks\`
  EOF
```

**Solución:**
```yaml
# ✅ CORRECTO - Usar variables concatenadas
run: |
  TEXT="Texto con \`backticks\`"
  echo "$TEXT"
```

### 📋 Checklist de Validación por Proyecto

Antes de hacer push en cada proyecto hermano:

- [ ] ✅ `actionlint` instalado
- [ ] ✅ Todos los workflows validados localmente
- [ ] ✅ No hay errores `syntax-check`
- [ ] ✅ Warnings de `shellcheck` revisados (opcionales)
- [ ] ✅ Commit messages sin caracteres especiales problemáticos
- [ ] ✅ Backticks escapados correctamente en scripts bash
- [ ] ✅ Heredocs evitados o simplificados

### 🎯 Comando de Validación Rápida

```bash
# Validar todos los workflows de una vez
actionlint .github/workflows/*.yml

# Ver solo errores críticos (ignorar warnings)
actionlint .github/workflows/*.yml 2>&1 | grep "syntax-check" || echo "✅ OK"

# Validar + commit en un solo paso
actionlint .github/workflows/*.yml && \
  git add .github/workflows/*.yml && \
  git commit -m "feat: agregar workflows validados"
```

### 📊 Resultados de Validación en edugo-api-mobile

| Workflow | Errores Encontrados | Solución Aplicada |
|----------|---------------------|-------------------|
| `auto-version.yml` | Commit multilinea (línea 47) | Múltiples flags `-m` |
| `sync-main-to-dev.yml` | Backticks en heredoc (línea 80) | Concatenación de strings |
| `ci.yml` | ✅ Sin errores | N/A |
| `test.yml` | ✅ Sin errores | N/A |
| `docker-only.yml` | ✅ Sin errores | N/A |
| `release.yml` | ⚠️ Warnings shellcheck | Ignorados (no críticos) |

### 🔗 Referencias

- **actionlint GitHub**: https://github.com/rhysd/actionlint
- **Documentación**: https://github.com/rhysd/actionlint/blob/main/docs/usage.md
- **Errores comunes**: Documentados en este plan

---

## ⭐ NUEVO: Workflow Manual Release (TODO-EN-UNO)

### 🎯 Descripción

**manual-release.yml** es un workflow completamente manual que hace TODO en un solo proceso:

1. ✅ Actualiza `version.txt`
2. ✅ Genera entrada en `CHANGELOG.md`
3. ✅ Crea commit de versión en main
4. ✅ Crea y pushea tag
5. ✅ Ejecuta tests completos
6. ✅ Construye imagen Docker multi-platform (amd64/arm64)
7. ✅ Publica imagen en GitHub Container Registry (GHCR)
8. ✅ Crea GitHub Release con changelog

### 🚀 Cómo Usarlo

```bash
# Desde GitHub UI:
1. Ir a: https://github.com/EduGoGroup/edugo-api-mobile/actions/workflows/manual-release.yml
2. Click "Run workflow"
3. Inputs:
   - Branch: main
   - Versión: 0.2.0 (sin 'v')
   - Tipo: minor / patch / major
4. Click "Run workflow"

# El workflow tarda ~20 minutos:
# - 1 min: Preparación (version.txt, CHANGELOG, tag)
# - 2 min: Tests
# - 17 min: Build Docker multi-platform
# - 1 min: GitHub Release
```

### ✅ Ventajas sobre Auto-Version

| Aspecto | auto-version.yml (❌ Eliminado) | manual-release.yml (✅ Nuevo) |
|---------|-------------------------------|-------------------------------|
| **Control** | Automático (impredecible) | Manual (tú decides cuándo) |
| **Confiabilidad** | Inestable (a veces no funciona) | 100% predecible |
| **Visibilidad** | Separado en múltiples workflows | TODO en un solo lugar |
| **Docker** | Depende de release.yml separado | Incluido en el mismo workflow |
| **Duración** | Desconocida | Predecible (~20 min) |
| **Debugging** | Difícil (múltiples workflows) | Fácil (un solo workflow) |

### 📋 Inputs del Workflow

**version** (required):
- Formato: `0.1.0` (sin 'v')
- Validación: Debe ser semver válido
- Verifica que el tag no exista

**bump_type** (required):
- `patch`: 0.1.0 → 0.1.1 (bugfix)
- `minor`: 0.1.0 → 0.2.0 (nueva feature)
- `major`: 0.1.0 → 1.0.0 (breaking change o producción)

### 📊 Outputs

| Componente | Descripción | Ubicación |
|------------|-------------|-----------|
| **Tag Git** | v0.1.0 | GitHub repository tags |
| **Commit** | chore: release v0.1.0 | Branch main |
| **Imagen Docker** | ghcr.io/edugogroup/edugo-api-mobile:v0.1.0 | GitHub Container Registry |
| **GitHub Release** | Release v0.1.0 | GitHub Releases |
| **CHANGELOG** | Entrada [0.1.0] | CHANGELOG.md |

### 🔧 Para Proyectos Hermanos

Al implementar en api-administracion y worker:

1. **Copiar** `.github/workflows/manual-release.yml` desde api-mobile
2. **Adaptar** nombre de imagen Docker (si es diferente)
3. **Mantener** todo lo demás igual
4. **NO copiar** auto-version.yml (fue eliminado)

---

## 📚 Lecciones Aprendidas del Proyecto Origen

### ✅ Lo que Funcionó Bien

1. **Copilot Custom Instructions en Español** - Excelente adopción
2. **Estrategia por Branch** - Elimina falsos positivos
3. **Validación Local con actionlint** - Previene errores
4. **Documentación Detallada** - Facilita replicación
5. **Plan con Checkboxes** - Tracking efectivo
6. **Workflow Manual Release (TODO-EN-UNO)** - Control total, predecible, 100% funcional

### ⚠️ Problemas Encontrados y Soluciones

| Problema | Causa | Solución | Prevención |
|----------|-------|----------|------------|
| Workflows fallando en push | Sintaxis YAML incorrecta | Usar actionlint | Validar antes de push |
| Backticks causan errores | Command substitution | Escapar o concatenar | Evitar en heredocs |
| Commit multilinea | Parsing YAML | Múltiples `-m` flags | Simplificar mensajes |
| Workflows ejecutándose en feature/* | Triggers incorrectos | Filtrar por branch | Documentar triggers |
| auto-version.yml inestable | Timing impredecible, fallos aleatorios | manual-release.yml TODO-EN-UNO | Usar workflows manuales controlados |

### 🎯 Recomendaciones para Proyectos Hermanos

1. **SIEMPRE usar actionlint** antes de push
2. **Copiar workflows validados** desde edugo-api-mobile como base
3. **Adaptar nombres** de imágenes Docker y variables
4. **Testear manualmente** con `workflow_dispatch` primero
5. **Documentar cambios** específicos del proyecto

---

## 🔢 IMPORTANTE: Esquema de Versionado Correcto

### ⚠️ Corrección Aplicada (2025-11-01)

**Problema Detectado**: Todos los proyectos estaban usando versiones v1.x.x y v2.x.x cuando deberían usar v0.x.x (proyectos en desarrollo, sin producción).

**Solución Implementada**: Reseteo a v0.x.x en todos los proyectos.

### Versionado por Proyecto

| Proyecto | Versión Anterior (Incorrecta) | Versión Nueva (Correcta) | Estado |
|----------|-------------------------------|--------------------------|--------|
| **edugo-shared** | v2.0.6 | v0.3.0 | ✅ Corregido |
| **edugo-api-mobile** | v1.0.2 | v0.1.0 | ⏸️ Pendiente |
| **edugo-api-administracion** | TBD | v0.1.0 | ⏸️ Pendiente |
| **edugo-worker** | TBD | v0.1.0 | ⏸️ Pendiente |

### Regla de Versionado para Proyectos en Desarrollo

```
v0.1.0 → Primera versión funcional
v0.2.0 → Nueva feature
v0.x.x → Desarrollo continuo (pueden haber breaking changes)

v1.0.0 → SOLO cuando salga a PRODUCCIÓN (primer release estable)
```

### Impacto en Workflows

- ✅ Tags en workflows deben ser v0.x.x
- ✅ Instrucciones de instalación usan v0.x.x
- ✅ CHANGELOG documenta versiones v0.x.x
- ✅ GitHub Releases usan v0.x.x

### Referencia

Ver informe detallado: `INFORME_VERSIONADO_CRITICO.md`

---

---

## 🚀 Guía Rápida: Crear Release e Imagen Docker (Nuevo Proceso Manual)

### Para api-mobile, api-administracion y worker:

```bash
# 1. Ir a GitHub Actions
https://github.com/EduGoGroup/[PROYECTO]/actions/workflows/manual-release.yml

# 2. Click "Run workflow"
Inputs:
  - Branch: main
  - Versión: 0.x.x (formato semver, sin 'v')
  - Tipo: patch / minor / major

# 3. Esperar ~20 minutos

# 4. Verificar resultados:
- Tag Git: creado ✅
- Imagen Docker: ghcr.io/edugogroup/[proyecto]:v0.x.x ✅
- GitHub Release: publicado ✅
- CHANGELOG: actualizado ✅
```

### Notas Importantes:

- ✅ **TODO en un solo workflow** (version + tag + tests + Docker + release)
- ✅ **Control total** (tú decides cuándo ejecutar)
- ✅ **Predecible** (siempre funciona igual)
- ⏱️ **Duración**: ~20 minutos (build multi-platform)
- ❌ **NO usar auto-version.yml** (fue eliminado por inestable)

---

**Responsable:** Claude Code + Jhoan Medina
**Siguiente acción:** Implementar CI/CD en edugo-api-administracion y edugo-worker con manual-release.yml
**Herramientas requeridas:** `actionlint`, `gh`, `git`
