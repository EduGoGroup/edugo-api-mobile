# Templates de Workflows - edugo-api-mobile

Templates de workflows generados durante Sprint 4 para implementar workflows reusables.

---

## 📋 Workflows en Este Directorio

### Workflows Reusables (en infrastructure)

Estos workflows se crean en `edugo-infrastructure/.github/workflows/reusable/`:

1. **pr-validation.yml**
   - Tipo: Reusable
   - Función: Validación completa de PRs (lint, test, security, build)
   - Parámetros: 6 inputs configurables
   - Usado por: pr-to-dev.yml, pr-to-main.yml

2. **sync-branches.yml**
   - Tipo: Reusable
   - Función: Sincronización automática main→dev
   - Parámetros: 3 inputs configurables
   - Usado por: sync-main-to-dev.yml

### Workflows Caller (en api-mobile)

Estos workflows llaman a los reusables:

3. **pr-to-dev.yml**
   - Tipo: Caller
   - Llama: pr-validation.yml
   - Reducción: ~150 líneas → ~15 líneas (-90%)

4. **pr-to-main.yml**
   - Tipo: Caller
   - Llama: pr-validation.yml (con security scan)
   - Reducción: ~180 líneas → ~18 líneas (-90%)

5. **sync-main-to-dev.yml**
   - Tipo: Caller
   - Llama: sync-branches.yml
   - Reducción: ~80 líneas → ~10 líneas (-87%)

---

## 📦 Workflows NO Convertidos

Estos workflows permanecen personalizados en api-mobile:

- **manual-release.yml** - Lógica de release específica del proyecto
- **test.yml** - Tests manuales con opciones específicas

---

## 🚀 Estado de Generación

| Workflow | Estado | Ubicación Definitiva |
|----------|--------|---------------------|
| pr-validation.yml | ⏳ Por crear | `edugo-infrastructure/.github/workflows/reusable/` |
| sync-branches.yml | ⏳ Por crear | `edugo-infrastructure/.github/workflows/reusable/` |
| pr-to-dev.yml | ⏳ Por crear | `edugo-api-mobile/.github/workflows/` |
| pr-to-main.yml | ⏳ Por crear | `edugo-api-mobile/.github/workflows/` |
| sync-main-to-dev.yml | ⏳ Por crear | `edugo-api-mobile/.github/workflows/` |

**Nota:** Los templates están documentados en SPRINT-4-TASKS.md. Se crearán durante la ejecución del Sprint 4.

---

## 📖 Cómo Usar Este Directorio

### Durante Sprint 4

1. **Día 1:** Crear workflows reusables en infrastructure
   - Copiar templates de SPRINT-4-TASKS.md Tareas 4.2-4.3
   - Guardar en infrastructure

2. **Día 2:** Crear workflows caller en api-mobile
   - Copiar templates de SPRINT-4-TASKS.md Tareas 4.6-4.8
   - Guardar en api-mobile

3. **Día 3:** Validar que funcionan correctamente

---

## 🔄 Flujo de Trabajo

```
┌─────────────────────────────────────┐
│  edugo-infrastructure               │
│  .github/workflows/reusable/        │
│                                     │
│  ├── pr-validation.yml ←──────┐    │
│  │   (define lógica)          │    │
│  │                            │    │
│  └── sync-branches.yml ←──────┼──┐ │
│      (define lógica)          │  │ │
└──────────────────────────────┼┼──┼─┘
                               ││  │
                     usa       ││  │  usa
                               ││  │
┌──────────────────────────────┼┼──┼─┐
│  edugo-api-mobile            ││  │ │
│  .github/workflows/          ││  │ │
│                              ││  │ │
│  ├── pr-to-dev.yml ──────────┘│  │ │
│  │   (solo config)            │  │ │
│  │                             │  │ │
│  ├── pr-to-main.yml ───────────┘  │ │
│  │   (solo config)                │ │
│  │                                │ │
│  └── sync-main-to-dev.yml ────────┘ │
│      (solo config)                  │
└─────────────────────────────────────┘
```

---

## 📊 Métricas de Mejora

| Workflow | Antes | Después | Reducción |
|----------|-------|---------|-----------|
| pr-to-dev.yml | ~150 líneas | ~15 líneas | 90% |
| pr-to-main.yml | ~180 líneas | ~18 líneas | 90% |
| sync-main-to-dev.yml | ~80 líneas | ~10 líneas | 87% |
| **TOTAL** | **~410 líneas** | **~43 líneas** | **~90%** |

**Beneficios adicionales:**
- ✅ Centralización (cambios en 1 lugar)
- ✅ Consistencia (mismo comportamiento)
- ✅ Mantenibilidad (más fácil de mantener)
- ✅ Escalabilidad (fácil agregar proyectos)

---

## 📚 Referencias

- **Documentación completa:** SPRINT-4-TASKS.md
- **GitHub Docs:** [Reusable Workflows](https://docs.github.com/en/actions/using-workflows/reusing-workflows)
- **Ejemplo práctico:** Ver edugo-shared (ya implementado)

---

**Última actualización:** 19 de Noviembre, 2025  
**Generado por:** Claude Code
