# Fase 6: Mejoras de Observabilidad

> **Estado:** ⏳ PENDIENTE
> **PR:** -
> **Branch:** `feature/observability` (por crear)
> **Estimado:** 4-6 horas

---

## Tareas

### ⏳ REF-005: Agregar Request ID y Tracing
- **Archivo:** `internal/infrastructure/http/middleware/request_id.go` (nuevo)
- **Problema:** No hay correlación entre logs de una misma request
- **Solución:** Middleware que genera/propaga request_id
- **Commit:** Pendiente

**Implementación requerida:**
- Middleware `RequestID()` que:
  - Lee `X-Request-ID` del header si existe
  - Genera UUID si no existe
  - Agrega al contexto de Gin
  - Propaga en header de respuesta
- Propagar request_id en logs
- Propagar en headers de RabbitMQ

---

### ⏳ Mejorar logging estructurado
- **Problema:** Logs no tienen contexto suficiente para debugging
- **Solución:** Agregar información contextual a todos los logs
- **Commit:** Pendiente

**Información a agregar:**
- `request_id` en todos los logs de handlers
- `endpoint` y `method`
- `duration` al finalizar request
- `user_id` cuando disponible

---

### ⏳ Agregar métricas básicas
- **Problema:** No hay visibilidad de performance en producción
- **Solución:** Implementar métricas con Prometheus
- **Commit:** Pendiente

**Métricas a implementar:**
- `http_requests_total` - Contador por endpoint/method/status
- `http_request_duration_seconds` - Histograma de latencias
- `http_requests_in_flight` - Gauge de requests activos
- `errors_total` - Contador de errores por tipo

**Endpoint:** `GET /metrics` (Prometheus format)
