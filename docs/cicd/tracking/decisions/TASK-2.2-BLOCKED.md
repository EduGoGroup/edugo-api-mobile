# Decisión: Tarea 2.2 Bloqueada - Migración a Go 1.25

**Fecha:** 2025-11-21
**Tarea:** 2.2 - Migrar a Go 1.25
**Fase:** FASE 1
**Sprint:** SPRINT-2
**Razón del Bloqueo:** Go no disponible en el entorno

---

## Contexto

La Tarea 2.2 requiere:
1. Actualizar `go.mod` a Go 1.25
2. Actualizar workflows (.github/workflows/*.yml) a Go 1.25
3. Actualizar Dockerfile a golang:1.25-alpine
4. Ejecutar `go mod tidy` para actualizar dependencias
5. Compilar con `go build ./...` para validar
6. Ejecutar tests con `go test ./...` para validar

**Problema:** Go no está disponible en el entorno actual (problema de red para descargar).

---

## Decisión

Implementar como **STUB** creando todos los cambios en archivos pero sin ejecutar comandos de Go.

### Estrategia del Stub

1. ✅ Actualizar go.mod manualmente (cambiar versión)
2. ✅ Actualizar todos los workflows manualmente
3. ✅ Actualizar Dockerfile manualmente (si existe)
4. 🟡 go mod tidy: STUB (documentar comando para FASE 2)
5. 🟡 go build: STUB (documentar validación para FASE 2)
6. 🟡 go test: STUB (documentar validación para FASE 2)

---

## Implementación del Stub

### Archivo 1: go.mod

**Cambio necesario:**
```diff
-go 1.24
+go 1.25
```

O si tiene patch version:
```diff
-go 1.24.10
+go 1.25
```

**Archivo:** `/home/user/edugo-api-mobile/go.mod`

---

### Archivo 2: Workflows

**Archivos a actualizar:** Todos los archivos en `.github/workflows/*.yml`

**Cambio necesario:**
```diff
env:
-  GO_VERSION: "1.24.10"
+  GO_VERSION: "1.25"
```

O:
```diff
env:
-  GO_VERSION: "1.24"
+  GO_VERSION: "1.25"
```

**Archivos esperados:**
- `.github/workflows/pr-to-dev.yml`
- `.github/workflows/pr-to-main.yml`
- `.github/workflows/test.yml`
- `.github/workflows/manual-release.yml`
- `.github/workflows/sync-main-to-dev.yml`

---

### Archivo 3: Dockerfile

**Cambio necesario:**
```diff
-FROM golang:1.24-alpine AS builder
+FROM golang:1.25-alpine AS builder
```

O:
```diff
-FROM golang:1.24.10-alpine AS builder
+FROM golang:1.25-alpine AS builder
```

**Archivo:** `/home/user/edugo-api-mobile/Dockerfile` (si existe)

---

## Para FASE 2: Validación Real

Cuando Go esté disponible, ejecutar:

### Paso 1: Verificar Go 1.25 instalado
```bash
go version
# Debe mostrar: go version go1.25.x linux/amd64
```

Si no está instalado:
```bash
# Descargar e instalar Go 1.25
wget https://go.dev/dl/go1.25.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
go version
```

### Paso 2: Actualizar dependencias
```bash
cd /home/user/edugo-api-mobile
go mod tidy
```

**Resultado esperado:** Sin errores, go.sum actualizado

### Paso 3: Compilar proyecto
```bash
go build ./...
```

**Resultado esperado:** Compilación exitosa sin errores

### Paso 4: Ejecutar tests unitarios
```bash
go test -short ./...
```

**Resultado esperado:** Todos los tests pasan

### Paso 5: Ejecutar tests completos (con integración)
```bash
go test ./...
```

**Resultado esperado:** Todos los tests pasan (requiere Docker)

### Paso 6: Verificar race detector
```bash
go test -race -short ./...
```

**Resultado esperado:** Sin race conditions detectadas

### Paso 7: Linter (opcional - se corregirá en tarea 2.10)
```bash
golangci-lint run --timeout=5m
```

**Resultado esperado:** Puede tener errores (se corregirán después)

---

## Criterios de Aceptación (FASE 2)

- [x] go.mod actualizado a `go 1.25`
- [x] Todos los workflows tienen `GO_VERSION: "1.25"`
- [x] Dockerfile actualizado a `golang:1.25-alpine` (si existe)
- [ ] `go mod tidy` ejecuta sin errores
- [ ] `go build ./...` compila exitosamente
- [ ] `go test -short ./...` pasa sin errores
- [ ] `go test ./...` pasa sin errores (con Docker)
- [ ] `go test -race -short ./...` pasa sin race conditions

Los primeros 3 se completan en FASE 1 (stub).
Los últimos 5 se validan en FASE 2 (implementación real).

---

## Rollback Plan (Si falla en FASE 2)

Si la migración a Go 1.25 falla en FASE 2:

```bash
# Revertir cambios
cd /home/user/edugo-api-mobile
git revert <commit-hash>
git push origin claude/sprint-2-phase-1-stubs-015ChMUC8gi8G1Rd21xAMWs1

# O revertir manualmente archivos
sed -i 's/go 1\.25/go 1.24/g' go.mod
sed -i 's/GO_VERSION: "1.25"/GO_VERSION: "1.24"/g' .github/workflows/*.yml
sed -i 's/golang:1\.25-alpine/golang:1.24-alpine/g' Dockerfile

# Restaurar dependencias
go mod tidy
go build ./...
go test -short ./...
```

---

## Estado del Stub

**FASE 1:** ✅ Completado
- Archivos modificados con cambios necesarios
- Documentación completa de validaciones

**FASE 2:** ⏳ Pendiente
- Ejecutar comandos de Go
- Validar compilación y tests
- Confirmar migración exitosa

**Último update:** 2025-11-21
**Por:** Claude Code
