# Infrastructure Stubs (TEMPORAL - FASE 1)

Este directorio contiene **stubs temporales** que simulan las entities de `github.com/EduGoGroup/edugo-infrastructure/postgres/entities`.

## ⚠️ IMPORTANTE: SOLO PARA FASE 1

Estos stubs fueron creados porque en **Fase 1** no hay:
- Conexión a internet para descargar dependencias
- Acceso al módulo real de infrastructure
- Ambiente completo para validación

## 📋 Entities Stubbed

### Postgres Entities
1. `Material` - Material educativo (PDF)
2. `User` - Usuario del sistema
3. `MaterialVersion` - Versión de material
4. `Progress` - Progreso de usuario en material
5. `Assessment` - Evaluación de material
6. `AssessmentAnswer` - Respuesta de evaluación
7. `AssessmentAttempt` - Intento de evaluación

## 🔄 FASE 2: Reemplazo

En **Fase 2** (con ambiente completo), estos stubs deben ser:

1. **Eliminados completamente**:
   ```bash
   rm -rf internal/infrastructure_stubs/
   ```

2. **Reemplazados por imports reales**:
   ```go
   // Reemplazar:
   import pgentities "github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure_stubs/postgres/entities"

   // Por:
   import pgentities "github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
   ```

3. **Actualizar go.mod**:
   ```bash
   go get github.com/EduGoGroup/edugo-infrastructure/postgres/entities@latest
   go mod tidy
   ```

4. **Validar compilación**:
   ```bash
   go build ./...
   go test ./...
   ```

## 📝 Diferencias con Infrastructure Real

Los stubs pueden tener pequeñas diferencias con las entities reales de infrastructure:

- Tags GORM pueden variar
- Campos adicionales pueden faltar
- Métodos helper de infrastructure no están incluidos
- Validaciones de infrastructure no están presentes

**Esto es aceptable para Fase 1**, ya que el objetivo es avanzar con la estructura y lógica del sprint.

## ✅ Uso en Fase 1

Para usar estos stubs en fase 1:

```go
import pgentities "github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure_stubs/postgres/entities"

// Usar normalmente:
material := &pgentities.Material{
    ID:    uuid.New(),
    Title: "Material de prueba",
    // ...
}
```

## 🚨 NO COMMITEAR EN FASE 2

Este directorio es **temporal** y **NO debe ser commiteado** en la versión final del sprint. Solo existe para facilitar el trabajo en Fase 1.

---

**Creado**: Sprint ENTITIES - Fase 1
**Estado**: Temporal - Eliminar en Fase 2
**Autor**: Claude Code (Sprint Entities Adaptation)
