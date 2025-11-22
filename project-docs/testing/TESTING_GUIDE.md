# Guía de Testing - EduGo API Mobile

**Versión**: 1.0  
**Última actualización**: 2025-11-09

---

## 📋 Tabla de Contenidos

1. [Filosofía de Testing](#filosofía-de-testing)
2. [Tipos de Tests](#tipos-de-tests)
3. [Estructura de Testing](#estructura-de-testing)
4. [Comandos Disponibles](#comandos-disponibles)
5. [Mejores Prácticas](#mejores-prácticas)
6. [Referencias](#referencias)

---

## 🎯 Filosofía de Testing

Este proyecto sigue la **pirámide de testing** con énfasis en tests unitarios rápidos y confiables:

```
            /\
           /  \
          / E2E \ ← 10% (Tests de integración completos)
         /------\
        /        \
       /Integration\ ← 20% (Tests de componentes)
      /------------\
     /              \
    /  Unit Tests   \ ← 70% (Tests unitarios)
   /------------------\
```

### Principios

1. **Tests rápidos**: Tests unitarios < 100ms
2. **Tests aislados**: Cada test es independiente
3. **Tests legibles**: Patrón AAA (Arrange-Act-Assert)
4. **Cobertura inteligente**: Excluir código generado y DTOs
5. **Automatización**: Tests en CI/CD

---

## 🧪 Tipos de Tests

### 1. Tests Unitarios (70%)

**Qué testean**: Funciones puras, lógica de negocio, validaciones

**Ubicación**: Junto al código fuente (`*_test.go` en el mismo paquete)

**Características**:
- ✅ Muy rápidos (< 100ms por test)
- ✅ Usan mocks para dependencias
- ✅ No requieren Docker ni recursos externos
- ✅ Se ejecutan en paralelo

**Ejemplos**:
```go
// internal/domain/valueobject/email_test.go
func TestNewEmail_ValidEmails(t *testing.T) {
    t.Parallel()

    email, err := NewEmail("test@example.com")

    assert.NoError(t, err)
    assert.Equal(t, "test@example.com", email.String())
}
```

**Ejecutar**:
```bash
make test-unit              # Todos los tests unitarios
make test-unit-coverage     # Con reporte de cobertura
make test-watch             # Watch mode (requiere entr)
```

### 2. Tests de Integración (20%)

**Qué testean**: Interacción entre componentes, acceso a bases de datos

**Ubicación**: `test/integration/` con build tag `//go:build integration`

**Características**:
- ⚙️ Medio-lentos (1-5s por test)
- 🐳 Usan testcontainers para BD reales
- 🎯 Testean flujos end-to-end
- 🔧 Requieren Docker corriendo

**Ejemplos**:
```go
// test/integration/auth_flow_test.go
//go:build integration

func TestAuthFlow_LoginSuccess(t *testing.T) {
    SkipIfIntegrationTestsDisabled(t)

    app, cleanup := SetupTestApp(t)
    defer cleanup()

    // ... test de flujo completo
}
```

**Ejecutar**:
```bash
make test-integration          # Todos los tests de integración
make test-integration-verbose  # Con logs detallados
make test-all                  # Unitarios + Integración
```

### 3. Tests End-to-End (10%)

**Qué testean**: Flujos completos de usuario con todos los servicios

**Ubicación**: `test/integration/` (mismo que integración)

**Características**:
- 🐌 Lentos (5-20s por test)
- 🌐 Todos los servicios reales
- 🎬 Escenarios completos de usuario

---

## 📁 Estructura de Testing

```
edugo-api-mobile/
├── internal/
│   ├── domain/
│   │   ├── entity/
│   │   │   ├── user.go
│   │   │   └── user_test.go          ← Tests unitarios
│   │   └── valueobject/
│   │       ├── email.go
│   │       └── email_test.go         ← Tests unitarios
│   ├── application/service/
│   │   ├── auth_service.go
│   │   └── auth_service_test.go      ← Tests unitarios (con mocks)
│   └── infrastructure/
│       └── http/handler/
│           ├── auth_handler.go
│           └── auth_handler_test.go  ← Tests unitarios (con mocks)
│
└── test/
    ├── integration/                  ← Tests de integración
    │   ├── config.go                 ← Control de ejecución
    │   ├── setup.go                  ← Testcontainers setup
    │   ├── testhelpers.go            ← Helpers y seeds
    │   ├── auth_flow_test.go         ← Tests E2E
    │   └── ...
    └── scripts/                      ← Scripts de desarrollo
        ├── setup_dev_env.sh
        └── teardown_dev_env.sh
```

---

## 🚀 Comandos Disponibles

### Testing Básico

```bash
# Tests unitarios (rápido, sin Docker)
make test-unit

# Tests de integración (requiere Docker)
make test-integration

# Todos los tests
make test-all

# Validar que todos pasan
make test-validate
```

### Cobertura

```bash
# Reporte completo con filtrado
make coverage-report

# Verificar umbral mínimo (60%)
make coverage-check

# Reporte de solo tests unitarios
make test-unit-coverage
```

### Análisis

```bash
# Analizar estructura de tests
make test-analyze

# Identificar módulos sin tests
make test-missing
```

### Desarrollo Local

```bash
# Levantar ambiente completo (PostgreSQL, MongoDB, RabbitMQ)
make dev-setup

# Detener ambiente
make dev-teardown

# Resetear ambiente
make dev-reset

# Ver logs
make dev-logs
```

---

## ✨ Mejores Prácticas

### 1. Patrón AAA (Arrange-Act-Assert)

```go
func TestSomething(t *testing.T) {
    // Arrange - Preparar datos y mocks
    input := "test"
    expected := "TEST"

    // Act - Ejecutar función a testear
    result := ToUpper(input)

    // Assert - Verificar resultados
    assert.Equal(t, expected, result)
}
```

### 2. Usar t.Parallel() cuando sea posible

```go
func TestSomething(t *testing.T) {
    t.Parallel() // Ejecutar en paralelo

    // ... resto del test
}
```

### 3. Table-Driven Tests

```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid", "test@example.com", false},
        {"invalid", "invalid", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            err := Validate(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### 4. Usar Mocks Apropiadamente

```go
// Crear mock
mockRepo := new(MockRepository)
mockRepo.On("FindByID", mock.Anything, userID).Return(user, nil)

// Usar en test
service := NewService(mockRepo)
result, err := service.GetUser(userID)

// Verificar que se llamó
mockRepo.AssertExpectations(t)
```

### 5. Cleanup con t.Cleanup()

```go
func TestSomething(t *testing.T) {
    resource := setup()
    t.Cleanup(func() {
        resource.Close()
    })

    // ... test code
}
```

---

## 🔧 Configuración

### Variables de Entorno para Tests

```bash
# Habilitar tests de integración
export RUN_INTEGRATION_TESTS=true

# Nivel de log
export LOG_LEVEL=debug
```

### Exclusiones de Cobertura

El archivo `.coverignore` define qué código excluir:

```
# Archivos generados
docs/
*_mock.go

# DTOs simples
internal/application/dto/

# Entry points
cmd/
```

### Umbral de Cobertura

**Umbral mínimo general**: 60%

**Umbrales por módulo**:
- Services: 70%+
- Domain (ValueObjects, Entities): 80%+
- Handlers: 60%+

---

## 📚 Referencias

- [Guía de Tests Unitarios](./TESTING_UNIT_GUIDE.md)
- [Guía de Tests de Integración](./TESTING_INTEGRATION_GUIDE.md)
- [Reporte de Análisis](./TEST_ANALYSIS_REPORT.md)
- [Plan de Cobertura](./TEST_COVERAGE_PLAN.md)

---

## 💡 Tips y Troubleshooting

### Docker no está corriendo

```bash
# Verificar Docker
make docker-check

# Iniciar Docker Desktop manualmente
```

### Tests de integración lentos

```bash
# Ejecutar solo tests unitarios (más rápidos)
make test-unit

# Skip tests de integración
RUN_INTEGRATION_TESTS=false make test-integration
```

### Ver cobertura en navegador

```bash
make coverage-report
open coverage/coverage.html
```

---

**Generado por**: Sistema de Testing EduGo  
**Mantenido por**: Equipo de Desarrollo
