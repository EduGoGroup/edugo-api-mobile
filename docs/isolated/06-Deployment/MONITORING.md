# Monitoreo y Observabilidad
# Sistema de Evaluaciones - EduGo

**Versión:** 1.0.0  
**Fecha:** 14 de Noviembre, 2025

---

## 1. MÉTRICAS CLAVE

### Application Metrics

| Métrica | Objetivo | Alerta Si |
|---------|----------|-----------|
| **Latencia p50** | <200ms | >500ms |
| **Latencia p95** | <2000ms | >5000ms |
| **Latencia p99** | <5000ms | >10000ms |
| **Throughput** | Variable | - |
| **Error Rate** | <1% | >5% |
| **Success Rate** | >99% | <95% |

### Database Metrics

| Métrica | Objetivo | Alerta Si |
|---------|----------|-----------|
| **DB Connections** | <20 (de 25 max) | >23 |
| **Query Time p95** | <100ms | >500ms |
| **Slow Queries** | <10/min | >50/min |

### Business Metrics

| Métrica | Descripción | Tracking |
|---------|-------------|----------|
| **assessment_attempts_total** | Total intentos creados | Counter |
| **assessment_attempts_passed** | Intentos aprobados | Counter |
| **assessment_score_avg** | Score promedio | Gauge |
| **assessment_duration_seconds** | Duración promedio | Histogram |

---

## 2. PROMETHEUS METRICS

### Exponer Métricas

Archivo: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/internal/infrastructure/metrics/prometheus.go`

```go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    AttemptsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "assessment_attempts_total",
            Help: "Total number of assessment attempts",
        },
        []string{"assessment_id", "passed"},
    )
    
    AttemptDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "assessment_duration_seconds",
            Help: "Time spent on assessment",
            Buckets: []float64{60, 300, 600, 1200, 1800, 3600},
        },
        []string{"assessment_id"},
    )
    
    HTTPRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "HTTP request latency",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint", "status"},
    )
)
```

**Endpoint:**
```go
import "github.com/prometheus/client_golang/prometheus/promhttp"

router.GET("/metrics", gin.WrapH(promhttp.Handler()))
```

---

## 3. ALERTAS CRÍTICAS

### Configuración Prometheus

```yaml
# prometheus/alerts.yml
groups:
  - name: edugo_api_mobile
    rules:
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
        for: 5m
        annotations:
          summary: "Error rate >5% en API Mobile"
      
      - alert: HighLatency
        expr: histogram_quantile(0.95, http_request_duration_seconds) > 2
        for: 5m
        annotations:
          summary: "Latencia p95 >2s"
      
      - alert: DatabaseConnectionPoolExhausted
        expr: db_connections_active / db_connections_max > 0.9
        for: 2m
        annotations:
          summary: "Pool de conexiones DB >90%"
```

---

## 4. LOGS ESTRUCTURADOS

### Usando edugo-shared Logger

```go
import "edugo-shared/logger"

// Configurar logger
logger.Init(logger.Config{
    Level:  "info",      // debug|info|warn|error
    Format: "json",      // json|text
    Output: "stdout",
})

// Usar en código
logger.Info("Assessment retrieved",
    "material_id", materialID,
    "student_id", studentID,
    "total_questions", assessment.TotalQuestions,
)

logger.Error("Failed to save attempt",
    "error", err,
    "student_id", studentID,
    "assessment_id", assessmentID,
)
```

**Output JSON:**
```json
{
  "level": "info",
  "timestamp": "2025-11-14T12:00:00Z",
  "message": "Assessment retrieved",
  "material_id": "01936d9a...",
  "student_id": "01936d9b...",
  "total_questions": 5
}
```

---

## 5. DASHBOARDS

### Grafana Dashboard Básico

**Panels:**
1. Request Rate (req/s)
2. Latency (p50, p95, p99)
3. Error Rate (%)
4. Assessment Attempts (total, passed, failed)
5. Database Connections
6. Memory/CPU Usage

**Importar dashboard:**
```bash
# Dashboard ID en Grafana.com: 15398 (ejemplo: Go Application Dashboard)
# O crear manual con queries:
rate(http_requests_total[5m])
histogram_quantile(0.95, http_request_duration_seconds)
```

---

**Generado con:** Claude Code
