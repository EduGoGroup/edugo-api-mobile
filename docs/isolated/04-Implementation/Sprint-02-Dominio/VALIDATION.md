# Validación del Sprint 02 - Capa de Dominio

## Pre-validación

### Verificar Estado del Proyecto
```bash
# Cambiar al directorio del proyecto
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Verificar estado de Git
git status

# Verificar rama actual
git branch --show-current
# Output esperado: develop, main, o feature/sprint-02-domain

# Verificar que go.mod existe
ls -la go.mod
```

### Verificar Estructura de Directorios
```bash
# Verificar que estructura de dominio existe
ls -la internal/domain/entities/
ls -la internal/domain/valueobjects/
ls -la internal/domain/repositories/
ls -la internal/domain/errors/

# Esperado: Directorios existen
```

---

## Checklist de Validación

### 1. Tests Unitarios

#### 1.1 Ejecutar Tests de Entities
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Tests de Assessment
go test ./internal/domain/entities -v -run TestAssessment

# Tests de Attempt
go test ./internal/domain/entities -v -run TestAttempt

# Tests de Answer
go test ./internal/domain/entities -v -run TestAnswer

# Todos los tests juntos
go test ./internal/domain/entities -v
```

**Criterio de éxito:** Todos los tests pasan sin errores (0 FAIL)

**Output esperado:**
```
=== RUN   TestNewAssessment_Success
--- PASS: TestNewAssessment_Success (0.00s)
=== RUN   TestNewAssessment_InvalidMaterialID
--- PASS: TestNewAssessment_InvalidMaterialID (0.00s)
...
PASS
ok      edugo-api-mobile/internal/domain/entities    0.XXXs
```

#### 1.2 Ejecutar Tests de Value Objects
```bash
# Tests de todos los value objects
go test ./internal/domain/valueobjects -v
```

**Criterio de éxito:** Todos los tests pasan (PASS)

---

### 2. Coverage (Cobertura de Código)

#### 2.1 Coverage de Entities
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Coverage de entities
go test ./internal/domain/entities -cover -coverprofile=coverage_entities.out

# Ver coverage detallado
go tool cover -func=coverage_entities.out

# Ver solo totales
go tool cover -func=coverage_entities.out | grep total
```

**Criterio de éxito:** Coverage total **>90%**

**Output esperado:**
```
total:    (statements)    92.5%
```

#### 2.2 Coverage de Value Objects
```bash
# Coverage de value objects
go test ./internal/domain/valueobjects -cover -coverprofile=coverage_vo.out

# Ver totales
go tool cover -func=coverage_vo.out | grep total
```

**Criterio de éxito:** Coverage total **>90%**

#### 2.3 Reporte HTML de Coverage
```bash
# Generar reporte visual
go test ./internal/domain/... -coverprofile=coverage_domain.out
go tool cover -html=coverage_domain.out -o coverage_domain.html

# Abrir en navegador
open coverage_domain.html  # macOS
# O: xdg-open coverage_domain.html  # Linux
```

**Criterio de éxito:** 
- Visualmente verificar que líneas críticas están cubiertas (verdes)
- Business rules tienen tests (ej: CanAttempt, IsPassed, cálculo de score)
- Validaciones tienen tests (ej: errores por parámetros inválidos)

---

### 3. Linting (Calidad de Código)

#### 3.1 golangci-lint
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Ejecutar linter en dominio
golangci-lint run ./internal/domain/...

# Con configuración específica (si existe .golangci.yml)
golangci-lint run --config .golangci.yml ./internal/domain/...
```

**Criterio de éxito:** **0 errores, 0 warnings**

**Output esperado:**
```
(vacío - sin output significa sin errores)
```

**Si hay warnings menores:**
```bash
# Permitir algunos warnings pero 0 errores
golangci-lint run ./internal/domain/... | grep -c "Error:"
# Output esperado: 0
```

#### 3.2 gofmt (Formato de Código)
```bash
# Verificar que código está formateado
gofmt -l ./internal/domain/

# Si output vacío = todo está bien formateado
# Si lista archivos = necesitan formateo
```

**Criterio de éxito:** Output vacío (no lista archivos)

**Si necesita formateo:**
```bash
# Formatear automáticamente
gofmt -w ./internal/domain/
```

#### 3.3 go vet (Análisis Estático)
```bash
# Analizar código en busca de errores comunes
go vet ./internal/domain/...
```

**Criterio de éxito:** Sin errores

---

### 4. Build (Compilación)

#### 4.1 Compilar Paquetes de Dominio
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Compilar entities
go build ./internal/domain/entities

# Compilar value objects
go build ./internal/domain/valueobjects

# Compilar repositories (interfaces)
go build ./internal/domain/repositories

# Compilar errors
go build ./internal/domain/errors

# Compilar todo el dominio
go build ./internal/domain/...
```

**Criterio de éxito:** Build exitoso sin errores de compilación

**Output esperado:**
```
(vacío - sin output significa compilación exitosa)
```

#### 4.2 Verificar Dependencias
```bash
# Verificar que go.mod está limpio
go mod tidy

# Verificar módulos
go mod verify
```

**Criterio de éxito:** 
- `go mod tidy` no hace cambios
- `go mod verify` retorna: `all modules verified`

---

### 5. Validación de Reglas de Negocio

#### 5.1 Verificar Business Rules en Tests
```bash
# Buscar tests de business rules específicas
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# ¿Hay tests de CanAttempt?
grep -r "TestAssessment_CanAttempt" ./internal/domain/entities/
# Esperado: Encuentra archivo de test

# ¿Hay tests de cálculo de score?
grep -r "TestNewAttempt_ScoreCalculation\|TestAttempt.*Score" ./internal/domain/entities/
# Esperado: Encuentra tests

# ¿Hay tests de IsPassed?
grep -r "TestAttempt_IsPassed" ./internal/domain/entities/
# Esperado: Encuentra test
```

**Criterio de éxito:** Todos los métodos de business rules tienen tests

#### 5.2 Ejecutar Tests de Casos Límite
```bash
# Ejecutar solo tests que verifican validaciones
go test ./internal/domain/entities -v -run "Invalid|Empty|Negative"

# Verificar que errores de dominio se retornan correctamente
go test ./internal/domain/entities -v -run "Error"
```

**Criterio de éxito:** Tests de validaciones pasan

---

### 6. Verificación de Inmutabilidad

#### 6.1 Verificar que Attempt es Inmutable
```bash
# Buscar setters en Attempt (NO deben existir)
grep -n "func (a \*Attempt) Set" /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/internal/domain/entities/attempt.go

# Output esperado: (vacío) - no debe haber setters
```

**Criterio de éxito:** No se encuentran métodos `Set*` en Attempt

#### 6.2 Verificar Value Objects Inmutables
```bash
# Value objects no deben tener setters
grep -rn "func (.*) Set" ./internal/domain/valueobjects/

# Output esperado: (vacío)
```

**Criterio de éxito:** Value objects sin setters

---

### 7. Verificación de Dependencias Externas

#### 7.1 Verificar que Dominio NO Depende de Frameworks
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Buscar imports no permitidos en dominio
# ❌ NO debe importar: gorm, gin, mongo driver, etc.

grep -r "gorm.io\|gin-gonic\|mongo-driver" ./internal/domain/

# Output esperado: (vacío) - no debe haber imports de frameworks
```

**Criterio de éxito:** Dominio solo importa:
- ✅ Stdlib de Go (time, errors, fmt, context)
- ✅ github.com/google/uuid (IDs)
- ✅ Packages de testing en archivos _test.go

**Imports permitidos:**
```go
import (
    "context"
    "errors"
    "fmt"
    "time"
    
    "github.com/google/uuid"
)

// En tests (_test.go):
import (
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)
```

#### 7.2 Verificar go.mod
```bash
# Ver dependencias directas del módulo
go list -m all | grep -v "indirect"

# Verificar que NO hay dependencias pesadas en dominio
# ✅ Debe tener: google/uuid, testify
# ❌ NO debe tener (todavía): gorm, gin, mongo
```

**Criterio de éxito:** Solo dependencias mínimas necesarias

---

## Criterios de Éxito Globales del Sprint

### Checklist Final

- [ ] **Entities Implementadas**
  - [ ] `assessment.go` con 3+ business rules
  - [ ] `attempt.go` inmutable con cálculo de score
  - [ ] `answer.go` con validaciones básicas

- [ ] **Value Objects Implementados**
  - [ ] Mínimo 5 value objects creados
  - [ ] Todos inmutables (sin setters)
  - [ ] Con método `Equals()` y `String()`

- [ ] **Repository Interfaces Definidas**
  - [ ] `AssessmentRepository` con 4+ métodos
  - [ ] `AttemptRepository` con 5+ métodos
  - [ ] `AnswerRepository` con 2+ métodos

- [ ] **Errores de Dominio**
  - [ ] `errors.go` con 10+ errores sentinel
  - [ ] Nombres con prefijo `Err`
  - [ ] Mensajes con prefijo `"domain:"`

- [ ] **Tests Unitarios**
  - [ ] Coverage >90% en entities
  - [ ] Coverage >90% en value objects
  - [ ] Tests de casos exitosos y fallidos
  - [ ] Tests de business rules
  - [ ] Tests de validaciones

- [ ] **Calidad de Código**
  - [ ] golangci-lint sin errores
  - [ ] gofmt sin archivos pendientes
  - [ ] go vet sin warnings
  - [ ] Build exitoso

- [ ] **Arquitectura Limpia**
  - [ ] Dominio sin dependencias a frameworks
  - [ ] Solo imports permitidos (stdlib, uuid, testing)
  - [ ] Entities con lógica de negocio (no anémicas)

---

## Validación Automatizada

### Script de Validación Completa

Crear archivo: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/validate_sprint02.sh`

```bash
#!/bin/bash
# validate_sprint02.sh - Validación automatizada de Sprint 02

set -e

cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

echo "========================================="
echo "VALIDACIÓN SPRINT 02 - CAPA DE DOMINIO"
echo "========================================="

PASSED=0
FAILED=0

# 1. Tests
echo ""
echo "1. Ejecutando tests..."
if go test ./internal/domain/... -v > /tmp/test_output.txt 2>&1; then
    echo "✅ Tests: PASSED"
    ((PASSED++))
else
    echo "❌ Tests: FAILED"
    cat /tmp/test_output.txt
    ((FAILED++))
fi

# 2. Coverage
echo ""
echo "2. Verificando coverage..."
go test ./internal/domain/entities -cover -coverprofile=/tmp/coverage.out > /dev/null 2>&1
coverage=$(go tool cover -func=/tmp/coverage.out | grep total | awk '{print $3}' | sed 's/%//')

if (( $(echo "$coverage >= 90" | bc -l) )); then
    echo "✅ Coverage: $coverage% (>=90%)"
    ((PASSED++))
else
    echo "❌ Coverage: $coverage% (<90%)"
    ((FAILED++))
fi

# 3. Linting
echo ""
echo "3. Ejecutando linter..."
if golangci-lint run ./internal/domain/... > /dev/null 2>&1; then
    echo "✅ Linting: PASSED"
    ((PASSED++))
else
    echo "❌ Linting: FAILED"
    golangci-lint run ./internal/domain/...
    ((FAILED++))
fi

# 4. Build
echo ""
echo "4. Verificando build..."
if go build ./internal/domain/... > /dev/null 2>&1; then
    echo "✅ Build: PASSED"
    ((PASSED++))
else
    echo "❌ Build: FAILED"
    ((FAILED++))
fi

# 5. Dependencias limpias
echo ""
echo "5. Verificando dependencias..."
if ! grep -r "gorm.io\|gin-gonic\|mongo-driver" ./internal/domain/ > /dev/null 2>&1; then
    echo "✅ Sin dependencias de frameworks"
    ((PASSED++))
else
    echo "❌ Dominio tiene dependencias no permitidas"
    grep -r "gorm.io\|gin-gonic\|mongo-driver" ./internal/domain/
    ((FAILED++))
fi

# Resumen
echo ""
echo "========================================="
echo "RESUMEN DE VALIDACIÓN"
echo "========================================="
echo "✅ Pasadas: $PASSED"
echo "❌ Fallidas: $FAILED"
echo "Total: $((PASSED + FAILED))"

if [ "$FAILED" -eq 0 ]; then
    echo ""
    echo "🎉 SPRINT 02 VALIDADO EXITOSAMENTE"
    exit 0
else
    echo ""
    echo "⚠️  HAY VALIDACIONES FALLIDAS"
    exit 1
fi
```

**Ejecutar validación:**
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile
chmod +x scripts/validate_sprint02.sh
./scripts/validate_sprint02.sh
```

**Criterio de éxito:** Script retorna exit code 0 (todas las validaciones pasan)

---

## Comandos de Rollback

### Si Algo Falla Durante el Sprint

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Ver cambios actuales
git status

# Si necesitas revertir cambios no commiteados
git checkout -- internal/domain/

# Si ya hiciste commit pero quieres volver atrás
git log --oneline -5  # Ver últimos commits
git revert <commit-hash>  # Revertir commit específico

# O crear branch de backup
git checkout -b backup/sprint-02-$(date +%Y%m%d)
git checkout main
git branch -D feature/sprint-02-domain  # Eliminar branch problemática
```

### Restaurar desde Backup
```bash
# Si creaste backup antes de empezar
git checkout backup/pre-sprint-02
git checkout -b feature/sprint-02-domain-retry
```

---

## Métricas de Éxito

Al completar Sprint 02, debes tener:

| Métrica | Objetivo | Comando Verificación |
|---------|----------|---------------------|
| Entities creadas | 3 | `ls internal/domain/entities/*.go \| grep -v test \| wc -l` |
| Value objects | >=5 | `ls internal/domain/valueobjects/*.go \| grep -v test \| wc -l` |
| Repository interfaces | 3 | `ls internal/domain/repositories/*.go \| wc -l` |
| Tests unitarios | Todos pasando | `go test ./internal/domain/... -v` |
| Coverage | >90% | `go test ./internal/domain/... -cover` |
| Errores linter | 0 | `golangci-lint run ./internal/domain/...` |
| Build | Exitoso | `go build ./internal/domain/...` |
| Dependencias limpias | Sí | `grep -r "gorm\|gin" ./internal/domain/` (vacío) |

---

## Reporte de Validación

Al completar validación, crear reporte:

```bash
cat > /Users/jhoanmedina/source/EduGo/Analisys/AnalisisEstandarizado/spec-01-evaluaciones/03-Sprints/Sprint-02-Dominio/VALIDATION_REPORT.md << EOF
# Reporte de Validación - Sprint 02

**Fecha:** $(date +%Y-%m-%d)
**Ejecutado por:** $(whoami)

## Resultados

- ✅ Tests unitarios: PASSED
- ✅ Coverage: $(go test ./internal/domain/entities -cover 2>/dev/null | grep coverage | awk '{print $5}')
- ✅ Linting: PASSED
- ✅ Build: PASSED
- ✅ Arquitectura limpia: PASSED

## Archivos Creados

\`\`\`
$(find internal/domain -name "*.go" | wc -l) archivos Go
$(find internal/domain -name "*_test.go" | wc -l) archivos de test
\`\`\`

## Estado

**SPRINT 02: COMPLETADO ✅**
EOF

cat /Users/jhoanmedina/source/EduGo/Analisys/AnalisisEstandarizado/spec-01-evaluaciones/03-Sprints/Sprint-02-Dominio/VALIDATION_REPORT.md
```

---

**Generado con:** Claude Code  
**Sprint:** 02/06  
**Última actualización:** 2025-11-14
