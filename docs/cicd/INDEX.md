# Índice - Plan de Implementación edugo-api-mobile

**🎯 Proyecto PILOTO para Optimización de CI/CD**

---

## 🗺️ Navegación Rápida

### Para Empezar
1. **[README.md](./README.md)** ⭐ - Contexto completo del proyecto (15-20 min)
2. **Este archivo (INDEX.md)** - Navegación rápida (5 min)

### Para Implementar
3. **[SPRINT-2-TASKS.md](./SPRINT-2-TASKS.md)** ⭐⭐⭐ - Sprint 2: Migración Go 1.25 + Optimización
4. **[SPRINT-4-TASKS.md](./SPRINT-4-TASKS.md)** ⭐⭐ - Sprint 4: Workflows Reusables

---

## 📊 Resumen Ultra-Rápido

```
Proyecto: edugo-api-mobile
Tipo: A (API desplegable con Docker)
Puerto: 8080
Success Rate: 90% (9/10) ✅ Muy Bueno

Prioridades:
├── Sprint 2: MIGRACIÓN + OPTIMIZACIÓN
│   ├── 🟡 P1: Migrar a Go 1.25 (PILOTO)
│   ├── 🟡 P1: Implementar paralelismo
│   ├── 🟡 P1: Pre-commit hooks
│   ├── 🟢 P2: Corregir 23 errores lint
│   ├── 🟢 P2: Control releases por variable
│   └── 🟢 P2: Mejorar coverage reporting
│   ⏱️ Estimado: 3-4 días / 12-16 horas
│
└── Sprint 4: WORKFLOWS REUSABLES
    ├── 🟢 P2: Crear workflows reusables base
    ├── 🟢 P2: Migrar api-mobile a reusables
    ├── 🟢 P2: Validar en staging
    └── 🟢 P2: Documentar patrón
    ⏱️ Estimado: 3-4 días / 12-15 horas

Total Estimado: 24-31 horas en 6-8 días
```

---

## 🚀 ¿Por Qué api-mobile es el PILOTO?

### Ventajas Clave

1. **✅ Ya está muy bien estructurado**
   - Success rate: 90% (el mejor después de shared)
   - 5 workflows bien organizados
   - Tests de integración con testcontainers
   - Security scan implementado
   - GitHub App tokens en uso

2. **✅ Menor riesgo de fallos**
   - Solo 1 fallo en últimas 10 ejecuciones
   - Tests confiables
   - Docker builds estables

3. **✅ Es el más representativo**
   - Tiene todos los workflows necesarios
   - Patrón aplicable a api-administracion
   - Usa todas las mejores prácticas

4. **✅ Validación rápida**
   - Ciclos de CI rápidos (~2-5 min)
   - Fácil detectar problemas temprano

### Lo Que Validaremos Aquí

- ✅ Go 1.25 funciona en CI/CD (ya validado localmente)
- ✅ Paralelismo mejora tiempos sin romper nada
- ✅ Pre-commit hooks son útiles sin ser molestos
- ✅ Workflows reusables son mantenibles

**Una vez validado aquí → replicar a api-administracion y worker**

---

## 📁 Estructura de Archivos

```
03-api-mobile/
├── INDEX.md                    ← Estás aquí
├── README.md                   ← Contexto del proyecto
├── SPRINT-2-TASKS.md          ← ⭐ Sprint 2 completo
├── SPRINT-4-TASKS.md          ← Sprint 4 completo
├── SCRIPTS/                    ← Scripts bash reutilizables
│   ├── migrate-go-1.25.sh
│   ├── setup-precommit.sh
│   ├── validate-workflows.sh
│   └── README.md
└── WORKFLOWS/                  ← Templates de workflows
    ├── pr-to-dev.yml
    ├── pr-to-main.yml
    ├── manual-release.yml
    ├── sync-main-to-dev.yml
    └── test.yml
```

---

## 🎯 Por Rol

### Soy el Implementador
→ **Ruta:** INDEX.md → README.md → SPRINT-2-TASKS.md  
→ **Ejecuto:** Tareas una por una, validando en cada paso  
→ **Tiempo:** 12-16 horas Sprint 2 (3-4 días)

### Soy el DevOps Lead
→ **Ruta:** README.md → SPRINT-2-TASKS.md (estructura) → SPRINT-4-TASKS.md  
→ **Reviso:** Estimaciones, riesgos, estrategia  
→ **Tiempo:** 1-2 horas de lectura

### Soy el QA/Tester
→ **Ruta:** README.md → Secciones de validación en cada SPRINT  
→ **Valido:** Tests pasan, cobertura mantiene, CI no rompe  
→ **Tiempo:** 30-60 min por PR

### Quiero Replicar en api-administracion
→ **Ruta:** README.md → SPRINT-2-TASKS.md completo  
→ **Adapto:** Scripts y comandos (cambiar rutas)  
→ **Tiempo:** 10-12 horas (más rápido, patrón ya validado)

---

## 📈 Roadmap de Lectura

### Nivel 1: Overview (20 min)
1. **INDEX.md** (este archivo) - 5 min
2. **README.md** (secciones resumen) - 10 min
3. **SPRINT-2-TASKS.md** (solo índice) - 5 min

### Nivel 2: Preparación (1 hora)
1. **README.md** completo - 20 min
2. **SPRINT-2-TASKS.md** (estructura + Día 1) - 30 min
3. **SCRIPTS/** (revisar scripts disponibles) - 10 min

### Nivel 3: Implementación (2-3 horas lectura + ejecución)
1. **SPRINT-2-TASKS.md** completo - 1-2 horas
2. Ejecutar tareas mientras lees - 1 hora
3. Validar resultados - 30 min

---

## 🔥 Top 5 Tareas Críticas

Si solo tienes tiempo limitado, prioriza:

### 1. **Migrar a Go 1.25** (60 min) 🟡 P1
- **Archivo:** SPRINT-2-TASKS.md → Tarea 2.1
- **Por qué:** PILOTO, validar aquí primero
- **Riesgo:** Bajo (ya validado localmente)
- **Impacto:** Alto (última versión, mejoras performance)

### 2. **Implementar paralelismo** (90 min) 🟡 P1
- **Archivo:** SPRINT-2-TASKS.md → Tarea 2.2
- **Por qué:** Reducir tiempos de CI ~30%
- **Riesgo:** Bajo (APIs de GitHub estables)
- **Impacto:** Alto (ahorro de tiempo)

### 3. **Pre-commit hooks** (60-90 min) 🟡 P1
- **Archivo:** SPRINT-2-TASKS.md → Tarea 2.3
- **Por qué:** Prevenir errores antes de push
- **Riesgo:** Bajo
- **Impacto:** Alto (calidad de código)

### 4. **Corregir 23 errores lint** (45 min) 🟢 P2
- **Archivo:** SPRINT-2-TASKS.md → Tarea 2.4
- **Por qué:** Limpieza de código, CI más limpio
- **Riesgo:** Muy bajo
- **Impacto:** Medio (calidad)

### 5. **Control releases por variable** (30 min) 🟢 P2
- **Archivo:** SPRINT-2-TASKS.md → Tarea 2.5
- **Por qué:** Evitar releases accidentales
- **Riesgo:** Muy bajo
- **Impacto:** Medio (control)

**Total Top 5:** ~5-6 horas (en lugar de 12-16h completo)

---

## 💾 Estado Actual vs Objetivo

### Estado Actual (api-mobile)

```yaml
Go Version: 1.24.10
Workflows: 5
  - pr-to-dev.yml
  - pr-to-main.yml
  - test.yml (manual)
  - manual-release.yml
  - sync-main-to-dev.yml
  
Success Rate: 90%
Lint Errors: 23 (20 errcheck + 3 govet)
Paralelismo: No
Pre-commit: No
Tests Integración: Sí (testcontainers)
Security Scan: Sí (Gosec)
GitHub App Token: Sí (solo release)
Coverage Threshold: 33%
```

### Estado Objetivo (Post Sprint 2)

```yaml
Go Version: 1.25 ✅
Workflows: 5 (mismos)
  
Success Rate: >95%
Lint Errors: 0 ✅
Paralelismo: Sí ✅
Pre-commit: Sí ✅
Tests Integración: Sí (sin cambios)
Security Scan: Sí (mejorado)
GitHub App Token: Sí (en más lugares)
Coverage Threshold: 33% (reportes mejorados)
Control Releases: Por variable ✅
```

### Estado Objetivo (Post Sprint 4)

```yaml
Workflows: 5 → 4 (usando reusables)
  - pr-to-dev.yml (llamando reusable)
  - pr-to-main.yml (llamando reusable)
  - manual-release.yml (personalizado)
  - sync-main-to-dev.yml (llamando reusable)
  
Duplicación: -60% en código
Mantenibilidad: +80%
Reusabilidad: Base para api-admin y worker
```

---

## 🆘 Ayuda Rápida

### ¿Por dónde empiezo?
**Respuesta:** README.md → SPRINT-2-TASKS.md línea ~100 (Tarea 2.1)

### ¿Cuánto tiempo necesito?
**Respuesta:**
- Sprint 2 completo: 12-16h en 3-4 días
- Sprint 2 modo rápido (Top 5): 5-6h en 2 días
- Sprint 4 completo: 12-15h en 3-4 días

### ¿Puedo saltar Sprint 2 e ir directo a Sprint 4?
**Respuesta:** **NO**. Sprint 4 depende de Sprint 2. Primero optimizar, luego reutilizar.

### ¿Los scripts funcionan?
**Respuesta:** Sí, diseñados para copiar/pegar. Ver `/SCRIPTS/` para todos los disponibles.

### ¿Qué hago si Go 1.25 falla en CI?
**Respuesta:** Ver SPRINT-2-TASKS.md → Tarea 2.1 → "Solución de Problemas". Incluye rollback automático.

### ¿Debo hacer PR por cada tarea?
**Respuesta:** No. Ver estrategia de commits en SPRINT-2-TASKS.md. Se agrupa lógicamente.

### ¿Cómo valido que no rompí nada?
**Respuesta:** Cada tarea tiene sección "Criterios de Validación" + "Checkpoint". Ejecutar antes de continuar.

---

## 📞 Referencias Externas

### Documentación Base
- [Análisis Estado Actual](../../01-ANALISIS-ESTADO-ACTUAL.md)
- [Matriz Comparativa](../../04-MATRIZ-COMPARATIVA.md)
- [Quick Wins](../../05-QUICK-WINS.md)
- [Pruebas Go 1.25](../../08-RESULTADO-PRUEBAS-GO-1.25.md) ✅

### Repositorio
- **URL:** https://github.com/EduGoGroup/edugo-api-mobile
- **Ruta Local:** `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile`
- **Branch Principal:** `main`
- **Branch Desarrollo:** `dev`

### Ejemplo de Referencia (shared)
- [Plan shared](../01-shared/) - Proyecto anterior con estructura similar

---

## ✅ Checklist Pre-Lectura

Antes de comenzar:
- [x] Estás en el directorio correcto
- [x] Tienes acceso al repo local
- [ ] Has leído el análisis general en `00-RESUMEN-EJECUTIVO.md`
- [ ] Tienes tiempo para leer (mínimo 20 min)
- [ ] Editor de markdown disponible
- [ ] Terminal lista para ejecutar comandos
- [ ] Decidido en qué rol estás (implementador/lead/tester)

---

## 🎯 Próxima Acción

```bash
# Opción A: Comenzar a implementar HOY
cd /Users/jhoanmedina/source/EduGo/Analisys/00-Projects-Isolated/cicd-analysis/implementation-plans/03-api-mobile
open SPRINT-2-TASKS.md

# Opción B: Solo entender el contexto
open README.md

# Opción C: Ver solo los scripts
cd SCRIPTS/
ls -la
cat README.md

# Opción D: Validar que tengo todo
./SCRIPTS/validate-prerequisites.sh
```

---

## 📊 Métricas del Plan

| Métrica | Valor |
|---------|-------|
| Archivos principales | 4 markdown |
| Scripts incluidos | ~8-10 bash scripts |
| Tareas Sprint 2 | ~15 tareas |
| Tareas Sprint 4 | ~12 tareas |
| Tiempo estimado total | 24-31 horas |
| Sprints cubiertos | 2 (Sprint 2 y Sprint 4) |
| Nivel de detalle | Ultra-alto |
| PRs esperados | 2-3 PRs |

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
open SPRINT-2-TASKS.md
# Ir directamente a Tarea 2.1: Migrar a Go 1.25
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
**Fecha:** 19 de Noviembre, 2025  
**Versión:** 1.0  
**Proyecto:** edugo-api-mobile (PILOTO)
