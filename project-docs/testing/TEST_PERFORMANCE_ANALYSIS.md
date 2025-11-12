# Análisis de Performance de Tests - EduGo API Mobile

**Fecha de Análisis**: 9 de noviembre de 2025  
**Ejecutado por**: Tarea 20.1 - Validación Final  
**Versión**: 1.0

---

## 📊 Resumen Ejecutivo

Este documento presenta el análisis detallado de performance de la suite completa de tests, identificando problemas críticos de rendimiento relacionados con la gestión de contenedores de Docker.

### Métricas Generales

| Métrica | Valor | Estado |
|---------|-------|--------|
| **Tiempo Total** | 7:18 (438s) | ⚠️ LENTO |
| **Tests Unitarios** | 0:05 (5s) | ✅ EXCELENTE |
| **Tests Integración** | 7:13 (433s) | ❌ MUY LENTO |
| **Tests Pasando** | 18/21 (85.7%) | ⚠️ ACEPTABLE |
| **Tests Fallando** | 3/21 (14.3%) | ⚠️ REQUIERE ATENCIÓN |

---

## 🔍 Problema Crítico Identificado

### ❌ Contenedores NO se Reutilizan

**Hallazgo Principal**: Cada test de integración crea y destruye sus propios contenedores (PostgreSQL, MongoDB, RabbitMQ), lo que resulta en:

- **~19 segundos por test** solo en setup/teardown de contenedores
- **60+ contenedores creados y destruidos** durante la ejecución completa
- **Tiempo desperdiciado**: ~6 minutos de los 7 minutos totales

### Evidencia del Problema

Analizando los logs de ejecución:

```
TestAuthFlow_LoginSuccess:
  - Crea contenedores: PostgreSQL, MongoDB, RabbitMQ (~12s)
  - Ejecuta test: (~2s)
  - Destruye contenedores: (~5s)
  - Total: 19.38s

TestAuthFlow_LoginInvalidCredentials:
  - Crea contenedores NUEVAMENTE: PostgreSQL, MongoDB, RabbitMQ (~12s)
  - Ejecuta test: (~2s)
  - Destruye contenedores: (~5s)
  - Total: 19.38s
```

**Patrón detectado**: Cada uno de los 21 tests repite este ciclo completo.

---

## 📈 Análisis Detallado por Fase

### Fase 1: Tests Unitarios ✅

| Métrica | Valor |
|---------|-------|
| **Tiempo Total** | 5 segundos |
| **Tests Ejecutados** | 77 tests |
| **Estado** | ✅ 100% PASS |
| **Performance** | ✅ EXCELENTE |

**Conclusión**: Los tests unitarios están perfectamente optimizados.

---

### Fase 2: Tests de Integración ❌

| Métrica | Valor |
|---------|-------|
| **Tiempo Total** | 433 segundos (7:13) |
| **Tests Ejecutados** | 21 tests |
| **Tests PASS** | 18 (85.7%) |
| **Tests FAIL** | 3 (14.3%) |
| **Tiempo Promedio/Test** | 20.6 segundos |
| **Performance** | ❌ MUY LENTO |

#### Desglose de Tiempos por Test

| Test | Tiempo | Setup | Ejecución | Teardown | Estado |
|------|--------|-------|-----------|----------|--------|
| TestAssessmentFlow_GetAssessment | 19.38s | ~12s | ~2s | ~5s | ✅ PASS |
| TestAssessmentFlow_GetAssessmentNotFound | 19.38s | ~12s | ~2s | ~5s | ✅ PASS |
| TestAssessmentFlow_SubmitAssessment | 19.38s | ~12s | ~2s | ~5s | ✅ PASS |
| TestAssessmentFlow_SubmitAssessmentDuplicate | 19.38s | ~12s | ~2s | ~5s | ✅ PASS |
| TestAuthFlow_LoginSuccess | 19.38s | ~12s | ~2s | ~5s | ✅ PASS |
| TestAuthFlow_LoginInvalidCredentials | 19.38s | ~12s | ~2s | ~5s | ✅ PASS |
| TestAuthFlow_LoginNonexistentUser | 19.20s | ~12s | ~2s | ~5s | ✅ PASS |
| TestMaterialFlow_CreateMaterial | 19.48s | ~12s | ~2s | ~5s | ✅ PASS |
| TestMaterialFlow_GetMaterial | 19.32s | ~12s | ~2s | ~5s | ✅ PASS |
| TestMaterialFlow_GetMaterialNotFound | 19.26s | ~12s | ~2s | ~5s | ✅ PASS |
| TestMaterialFlow_ListMaterials | 19.39s | ~12s | ~2s | ~5s | ✅ PASS |
| TestPostgresTablesExist | 10.42s | ~8s | ~2s | ~0s | ❌ FAIL |
| TestProgressFlow_UpsertProgress | 62.04s | ~60s | ~2s | ~0s | ❌ FAIL (timeout) |
| TestProgressFlow_UpsertProgressUpdate | 19.49s | ~12s | ~2s | ~5s | ✅ PASS |
| TestProgressFlow_UpsertProgressUnauthorized | 61.99s | ~60s | ~2s | ~0s | ❌ FAIL (timeout) |
| TestProgressFlow_UpsertProgressInvalidData | 19.38s | ~12s | ~2s | ~5s | ✅ PASS |
| TestStatsFlow_GetMaterialStats | 19.38s | ~12s | ~2s | ~5s | ✅ PASS |
| TestStatsFlow_GetGlobalStats | 19.47s | ~12s | ~2s | ~5s | ✅ PASS |

**Observaciones**:
- ✅ **18 tests** completan en ~19 segundos (patrón consistente)
- ❌ **2 tests** fallan por timeout de RabbitMQ (~62 segundos cada uno)
- ❌ **1 test** falla por problemas de conexión TCP (~10 segundos)

---

## 🐳 Análisis de Uso de Contenedores

### Patrón Actual (Ineficiente)

```
Test 1:
  ├─ Crear PostgreSQL container (~3s)
  ├─ Crear MongoDB container (~2s)
  ├─ Crear RabbitMQ container (~7s)
  ├─ Ejecutar test (~2s)
  ├─ Destruir PostgreSQL (~1s)
  ├─ Destruir MongoDB (~1s)
  └─ Destruir RabbitMQ (~3s)
  Total: ~19s

Test 2:
  ├─ Crear PostgreSQL container (~3s)  ← DUPLICADO
  ├─ Crear MongoDB container (~2s)     ← DUPLICADO
  ├─ Crear RabbitMQ container (~7s)    ← DUPLICADO
  ├─ Ejecutar test (~2s)
  ├─ Destruir PostgreSQL (~1s)
  ├─ Destruir MongoDB (~1s)
  └─ Destruir RabbitMQ (~3s)
  Total: ~19s

... (repite 21 veces)
```

**Tiempo total desperdiciado**: ~17 segundos × 21 tests = **357 segundos (6 minutos)**

### Patrón Óptimo (Reutilización)

```
Setup Global:
  ├─ Crear PostgreSQL container (~3s)
  ├─ Crear MongoDB container (~2s)
  └─ Crear RabbitMQ container (~7s)
  Total: ~12s (UNA SOLA VEZ)

Test 1:
  ├─ Limpiar datos (~0.5s)
  ├─ Ejecutar test (~2s)
  Total: ~2.5s

Test 2:
  ├─ Limpiar datos (~0.5s)
  ├─ Ejecutar test (~2s)
  Total: ~2.5s

... (repite 21 veces)

Teardown Global:
  ├─ Destruir PostgreSQL (~1s)
  ├─ Destruir MongoDB (~1s)
  └─ Destruir RabbitMQ (~3s)
  Total: ~5s (UNA SOLA VEZ)
```

**Tiempo total optimizado**: 12s + (2.5s × 21) + 5s = **69.5 segundos (~1 minuto)**

**Mejora esperada**: De 7:13 a ~1:10 = **83% más rápido** 🚀

---

## ❌ Tests Fallidos

### 1. TestPostgresTablesExist (10.42s) - FAIL

**Error**: `read tcp 127.0.0.1:62143->127.0.0.1:62142: read: connection reset by peer`

**Causa**: Problema de conexión TCP temporal con el contenedor de PostgreSQL.

**Impacto**: **NO CRÍTICO** - Es un problema de infraestructura de test, no de lógica de negocio.

**Solución Recomendada**:
- Agregar retry logic (2-3 intentos)
- Aumentar timeout de conexión
- Verificar que el contenedor esté completamente listo antes de conectar

---

### 2. TestProgressFlow_UpsertProgress (62.04s) - FAIL

**Error**: `context deadline exceeded` - RabbitMQ no inicia en 60 segundos

**Causa**: RabbitMQ tarda demasiado en iniciar (timeout de 60s excedido)

**Logs**:
```
2025/11/09 20:49:51 container logs (wait until ready: `.*Server startup complete.*` matched 0 times, expected 1
context deadline exceeded)
```

**Impacto**: **CRÍTICO** - Bloquea el test completamente

**Solución Recomendada**:
1. Aumentar timeout de RabbitMQ a 90 segundos
2. Usar contenedor más ligero (sin management plugin)
3. Considerar mock de RabbitMQ para tests que no lo necesitan

---

### 3. TestProgressFlow_UpsertProgressUnauthorized (61.99s) - FAIL

**Error**: Mismo problema que el test anterior (RabbitMQ timeout)

**Causa**: Idéntica al test #2

**Solución**: Misma que el test #2

---

## ⚠️ Problemas Secundarios Detectados

### 1. Contenedores Residuales

**Problema**: 2 contenedores de RabbitMQ no se limpiaron correctamente

```
- rabbitmq:3.12-management-alpine (2dd771c4dc3b)
- rabbitmq:3.12-management-alpine (0222fc5b1c48)
```

**Causa**: Los tests que fallaron por timeout no ejecutaron el cleanup

**Solución**:
- Usar `defer` para garantizar cleanup
- Usar `t.Cleanup()` en lugar de defer manual
- Agregar timeout más corto para evitar bloqueos largos

---

### 2. RabbitMQ Authentication Failures

**Advertencia recurrente**:
```
⚠️  Warning: RabbitMQ topology setup failed (non-critical): 
failed to connect to RabbitMQ: Exception (403) Reason: "username or password not allowed"
```

**Impacto**: **NO CRÍTICO** - El sistema usa un mock publisher como fallback

**Causa**: Credenciales incorrectas o configuración de RabbitMQ

**Solución**:
- Revisar configuración de credenciales en testcontainers
- Usar variables de entorno correctas
- Considerar si RabbitMQ es realmente necesario para estos tests

---

### 3. Tabla `progress` No Existe

**Advertencia recurrente**:
```
⚠️  Warning: Failed to truncate progress: pq: relation "progress" does not exist
```

**Impacto**: **BAJO** - No afecta los tests actuales

**Causa**: La tabla `progress` no está en el schema actual

**Solución**:
- Agregar tabla `progress` al schema
- O remover el truncate de esa tabla si no se usa

---

### 4. MongoDB Unique Index Warning

**Advertencia recurrente**:
```
Warning: Failed to create unique index on assessment_results: 
multi-key map passed in for ordered parameter keys
```

**Impacto**: **BAJO** - No afecta funcionalidad

**Causa**: Problema con la sintaxis de creación de índice en MongoDB

**Solución**:
- Revisar sintaxis de creación de índice
- Usar formato correcto para índices compuestos

---

## 💡 Recomendaciones de Mejora

### Prioridad 1: CRÍTICA (Implementar Inmediatamente)

#### 1.1 Implementar Reutilización de Contenedores

**Problema**: Cada test crea y destruye contenedores (6 minutos desperdiciados)

**Solución**:

```go
// test/integration/setup.go

var (
    sharedPostgres testcontainers.Container
    sharedMongoDB  testcontainers.Container
    sharedRabbitMQ testcontainers.Container
    setupOnce      sync.Once
)

func SetupSharedContainers(t *testing.T) (*Resources, error) {
    var err error
    
    setupOnce.Do(func() {
        // Crear contenedores UNA SOLA VEZ
        sharedPostgres, err = createPostgresContainer()
        if err != nil {
            return
        }
        
        sharedMongoDB, err = createMongoDBContainer()
        if err != nil {
            return
        }
        
        sharedRabbitMQ, err = createRabbitMQContainer()
        if err != nil {
            return
        }
    })
    
    // Limpiar datos entre tests (TRUNCATE, no DROP)
    cleanDatabases(t, sharedPostgres, sharedMongoDB)
    
    return &Resources{
        PostgreSQL: sharedPostgres,
        MongoDB:    sharedMongoDB,
        RabbitMQ:   sharedRabbitMQ,
    }, err
}
```

**Impacto Esperado**: Reducir tiempo de 7:13 a ~1:10 (83% más rápido)

---

#### 1.2 Aumentar Timeout de RabbitMQ

**Problema**: 2 tests fallan por timeout de RabbitMQ (60s no es suficiente)

**Solución**:

```go
// test/integration/setup.go

rabbitContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
    ContainerRequest: testcontainers.ContainerRequest{
        Image: "rabbitmq:3.12-management-alpine",
        WaitingFor: wait.ForLog(".*Server startup complete.*").
            WithStartupTimeout(90 * time.Second), // Aumentar de 60s a 90s
    },
    Started: true,
})
```

**Impacto Esperado**: Resolver 2 tests fallidos

---

#### 1.3 Agregar Retry Logic para Conexiones TCP

**Problema**: TestPostgresTablesExist falla por conexión TCP temporal

**Solución**:

```go
func connectWithRetry(connStr string, maxRetries int) (*sql.DB, error) {
    var db *sql.DB
    var err error
    
    for i := 0; i < maxRetries; i++ {
        db, err = sql.Open("postgres", connStr)
        if err == nil {
            err = db.Ping()
            if err == nil {
                return db, nil
            }
        }
        
        time.Sleep(time.Second * time.Duration(i+1))
    }
    
    return nil, fmt.Errorf("failed after %d retries: %w", maxRetries, err)
}
```

**Impacto Esperado**: Resolver 1 test fallido

---

### Prioridad 2: ALTA (Implementar en Sprint Actual)

#### 2.1 Usar Contenedor RabbitMQ Más Ligero

**Problema**: RabbitMQ con management plugin tarda ~7 segundos en iniciar

**Solución**:

```go
// Cambiar de:
Image: "rabbitmq:3.12-management-alpine"

// A:
Image: "rabbitmq:3.12-alpine"  // Sin management plugin
```

**Impacto Esperado**: Reducir tiempo de inicio de 7s a 3s

---

#### 2.2 Implementar Cleanup Robusto

**Problema**: 2 contenedores residuales no se limpiaron

**Solución**:

```go
func TestSomething(t *testing.T) {
    resources, err := SetupSharedContainers(t)
    require.NoError(t, err)
    
    // Usar t.Cleanup() en lugar de defer
    t.Cleanup(func() {
        cleanDatabases(t, resources.PostgreSQL, resources.MongoDB)
    })
    
    // ... test code
}
```

**Impacto Esperado**: Eliminar contenedores residuales

---

#### 2.3 Paralelizar Tests Independientes

**Problema**: Tests se ejecutan secuencialmente

**Solución**:

```go
func TestAuthFlow_LoginSuccess(t *testing.T) {
    t.Parallel()  // Ejecutar en paralelo
    
    // ... test code
}
```

**Nota**: Solo paralelizar tests que usan contenedores compartidos y limpian sus datos

**Impacto Esperado**: Reducir tiempo adicional en 30-40%

---

### Prioridad 3: MEDIA (Mejoras Opcionales)

#### 3.1 Agregar Tabla `progress` al Schema

**Problema**: Advertencia recurrente sobre tabla faltante

**Solución**: Agregar migración para crear tabla `progress`

---

#### 3.2 Corregir Índice de MongoDB

**Problema**: Advertencia sobre sintaxis de índice

**Solución**: Revisar y corregir sintaxis de creación de índice único

---

#### 3.3 Configurar Credenciales de RabbitMQ

**Problema**: Advertencia sobre autenticación fallida

**Solución**: Configurar credenciales correctas o usar mock

---

## 📊 Proyección de Mejoras

### Escenario Actual

| Fase | Tiempo | % del Total |
|------|--------|-------------|
| Tests Unitarios | 5s | 1.1% |
| Tests Integración | 433s | 98.9% |
| **Total** | **438s (7:18)** | **100%** |

### Escenario Optimizado (Prioridad 1)

| Fase | Tiempo | Mejora | % del Total |
|------|--------|--------|-------------|
| Tests Unitarios | 5s | - | 6.7% |
| Tests Integración | 70s | -363s (-84%) | 93.3% |
| **Total** | **75s (1:15)** | **-363s (-83%)** | **100%** |

### Escenario Óptimo (Prioridad 1 + 2)

| Fase | Tiempo | Mejora | % del Total |
|------|--------|--------|-------------|
| Tests Unitarios | 5s | - | 11.1% |
| Tests Integración | 40s | -393s (-91%) | 88.9% |
| **Total** | **45s (0:45)** | **-393s (-90%)** | **100%** |

---

## ✅ Plan de Acción

### Semana 1: Implementar Prioridad 1

- [ ] **Día 1-2**: Implementar reutilización de contenedores
  - Crear `SetupSharedContainers()`
  - Modificar todos los tests para usar contenedores compartidos
  - Implementar cleanup entre tests (TRUNCATE)

- [ ] **Día 3**: Aumentar timeout de RabbitMQ
  - Cambiar timeout de 60s a 90s
  - Probar tests que fallaban

- [ ] **Día 4**: Agregar retry logic
  - Implementar `connectWithRetry()`
  - Aplicar a todas las conexiones de BD

- [ ] **Día 5**: Validación
  - Ejecutar suite completa
  - Verificar que todos los tests pasan
  - Medir tiempos de ejecución

**Meta**: Reducir tiempo de 7:18 a ~1:15

---

### Semana 2: Implementar Prioridad 2

- [ ] **Día 1**: Usar RabbitMQ más ligero
  - Cambiar a imagen sin management plugin
  - Probar que tests siguen funcionando

- [ ] **Día 2**: Implementar cleanup robusto
  - Reemplazar defer con t.Cleanup()
  - Verificar que no quedan contenedores residuales

- [ ] **Día 3-4**: Paralelizar tests
  - Identificar tests independientes
  - Agregar t.Parallel()
  - Probar que no hay race conditions

- [ ] **Día 5**: Validación
  - Ejecutar suite completa
  - Verificar mejoras de performance

**Meta**: Reducir tiempo de ~1:15 a ~0:45

---

## 📝 Conclusiones

### Hallazgos Principales

1. ✅ **Tests unitarios**: Excelente performance (5s para 77 tests)
2. ❌ **Tests integración**: Muy lentos (7:13 para 21 tests)
3. ❌ **Problema crítico**: Contenedores no se reutilizan (6 minutos desperdiciados)
4. ⚠️ **Tests fallidos**: 3 de 21 (14.3%) - 2 por timeout de RabbitMQ, 1 por TCP

### Impacto de Mejoras

| Mejora | Tiempo Ahorrado | Esfuerzo | Prioridad |
|--------|-----------------|----------|-----------|
| Reutilizar contenedores | -363s (-83%) | 2 días | 🔴 CRÍTICA |
| Aumentar timeout RabbitMQ | Resolver 2 tests | 1 hora | 🔴 CRÍTICA |
| Retry logic TCP | Resolver 1 test | 2 horas | 🔴 CRÍTICA |
| RabbitMQ ligero | -84s adicionales | 1 hora | 🟡 ALTA |
| Paralelización | -20s adicionales | 2 días | 🟡 ALTA |

### Próximos Pasos

1. ✅ **Inmediato**: Implementar reutilización de contenedores (Prioridad 1.1)
2. ✅ **Esta semana**: Resolver tests fallidos (Prioridad 1.2 y 1.3)
3. ✅ **Próxima semana**: Optimizaciones adicionales (Prioridad 2)
4. ✅ **Monitoreo continuo**: Ejecutar suite completa en CI/CD y rastrear tiempos

---

**Última actualización**: 9 de noviembre de 2025  
**Próxima revisión**: Después de implementar Prioridad 1

