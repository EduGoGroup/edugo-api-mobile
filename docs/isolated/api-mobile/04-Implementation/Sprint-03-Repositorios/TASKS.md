# Tareas del Sprint 03 - Repositorios

## Objetivo
Implementar la capa de persistencia con repositorios PostgreSQL (usando GORM) y MongoDB (usando driver oficial) que implementan las interfaces definidas en Sprint-02. Incluir tests de integración con Testcontainers para validar queries reales contra bases de datos.

---

## Tareas

### TASK-03-001: Implementar PostgresAssessmentRepository
**Tipo:** feature  
**Prioridad:** HIGH  
**Estimación:** 3h  
**Asignado a:** @ai-executor

#### Descripción
Implementar repositorio PostgreSQL para entity Assessment usando GORM como ORM.

#### Pasos de Implementación

1. Crear archivo:
   `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/internal/infrastructure/persistence/postgres_assessment_repository.go`

2. Implementar repositorio con GORM:
   ```go
   package persistence
   
   import (
       "context"
       "errors"
       
       "github.com/google/uuid"
       "gorm.io/gorm"
       
       "edugo-api-mobile/internal/domain/entities"
       "edugo-api-mobile/internal/domain/repositories"
   )
   
   // PostgresAssessmentRepository implementa repositories.AssessmentRepository
   type PostgresAssessmentRepository struct {
       db *gorm.DB
   }
   
   // NewPostgresAssessmentRepository crea un nuevo repositorio
   func NewPostgresAssessmentRepository(db *gorm.DB) repositories.AssessmentRepository {
       return &PostgresAssessmentRepository{db: db}
   }
   
   // FindByID busca un assessment por ID
   func (r *PostgresAssessmentRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Assessment, error) {
       var assessment entities.Assessment
       
       result := r.db.WithContext(ctx).
           Where("id = ?", id).
           First(&assessment)
       
       if result.Error != nil {
           if errors.Is(result.Error, gorm.ErrRecordNotFound) {
               return nil, nil // No encontrado
           }
           return nil, result.Error
       }
       
       return &assessment, nil
   }
   
   // FindByMaterialID busca assessment por material ID
   func (r *PostgresAssessmentRepository) FindByMaterialID(ctx context.Context, materialID uuid.UUID) (*entities.Assessment, error) {
       var assessment entities.Assessment
       
       result := r.db.WithContext(ctx).
           Where("material_id = ?", materialID).
           First(&assessment)
       
       if result.Error != nil {
           if errors.Is(result.Error, gorm.ErrRecordNotFound) {
               return nil, nil
           }
           return nil, result.Error
       }
       
       return &assessment, nil
   }
   
   // Save guarda un assessment (INSERT o UPDATE)
   func (r *PostgresAssessmentRepository) Save(ctx context.Context, assessment *entities.Assessment) error {
       return r.db.WithContext(ctx).Save(assessment).Error
   }
   
   // Delete elimina un assessment
   func (r *PostgresAssessmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
       result := r.db.WithContext(ctx).
           Delete(&entities.Assessment{}, "id = ?", id)
       
       if result.Error != nil {
           return result.Error
       }
       
       if result.RowsAffected == 0 {
           return errors.New("assessment not found")
       }
       
       return nil
   }
   ```

#### Criterios de Aceptación
- [ ] Implementa interface `repositories.AssessmentRepository`
- [ ] Usa `context.Context` en todos los métodos
- [ ] Maneja `gorm.ErrRecordNotFound` correctamente (retorna nil, no error)
- [ ] Usa `WithContext()` para timeout/cancellation
- [ ] Métodos retornan errores de GORM sin wrapping innecesario

#### Comandos de Validación
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Compilar
go build ./internal/infrastructure/persistence

# Verificar que implementa interface
go test ./internal/infrastructure/persistence -v -run TestPostgresAssessmentRepository_Interface
```

#### Dependencias
- Requiere: Sprint-02 completado (interface AssessmentRepository)
- Requiere: GORM instalado (`go get gorm.io/gorm@v1.25.5`)
- Usa: `gorm.io/driver/postgres` para conexión

#### Tiempo Estimado
3 horas

---

### TASK-03-002: Implementar PostgresAttemptRepository con Transacciones
**Tipo:** feature  
**Prioridad:** HIGH  
**Estimación:** 4h  
**Asignado a:** @ai-executor

#### Descripción
Implementar repositorio para Attempt con soporte de transacciones ACID. Al guardar un Attempt, debe guardar también todas sus Answers en una sola transacción.

#### Pasos de Implementación

1. Crear archivo:
   `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/internal/infrastructure/persistence/postgres_attempt_repository.go`

2. Implementar con transacciones:
   ```go
   package persistence
   
   import (
       "context"
       "errors"
       
       "github.com/google/uuid"
       "gorm.io/gorm"
       
       "edugo-api-mobile/internal/domain/entities"
       "edugo-api-mobile/internal/domain/repositories"
   )
   
   type PostgresAttemptRepository struct {
       db *gorm.DB
   }
   
   func NewPostgresAttemptRepository(db *gorm.DB) repositories.AttemptRepository {
       return &PostgresAttemptRepository{db: db}
   }
   
   // FindByID busca un intento por ID
   func (r *PostgresAttemptRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Attempt, error) {
       var attempt entities.Attempt
       
       // Cargar attempt con sus answers (eager loading)
       result := r.db.WithContext(ctx).
           Preload("Answers").  // GORM preload de relación
           Where("id = ?", id).
           First(&attempt)
       
       if result.Error != nil {
           if errors.Is(result.Error, gorm.ErrRecordNotFound) {
               return nil, nil
           }
           return nil, result.Error
       }
       
       return &attempt, nil
   }
   
   // Save guarda un attempt CON sus answers en transacción ACID
   func (r *PostgresAttemptRepository) Save(ctx context.Context, attempt *entities.Attempt) error {
       // Usar transacción para atomicidad
       return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
           // 1. Guardar attempt
           if err := tx.Create(attempt).Error; err != nil {
               return err // Rollback automático
           }
           
           // 2. Guardar todas las answers
           for _, answer := range attempt.Answers {
               answer.AttemptID = attempt.ID // Asegurar FK
               if err := tx.Create(answer).Error; err != nil {
                   return err // Rollback automático
               }
           }
           
           // 3. Commit si todo OK
           return nil
       })
   }
   
   // FindByStudentAndAssessment busca intentos de un estudiante en un assessment
   func (r *PostgresAttemptRepository) FindByStudentAndAssessment(
       ctx context.Context,
       studentID, assessmentID uuid.UUID,
   ) ([]*entities.Attempt, error) {
       var attempts []*entities.Attempt
       
       result := r.db.WithContext(ctx).
           Preload("Answers").
           Where("student_id = ? AND assessment_id = ?", studentID, assessmentID).
           Order("created_at DESC").
           Find(&attempts)
       
       if result.Error != nil {
           return nil, result.Error
       }
       
       return attempts, nil
   }
   
   // CountByStudentAndAssessment cuenta intentos de un estudiante
   func (r *PostgresAttemptRepository) CountByStudentAndAssessment(
       ctx context.Context,
       studentID, assessmentID uuid.UUID,
   ) (int, error) {
       var count int64
       
       result := r.db.WithContext(ctx).
           Model(&entities.Attempt{}).
           Where("student_id = ? AND assessment_id = ?", studentID, assessmentID).
           Count(&count)
       
       if result.Error != nil {
           return 0, result.Error
       }
       
       return int(count), nil
   }
   
   // FindByStudent busca todos los intentos de un estudiante (historial)
   func (r *PostgresAttemptRepository) FindByStudent(
       ctx context.Context,
       studentID uuid.UUID,
       limit, offset int,
   ) ([]*entities.Attempt, error) {
       var attempts []*entities.Attempt
       
       result := r.db.WithContext(ctx).
           Preload("Answers").
           Where("student_id = ?", studentID).
           Order("created_at DESC").
           Limit(limit).
           Offset(offset).
           Find(&attempts)
       
       if result.Error != nil {
           return nil, result.Error
       }
       
       return attempts, nil
   }
   ```

#### Criterios de Aceptación
- [ ] Usa `tx.Transaction()` para ACID al guardar Attempt+Answers
- [ ] Rollback automático si cualquier INSERT falla
- [ ] Usa `Preload("Answers")` para eager loading
- [ ] Método `CountByStudentAndAssessment` para verificar límite de intentos
- [ ] Paginación en `FindByStudent` (limit/offset)

#### Comandos de Validación
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Test de transacción (debe crear attempt + answers atómicamente)
go test ./internal/infrastructure/persistence -v -run TestAttemptRepository_Transaction
```

#### Dependencias
- Requiere: TASK-03-001 completada
- Requiere: Entities Attempt y Answer de Sprint-02
- Usa: `gorm.io/gorm` transacciones

#### Tiempo Estimado
4 horas

---

### TASK-03-003: Implementar MongoQuestionRepository
**Tipo:** feature  
**Prioridad:** HIGH  
**Estimación:** 3h  
**Asignado a:** @ai-executor

#### Descripción
Implementar repositorio MongoDB para leer preguntas de evaluaciones desde colección `material_assessment`.

#### Pasos de Implementación

1. Crear archivo:
   `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/internal/infrastructure/persistence/mongo_question_repository.go`

2. Implementar con MongoDB driver:
   ```go
   package persistence
   
   import (
       "context"
       
       "go.mongodb.org/mongo-driver/bson"
       "go.mongodb.org/mongo-driver/bson/primitive"
       "go.mongodb.org/mongo-driver/mongo"
   )
   
   // Question representa una pregunta de MongoDB
   type Question struct {
       ID          string    `bson:"id"`
       Text        string    `bson:"text"`
       Type        string    `bson:"type"`
       Options     []Option  `bson:"options"`
       // NO incluir correct_answer (seguridad)
   }
   
   type Option struct {
       ID   string `bson:"id"`
       Text string `bson:"text"`
   }
   
   // MaterialAssessment representa documento de MongoDB
   type MaterialAssessment struct {
       ID              primitive.ObjectID `bson:"_id"`
       MaterialID      string             `bson:"material_id"`
       Questions       []Question         `bson:"questions"`
       TotalQuestions  int                `bson:"total_questions"`
   }
   
   type MongoQuestionRepository struct {
       collection *mongo.Collection
   }
   
   func NewMongoQuestionRepository(db *mongo.Database) *MongoQuestionRepository {
       return &MongoQuestionRepository{
           collection: db.Collection("material_assessment"),
       }
   }
   
   // FindQuestionsByMaterialID obtiene preguntas (SIN respuestas correctas)
   func (r *MongoQuestionRepository) FindQuestionsByMaterialID(
       ctx context.Context,
       materialID string,
   ) ([]Question, error) {
       var assessment MaterialAssessment
       
       // Buscar por material_id
       filter := bson.M{"material_id": materialID}
       
       // Proyección: EXCLUIR respuestas correctas por seguridad
       projection := bson.M{
           "questions.correct_answer": 0,  // NO enviar al cliente
           "questions.feedback":       0,  // NO enviar feedback
       }
       
       err := r.collection.FindOne(ctx, filter).Decode(&assessment)
       if err != nil {
           if err == mongo.ErrNoDocuments {
               return nil, nil // No encontrado
           }
           return nil, err
       }
       
       return assessment.Questions, nil
   }
   
   // FindByMongoDocumentID busca por ObjectID de MongoDB
   func (r *MongoQuestionRepository) FindByMongoDocumentID(
       ctx context.Context,
       mongoDocID string,
   ) (*MaterialAssessment, error) {
       objectID, err := primitive.ObjectIDFromHex(mongoDocID)
       if err != nil {
           return nil, err
       }
       
       var assessment MaterialAssessment
       filter := bson.M{"_id": objectID}
       projection := bson.M{
           "questions.correct_answer": 0,
           "questions.feedback":       0,
       }
       
       opts := mongo.NewSingleResultOptions().SetProjection(projection)
       err = r.collection.FindOne(ctx, filter, opts).Decode(&assessment)
       if err != nil {
           if err == mongo.ErrNoDocuments {
               return nil, nil
           }
           return nil, err
       }
       
       return &assessment, nil
   }
   ```

#### Criterios de Aceptación
- [ ] NUNCA expone `correct_answer` al cliente (proyección MongoDB)
- [ ] NUNCA expone `feedback` al cliente
- [ ] Maneja `mongo.ErrNoDocuments` correctamente (retorna nil)
- [ ] Usa `context.Context` para timeout
- [ ] Convierte ObjectID correctamente con `primitive.ObjectIDFromHex()`

#### Comandos de Validación
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Test de seguridad: verificar que correct_answer NO se retorna
go test ./internal/infrastructure/persistence -v -run TestMongoRepository_Security
```

#### Dependencias
- Requiere: MongoDB 7.0+ corriendo
- Requiere: `go.mongodb.org/mongo-driver/mongo` v1.13.1+
- Usa: Colección `material_assessment` poblada

#### Tiempo Estimado
3 horas

---

### TASK-03-004: Tests de Integración con Testcontainers
**Tipo:** test  
**Prioridad:** HIGH  
**Estimación:** 4h  
**Asignado a:** @ai-executor

#### Descripción
Crear tests de integración que usan PostgreSQL y MongoDB reales (via Testcontainers) en lugar de mocks.

#### Pasos de Implementación

1. Crear archivo:
   `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/tests/integration/repository_test.go`

2. Implementar con Testcontainers:
   ```go
   //go:build integration
   // +build integration
   
   package integration_test
   
   import (
       "context"
       "testing"
       
       "github.com/google/uuid"
       "github.com/stretchr/testify/assert"
       "github.com/stretchr/testify/require"
       "github.com/testcontainers/testcontainers-go"
       "github.com/testcontainers/testcontainers-go/modules/postgres"
       "gorm.io/driver/postgres"
       "gorm.io/gorm"
       
       "edugo-api-mobile/internal/domain/entities"
       "edugo-api-mobile/internal/infrastructure/persistence"
   )
   
   func TestPostgresAssessmentRepository_Integration(t *testing.T) {
       ctx := context.Background()
       
       // 1. Iniciar contenedor PostgreSQL
       pgContainer, err := postgres.RunContainer(ctx,
           testcontainers.WithImage("postgres:15-alpine"),
           postgres.WithDatabase("testdb"),
           postgres.WithUsername("test"),
           postgres.WithPassword("test"),
       )
       require.NoError(t, err)
       defer pgContainer.Terminate(ctx)
       
       // 2. Obtener connection string
       connStr, err := pgContainer.ConnectionString(ctx)
       require.NoError(t, err)
       
       // 3. Conectar con GORM
       db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
       require.NoError(t, err)
       
       // 4. Migrar schema
       err = db.AutoMigrate(&entities.Assessment{})
       require.NoError(t, err)
       
       // 5. Crear repositorio
       repo := persistence.NewPostgresAssessmentRepository(db)
       
       // 6. Test: Guardar y recuperar
       assessment, _ := entities.NewAssessment(
           uuid.New(),
           "507f1f77bcf86cd799439011",
           "Test Assessment",
           5,
           70,
       )
       
       err = repo.Save(ctx, assessment)
       require.NoError(t, err)
       
       // 7. Recuperar
       found, err := repo.FindByID(ctx, assessment.ID)
       require.NoError(t, err)
       require.NotNil(t, found)
       
       assert.Equal(t, assessment.Title, found.Title)
       assert.Equal(t, assessment.TotalQuestions, found.TotalQuestions)
   }
   
   func TestPostgresAttemptRepository_Transaction_Integration(t *testing.T) {
       // ... similar setup con Testcontainers ...
       
       // Test que transacción funciona: crear attempt con 5 answers
       attempt, _ := entities.NewAttempt(...)
       
       err := repo.Save(ctx, attempt)
       require.NoError(t, err)
       
       // Verificar que attempt Y answers se guardaron
       found, err := repo.FindByID(ctx, attempt.ID)
       require.NoError(t, err)
       assert.Len(t, found.Answers, 5) // Todas las answers guardadas
   }
   ```

#### Criterios de Aceptación
- [ ] Usa tag `//go:build integration` para separar de unit tests
- [ ] Levanta PostgreSQL con Testcontainers
- [ ] Ejecuta migraciones automáticamente (AutoMigrate)
- [ ] Tests de transacciones verifican atomicidad
- [ ] Limpia contenedores después de tests (defer Terminate)

#### Comandos de Validación
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Ejecutar SOLO tests de integración (requiere Docker)
go test ./tests/integration -v -tags=integration

# Verificar que contenedores se limpian
docker ps -a | grep testcontainers
# Debe estar vacío después de tests
```

#### Dependencias
- Requiere: Docker instalado y corriendo
- Requiere: `github.com/testcontainers/testcontainers-go` v0.27.0+
- Usa: `testcontainers/modules/postgres`

#### Tiempo Estimado
4 horas

---

### TASK-03-005: Configuración de Connection Pool
**Tipo:** feature  
**Prioridad:** MEDIUM  
**Estimación:** 2h  
**Asignado a:** @ai-executor

#### Descripción
Configurar pool de conexiones para PostgreSQL y MongoDB con parámetros optimizados por ambiente.

#### Pasos de Implementación

1. Crear archivo:
   `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/internal/infrastructure/database/postgres.go`

2. Implementar configuración:
   ```go
   package database
   
   import (
       "fmt"
       "time"
       
       "gorm.io/driver/postgres"
       "gorm.io/gorm"
   )
   
   type PostgresConfig struct {
       Host         string
       Port         int
       User         string
       Password     string
       Database     string
       MaxOpenConns int
       MaxIdleConns int
       MaxLifetime  time.Duration
   }
   
   func NewPostgresConnection(cfg PostgresConfig) (*gorm.DB, error) {
       dsn := fmt.Sprintf(
           "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
           cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database,
       )
       
       db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
       if err != nil {
           return nil, err
       }
       
       // Configurar connection pool
       sqlDB, err := db.DB()
       if err != nil {
           return nil, err
       }
       
       sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
       sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
       sqlDB.SetConnMaxLifetime(cfg.MaxLifetime)
       
       // Ping para verificar conexión
       if err := sqlDB.Ping(); err != nil {
           return nil, err
       }
       
       return db, nil
   }
   ```

3. Configuración por ambiente:
   ```go
   // config/database.go
   func GetPostgresConfig(env string) PostgresConfig {
       switch env {
       case "production":
           return PostgresConfig{
               MaxOpenConns: 25,
               MaxIdleConns: 10,
               MaxLifetime:  5 * time.Minute,
           }
       case "development":
           return PostgresConfig{
               MaxOpenConns: 10,
               MaxIdleConns: 5,
               MaxLifetime:  3 * time.Minute,
           }
       default: // test
           return PostgresConfig{
               MaxOpenConns: 5,
               MaxIdleConns: 2,
               MaxLifetime:  1 * time.Minute,
           }
       }
   }
   ```

#### Criterios de Aceptación
- [ ] Connection pool configurado con `SetMaxOpenConns`, `SetMaxIdleConns`, `SetConnMaxLifetime`
- [ ] Configuración diferente por ambiente (dev, prod, test)
- [ ] Ping inicial para verificar conexión
- [ ] Similar para MongoDB con `mongo.Client.SetMaxPoolSize()`

#### Comandos de Validación
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Verificar que conexión funciona
go run cmd/api/main.go --check-db
```

#### Tiempo Estimado
2 horas

---

## Resumen del Sprint

**Total de Tareas:** 5  
**Estimación Total:** 16 horas  
**Archivos a Crear:** ~8 archivos Go + tests

**Entregables:**
1. 3 repositorios implementados (PostgreSQL × 2, MongoDB × 1)
2. Tests de integración con Testcontainers
3. Connection pooling configurado
4. Transacciones ACID en Attempt+Answers

**Criterios de Éxito:**
- [ ] Repositorios implementan interfaces de Sprint-02
- [ ] Tests de integración pasando (coverage >70%)
- [ ] Transacciones funcionando atómicamente
- [ ] Connection pool optimizado
- [ ] Seguridad: respuestas correctas NUNCA expuestas

---

**Generado con:** Claude Code  
**Sprint:** 03/06  
**Última actualización:** 2025-11-14
