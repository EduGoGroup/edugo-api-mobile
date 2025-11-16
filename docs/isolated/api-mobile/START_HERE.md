# 🚀 START HERE - API Mobile (Sistema de Evaluaciones)

## ⭐ PUNTO DE ENTRADA ÚNICO

**Bienvenido a la documentación COMPLETA y AUTÓNOMA de edugo-api-mobile.**

Esta carpeta contiene TODO lo necesario para implementar el sistema de evaluaciones de EduGo sin depender de archivos externos.

---

## 📍 ¿Qué es edugo-api-mobile?

**API REST** que sirve a la aplicación móvil de EduGo (estudiantes, profesores, tutores).

**Puerto:** 8080  
**Tecnología:** Go + Gin + GORM + PostgreSQL + MongoDB  
**Arquitectura:** Clean Architecture (Hexagonal)

### Funcionalidades Principales
- ✅ **Autenticación JWT** (login, registro, refresh tokens)
- ✅ **Gestión de Materiales** (CRUD, upload, download)
- ✅ **Sistema de Progreso** (tracking de avance del estudiante)
- 🎯 **Sistema de Evaluaciones** (NUEVO - Lo que implementaremos)

---

## 🎯 ¿Qué Vamos a Implementar?

**Sistema Completo de Evaluaciones:**

1. **Assessments (Cuestionarios)**
   - Obtener quiz generado por IA
   - Metadata de evaluación (título, # preguntas, umbral)
   
2. **Attempts (Intentos)**
   - Crear intento de evaluación
   - Enviar respuestas
   - Calificación automática
   - Feedback detallado

3. **Resultados y Progreso**
   - Historial de intentos del estudiante
   - Estadísticas de performance
   - Tracking de aprendizaje

---

## 📂 Estructura de Esta Carpeta

```
api-mobile/
│
├── START_HERE.md                ⭐ Este archivo - LEER PRIMERO
├── EXECUTION_PLAN.md            Plan paso a paso de ejecución
│
├── 01-Context/                  Contexto del proyecto
│   ├── PROJECT_OVERVIEW.md      Overview detallado
│   ├── ECOSYSTEM_CONTEXT.md     Cómo encaja en el ecosistema
│   ├── DEPENDENCIES.md          Qué necesita de otros proyectos
│   └── TECH_STACK.md            Stack tecnológico
│
├── 02-Requirements/             Requisitos funcionales y técnicos
│   ├── PRD.md                   Product Requirements Document
│   ├── FUNCTIONAL_SPECS.md      Especificaciones funcionales
│   ├── TECHNICAL_SPECS.md       Especificaciones técnicas
│   └── ACCEPTANCE_CRITERIA.md   Criterios de aceptación
│
├── 03-Design/                   Diseño arquitectónico
│   ├── ARCHITECTURE.md          Arquitectura Clean detallada
│   ├── DATA_MODEL.md            Modelo de datos completo
│   ├── API_CONTRACTS.md         Contratos de API (OpenAPI)
│   └── SECURITY_DESIGN.md       Diseño de seguridad
│
├── 04-Implementation/           Implementación (6 sprints)
│   ├── Sprint-01-Schema-BD/     Schema PostgreSQL + MongoDB
│   ├── Sprint-02-Dominio/       Entities, Value Objects, Interfaces
│   ├── Sprint-03-Repositorios/  Implementación de repositorios
│   ├── Sprint-04-Services-API/  Services y endpoints REST
│   ├── Sprint-05-Testing/       Tests unitarios e integración
│   └── Sprint-06-CI-CD/         CI/CD y deployment
│
├── 05-Testing/                  Estrategia de testing
│   ├── TEST_STRATEGY.md
│   ├── TEST_CASES.md
│   └── COVERAGE_REPORT.md
│
├── 06-Deployment/               Deployment y monitoreo
│   ├── DEPLOYMENT_GUIDE.md
│   ├── INFRASTRUCTURE.md
│   └── MONITORING.md
│
└── PROGRESS.json                Tracking de progreso (JSON)
```

---

## 🚦 Flujo de Inicio Rápido

### Paso 1: Leer Contexto (15 min)
```bash
# Entender qué es este proyecto y cómo encaja
cat 01-Context/PROJECT_OVERVIEW.md
cat 01-Context/ECOSYSTEM_CONTEXT.md
cat 01-Context/DEPENDENCIES.md
```

### Paso 2: Revisar Requisitos (30 min)
```bash
# Entender QUÉ vamos a construir
cat 02-Requirements/PRD.md
cat 02-Requirements/FUNCTIONAL_SPECS.md
cat 02-Requirements/ACCEPTANCE_CRITERIA.md
```

### Paso 3: Estudiar Arquitectura (45 min)
```bash
# Entender CÓMO lo vamos a construir
cat 03-Design/ARCHITECTURE.md
cat 03-Design/DATA_MODEL.md
cat 03-Design/API_CONTRACTS.md
```

### Paso 4: Ejecutar Plan (Ver EXECUTION_PLAN.md)
```bash
# Plan detallado de implementación
cat EXECUTION_PLAN.md
```

### Paso 5: Implementar Sprint por Sprint (15 días estimados)
```bash
cd 04-Implementation/Sprint-01-Schema-BD/
cat README.md
cat TASKS.md
# ... ejecutar tareas ...
# Repetir para cada sprint
```

---

## 🔗 Dependencias Externas

Este proyecto **NECESITA** de otros componentes del ecosistema:

### 1. edugo-infrastructure v0.1.1 (NUEVO)
**Versión requerida:** v0.1.1  
**Qué usar:**
- `database/migrations/` - Migraciones SQL (materials, assessment, progress)
- `schemas/events/` - JSON Schemas (material.uploaded, evaluation.submitted)
- `docker/docker-compose.yml` - Servicios Docker (PostgreSQL, MongoDB, RabbitMQ)

**Estado:** ✅ COMPLETADO (96%)

### 2. edugo-shared v0.7.0 (FROZEN)
**Versión requerida:** v0.7.0 (FROZEN hasta post-MVP)  
**❌ NO USAR:** v1.3.0+ (no existen)

**Módulos usados:**
- `config` - Configuración multi-ambiente
- `database` - Conexiones PostgreSQL/MongoDB
- `auth` - JWT y autenticación
- `logger` - Logging estructurado
- `messaging/rabbit` - RabbitMQ con DLQ (NUEVO en v0.7.0)
- `evaluation` - Modelos de evaluación (NUEVO en v0.7.0)

**Estado:** ✅ COMPLETADO - 12 módulos publicados, ~75% coverage

### 3. PostgreSQL 15+
**Uso:** Base de datos principal (intentos, usuarios, materiales)  
**Tablas previas requeridas:**
- `users` (autenticación)
- `materials` (contenido educativo)

**Tablas nuevas:** `assessment`, `assessment_attempt`, `assessment_attempt_answer`

### 4. MongoDB 7.0+
**Uso:** Almacenamiento de preguntas generadas por IA  
**Colección:** `material_assessment`  
**Escritor:** edugo-worker (proceso separado)  
**Lector:** edugo-api-mobile (este proyecto)

### 5. RabbitMQ 3.12+ (Opcional para evaluaciones MVP)
**Uso:** Comunicación asíncrona (publicar eventos)  
**Eventos publicados:**
- `evaluation.submitted` → Worker procesa y genera analytics

### 6. edugo-worker (Proceso Separado)
**Responsabilidad:** Generar preguntas de evaluación con IA  
**Flujo:** Material PDF → Worker (OpenAI) → MongoDB (`material_assessment`)  
**Estado:** ✅ Debe estar funcionando para tener preguntas disponibles

---

## ⚙️ Configuración Requerida

### Variables de Entorno
```bash
# PostgreSQL
DATABASE_URL=postgres://user:pass@localhost:5432/edugo_dev?sslmode=disable

# MongoDB
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=edugo_dev

# RabbitMQ (opcional para MVP)
RABBITMQ_URL=amqp://guest:guest@localhost:5672/

# Auth
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRATION=24h

# Server
PORT=8080
ENVIRONMENT=local  # local, dev, qa, prod
LOG_LEVEL=debug
```

### Prerequisitos de Sistema
```bash
# Go 1.21+
go version

# PostgreSQL 15+
psql --version

# MongoDB 7.0+
mongosh --version

# Docker (para desarrollo)
docker --version
```

---

## 📋 Plan de Implementación

Ver archivo **EXECUTION_PLAN.md** para el plan detallado.

Resumen:
1. **Sprint 01:** Schema de base de datos (3 días)
2. **Sprint 02:** Dominio (entities, value objects) (3 días)
3. **Sprint 03:** Repositorios (3 días)
4. **Sprint 04:** Services y API REST (4 días)
5. **Sprint 05:** Testing (2 días)
6. **Sprint 06:** CI/CD (2 días)

**Total estimado:** 15-17 días laborables

---

## ✅ Checklist Pre-Implementación

Antes de comenzar Sprint 01, verifica:

### Ambiente de Desarrollo
- [ ] Go 1.21+ instalado
- [ ] PostgreSQL 15+ corriendo
- [ ] MongoDB 7.0+ corriendo
- [ ] Repositorio edugo-api-mobile clonado
- [ ] Rama feature creada: `git checkout -b feature/evaluations`

### Dependencias
- [ ] edugo-shared v1.3.0 publicado en GitHub
- [ ] Tabla `users` existe en PostgreSQL
- [ ] Tabla `materials` existe en PostgreSQL
- [ ] Al menos 1 material con `processing_status = 'completed'`

### Configuración
- [ ] Archivo `.env.local` creado con variables necesarias
- [ ] Conexión a PostgreSQL verificada: `psql -U user -d edugo_dev`
- [ ] Conexión a MongoDB verificada: `mongosh "mongodb://localhost:27017"`

### Opcional (para testing completo)
- [ ] RabbitMQ 3.12+ corriendo
- [ ] edugo-worker generó al menos 1 assessment en MongoDB

---

## 🎯 Resultado Esperado

Al completar los 6 sprints, tendrás:

### Funcionalidades
- ✅ API REST completa de evaluaciones
- ✅ 4 endpoints principales funcionando
- ✅ Calificación automática
- ✅ Feedback detallado por pregunta
- ✅ Historial de intentos

### Calidad
- ✅ Cobertura de tests >85%
- ✅ Tests de integración con Testcontainers
- ✅ Documentación Swagger actualizada
- ✅ CI/CD funcionando (GitHub Actions)

### Arquitectura
- ✅ Clean Architecture implementada
- ✅ Domain Layer independiente
- ✅ Repositorios con interfaces
- ✅ Código mantenible y testeable

---

## 📞 Soporte y Recursos

### Dentro de Esta Carpeta
- **Dudas de arquitectura:** `03-Design/ARCHITECTURE.md`
- **Dudas de requisitos:** `02-Requirements/`
- **Dudas de implementación:** `04-Implementation/Sprint-XX/TASKS.md`
- **Dudas de testing:** `05-Testing/TEST_STRATEGY.md`

### Contexto del Ecosistema
- **Cómo encaja este proyecto:** `01-Context/ECOSYSTEM_CONTEXT.md`
- **Qué depende de qué:** `01-Context/DEPENDENCIES.md`
- **Stack tecnológico:** `01-Context/TECH_STACK.md`

---

## 🚀 Comenzar AHORA

```bash
# 1. Lee el overview del proyecto
cat 01-Context/PROJECT_OVERVIEW.md

# 2. Lee el plan de ejecución
cat EXECUTION_PLAN.md

# 3. Inicia Sprint 01
cd 04-Implementation/Sprint-01-Schema-BD/
cat README.md
cat TASKS.md

# 4. Ejecuta las tareas paso a paso
# ... sigue las instrucciones de TASKS.md
```

---

**Última actualización:** 15 de Noviembre, 2025  
**Generado con:** Claude Code  
**Proyecto:** edugo-api-mobile - Sistema de Evaluaciones  
**Tipo de documentación:** Aislada y autónoma

---

## 🎓 Filosofía de Esta Documentación

> **"Todo lo que necesitas está aquí. No necesitas buscar en archivos externos. Esta carpeta es autónoma."**

**Si encuentras que falta algo, es un bug en la documentación. Repórtalo.**

---

¡Éxito en tu implementación! 🚀
