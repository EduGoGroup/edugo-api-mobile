# Sprint 05: Testing Completo (ORIGINAL)
# Sistema de Evaluaciones - EduGo

**NOTA:** Este sprint fue dividido en Sprint 05-A y Sprint 05-B debido a su alcance extenso.

Ver:
- [Sprint-05-A-Testing-Critico/README.md](../Sprint-05-A-Testing-Critico/README.md) - Tests críticos + coverage >60%
- [Sprint-05-B-Testing-Avanzado/README.md](../Sprint-05-B-Testing-Avanzado/README.md) - Tests avanzados + coverage >80%

**Duración Original:** 2 días  
**Duración Real:** 4-5 días (por eso se dividió)

**Objetivo Original:** Suite completa de tests (unitarios, integración, E2E) con coverage >80%.

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

**Sprint:** 05/06 (dividido en 05-A y 05-B)
