# 🔄 Workflows de CI/CD - edugo-api-mobile

## 📋 Workflows Configurados

### 1️⃣ **ci.yml** - Pipeline de Integración Continua

**Trigger:**
- ✅ Pull Requests a `main` o `develop`
- ✅ Push directo a `main` (red de seguridad)

**Ejecuta:**
- ✅ Verificación de formato (gofmt)
- ✅ Verificación de go.mod y go.sum sincronizados
- ✅ Análisis estático (go vet)
- ✅ Tests con race detection
- ✅ Build verification
- ✅ Verificación de Swagger docs
- ✅ Linter (opcional, no bloquea)
- ✅ Security scan con gosec

**Cuándo se ejecuta:**
```bash
# Cuando creas un PR
gh pr create --title "..." --body "..."  # ← AQUÍ se ejecuta

# O cuando alguien hace push directo a main (no recomendado)
git push origin main  # ← AQUÍ se ejecuta
```

**Duración estimada:** 3-4 minutos

---

### 2️⃣ **test.yml** - Tests con Cobertura

**Trigger:**
- ✅ Manual (workflow_dispatch desde GitHub UI)
- ✅ Pull Requests a `main` o `develop`

**Ejecuta:**
- ✅ Tests unitarios con cobertura detallada
- ✅ Generación de reporte HTML
- ✅ Upload a Codecov
- ✅ Comentario en PR con porcentaje de cobertura
- ✅ Tests de integración con PostgreSQL y MongoDB (opcional)

**Cuándo se ejecuta:**
```bash
# Manual desde GitHub UI:
# Actions → Tests with Coverage → Run workflow

# O automáticamente en PRs
gh pr create  # ← AQUÍ se ejecuta
```

**Duración estimada:** 4-5 minutos

---

### 3️⃣ **build-and-push.yml** - Build y Push de Docker

**Trigger:**
- ✅ Manual (workflow_dispatch con selección de ambiente)
- ✅ Push a `main` (automático)

**Ejecuta:**
- ✅ Tests antes del build
- ✅ Build de imagen Docker
- ✅ Push a GitHub Container Registry (ghcr.io)
- ✅ Tags automáticos (latest, branch, sha, environment)
- ✅ Resumen detallado del deployment

**Cuándo se ejecuta:**
```bash
# Automático cuando haces push a main
git push origin main  # ← AQUÍ se ejecuta

# Manual desde GitHub UI con selección de ambiente:
# Actions → Build and Push Docker Image → Run workflow
# Seleccionar: development, staging, o production
```

**Tags generados:**
- `latest` - Último build de main
- `main-<sha>` - Build específico por commit
- `<environment>` - Tag del ambiente seleccionado (manual)
- `<environment>-YYYYMMDD-HHmmss` - Tag con timestamp (manual)

**Duración estimada:** 5-8 minutos

---

### 4️⃣ **release.yml** - Release Completo (TAGS)

**Trigger:** Solo cuando creas un tag `v*` (ej: `v1.0.0`, `v2.1.3`)

**Ejecuta:**
- ✅ Validación completa del código
- ✅ Tests con cobertura
- ✅ Build de imagen Docker con tags versionados
- ✅ Creación automática de GitHub Release
- ✅ Generación de changelog desde commits o CHANGELOG.md
- ✅ Documentación de deployment en el release

**Cuándo se ejecuta:**
```bash
# Crear y pushear tag
git tag -a v1.0.0 -m "Release 1.0.0: Primera versión estable"
git push origin v1.0.0  # ← AQUÍ se ejecuta
```

**Tags Docker generados:**
- `v1.0.0` - Versión semántica completa
- `v1.0` - Major.Minor
- `v1` - Major
- `latest` - Última versión
- `v1.0.0-<sha>` - Con commit hash

**Duración estimada:** 6-10 minutos

---

## 🎯 Estrategia de CI/CD Optimizada

### **Flujo Normal de Desarrollo:**

```
┌─────────────────────────────────────────────────────────────┐
│  1. Desarrollo Local                                        │
│     - Hacer cambios en tu branch                           │
│     - Ejecutar tests localmente: go test ./...             │
│     - Verificar formato: gofmt -w .                        │
│     - git commit                                            │
│     ✅ NO GASTA MINUTOS DE GITHUB                           │
└─────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  2. Crear Pull Request                                      │
│     - gh pr create                                          │
│     - CI automático (ci.yml + test.yml)                     │
│     - Revisar resultados y cobertura                        │
│     - Aprobar y mergear                                     │
│     ✅ VALIDA ANTES DE MERGE                                │
└─────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  3. Merge a Main                                            │
│     - gh pr merge                                           │
│     - CI de seguridad (ci.yml)                             │
│     - Build automático de imagen Docker                     │
│     ✅ CÓDIGO VALIDADO + IMAGEN EN GHCR                     │
└─────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  4. Crear Release (cuando estés listo)                     │
│     - Actualizar CHANGELOG.md (opcional)                    │
│     - git tag -a v1.2.0 -m "Release 1.2.0"                  │
│     - git push origin v1.2.0                                │
│     - Release automático (release.yml)                      │
│     - Imagen Docker con tags versionados                    │
│     ✅ RELEASE COMPLETO CON DOCUMENTACIÓN                   │
└─────────────────────────────────────────────────────────────┘
```

---

## 🐳 Gestión de Imágenes Docker

### **Después de cada push a main:**
```bash
# La imagen se publica automáticamente como:
docker pull ghcr.io/edugogroup/edugo-api-mobile:latest
docker pull ghcr.io/edugogroup/edugo-api-mobile:main-abc1234
```

### **Cuando creas un release (tag):**
```bash
# Se publican múltiples tags versionados:
docker pull ghcr.io/edugogroup/edugo-api-mobile:v1.2.0
docker pull ghcr.io/edugogroup/edugo-api-mobile:v1.2
docker pull ghcr.io/edugogroup/edugo-api-mobile:v1
docker pull ghcr.io/edugogroup/edugo-api-mobile:latest
```

### **Deploy manual de ambiente específico:**
```bash
# Desde GitHub UI: Actions → Build and Push → Run workflow
# Seleccionar ambiente: production

# Resultado:
docker pull ghcr.io/edugogroup/edugo-api-mobile:production
docker pull ghcr.io/edugogroup/edugo-api-mobile:production-20251031-143000
```

---

## 💰 Ahorro de Minutos de GitHub Actions

### **Estrategia Optimizada:**

| Escenario | Workflows Ejecutados | Minutos Estimados |
|-----------|---------------------|-------------------|
| Push a branch feature | 0 (no ejecuta nada) | 0 min |
| Crear PR | ci.yml + test.yml | ~8 min |
| Merge a main | ci.yml + build-and-push.yml | ~12 min |
| Crear release (tag) | release.yml | ~10 min |

**Mes típico (10 PRs, 3 releases):**
- 10 PRs = 80 minutos
- 10 merges a main = 120 minutos
- 3 releases = 30 minutos
- **Total = 230 minutos/mes** (✅ Solo ~10% del plan gratuito de 2,000 min)

---

## 🚀 Guía Rápida

### **Para desarrollo normal:**
```bash
# 1. Crear branch de feature
git checkout -b feature/nueva-funcionalidad

# 2. Desarrollar y probar localmente
go test ./...
gofmt -w .

# 3. Commit y push
git commit -m "feat: nueva funcionalidad"
git push origin feature/nueva-funcionalidad

# 4. Crear PR (ejecuta ci.yml + test.yml automáticamente)
gh pr create --title "Nueva funcionalidad" --body "..."

# 5. Esperar aprobación y merge
# Al hacer merge, se ejecuta automáticamente build-and-push.yml
```

### **Para crear una release:**
```bash
# 1. Asegurarse de estar en main actualizado
git checkout main
git pull origin main

# 2. Actualizar CHANGELOG.md (opcional pero recomendado)
vim CHANGELOG.md
git add CHANGELOG.md
git commit -m "chore: actualizar changelog para v1.2.0"
git push origin main

# 3. Crear y pushear tag (ejecuta release.yml automáticamente)
git tag -a v1.2.0 -m "Release 1.2.0: Nuevas funcionalidades X, Y, Z"
git push origin v1.2.0

# 4. GitHub Actions:
#    - Valida todo el código
#    - Ejecuta tests
#    - Construye imagen Docker
#    - Crea GitHub Release automáticamente
#    - Publica documentación
```

### **Para deploy manual a un ambiente:**
```bash
# Opción 1: Desde GitHub UI
# 1. Ir a Actions → Build and Push Docker Image
# 2. Click en "Run workflow"
# 3. Seleccionar ambiente (development/staging/production)
# 4. Click "Run workflow"

# Opción 2: Desde CLI con gh
gh workflow run build-and-push.yml -f environment=production
```

---

## 📊 Monitoreo de Workflows

### **Ver estado de workflows:**
```bash
# Listar últimos workflows ejecutados
gh run list --limit 10

# Ver detalles de un workflow específico
gh run view <run-id>

# Ver logs de un workflow
gh run view <run-id> --log

# Re-ejecutar un workflow fallido
gh run rerun <run-id>

# Ver workflows en ejecución
gh run watch
```

### **Ver imagen Docker publicada:**
```bash
# Autenticarse en GHCR
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin

# Ver tags disponibles
gh api /orgs/EduGoGroup/packages/container/edugo-api-mobile/versions

# Pull de la imagen
docker pull ghcr.io/edugogroup/edugo-api-mobile:latest
```

---

## 🛡️ Branch Protection (Recomendado)

Para forzar el uso de PRs y garantizar calidad:

1. GitHub → Settings → Branches → Add rule
2. Branch name pattern: `main`
3. Configurar:
   - ✅ Require pull request before merging
   - ✅ Require approvals: 1
   - ✅ Require status checks to pass:
     - `Tests and Validation`
     - `Tests with Coverage`
   - ✅ Require branches to be up to date
   - ✅ Do not allow bypassing the above settings

---

## 🔍 Troubleshooting

### **Error: "GOPRIVATE no configurado"**
```bash
# Asegúrate de que el workflow tiene acceso a repos privados
# Ya está configurado en los workflows con:
git config --global url."https://${{ secrets.GITHUB_TOKEN }}@github.com/".insteadOf "https://github.com/"
```

### **Error: "No se puede pushear imagen Docker"**
```bash
# Verifica permisos del workflow
# Los workflows necesitan: permissions.packages: write
# Ya está configurado en build-and-push.yml y release.yml
```

### **Workflow no se ejecuta en tag:**
```bash
# Asegúrate de que el tag tenga el prefijo 'v'
git tag v1.0.0  # ✅ Correcto
git tag 1.0.0   # ❌ No ejecutará release.yml

# Push del tag
git push origin v1.0.0
```

---

## 📚 Recursos Adicionales

- [Documentación de GitHub Actions](https://docs.github.com/en/actions)
- [GitHub Container Registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [Codecov Documentation](https://docs.codecov.com/)
- [Guía de Migración edugo-shared](../../MIGRACION_EDUGO_SHARED_V2.0.5.md)

---

## 📝 Checklist para Nuevos Proyectos

Si vas a replicar estos workflows en otros proyectos:

- [ ] Copiar los 4 archivos de workflows
- [ ] Actualizar `GO_VERSION` a la versión de Go del proyecto
- [ ] Actualizar `IMAGE_NAME` si es necesario
- [ ] Verificar que existe Swagger (o comentar esa sección)
- [ ] Configurar branch protection en GitHub
- [ ] Hacer un PR de prueba para validar ci.yml y test.yml
- [ ] Crear un tag de prueba para validar release.yml
- [ ] Documentar workflows específicos del proyecto

---

**Última actualización:** 2025-10-31
**Mantenedor:** Equipo EduGo
**Proyecto:** edugo-api-mobile
