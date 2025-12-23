# TASKS Sprint-00: Integrar con infrastructure

## TASK-001: Actualizar go.mod

Agregar infrastructure y actualizar shared a v0.7.0:

```bash
go get github.com/EduGoGroup/edugo-infrastructure/database@v0.2.0
go get github.com/EduGoGroup/edugo-infrastructure/schemas@v0.2.0
go get github.com/EduGoGroup/edugo-shared/auth@v0.7.0
go get github.com/EduGoGroup/edugo-shared/logger@v0.7.0
go get github.com/EduGoGroup/edugo-shared/evaluation@v0.7.0
go get github.com/EduGoGroup/edugo-shared/messaging/rabbit@v0.7.0
go get github.com/EduGoGroup/edugo-shared/database/postgres@v0.7.0
go get github.com/EduGoGroup/edugo-shared/database/mongodb@v0.7.0
go mod tidy
```

## TASK-002: Integrar validador de eventos

Usar infrastructure/schemas para validar eventos:

```go
import "github.com/EduGoGroup/edugo-infrastructure/schemas"

// Inicializar en main
validator, _ := schemas.NewValidator()

// Al publicar eventos
if err := validator.Validate(event, "material-uploaded-v1"); err != nil {
    return fmt.Errorf("invalid event: %w", err)
}
```

## TASK-003: Usar migraciones de infrastructure

Actualizar README.md:

```markdown
## Setup Base de Datos

# Usar infrastructure
cd /path/to/edugo-infrastructure
make dev-setup
```

## TASK-004: Eliminar scripts locales obsoletos

Si existe scripts/postgresql/ o scripts/mongodb/:
```bash
rm -rf scripts/postgresql/ scripts/mongodb/
```

## TASK-005: Actualizar CI/CD

Usar migraciones de infrastructure en GitHub Actions.

---

Estimación: 1-2 horas
