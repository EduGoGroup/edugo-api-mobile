# Instrucciones Específicas del Proyecto - EduGo API Mobile

## 📋 Plan de Trabajo Activo

**IMPORTANTE**: Este proyecto tiene un plan de trabajo en curso documentado en:

👉 **[sprint/README.md](../sprint/README.md)**

Antes de realizar cualquier tarea, **SIEMPRE**:
1. Leer el archivo `sprint/README.md` para entender el contexto y fase actual
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

## 📁 Archivos de Configuración

- `config/config.yaml` - Configuración base
- `config/config-{env}.yaml` - Configuración por ambiente
- Las contraseñas y secrets vienen de variables de entorno

---

## 🚨 Reglas Específicas de Este Proyecto

### Commits
1. **NUNCA** hacer commit si el proyecto tiene errores de compilación
2. Solo hacer commits atómicos según lo planeado en `sprint/README.md`
3. Seguir el formato de commit establecido (feat, fix, refactor, test, etc.)
4. Incluir siempre el footer de Claude Code en commits
5. Actualizar el `sprint/README.md` marcando casillas al completar tareas

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
1. git status                    # Ver estado actual
2. cat sprint/README.md          # Revisar plan de trabajo
3. git log -1 --oneline          # Ver último commit
```

### Durante el desarrollo:
1. Seguir las tareas del `sprint/README.md` en orden
2. Marcar casillas completadas
3. Hacer commits atómicos según lo planeado
4. **NO HACER PUSH** sin autorización del usuario

### Al finalizar una fase:
1. Actualizar `sprint/README.md` con estado ✅
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

**Última actualización**: 2025-10-31 (v2 - Agregado flujo edugo-shared)
**Responsable**: Claude Code + Jhoan Medina
