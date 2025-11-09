# Propuesta: Estrategia de Ramas Transparente y Confiable

**Fecha**: 2025-11-08  
**Problema**: Historial confuso con merges bidireccionales que genera desconfianza  
**Solución**: Fast-Forward Only - main y dev siempre en el mismo commit después de sync

---

## 🔴 Problema Actual

### Síntomas:
- ✅ Contenido sincronizado
- ❌ Historial confuso (31 commits de diferencia)
- ❌ Imposible verificar visualmente si están iguales
- ❌ "El lobo viene" - pierdes confianza en el sistema
- ❌ Commits de sync crean ruido

### Historial Actual (Enredado):
```
*   107b55c (dev) chore: sync main v0.1.6 to dev
|\  
| * f0f9a63 (main) chore: release v0.1.6
* | f5dc923 chore: sync main v0.1.5 to dev
|\| 
| * b99d439 the main (#32)
* | f13c21a Feature/infrastructure bootstrap refactor
```

**Problema**: Merges bidireccionales crean historial imposible de seguir.

---

## ✅ Solución Propuesta: Fast-Forward Only

### Principio Fundamental:
**"main y dev SIEMPRE apuntan al mismo commit después de sincronización"**

### Garantía:
```bash
# Después de cada sync:
git rev-parse main == git rev-parse dev  # ← MISMO SHA
git diff main dev                         # ← Sin diferencias
git log main..dev                         # ← Vacío
```

---

## 📋 Flujo de Trabajo Propuesto

### 1. Estado Inicial (Después de Release)
```
main: A---B---C (v0.1.6)
dev:  A---B---C (v0.1.6) ← MISMO commit, MISMO SHA
```

**Verificación**:
```bash
git rev-parse main  # f0f9a63...
git rev-parse dev   # f0f9a63... ← IGUAL
```

---

### 2. Desarrollo de Nueva Feature
```
main: A---B---C (sin cambios)
dev:  A---B---C---D---E (nueva feature)
```

**Flujo**:
```bash
# Crear feature branch desde dev
git checkout dev
git pull origin dev
git checkout -b feature/nueva-funcionalidad

# Desarrollar...
git commit -m "feat: nueva funcionalidad"
git push origin feature/nueva-funcionalidad

# PR a dev
gh pr create --base dev --head feature/nueva-funcionalidad
# Merge después de aprobación
```

---

### 3. Release (PR de dev → main)
```
# Antes del PR:
main: A---B---C
dev:  A---B---C---D---E

# Después del PR (fast-forward):
main: A---B---C---D---E ← Mismo commit que dev
dev:  A---B---C---D---E
```

**Flujo**:
```bash
# Cuando estés listo para release
gh pr create --base main --head dev --title "Release v0.1.7"

# Merge (debe ser fast-forward)
gh pr merge --merge  # NO squash, NO rebase
```

---

### 4. Crear Release y Sincronizar
```
# Después de manual-release:
main: A---B---C---D---E---F (v0.1.7)
dev:  A---B---C---D---E (atrás por 1 commit)

# Workflow automático hace fast-forward:
dev:  A---B---C---D---E---F (v0.1.7) ← MISMO commit
```

**Automático**: El workflow `sync-main-to-dev-ff.yml` hace:
```bash
git checkout dev
git merge --ff-only main  # Fast-forward, sin merge commit
git push origin dev
```

---

### 5. Hotfix en Main
```
# Bug crítico en producción
main: A---B---C---D---E---F (v0.1.7)
dev:  A---B---C---D---E---F (v0.1.7)

# Fix en main:
git checkout -b hotfix/critical-bug main
git commit -m "fix: critical bug"
gh pr create --base main --head hotfix/critical-bug

# Después del merge:
main: A---B---C---D---E---F---G (hotfix)
dev:  A---B---C---D---E---F (atrás por 1)

# Workflow automático sincroniza:
dev:  A---B---C---D---E---F---G ← MISMO commit
```

---

## 🎯 Ventajas de Esta Estrategia

### 1. **Transparencia Total**
```bash
# Siempre puedes verificar:
git log --oneline main
git log --oneline dev
# ← Historial IDÉNTICO después de sync
```

### 2. **Confianza**
- ✅ Mismo SHA = mismo contenido GARANTIZADO
- ✅ No hay "commits adelante" confusos
- ✅ `git diff main dev` siempre vacío después de sync

### 3. **Simplicidad**
- ✅ Historial lineal, fácil de leer
- ✅ Sin merges bidireccionales
- ✅ Sin commits de sync que crean ruido

### 4. **Verificable**
```bash
# Script de verificación:
MAIN_SHA=$(git rev-parse main)
DEV_SHA=$(git rev-parse dev)

if [ "$MAIN_SHA" = "$DEV_SHA" ]; then
  echo "✅ Sincronizadas"
else
  echo "❌ ALERTA: Divergencia detectada"
fi
```

---

## 🔧 Implementación

### Paso 1: Limpiar Estado Actual

**Opción A: Reset dev a main (Recomendado)**
```bash
# Hacer backup por si acaso
git branch dev-backup dev

# Reset dev a main
git checkout dev
git reset --hard origin/main
git push --force origin dev

# Verificar
git log --oneline main..dev  # ← Debe estar vacío
```

**Opción B: Merge final y luego reset**
```bash
# Si hay commits en dev que quieres preservar
git checkout main
git merge dev  # Último merge
git push origin main

# Luego reset dev
git checkout dev
git reset --hard main
git push --force origin dev
```

### Paso 2: Actualizar Workflows

**Eliminar**: `sync-main-to-dev.yml` (ya eliminado)  
**Agregar**: `sync-main-to-dev-ff.yml` (ya creado)

**Características del nuevo workflow**:
- ✅ Solo hace fast-forward (sin merge commits)
- ✅ Falla si dev tiene commits que main no tiene
- ✅ Verifica que main y dev queden con mismo SHA
- ✅ Genera reporte de sincronización

### Paso 3: Actualizar Documentación

Actualizar `.github/workflows/README.md` con:
- Nuevo flujo de trabajo
- Garantías de sincronización
- Comandos de verificación

---

## 📊 Comparación: Antes vs Después

### Antes (Actual):
```bash
git log --oneline main..dev
# 107b55c chore: sync main v0.1.6 to dev
# f5dc923 chore: sync main v0.1.5 to dev
# f13c21a Feature/infrastructure bootstrap
# ... 31 commits

git diff main dev --stat
# (vacío, pero no es obvio)
```

**Problema**: Contenido igual, historial confuso.

### Después (Propuesto):
```bash
git log --oneline main..dev
# (vacío) ← CLARO

git diff main dev --stat
# (vacío) ← CLARO

git rev-parse main
# f0f9a631c2fc2cbf82be297d175ec202a55b39f9
git rev-parse dev
# f0f9a631c2fc2cbf82be297d175ec202a55b39f9 ← MISMO
```

**Resultado**: Contenido igual, historial igual, SHA igual.

---

## ⚠️ Consideraciones

### 1. Force Push Inicial
- Necesario para limpiar el historial actual
- Solo una vez, después no será necesario

### 2. Coordinación del Equipo
- Avisar antes de hacer el reset
- Todos deben hacer `git pull --rebase` después

### 3. PRs Abiertos
- Verificar que no haya PRs abiertos a dev
- Cerrar o mergear antes del reset

---

## 🚀 Plan de Migración

### Fase 1: Preparación (Ahora)
- [x] Crear workflow `sync-main-to-dev-ff.yml`
- [x] Documentar estrategia
- [ ] Revisar y aprobar propuesta

### Fase 2: Limpieza (Próximo)
- [ ] Verificar que no hay PRs abiertos
- [ ] Hacer backup: `git branch dev-backup dev`
- [ ] Reset dev a main
- [ ] Verificar sincronización

### Fase 3: Validación
- [ ] Probar flujo con feature pequeña
- [ ] Verificar que sync automático funciona
- [ ] Actualizar documentación

### Fase 4: Adopción
- [ ] Comunicar nuevo flujo al equipo
- [ ] Actualizar guías de contribución
- [ ] Monitorear primeros releases

---

## 📝 Comandos de Verificación

### Verificar Sincronización:
```bash
# Opción 1: Comparar SHAs
git rev-parse main
git rev-parse dev
# Deben ser idénticos

# Opción 2: Ver diferencias
git log --oneline main..dev
# Debe estar vacío

# Opción 3: Diff de contenido
git diff main dev
# Debe estar vacío
```

### Script de Verificación Automática:
```bash
#!/bin/bash
# verify-sync.sh

MAIN_SHA=$(git rev-parse origin/main)
DEV_SHA=$(git rev-parse origin/dev)

echo "main: $MAIN_SHA"
echo "dev:  $DEV_SHA"

if [ "$MAIN_SHA" = "$DEV_SHA" ]; then
  echo "✅ SINCRONIZADAS"
  exit 0
else
  echo "❌ DIVERGENCIA DETECTADA"
  echo ""
  echo "Commits en dev que NO están en main:"
  git log --oneline origin/main..origin/dev
  echo ""
  echo "Commits en main que NO están en dev:"
  git log --oneline origin/dev..origin/main
  exit 1
fi
```

---

## 🎯 Resultado Final

### Garantías:
1. **Transparencia**: Historial idéntico, verificable visualmente
2. **Confianza**: Mismo SHA = mismo contenido, sin ambigüedad
3. **Simplicidad**: Sin merges bidireccionales, sin commits de sync
4. **Verificable**: Scripts automáticos pueden validar sincronización

### Flujo Claro:
```
feature → dev (desarrollo)
dev → main (release)
main → dev (sync automático, fast-forward)
```

### Sin Confusión:
- ✅ main = cara hacia afuera (producción)
- ✅ dev = cara hacia desarrollo (trabajo diario)
- ✅ Siempre sincronizados después de release
- ✅ Historial limpio y lineal

---

## ❓ Preguntas Frecuentes

### ¿Qué pasa si dev tiene commits que main no tiene?
El workflow falla y requiere intervención manual. Esto es BUENO porque:
- Te alerta de divergencia
- Evita pérdida accidental de trabajo
- Requiere decisión consciente

### ¿Puedo seguir trabajando en dev mientras espero release?
Sí, pero:
- Crea feature branches desde dev
- No hagas commit directo a dev
- Espera a que se sincronice después del release

### ¿Qué pasa con hotfixes?
Hotfix en main → sync automático a dev (fast-forward)
Todo transparente y verificable.

---

**¿Proceder con esta estrategia?**
