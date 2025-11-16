# Sprint 03: Repositorios
# Sistema de Evaluaciones - EduGo

**Duración:** 3 días  
**Objetivo:** Implementar repositorios PostgreSQL y MongoDB para persistir entities de dominio usando GORM y MongoDB driver.

---

## 🎯 Objetivo del Sprint

Crear la capa de infraestructura/persistencia que implementa las interfaces de repositorio definidas en Sprint-02. Incluye:
- 2 repositorios PostgreSQL con GORM (Assessment, Attempt)
- 1 repositorio MongoDB para preguntas
- Tests de integración con Testcontainers
- Pool de conexiones y manejo de transacciones

---

## 📋 Tareas del Sprint

Ver archivo [TASKS.md](./TASKS.md) para lista detallada.

**Resumen:**
- TASK-03-001: PostgresAssessmentRepository
- TASK-03-002: PostgresAttemptRepository (con transacciones ACID)
- TASK-03-003: MongoQuestionRepository
- TASK-03-004: Tests de integración con Testcontainers
- TASK-03-005: Connection pooling y configuración

---

## 🔗 Dependencias

Ver archivo [DEPENDENCIES.md](./DEPENDENCIES.md).

**Críticas:**
- Sprint-02 completado (interfaces de repositorios definidas)
- PostgreSQL 15+ corriendo
- MongoDB 7.0+ corriendo
- Docker para Testcontainers
- GORM v1.25.5+
- MongoDB driver v1.13.1+

---

## ❓ Decisiones y Preguntas

Ver archivo [QUESTIONS.md](./QUESTIONS.md).

**Decisiones clave:**
- Usar GORM (no SQL puro)
- Transacciones explícitas para Attempt+Answers
- Testcontainers para tests de integración (no mocks)
- Connection pooling con configuración por ambiente

---

## ✅ Validación

Ver archivo [VALIDATION.md](./VALIDATION.md) para checklist completo.

**Criterios de éxito:**
- [ ] 3 repositorios implementados
- [ ] Tests de integración con Testcontainers pasando
- [ ] Transacciones ACID funcionando
- [ ] Connection pool configurado
- [ ] Coverage >70% en repositorios

---

## 📊 Entregables

1. `internal/infrastructure/persistence/postgres_assessment_repository.go`
2. `internal/infrastructure/persistence/postgres_attempt_repository.go`
3. `internal/infrastructure/persistence/mongo_question_repository.go`
4. Tests de integración en `tests/integration/`
5. Configuración de conexiones

---

## 🚀 Comandos Rápidos

```bash
# Tests de integración (requiere Docker)
go test ./internal/infrastructure/persistence -v -tags=integration

# Tests con Testcontainers
docker ps  # Verificar contenedores de test
go test ./tests/integration/... -v

# Verificar conexión a PostgreSQL
psql -U postgres -d edugo_test -c "SELECT COUNT(*) FROM assessment;"

# Verificar conexión a MongoDB
mongosh --eval "db.material_assessment.countDocuments()"
```

---

**Generado con:** Claude Code  
**Sprint:** 03/06  
**Última actualización:** 2025-11-14
