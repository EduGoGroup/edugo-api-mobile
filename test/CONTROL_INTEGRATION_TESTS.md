# 🎛️ Sistema de Control para Tests de Integración

## 🎯 Problema Resuelto

**Antes**: Tests de integración siempre corrían, causando problemas cuando:
- Docker no está disponible
- Hay problemas con testcontainers
- Quieres skipearlos temporalmente en CI

**Ahora**: Sistema flexible con **múltiples niveles de control**.

---

## ⚡ Uso Rápido

### **Comando Rápido (Recomendado)**

```bash
# ✅ Ejecutar tests
make test-integration

# ⏭️ Skipear tests
make test-integration-skip

# 🐳 Verificar Docker primero
make docker-check
```

### **Control Manual**

```bash
# Habilitar
export RUN_INTEGRATION_TESTS=true
go test -tags=integration ./test/integration/...

# Deshabilitar
export RUN_INTEGRATION_TESTS=false
go test -tags=integration ./test/integration/...
```

---

## 🔧 Variables Disponibles

| Variable | Dónde | Efecto |
|----------|-------|--------|
| **`RUN_INTEGRATION_TESTS`** | Local/CI | `true` = corre, `false` = skip |
| **`INTEGRATION_TESTS`** | Local/CI | `1`/`true` = corre |
| **`SKIP_INTEGRATION_TESTS`** | CI | `true` = forzar skip |
| **`CI`** | GitHub Actions | Auto-detectado, corre por defecto |

---

## 🎯 Casos de Uso

### 1. **Desarrollo Local Normal**

```bash
# Quiero correr tests HOY
make test-integration

# NO tengo tiempo / Docker no funciona
make test-integration-skip
# O simplemente no ejecutes nada
```

### 2. **Problema Temporal (Docker caído)**

```bash
# Opción A: Usar comando skip
make test-integration-skip

# Opción B: Deshabilitar para toda la sesión
export RUN_INTEGRATION_TESTS=false
make test  # Ahora todos skipean automáticamente
```

### 3. **CI/CD con Problemas**

**GitHub Actions**:

```yaml
# Opción 1: Variable en workflow
- name: Run tests
  env:
    SKIP_INTEGRATION_TESTS: true  # Forzar skip
  run: go test -tags=integration ./test/integration/...

# Opción 2: Variable en GitHub Settings
# Settings → Variables → ENABLE_INTEGRATION_TESTS = false
jobs:
  test:
    if: vars.ENABLE_INTEGRATION_TESTS == 'true'
```

### 4. **Hook Local Inteligente (Opcional)**

```bash
# .git/hooks/pre-push (ejemplo)
if docker ps > /dev/null 2>&1; then
  echo "✅ Docker disponible, corriendo integration tests"
  RUN_INTEGRATION_TESTS=true make test-integration
else
  echo "⏭️  Docker no disponible, skipping"
fi
```

---

## 📊 Comportamiento por Defecto

| Contexto | Variable | Resultado |
|----------|----------|-----------|
| Local sin variable | - | ⏭️ **SKIP** |
| Local con `RUN_INTEGRATION_TESTS=true` | ✅ | **CORREN** |
| Local con `RUN_INTEGRATION_TESTS=false` | ❌ | **SKIP** |
| CI (GitHub Actions) | `CI=true` | **CORREN** |
| CI con `SKIP_INTEGRATION_TESTS=true` | ❌ | **SKIP** |

**Filosofía**:
- **Local**: Skip por defecto (habilitar explícitamente)
- **CI**: Corre por defecto (deshabilitar solo si es necesario)

---

## 🛠️ Implementación en Tests

```go
// +build integration

package integration

import "testing"

func TestAuthFlow(t *testing.T) {
    // ✅ SIEMPRE incluir esto al inicio
    SkipIfIntegrationTestsDisabled(t)

    // Si llegamos aquí, tests están habilitados
    // ... tu código de test
}
```

---

## 🔄 Cambiar Comportamiento

### **Temporalmente (1 vez)**
```bash
RUN_INTEGRATION_TESTS=true make test-integration
```

### **Para Toda la Sesión**
```bash
export RUN_INTEGRATION_TESTS=true
# Ahora todos los comandos usarán este valor
```

### **Permanentemente (shell config)**
```bash
# En ~/.zshrc o ~/.bashrc
export RUN_INTEGRATION_TESTS=false  # Skip por defecto
```

### **En CI (GitHub Settings)**
```
Settings → Secrets and variables → Actions → Variables
Crear: ENABLE_INTEGRATION_TESTS = false
```

---

## ✅ Validación

```bash
# Test 1: Sin variable (debe skipear)
go test -v -tags=integration ./test/integration/example_test.go
# Esperado: "⏭️  Integration tests disabled"

# Test 2: Con variable (debe correr)
RUN_INTEGRATION_TESTS=true go test -v -tags=integration ./test/integration/example_test.go
# Esperado: "✅ Integration tests están HABILITADOS"

# Test 3: Makefile skip
make test-integration-skip
# Esperado: Tests se skipean automáticamente

# Test 4: Makefile run
make test-integration
# Esperado: Tests corren (si Docker está disponible)
```

---

## 📚 Archivos Relacionados

- **Implementación**: `test/integration/config.go`
- **Documentación**: `test/integration/README.md`
- **Ejemplo CI**: `.github/workflows/integration-tests.yml.example`
- **Makefile**: Ver comandos `test-integration*`

---

## 🎉 Beneficios

1. ✅ **Flexibilidad**: Habilitar/deshabilitar sin cambiar código
2. ✅ **CI robusto**: Skipear temporalmente si hay problemas
3. ✅ **Developer friendly**: No forzar tests si Docker no funciona
4. ✅ **Zero config**: Por defecto skipea en local, corre en CI
5. ✅ **Múltiples niveles**: Variable, Makefile, GitHub Settings

---

**¿Dudas?** Ver: `test/integration/README.md` para más detalles.
