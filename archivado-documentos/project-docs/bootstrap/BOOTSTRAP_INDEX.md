# Índice de Documentación del Sistema de Bootstrap

Este documento proporciona un índice completo de toda la documentación relacionada con el sistema de bootstrap de infraestructura.

## 📚 Documentación Principal

### Para Usuarios Nuevos

1. **[BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md)** - Guía completa de uso
   - Conceptos básicos del sistema de bootstrap
   - Inicialización de recursos
   - Recursos opcionales
   - Inyección de mocks
   - Gestión del ciclo de vida
   - Ejemplos prácticos
   - Best practices

### Para Migración

2. **[BOOTSTRAP_MIGRATION_GUIDE.md](BOOTSTRAP_MIGRATION_GUIDE.md)** - Guía de migración
   - ¿Por qué migrar?
   - Cambios principales
   - Migración paso a paso
   - Ejemplos de código antes/después
   - Testing con bootstrap
   - Troubleshooting

### Para Testing

3. **[../internal/bootstrap/INTEGRATION_TESTS.md](../internal/bootstrap/INTEGRATION_TESTS.md)** - Guía de testing
   - Estrategia de testing
   - Tests unitarios con mocks
   - Tests de integración
   - Ejemplos de tests
   - Helpers de testing

### Para Configuración

4. **[../config/OPTIONAL_RESOURCES.md](../config/OPTIONAL_RESOURCES.md)** - Recursos opcionales
   - Configuración de recursos opcionales
   - Recursos disponibles
   - Implementaciones noop
   - Ejemplos de configuración

## 🎯 Guías Rápidas

### Quick Start

- **[../README.md](../README.md)** - README principal con sección de bootstrap
- **[../QUICKSTART.md](../QUICKSTART.md)** - Guía de inicio rápido

### Referencia Rápida

**Inicialización Básica:**
```go
b := bootstrap.New(cfg)
resources, cleanup, err := b.InitializeInfrastructure(ctx)
if err != nil {
    log.Fatal(err)
}
defer cleanup()
```

**Recursos Opcionales:**
```go
b := bootstrap.New(cfg,
    bootstrap.WithOptionalResource("rabbitmq"),
    bootstrap.WithOptionalResource("s3"),
)
```

**Testing con Mocks:**
```go
b := bootstrap.New(cfg,
    bootstrap.WithLogger(mockLogger),
    bootstrap.WithPostgreSQL(mockDB),
    bootstrap.WithMongoDB(mockMongoDB),
    bootstrap.WithRabbitMQ(mockPublisher),
    bootstrap.WithS3Client(mockS3),
)
```

## 📖 Documentación por Caso de Uso

### Desarrollo Local

**Objetivo**: Ejecutar la aplicación sin toda la infraestructura

**Documentos relevantes**:
1. [BOOTSTRAP_USAGE.md - Recursos Opcionales](BOOTSTRAP_USAGE.md#recursos-opcionales)
2. [../config/OPTIONAL_RESOURCES.md](../config/OPTIONAL_RESOURCES.md)
3. [../QUICKSTART.md - Desarrollo sin Infraestructura Completa](../QUICKSTART.md#-desarrollo-sin-infraestructura-completa)

**Ejemplo**:
```go
b := bootstrap.New(cfg,
    bootstrap.WithOptionalResource("rabbitmq"),
    bootstrap.WithOptionalResource("s3"),
)
```

### Testing Unitario

**Objetivo**: Escribir tests con mocks de todos los recursos

**Documentos relevantes**:
1. [BOOTSTRAP_USAGE.md - Inyección de Mocks](BOOTSTRAP_USAGE.md#inyección-de-mocks)
2. [../internal/bootstrap/INTEGRATION_TESTS.md - Tests Unitarios](../internal/bootstrap/INTEGRATION_TESTS.md#tests-unitarios)
3. [BOOTSTRAP_USAGE.md - Ejemplo 3: Tests Unitarios](BOOTSTRAP_USAGE.md#ejemplo-3-tests-unitarios)

**Ejemplo**:
```go
func TestService(t *testing.T) {
    b := bootstrap.New(cfg,
        bootstrap.WithLogger(mockLogger),
        bootstrap.WithPostgreSQL(mockDB),
        bootstrap.WithMongoDB(mockMongoDB),
        bootstrap.WithRabbitMQ(mockPublisher),
        bootstrap.WithS3Client(mockS3),
    )
    resources, cleanup, _ := b.InitializeInfrastructure(ctx)
    defer cleanup()
    // Tests...
}
```

### Testing de Integración

**Objetivo**: Ejecutar tests con recursos reales o testcontainers

**Documentos relevantes**:
1. [../internal/bootstrap/INTEGRATION_TESTS.md - Tests de Integración](../internal/bootstrap/INTEGRATION_TESTS.md#tests-de-integración)
2. [BOOTSTRAP_USAGE.md - Ejemplo 4: Tests de Integración](BOOTSTRAP_USAGE.md#ejemplo-4-tests-de-integración)

**Ejemplo**:
```go
func TestIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    b := bootstrap.New(cfg)
    resources, cleanup, err := b.InitializeInfrastructure(ctx)
    require.NoError(t, err)
    defer cleanup()
    // Tests con recursos reales...
}
```

### Migración de Código Legacy

**Objetivo**: Actualizar código existente al sistema de bootstrap

**Documentos relevantes**:
1. [BOOTSTRAP_MIGRATION_GUIDE.md](BOOTSTRAP_MIGRATION_GUIDE.md)
2. [BOOTSTRAP_MIGRATION_GUIDE.md - Migración Paso a Paso](BOOTSTRAP_MIGRATION_GUIDE.md#migración-paso-a-paso)
3. [../README.md - Guía de Migración](../README.md#guía-de-migración-al-sistema-de-bootstrap)

**Pasos**:
1. Actualizar `main.go`
2. Actualizar `container`
3. Actualizar tests

### Producción

**Objetivo**: Desplegar aplicación con todos los recursos

**Documentos relevantes**:
1. [BOOTSTRAP_USAGE.md - Ejemplo 1: Aplicación de Producción](BOOTSTRAP_USAGE.md#ejemplo-1-aplicación-de-producción)
2. [../README.md - Configuración](../README.md#configuración)

**Ejemplo**:
```go
func main() {
    ctx := context.Background()
    cfg, _ := config.Load()

    b := bootstrap.New(cfg)
    resources, cleanup, err := b.InitializeInfrastructure(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer cleanup()

    container := container.NewContainer(resources)
    router := router.SetupRouter(container)
    router.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}
```

## 🔍 Búsqueda por Tema

### Conceptos

| Tema | Documento | Sección |
|------|-----------|---------|
| ¿Qué es el bootstrap? | [BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md) | Introducción |
| Resources struct | [BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md) | Conceptos Básicos - Resources |
| Bootstrapper | [BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md) | Conceptos Básicos - Bootstrapper |
| BootstrapOptions | [BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md) | Conceptos Básicos - BootstrapOptions |

### Recursos

| Recurso | Documento | Sección |
|---------|-----------|---------|
| Logger | [BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md) | Acceso a Recursos |
| PostgreSQL | [BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md) | Acceso a Recursos |
| MongoDB | [BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md) | Acceso a Recursos |
| RabbitMQ | [BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md) | Acceso a Recursos |
| S3 | [BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md) | Acceso a Recursos |

### Recursos Opcionales

| Tema | Documento | Sección |
|------|-----------|---------|
| ¿Qué son? | [BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md) | Recursos Opcionales |
| Configuración YAML | [../config/OPTIONAL_RESOURCES.md](../config/OPTIONAL_RESOURCES.md) | Configuración |
| Configuración código | [BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md) | Recursos Opcionales - Configuración |
| Implementaciones noop | [BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md) | Recursos Opcionales - Implementaciones Noop |
| NoopPublisher | [../internal/bootstrap/noop/publisher.go](../internal/bootstrap/noop/publisher.go) | Código fuente |
| NoopS3Storage | [../internal/bootstrap/noop/storage.go](../internal/bootstrap/noop/storage.go) | Código fuente |

### Testing

| Tema | Documento | Sección |
|------|-----------|---------|
| Inyección de mocks | [BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md) | Inyección de Mocks |
| Tests unitarios | [../internal/bootstrap/INTEGRATION_TESTS.md](../internal/bootstrap/INTEGRATION_TESTS.md) | Tests Unitarios |
| Tests de integración | [../internal/bootstrap/INTEGRATION_TESTS.md](../internal/bootstrap/INTEGRATION_TESTS.md) | Tests de Integración |
| Helpers de testing | [BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md) | Best Practices - Crear Helpers |
| Ejemplos de tests | [BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md) | Ejemplos Prácticos |

### Migración

| Tema | Documento | Sección |
|------|-----------|---------|
| ¿Por qué migrar? | [BOOTSTRAP_MIGRATION_GUIDE.md](BOOTSTRAP_MIGRATION_GUIDE.md) | ¿Por Qué Migrar? |
| Cambios principales | [BOOTSTRAP_MIGRATION_GUIDE.md](BOOTSTRAP_MIGRATION_GUIDE.md) | Cambios Principales |
| Migrar main.go | [BOOTSTRAP_MIGRATION_GUIDE.md](BOOTSTRAP_MIGRATION_GUIDE.md) | Paso 1: Actualizar main.go |
| Migrar container | [BOOTSTRAP_MIGRATION_GUIDE.md](BOOTSTRAP_MIGRATION_GUIDE.md) | Paso 2: Actualizar Container |
| Migrar servicios | [BOOTSTRAP_MIGRATION_GUIDE.md](BOOTSTRAP_MIGRATION_GUIDE.md) | Paso 3: Actualizar Referencias |
| Troubleshooting | [BOOTSTRAP_MIGRATION_GUIDE.md](BOOTSTRAP_MIGRATION_GUIDE.md) | Troubleshooting |

### Configuración

| Tema | Documento | Sección |
|------|-----------|---------|
| Configuración básica | [BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md) | Uso Básico |
| Recursos opcionales | [../config/OPTIONAL_RESOURCES.md](../config/OPTIONAL_RESOURCES.md) | Todo el documento |
| Variables de entorno | [../README.md](../README.md) | Configuración |
| Archivos YAML | [../config/README.md](../config/README.md) | Configuración |

### Ciclo de Vida

| Tema | Documento | Sección |
|------|-----------|---------|
| Cleanup automático | [BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md) | Gestión del Ciclo de Vida |
| Orden de cleanup | [BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md) | Gestión del Ciclo de Vida - Orden |
| Errores en cleanup | [BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md) | Gestión del Ciclo de Vida - Manejo de Errores |
| Cleanup manual | [BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md) | Gestión del Ciclo de Vida - Cleanup Manual |

## 🛠️ Especificación Técnica

Para desarrolladores que necesitan entender la implementación interna:

- **[../.kiro/specs/infrastructure-bootstrap-refactor/requirements.md](../.kiro/specs/infrastructure-bootstrap-refactor/requirements.md)** - Requisitos del sistema
- **[../.kiro/specs/infrastructure-bootstrap-refactor/design.md](../.kiro/specs/infrastructure-bootstrap-refactor/design.md)** - Diseño detallado
- **[../.kiro/specs/infrastructure-bootstrap-refactor/tasks.md](../.kiro/specs/infrastructure-bootstrap-refactor/tasks.md)** - Plan de implementación

## 📝 Código Fuente

Archivos principales del sistema de bootstrap:

```
internal/bootstrap/
├── bootstrap.go          # Orquestador principal
├── interfaces.go         # Interfaces de recursos
├── config.go             # Configuración y opciones
├── factories.go          # Factories de recursos
├── lifecycle.go          # Gestión de ciclo de vida
├── bootstrap_test.go     # Tests unitarios
├── lifecycle_test.go     # Tests de lifecycle
├── bootstrap_integration_test.go  # Tests de integración
├── INTEGRATION_TESTS.md  # Documentación de tests
└── noop/
    ├── publisher.go      # Implementación noop de RabbitMQ
    └── storage.go        # Implementación noop de S3
```

## 🎓 Tutoriales y Ejemplos

### Tutorial 1: Primera Aplicación con Bootstrap

1. Lee [BOOTSTRAP_USAGE.md - Uso Básico](BOOTSTRAP_USAGE.md#uso-básico)
2. Sigue [QUICKSTART.md](../QUICKSTART.md)
3. Revisa [BOOTSTRAP_USAGE.md - Ejemplo 1](BOOTSTRAP_USAGE.md#ejemplo-1-aplicación-de-producción)

### Tutorial 2: Desarrollo Local sin Infraestructura Completa

1. Lee [BOOTSTRAP_USAGE.md - Recursos Opcionales](BOOTSTRAP_USAGE.md#recursos-opcionales)
2. Configura según [config/OPTIONAL_RESOURCES.md](../config/OPTIONAL_RESOURCES.md)
3. Revisa [BOOTSTRAP_USAGE.md - Ejemplo 2](BOOTSTRAP_USAGE.md#ejemplo-2-desarrollo-local)

### Tutorial 3: Escribir Tests con Bootstrap

1. Lee [BOOTSTRAP_USAGE.md - Inyección de Mocks](BOOTSTRAP_USAGE.md#inyección-de-mocks)
2. Sigue [internal/bootstrap/INTEGRATION_TESTS.md](../internal/bootstrap/INTEGRATION_TESTS.md)
3. Revisa [BOOTSTRAP_USAGE.md - Ejemplo 3](BOOTSTRAP_USAGE.md#ejemplo-3-tests-unitarios)

### Tutorial 4: Migrar Código Existente

1. Lee [BOOTSTRAP_MIGRATION_GUIDE.md - ¿Por Qué Migrar?](BOOTSTRAP_MIGRATION_GUIDE.md#por-qué-migrar)
2. Sigue [BOOTSTRAP_MIGRATION_GUIDE.md - Migración Paso a Paso](BOOTSTRAP_MIGRATION_GUIDE.md#migración-paso-a-paso)
3. Revisa ejemplos en [README.md - Guía de Migración](../README.md#guía-de-migración-al-sistema-de-bootstrap)

## ❓ FAQ

### ¿Dónde empiezo?

Empieza con [BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md) para entender los conceptos básicos.

### ¿Cómo migro mi código existente?

Sigue [BOOTSTRAP_MIGRATION_GUIDE.md](BOOTSTRAP_MIGRATION_GUIDE.md) paso a paso.

### ¿Cómo escribo tests?

Lee [internal/bootstrap/INTEGRATION_TESTS.md](../internal/bootstrap/INTEGRATION_TESTS.md) y revisa los ejemplos en [BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md).

### ¿Cómo desarrollo sin RabbitMQ o S3?

Configura recursos opcionales según [config/OPTIONAL_RESOURCES.md](../config/OPTIONAL_RESOURCES.md).

### ¿Dónde está el código fuente?

En `internal/bootstrap/`. Ver [Código Fuente](#-código-fuente) arriba.

### ¿Hay ejemplos de código?

Sí, muchos en [BOOTSTRAP_USAGE.md - Ejemplos Prácticos](BOOTSTRAP_USAGE.md#ejemplos-prácticos).

## 🔗 Enlaces Rápidos

- **Documentación Principal**: [BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md)
- **Migración**: [BOOTSTRAP_MIGRATION_GUIDE.md](BOOTSTRAP_MIGRATION_GUIDE.md)
- **Testing**: [../internal/bootstrap/INTEGRATION_TESTS.md](../internal/bootstrap/INTEGRATION_TESTS.md)
- **Recursos Opcionales**: [../config/OPTIONAL_RESOURCES.md](../config/OPTIONAL_RESOURCES.md)
- **Quick Start**: [../QUICKSTART.md](../QUICKSTART.md)
- **README**: [../README.md](../README.md)

## 📧 Soporte

Si tienes preguntas o encuentras problemas:

1. Revisa la sección de [Troubleshooting](BOOTSTRAP_MIGRATION_GUIDE.md#troubleshooting) en la guía de migración
2. Revisa los ejemplos en [BOOTSTRAP_USAGE.md](BOOTSTRAP_USAGE.md)
3. Consulta con el equipo de desarrollo
