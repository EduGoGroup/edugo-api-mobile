# Sprint 06: CI/CD y Documentación
# Sistema de Evaluaciones - EduGo

**Duración:** 2 días  
**Objetivo:** Pipeline CI/CD completo con GitHub Actions, linting automático, tests en CI, build Docker y documentación final.

---

## 🎯 Objetivo

Automatizar calidad y deployment:
- GitHub Actions workflow completo
- Linting automático en CI
- Tests automáticos (unit, integration, E2E)
- Build y publish de imagen Docker
- Documentación README actualizada

---

## 📋 Tareas

Ver [TASKS.md](./TASKS.md)

---

## ✅ Validación

- [ ] Pipeline verde en GitHub Actions
- [ ] Linting automático
- [ ] Tests automáticos
- [ ] Docker image publicada

```bash
# Ver actions
gh run list --workflow=ci.yml

# Ejecutar localmente
docker build -t edugo-api-mobile .
```

---

**Sprint:** 06/06
