# Sprint 01: Schema de Base de Datos
# Sistema de Evaluaciones - EduGo

**Duración:** 2 días  
**Objetivo:** Crear schema PostgreSQL completo para evaluaciones con migraciones, índices y seeds.

---

## 🎯 Objetivo del Sprint

Implementar la capa de persistencia para el Sistema de Evaluaciones, creando 4 tablas en PostgreSQL con sus respectivas migraciones, índices optimizados y datos de prueba.

---

## 📋 Tareas del Sprint

Ver archivo [TASKS.md](./TASKS.md) para lista detallada.

**Resumen:**
- 4 tablas PostgreSQL
- Migraciones ejecutables e idempotentes
- 15+ índices optimizados
- Seeds de datos de prueba
- Tests de integridad

---

## 🔗 Dependencias

Ver archivo [DEPENDENCIES.md](./DEPENDENCIES.md).

**Críticas:**
- PostgreSQL 15+ instalado y configurado
- Tabla `materials` existente
- Tabla `users` existente
- Función `gen_uuid_v7()` disponible

---

## ❓ Decisiones y Preguntas

Ver archivo [QUESTIONS.md](./QUESTIONS.md).

**Decisiones clave:**
- Usar UUIDv7 para IDs (ordenamiento cronológico)
- Intentos inmutables (no UPDATE permitido)
- Particionamiento opcional (Post-MVP)

---

## ✅ Validación

Ver archivo [VALIDATION.md](./VALIDATION.md) para checklist completo.

**Criterios de éxito:**
- [ ] Migraciones ejecutan sin errores
- [ ] Todos los constraints funcionan
- [ ] Índices creados correctamente
- [ ] Seeds insertan datos exitosamente
- [ ] Rollback funciona

---

## 📊 Entregables

1. `scripts/postgresql/06_assessments.sql` - Migración completa
2. `scripts/postgresql/06_assessments_rollback.sql` - Rollback
3. `scripts/postgresql/seeds/assessment_seeds.sql` - Datos de prueba
4. Tests de integración pasando

---

## 🚀 Comandos Rápidos

```bash
# Ejecutar migración
psql -U postgres -d edugo < scripts/postgresql/06_assessments.sql

# Insertar seeds
psql -U postgres -d edugo < scripts/postgresql/seeds/assessment_seeds.sql

# Rollback (si es necesario)
psql -U postgres -d edugo < scripts/postgresql/06_assessments_rollback.sql

# Validar schema
psql -U postgres -d edugo -c "\d assessment"
psql -U postgres -d edugo -c "\d assessment_attempt"
```

---

**Generado con:** Claude Code  
**Sprint:** 01/06  
**Última actualización:** 2025-11-14
