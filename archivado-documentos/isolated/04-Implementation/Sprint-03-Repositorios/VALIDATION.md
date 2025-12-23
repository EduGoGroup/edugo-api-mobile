# Validación del Sprint 03 - Repositorios

## Pre-validación

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Verificar que PostgreSQL está corriendo
pg_isready -h localhost -p 5432

# Verificar que MongoDB está corriendo
mongosh --eval "db.runCommand({ ping: 1 })"

# Verificar que Docker está corriendo (para Testcontainers)
docker ps
```

---

## Checklist de Validación

### 1. Tests de Integración con Testcontainers
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Ejecutar tests de integración
go test ./tests/integration -v -tags=integration

# Verificar que contenedores se limpian
docker ps -a | grep testcontainers
# Esperado: vacío (contenedores eliminados después de tests)
```

**Criterio de éxito:** Tests de integración pasan, contenedores limpios

### 2. Verificar Transacciones ACID
```bash
# Test específico de transacciones
go test ./internal/infrastructure/persistence -v -run TestAttemptRepository_Transaction -tags=integration
```

**Criterio de éxito:** Attempt + Answers se guardan atómicamente

### 3. Coverage de Repositorios
```bash
go test ./internal/infrastructure/persistence -cover -tags=integration
# Esperado: >70%
```

### 4. Seguridad: Verificar que correct_answer NO se expone
```bash
go test ./internal/infrastructure/persistence -v -run TestMongoRepository_Security
```

**Criterio de éxito:** Test verifica que campo `correct_answer` no está en respuesta

### 5. Build
```bash
go build ./internal/infrastructure/...
```

**Criterio de éxito:** Build sin errores

---

## Criterios de Éxito Globales

- [ ] 3 repositorios implementados
- [ ] Tests de integración pasando
- [ ] Transacciones ACID funcionando
- [ ] Connection pool configurado
- [ ] Coverage >70%
- [ ] Build exitoso
- [ ] Seguridad verificada (no exponer respuestas)

---

## Comandos de Rollback

```bash
git checkout -- internal/infrastructure/
git checkout -- tests/integration/
```

---

**Generado con:** Claude Code  
**Sprint:** 03/06
