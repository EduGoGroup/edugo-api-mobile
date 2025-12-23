# Reporte de Coverage
# Sistema de Evaluaciones - EduGo

**Versión:** 1.0.0  
**Fecha:** 14 de Noviembre, 2025

---

## 1. TEMPLATE DE REPORTE

### Coverage por Package

| Package | Statements | Coverage | Status |
|---------|-----------|----------|--------|
| internal/domain/entities | XXX | XX.X% | ✅/❌ |
| internal/domain/valueobjects | XXX | XX.X% | ✅/❌ |
| internal/domain/repositories | XXX | XX.X% | ✅/❌ |
| internal/application/services | XXX | XX.X% | ✅/❌ |
| internal/infrastructure/persistence | XXX | XX.X% | ✅/❌ |
| internal/interfaces/http/handlers | XXX | XX.X% | ✅/❌ |
| **TOTAL** | **XXX** | **XX.X%** | **✅/❌** |

**Objetivo:** >80% global

### Coverage por Capa

| Capa | Coverage | Objetivo | Delta |
|------|----------|----------|-------|
| Domain | XX.X% | >90% | +X.X% |
| Application | XX.X% | >85% | +X.X% |
| Infrastructure | XX.X% | >70% | +X.X% |
| Handlers | XX.X% | >80% | +X.X% |

---

## 2. COMANDOS PARA GENERAR REPORTE

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Generar coverage
go test ./... -coverprofile=coverage.out

# Ver por función
go tool cover -func=coverage.out

# Generar HTML
go tool cover -html=coverage.out -o coverage.html

# Coverage por package (tabla)
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out | grep -E "^(.*/).*\.go"

# Solo totales
go tool cover -func=coverage.out | grep total
```

---

## 3. GAPS DE COVERAGE

### Áreas Sin Tests (Ejemplo)

**Gap 1:** Error handling en handlers  
**Prioridad:** Alta  
**Plan:** Agregar tests de error responses

**Gap 2:** Edge cases en validators  
**Prioridad:** Media  
**Plan:** Tests con valores límite

---

## 4. PLAN DE MEJORA

Si coverage <80%:

1. Identificar packages con coverage bajo
2. Priorizar dominio (crítico)
3. Agregar tests de casos faltantes
4. Re-ejecutar coverage

```bash
# Identificar gaps
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out | awk '$3 < 80 {print $1, $3}'
```

---

**Generado con:** Claude Code
