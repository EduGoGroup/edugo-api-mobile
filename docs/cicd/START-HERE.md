# 🎯 COMIENZA AQUÍ - edugo-api-mobile

⚠️ **UBICACIÓN Y CONTEXTO DE TRABAJO:**

```
┌─────────────────────────────────────────────────────────────┐
│ 📍 Estás en: docs/cicd/ (documentación CI/CD)             │
│ 📂 Ruta: edugo-api-mobile/docs/cicd/                     │
│ ⚠️ Toda la documentación CI/CD está en esta carpeta      │
│ ✅ Usa rutas relativas a docs/cicd/                      │
└─────────────────────────────────────────────────────────────┘
```

**Última actualización:** 20 Nov 2025

---

## 🗺️ MAPA DE UBICACIÓN

```
edugo-api-mobile/
│
└── docs/
    └── cicd/                                   ← 👉 ESTÁS AQUÍ
        ├── START-HERE.md                       ← Este archivo
        ├── INDEX.md                            ← Navegación completa
        ├── PROMPTS.md                          ← Prompts para cada fase
        ├── README.md                           ← Contexto del proyecto
        ├── WORKFLOWS-REUSABLES-GUIDE.md        ← Guía de workflows
        ├── docs/                               ← Documentación adicional
        ├── sprints/                            ← ⭐ Planes de sprint
        │   ├── SPRINT-2-TASKS.md               ← Migración + Optimización
        │   ├── SPRINT-4-TASKS.md               ← Workflows Reusables
        │   └── SPRINT-ENTITIES-ADAPTATION.md
        ├── tracking/                           ← Estado y seguimiento
        │   ├── decisions/                      ← Decisiones técnicas
        │   ├── SPRINT-2-COMPLETE.md
        │   ├── SPRINT-4-METRICAS-FINALES.md
        │   └── FASE-2-VALIDATION.md
        └── assets/                             ← Scripts y recursos
            ├── scripts/
            └── workflows/
```

---

## 🎯 ¿QUÉ QUIERES HACER?

### 🔍 Opción 1: Ver Estado Actual del Proyecto
```bash
cat docs/cicd/tracking/SPRINT-2-COMPLETE.md
cat docs/cicd/tracking/SPRINT-4-METRICAS-FINALES.md
```

Lee los archivos para saber:
- Sprints completados
- Estado de implementación
- Métricas finales
- Próximos pasos

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
| 📊 Estado actual | [tracking/SPRINT-2-COMPLETE.md](tracking/SPRINT-2-COMPLETE.md) |
| 📈 Métricas finales | [tracking/SPRINT-4-METRICAS-FINALES.md](tracking/SPRINT-4-METRICAS-FINALES.md) |
| 📖 Contexto del proyecto | [README.md](README.md) |
| 🎯 Ver tareas del sprint | [sprints/](sprints/) |
| 🔍 Decisiones técnicas | [tracking/decisions/](tracking/decisions/) |

---

## 🤖 PARA CLAUDE CODE (INSTRUCCIONES CRÍTICAS)

### ⚠️ Antes de Hacer CUALQUIER COSA:

1. **Lee SIEMPRE:** `docs/cicd/INDEX.md`
2. **Verifica ubicación:**
   ```bash
   pwd
   # Debe estar en el root del proyecto: edugo-api-mobile/
   ```
3. **Lee estado:** `docs/cicd/tracking/SPRINT-4-METRICAS-FINALES.md`
4. **Identifica:**
   - Sprints completados
   - Estado de workflows
   - Próximas mejoras

### ⚠️ Estructura de la Documentación CI/CD

**Documentación está en:**
- ✅ `docs/cicd/sprints/SPRINT-X-TASKS.md` - Planes de sprint
- ✅ `docs/cicd/tracking/*` - Estado y métricas
- ✅ `docs/cicd/docs/*` - Documentación adicional
- ✅ `docs/cicd/assets/*` - Scripts y recursos

**Código del proyecto está en:**
- ✅ `.github/workflows/*` - Workflows de CI/CD
- ✅ `internal/*` - Código fuente
- ✅ `cmd/*` - Aplicación principal

### ⚠️ Cómo Verificar que Estás en el Archivo Correcto:

```bash
# Al abrir un archivo de documentación CI/CD, verifica:
ls -la docs/cicd/sprints/SPRINT-2-TASKS.md
# Debe existir en: edugo-api-mobile/docs/cicd/sprints/

# Verificar workflows implementados:
ls -la .github/workflows/
```

---

## 🔗 Enlaces Importantes

- **Plan general del proyecto:** [README.md](README.md)
- **Navegación completa:** [INDEX.md](INDEX.md)
- **Prompts para ejecutar:** [PROMPTS.md](PROMPTS.md) ⭐
- **Guía de workflows reusables:** [WORKFLOWS-REUSABLES-GUIDE.md](WORKFLOWS-REUSABLES-GUIDE.md)
- **Estado de sprints:** [tracking/](tracking/)

---

## 📊 COMANDOS RÁPIDOS

### Ver estado actual del proyecto:
```bash
# Ver métricas finales del Sprint 4
cat docs/cicd/tracking/SPRINT-4-METRICAS-FINALES.md

# Ver completitud del Sprint 2
cat docs/cicd/tracking/SPRINT-2-COMPLETE.md
```

### Ver workflows implementados:
```bash
# Listar todos los workflows
ls -1 .github/workflows/

# Ver workflows reusables
ls -1 .github/workflows/reusable-*.yml
```

### Ver sprints disponibles:
```bash
ls -1 docs/cicd/sprints/
```

### Ver decisiones técnicas:
```bash
ls -1 docs/cicd/tracking/decisions/
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
# Leer contexto completo del proyecto
cat docs/cicd/README.md

# Ver guía de workflows reusables
cat docs/cicd/WORKFLOWS-REUSABLES-GUIDE.md

# Ver sprints completados
cat docs/cicd/tracking/SPRINT-2-COMPLETE.md
cat docs/cicd/tracking/SPRINT-4-METRICAS-FINALES.md
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
