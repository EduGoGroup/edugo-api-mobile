# Plan de Trabajo: Unificación de api-mobile con estándar api-admin

## Objetivo
Unificar api-mobile con el estándar de api-admin para:
1. Eliminar entrypoint innecesario (responsabilidad de docker-compose)
2. Implementar estructura de configuración Viper igual a api-admin

---

## Fase 1: Eliminar Entrypoint

### Justificación
- El entrypoint wait-for es **redundante** porque docker-compose ya maneja esto con `depends_on: condition: service_healthy`
- api-admin no tiene entrypoint y funciona correctamente
- El entrypoint causa problemas con modo mock

### Pasos

#### 1.1 Modificar Dockerfile
- Eliminar COPY de scripts de entrypoint
- Eliminar ENTRYPOINT, dejar solo CMD

**Antes:**
```dockerfile
COPY scripts/docker-entrypoint.sh /scripts/
COPY scripts/wait-for.sh /scripts/
ENTRYPOINT ["/scripts/docker-entrypoint.sh"]
CMD ["./main"]
```

**Después:**
```dockerfile
CMD ["./main"]
```

#### 1.2 Eliminar archivos de scripts
- `/scripts/docker-entrypoint.sh`
- `/scripts/wait-for.sh`

#### 1.3 Actualizar documentación si aplica

---

## Fase 2: Agregar Estructura Config Viper

### Justificación
- api-admin usa archivos config/*.yaml + override con ENV vars (patrón Viper estándar)
- api-mobile solo usa ENV vars, inconsistente con api-admin
- Viper permite configuración por ambiente (local, dev, qa, prod)

### Pasos

#### 2.1 Crear estructura de carpeta config/
```
config/
├── README.md
├── config.yaml          # Configuración base
├── config-local.yaml    # Desarrollo local
├── config-dev.yaml      # Ambiente dev
├── config-qa.yaml       # Ambiente QA
├── config-prod.yaml     # Producción
└── config-test.yaml     # Testing
```

#### 2.2 Crear config.yaml base
Basado en la estructura actual de ENV vars:
- server (port, host, timeouts)
- database (postgres, mongodb)
- messaging (rabbitmq)
- storage (s3)
- auth (jwt)
- logging
- development (use_mock_repositories)

#### 2.3 Modificar loader de configuración
- Actualizar `internal/config/loader.go` para cargar archivos YAML
- Mantener compatibilidad con ENV vars (override)

#### 2.4 Actualizar Dockerfile
```dockerfile
COPY config/ /root/config/
```

---

## Fase 3: Testing y Validación

### 3.1 Tests locales
```bash
go test ./... -short
```

### 3.2 Validar docker-compose principal
```bash
cd edugo-dev-environment/docker
docker-compose up -d
```

### 3.3 Validar docker-compose-mock
- Verificar que ya no necesita entrypoint override

---

## Fase 4: PR y Release

### 4.1 PR a dev
- Branch: `feature/unify-config-standard`
- Monitorear pipeline (max 10 min)
- Merge a dev

### 4.2 PR a main
- Crear PR dev → main
- Monitorear pipeline
- Merge a main

### 4.3 Manual Release
- Version: 0.14.0 (minor - refactor)

---

## Archivos a Modificar

| Archivo | Acción |
|---------|--------|
| Dockerfile | Eliminar entrypoint, agregar COPY config/ |
| scripts/docker-entrypoint.sh | Eliminar |
| scripts/wait-for.sh | Eliminar |
| config/*.yaml | Crear (nuevos) |
| internal/config/loader.go | Modificar para cargar YAML |

---

## Estimación: ~1.5 horas
