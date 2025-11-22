# Pre-commit Hooks - edugo-api-mobile

Configuración de hooks de pre-commit para garantizar calidad de código antes de hacer commits.

---

## 📋 ¿Qué son los Pre-commit Hooks?

Los **pre-commit hooks** son validaciones automáticas que se ejecutan **antes** de cada commit. Ayudan a:

- ✅ Prevenir errores comunes antes de push
- ✅ Mantener código formateado consistentemente
- ✅ Detectar problemas de seguridad temprano
- ✅ Validar sintaxis de archivos
- ✅ Reducir errores en CI/CD

---

## 🚀 Instalación

### 1. Instalar pre-commit

```bash
# macOS
brew install pre-commit

# Linux
pip install pre-commit

# Verificar instalación
pre-commit --version
```

### 2. Activar hooks en el repositorio

```bash
# Desde la raíz del proyecto
cd /path/to/edugo-api-mobile

# Instalar hooks
pre-commit install

# Verificar instalación
pre-commit --version
```

**Resultado esperado:**
```
pre-commit installed at .git/hooks/pre-commit
```

---

## 🔧 Configuración

Los hooks están configurados en `.pre-commit-config.yaml` con las siguientes validaciones:

### Validaciones Generales (7)

| Hook | Descripción | Ejemplo de Error |
|------|-------------|------------------|
| **no-commit-to-branch** | Previene commits directos a `main`/`dev` | `❌ No puedes hacer commit directo a main` |
| **trailing-whitespace** | Remueve espacios al final de líneas | `❌ Línea 45 tiene espacios finales` |
| **end-of-file-fixer** | Asegura salto de línea al final | `✅ Agregado \n al final de archivo` |
| **check-added-large-files** | Previene archivos >500KB | `❌ video.mp4 es muy grande (2MB)` |
| **check-yaml** | Valida sintaxis YAML | `❌ workflow.yml tiene sintaxis inválida` |
| **check-json** | Valida sintaxis JSON | `❌ config.json falta una coma` |
| **check-merge-conflict** | Detecta markers de conflictos | `❌ Encontrado <<<<<<< HEAD` |
| **detect-private-key** | Detecta claves privadas | `❌ ⚠️  Clave SSH detectada` |

### Validaciones Go (4)

| Hook | Descripción | Ejemplo de Corrección |
|------|-------------|----------------------|
| **go fmt** | Formatea código automáticamente | `✅ main.go formateado` |
| **go vet** | Detecta errores comunes | `❌ Printf tiene argumentos incorrectos` |
| **go mod tidy** | Limpia dependencias | `✅ go.mod actualizado` |
| **golangci-lint** | Linting completo (opcional) | `❌ 3 errores de errcheck detectados` |

### Validaciones Adicionales (1)

| Hook | Descripción | Nota |
|------|-------------|------|
| **dockerfile-lint** | Valida Dockerfile | Requiere `hadolint` instalado (opcional) |

---

## 💻 Uso

### Automático (Recomendado)

Los hooks se ejecutan **automáticamente** en cada `git commit`:

```bash
# Hacer cambios
vim internal/handler/user_handler.go

# Agregar al staging
git add .

# Commit (hooks se ejecutan automáticamente)
git commit -m "feat: agregar endpoint de usuarios"
```

**Ejemplo de salida:**
```
no-commit-to-branch..................................................Passed
trailing-whitespace..................................................Passed
end-of-file-fixer....................................................Passed
check-yaml...........................................................Passed
go fmt...............................................................Passed
go vet...............................................................Passed
go mod tidy..........................................................Passed
golangci-lint........................................................Passed
[feature/users 3a4b5c6] feat: agregar endpoint de usuarios
 2 files changed, 45 insertions(+)
```

### Manual (Testing)

```bash
# Ejecutar todos los hooks manualmente
pre-commit run --all-files

# Ejecutar hook específico
pre-commit run go-fmt --all-files
pre-commit run golangci-lint --all-files

# Ver lista de hooks configurados
pre-commit run --list
```

---

## ⚙️ Configuración Personalizada

### Desactivar Temporalmente

```bash
# Opción 1: Skip hooks para un commit específico
git commit --no-verify -m "WIP: trabajo en progreso"

# Opción 2: Desactivar permanentemente
git config core.hooksPath .git/hooks
```

### Desactivar Hook Específico

Editar `.pre-commit-config.yaml` y comentar el hook:

```yaml
# - id: golangci-lint  # ← Comentar para desactivar
#   name: golangci-lint
#   entry: bash -c 'golangci-lint run --fast'
```

### Cambiar Severidad de golangci-lint

Si `golangci-lint` es muy lento, moverlo a ejecución manual:

```yaml
- id: golangci-lint
  name: golangci-lint
  entry: bash -c 'golangci-lint run --fast'
  language: system
  files: \.go$
  pass_filenames: false
  stages: [manual]  # ← Solo manual, no automático
```

Luego ejecutar manualmente:
```bash
pre-commit run golangci-lint --all-files
```

---

## 🐛 Solución de Problemas

### Error: "command not found: pre-commit"

**Solución:**
```bash
# Instalar pre-commit
brew install pre-commit  # macOS
pip install pre-commit   # Linux/Windows
```

### Error: "golangci-lint: command not found"

**Solución:**
```bash
# Instalar golangci-lint
brew install golangci-lint  # macOS
# O
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

**Alternativa:** Comentar el hook `golangci-lint` en `.pre-commit-config.yaml`

### Hooks muy lentos

**Solución:**
1. Desactivar `golangci-lint` (el más lento)
2. Usar `--fast` en golangci-lint
3. Ejecutar solo en archivos modificados (default)

### Error: "go: cannot find main module"

**Solución:**
```bash
# Verificar que estás en la raíz del proyecto
pwd  # Debe mostrar: /path/to/edugo-api-mobile

# Verificar que existe go.mod
ls go.mod  # Debe existir
```

---

## 📊 Tiempos de Ejecución

Tiempos aproximados en MacBook Pro M1:

| Hook | Tiempo | Impacto |
|------|--------|---------|
| Validaciones generales | <1s | ⚡ Muy rápido |
| go fmt | <1s | ⚡ Muy rápido |
| go vet | 2-3s | ⚡ Rápido |
| go mod tidy | 1-2s | ⚡ Rápido |
| golangci-lint (--fast) | 5-10s | ⚠️  Moderado |
| golangci-lint (completo) | 20-30s | 🐢 Lento |

**Tiempo total sin golangci-lint:** ~5 segundos  
**Tiempo total con golangci-lint --fast:** ~10-15 segundos

---

## 🎯 Mejores Prácticas

### ✅ Recomendado

- ✅ Mantener hooks habilitados todo el tiempo
- ✅ Ejecutar `pre-commit run --all-files` después de pull
- ✅ Agregar validaciones específicas del proyecto según necesidad
- ✅ Usar `--no-verify` solo en casos excepcionales (WIP, rebases)

### ❌ Evitar

- ❌ Desactivar hooks permanentemente
- ❌ Hacer `git commit --no-verify` habitualmente
- ❌ Ignorar errores de hooks

---

## 🔄 Actualizar Hooks

```bash
# Actualizar versiones de hooks
pre-commit autoupdate

# Limpiar cache
pre-commit clean

# Reinstalar hooks
pre-commit uninstall
pre-commit install
```

---

## 📚 Recursos Adicionales

- [Pre-commit Documentation](https://pre-commit.com/)
- [Supported Hooks](https://pre-commit.com/hooks.html)
- [golangci-lint Configuration](https://golangci-lint.run/usage/configuration/)
- [go fmt Documentation](https://pkg.go.dev/cmd/gofmt)

---

## 💬 Preguntas Frecuentes

**P: ¿Son obligatorios los hooks?**  
R: No, son opcionales. Puedes usar `--no-verify` para saltarlos, pero no es recomendado.

**P: ¿Afectan el rendimiento del commit?**  
R: Sí, agregan ~5-15 segundos por commit, pero previenen errores costosos en CI/CD.

**P: ¿Puedo personalizar los hooks?**  
R: Sí, edita `.pre-commit-config.yaml` según tus necesidades.

**P: ¿Funcionan en todos los sistemas operativos?**  
R: Sí, pre-commit es compatible con macOS, Linux y Windows.

---

**Última actualización:** 2025-11-21  
**Versión de pre-commit:** 4.5.0  
**Generado por:** Claude Code
