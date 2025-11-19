# Plan de Implementación - edugo-api-mobile

**Proyecto PILOTO para Optimización de CI/CD**

---

## 📋 Información del Proyecto

| Campo | Valor |
|-------|-------|
| **Nombre** | edugo-api-mobile |
| **Tipo** | A (API REST desplegable con Docker) |
| **Puerto** | 8080 |
| **Tecnología** | Go + Gin + GORM + Swagger |
| **Base de Datos** | PostgreSQL 15 |
| **Repositorio** | https://github.com/EduGoGroup/edugo-api-mobile |
| **Ruta Local** | `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile` |
| **Workflows Actuales** | 5 |
| **Success Rate** | 90% (9/10 últimas ejecuciones) ✅ |
| **Estado** | ✅ Muy bueno (mejor después de shared) |

---

## 🎯 Por Qué Este Proyecto es el PILOTO

### Razones Estratégicas

#### 1. **Ya Tiene Excelente Base** ✅
- Success rate: 90% - El mejor de los proyectos Tipo A
- Solo 1 fallo en las últimas 10 ejecuciones
- Workflows bien estructurados y organizados
- Tests de integración funcionando (testcontainers)
- Security scan implementado (Gosec)
- GitHub App tokens en uso correcto

#### 2. **Menor Riesgo de Romper Producción** 🛡️
- Tests confiables que detectan problemas temprano
- Docker builds estables
- Ciclos de CI rápidos (~2-5 min)
- Fácil hacer rollback si algo falla

#### 3. **Es el Más Representativo** 📊
- Tiene TODOS los workflows que necesitamos
- Usa todas las mejores prácticas actuales
- Patrón directamente aplicable a:
  - `edugo-api-administracion` (gemelo)
  - `edugo-worker` (similar)

#### 4. **Validación Rápida de Cambios** ⚡
- CI completo en ~5 minutos
- Tests rápidos de ejecutar
- Feedback loop corto

### Lo Que Validaremos Aquí (Para Luego Replicar)

- ✅ **Go 1.25:** Funciona en CI/CD sin problemas (ya validado localmente)
- ✅ **Paralelismo:** Mejora tiempos sin romper tests
- ✅ **Pre-commit hooks:** Son útiles sin ser molestos para developers
- ✅ **Workflows reusables:** Son mantenibles y escalables

**Estrategia:** Una vez validado aquí → replicar confiadamente a los demás proyectos.

---

## 📊 Estado Actual Detallado

### Workflows Existentes (5)

#### 1. **pr-to-dev.yml** - CI para Pull Requests a dev
```yaml
Trigger: Pull request a dev
Jobs: 3
  - lint (golangci-lint)
  - test (unit + integration)
  - build-docker
Duración: ~2 min
Success: ✅ Muy confiable
```

**Fortalezas:**
- ✅ Tests de integración con testcontainers
- ✅ Coverage threshold 33%
- ✅ Reporte de cobertura en PR

**Oportunidades de mejora:**
- 🟡 Agregar paralelismo (lint + test + build en paralelo)
- 🟡 Cache de dependencias Go

#### 2. **pr-to-main.yml** - CI para Pull Requests a main
```yaml
Trigger: Pull request a main
Jobs: 4
  - lint
  - test
  - security-scan (Gosec)
  - build-docker
Duración: ~5 min
Success: ✅ Muy confiable
```

**Fortalezas:**
- ✅ Security scan adicional (Gosec)
- ✅ Validación más estricta

**Oportunidades de mejora:**
- 🟡 Paralelismo (4 jobs pueden correr en paralelo)
- 🟡 Cache más agresivo

#### 3. **test.yml** - Tests manuales
```yaml
Trigger: workflow_dispatch (manual)
Jobs: 1
  - test (con opciones)
Duración: ~2 min
Success: ✅ Confiable
```

**Fortalezas:**
- ✅ Útil para debugging
- ✅ Opciones configurables

**Oportunidades de mejora:**
- 🟡 Agregar opción para skip integration tests

#### 4. **manual-release.yml** - Release manual
```yaml
Trigger: workflow_dispatch (manual)
Jobs: 3
  - validate (tests)
  - release (GitHub release)
  - build-docker (multi-arch)
Duración: ~8 min
Success: ✅ Confiable
Características especiales:
  - GitHub App Token (dispara sync-main-to-dev)
  - Multi-platform (amd64, arm64)
  - Actualiza version.txt
  - Genera changelog
```

**Fortalezas:**
- ✅ Usa GitHub App para disparar workflows subsecuentes
- ✅ Multi-platform Docker builds
- ✅ Control manual evita releases accidentales

**Oportunidades de mejora:**
- 🟡 Agregar variable ENABLE_AUTO_RELEASE para control fino
- 🟡 Validación pre-release más estricta

#### 5. **sync-main-to-dev.yml** - Sincronización automática
```yaml
Trigger: 
  - Push a main
  - Tag v*
Jobs: 1
  - sync
Duración: ~1 min
Success: ⚠️ 1 fallo temporal
```

**Fortalezas:**
- ✅ Mantiene dev actualizado automáticamente
- ✅ Maneja conflictos gracefully

**Oportunidades de mejora:**
- 🟡 Usar GitHub App token también aquí
- 🟡 Mejorar manejo de conflictos

### Configuración Actual

```yaml
Go Version: 1.24.10
golangci-lint: v1.64.7
Docker Registry: ghcr.io
Docker Platforms: linux/amd64, linux/arm64
Coverage Threshold: 33%
Coverage Reporting: ✅ PR comments
Security Scan: ✅ Gosec (solo pr-to-main)
Pre-commit Hooks: ❌ No
Paralelismo: ❌ No
```

### Problemas Conocidos

#### 🔴 Errores de Lint (23 total)
```
20 errores errcheck:
  - defer stmt.Close() sin verificar error
  - defer resp.Body.Close() sin verificar error
  
3 errores govet:
  - Build tags obsoletos (// +build en lugar de //go:build)
```

**Impacto:** Bajo (no bloquean CI porque lint continúa)
**Prioridad:** 🟢 P2 (corregir pero no urgente)

#### 🟡 Versión Go Desactualizada
```
Actual: Go 1.24.10
Disponible: Go 1.25.4
```

**Impacto:** Medio (perdemos mejoras de performance)
**Prioridad:** 🟡 P1 (migrar en Sprint 2)

---

## 🎯 Objetivos de los Sprints

### Sprint 2: Migración + Optimización (Este Sprint)

**Duración:** 3-4 días  
**Esfuerzo:** 12-16 horas  
**Prioridad:** 🟡 P1 (Alta)

#### Objetivos Principales

1. **Migrar a Go 1.25** (PILOTO)
   - Validar en api-mobile primero
   - Si funciona → replicar a demás proyectos
   - Rollback automático si falla

2. **Implementar Paralelismo**
   - Reducir tiempo CI ~30-40%
   - Aprovechar runners de GitHub mejor
   - Mantener confiabilidad

3. **Configurar Pre-commit Hooks**
   - Prevenir errores antes de push
   - Formateo automático
   - Lint local

4. **Corregir Errores de Lint**
   - 23 errores actuales
   - Limpieza de código
   - CI más limpio

5. **Mejorar Control de Releases**
   - Variable ENABLE_AUTO_RELEASE
   - Prevenir releases accidentales
   - Mayor control

#### Resultado Esperado

```yaml
✅ Go 1.25 validado en CI
✅ Tiempos de CI reducidos ~30%
✅ Pre-commit hooks configurados
✅ 0 errores de lint
✅ Control de releases mejorado
✅ Documentación actualizada
✅ Success rate: >95%
```

#### Métricas de Éxito

| Métrica | Antes | Después | Mejora |
|---------|-------|---------|--------|
| **Go Version** | 1.24.10 | 1.25 | ✅ Latest |
| **Tiempo PR→dev** | ~2 min | ~1.5 min | -25% |
| **Tiempo PR→main** | ~5 min | ~3 min | -40% |
| **Errores Lint** | 23 | 0 | -100% |
| **Success Rate** | 90% | >95% | +5% |
| **Pre-commit** | No | Sí | ✅ |

---

### Sprint 4: Workflows Reusables (Sprint Futuro)

**Duración:** 3-4 días  
**Esfuerzo:** 12-15 horas  
**Prioridad:** 🟢 P2 (Media)  
**Prerequisito:** Sprint 2 completado

#### Objetivos Principales

1. **Crear Workflows Reusables Base**
   - pr-validation.yml (reusable)
   - release.yml (reusable)
   - sync-branches.yml (reusable)

2. **Migrar api-mobile a Reusables**
   - Convertir pr-to-dev.yml
   - Convertir pr-to-main.yml
   - Convertir sync-main-to-dev.yml
   - Mantener manual-release.yml personalizado

3. **Validar y Documentar**
   - Tests exhaustivos
   - Documentación de uso
   - Patrón para otros proyectos

#### Resultado Esperado

```yaml
✅ 3 workflows reusables creados
✅ api-mobile usa reusables
✅ Código duplicado -60%
✅ Mantenibilidad +80%
✅ Patrón documentado para replicar
```

---

## 📅 Cronograma Sugerido

### Semana 1: Sprint 2 - Migración + Optimización

#### Día 1 (Lunes): Preparación + Go 1.25
**Tiempo:** 4 horas
```
09:00-10:00  Tarea 2.1: Preparación y backup
10:00-11:30  Tarea 2.2: Migrar a Go 1.25
11:30-12:00  Tarea 2.3: Validar en CI
---
14:00-15:00  Tarea 2.4: Monitorear y validar
```

#### Día 2 (Martes): Paralelismo
**Tiempo:** 4 horas
```
09:00-10:30  Tarea 2.5: Implementar paralelismo PR→dev
10:30-12:00  Tarea 2.6: Implementar paralelismo PR→main
---
14:00-15:30  Tarea 2.7: Validar tiempos
```

#### Día 3 (Miércoles): Pre-commit + Lint
**Tiempo:** 4 horas
```
09:00-10:30  Tarea 2.8: Pre-commit hooks
10:30-11:30  Tarea 2.9: Validar hooks
---
14:00-15:00  Tarea 2.10: Corregir errores lint
15:00-16:00  Tarea 2.11: Validar lint limpio
```

#### Día 4 (Jueves): Control + Documentación
**Tiempo:** 3 horas
```
09:00-09:30  Tarea 2.12: Control releases
09:30-10:30  Tarea 2.13: Documentación
10:30-12:00  Tarea 2.14: Testing final
```

#### Día 5 (Viernes): PR y Merge
**Tiempo:** 2 horas
```
09:00-10:00  Tarea 2.15: Crear PR
10:00-11:00  Esperar review + CI
11:00-11:30  Merge a dev
```

**Total Semana 1:** 17 horas en 5 días

---

### Semana 3-4: Sprint 4 - Workflows Reusables

#### Día 1: Crear Workflows Reusables Base
**Tiempo:** 4 horas
```
Crear estructura en infrastructure
Implementar pr-validation.yml reusable
Implementar sync-branches.yml reusable
```

#### Día 2: Migrar api-mobile
**Tiempo:** 4 horas
```
Convertir pr-to-dev.yml a llamar reusable
Convertir pr-to-main.yml a llamar reusable
Convertir sync-main-to-dev.yml a llamar reusable
```

#### Día 3: Testing
**Tiempo:** 3 horas
```
Tests exhaustivos de reusables
Validar todos los casos
Ajustes finales
```

#### Día 4: Documentación
**Tiempo:** 2 horas
```
Documentar workflows reusables
Documentar cómo usarlos
Crear ejemplos
```

**Total Sprint 4:** 13 horas en 4 días

---

## 🔧 Stack Tecnológico

### Lenguaje y Framework
```
Go: 1.24.10 → 1.25 (en Sprint 2)
Framework: Gin
ORM: GORM
Swagger: gin-swagger
```

### CI/CD
```
Plataforma: GitHub Actions
Actions usadas:
  - actions/checkout@v4
  - actions/setup-go@v5
  - golangci/golangci-lint-action@v6
  - docker/setup-buildx-action@v3
  - docker/login-action@v3
  - docker/build-push-action@v5
```

### Testing
```
Framework: testing (stdlib)
Integration: testcontainers-go
Coverage: go test -coverprofile
Threshold: 33%
```

### Docker
```
Registry: ghcr.io
Platforms: linux/amd64, linux/arm64
Base image: golang:1.24-alpine → golang:1.25-alpine
```

---

## 📁 Estructura del Repositorio

```
edugo-api-mobile/
├── .github/
│   ├── workflows/
│   │   ├── pr-to-dev.yml
│   │   ├── pr-to-main.yml
│   │   ├── test.yml
│   │   ├── manual-release.yml
│   │   └── sync-main-to-dev.yml
│   └── version.txt
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── handlers/
│   ├── services/
│   ├── repositories/
│   ├── models/
│   └── middleware/
├── pkg/
│   └── ... (utilidades)
├── docs/
│   └── swagger/
├── Dockerfile
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

---

## 🎓 Aprendizajes de Proyectos Anteriores

### De shared (Sprint 1)

✅ **Funcionó bien:**
- Migración Go 1.25 validada localmente primero
- Scripts bash reutilizables
- Checkpoint después de cada tarea
- Documentación ultra-detallada

⚠️ **Mejorar:**
- Validar en CI más temprano
- Tests de integración más exhaustivos
- Rollback automático si falla

### De infrastructure

🔴 **Evitar:**
- Cambios grandes sin validación previa
- Múltiples cambios en un solo commit
- No documentar razones de cambios

✅ **Hacer:**
- Cambios pequeños e incrementales
- Validar cada cambio individualmente
- Documentar por qué, no solo qué

---

## 🚨 Riesgos y Mitigaciones

### Riesgo 1: Go 1.25 Falla en CI

**Probabilidad:** 🟡 Baja (10%)  
**Impacto:** 🔴 Alto

**Mitigación:**
- ✅ Ya validado localmente
- ✅ Script de rollback automático incluido
- ✅ Hacer en rama separada
- ✅ Validar en PR antes de merge

**Plan B:**
```bash
# Rollback automático si falla
git revert <commit-go-1.25>
git push origin feature/cicd-sprint-2
```

---

### Riesgo 2: Paralelismo Rompe Tests

**Probabilidad:** 🟡 Media (20%)  
**Impacto:** 🟡 Medio

**Mitigación:**
- ✅ Tests son independientes (testcontainers aísla)
- ✅ Validar localmente con `act`
- ✅ Hacer en commit separado (fácil revertir)

**Plan B:**
```yaml
# Revertir a secuencial si falla
needs: [lint]  # ← Readd dependencies
```

---

### Riesgo 3: Pre-commit Hooks Molestos

**Probabilidad:** 🟢 Baja (5%)  
**Impacto:** 🟢 Bajo

**Mitigación:**
- ✅ Opcional (no obligatorio)
- ✅ Documentar cómo desactivar
- ✅ Hooks rápidos (<5 seg)

**Plan B:**
```bash
# Desactivar si es muy molesto
git config core.hooksPath .git/hooks  # Volver a default
```

---

### Riesgo 4: Errores de Lint Introducen Bugs

**Probabilidad:** 🟢 Muy baja (<5%)  
**Impacto:** 🟡 Medio

**Mitigación:**
- ✅ Tests exhaustivos después de corrección
- ✅ Review cuidadoso de cada cambio
- ✅ Commit separado para lint fixes

**Plan B:**
```bash
# Revertir cambios de lint si causan problemas
git revert <commit-lint-fixes>
```

---

## ✅ Criterios de Éxito Global

### Sprint 2 Completado Cuando:

- ✅ Go 1.25 funcionando en CI sin errores
- ✅ Paralelismo implementado y probado
- ✅ Tiempos de CI reducidos al menos 25%
- ✅ Pre-commit hooks configurados y documentados
- ✅ 0 errores de lint en codebase
- ✅ Control de releases por variable funcional
- ✅ Toda la documentación actualizada
- ✅ PR mergeado a dev
- ✅ Success rate mantiene >90% o mejora

### Sprint 4 Completado Cuando:

- ✅ 3 workflows reusables creados y probados
- ✅ api-mobile usa workflows reusables
- ✅ Código duplicado reducido >60%
- ✅ Documentación completa de reusables
- ✅ Patrón listo para replicar
- ✅ Tests pasan en todos los escenarios
- ✅ PR mergeado a main

---

## 📚 Referencias

### Documentación Interna
- [Análisis Estado Actual](../../01-ANALISIS-ESTADO-ACTUAL.md)
- [Propuestas de Mejora](../../02-PROPUESTAS-MEJORA.md)
- [Matriz Comparativa](../../04-MATRIZ-COMPARATIVA.md)
- [Quick Wins](../../05-QUICK-WINS.md)
- [Resultado Pruebas Go 1.25](../../08-RESULTADO-PRUEBAS-GO-1.25.md)

### Documentación Externa
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Reusable Workflows](https://docs.github.com/en/actions/using-workflows/reusing-workflows)
- [Go 1.25 Release Notes](https://go.dev/doc/go1.25)
- [golangci-lint](https://golangci-lint.run/)

---

## 🎯 Próximos Pasos

1. **Leer:** [SPRINT-2-TASKS.md](./SPRINT-2-TASKS.md) - Plan detallado
2. **Preparar:** Entorno local, verificar acceso
3. **Ejecutar:** Tareas del Sprint 2 una por una
4. **Validar:** Cada tarea con sus checkpoints
5. **Documentar:** Ajustes y aprendizajes

---

**Generado por:** Claude Code  
**Fecha:** 19 de Noviembre, 2025  
**Versión:** 1.0  
**Estado:** Listo para Ejecución
