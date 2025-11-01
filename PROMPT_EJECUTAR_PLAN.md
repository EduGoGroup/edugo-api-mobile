# 🚀 Prompt para Ejecutar Plan CI/CD en Proyectos Hermanos

## 📋 Contexto

Este prompt está diseñado para iniciar una nueva conversación de Claude Code que ejecute el plan documentado en `PLAN_CICD_PROYECTOS_HERMANOS.md`.

---

## 🎯 PROMPT PARA NUEVA CONVERSACIÓN

Copia y pega el siguiente texto en una **nueva conversación** de Claude Code:

```
Hola, necesito que ejecutes el plan de trabajo para replicar la configuración de CI/CD y GitHub Copilot en los proyectos hermanos de EduGo.

📚 CONTEXTO IMPORTANTE:

1. El plan completo está en: /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/PLAN_CICD_PROYECTOS_HERMANOS.md

2. Proyecto origen (ya completado):
   - Ruta: /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile
   - Branch: dev
   - Estado: ✅ Completado con Copilot instructions, workflows optimizados y validación con actionlint

3. Proyectos pendientes (en orden de prioridad):
   a) edugo-shared (CRÍTICO - empezar por este)
      Ruta: /Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared
      Tipo: Librería Go compartida
      Workflows: CI, Tests, Release (SIN Docker)

   b) edugo-api-administracion
      Ruta: /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion
      Tipo: API REST (Go)
      Workflows: Completo (6 workflows + Docker)

   c) edugo-worker
      Ruta: /Users/jhoanmedina/source/EduGo/repos-separados/edugo-worker
      Tipo: Worker (Go)
      Workflows: 4 workflows + Docker (sin auto-version ni sync)

   d) edugo-dev-environment
      Ruta: /Users/jhoanmedina/source/EduGo/repos-separados/edugo-dev-environment
      Tipo: Docker Compose
      Workflows: Solo validación simple

🎯 OBJETIVO:

Ejecutar la FASE 3 del plan (edugo-shared) siguiendo estos pasos:

1. Leer el plan completo en PLAN_CICD_PROYECTOS_HERMANOS.md
2. Instalar actionlint si no está instalado (brew install actionlint)
3. Cambiar al directorio de edugo-shared
4. Seguir los pasos 3.1 a 3.12 del plan (FASE 3)
5. Validar workflows con actionlint ANTES de cada commit
6. Crear PR a dev
7. Verificar ejecución correcta de workflows
8. Documentar hallazgos en el plan

⚠️ REGLAS CRÍTICAS:

- SIEMPRE validar workflows con actionlint antes de push
- NO usar heredocs con backticks en workflows
- Usar múltiples -m flags para commits multilinea
- Crear copilot-instructions.md enfocado en librería (no API)
- edugo-shared NO debe tener workflows de Docker
- Enfocarse en retrocompatibilidad y semantic versioning
- Todos los comentarios y mensajes en ESPAÑOL

📝 DOCUMENTACIÓN CLAVE:

- Plan completo: PLAN_CICD_PROYECTOS_HERMANOS.md
- Errores comunes: Sección "Herramientas de Validación Pre-Commit"
- Workflows de referencia: edugo-api-mobile/.github/workflows/
- Copilot instructions base: edugo-api-mobile/.github/copilot-instructions.md

🚀 COMIENZA:

Por favor, confirma que entiendes el plan y comienza con el análisis previo de edugo-shared (Paso 3.1).
```

---

## 🔄 PROMPT ALTERNATIVO (Más Conciso)

Si prefieres un prompt más corto:

```
Ejecuta la FASE 3 del plan CI/CD en: /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/PLAN_CICD_PROYECTOS_HERMANOS.md

Proyecto objetivo: edugo-shared
Ruta: /Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared

Sigue los pasos 3.1 a 3.12 del plan. IMPORTANTE:
- Validar workflows con actionlint
- NO Docker (es librería)
- Enfoque en semantic versioning
- Todo en español

Comienza con análisis previo (Paso 3.1).
```

---

## 📊 Estado Actual

| Proyecto | Estado | Branch Actual |
|----------|--------|---------------|
| edugo-api-mobile | ✅ Completado | dev |
| edugo-shared | ⏳ Siguiente | - |
| edugo-api-administracion | ⏸️ Pendiente | - |
| edugo-worker | ⏸️ Pendiente | - |
| edugo-dev-environment | ⏸️ Pendiente | - |

---

## 🛠️ Pre-requisitos Instalados

- ✅ `actionlint` - Validación de workflows
- ✅ `gh` - GitHub CLI
- ✅ `git` - Control de versiones

---

## 📚 Archivos de Referencia

Los siguientes archivos del proyecto origen están listos para copiar/adaptar:

```
edugo-api-mobile/
├── .github/
│   ├── copilot-instructions.md      # Adaptar para cada tipo de proyecto
│   └── workflows/
│       ├── ci.yml                    # Base para todos
│       ├── test.yml                  # Base para todos
│       ├── release.yml               # Adaptar según proyecto
│       ├── auto-version.yml          # Solo APIs
│       ├── sync-main-to-dev.yml      # Solo APIs
│       └── docker-only.yml           # Solo proyectos con Docker
└── PLAN_CICD_PROYECTOS_HERMANOS.md  # Este plan
```

---

## ⚠️ Notas Importantes

1. **edugo-shared es CRÍTICO** - Todos los demás proyectos dependen de él
2. **Release workflow en shared es ESENCIAL** - Debe crear GitHub Releases con tags
3. **NO agregar Docker workflows** a edugo-shared
4. **Validar SIEMPRE con actionlint** antes de push
5. **Documentar cambios** en el plan después de cada fase

---

## 🎯 Resultado Esperado

Al finalizar la FASE 3, edugo-shared debe tener:

- ✅ `.github/copilot-instructions.md` (enfocado en librería)
- ✅ `.github/workflows/ci.yml` (tests exhaustivos)
- ✅ `.github/workflows/test.yml` (cobertura >80%)
- ✅ `.github/workflows/release.yml` (crea GitHub Release + tag)
- ✅ `.github/workflows/README.md` (documentación)
- ✅ PR mergeado a dev
- ✅ Workflows validados y funcionando
- ✅ Plan actualizado con estado ✅

---

## 📞 Soporte

Si encuentras problemas durante la ejecución:

1. Consulta la sección "Herramientas de Validación Pre-Commit" del plan
2. Revisa la sección "Lecciones Aprendidas"
3. Compara con los workflows validados de edugo-api-mobile
4. Usa `actionlint` para identificar errores específicos

---

**Creado:** 2025-11-01
**Última actualización:** 2025-11-01
**Autor:** Claude Code
