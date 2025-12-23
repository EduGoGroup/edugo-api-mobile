# Sprint 05: Testing Completo
# Sistema de Evaluaciones - EduGo

**Duración:** 2 días  
**Objetivo:** Suite completa de tests (unitarios, integración, E2E) con coverage >80%.

---

## 🎯 Objetivo

Asegurar calidad del código con:
- Tests unitarios dominio (>90%)
- Tests integración repositorios (>70%)
- Tests E2E flujos completos
- Tests de seguridad
- Tests de performance

---

## 📋 Tareas

Ver [TASKS.md](./TASKS.md)

---

## ✅ Validación

- [ ] Coverage global >80%
- [ ] Tests de seguridad pasando
- [ ] Tests de performance <2s p95

```bash
go test ./... -cover
go test ./tests/e2e -v -tags=e2e
```

---

**Sprint:** 05/06
