# Tests de Integración - API Mobile

Tests con testcontainers (PostgreSQL + MongoDB + RabbitMQ).

---

## 🎛️ Control de Ejecución

### **Variables de Entorno**

Los tests de integración se pueden habilitar/deshabilitar con variables de entorno:

| Variable | Valores | Descripción |
|----------|---------|-------------|
| `RUN_INTEGRATION_TESTS` | `true`/`false` | Control principal (recomendado) |
| `INTEGRATION_TESTS` | `1`/`0`/`true`/`false` | Formato alternativo |
| `SKIP_INTEGRATION_TESTS` | `true` | Forzar skip (útil en CI) |
| `CI` | `true` | En CI se ejecutan por defecto |

---

## 🚀 Ejecutar Tests

### **Opción 1: Con Makefile (Recomendado)**

```bash
# ✅ Tests HABILITADOS (se ejecutan)
make test-integration

# ⏭️ Tests DESHABILITADOS (se skipean)
make test-integration-skip

# 📊 Con coverage
make test-integration-coverage

# 🐳 Verificar Docker primero
make docker-check
```

### **Opción 2: Directo con Go**

```bash
# ✅ Habilitar y ejecutar
RUN_INTEGRATION_TESTS=true go test -v -tags=integration ./test/integration/...

# ⏭️ Deshabilitar (skip)
RUN_INTEGRATION_TESTS=false go test -v -tags=integration ./test/integration/...

# 🔥 Sin variable (skip por defecto en local)
go test -v -tags=integration ./test/integration/...
```

### **Opción 3: Exportar Variable**

```bash
# Habilitar para toda la sesión
export RUN_INTEGRATION_TESTS=true

# Ahora todos los comandos ejecutarán tests
go test -v -tags=integration ./test/integration/...
make test-integration

# Deshabilitar
export RUN_INTEGRATION_TESTS=false
```

---

## 🎯 Casos de Uso

### **Desarrollo Local**

```bash
# Quiero correr tests de integración HOY
make test-integration

# NO quiero correrlos (problema con Docker, sin tiempo, etc)
make test-integration-skip
# O simplemente no ejecutes el comando
```

### **CI/CD**

```bash
# GitHub Actions - Siempre habilitados por defecto
CI=true go test -tags=integration ./test/integration/...

# Pero si hay problemas, deshabilitarlos temporalmente
SKIP_INTEGRATION_TESTS=true go test -tags=integration ./test/integration/...
```

### **Pre-commit Hook**

```bash
# Solo correr si Docker está disponible
if docker ps > /dev/null 2>&1; then
  RUN_INTEGRATION_TESTS=true go test -tags=integration ./test/integration/...
else
  echo "⏭️  Docker no disponible, skipping integration tests"
fi
```

---

## 📋 Comportamiento por Defecto

| Contexto | Variable | Resultado |
|----------|----------|-----------|
| **Local sin variable** | - | ⏭️ SKIP (no corren) |
| **Local con RUN_INTEGRATION_TESTS=true** | `true` | ✅ CORREN |
| **Local con RUN_INTEGRATION_TESTS=false** | `false` | ⏭️ SKIP |
| **CI sin variable** | `CI=true` | ✅ CORREN |
| **CI con SKIP_INTEGRATION_TESTS=true** | ambas | ⏭️ SKIP |

---

## 🛠️ Implementación en Tests

Cada test debe incluir al inicio:

```go
func TestAuthFlow(t *testing.T) {
    integration.SkipIfIntegrationTestsDisabled(t)
    // ... resto del test
}
```

Esto verifica automáticamente las variables y skipea si es necesario.

---

## ⚠️ Requisitos

- **Docker**: Debe estar corriendo
- **Testcontainers**: Instalado automáticamente con `go mod download`
- **Build tag**: Usar `-tags=integration`

**Verificar Docker**:
```bash
make docker-check
# O manualmente:
docker ps
```

---

**Nota**: Los tests de integración pueden tardar 1-2 minutos debido a que levantan contenedores reales.
