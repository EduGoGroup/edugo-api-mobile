# GitHub Copilot - Instrucciones Personalizadas: EduGo API Mobile

## 🌍 IDIOMA / LANGUAGE

**IMPORTANTE**: Todos los comentarios, sugerencias, code reviews y respuestas en chat deben estar **SIEMPRE EN ESPAÑOL**.

- ✅ Comentarios en Pull Requests: **español**
- ✅ Sugerencias de código: **español**
- ✅ Explicaciones en chat: **español**
- ✅ Mensajes de error: **español**

---

## 🏗️ Arquitectura del Proyecto

Este proyecto implementa **Clean Architecture (Hexagonal)** con Go 1.25:

```
internal/
├── domain/              # Entidades, Value Objects, Interfaces
├── application/         # Servicios, DTOs, Casos de uso
├── infrastructure/      # Implementaciones concretas
│   ├── http/           # Handlers, Middleware
│   ├── persistence/    # Repositorios (PostgreSQL, MongoDB)
│   └── messaging/      # RabbitMQ (pendiente implementar)
├── container/          # Inyección de Dependencias
└── config/             # Configuración con Viper
```

### Principios Arquitectónicos
- **Dependency Inversion**: El dominio NO depende de infraestructura
- **Separation of Concerns**: Cada capa tiene responsabilidades claras
- **Dependency Injection**: Usar container/container.go para DI
- **Interface Segregation**: Interfaces pequeñas y específicas

---

## 📦 Dependencia Compartida: edugo-shared

Usamos el módulo `github.com/EduGoGroup/edugo-shared` para funcionalidad compartida:

### Paquetes Disponibles
- **logger**: Logger Zap estructurado (`edugo-shared/logger`)
- **auth**: JWT Manager y autenticación (`edugo-shared/auth`)
- **errors**: Tipos de error de aplicación (`edugo-shared/common/errors`)
- **middleware**: Middleware reutilizable (en desarrollo)

### ⚠️ REGLA CRÍTICA: NO Reimplementar Funcionalidad

```go
// ❌ INCORRECTO: Reimplementar funcionalidad existente
type MyJWTManager struct { ... }
func (m *MyJWTManager) GenerateToken() { ... }

// ✅ CORRECTO: Usar edugo-shared
import "github.com/EduGoGroup/edugo-shared/auth"
jwtManager := auth.NewJWTManager(secret, expiration)
```

---

## 🎯 Convenciones de Código

### Naming Conventions

```go
// DTOs
type UserDTO struct { ... }          // ✅ Termina en DTO
type CreateCourseDTO struct { ... }  // ✅ Termina en DTO

// Servicios
type UserService struct { ... }      // ✅ Termina en Service
type AuthService struct { ... }      // ✅ Termina en Service

// Repositorios
type UserRepository interface { ... } // ✅ Termina en Repository
type PostgresUserRepository struct { ... } // ✅ Implementación específica

// Handlers
type UserHandler struct { ... }      // ✅ Termina en Handler
```

### Manejo de Errores

```go
// ✅ CORRECTO: Usar tipos de error de edugo-shared
import "github.com/EduGoGroup/edugo-shared/common/errors"

func (s *UserService) GetUser(ctx context.Context, id string) (*UserDTO, error) {
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        if errors.IsNotFound(err) {
            return nil, errors.NewNotFoundError("user", id)
        }
        return nil, errors.NewInternalError("failed to get user", err)
    }
    return user, nil
}

// ❌ INCORRECTO: NO usar fmt.Errorf directamente
return nil, fmt.Errorf("user not found: %s", id)

// ❌ INCORRECTO: NO usar errors.New
return nil, errors.New("user not found")
```

### Context en Todas las Funciones

```go
// ✅ CORRECTO: Siempre recibir context.Context como primer parámetro
func (s *UserService) CreateUser(ctx context.Context, dto CreateUserDTO) (*UserDTO, error)
func (r *PostgresUserRepository) Save(ctx context.Context, user *domain.User) error
func (h *UserHandler) CreateUser(c *gin.Context)  // Gin ya provee context

// ❌ INCORRECTO: Métodos sin context
func (s *UserService) CreateUser(dto CreateUserDTO) (*UserDTO, error)
```

### Logging Estructurado

```go
// ✅ CORRECTO: Usar logger de edugo-shared con campos estructurados
import (
    "github.com/EduGoGroup/edugo-shared/logger"
    "go.uber.org/zap"
)

func (s *UserService) CreateUser(ctx context.Context, dto CreateUserDTO) (*UserDTO, error) {
    logger.Info(ctx, "creating user",
        zap.String("email", dto.Email),
        zap.String("role", dto.Role),
    )

    // ... lógica ...

    if err != nil {
        logger.Error(ctx, "failed to create user",
            zap.Error(err),
            zap.String("email", dto.Email),
        )
        return nil, err
    }

    logger.Info(ctx, "user created successfully", zap.String("user_id", user.ID))
    return user, nil
}

// ❌ INCORRECTO: NO usar log estándar
log.Println("user created:", userID)
log.Printf("error: %v", err)

// ❌ INCORRECTO: NO usar fmt.Println
fmt.Println("creating user...")
```

---

## 🔐 Autenticación y Seguridad

### JWT con edugo-shared

```go
// ✅ CORRECTO: Usar JWTManager de edugo-shared
import "github.com/EduGoGroup/edugo-shared/auth"

jwtManager := auth.NewJWTManager(jwtSecret, 15*time.Minute)
token, err := jwtManager.GenerateToken(userID, email, roles)
```

### Passwords con Bcrypt

```go
// ✅ CORRECTO: bcrypt para hash de passwords
import "golang.org/x/crypto/bcrypt"

hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// Verificar
err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))

// ❌ INCORRECTO: NO usar SHA256 para passwords
hash := sha256.Sum256([]byte(password))  // ❌ Inseguro
```

### Middleware JWT

```go
// ✅ CORRECTO: Usar middleware de edugo-shared
import "github.com/EduGoGroup/edugo-shared/middleware"

protected := router.Group("/api/v1")
protected.Use(middleware.JWTAuthMiddleware(jwtManager))
{
    protected.POST("/courses", courseHandler.CreateCourse)
}
```

### Rate Limiting Anti Fuerza Bruta

```go
// ✅ Implementado en login
// Máximo 5 intentos fallidos en 15 minutos
// Bloqueo temporal de 30 minutos después
// Ver: internal/application/auth_service.go
```

---

## 🗄️ Bases de Datos

### PostgreSQL (Datos Relacionales)

```go
// ✅ Usar sqlx para queries
type PostgresUserRepository struct {
    db *sqlx.DB
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
    var user domain.User
    query := `SELECT id, email, password_hash, created_at FROM users WHERE id = $1`
    err := r.db.GetContext(ctx, &user, query, id)
    if err == sql.ErrNoRows {
        return nil, errors.NewNotFoundError("user", id)
    }
    return &user, err
}
```

### MongoDB (Datos No Estructurados)

```go
// ✅ Usar mongo-driver oficial
import "go.mongodb.org/mongo-driver/mongo"

type MongoLogRepository struct {
    collection *mongo.Collection
}
```

### RabbitMQ (Messaging)

```
⚠️ PENDIENTE DE IMPLEMENTAR
- Ver TODO en: internal/infrastructure/messaging/rabbitmq.go
```

---

## ✅ Testing

### Principios de Testing

```go
// ✅ Tests de integración con testcontainers
import (
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestUserRepository_Integration(t *testing.T) {
    // Setup: Levantar PostgreSQL container
    ctx := context.Background()
    container, err := postgres.RunContainer(ctx, ...)
    require.NoError(t, err)
    defer container.Terminate(ctx)

    // Test: Usar repositorio real
    repo := NewPostgresUserRepository(db)
    // ...

    // Cleanup: Automático con defer
}

// ✅ Tests unitarios con mocks para dependencias externas
type MockUserRepository struct {
    mock.Mock
}

// ✅ Tests deben ser independientes y ejecutarse en paralelo
func TestUserService_CreateUser(t *testing.T) {
    t.Parallel()  // ✅ Permite ejecución paralela
    // ...
}
```

### Cobertura de Tests

- **Objetivo**: >70% de cobertura
- **Actual**: 0.0% (proyecto en construcción)
- **Prioridad**: Servicios de aplicación y repositorios

---

## 🚨 Deuda Técnica Conocida

### 1. Handlers Duplicados (ALTA PRIORIDAD - Fase 3 del Sprint)

```
❌ internal/handlers/              # Handlers VIEJOS con mocks
   ├── user_handler.go            # NO USAR
   ├── course_handler.go          # NO USAR
   └── auth_handler.go            # NO USAR

✅ internal/infrastructure/http/handler/   # Handlers NUEVOS reales
   ├── user_handler.go            # ✅ USAR ESTOS
   ├── course_handler.go          # ✅ USAR ESTOS
   └── auth_handler.go            # ✅ USAR ESTOS
```

**Acción**: Los handlers viejos serán eliminados en Fase 3.

### 2. TODOs Pendientes

#### Funcionalidad S3 (AWS Storage)
```go
// Ver: internal/infrastructure/storage/aws-storage.go
// TODO: Implementar upload, download, delete de archivos
```

#### RabbitMQ Messaging
```go
// Ver: internal/infrastructure/messaging/rabbitmq.go
// TODO: Conectar producer y consumer
```

#### Queries Complejas en Repositorios
```go
// Ver: internal/infrastructure/persistence/
// TODO: Implementar búsquedas avanzadas, filtros, paginación
```

---

## 🎯 Flujo de Trabajo y Sprint Actual

### Antes de Sugerir Cambios

1. ✅ **SIEMPRE** revisar `sprint/README.md` para conocer fase actual
2. ✅ Verificar que el cambio esté alineado con el sprint
3. ✅ NO sugerir cambios que rompan la arquitectura limpia
4. ✅ Priorizar uso de edugo-shared sobre reimplementación

### Sprint Actual: Conectar Implementación Real

```
Fase 1: ✅ Completada (Conectar Container DI)
Fase 2: ⏳ En progreso (Completar TODOs de servicios)
   └─ Subtarea actual: Implementar S3
Fase 3: ⏸️ Pendiente (Eliminar handlers duplicados)
```

### Estado del Branch

- **Branch actual**: `feature/conectar`
- **Último commit**: `5d9e3ca` - "chore: sincronizar go.mod y go.sum"
- **Progreso**: 16.6% del sprint (1/6 commits completados)

---

## 🛠️ Tecnologías y Stack

### Framework y Bibliotecas Core
- **Framework Web**: Gin Gonic
- **Config Management**: Viper
- **Logging**: Zap (via edugo-shared)
- **Database Drivers**:
  - PostgreSQL: `lib/pq` + `sqlx`
  - MongoDB: `mongo-driver`

### Autenticación y Seguridad
- **JWT**: Via `edugo-shared/auth`
- **Password Hashing**: `bcrypt`
- **Rate Limiting**: Redis (implementado)

### Testing
- **Framework**: Testing estándar de Go
- **Containers**: Testcontainers
- **Mocking**: Testify/mock

### DevOps
- **Containerización**: Docker + Docker Compose
- **CI/CD**: GitHub Actions
- **Registry**: GitHub Container Registry (ghcr.io)

---

## 📚 Documentación API

### Swagger/OpenAPI

```go
// ✅ CORRECTO: Agregar anotaciones Swagger en handlers
// @Summary Crear nuevo usuario
// @Description Crea un usuario en el sistema
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserDTO true "Datos del usuario"
// @Success 201 {object} UserDTO
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
    // ...
}
```

### Generar Documentación

```bash
# Regenerar docs Swagger
swag init -g cmd/main.go --output docs

# Acceder a Swagger UI
# http://localhost:8080/swagger/index.html
```

---

## 🌐 Variables de Entorno

### Variables Requeridas

```bash
# Base de datos
POSTGRES_PASSWORD=<contraseña>
MONGODB_URI=mongodb://<user>:<pass>@<host>:<port>

# Messaging
RABBITMQ_URL=amqp://<user>:<pass>@<host>:<port>

# Autenticación
JWT_SECRET=<secret_seguro>

# Ambiente
APP_ENV=local|dev|qa|prod
```

### NO Hardcodear Secrets

```go
// ❌ INCORRECTO: Secrets hardcodeados
const jwtSecret = "mi-secret-super-seguro"
const dbPassword = "postgres123"

// ✅ CORRECTO: Leer de variables de entorno
jwtSecret := viper.GetString("jwt.secret")
dbPassword := viper.GetString("database.password")
```

---

## 🎨 Estilo de Código

### Formato

```bash
# ✅ SIEMPRE formatear con gofmt antes de commit
gofmt -w .

# ✅ Verificar con linter
golangci-lint run
```

### Comentarios

```go
// ✅ CORRECTO: Comentarios en español, explicativos
// CreateUser crea un nuevo usuario en el sistema y envía un email de bienvenida.
// Valida que el email sea único antes de crear el registro.
func (s *UserService) CreateUser(ctx context.Context, dto CreateUserDTO) (*UserDTO, error)

// ❌ INCORRECTO: Comentarios obvios o redundantes
// CreateUser crea un usuario
func (s *UserService) CreateUser(...)
```

### Imports

```go
// ✅ CORRECTO: Agrupar imports
import (
    // Standard library
    "context"
    "fmt"
    "time"

    // Third party
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"

    // Internal - edugo-shared
    "github.com/EduGoGroup/edugo-shared/auth"
    "github.com/EduGoGroup/edugo-shared/logger"

    // Internal - este proyecto
    "github.com/EduGoGroup/edugo-api-mobile/internal/domain"
    "github.com/EduGoGroup/edugo-api-mobile/internal/application"
)
```

---

## ⚡ Mejores Prácticas Adicionales

### 1. Inyección de Dependencias

```go
// ✅ CORRECTO: Constructor con dependencias explícitas
func NewUserService(
    repo UserRepository,
    logger logger.Logger,
    emailService EmailService,
) *UserService {
    return &UserService{
        repo:         repo,
        logger:       logger,
        emailService: emailService,
    }
}

// ❌ INCORRECTO: Dependencias globales o singleton
var globalDB *sql.DB  // ❌ Evitar
```

### 2. Validación de DTOs

```go
// ✅ CORRECTO: Usar validaciones explícitas
import "github.com/go-playground/validator/v10"

type CreateUserDTO struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Name     string `json:"name" validate:"required,min=2"`
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var dto CreateUserDTO
    if err := c.ShouldBindJSON(&dto); err != nil {
        c.JSON(400, gin.H{"error": "invalid request body"})
        return
    }

    if err := validate.Struct(dto); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    // ...
}
```

### 3. Transacciones de Base de Datos

```go
// ✅ CORRECTO: Usar transacciones para operaciones múltiples
func (s *UserService) CreateUserWithProfile(ctx context.Context, dto CreateUserDTO) error {
    tx, err := s.db.BeginTxx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()  // Rollback automático si no hay commit

    // Operación 1
    user, err := s.userRepo.SaveTx(ctx, tx, user)
    if err != nil {
        return err
    }

    // Operación 2
    err = s.profileRepo.SaveTx(ctx, tx, profile)
    if err != nil {
        return err
    }

    return tx.Commit()
}
```

---

## 🎓 Recursos de Referencia

- **Clean Architecture**: [sprint/plan-autenticacion-oauth2.md](../sprint/plan-autenticacion-oauth2.md)
- **Sprint Actual**: [sprint/README.md](../sprint/README.md)
- **Configuración Proyecto**: [.claude/CLAUDE.md](../.claude/CLAUDE.md)
- **Workflows CI/CD**: [.github/workflows/README.md](workflows/README.md)

---

## 📝 Notas Finales para Copilot

### Al Revisar Pull Requests

1. ✅ Verificar que se usen tipos de error de `edugo-shared`
2. ✅ Confirmar que todos los métodos reciben `context.Context`
3. ✅ Validar que se use logging estructurado
4. ✅ Detectar código duplicado entre handlers viejos y nuevos
5. ✅ Señalar TODOs o funcionalidad incompleta
6. ✅ Verificar que no se reimplemente funcionalidad de `edugo-shared`

### Al Sugerir Código

1. ✅ Seguir Clean Architecture (no mezclar capas)
2. ✅ Usar dependencias de `edugo-shared` cuando corresponda
3. ✅ Incluir logging adecuado
4. ✅ Manejar errores con tipos apropiados
5. ✅ Agregar validaciones necesarias
6. ✅ Escribir código testeable

### Recordatorio de Idioma

🌍 **TODOS los comentarios, sugerencias y explicaciones deben estar en ESPAÑOL.**

---

**Última actualización**: 2025-11-01
**Versión del proyecto**: En desarrollo (Sprint Fase 2)
**Go Version**: 1.25.3
**edugo-shared Version**: v2.0.5
