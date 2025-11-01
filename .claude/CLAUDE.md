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

**Última actualización**: 2025-10-31
**Responsable**: Claude Code + Jhoan Medina
