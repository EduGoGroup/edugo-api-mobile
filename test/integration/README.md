# Tests de Integración - API Mobile

Tests con testcontainers (PostgreSQL + MongoDB + RabbitMQ).

## Ejecutar

```bash
# Todos los integration tests
make test-integration

# Con coverage
make test-integration-coverage
```

**Nota**: Se ejecutan SOLO con `-tags=integration`
