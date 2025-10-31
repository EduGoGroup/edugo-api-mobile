# 🐛 Problema: Submódulos no pueden descargarse con versión específica v2.0.5

## Descripción del Problema

Los submódulos de `edugo-shared` v2.0.5 (common, auth, logger) **no pueden instalarse** porque faltan los tags específicos por submódulo. Actualmente solo existe el tag `v2.0.5` a nivel raíz del repositorio, pero Go requiere tags individuales para cada submódulo en un monorepo.

## Reproducción del Error

Al intentar instalar los submódulos:

```bash
# Eliminar la versión monolítica
go mod edit -droprequire github.com/EduGoGroup/edugo-shared/v2

# Intentar instalar submódulos con v2.0.5
go get github.com/EduGoGroup/edugo-shared/common@v2.0.5
# ❌ Error: module github.com/EduGoGroup/edugo-shared@v2.0.5 found (v2.0.5+incompatible),
#           but does not contain package github.com/EduGoGroup/edugo-shared/common

# Intentar con @latest
go get github.com/EduGoGroup/edugo-shared/common@latest
# ✅ Se descarga pero con pseudo-versión: v0.0.0-20251031205144-a9148968daba

go get github.com/EduGoGroup/edugo-shared/auth@latest
# ❌ Error: requires github.com/EduGoGroup/edugo-shared/common@v0.0.0:
#           unknown revision common/v0.0.0
```

## Causa Raíz

En monorepos con múltiples módulos Go, cada submódulo necesita su propio tag versionado:

### Estado Actual ❌
```bash
git tag -l
# Salida:
v2.0.5
v2.0.1
v2.0.0
v1.0.0
v0.1.0
```

### Estado Esperado ✅
```bash
git tag -l
# Debería incluir:
v2.0.5
common/v2.0.5
auth/v2.0.5
logger/v2.0.5
messaging/v2.0.5
postgres/v2.0.5
mongodb/v2.0.5
```

## Problema Adicional: Directivas `replace`

Los archivos `go.mod` de los submódulos contienen directivas `replace` apuntando a paths locales:

**auth/go.mod:**
```go
module github.com/EduGoGroup/edugo-shared/auth

require (
    github.com/EduGoGroup/edugo-shared/common v0.0.0
    // ... otras dependencias
)

replace github.com/EduGoGroup/edugo-shared/common => ../common
```

**Problema:** Estas directivas `replace` funcionan en desarrollo local, pero causan conflictos cuando otros proyectos intentan usar los submódulos como dependencias remotas.

## Solución Requerida

### 1. Crear Tags Específicos para cada Submódulo

Para la versión **v2.0.5**, ejecutar:

```bash
# Asegurarse de estar en el commit correcto de v2.0.5
git checkout v2.0.5

# Crear tags para cada submódulo
git tag common/v2.0.5
git tag auth/v2.0.5
git tag logger/v2.0.5
git tag messaging/v2.0.5
git tag postgres/v2.0.5
git tag mongodb/v2.0.5

# Publicar todos los tags
git push origin --tags
```

### 2. Eliminar las Directivas `replace` de los go.mod publicados

**Opción A: Eliminar completamente las directivas `replace`**

Editar cada `go.mod` de los submódulos y eliminar las líneas `replace`:

```diff
# auth/go.mod
module github.com/EduGoGroup/edugo-shared/auth

require (
-   github.com/EduGoGroup/edugo-shared/common v0.0.0
+   github.com/EduGoGroup/edugo-shared/common v2.0.5
    github.com/golang-jwt/jwt/v5 v5.3.0
    // ... otras dependencias
)

- replace github.com/EduGoGroup/edugo-shared/common => ../common
```

**Opción B: Mantener `replace` solo para desarrollo local**

Si quieren mantener las directivas `replace` para desarrollo local, documentar en el README que los desarrolladores deben ejecutar:

```bash
# Para desarrollo local
go mod edit -replace github.com/EduGoGroup/edugo-shared/common=../common

# Antes de hacer commit/release, eliminar replace
go mod edit -dropreplace github.com/EduGoGroup/edugo-shared/common
```

### 3. Actualizar Referencias de Versión entre Submódulos

Asegurarse de que cuando un submódulo dependa de otro, use la versión correcta:

```go
// auth/go.mod
require (
    github.com/EduGoGroup/edugo-shared/common v2.0.5  // ← Versión específica, no v0.0.0
)
```

### 4. Proceso para Futuros Releases

Documentar el proceso de release para mantener consistencia:

```bash
# 1. Hacer cambios y commits
git add .
git commit -m "feat: nuevas funcionalidades para v2.0.6"

# 2. Crear tag principal
git tag v2.0.6

# 3. Crear tags de submódulos
for module in common auth logger messaging postgres mongodb; do
  git tag ${module}/v2.0.6
done

# 4. Publicar todos los tags
git push origin --tags
```

## Verificación Post-Fix

Después de aplicar la solución, verificar que funcione:

```bash
# Crear proyecto de prueba
mkdir test-edugo-shared
cd test-edugo-shared
go mod init test

# Instalar submódulos
go get github.com/EduGoGroup/edugo-shared/common@v2.0.5
go get github.com/EduGoGroup/edugo-shared/auth@v2.0.5
go get github.com/EduGoGroup/edugo-shared/logger@v2.0.5

# Verificar que se instalaron correctamente
go list -m github.com/EduGoGroup/edugo-shared/common
# Salida esperada: github.com/EduGoGroup/edugo-shared/common v2.0.5

go list -m github.com/EduGoGroup/edugo-shared/auth
# Salida esperada: github.com/EduGoGroup/edugo-shared/auth v2.0.5
```

## Referencias

- [Go Modules: Multi-module repositories](https://github.com/golang/go/wiki/Modules#multi-module-repositories)
- [Go Module Version Control](https://go.dev/ref/mod#vcs)
- [Semantic Versioning for Multi-Module Repos](https://go.dev/blog/module-compatibility)

## Impacto

Este problema bloquea a **todos los proyectos** que intenten migrar de v2.0.1 (monolítico) a v2.0.5 (modular), anulando el beneficio principal de la arquitectura modular anunciado en el release.

### Proyectos Afectados
- `edugo-api-mobile` - Bloqueado en migración a v2.0.5
- Cualquier otro proyecto que intente usar los submódulos individualmente

---

**Prioridad:** 🔴 Alta - Bloquea adopción de v2.0.5

**Tipo:** 🐛 Bug - Configuración de release

**Componente:** DevOps / Release Management
