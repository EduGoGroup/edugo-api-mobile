# 🧹 Análisis de Limpieza - Directorio Raíz

Análisis completo de archivos y carpetas en la raíz del proyecto para identificar qué se puede eliminar, mover o configurar mejor.

**Fecha del análisis**: 11 de noviembre de 2025
**Tamaño total del directorio**: ~316 MB
**Archivos analizados**: 49 items

---

## 📊 Resumen Ejecutivo

| Categoría | Cantidad | Tamaño | Acción |
|-----------|----------|--------|--------|
| **🗑️ Para Eliminar** | 5 archivos | ~110 MB | Borrar archivos temporales/build |
| **📁 Para Mover** | 3 items + 15 archivos .md | ~365 KB | Reorganizar en `/project-docs/` |
| **⚙️ Para Configurar** | 4 items | Variable | Configurar salidas automáticas |
| **✅ Mantener en Raíz** | 37 items | ~206 MB | Son necesarios o correctos |

**⚠️ Hallazgo importante**: `/docs/` contiene 225 KB de documentación del proyecto mezclada con archivos de Swagger. Ver **Apéndice** al final para plan de reorganización detallado.

---

## 🗑️ 1. ARCHIVOS PARA ELIMINAR

### Binarios de Compilación/Debug (110 MB)

| Archivo | Tamaño | Tipo | Razón |
|---------|--------|------|-------|
| `__debug_bin2057900211` | 55 MB | Binario ejecutable Go | Debug de Delve (GoLand/VSCode) |
| `integration.test` | 52 MB | Binario ejecutable de tests | Test compilado temporal |
| `coverage-filtered.out` | 92 KB | Archivo de texto | Duplicado (existe en `/coverage/`) |
| `coverage.out` | 24 KB | Archivo de texto | Duplicado (existe en `/coverage/`) |

#### ¿Por qué eliminar?

1. **`__debug_bin2057900211`**:
   - Generado automáticamente por Delve (debugger de Go)
   - Se regenera cada vez que debuggeas
   - Ya está en `.gitignore` pero quedó en disco

2. **`integration.test`**:
   - Binario temporal de `go test`
   - Se regenera automáticamente
   - Ya está en `.gitignore` (por `*.test`)

3. **`coverage*.out`**:
   - Duplicados de archivos en `/coverage/`
   - No deberían estar en raíz

#### Comando para eliminar:

```bash
rm __debug_bin2057900211 integration.test coverage-filtered.out coverage.out
```

#### Prevenir en el futuro:

Ya están cubiertos por `.gitignore`:
```gitignore
# Archivos de test
*.test
*.out

# Cobertura de tests
*.coverprofile
```

**Recomendación**: Ejecuta el comando de limpieza periódicamente:
```bash
# Agregar al Makefile
make clean-debug
```

---

## 📁 2. ARCHIVOS/CARPETAS PARA MOVER

> **⚠️ IMPORTANTE - Diferencia entre `/docs/` y `/project-docs/`**
>
> - **`/docs/`**: Carpeta EXCLUSIVA para documentación de API generada automáticamente por **Swagger/swag**
>   - Contiene: `docs.go`, `swagger.json`, `swagger.yaml`
>   - **NO mezclar** con documentación del proyecto
>   - Generado por: `swag init -g cmd/main.go`
>
> - **`/project-docs/`**: Carpeta para documentación del PROYECTO (manual)
>   - Informes técnicos
>   - Análisis arquitectónico
>   - Diagramas del proyecto
>   - Guías de desarrollo
>   - Migraciones
>
> **Razón**: Evitar confusión y mantener separado lo generado automáticamente de lo escrito manualmente.

---

### B. Carpeta `analisis-arquitectonico/` (140 KB)

**De**: `/analisis-arquitectonico/`
**A**: `/project-docs/analisis-arquitectonico/`

**Razón**: Es documentación técnica del proyecto, no código fuente ni docs de API

```bash
mkdir -p project-docs
mv analisis-arquitectonico project-docs/
```

---

### C. Carpeta `documentation/` (116 KB) → Renombrar

**De**: `/documentation/`
**A**: `/project-docs/`

**Razón**: Evitar confusión con `/docs/` que es exclusivo de Swagger

**Contenido actual de `/documentation/`**:
- `INFORME-ANALISIS-ECOSISTEMA.md`
- `INFORME_VERSIONADO_CRITICO.md`
- `ISSUE_EDUGO_SHARED_TAGGING.md`
- `MEJORAS_SISTEMA_COMANDOS_AGENTES.md`
- `MIGRACION_EDUGO_SHARED_V2.0.5.md`

```bash
# Renombrar carpeta
mv documentation project-docs
```

**Nota**: Si `project-docs/` ya existe (por pasos A y B), fusionar contenido:

```bash
# Alternativa: Fusionar
mv documentation/* project-docs/
rmdir documentation
```

---

## ⚙️ 3. SALIDAS AUTOMÁTICAS - CONFIGURAR

### A. Binarios Compilados → `/bin/` ✅ YA CONFIGURADO

**Estado**: La carpeta `/bin/` ya existe y está en uso (163 MB)

**Contenido actual**:
- `api-mobile` (59 MB)
- `edugo-api-mobile` (47 MB)
- `test-api` (59 MB)
- `configctl` (4 MB)

**Configuración actual en Makefile**:
```makefile
# Ya está correcto
build:
	@go build -o bin/api-mobile cmd/main.go
```

**✅ No requiere acción**. Ya está bien configurado.

---

### B. Archivos de Cobertura → `/coverage/` ✅ YA CONFIGURADO

**Estado**: La carpeta `/coverage/` ya existe (1.1 MB)

**Problema**: Archivos duplicados en raíz (`coverage.out`, `coverage-filtered.out`)

**Solución**:
1. Eliminar duplicados de raíz (ver sección 1)
2. Actualizar scripts de testing para que SIEMPRE usen `/coverage/`

**Verificar en scripts**:
```bash
# Buscar scripts que generan cobertura
grep -r "coverage.out" scripts/ Makefile test_monitor.sh
```

**Actualizar** para que use:
```bash
# Antes
go test -coverprofile=coverage.out

# Después
mkdir -p coverage
go test -coverprofile=coverage/coverage.out
```

---

### C. Binarios de Test → `/bin/test/` (NUEVO)

**Problema**: `integration.test` se genera en raíz

**Solución**: Configurar salida específica

```bash
# Crear carpeta para binarios de test
mkdir -p bin/test
```

**Actualizar en Makefile**:
```makefile
test-integration:
	@go test -c -o bin/test/integration.test ./test/integration/...
	@./bin/test/integration.test
```

**Actualizar `.gitignore`**:
```gitignore
# Binarios de test
bin/test/
*.test
```

---

### D. Archivos de Debug → `/tmp/` o Ignorar

**Problema**: `__debug_bin*` se genera en raíz

**Solución**: Configurar IDE para usar carpeta temporal

#### Para GoLand/IntelliJ:
1. Settings → Go → Build Tags & Vendoring
2. Output directory: `${PROJECT_DIR}/tmp/debug`

#### Para VSCode (`launch.json`):
```json
{
  "configurations": [
    {
      "name": "Launch",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/cmd/main.go",
      "output": "${workspaceFolder}/tmp/debug/__debug_bin"
    }
  ]
}
```

**Actualizar `.gitignore`**:
```gitignore
# Archivos temporales de debug
tmp/
__debug_bin*
```

---

## ✅ 4. MANTENER EN RAÍZ

Estos archivos/carpetas **DEBEN** permanecer en la raíz:

### Archivos de Configuración del Proyecto

| Archivo | Propósito | Estado |
|---------|-----------|--------|
| `.env`, `.env.example`, `.env.docker` | Variables de entorno | ✅ Correcto |
| `.envrc` | Direnv (carga automática de env) | ✅ Correcto |
| `.gitignore` | Git ignore rules | ✅ Correcto |
| `.dockerignore` | Docker ignore rules | ✅ Correcto |
| `.coverignore` | Coverage ignore rules | ✅ Correcto |
| `go.mod`, `go.sum` | Dependencias de Go | ✅ Correcto |
| `Dockerfile` | Construcción de imagen Docker | ✅ Correcto |
| `docker-compose*.yml` | Orquestación de contenedores | ✅ Correcto |
| `Makefile` | Automatización de tareas | ✅ Correcto |

### Documentación Principal

| Archivo | Propósito | Estado |
|---------|-----------|--------|
| `README.md` | Documentación principal | ✅ Correcto |
| `CHANGELOG.md` | Historial de cambios | ✅ Correcto |
| `QUICKSTART.md` | Inicio rápido | ✅ Correcto |
| `CONFIG.md` | Guía de configuración | ✅ Correcto |
| `COMMIT_GUIDE.md` | Guía de commits | ✅ Correcto |

### Scripts y Herramientas

| Archivo | Propósito | Estado |
|---------|-----------|--------|
| `token` | Atajo para obtener JWT tokens | ✅ Correcto |
| `test_monitor.sh` | Monitor de tests | ✅ Correcto |

### Carpetas de Código Fuente

| Carpeta | Propósito | Estado |
|---------|-----------|--------|
| `cmd/` | Entrypoints de aplicación | ✅ Correcto |
| `internal/` | Código interno de la aplicación | ✅ Correcto |
| `config/` | Archivos de configuración | ✅ Correcto |
| `test/` | Tests del proyecto | ✅ Correcto |
| `scripts/` | Scripts de utilidad | ✅ Correcto |
| `tools/` | Herramientas de desarrollo | ✅ Correcto |

### Carpetas de Proyecto/Organización

| Carpeta | Propósito | Estado |
|---------|-----------|--------|
| `api-tests/` | Testing HTTP con httpyac | ✅ Correcto |
| `sprint/` | Documentación de sprints | ✅ Correcto |
| `docs/` | **Swagger/OpenAPI docs (generados por swag)** | ✅ Correcto |
| `documentation/` | Informes y docs del proyecto | ⚠️ Renombrar a `project-docs/` |
| `analisis-arquitectonico/` | Análisis arquitectónico | ⚠️ Mover a `project-docs/` |
| `bin/` | Binarios compilados | ✅ Correcto |
| `coverage/` | Reportes de cobertura | ✅ Correcto |

### Carpetas de IDE (Ignoradas en Git)

| Carpeta | Propósito | Estado |
|---------|-----------|--------|
| `.vscode/` | Configuración de VSCode | ✅ Ignorado en git |
| `.idea/` | Configuración de GoLand/IntelliJ | ✅ Ignorado en git |
| `.zed/` | Configuración de Zed editor | ✅ Ignorado en git |
| `.claude/` | Configuración de Claude Code | ✅ Versionado (correcto) |
| `.github/` | GitHub Actions, templates | ✅ Versionado (correcto) |
| `.kiro/` | Configuración de Kiro (?) | ⚠️ Verificar si es necesario |

### Archivos de Sistema (Ignorar)

| Archivo | Propósito | Estado |
|---------|-----------|--------|
| `.DS_Store` | Metadata de macOS | ⚠️ Ya está en `.gitignore` |
| `.git/` | Repositorio Git | ✅ Correcto |

---

## 🎯 Plan de Acción Recomendado

### Paso 1: Limpieza Inmediata (Eliminar 110 MB)

```bash
# Desde la raíz del proyecto
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Eliminar binarios temporales
rm -f __debug_bin2057900211 integration.test

# Eliminar archivos de cobertura duplicados
rm -f coverage-filtered.out coverage.out

# Verificar
ls -lh
```

**Tiempo estimado**: 1 minuto
**Liberación de espacio**: ~110 MB

---

### Paso 2: Reorganizar Documentación

```bash
# Crear carpeta para documentación del proyecto
mkdir -p project-docs/diagrams

# Mover diagrama del proyecto
mv architecture-diagram.svg project-docs/diagrams/

# Mover análisis arquitectónico
mv analisis-arquitectonico project-docs/

# Renombrar documentation/ a project-docs/
# PRIMERO verificar si project-docs/ ya existe
if [ -d "project-docs" ]; then
  # Fusionar contenido
  mv documentation/* project-docs/
  rmdir documentation
else
  # Simplemente renombrar
  mv documentation project-docs
fi
```

**⚠️ IMPORTANTE**:
- **NO tocar** la carpeta `/docs/` (es para Swagger)
- Solo hay que verificar si dentro de `/docs/` hay archivos que NO sean de Swagger y moverlos a `/project-docs/`

**Verificar archivos mal clasificados en /docs/**:

```bash
# Listar archivos que NO son de Swagger
ls -la docs/ | grep -v -E "swagger|docs.go|development"

# Si encuentras .md que no sean BOOTSTRAP o TEST, moverlos a project-docs/
# Ejemplo (ajustar según hallazgos):
# mv docs/REFACTORING_*.md project-docs/
# mv docs/IDE_SETUP.md project-docs/
```

**Archivos que SÍ están bien en /docs/**:
- ✅ `docs.go`, `swagger.json`, `swagger.yaml` - Generados por swag
- ✅ Carpeta `development/` - Parece razonable (verificar)
- ⚠️ Archivos `BOOTSTRAP_*.md`, `TESTING_*.md`, `REFACTORING_*.md`, `IDE_SETUP.md` - **Considerar mover a project-docs/**

**Tiempo estimado**: 5 minutos

---

### Paso 3: Configurar Salidas Automáticas

#### A. Actualizar Makefile para Cobertura

```makefile
# Agregar target para limpiar debug
.PHONY: clean-debug
clean-debug:
	@echo "Limpiando archivos de debug..."
	@rm -f __debug_bin* *.test
	@rm -f coverage.out coverage-filtered.out
	@echo "✓ Debug limpiado"

# Actualizar targets de cobertura
coverage:
	@mkdir -p coverage
	@go test -coverprofile=coverage/coverage.out ./...
	@go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	@echo "✓ Reporte de cobertura generado en coverage/"

test-integration:
	@mkdir -p bin/test
	@go test -c -o bin/test/integration.test ./test/integration/...
	@./bin/test/integration.test
	@rm -f bin/test/integration.test
	@echo "✓ Tests de integración completados"
```

#### B. Actualizar `.gitignore`

```gitignore
# Archivos temporales de debug
tmp/
__debug_bin*

# Binarios de test
bin/test/
*.test

# Asegurar que cobertura en raíz se ignore
/coverage.out
/coverage-*.out
```

#### C. Configurar VSCode para Debug

Archivo: `.vscode/launch.json`

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch API",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/cmd/main.go",
      "output": "${workspaceFolder}/tmp/debug/__debug_bin",
      "cwd": "${workspaceFolder}"
    }
  ]
}
```

**Tiempo estimado**: 5 minutos

---

### Paso 4: Agregar Target de Limpieza al Makefile

```makefile
# Target completo de limpieza
.PHONY: clean-all
clean-all: clean-debug
	@echo "Limpiando binarios..."
	@rm -rf bin/*
	@echo "Limpiando cobertura..."
	@rm -rf coverage/*
	@echo "✓ Proyecto completamente limpio"

# Uso:
# make clean-debug  → Solo archivos de debug
# make clean-all    → Todo (binarios + cobertura + debug)
```

---

## 📋 Verificaciones Post-Limpieza

Después de aplicar los cambios, verifica:

### 1. Tamaño del Directorio

```bash
du -sh .
# Esperado: ~206 MB (reducción de ~110 MB)
```

### 2. Archivos en Raíz

```bash
ls -lh | wc -l
# Esperado: ~43 items (reducción de 6 items)
```

### 3. Git Status Limpio

```bash
git status
# No deberían aparecer archivos binarios
```

### 4. Compilación Funcional

```bash
make build
# Debe generar binario en bin/
```

### 5. Tests Funcionales

```bash
make test
# Debe generar cobertura en coverage/
```

---

## 🔄 Mantenimiento Periódico

### Comando Semanal

Agrega esto a tu rutina de desarrollo:

```bash
# Limpiar archivos temporales
make clean-debug

# O crear script en scripts/
./scripts/weekly-cleanup.sh
```

### Script Propuesto: `scripts/weekly-cleanup.sh`

```bash
#!/usr/bin/env bash
set -euo pipefail

echo "🧹 Limpieza semanal del proyecto..."

# Eliminar binarios de debug
rm -f __debug_bin* *.test coverage.out coverage-*.out

# Limpiar binarios viejos (más de 7 días)
find bin/ -type f -mtime +7 -delete 2>/dev/null || true

# Limpiar cobertura vieja
find coverage/ -type f -mtime +14 -delete 2>/dev/null || true

echo "✓ Limpieza completada"
```

---

## 📊 Antes vs Después

| Métrica | Antes | Después | Mejora |
|---------|-------|---------|--------|
| **Tamaño total** | 316 MB | ~206 MB | -35% |
| **Items en raíz** | 49 | 43 | -12% |
| **Archivos temporales** | 4 | 0 | -100% |
| **Carpetas de docs** | 3 dispersas | 1 centralizada | +organización |
| **Binarios en raíz** | 2 (110 MB) | 0 | -100% |

---

## ⚠️ Advertencias

1. **NO ELIMINES** archivos sin verificar antes con `git status`
2. **RESPALDA** antes de mover carpetas grandes
3. **VERIFICA** que los paths actualizados funcionen en scripts
4. **PRUEBA** compilación y tests después de reorganizar

---

## 📝 Notas Adicionales

### Carpeta `.kiro/` (484 KB)

No tengo contexto sobre esta carpeta. Recomendaciones:

- Si es un IDE o herramienta que no usas → Agregar a `.gitignore`
- Si es necesaria para el proyecto → Documentar su propósito
- Si está en desuso → Eliminar

**Acción**: Revisar manualmente qué contiene:
```bash
ls -la .kiro/
```

---

## 📂 Apéndice: Archivos en `/docs/` que NO son de Swagger

### Problema Detectado

La carpeta `/docs/` contiene **archivos mixtos**:
- ✅ Archivos de Swagger (generados automáticamente)
- ⚠️ Archivos de documentación del proyecto (escritos manualmente)

### Contenido Actual de `/docs/`

#### Archivos de Swagger (CORRECTOS - mantener aquí)
- `docs.go` (53 KB)
- `swagger.json` (52 KB)
- `swagger.yaml` (27 KB)

#### Archivos de Documentación del Proyecto (CONSIDERAR MOVER)

**Grupo 1: Bootstrap** (63 KB)
- `BOOTSTRAP_INDEX.md`
- `BOOTSTRAP_MIGRATION_GUIDE.md`
- `BOOTSTRAP_README.md`
- `BOOTSTRAP_USAGE.md`

**Grupo 2: Refactoring** (38 KB)
- `REFACTORING_MAIN.md`
- `REFACTORING_STRUCTURE.md`
- `REFACTORING_SUMMARY.md`

**Grupo 3: Testing** (87 KB)
- `TESTING_GUIDE.md`
- `TESTING_INTEGRATION_GUIDE.md`
- `TESTING_UNIT_GUIDE.md`
- `TEST_ANALYSIS_REPORT.md`
- `TEST_COVERAGE_PLAN.md`
- `TEST_PERFORMANCE_ANALYSIS.md`
- `TEST_PERFORMANCE_RESULTS.md`

**Grupo 4: Desarrollo** (37 KB)
- `IDE_SETUP.md`
- Subcarpeta `development/`
  - `CREDENTIALS.md`
  - `ENVIRONMENT_SETUP.md`
  - `TROUBLESHOOTING.md`

**Total**: ~225 KB de documentación NO-Swagger en `/docs/`

---

### Recomendaciones de Reorganización

#### **Opción A: Separación Total** ⭐ RECOMENDADA

Mover toda la documentación del proyecto a `/project-docs/`:

```bash
# Crear estructura organizada
mkdir -p project-docs/{bootstrap,refactoring,testing,development}

# Mover archivos agrupados
mv docs/BOOTSTRAP_*.md project-docs/bootstrap/
mv docs/REFACTORING_*.md project-docs/refactoring/
mv docs/TEST*.md project-docs/testing/
mv docs/TESTING_*.md project-docs/testing/
mv docs/IDE_SETUP.md project-docs/
mv docs/development project-docs/

# Resultado: /docs/ contiene SOLO Swagger
ls docs/
# docs.go  swagger.json  swagger.yaml
```

**Ventajas**:
- ✅ Claridad total: `/docs/` = solo Swagger
- ✅ Evita confusión futura
- ✅ Facilita regeneración de Swagger sin tocar docs del proyecto
- ✅ Mejor organización por temas

---

#### **Opción B: Subcarpetas dentro de `/docs/`**

Mantener todo en `/docs/` pero organizado:

```bash
# Crear subcarpetas
mkdir -p docs/api docs/project/{bootstrap,refactoring,testing}

# Mover Swagger a /api/
mv docs/{docs.go,swagger.json,swagger.yaml} docs/api/

# Mover docs del proyecto
mv docs/BOOTSTRAP_*.md docs/project/bootstrap/
mv docs/REFACTORING_*.md docs/project/refactoring/
mv docs/TEST*.md docs/project/testing/
mv docs/TESTING_*.md docs/project/testing/
mv docs/IDE_SETUP.md docs/project/
# development/ ya es subcarpeta, moverla
mv docs/development docs/project/
```

**Ventajas**:
- ✅ Todo centralizado en `/docs/`
- ⚠️ Requiere actualizar configuración de Swagger para buscar en `/docs/api/`

---

#### **Opción C: Dejar Como Está**

Mantener la estructura actual y agregar un README explicativo.

Crear `docs/README.md`:

```markdown
# Carpeta /docs/

Esta carpeta contiene documentación mixta:

## 📡 Documentación de API (Swagger)
- `docs.go`, `swagger.json`, `swagger.yaml`
- Generados automáticamente por: `swag init -g cmd/main.go`
- **NO editar manualmente**

## 📚 Documentación del Proyecto
- `BOOTSTRAP_*.md` - Guías de bootstrap
- `REFACTORING_*.md` - Documentos de refactoring
- `TEST*.md` - Análisis y guías de testing
- `IDE_SETUP.md` - Configuración de IDE
- `development/` - Guías de desarrollo
```

**Ventajas**:
- ✅ Sin cambios (rápido)
- ⚠️ Mantiene la confusión entre generado y manual

---

### Decisión Recomendada

**Implementar Opción A** porque:

1. ✅ **Claridad**: `/docs/` exclusivo para Swagger
2. ✅ **Separación**: Generado automático vs. manual
3. ✅ **Escalabilidad**: Fácil agregar más docs del proyecto
4. ✅ **Mantenimiento**: Regenerar Swagger no afecta docs del proyecto
5. ✅ **Convención**: Muchos proyectos usan esta estructura

### Comando Completo para Opción A

```bash
# Crear estructura
mkdir -p project-docs/{bootstrap,refactoring,testing}

# Mover archivos de bootstrap
mv docs/BOOTSTRAP_*.md project-docs/bootstrap/

# Mover archivos de refactoring
mv docs/REFACTORING_*.md project-docs/refactoring/

# Mover archivos de testing
mv docs/TEST*.md docs/TESTING_*.md project-docs/testing/

# Mover IDE setup
mv docs/IDE_SETUP.md project-docs/

# Mover carpeta development
mv docs/development project-docs/

# Verificar que solo queden archivos de Swagger
ls -la docs/
# Debería mostrar solo: docs.go, swagger.json, swagger.yaml

# Crear README en docs/ para documentar
cat > docs/README.md << 'EOF'
# Swagger/OpenAPI Documentation

Esta carpeta contiene la documentación de la API generada automáticamente.

## Archivos

- `docs.go` - Metadata de Swagger
- `swagger.json` - Especificación OpenAPI en JSON
- `swagger.yaml` - Especificación OpenAPI en YAML

## Regenerar Documentación

```bash
swag init -g cmd/main.go
```

## Ver Documentación

Servidor corriendo en: http://localhost:9090

Swagger UI: http://localhost:9090/swagger/index.html

---

**IMPORTANTE**: Esta carpeta es para documentación de API únicamente.
Para documentación del proyecto, ver `/project-docs/`.
EOF

# Crear índice en project-docs/
cat > project-docs/README.md << 'EOF'
# Documentación del Proyecto

Documentación técnica, guías y análisis del proyecto EduGo API Mobile.

## 📂 Estructura

- `bootstrap/` - Guías de bootstrap y migración
- `refactoring/` - Documentos de refactoring
- `testing/` - Análisis y guías de testing
- `development/` - Guías de desarrollo local
- `analisis-arquitectonico/` - Análisis arquitectónico del proyecto
- `diagrams/` - Diagramas del proyecto

## 📄 Documentos Principales

### Bootstrap
- `bootstrap/BOOTSTRAP_INDEX.md`
- `bootstrap/BOOTSTRAP_MIGRATION_GUIDE.md`
- `bootstrap/BOOTSTRAP_README.md`
- `bootstrap/BOOTSTRAP_USAGE.md`

### Refactoring
- `refactoring/REFACTORING_MAIN.md`
- `refactoring/REFACTORING_STRUCTURE.md`
- `refactoring/REFACTORING_SUMMARY.md`

### Testing
- `testing/TESTING_GUIDE.md`
- `testing/TESTING_INTEGRATION_GUIDE.md`
- `testing/TESTING_UNIT_GUIDE.md`
- `testing/TEST_ANALYSIS_REPORT.md`
- `testing/TEST_COVERAGE_PLAN.md`
- `testing/TEST_PERFORMANCE_ANALYSIS.md`
- `testing/TEST_PERFORMANCE_RESULTS.md`

### Desarrollo
- `IDE_SETUP.md`
- `development/CREDENTIALS.md`
- `development/ENVIRONMENT_SETUP.md`
- `development/TROUBLESHOOTING.md`

---

**Última actualización**: 11 de noviembre de 2025
EOF

echo "✅ Reorganización completada"
echo "📁 /docs/ ahora contiene SOLO Swagger"
echo "📁 /project-docs/ contiene toda la documentación del proyecto"
```

---

**Última actualización**: 11 de noviembre de 2025
**Autor**: Claude Code
**Próxima revisión**: Después de aplicar cambios
