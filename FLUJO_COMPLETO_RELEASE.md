# Flujo Completo: PR → Manual Release → Sync Automático

**Respuestas a tus preguntas**:
1. ✅ PR de dev → main: **Merge normal** (no squash, no rebase)
2. ✅ Sync automático: **Transparente para ti** (workflow lo hace)
3. ✅ manual-release: **Crea commit en main** → sync lo lleva a dev

---

## 🔄 Flujo Completo Paso a Paso

### Paso 1: Estado Inicial (Sincronizado)

```
main: A---B---C (v0.1.6)
dev:  A---B---C (v0.1.6) ← MISMO commit
```

**Verificación**:
```bash
git rev-parse main  # abc123
git rev-parse dev   # abc123 ← IGUAL
```

---

### Paso 2: Desarrollo en Dev

```bash
# Crear feature
git checkout dev
git checkout -b feature/nueva-funcionalidad

# Desarrollar...
git commit -m "feat: nueva funcionalidad"
git push origin feature/nueva-funcionalidad

# PR a dev
gh pr create --base dev --head feature/nueva-funcionalidad
```

**Estado después del merge**:
```
main: A---B---C (v0.1.6, sin cambios)
dev:  A---B---C---D (nueva feature)
```

---

### Paso 3: PR de dev → main (Cuando estés listo para release)

```bash
# Crear PR
gh pr create --base main --head dev --title "Release v0.1.7"
```

**Estrategia de merge**: **MERGE NORMAL** (no squash, no rebase)

```bash
# En GitHub UI o CLI:
gh pr merge --merge  # ← Importante: --merge (no squash, no rebase)
```

**¿Por qué merge normal?**
- ✅ Preserva historial completo
- ✅ Permite fast-forward después
- ✅ Mantiene commits individuales

**Estado después del merge**:
```
main: A---B---C---D (ahora tiene la feature)
dev:  A---B---C---D (mismo commit)
```

**Verificación**:
```bash
git rev-parse main  # def456
git rev-parse dev   # def456 ← IGUAL (por ahora)
```

---

### Paso 4: Ejecutar Manual Release

```bash
# Desde GitHub UI:
# Actions → Manual Release → Run workflow
# - Branch: main
# - Version: 0.1.7
# - Type: minor
```

**Lo que hace `manual-release.yml`**:

#### 4.1. Actualiza archivos
```bash
# Actualiza version.txt
echo "0.1.7" > .github/version.txt

# Actualiza CHANGELOG.md
# (genera entrada automáticamente)
```

#### 4.2. Crea commit en main
```bash
git add .github/version.txt CHANGELOG.md
git commit -m "chore: release v0.1.7"
git push origin main
```

#### 4.3. Crea tag
```bash
git tag -a "v0.1.7" -m "Release v0.1.7"
git push origin "v0.1.7"
```

#### 4.4. Construye Docker
```bash
# Build imagen con tags: v0.1.7, 0.1.7, latest
```

**Estado después de manual-release**:
```
main: A---B---C---D---E (v0.1.7) ← Commit E = "chore: release v0.1.7"
                      ↑
                    tag v0.1.7

dev:  A---B---C---D (atrás por 1 commit)
```

**Ahora main y dev NO están sincronizados** (main tiene commit E que dev no tiene)

---

### Paso 5: Sync Automático (Transparente para ti)

**Trigger**: El push a main (paso 4.2) dispara `sync-main-to-dev-ff.yml`

**Lo que hace el workflow automáticamente**:

```bash
# 1. Checkout dev
git checkout dev

# 2. Fast-forward a main
git merge --ff-only main

# 3. Push
git push origin dev
```

**Estado después del sync**:
```
main: A---B---C---D---E (v0.1.7)
dev:  A---B---C---D---E (v0.1.7) ← MISMO commit
```

**Verificación automática**:
```bash
git rev-parse main  # xyz789
git rev-parse dev   # xyz789 ← IGUAL
```

**✅ Transparente para ti**: No tienes que hacer nada, el workflow lo hace automáticamente.

---

## 📋 Resumen del Flujo Completo

```
┌─────────────────────────────────────────────────────────────┐
│  1. Desarrollo en feature branch                            │
│     feature → PR → dev                                      │
│     Estado: main (C), dev (C-D)                             │
└─────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  2. PR de dev → main (MERGE NORMAL)                         │
│     gh pr merge --merge                                     │
│     Estado: main (C-D), dev (C-D) ← IGUALES                 │
└─────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  3. Manual Release (TÚ ejecutas)                            │
│     - Actualiza version.txt y CHANGELOG.md                  │
│     - Crea commit "chore: release v0.1.7" en main           │
│     - Crea tag v0.1.7                                       │
│     - Construye Docker                                      │
│     Estado: main (C-D-E), dev (C-D) ← DIFERENTES            │
└─────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│  4. Sync Automático (WORKFLOW lo hace)                      │
│     - git merge --ff-only main                              │
│     - git push origin dev                                   │
│     Estado: main (C-D-E), dev (C-D-E) ← IGUALES             │
└─────────────────────────────────────────────────────────────┘
```

---

## 🎯 Respuestas a tus Preguntas

### 1. ¿Cómo es la estrategia de merge del PR dev → main?

**Respuesta**: **MERGE NORMAL** (no squash, no rebase)

```bash
gh pr merge --merge  # ← Importante
```

**Por qué**:
- Preserva historial completo
- Permite fast-forward después
- Mantiene trazabilidad de commits

**Configuración recomendada en GitHub**:
```
Settings → Branches → main → Branch protection rules:
- ✅ Require pull request before merging
- ✅ Allow merge commits
- ❌ Allow squash merging (deshabilitar)
- ❌ Allow rebase merging (deshabilitar)
```

---

### 2. ¿El workflow se encarga de actualizar main a dev? ¿Transparente para mí?

**Respuesta**: **SÍ, 100% transparente**

**Qué hace automáticamente**:
1. Detecta push a main (del manual-release)
2. Hace fast-forward de dev a main
3. Verifica que queden sincronizados
4. Genera reporte

**Tú no haces nada**, solo:
1. Ejecutas manual-release
2. Esperas que termine
3. Verificas (opcional): `git log main..dev` (debe estar vacío)

**Notificación**:
- ✅ Recibes notificación de GitHub Actions
- ✅ Puedes ver el reporte en Actions → Sync Main to Dev

---

### 3. ¿Cómo queda manual-release? ¿El commit de release debe irse a dev?

**Respuesta**: **SÍ, el commit de release va a dev automáticamente**

**Flujo detallado**:

```
Antes de manual-release:
main: A---B---C---D
dev:  A---B---C---D ← IGUALES

Durante manual-release:
1. Actualiza version.txt → "0.1.7"
2. Actualiza CHANGELOG.md → entrada v0.1.7
3. Commit: "chore: release v0.1.7" → commit E
4. Push a main
5. Crea tag v0.1.7
6. Construye Docker

Después de manual-release:
main: A---B---C---D---E (v0.1.7)
                      ↑
                  commit de release
dev:  A---B---C---D (sin commit E todavía)

Sync automático (inmediatamente después):
dev:  A---B---C---D---E (v0.1.7) ← Ahora tiene commit E
```

**El commit E incluye**:
- `.github/version.txt` → "0.1.7"
- `CHANGELOG.md` → entrada de v0.1.7

**Este commit SÍ va a dev** porque:
- Es parte del historial de main
- El sync hace fast-forward
- dev recibe TODO lo de main

---

## 🔍 Verificación Manual (Opcional)

Después de cada release, puedes verificar:

```bash
# 1. Verificar que están en el mismo commit
git fetch origin
git rev-parse origin/main
git rev-parse origin/dev
# Deben ser idénticos

# 2. Verificar que no hay diferencias
git log --oneline origin/main..origin/dev
# Debe estar vacío

# 3. Verificar contenido
git diff origin/main origin/dev
# Debe estar vacío

# 4. Verificar version.txt
git show origin/main:.github/version.txt
git show origin/dev:.github/version.txt
# Deben mostrar la misma versión
```

---

## 📊 Ejemplo Completo con Comandos Reales

### Escenario: Liberar v0.1.7

```bash
# ============================================
# PASO 1: Desarrollo (tú)
# ============================================
git checkout dev
git pull origin dev
git checkout -b feature/nueva-funcionalidad

# ... desarrollar ...
git add .
git commit -m "feat: nueva funcionalidad"
git push origin feature/nueva-funcionalidad

# PR a dev
gh pr create --base dev --head feature/nueva-funcionalidad --title "Nueva funcionalidad"
gh pr merge  # Después de aprobación

# ============================================
# PASO 2: PR a main (tú)
# ============================================
git checkout dev
git pull origin dev

# Crear PR de dev a main
gh pr create --base main --head dev --title "Release v0.1.7"

# Merge (IMPORTANTE: merge normal, no squash)
gh pr merge --merge

# ============================================
# PASO 3: Manual Release (tú)
# ============================================
# Ir a GitHub UI:
# Actions → Manual Release → Run workflow
# - Branch: main
# - Version: 0.1.7
# - Type: minor
# Click "Run workflow"

# Esperar que termine (2-3 minutos)

# ============================================
# PASO 4: Sync Automático (workflow)
# ============================================
# NO HACES NADA
# El workflow automáticamente:
# - Hace fast-forward de dev a main
# - Verifica sincronización
# - Genera reporte

# ============================================
# PASO 5: Verificación (opcional, tú)
# ============================================
git fetch origin

# Verificar que están sincronizados
git log --oneline origin/main..origin/dev
# Salida: (vacío) ← ✅

git rev-parse origin/main
# Salida: abc123def456...

git rev-parse origin/dev
# Salida: abc123def456... ← ✅ IGUAL

echo "✅ Sincronizados correctamente"
```

---

## 🎯 Ventajas de Este Flujo

### 1. **Transparente**
- ✅ No tienes que sincronizar manualmente
- ✅ Workflow lo hace automáticamente
- ✅ Recibes notificación si algo falla

### 2. **Confiable**
- ✅ Mismo commit = mismo contenido
- ✅ Verificación automática
- ✅ Falla si hay divergencia

### 3. **Simple**
- ✅ Tú solo ejecutas manual-release
- ✅ El resto es automático
- ✅ Sin pasos manuales de sync

### 4. **Verificable**
- ✅ `git log main..dev` siempre vacío
- ✅ `git rev-parse` muestra mismo SHA
- ✅ Sin ambigüedad

---

## ⚠️ Caso Especial: ¿Qué pasa si el sync falla?

### Escenario: Dev tiene commits que main no tiene

```
main: A---B---C---D---E (v0.1.7)
dev:  A---B---C---D---X (commit X no está en main)
```

**El workflow falla** con mensaje:
```
⚠️ ADVERTENCIA: dev tiene commits que main no tiene
⚠️ No se puede hacer fast-forward automático
⚠️ Acción manual requerida

Commits en dev que NO están en main:
X - feat: trabajo en progreso
```

**Solución**:
```bash
# Opción 1: Llevar commit X a main primero
git checkout main
git cherry-pick X
git push origin main
# Luego el sync funcionará

# Opción 2: Descartar commit X (si era experimental)
git checkout dev
git reset --hard origin/main
git push --force origin dev
```

**Esto es BUENO** porque:
- ✅ Te alerta de divergencia
- ✅ Evita pérdida accidental de trabajo
- ✅ Requiere decisión consciente

---

## 📝 Checklist de Release

```
[ ] 1. Desarrollo completo en dev
[ ] 2. PR de dev → main (merge normal)
[ ] 3. Ejecutar Manual Release en GitHub UI
[ ] 4. Esperar que termine manual-release
[ ] 5. Verificar que sync automático terminó exitosamente
[ ] 6. (Opcional) Verificar: git log main..dev (vacío)
[ ] 7. Continuar desarrollo en dev
```

---

## 🚀 Resultado Final

**Después de cada release**:
```bash
git log --oneline main
# A---B---C---D---E (v0.1.7)

git log --oneline dev
# A---B---C---D---E (v0.1.7) ← IDÉNTICO

git diff main dev
# (vacío) ← Sin diferencias

git rev-parse main
# abc123...

git rev-parse dev
# abc123... ← MISMO SHA
```

**Garantía**: main y dev siempre sincronizados después de release, sin intervención manual.

---

**¿Esto responde tus preguntas? ¿Proceder con la implementación?**
