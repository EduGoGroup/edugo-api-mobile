# Dependencias del Sprint 02 - Capa de Dominio

## Dependencias Técnicas Previas

### Go 1.21+
- [ ] **Go 1.21 o superior** instalado
- [ ] Soporte para generics (usado en value objects)
- [ ] Soporte para error wrapping

```bash
# Verificar versión de Go
go version
# Output esperado: go version go1.21.x o superior

# Verificar que GOPATH está configurado
echo $GOPATH
go env GOPATH
```

### Sprint 01 Completado
- [ ] **Schema PostgreSQL** creado (migración 06_assessments.sql ejecutada)
- [ ] Tablas: `assessment`, `assessment_attempt`, `assessment_attempt_answer` existen
- [ ] Función `gen_uuid_v7()` disponible

```bash
# Verificar que Sprint-01 está completado
psql -U postgres -d edugo_test -c "\dt assessment*"
# Output esperado: 3-4 tablas (assessment, assessment_attempt, assessment_attempt_answer, material_summary_link)
```

### Estructura Clean Architecture
- [ ] Directorios base creados en edugo-api-mobile

```bash
# Crear estructura si no existe
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

mkdir -p internal/domain/entities
mkdir -p internal/domain/valueobjects
mkdir -p internal/domain/repositories
mkdir -p internal/domain/errors
mkdir -p internal/application/services
mkdir -p internal/infrastructure/persistence

echo "✓ Estructura Clean Architecture creada"
```

---

## Dependencias de Código

### Packages Go Requeridos

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# UUID generation
go get github.com/google/uuid@latest

# Testing framework
go get github.com/stretchr/testify@v1.8.4

# Verificar instalación
go mod tidy
go mod verify
```

**Packages instalados:**
- `github.com/google/uuid` - Generación de UUIDs v4
- `github.com/stretchr/testify/assert` - Assertions en tests
- `github.com/stretchr/testify/require` - Assertions con early exit

### Verificar go.mod

```bash
# El archivo go.mod debe contener mínimo:
cat go.mod | grep -E "(google/uuid|testify)"

# Output esperado:
# github.com/google/uuid vX.X.X
# github.com/stretchr/testify v1.8.4
```

---

## Herramientas de Desarrollo

### golangci-lint (Opcional pero Recomendado)
```bash
# Instalar golangci-lint
brew install golangci-lint  # macOS
# O
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Verificar instalación
golangci-lint --version
# Output esperado: golangci-lint has version X.X.X
```

### gofmt (Incluido con Go)
```bash
# Verificar que gofmt está disponible
which gofmt
# Output esperado: /usr/local/go/bin/gofmt

# Formatear código
gofmt -w ./internal/domain/
```

### go test (Incluido con Go)
```bash
# Verificar que go test funciona
go test -h | head -5
# Output esperado: usage: go test [build/test flags] [packages] [build/test flags & test binary flags]
```

---

## Variables de Entorno

**Nota:** Sprint 02 (capa de dominio) NO requiere variables de entorno porque no tiene dependencias externas. Las variables se necesitarán en Sprint 03 (repositorios).

Para tests de dominio:
```bash
# Opcional: ejecutar tests en modo verbose
export GO_TEST_VERBOSE=1

# Opcional: coverage threshold
export GO_TEST_COVERAGE_THRESHOLD=90
```

---

## Dependencias entre Archivos del Sprint

### Orden de Implementación Recomendado

1. **Primero:** `internal/domain/errors/errors.go`
   - Requerido por todas las entities

2. **Segundo:** Entity Answer
   - Más simple, sin dependencias
   - Requerido por: Entity Attempt

3. **Tercero:** Entity Assessment
   - Sin dependencias de otras entities
   - Puede implementarse en paralelo con Answer

4. **Cuarto:** Entity Attempt
   - **Requiere:** Entity Answer (usa slice de Answer)
   - Depende de Answer para calcular score

5. **Quinto:** Value Objects
   - Opcional, pueden implementarse en paralelo
   - No tienen dependencias entre sí

6. **Sexto:** Repository Interfaces
   - **Requieren:** Entities completas (usan los tipos)
   - Solo definen interfaces, no implementación

```bash
# Verificar que errores están definidos antes de crear entities
ls -la internal/domain/errors/errors.go
# Si no existe, crearlo primero

# Verificar que Answer existe antes de crear Attempt
ls -la internal/domain/entities/answer.go
# Si no existe, crearlo primero
```

---

## Verificación de Dependencias

### Script de Verificación Automática

Crear archivo: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/verify_sprint02_deps.sh`

```bash
#!/bin/bash
# verify_sprint02_deps.sh - Verifica dependencias de Sprint 02

set -e

echo "=== VERIFICANDO DEPENDENCIAS SPRINT 02 ==="

# 1. Verificar Go version
echo ""
echo "1. Verificando Go..."
go_version=$(go version | awk '{print $3}' | sed 's/go//')
min_version="1.21"

if [ "$(printf '%s\n' "$min_version" "$go_version" | sort -V | head -n1)" = "$min_version" ]; then
    echo "✓ Go $go_version >= $min_version"
else
    echo "✗ Go $go_version < $min_version (requerido: $min_version+)"
    exit 1
fi

# 2. Verificar estructura de directorios
echo ""
echo "2. Verificando estructura de directorios..."
for dir in internal/domain/entities internal/domain/valueobjects internal/domain/repositories internal/domain/errors; do
    if [ -d "$dir" ]; then
        echo "✓ $dir existe"
    else
        echo "✗ $dir no existe - creando..."
        mkdir -p "$dir"
    fi
done

# 3. Verificar dependencias Go
echo ""
echo "3. Verificando dependencias Go..."
if go list -m github.com/google/uuid > /dev/null 2>&1; then
    echo "✓ github.com/google/uuid instalado"
else
    echo "✗ github.com/google/uuid NO instalado"
    echo "  Ejecutar: go get github.com/google/uuid@latest"
    exit 1
fi

if go list -m github.com/stretchr/testify > /dev/null 2>&1; then
    echo "✓ github.com/stretchr/testify instalado"
else
    echo "✗ github.com/stretchr/testify NO instalado"
    echo "  Ejecutar: go get github.com/stretchr/testify@v1.8.4"
    exit 1
fi

# 4. Verificar Sprint-01 (schema PostgreSQL)
echo ""
echo "4. Verificando Sprint-01 completado..."
if psql -U postgres -d edugo_test -c "\dt assessment" > /dev/null 2>&1; then
    echo "✓ Schema PostgreSQL existe (Sprint-01 completado)"
else
    echo "⚠ Schema PostgreSQL no encontrado (ejecutar Sprint-01 primero)"
    echo "  Nota: Esto no bloquea Sprint-02 (dominio puro), pero es recomendado"
fi

# 5. Verificar herramientas opcionales
echo ""
echo "5. Verificando herramientas opcionales..."
if command -v golangci-lint > /dev/null 2>&1; then
    echo "✓ golangci-lint instalado"
else
    echo "⚠ golangci-lint no instalado (opcional pero recomendado)"
fi

if command -v gofmt > /dev/null 2>&1; then
    echo "✓ gofmt disponible"
else
    echo "✗ gofmt no disponible (debería venir con Go)"
fi

echo ""
echo "==================================="
echo "✅ TODAS LAS DEPENDENCIAS CRÍTICAS SATISFECHAS"
echo "==================================="
```

**Ejecutar script:**
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile
chmod +x scripts/verify_sprint02_deps.sh
./scripts/verify_sprint02_deps.sh
```

---

## Dependencias de Sprint-01

### Validar que Sprint-01 Está Completo

```bash
# Verificar archivos de Sprint-01
ls -la /Users/jhoanmedina/source/EduGo/Analisys/AnalisisEstandarizado/spec-01-evaluaciones/03-Sprints/Sprint-01-Schema-BD/

# Esperado: 5 archivos
# - README.md
# - TASKS.md
# - DEPENDENCIES.md
# - QUESTIONS.md
# - VALIDATION.md

# Verificar que schema SQL existe
ls -la /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/scripts/postgresql/06_assessments.sql

# Si Sprint-01 no está completo, ejecutarlo primero
```

---

## Troubleshooting

### Error: "package github.com/google/uuid not found"
**Solución:**
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile
go get github.com/google/uuid@latest
go mod tidy
```

### Error: "package github.com/stretchr/testify/assert not found"
**Solución:**
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile
go get github.com/stretchr/testify@v1.8.4
go mod tidy
```

### Error: "go.mod not found"
**Solución:**
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile
go mod init edugo-api-mobile
go mod tidy
```

### Error: "cannot find package internal/domain/errors"
**Solución:**
```bash
# Crear el archivo de errores primero (TASK-02-001 paso 3)
mkdir -p internal/domain/errors
# Luego implementar errors.go según TASKS.md
```

### Warning: "golangci-lint not found"
**Solución (Opcional):**
```bash
# macOS
brew install golangci-lint

# Linux
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Verificar
golangci-lint --version
```

---

## Checklist de Pre-Ejecución

Antes de comenzar Sprint-02, verificar:

- [ ] Go 1.21+ instalado (`go version`)
- [ ] `github.com/google/uuid` instalado (`go list -m github.com/google/uuid`)
- [ ] `github.com/stretchr/testify` instalado (`go list -m github.com/stretchr/testify`)
- [ ] Estructura de directorios creada (`internal/domain/*`)
- [ ] go.mod existe (`ls go.mod`)
- [ ] Sprint-01 completado (opcional pero recomendado)
- [ ] Script de verificación ejecutado y pasó

**Comando único de verificación:**
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile
go version && go mod verify && echo "✓ Listo para Sprint-02"
```

---

## Próximos Sprints (Dependencias Futuras)

**Sprint-03 requerirá:**
- GORM (ORM para PostgreSQL)
- MongoDB driver
- Testcontainers

**Sprint-04 requerirá:**
- Gin framework
- Validator
- Swag (Swagger)

**Nota:** NO instalar dependencias de sprints futuros todavía. Mantener go.mod limpio.

---

**Generado con:** Claude Code  
**Sprint:** 02/06  
**Última actualización:** 2025-11-14
