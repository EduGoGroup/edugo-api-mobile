# Guía de Tests de Integración

## 🎯 Objetivo

Los tests de integración verifican la interacción entre múltiples componentes del sistema, incluyendo bases de datos reales mediante testcontainers.

## 📍 Ubicación

`test/integration/` con build tag `//go:build integration`

## 🐳 Testcontainers

### Setup Automático

```go
//go:build integration

func TestSomething(t *testing.T) {
    SkipIfIntegrationTestsDisabled(t)

    // Setup completo de app con testcontainers
    app, cleanup := SetupTestApp(t)
    defer cleanup()

    // app.Container - Acceso a servicios
    // app.DB - PostgreSQL real
    // app.MongoDB - MongoDB real
}
```

### Servicios Disponibles

- **PostgreSQL 15**: Base de datos relacional
- **MongoDB 7**: Base de datos NoSQL
- **RabbitMQ 3.12**: Sistema de mensajería

## 🌱 Helpers de Seed

### Usuario de Prueba

```go
// Usuario simple (test@edugo.com, password: Test1234!)
userID, email := SeedTestUser(t, db)

// Usuario con email personalizado
userID, email := SeedTestUserWithEmail(t, db, "custom@example.com")

// Usuario con rol específico
userID, email := SeedTestUserWithRole(t, db, "teacher@example.com", "teacher")

// Múltiples usuarios
users := SeedTestUsers(t, db, 5, "student")
// users[0].Email, users[0].Password, users[0].ID
```

### Escenario Completo

```go
// Crear teacher, students, materials, assessments
scenario := SeedCompleteTestScenario(t, db, mongodb, 3)

// Acceder a datos
teacherID := scenario.Teacher.ID
studentID := scenario.Students[0].ID
materialID := scenario.Materials[0]
password := scenario.Teacher.Password // "Test1234!"
```

### Material y Assessment

```go
// Crear material
materialID := SeedTestMaterial(t, db, authorID)

// Material con título personalizado
materialID := SeedTestMaterialWithTitle(t, db, authorID, "Mi Material")

// Crear assessment
assessmentID := SeedTestAssessment(t, mongodb, materialID)
```

## 🧹 Limpieza de Datos

```go
// Limpiar PostgreSQL (entre tests)
CleanDatabase(t, db)

// Limpiar MongoDB
CleanMongoCollections(t, mongodb)
```

## 📝 Ejemplo Completo

```go
//go:build integration

package integration

import (
    "testing"
)

func TestAuthFlow_LoginSuccess(t *testing.T) {
    // 1. Skip si tests deshabilitados
    SkipIfIntegrationTestsDisabled(t)

    // 2. Setup de app con testcontainers
    app, cleanup := SetupTestApp(t)
    defer cleanup()

    // 3. Limpiar datos
    CleanDatabase(t, app.DB)

    // 4. Seed de datos de prueba
    userID, email := SeedTestUser(t, app.DB)

    // 5. Ejecutar test
    router := app.Container.Router()
    w := httptest.NewRecorder()

    reqBody := `{"email":"test@edugo.com","password":"Test1234!"}`
    req := httptest.NewRequest("POST", "/api/v1/auth/login",
        strings.NewReader(reqBody))
    req.Header.Set("Content-Type", "application/json")

    router.ServeHTTP(w, req)

    // 6. Assertions
    assert.Equal(t, 200, w.Code)
    assert.Contains(t, w.Body.String(), "access_token")
}
```

## ⚡ Comandos

```bash
# Ejecutar tests de integración
RUN_INTEGRATION_TESTS=true make test-integration

# Con verbose
make test-integration-verbose

# Verificar Docker
make docker-check
```

## 🔧 Troubleshooting

### Docker no está corriendo

```bash
make docker-check
# Si falla, iniciar Docker Desktop
```

### Tests muy lentos

Los tests de integración son lentos por naturaleza (testcontainers):
- Promedio: 16s por test
- Suite completa: ~5 minutos

Para desarrollo rápido, usar `make test-unit`

### Error "connection reset by peer"

Error temporal de testcontainers. Reintentar el test.

---

**Ver también**:
- [TESTING_UNIT_GUIDE.md](./TESTING_UNIT_GUIDE.md)
- [TESTING_GUIDE.md](./TESTING_GUIDE.md)
