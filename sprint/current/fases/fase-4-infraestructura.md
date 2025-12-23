# Fase 4: Refactorizaciones de Infraestructura

> **Estado:** ⏳ PENDIENTE
> **PR:** -
> **Branch:** `feature/infra-refactor` (por crear)
> **Estimado:** 4-6 horas

---

## Tareas

### ⏳ REF-004: Implementar Circuit Breaker para servicios externos
- **Archivo:** `internal/infrastructure/messaging/resilient_publisher.go` (nuevo)
- **Problema:** No hay protección contra fallos en cascada de RabbitMQ
- **Solución:** Usar `sony/gobreaker` (ya en go.mod)
- **Commit:** Pendiente

**Implementación requerida:**
- Crear wrapper `ResilientPublisher` con circuit breaker
- Configurar thresholds (failures, timeout, half-open)
- Integrar en bootstrap

---

### ⏳ REF-006: Implementar Healthcheck detallado
- **Archivo:** `internal/infrastructure/http/handler/health_handler.go`
- **Problema:** Healthcheck actual es básico (solo "ok")
- **Solución:** Agregar checks individuales por servicio
- **Commit:** Pendiente

**Implementación requerida:**
- Parámetro `?detail=1` para info completa
- Check de PostgreSQL (ping)
- Check de MongoDB (ping)
- Check de RabbitMQ (connection)
- Check de S3 (optional)
- Latencias de cada servicio

**Ejemplo de respuesta:**
```json
{
  "status": "healthy",
  "checks": {
    "postgres": { "status": "up", "latency_ms": 5 },
    "mongodb": { "status": "up", "latency_ms": 12 },
    "rabbitmq": { "status": "up", "latency_ms": 3 }
  }
}
```

---

### ⏳ TODO-008: Implementar lógica de deshabilitación de recursos
- **Archivo:** `internal/bootstrap/config.go:96-97`
- **Problema:** `WithDisabledResource` no implementado completamente
- **Solución:** Completar lógica para deshabilitar recursos opcionales
- **Commit:** Pendiente

**Implementación requerida:**
- Validar que el recurso existe
- Marcar recurso como deshabilitado
- Omitir inicialización en bootstrap
