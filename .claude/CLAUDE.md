# Instrucciones Específicas del Proyecto - EduGo API Mobile

## 📋 Plan de Trabajo Activo

**IMPORTANTE**: Este proyecto tiene un plan de trabajo en curso documentado en:

👉 **[sprint/current/readme.md](../sprint/current/readme.md)**

Antes de realizar cualquier tarea, **SIEMPRE**:
1. Leer el archivo `sprint/current/readme.md` para entender el contexto y fase actual
2. Verificar qué tareas están completadas y cuáles están pendientes
3. Seguir el orden de las fases establecidas
4. Actualizar las casillas de verificación según el progreso
5. Documentar cualquier desviación o hallazgo en el plan

---

## 🎯 Estado Actual del Proyecto

### Fase Activa
- **Fase 1**: ✅ COMPLETADA (Conectar implementación real con Container DI)
- **Branch**: `feature/conectar`
- **Último commit**: `3332c05` - "feat: conectar implementación real con Container DI"

### Próxima Fase
- **Fase 2**: ⏳ PENDIENTE (Completar TODOs de Servicios)
  - Subtarea siguiente: Implementar funcionalidad S3

---

## 🏗️ Arquitectura del Proyecto

Este proyecto implementa **Clean Architecture (Hexagonal)** con las siguientes capas:

```
internal/
├── domain/              # Entidades, Value Objects, Interfaces de repositorio
├── application/         # Servicios, DTOs, Casos de uso
├── infrastructure/      # Implementaciones concretas
│   ├── http/           # Handlers, Middleware
│   ├── persistence/    # Repositorios (PostgreSQL, MongoDB)
│   └── messaging/      # RabbitMQ (pendiente implementar)
├── container/          # Inyección de Dependencias
└── config/             # Configuración con Viper
```

---

## ⚙️ Tecnologías y Dependencias

- **Framework Web**: Gin
- **Base de Datos**: PostgreSQL (driver: lib/pq)
- **Base de Datos NoSQL**: MongoDB (mongo-driver)
- **Messaging**: RabbitMQ (pendiente conectar)
- **Storage**: AWS S3 (pendiente configurar)
- **Logging**: Zap Logger (edugo-shared)
- **Autenticación**: JWT (edugo-shared/auth)
- **Testing**: Testcontainers

---

## 📦 Trabajo con edugo-shared

### Ubicación del Proyecto Shared
```
Ruta local: /Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared
```

### 🔄 Flujo Completo de Modificación de edugo-shared

**IMPORTANTE**: Cuando necesites modificar código en `edugo-shared`, debes seguir este flujo obligatorio para mantener las versiones sincronizadas:

#### **Paso 1: Modificar código en edugo-shared**

```bash
# Navegar al proyecto shared
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared

# Verificar branch y estado
git status
git log -1 --oneline

# Hacer los cambios necesarios (ej: agregar middleware)
# ... editar archivos ...

# Compilar y verificar que no hay errores
go build ./...
go test ./...

# Commit de cambios en shared
git add .
git commit -m "feat: agregar middleware JWT para Gin"
```

#### **Paso 2: Crear Tag de Versión (OBLIGATORIO)**

```bash
# Listar tags existentes para ver última versión
git tag -l | sort -V | tail -5

# Crear nuevo tag semántico (seguir Semantic Versioning)
# Formato: vMAJOR.MINOR.PATCH o v0.0.0-YYYYMMDDHHMMSS-commit
# Ejemplos:
# - Cambio menor (nueva feature): v0.1.0 → v0.2.0
# - Parche (bugfix): v0.1.0 → v0.1.1
# - Breaking change: v0.1.0 → v1.0.0

git tag v0.2.0  # Ajustar según el tipo de cambio

# Push del tag al remote (esto genera el release en GitHub)
git push origin v0.2.0
```

#### **Paso 3: Actualizar Dependencia en edugo-api-mobile**

```bash
# Navegar de vuelta al proyecto api-mobile
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Opción A: Actualizar a tag específico
go get github.com/EduGoGroup/edugo-shared@v0.2.0

# Opción B: Actualizar a última versión
go get -u github.com/EduGoGroup/edugo-shared

# Limpiar módulos
go mod tidy

# Verificar que se actualizó correctamente
go list -m github.com/EduGoGroup/edugo-shared
# Debe mostrar: github.com/EduGoGroup/edugo-shared v0.2.0

# Compilar para verificar compatibilidad
go build ./...
```

#### **Paso 4: Commit de Actualización en api-mobile**

```bash
# Agregar go.mod y go.sum al staging
git add go.mod go.sum

# Commit indicando la actualización
git commit -m "chore: actualizar edugo-shared a v0.2.0

- Actualizar dependencia de edugo-shared
- Incluye nuevo middleware JWT para Gin

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

### 📋 Checklist de Modificación de Shared

Cada vez que modifiques `edugo-shared`, verifica:

- [ ] Los cambios están en el branch correcto de shared
- [ ] Código compila sin errores (`go build ./...`)
- [ ] Tests pasan (`go test ./...`)
- [ ] Commit creado en shared con mensaje descriptivo
- [ ] **Tag de versión creado** (vMAJOR.MINOR.PATCH)
- [ ] Tag pusheado a GitHub (`git push origin <tag>`)
- [ ] Dependencia actualizada en api-mobile (`go get`)
- [ ] go.mod y go.sum actualizados
- [ ] api-mobile compila con nueva versión
- [ ] Commit de actualización creado en api-mobile

### ⚠️ Reglas de Versionado Semántico

| Tipo de Cambio | Ejemplo | Versión |
|----------------|---------|---------|
| **Breaking Change** | Cambiar firma de función pública | v0.1.0 → v1.0.0 |
| **Nueva Feature** | Agregar middleware nuevo | v0.1.0 → v0.2.0 |
| **Bugfix** | Corregir error en logger | v0.1.0 → v0.1.1 |
| **Desarrollo** | Cambios experimentales | v0.0.0-20251031... |

### 🚨 Errores Comunes a Evitar

❌ **NO HACER**:
- Modificar shared sin crear tag
- Olvidar hacer `go get` en api-mobile
- Pushear código que no compila
- Usar versiones en desarrollo (commit hash) en producción

✅ **SÍ HACER**:
- Siempre crear tag después de commit en shared
- Actualizar inmediatamente api-mobile
- Verificar compilación en ambos proyectos
- Documentar breaking changes en mensaje de commit

### 📚 Paquetes Disponibles en edugo-shared

```
edugo-shared/
├── auth/               # JWT Manager, autenticación
├── logger/             # Logger Zap estructurado
├── common/
│   └── errors/        # Error types de aplicación
└── (pendientes)
    ├── middleware/    # Middleware reutilizable (próximo)
    └── utils/         # Utilidades comunes
```

### 🔗 Referencias Útiles

- Repo shared: `https://github.com/EduGoGroup/edugo-shared`
- Go modules docs: `https://go.dev/ref/mod`
- Semantic Versioning: `https://semver.org`

---

## 🔐 Variables de Entorno Requeridas

El proyecto requiere las siguientes variables de entorno para funcionar:

```bash
# Base de datos
POSTGRES_PASSWORD=<contraseña_postgres>
MONGODB_URI=mongodb://<usuario>:<password>@<host>:<puerto>

# Messaging
RABBITMQ_URL=amqp://<usuario>:<password>@<host>:<puerto>

# Autenticación
JWT_SECRET=<secret_para_jwt>

# Ambiente
APP_ENV=local|dev|qa|prod
```

---

## 🎯 Sistema de Análisis de Sprint

El proyecto cuenta con un sistema flexible de análisis arquitectónico:

### Comando Principal: `/01-analysis`

```bash
# Sintaxis
/01-analysis [--source=sprint|current] [--phase=N] [--mode=full|quick]

# Ejemplos
/01-analysis                              # Análisis completo de sprint/current
/01-analysis --mode=quick                 # Análisis rápido sin diagramas
/01-analysis --source=sprint              # Analizar sprint/readme.md
/01-analysis --source=sprint --phase=3    # Solo fase 3 del sprint root
/01-analysis --phase=2 --mode=quick       # Fase 2 rápido de current
```

### Atajo: `/01-quick-analysis`

```bash
# Equivalente a /01-analysis --mode=quick
/01-quick-analysis
/01-quick-analysis --source=sprint
/01-quick-analysis --phase=3
```

### Parámetros

| Parámetro | Valores | Default | Descripción |
|-----------|---------|---------|-------------|
| `--source` | `sprint`, `current` | `current` | De dónde leer el readme |
| `--phase` | Número de fase | todas | Analizar solo una fase |
| `--mode` | `full`, `quick` | `full` | Con o sin diagramas |

### Modos de Análisis

**MODE=full** (Completo):
- `architecture.md` - Diagramas de arquitectura (Mermaid)
- `data-model.md` - Diagramas ER
- `process-diagram.md` - Diagramas de flujo
- `readme.md` - Resumen ejecutivo

**MODE=quick** (Rápido):
- `readme.md` - Solo análisis ejecutivo (sin diagramas)

### Alcance de Análisis

**SCOPE=complete**:
- Analiza todo el contenido del readme
- Archivos: `architecture.md`, etc.

**SCOPE=phase-N**:
- Enfoque en fase específica
- Archivos: `architecture-phase-3.md`, etc.
- Incluye resumen general como contexto

---

## 🔧 Sistema de Revisión Automática de PRs

El proyecto cuenta con un sistema automatizado para revisar y corregir Pull Requests.

### Comando: `/05-pr-fix`

Invoca al agente especializado **flow-pr-fixer** que analiza PRs, clasifica comentarios de reviewers y aplica correcciones automáticas.

```bash
# Sintaxis
/05-pr-fix [--pr=NUMBER] [--auto-fix] [--branch=NAME]

# Ejemplos
/05-pr-fix                      # Revisar PR del branch actual
/05-pr-fix --auto-fix           # Revisar y aplicar correcciones inmediatas
/05-pr-fix --pr=123             # Revisar PR específico
/05-pr-fix --pr=456 --auto-fix  # Revisar PR específico con auto-corrección
```

### Flujo de Trabajo

1. **Conectar al PR**: Obtiene información del PR (activo o especificado)
2. **Verificar Pipelines**: Revisa estado de checks (build, linting, tests)
3. **Obtener Comentarios**: Lee comentarios de Copilot, Claude Web, reviewers humanos
4. **Clasificar Comentarios**: Categoriza según criterios predefinidos
5. **Aplicar Correcciones**: Corrige automáticamente issues obvios (si --auto-fix)
6. **Generar Informe**: Crea reporte estructurado con clasificación completa

### Clasificación de Comentarios

| Categoría | Símbolo | Descripción | Acción |
|-----------|---------|-------------|--------|
| **2.1 - Corrección Inmediata** | 🟢 | Typos, formato, linting, imports | Corregir automáticamente |
| **2.2 - Traducciones/Docs** | 🔵 | Traducción texto, mejoras docs | Excluir (fuera de scope) |
| **2.3 - Deuda Técnica** | 🟡 | Refactorización, arquitectura | Documentar para después |
| **2.4 - No Relevantes** | ⚪ | Preferencias personales, ya implementados | Descartar con razón |
| **2.5 - Dudosos** | 🟣 | Ambiguos, múltiples opciones | Pedir decisión al usuario |

### Informe Generado

El agente genera un informe markdown con:

- **Resumen Ejecutivo**: Cantidad de comentarios por categoría
- **Estado de Pipelines**: Estado de todos los checks (build, linting, tests)
- **Correcciones Aplicadas**: Lista de fixes automáticos realizados
- **Deuda Técnica**: Items con justificación, impacto, esfuerzo y prioridad
- **Comentarios Dudosos**: Opciones para el usuario (inmediato, deuda, descartar)
- **Próximos Pasos**: Acciones recomendadas

### Ejemplo de Uso Típico

```bash
# 1. Crear PR y esperar reviews de Copilot/Claude
git push
gh pr create

# 2. Revisar comentarios (sin aplicar correcciones)
/05-pr-fix

# 3. Leer informe y decidir sobre comentarios dudosos
[Revisar informe generado]

# 4. Aplicar correcciones aprobadas
/05-pr-fix --auto-fix

# 5. Crear documento de deuda técnica si es necesario
[Usar informe para crear tech-debt.md]

# 6. Commit y push
git add .
git commit -m "fix: aplicar correcciones de PR review"
git push
```

### Requisitos

- **GitHub CLI** (`gh`) instalado y autenticado, O
- **MCP GitHub** configurado en `.claude/settings.json`
- **Permisos** de lectura/escritura en el repositorio
- **Branch** debe estar asociado a un PR abierto (si no se usa --pr)

### Documentación Completa

- Agente: `.claude/agents/flow-pr-fixer.md`
- Comando: `.claude/commands/05-pr-fix.md`

---

## 📁 Archivos de Configuración

- `config/config.yaml` - Configuración base
- `config/config-{env}.yaml` - Configuración por ambiente
- Las contraseñas y secrets vienen de variables de entorno

---

## 🚨 Reglas Específicas de Este Proyecto

### Commits
1. **NUNCA** hacer commit si el proyecto tiene errores de compilación
2. Solo hacer commits atómicos según lo planeado en `sprint/current/readme.md`
3. Seguir el formato de commit establecido (feat, fix, refactor, test, etc.)
4. Incluir siempre el footer de Claude Code en commits
5. Actualizar el `sprint/current/readme.md` marcando casillas al completar tareas

### Manejo de Errores
1. Usar los error types de `edugo-shared/common/errors`
2. Siempre hacer logging de errores con contexto
3. Retornar errores de aplicación apropiados en handlers
4. No silenciar errores, propagarlos hasta el handler

### Testing
1. Usar testcontainers para tests de integración
2. Los tests deben ser independientes y poder ejecutarse en paralelo
3. Limpiar recursos después de cada test

### Código Duplicado
- **IMPORTANTE**: Existen handlers duplicados:
  - `internal/handlers/` (VIEJOS, con mocks) - **NO USAR**
  - `internal/infrastructure/http/handler/` (NUEVOS, reales) - **USAR ESTOS**
- Los handlers viejos serán eliminados en Fase 3 del plan

---

## 🔄 Flujo de Trabajo

### Al comenzar una sesión:
```bash
1. git status                         # Ver estado actual
2. cat sprint/current/readme.md       # Revisar plan de trabajo
3. git log -1 --oneline               # Ver último commit
```

### Durante el desarrollo:
1. Seguir las tareas del `sprint/current/readme.md` en orden
2. Marcar casillas completadas
3. Hacer commits atómicos según lo planeado
4. **NO HACER PUSH** sin autorización del usuario

### Al finalizar una fase:
1. Actualizar `sprint/current/readme.md` con estado ✅
2. Documentar hallazgos o cambios al plan
3. Preparar contexto para próxima fase

---

## 📚 Documentación Adicional

- Swagger UI disponible en: `http://localhost:8080/swagger/index.html`
- Generar docs Swagger: `make swagger` o `swag init -g cmd/main.go`
- Health check endpoint: `GET /health`

---

## 🎯 Objetivo Final del Sprint

Completar la migración de handlers mock a implementación real, eliminando código duplicado y completando todas las funcionalidades pendientes (S3, RabbitMQ, queries complejas).

**Estado**: 1/6 commits completados (16.6% del sprint)

---

**Última actualización**: 2025-11-05 (v3 - Agregado sistema /pr-fix para revisión automática de PRs)
**Responsable**: Claude Code + Jhoan Medina
