# Product Requirements Document (PRD)
# Sistema de Evaluaciones - EduGo

**Versión:** 1.0.0  
**Fecha:** 14 de Noviembre, 2025  
**Proyecto:** EduGo - Sistema de Evaluaciones Automáticas  
**Prioridad:** 🔴 P0 (CRÍTICA)

---

## 1. VISIÓN DEL PRODUCTO

### 1.1 Propósito

Implementar un sistema completo de evaluaciones automáticas que permita a los estudiantes validar su comprensión de materiales educativos mediante cuestionarios generados por IA, recibiendo retroalimentación inmediata y personalizada.

### 1.2 Problema a Resolver

**Situación Actual:**
- El sistema EduGo puede procesar PDFs y generar resúmenes mediante IA (worker)
- Los quizzes generados por el worker se almacenan en MongoDB
- **NO EXISTE** infraestructura para que estudiantes:
  - Accedan a los cuestionarios
  - Respondan las preguntas
  - Reciban calificaciones automáticas
  - Obtengan retroalimentación educativa
  - Rastreen su historial de intentos

**Impacto del Problema:**
- Funcionalidad core del producto incompleta (40% implementado según GAP_ANALYSIS.md)
- Imposibilidad de validar aprendizaje estudiantil
- Sin métricas de rendimiento académico
- Docentes no pueden evaluar progreso real

### 1.3 Solución Propuesta

Completar el sistema de evaluaciones implementando:

1. **Base de Datos (PostgreSQL):**
   - Tabla `assessment` - Metadatos de evaluaciones
   - Tabla `assessment_attempt` - Intentos de estudiantes
   - Tabla `assessment_attempt_answer` - Respuestas individuales
   - Tabla `material_summary_link` - Enlaces a resúmenes MongoDB

2. **API REST (edugo-api-mobile):**
   - Endpoints para obtener cuestionarios
   - Endpoints para enviar respuestas
   - Endpoints para consultar resultados
   - Lógica de calificación automática

3. **Integración MongoDB + PostgreSQL:**
   - Lectura de preguntas desde MongoDB
   - Almacenamiento de intentos en PostgreSQL
   - Validación servidor-side (seguridad)

---

## 2. OBJETIVOS DE NEGOCIO

### 2.1 Objetivos Primarios

| # | Objetivo | Métrica de Éxito | Plazo |
|---|----------|------------------|-------|
| OB-01 | Habilitar evaluación automática de comprensión estudiantil | >60% de estudiantes completan quiz tras leer material | 2 semanas |
| OB-02 | Proporcionar retroalimentación inmediata y educativa | Tiempo de respuesta <2 seg, feedback personalizado | 2 semanas |
| OB-03 | Registrar historial completo de intentos para analytics | 100% de intentos registrados permanentemente | 2 semanas |
| OB-04 | Validar aprendizaje real (no solo lectura) | Puntaje promedio >70% | 1 mes post-lanzamiento |

### 2.2 Objetivos Secundarios

| # | Objetivo | Métrica de Éxito | Plazo |
|---|----------|------------------|-------|
| OB-05 | Notificar docentes de resultados estudiantiles | 100% de docentes reciben notificación asíncrona | Post-MVP |
| OB-06 | Permitir múltiples intentos de mejora | <10% abandono tras primer intento | Post-MVP |
| OB-07 | Generar reportes de rendimiento por unidad académica | Dashboard operativo | Post-MVP |

---

## 3. STAKEHOLDERS

### 3.1 Usuarios Directos

| Rol | Necesidad Principal | Pain Point Actual |
|-----|---------------------|-------------------|
| **Estudiante** | Validar comprensión del material, recibir feedback inmediato | No puede autoevaluarse, no sabe si comprendió |
| **Profesor** | Ver rendimiento de estudiantes, identificar temas difíciles | No tiene visibilidad del aprendizaje real |
| **Tutor/Padre** | Monitorear progreso académico del estudiante | Solo ve si leyó, no si comprendió |

### 3.2 Usuarios Indirectos

| Rol | Interés | Beneficio |
|-----|---------|-----------|
| **Administrador Escuela** | Métricas de calidad educativa | Reportes de rendimiento por sección |
| **Desarrolladores** | Sistema mantenible y escalable | Arquitectura clean, tests >80% |

---

## 4. ALCANCE DEL PROYECTO

### 4.1 EN ALCANCE (MVP - 2 semanas)

#### Funcionalidades Core

✅ **FC-01: Obtener Cuestionario**
- Endpoint `GET /v1/materials/:id/assessment`
- Retorna preguntas SIN respuestas correctas (seguridad)
- Integración con MongoDB (colección `material_assessment`)
- Validación de permisos (solo usuarios con acceso al material)

✅ **FC-02: Enviar Respuestas y Obtener Calificación**
- Endpoint `POST /v1/materials/:id/assessment/attempts`
- Validación de respuestas en servidor
- Cálculo automático de puntaje
- Generación de feedback educativo por pregunta
- Persistencia de intento y respuestas en PostgreSQL

✅ **FC-03: Consultar Resultados**
- Endpoint `GET /v1/attempts/:id/results`
- Endpoint `GET /v1/users/me/attempts` (historial)
- Detalle completo: puntaje, respuestas correctas/incorrectas, feedback

✅ **FC-04: Schema de Base de Datos**
- Crear 4 tablas PostgreSQL
- Migraciones ejecutables
- Seeds de datos de prueba
- Índices optimizados

✅ **FC-05: Tests Completos**
- Tests unitarios de dominio (>85% coverage)
- Tests de integración con testcontainers
- Tests end-to-end de flujo completo

### 4.2 FUERA DE ALCANCE (Post-MVP)

❌ **Tipo de preguntas avanzadas:**
- Verdadero/Falso
- Selección múltiple (varias respuestas correctas)
- Respuesta corta con NLP
- Emparejamiento

❌ **Límite de reintentos:**
- Restricción de máximo N intentos por día
- Cooldown entre intentos

❌ **Banco de preguntas aleatorias:**
- Selección aleatoria de 5 de 20 preguntas
- Prevención de memorización

❌ **Retroalimentación adaptativa:**
- Feedback diferente según magnitud del error
- Sugerencias de secciones específicas para repasar

❌ **Notificación a docentes:**
- Worker consume evento `assessment_attempt_recorded`
- Email/push notification a docentes

❌ **Reportes y analytics:**
- Dashboard de rendimiento por unidad académica
- Identificación de preguntas problemáticas
- Curvas de aprendizaje

---

## 5. RESTRICCIONES Y SUPUESTOS

### 5.1 Restricciones Técnicas

| ID | Restricción | Impacto |
|----|-------------|---------|
| RT-01 | Usar arquitectura existente de api-mobile (Clean Architecture) | Mantener consistencia con codebase |
| RT-02 | Go 1.21+ con Gin framework | Tecnología ya definida |
| RT-03 | PostgreSQL para datos relacionales, MongoDB para documentos | Arquitectura híbrida existente |
| RT-04 | Reutilizar shared/testing para testcontainers | Evitar duplicación |
| RT-05 | Coverage mínimo 80% en CI/CD | Estándar del proyecto |

### 5.2 Restricciones de Negocio

| ID | Restricción | Razón |
|----|-------------|-------|
| RN-01 | Tiempo de respuesta <2 seg para calificación | Experiencia de usuario |
| RN-02 | NUNCA enviar respuestas correctas al cliente antes de enviar respuestas | Seguridad anti-trampa |
| RN-03 | Intentos inmutables (no editables) | Auditoría y trazabilidad |
| RN-04 | Validación de respuestas siempre en servidor | Seguridad |

### 5.3 Supuestos

| ID | Supuesto | Validación Necesaria |
|----|----------|----------------------|
| AS-01 | El worker ya genera quizzes y los guarda en MongoDB | ⚠️ **NO VERIFICADO** - Ver VERIFICACION_WORKER.md |
| AS-02 | La colección `material_assessment` existe en MongoDB | ⚠️ **NO VERIFICADO** |
| AS-03 | Las preguntas tienen estructura estándar (id, text, options, correct_answer) | Requiere inspección |
| AS-04 | Cada material tiene máximo 1 assessment (relación 1:1) | Simplificación MVP |

---

## 6. CRITERIOS DE ÉXITO

### 6.1 Criterios Funcionales

| ID | Criterio | Método de Validación |
|----|----------|----------------------|
| CF-01 | Estudiante puede obtener quiz de un material | Test E2E con material de prueba |
| CF-02 | Estudiante recibe calificación inmediata (<2 seg) | Test de performance |
| CF-03 | Feedback educativo presente en todas las respuestas incorrectas | Inspección manual de response |
| CF-04 | Historial de intentos accesible | Test E2E de endpoint /users/me/attempts |
| CF-05 | Respuestas correctas NUNCA expuestas antes de enviar respuestas | Test de seguridad |

### 6.2 Criterios Técnicos

| ID | Criterio | Método de Validación |
|----|----------|----------------------|
| CT-01 | Schema BD ejecutable sin errores | Migración en entorno limpio |
| CT-02 | Coverage de tests >80% | Reporte go test -cover |
| CT-03 | Tests de integración pasando con testcontainers | CI/CD pipeline |
| CT-04 | API documentada con Swagger | Endpoint /swagger/index.html |
| CT-05 | Sin errores de linting | golangci-lint run |

### 6.3 Criterios de Calidad

| ID | Criterio | Método de Validación |
|----|----------|----------------------|
| CQ-01 | Código sigue Clean Architecture (domain, application, infrastructure) | Code review |
| CQ-02 | Transacciones ACID para intentos (attempt + answers atómico) | Test de rollback |
| CQ-03 | Manejo de errores robusto | Test de failure scenarios |
| CQ-04 | Logging estructurado de todas las operaciones | Inspección de logs |

---

## 7. MÉTRICAS DE ÉXITO (KPIs)

### 7.1 Métricas de Producto

| KPI | Objetivo | Medición |
|-----|----------|----------|
| **Tasa de Completitud de Quiz** | >60% | COUNT(assessment_attempt) / COUNT(DISTINCT reading_log WHERE progress=100) |
| **Puntaje Promedio** | >70% | AVG(score) de assessment_attempt |
| **Tiempo Promedio de Completitud** | 8-12 min para 5 preguntas | AVG(time_spent_seconds) |
| **Tasa de Abandono** | <10% | (quiz_started - quiz_completed) / quiz_started |

### 7.2 Métricas Técnicas

| KPI | Objetivo | Medición |
|-----|----------|----------|
| **Latencia p95** | <500ms | Prometheus metrics |
| **Tasa de Error** | <1% | HTTP 5xx / total requests |
| **Test Coverage** | >80% | go test -cover |
| **Uptime** | >99.9% | Monitoreo de health endpoint |

---

## 8. RIESGOS Y MITIGACIONES

### 8.1 Riesgos Técnicos

| ID | Riesgo | Probabilidad | Impacto | Mitigación |
|----|--------|--------------|---------|------------|
| RT-01 | Worker NO genera quizzes en MongoDB | Alta | 🔴 Crítico | **Acción inmediata:** Verificar código del worker. Plan B: Generar manualmente para MVP |
| RT-02 | Esquema de preguntas en MongoDB incompatible | Media | 🟡 Alto | Inspeccionar colección, crear adapter si necesario |
| RT-03 | Performance de queries de validación <2 seg | Baja | 🟡 Alto | Índices en PostgreSQL, caché de preguntas |
| RT-04 | Integración con MongoDB falla en producción | Baja | 🔴 Crítico | Tests exhaustivos con testcontainers, retry logic |

### 8.2 Riesgos de Negocio

| ID | Riesgo | Probabilidad | Impacto | Mitigación |
|----|--------|--------------|---------|------------|
| RN-01 | Estudiantes no usan cuestionarios | Media | 🟡 Alto | UX simple, gamificación (puntajes visibles) |
| RN-02 | Feedback no suficientemente educativo | Media | 🟡 Alto | Iterar con profesores, mejorar prompts de IA |
| RN-03 | Trampas (compartir respuestas) | Baja | 🟢 Bajo | Post-MVP: banco aleatorio, límite de intentos |

---

## 9. DEPENDENCIAS

### 9.1 Dependencias Externas

| Dependencia | Proveedor | Criticidad | Estado |
|-------------|-----------|------------|--------|
| Worker funcionando | edugo-worker | 🔴 Bloqueante | ⚠️ No verificado |
| MongoDB con colección `material_assessment` | edugo-worker | 🔴 Bloqueante | ⚠️ No verificado |
| PostgreSQL disponible | Infraestructura | 🔴 Bloqueante | ✅ Funcional |
| shared/testing | edugo-shared | 🟡 Alta | ✅ Funcional (v0.6.2) |

### 9.2 Dependencias Internas

| Dependencia | Proyecto | Criticidad | Estado |
|-------------|----------|------------|--------|
| Tabla `materials` existente | api-mobile | 🔴 Bloqueante | ✅ Implementada |
| Tabla `users` existente | api-mobile | 🔴 Bloqueante | ✅ Implementada |
| Auth JWT funcionando | api-mobile | 🔴 Bloqueante | ✅ Funcional |
| Middleware de autenticación | shared/middleware | 🔴 Bloqueante | ✅ Funcional |

---

## 10. CRONOGRAMA DE ALTO NIVEL

### Sprint Mobile-1: Sistema de Evaluaciones (2 semanas)

| Semana | Entregables | Responsable |
|--------|-------------|-------------|
| **Semana 1** | - Schema BD (4 tablas) <br> - Dominio (entities, value objects, repositories) <br> - Infraestructura (repositorios PostgreSQL + MongoDB) | Desarrollador Go Senior |
| **Semana 2** | - Services de aplicación <br> - API REST (4 endpoints) <br> - Tests completos (unitarios + integración + E2E) <br> - Documentación Swagger | Desarrollador Go Senior |

---

## 11. APROBACIONES

### 11.1 Stakeholders Clave

| Rol | Nombre | Firma | Fecha |
|-----|--------|-------|-------|
| Product Owner | [PENDIENTE] | | |
| Tech Lead | [PENDIENTE] | | |
| QA Lead | [PENDIENTE] | | |

### 11.2 Criterios de Aprobación de PRD

- ✅ Alcance claramente definido (MVP vs Post-MVP)
- ✅ Riesgos identificados con mitigaciones
- ✅ Dependencias documentadas
- ✅ KPIs medibles definidos
- ✅ Cronograma realista

---

## 12. HISTORIAL DE CAMBIOS

| Versión | Fecha | Autor | Cambios |
|---------|-------|-------|---------|
| 1.0.0 | 2025-11-14 | Claude Code | Versión inicial basada en documentación existente |

---

## 13. REFERENCIAS

- [Plan de Implementación](../../../docs/roadmap/PLAN_IMPLEMENTACION.md)
- [GAP Analysis](../../../docs/analisis/GAP_ANALYSIS.md)
- [Historia de Usuario HU-MOB-EVA-01](../../../docs/historias_usuario/api_mobile/evaluacion/HU_MOB_EVA_01_realizar_quiz.md)
- [Flujo de Evaluación](../../../docs/diagramas/procesos/03_evaluacion.md)
- [Distribución de Responsabilidades](../../../docs/analisis/DISTRIBUCION_RESPONSABILIDADES.md)

---

**Generado con:** Claude Code  
**Proyecto:** EduGo - Sistema de Evaluaciones  
**Prioridad:** 🔴 P0 (CRÍTICA)
