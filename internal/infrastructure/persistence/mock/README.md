# Mock Repositories - Desarrollo Frontend Sin Docker

## 📋 Descripción

Este paquete proporciona implementaciones mock de todos los repositorios del sistema, permitiendo desarrollo frontend sin necesidad de levantar Docker (PostgreSQL, MongoDB, RabbitMQ).

## 🎯 Casos de Uso

- **Desarrollo frontend**: Levantar API en ~1.5s sin Docker
- **Testing rápido**: Pruebas de UI/UX con datos predecibles
- **CI/CD**: Tests más rápidos sin dependencias de infraestructura
- **Demos**: Presentaciones sin necesidad de setup complejo

## 🚀 Activación

### Opción 1: Variable de Entorno (Recomendado)

```bash
# Activar modo mock
export DEVELOPMENT_USE_MOCK_REPOSITORIES=true
go run cmd/main.go

# Desactivar (usar bases de datos reales)
export DEVELOPMENT_USE_MOCK_REPOSITORIES=false
docker-compose up -d
go run cmd/main.go
```

### Opción 2: Archivo de Configuración

Modificar `config/config.yaml`:

```yaml
development:
  use_mock_repositories: true  # Activar mocks
```

## 📦 Estructura

```
mock/
├── README.md                          # Este archivo
├── fixtures/                          # Datos predefinidos (seed data)
│   ├── users.go                       # Usuarios de prueba
│   ├── materials.go                   # Materiales de prueba
│   ├── assessments.go                 # Evaluaciones de prueba
│   ├── progress.go                    # Progreso de estudiantes
│   └── tokens.go                      # Tokens de sesión
├── postgres/                          # Mocks de repositorios PostgreSQL
│   ├── user_repository_mock.go
│   ├── material_repository_mock.go
│   ├── progress_repository_mock.go
│   ├── refresh_token_repository_mock.go
│   ├── login_attempt_repository_mock.go
│   ├── assessment_repository_mock.go
│   ├── attempt_repository_mock.go
│   └── answer_repository_mock.go
└── mongodb/                           # Mocks de repositorios MongoDB
    ├── summary_repository_mock.go
    ├── assessment_repository_mock.go
    └── assessment_document_repository_mock.go
```

## 🔧 Características de los Mocks

### Thread-Safe

Todos los mocks usan `sync.RWMutex` para operaciones concurrentes seguras:

```go
type MockUserRepository struct {
    users map[string]*pgentities.User
    mu    sync.RWMutex  // Protección de concurrencia
}
```

### Datos Predefinidos (Fixtures)

Cada mock se inicializa con datos realistas:

```go
// IDs predecibles para facilitar testing
teacherID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
studentID := uuid.MustParse("00000000-0000-0000-0000-000000000002")
```

### Inmutabilidad

Los mocks retornan copias de los datos para evitar mutaciones externas:

```go
func (r *MockUserRepository) FindByID(ctx context.Context, id valueobject.UserID) (*pgentities.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    user, exists := r.users[id.String()]
    if !exists {
        return nil, nil
    }

    return cloneUser(user), nil  // Retornar copia
}
```

## 📊 Datos de Prueba Disponibles

### Usuarios

| Email | Password | Rol | ID |
|-------|----------|-----|-----|
| `docente@edugo.com` | `password123` | teacher | `00000000-0000-0000-0000-000000000001` |
| `estudiante@edugo.com` | `password123` | student | `00000000-0000-0000-0000-000000000002` |

### Materiales

- 10 materiales de prueba con diferentes estados
- Relacionados con el docente de prueba
- Cubren todos los tipos (PDF, Video, Artículo, etc.)

### Evaluaciones

- 5 evaluaciones predefinidas
- Con preguntas de opción múltiple y verdadero/falso
- Asociadas a materiales existentes

## ⚠️ Limitaciones Importantes

### No Valida SQL

Los mocks **NO ejecutan queries SQL reales**, por lo que:

- ❌ Errores de sintaxis SQL no se detectan
- ❌ Constraints de base de datos no se validan
- ❌ Índices y performance no se prueban

### No Reemplaza Testing Real

Los mocks son para **desarrollo**, no para testing:

- ✅ Úsalos para: Desarrollo frontend rápido
- ❌ NO los uses para: Validar lógica crítica de negocio

### Sincronización Manual

Si modificas interfaces de repositorio, debes actualizar mocks manualmente.

## 🧪 Testing con Mocks

### Desarrollo Frontend

```bash
# Terminal 1: Levantar API con mocks
export DEVELOPMENT_USE_MOCK_REPOSITORIES=true
go run cmd/main.go

# Terminal 2: Ejecutar frontend
cd ../frontend
npm run dev
```

### CI/CD

Los mocks permiten tests más rápidos en pipelines:

```yaml
# .github/workflows/frontend-tests.yml
- name: Run Frontend Tests
  env:
    DEVELOPMENT_USE_MOCK_REPOSITORIES: true
  run: |
    go run cmd/main.go &
    npm run test:e2e
```

## 📝 Agregar Nuevos Mocks

### 1. Crear Fixture

```go
// fixtures/new_entity.go
package fixtures

func GetDefaultNewEntities() map[string]*Entity {
    return map[string]*Entity{
        "id-1": {
            ID:   "id-1",
            Name: "Test Entity",
        },
    }
}
```

### 2. Crear Mock Repository

```go
// postgres/new_entity_repository_mock.go
package postgres

type MockNewEntityRepository struct {
    entities map[string]*Entity
    mu       sync.RWMutex
}

func NewMockNewEntityRepository() repository.NewEntityRepository {
    return &MockNewEntityRepository{
        entities: fixtures.GetDefaultNewEntities(),
    }
}

// Implementar métodos de la interfaz...
```

### 3. Registrar en Container

```go
// internal/container/repositories.go
func newMockRepositoryContainer() *RepositoryContainer {
    return &RepositoryContainer{
        // ... otros repos
        NewEntityRepository: mockPostgres.NewMockNewEntityRepository(),
    }
}
```

## 🔍 Debugging

### Logs de Operaciones

Para ver qué operaciones ejecutan los mocks:

```yaml
# config/config.yaml
logging:
  level: "debug"  # Ver operaciones de mocks en logs
```

### Inspección de Datos

Los datos están en memoria, puedes inspeccionarlos con debugger:

```go
// En cualquier test o handler
repo := c.Repositories.UserRepository.(*mockPostgres.MockUserRepository)
fmt.Printf("Usuarios en memoria: %+v\n", repo.users)
```

## 📚 Referencias

- Diseño completo: `/docs/MOCK_REPOSITORY_ANALYSIS.md`
- Interfaces de dominio: `/internal/domain/repository/`
- Implementaciones reales: `/internal/infrastructure/persistence/postgres/`

## ⚡ Performance

| Métrica | Con Docker | Con Mocks | Mejora |
|---------|------------|-----------|---------|
| RAM | ~4 GB | ~200 MB | 95% ↓ |
| Tiempo inicio | ~30s | ~1.5s | 95% ↓ |
| Dependencias | Docker + PostgreSQL + MongoDB + RabbitMQ | Solo Go | 100% ↓ |

---

**Versión**: 1.0  
**Última actualización**: 2025-11-29
