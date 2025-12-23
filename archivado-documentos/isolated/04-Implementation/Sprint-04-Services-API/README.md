# Sprint 04: Services y API REST
# Sistema de Evaluaciones - EduGo

**Duración:** 4 días  
**Objetivo:** Implementar capa de aplicación (services) y handlers REST con 4 endpoints, middleware, validación y Swagger.

---

## 🎯 Objetivo del Sprint

Crear la capa de aplicación y API REST que expone el Sistema de Evaluaciones:
- 2 services de aplicación (AssessmentService, ScoringService)
- 4 endpoints REST funcionales
- Middleware de autenticación y validación
- Documentación Swagger/OpenAPI
- Tests E2E del flujo completo

---

## 📋 Tareas del Sprint

Ver [TASKS.md](./TASKS.md)

**Tareas principales:**
- AssessmentService con lógica de orquestación
- ScoringService con validación servidor-side
- AssessmentHandler con 4 endpoints
- Middleware y rutas Gin
- Swagger annotations
- Tests E2E

---

## 🔗 Dependencias

Ver [DEPENDENCIES.md](./DEPENDENCIES.md)

**Críticas:**
- Sprint-03 completado (repositorios)
- Gin framework v1.10+
- Go validator v10
- Swag para Swagger

---

## ✅ Validación

Ver [VALIDATION.md](./VALIDATION.md)

**Criterios:**
- [ ] 4 endpoints REST funcionando
- [ ] Swagger UI accesible en /swagger/index.html
- [ ] Tests E2E pasando
- [ ] Validación servidor-side de scores
- [ ] Middleware de auth aplicado

---

## 🚀 Comandos Rápidos

```bash
# Ejecutar API
go run cmd/api/main.go

# Generar Swagger docs
swag init -g cmd/api/main.go

# Tests E2E
go test ./tests/e2e -v

# Swagger UI
open http://localhost:8080/swagger/index.html
```

---

**Sprint:** 04/06
