# 🎯 COMIENZA AQUÍ - edugo-api-mobile

⚠️ **UBICACIÓN Y CONTEXTO DE TRABAJO:**

```
┌─────────────────────────────────────────────────────────────┐
│ 📍 Estás en: 03-api-mobile/ (plan de implementación)      │
│ 📂 Ruta: implementation-plans/03-api-mobile/              │
│ ⚠️ NO uses archivos de otros proyectos (01, 02, etc.)    │
│ ✅ SOLO usa archivos dentro de 03-api-mobile/            │
└─────────────────────────────────────────────────────────────┘
```

**Última actualización:** 20 Nov 2025

---

## 🗺️ MAPA DE UBICACIÓN

```
00-Projects-Isolated/cicd-analysis/implementation-plans/
│
├── 01-shared/                                  ← Otro proyecto
├── 02-infrastructure/                          ← Otro proyecto
│
└── 03-api-mobile/                              ← 👉 ESTÁS AQUÍ
    ├── START-HERE.md                           ← Este archivo
    ├── INDEX.md                                ← Navegación completa
    ├── PROMPTS.md                              ← Prompts para cada fase
    ├── README.md                               ← Contexto del proyecto
    ├── docs/                                   ← Documentación
    ├── sprints/                                ← ⭐ Planes de sprint
    │   ├── SPRINT-2-TASKS.md                   ← Migración + Optimización
    │   ├── SPRINT-4-TASKS.md                   ← Workflows Reusables
    │   └── SPRINT-ENTITIES-ADAPTATION.md
    ├── tracking/                               ← Estado y seguimiento
    │   ├── SPRINT-STATUS.md                    ← Estado actual
    │   ├── REGLAS.md                           ← Reglas de ejecución
    │   └── PR-TEMPLATE.md                      ← Template de PR
    └── assets/                                 ← Scripts y recursos
        ├── scripts/
        └── workflows/
```

---

## 🎯 ¿QUÉ QUIERES HACER?

### 🔍 Opción 1: Ver Estado Actual del Proyecto
```bash
cat tracking/SPRINT-STATUS.md | head -40
```

Lee el archivo para saber:
- Sprint activo
- Fase actual (1, 2, o 3)
- Progreso
- Próxima tarea

### ▶️ Opción 2: Continuar Trabajo desde donde quedó

**Prompt a usar:**
```
Continúa el trabajo de CI/CD en edugo-api-mobile desde donde quedó.
```

Ver detalles en: [PROMPTS.md](PROMPTS.md#continuar-desde-donde-quedó)

### 🆕 Opción 3: Iniciar Nuevo Sprint

**Prompt a usar:**
```
Ejecuta FASE 1 del SPRINT-X en edugo-api-mobile.
```

Reemplaza X con: 2 o 4  
Ver detalles en: [PROMPTS.md](PROMPTS.md#fase-1)

### 📚 Opción 4: Entender el Sistema Completo

**Lee en orden:**
1. [INDEX.md](INDEX.md) - Navegación general (5 min)
2. [README.md](README.md) - Contexto del proyecto (15 min)
3. [tracking/REGLAS.md](tracking/REGLAS.md) - Reglas detalladas (15 min)

---

## 📍 NAVEGACIÓN RÁPIDA

| Quiero... | Abrir... |
|-----------|----------|
| 🗺️ Navegar el proyecto | [INDEX.md](INDEX.md) |
| 🎯 Prompts para ejecutar | [PROMPTS.md](PROMPTS.md) ⭐ |
| 📊 Estado actual | [tracking/SPRINT-STATUS.md](tracking/SPRINT-STATUS.md) |
| 📜 Reglas de ejecución | [tracking/REGLAS.md](tracking/REGLAS.md) |
| 📖 Contexto del proyecto | [README.md](README.md) |
| 🎯 Ver tareas del sprint | [sprints/](sprints/) |
| 📈 Ver progreso | [tracking/SPRINT-STATUS.md](tracking/SPRINT-STATUS.md) |

---

## 🤖 PARA CLAUDE CODE (INSTRUCCIONES CRÍTICAS)

### ⚠️ Antes de Hacer CUALQUIER COSA:

1. **Lee SIEMPRE:** `INDEX.md`
2. **Verifica ubicación:**
   ```bash
   pwd
   # Debe contener: /03-api-mobile
   ```
3. **Lee estado:** `tracking/SPRINT-STATUS.md`
4. **Identifica:**
   - Sprint activo
   - Fase actual
   - Próxima tarea

### ⚠️ NO Uses Archivos Fuera de 03-api-mobile/

**Archivos PROHIBIDOS:**
- ❌ `/01-shared/*` (otro proyecto)
- ❌ `/02-infrastructure/*` (otro proyecto)
- ❌ `/04-api-admin/*` (otro proyecto)
- ❌ Cualquier archivo fuera de 03-api-mobile/

**Archivos PERMITIDOS:**
- ✅ `sprints/SPRINT-X-TASKS.md`
- ✅ `tracking/*`
- ✅ `docs/*`

### ⚠️ Cómo Verificar que Estás en el Archivo Correcto:

```bash
# Al abrir un archivo de sprint, verifica:
readlink -f sprints/SPRINT-2-TASKS.md
# Debe mostrar: .../03-api-mobile/sprints/SPRINT-2-TASKS.md

# Si muestra otra ruta, estás en el lugar equivocado
```

---

## 🔗 Enlaces Importantes

- **Plan general del proyecto:** [README.md](README.md)
- **Navegación completa:** [INDEX.md](INDEX.md)
- **Prompts para ejecutar:** [PROMPTS.md](PROMPTS.md) ⭐
- **Estado en tiempo real:** [tracking/SPRINT-STATUS.md](tracking/SPRINT-STATUS.md)
- **Reglas de ejecución:** [tracking/REGLAS.md](tracking/REGLAS.md)

---

## 📊 COMANDOS RÁPIDOS

### Ver estado actual:
```bash
cd /Users/jhoanmedina/source/EduGo/Analisys/00-Projects-Isolated/cicd-analysis/implementation-plans/03-api-mobile
cat tracking/SPRINT-STATUS.md | head -40
```

### Ver próxima tarea:
```bash
grep "⏳\|🔄" tracking/SPRINT-STATUS.md | head -1
```

### Ver sprints disponibles:
```bash
ls -1 sprints/
```

### Ver logs de sesiones anteriores:
```bash
ls -lt tracking/logs/ | head -5
```

---

## 🎉 ¡Listo para Comenzar!

Has llegado al final del índice. Ahora tienes:
- ✅ Visión completa del proyecto
- ✅ Entiendes por qué api-mobile es el PILOTO
- ✅ Sabes qué sprints hay y en qué orden
- ✅ Conoces las rutas según tu rol

**Siguiente paso recomendado:**
```bash
open README.md
# Leer contexto completo (15-20 min)
```

O si ya estás listo:
```bash
open sprints/SPRINT-2-TASKS.md
# Ir directamente a implementación
```

---

## 🔄 Dependencias Entre Sprints

```
Sprint 1 (shared)
    ↓ (completado previamente)
    ↓
Sprint 2 (api-mobile) ← ESTAMOS AQUÍ
    ↓ (migración + optimización)
    ↓
Sprint 3 (api-admin, worker)
    ↓ (replicar patrón validado)
    ↓
Sprint 4 (infrastructure + reusables) ← LUEGO AQUÍ
    ↓ (centralización)
    ↓
Sprint 5+ (todos)
    (mantenimiento)
```

---

## 📝 Notas Importantes

### ⚠️ Antes de Ejecutar Cualquier Script

1. **Leer el script completo**
2. **Verificar rutas** (ajustar si es necesario)
3. **Ejecutar en rama de desarrollo**, NO en main
4. **Hacer backup** antes de cambios grandes
5. **Validar resultado** antes de commit

### ⚠️ Sobre el Paralelismo

- Funciona muy bien en GitHub Actions
- Ahorra tiempo, pero consume más recursos
- Validar que no agota límites de plan

### ⚠️ Sobre Pre-commit Hooks

- Son locales, cada dev debe configurar
- Agregar a documentación de onboarding
- No son obligatorios, pero muy recomendados

---

**Generado por:** Claude Code  
**Fecha:** 20 de Noviembre, 2025  
**Versión:** 1.0  
**Proyecto:** edugo-api-mobile (PILOTO)
