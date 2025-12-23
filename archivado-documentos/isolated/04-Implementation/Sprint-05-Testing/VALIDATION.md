# Validación del Sprint 05

## Validación de Coverage

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Coverage global
go test ./... -cover -coverprofile=coverage.out
go tool cover -func=coverage.out | grep total
# Esperado: >80%

# Coverage por capa
go test ./internal/domain/... -cover  # Esperado: >90%
go test ./internal/application/... -cover  # Esperado: >80%
go test ./internal/infrastructure/... -cover -tags=integration  # Esperado: >70%

# Reporte HTML
go tool cover -html=coverage.out -o coverage.html
```

---

## Criterios de Éxito

- [ ] Coverage global >80%
- [ ] Coverage dominio >90%
- [ ] Tests de seguridad pasando
- [ ] Tests de performance <2s p95
- [ ] Todos los tests pasando

---

**Sprint:** 05/06
