# 📚 EduGo API Mobile - Documentación Técnica

> **Versión:** 1.0  
> **Última actualización:** Diciembre 2024  
> **Framework:** Go 1.25 + Gin  
> **Licencia:** MIT

---

## 🎯 Visión General

**EduGo API Mobile** es el backend REST API diseñado para operaciones frecuentes de docentes y estudiantes en la plataforma educativa EduGo. Esta API está optimizada para dispositivos móviles y maneja las funcionalidades core de la experiencia de aprendizaje.

### ¿Qué es EduGo?

EduGo es una plataforma educativa integral que permite a instituciones educativas digitalizar su proceso de enseñanza-aprendizaje. La plataforma se compone de múltiples servicios:

- **api-admin:** Gestión administrativa, autenticación centralizada, gestión de usuarios y escuelas
- **api-mobile:** (Este servicio) Operaciones frecuentes para docentes y estudiantes
- **worker:** Procesamiento asíncrono de PDFs, generación de resúmenes y quizzes con IA
- **frontend-web:** Panel administrativo web
- **frontend-mobile:** Aplicaciones móviles para iOS y Android

### Propósito Principal de API Mobile

| Funcionalidad | Descripción | Usuarios |
|---------------|-------------|----------|
| **Materiales Educativos** | CRUD de PDFs y documentos, URLs presignadas para S3, versionado | Docentes |
| **Evaluaciones** | Quizzes generados por IA, scoring automático, feedback detallado | Estudiantes |
| **Progreso de Lectura** | Tracking de avance, última página, tiempo de lectura | Estudiantes |
| **Resúmenes IA** | Consulta de resúmenes generados automáticamente | Todos |
| **Estadísticas** | Métricas de uso, completion rate, scores promedio | Docentes, Admins |

### Flujo Típico de Uso

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           FLUJO DE USO TÍPICO                               │
└─────────────────────────────────────────────────────────────────────────────┘

  DOCENTE                                              ESTUDIANTE
     │                                                      │
     │ 1. Sube PDF                                          │
     │ ──────────────────────────────────▶                  │
     │                                                      │
     │ 2. Worker procesa (IA)                               │
     │    • Genera resumen                                  │
     │    • Genera quiz                                     │
     │                                                      │
     │                                                      │ 3. Consulta materiales
     │                                                      │ ◀──────────────────────
     │                                                      │
     │                                                      │ 4. Descarga y lee PDF
     │                                                      │ ◀──────────────────────
     │                                                      │
     │                                                      │ 5. Actualiza progreso
     │                                                      │ ──────────────────────▶
     │                                                      │
     │                                                      │ 6. Completa quiz
     │                                                      │ ──────────────────────▶
     │                                                      │
     │                                                      │ 7. Recibe feedback
     │                                                      │ ◀──────────────────────
     │                                                      │
     │ 8. Ve estadísticas                                   │
     │ ◀──────────────────────────────────                  │
     │                                                      │
```

---

## 📖 Índice de Documentación

### Documentación Principal

| Documento | Descripción | Audiencia |
|-----------|-------------|-----------|
| [ARCHITECTURE.md](./ARCHITECTURE.md) | Arquitectura del sistema, patrones de diseño, capas, DI | Desarrolladores Backend |
| [DATABASE.md](./DATABASE.md) | Esquemas PostgreSQL y MongoDB, relaciones, índices, queries | Desarrolladores, DBAs |
| [API-REFERENCE.md](./API-REFERENCE.md) | Documentación completa de ~18 endpoints REST | Desarrolladores Frontend/Mobile |
| [SETUP.md](./SETUP.md) | Guía de instalación, Docker, configuración, troubleshooting | Nuevos desarrolladores |
| [FLOWS.md](./FLOWS.md) | Diagramas de flujos y procesos de negocio | Product, Desarrolladores |

### Documentación de Mejoras

| Documento | Descripción | Prioridad |
|-----------|-------------|-----------|
| [improvements/README.md](./improvements/README.md) | Índice de mejoras pendientes | - |
| [improvements/DEPRECATED-CODE.md](./improvements/DEPRECATED-CODE.md) | Código marcado para eliminación | 🔴 Alta |
| [improvements/TODO-ITEMS.md](./improvements/TODO-ITEMS.md) | TODOs pendientes en el código | 🟡 Media |
| [improvements/LEGACY-ENDPOINTS.md](./improvements/LEGACY-ENDPOINTS.md) | Endpoints legacy a migrar | 🟡 Media |
| [improvements/TECHNICAL-DEBT.md](./improvements/TECHNICAL-DEBT.md) | Deuda técnica acumulada | 🔴 Alta |
| [improvements/REFACTORING-OPPORTUNITIES.md](./improvements/REFACTORING-OPPORTUNITIES.md) | Oportunidades de mejora | 🟢 Baja |

---

## 🏗️ Stack Tecnológico

### Backend
| Tecnología | Uso |
|------------|-----|
| **Go 1.25** | Lenguaje principal |
| **Gin** | Framework HTTP |
| **GORM** | ORM para PostgreSQL |
| **Mongo Driver** | Driver oficial MongoDB |
| **Swagger/Swag** | Documentación API automática |

### Bases de Datos
| Servicio | Propósito |
|----------|-----------|
| **PostgreSQL 16** | Datos relacionales (usuarios, materiales, progreso) |
| **MongoDB 7.0** | Datos no estructurados (assessments, resúmenes IA) |

### Servicios de Soporte
| Servicio | Propósito |
|----------|-----------|
| **RabbitMQ 3.12** | Cola de mensajes para procesamiento async |
| **AWS S3** | Almacenamiento de archivos (PDFs) |

---

## 🔐 Autenticación

La API utiliza **JWT Bearer Token** para autenticación. Los tokens son emitidos por `api-admin` (servicio centralizado de autenticación EduGo).

```
Authorization: Bearer <jwt_token>
```

### Validación de Tokens
- **Modo Local:** Validación JWT con secreto compartido
- **Modo Producción:** Validación remota contra api-admin (opcional)

---

## 📊 Resumen de Endpoints

| Grupo | Endpoints | Auth | Descripción |
|-------|-----------|------|-------------|
| **Materials** | 8 | ✅ | CRUD de materiales, URLs presignadas S3 |
| **Assessments** | 5 | ✅ | Quizzes, intentos, resultados |
| **Progress** | 2 | ✅ | Progreso de lectura |
| **Stats** | 2 | ✅ | Estadísticas globales y por material |
| **Health** | 1 | ❌ | Health check del servicio |

**Total:** ~18 endpoints REST

---

## 🚀 Quick Start

```bash
# 1. Clonar y configurar
cp .env.example .env

# 2. Levantar infraestructura (PostgreSQL, MongoDB, RabbitMQ)
docker-compose up -d postgres mongodb rabbitmq

# 3. Ejecutar API
make run

# 4. Verificar
curl http://localhost:8080/health
```

**Swagger UI:** http://localhost:8080/swagger/index.html

---

## 📁 Estructura del Proyecto

```
edugo-api-mobile/
├── cmd/
│   └── main.go              # Entry point
├── config/
│   ├── config.yaml          # Configuración base
│   └── config-{env}.yaml    # Overrides por ambiente
├── internal/
│   ├── application/
│   │   ├── dto/             # Data Transfer Objects
│   │   └── service/         # Casos de uso
│   ├── bootstrap/           # Inicialización
│   ├── container/           # Dependency Injection
│   ├── domain/
│   │   ├── repository/      # Interfaces de repositorios
│   │   └── valueobject/     # Value Objects
│   └── infrastructure/
│       ├── http/
│       │   ├── handler/     # Controladores HTTP
│       │   ├── middleware/  # Auth, CORS, etc.
│       │   └── router/      # Configuración de rutas
│       ├── messaging/       # RabbitMQ
│       ├── persistence/     # PostgreSQL & MongoDB
│       └── storage/         # AWS S3
├── docs/                    # Swagger generado
├── documents/               # Documentación técnica
└── test/                    # Tests de integración
```

---

## 🔗 Relación con Otros Servicios

```
                    ┌──────────────────┐
                    │   api-admin      │
                    │ (autenticación)  │
                    └────────┬─────────┘
                             │ JWT validation
                             ▼
┌─────────────┐     ┌──────────────────┐     ┌─────────────┐
│   Mobile    │────▶│  api-mobile      │────▶│  Worker     │
│   Apps      │     │  (este servicio) │     │  (PDF proc) │
└─────────────┘     └────────┬─────────┘     └─────────────┘
                             │
              ┌──────────────┼──────────────┐
              ▼              ▼              ▼
        ┌──────────┐  ┌──────────┐  ┌──────────┐
        │PostgreSQL│  │ MongoDB  │  │   S3     │
        └──────────┘  └──────────┘  └──────────┘
```

---

## 🧪 Testing

### Tipos de Tests

| Tipo | Comando | Descripción | Tiempo |
|------|---------|-------------|--------|
| **Unitarios** | `make test-unit` | Tests rápidos sin dependencias externas | ~30s |
| **Integración** | `make test-integration` | Tests con testcontainers (PostgreSQL, MongoDB) | ~5min |
| **Todos** | `make test` | Unitarios + Integración | ~6min |
| **Cobertura** | `make coverage-report` | Genera reporte HTML de cobertura | ~7min |

### Cobertura Actual

```
Objetivo mínimo: 60%
Cobertura actual: ~65-70%
```

### Ejecutar Tests Específicos

```bash
# Tests de un paquete específico
go test -v ./internal/application/service/...

# Tests con nombre específico
go test -v -run TestMaterialService ./internal/application/service/...

# Tests de integración solamente
RUN_INTEGRATION_TESTS=true go test -v -tags=integration ./test/integration/...
```

---

## 📊 Métricas del Proyecto

### Estadísticas de Código

| Métrica | Valor |
|---------|-------|
| **Líneas de código Go** | ~15,000 |
| **Archivos .go** | ~120 |
| **Tests** | ~80 archivos |
| **Endpoints** | ~18 |
| **Handlers** | 5 (Material, Assessment, Progress, Stats, Summary) |
| **Services** | 6 |
| **Repositories** | 8 |

### Dependencias Principales

```
github.com/gin-gonic/gin          v1.11.0   # Framework HTTP
github.com/google/uuid            v1.6.0    # UUIDs
github.com/lib/pq                 v1.10.9   # PostgreSQL driver
go.mongodb.org/mongo-driver       v1.17.6   # MongoDB driver
github.com/rabbitmq/amqp091-go    v1.10.0   # RabbitMQ client
github.com/aws/aws-sdk-go-v2      v1.39.5   # AWS SDK
github.com/stretchr/testify       v1.11.1   # Testing
gorm.io/gorm                      v1.25.12  # ORM
```

---

## 🔒 Seguridad

### Prácticas Implementadas

- ✅ **JWT Authentication:** Tokens firmados con HMAC-SHA256
- ✅ **Input Validation:** Validación en DTOs con tags de binding
- ✅ **SQL Injection Prevention:** Uso de prepared statements y ORM
- ✅ **CORS Configuration:** Headers configurados en middleware
- ✅ **Rate Limiting:** (Pendiente - implementar en API Gateway)
- ✅ **Secrets Management:** Variables de entorno, nunca en código

### Reportar Vulnerabilidades

Si encuentras una vulnerabilidad de seguridad, por favor repórtala a: security@edugo.com

---

## 📝 Changelog

Ver [CHANGELOG.md](../CHANGELOG.md) para el historial de versiones.

### Versiones Recientes

| Versión | Fecha | Highlights |
|---------|-------|------------|
| Sprint-04 | Dic 2024 | Sistema de evaluaciones con PostgreSQL |
| Sprint-03 | Nov 2024 | Migración de assessments a PostgreSQL |
| Sprint-02 | Oct 2024 | Sistema de progreso idempotente |
| Sprint-01 | Sep 2024 | MVP con materiales y assessments |

---

## 🤝 Contribuir

### Guía de Contribución

1. **Fork** el repositorio
2. **Crear branch:** `git checkout -b feature/mi-feature`
3. **Hacer cambios** siguiendo las guías de estilo
4. **Ejecutar tests:** `make test`
5. **Ejecutar linter:** `make lint`
6. **Commit:** `git commit -m "feat: descripción"`
7. **Push:** `git push origin feature/mi-feature`
8. **Pull Request:** Crear PR con descripción detallada

### Convención de Commits

```
feat: nueva funcionalidad
fix: corrección de bug
docs: documentación
style: formato, sin cambios de código
refactor: refactorización
test: tests
chore: mantenimiento
```

---

## 📞 Contacto

- **Equipo:** EduGo Development Team
- **Email:** soporte@edugo.com
- **Repositorio:** github.com/EduGoGroup/edugo-api-mobile
- **Issues:** github.com/EduGoGroup/edugo-api-mobile/issues

---

## 📜 Licencia

Este proyecto está bajo la licencia MIT. Ver [LICENSE](../LICENSE) para más detalles.
