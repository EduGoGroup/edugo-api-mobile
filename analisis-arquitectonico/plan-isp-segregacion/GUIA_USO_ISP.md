# 📖 Guía de Uso: Interfaces Segregadas (ISP)

**Proyecto**: EduGo API Mobile  
**Fecha**: 2025-11-06  
**Principio**: Interface Segregation Principle (ISP)

---

## 🎯 Introducción

Todos los repositorios del proyecto implementan el **Principio de Segregación de Interfaces (ISP)**. Cada repositorio expone interfaces específicas (Reader, Writer, Stats, etc.) que permiten a los services depender solo de las operaciones que necesitan.

---

## 📋 Catálogo de Interfaces

### 1. UserRepository

**Ubicación**: `internal/domain/repository/user_repository.go`

#### UserReader
```go
type UserReader interface {
    FindByID(ctx context.Context, id valueobject.UserID) (*entity.User, error)
    FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error)
}
```

**Cuándo usar**:
- ✅ Services que solo consultan usuarios
- ✅ Servicios de perfiles
- ✅ Servicios de búsqueda

**Ejemplo**:
```go
type UserProfileService struct {
    userReader repository.UserReader  // Solo necesita leer
}
```

#### UserWriter
```go
type UserWriter interface {
    Update(ctx context.Context, user *entity.User) error
}
```

**Cuándo usar**:
- ✅ Services que solo modifican usuarios
- ✅ Servicios de actualización de perfil

#### UserRepository (Completo)
```go
type UserRepository interface {
    UserReader
    UserWriter
}
```

**Cuándo usar**:
- ✅ Services que necesitan leer Y escribir
- ✅ Servicio de autenticación (lee para validar, escribe para actualizar)

---

### 2. MaterialRepository

**Ubicación**: `internal/domain/repository/material_repository.go`

#### MaterialReader (4 métodos)
- `FindByID`: Buscar material por ID
- `FindByIDWithVersions`: Material con historial de versiones
- `List`: Listar con filtros
- `FindByAuthor`: Materiales de un autor

**Cuándo usar**:
- ✅ Servicios de consulta de materiales
- ✅ Servicios de búsqueda
- ✅ APIs de lectura

#### MaterialWriter (4 métodos)
- `Create`: Crear nuevo material
- `Update`: Actualizar material
- `UpdateStatus`: Cambiar estado (draft/published)
- `UpdateProcessingStatus`: Estado de procesamiento

**Cuándo usar**:
- ✅ Servicios de creación de materiales
- ✅ Servicios de procesamiento
- ✅ Workers que actualizan estados

#### MaterialStats (1 método)
- `CountPublishedMaterials`: Contar materiales publicados

**Cuándo usar**:
- ✅ Servicios de estadísticas
- ✅ Dashboards
- ✅ Reportes

---

### 3. ProgressRepository

**Ubicación**: `internal/domain/repository/progress_repository.go`

#### ProgressReader (1 método)
- `FindByMaterialAndUser`: Buscar progreso específico

**Cuándo usar**:
- ✅ Consultar progreso de un usuario
- ✅ Verificar si completó un material

#### ProgressWriter (3 métodos)
- `Save`: Guardar nuevo progreso
- `Update`: Actualizar progreso existente
- `Upsert`: INSERT o UPDATE idempotente

**Cuándo usar**:
- ✅ Actualizar progreso de usuario
- ✅ Registrar avances
- **Recomendado**: Usar `Upsert` para simplificar lógica

#### ProgressStats (2 métodos)
- `CountActiveUsers`: Usuarios activos (últimos 30 días)
- `CalculateAverageProgress`: Promedio de progreso global

**Cuándo usar**:
- ✅ Dashboards de engagement
- ✅ Métricas de adopción

---

### 4. AssessmentRepository

**Ubicación**: `internal/domain/repository/assessment_repository.go`

#### AssessmentReader (3 métodos)
- `FindAssessmentByMaterialID`: Obtener quiz de material
- `FindAttemptsByUser`: Historial de intentos
- `GetBestAttempt`: Mejor intento del usuario

#### AssessmentWriter (3 métodos)
- `SaveAssessment`: Guardar/actualizar quiz
- `SaveAttempt`: Registrar intento
- `SaveResult`: Guardar resultado final

#### AssessmentStats (2 métodos)
- `CountCompletedAssessments`: Total evaluaciones completadas
- `CalculateAverageScore`: Promedio de puntajes

---

### 5. RefreshTokenRepository

**Ubicación**: `internal/domain/repository/refresh_token_repository.go`

#### RefreshTokenReader (1 método)
- `FindByTokenHash`: Buscar token por hash

**Cuándo usar**: Validar refresh tokens

#### RefreshTokenWriter (3 métodos)
- `Store`: Guardar nuevo token
- `Revoke`: Revocar token (logout)
- `RevokeAllByUserID`: Cerrar todas las sesiones

#### RefreshTokenMaintenance (1 método)
- `DeleteExpired`: Limpieza de tokens expirados (cron job)

---

### 6. SummaryRepository

**Ubicación**: `internal/domain/repository/summary_repository.go`

#### SummaryReader (2 métodos)
- `FindByMaterialID`: Obtener resumen
- `Exists`: Verificar si existe

#### SummaryWriter (2 métodos)
- `Save`: Guardar/actualizar resumen
- `Delete`: Eliminar resumen

---

### 7. LoginAttemptRepository

**Ubicación**: `internal/domain/repository/login_attempt_repository.go`

#### LoginAttemptReader (2 métodos)
- `CountFailedAttempts`: Contar intentos fallidos recientes
- `IsRateLimited`: Verificar si está bloqueado

**Cuándo usar**: Rate limiting de login

#### LoginAttemptWriter (1 método)
- `RecordAttempt`: Registrar intento (exitoso o fallido)

---

## 🎨 Patrones de Uso

### Patrón 1: Service de Solo Lectura

```go
type MaterialSearchService struct {
    materialReader repository.MaterialReader  // Solo 4 métodos
    logger         logger.Logger
}

func NewMaterialSearchService(
    reader repository.MaterialReader,  // ← Interfaz pequeña
    logger logger.Logger,
) *MaterialSearchService {
    return &MaterialSearchService{
        materialReader: reader,
        logger:         logger,
    }
}

func (s *MaterialSearchService) Search(ctx context.Context, filters repository.ListFilters) ([]*entity.Material, error) {
    return s.materialReader.List(ctx, filters)
}
```

**Beneficios**:
- ✅ Dependencia mínima (solo 4 métodos en vez de 9)
- ✅ Tests más simples (mock de 4 métodos)
- ✅ Clara responsabilidad (solo lectura)

---

### Patrón 2: Service de Solo Escritura

```go
type MaterialCreationService struct {
    materialWriter repository.MaterialWriter  // Solo 4 métodos
    publisher      messaging.Publisher
    logger         logger.Logger
}

func (s *MaterialCreationService) CreateMaterial(ctx context.Context, material *entity.Material) error {
    // Crear material
    if err := s.materialWriter.Create(ctx, material); err != nil {
        return err
    }
    
    // Publicar evento
    return s.publisher.Publish("material.created", material)
}
```

**Beneficios**:
- ✅ No puede leer accidentalmente
- ✅ Mock de 4 métodos en tests
- ✅ Principio de mínimo privilegio

---

### Patrón 3: Service con Múltiples Capacidades

```go
type MaterialManagementService struct {
    materialRepo repository.MaterialRepository  // Interfaz completa (9 métodos)
    logger       logger.Logger
}

// Usa Reader para consultas
func (s *MaterialManagementService) GetMaterial(ctx context.Context, id valueobject.MaterialID) (*entity.Material, error) {
    return s.materialRepo.FindByID(ctx, id)
}

// Usa Writer para modificaciones
func (s *MaterialManagementService) UpdateMaterial(ctx context.Context, material *entity.Material) error {
    return s.materialRepo.Update(ctx, material)
}

// Usa Stats para métricas
func (s *MaterialManagementService) GetPublishedCount(ctx context.Context) (int64, error) {
    return s.materialRepo.CountPublishedMaterials(ctx)
}
```

**Cuándo usar**: Services "orchestrator" que necesitan acceso completo

---

### Patrón 4: Service con Stats

```go
type StatsService struct {
    materialStats   repository.MaterialStats     // Solo 1 método
    progressStats   repository.ProgressStats     // Solo 2 métodos
    assessmentStats repository.AssessmentStats   // Solo 2 métodos
}

func (s *StatsService) GetGlobalStats(ctx context.Context) (*dto.GlobalStatsResponse, error) {
    // Ejecutar queries en paralelo
    var wg sync.WaitGroup
    var mu sync.Mutex
    
    stats := &dto.GlobalStatsResponse{}
    
    wg.Add(3)
    
    // Material stats
    go func() {
        defer wg.Done()
        count, _ := s.materialStats.CountPublishedMaterials(ctx)
        mu.Lock()
        stats.TotalMaterials = count
        mu.Unlock()
    }()
    
    // Progress stats
    go func() {
        defer wg.Done()
        count, _ := s.progressStats.CountActiveUsers(ctx)
        mu.Lock()
        stats.ActiveUsers = count
        mu.Unlock()
    }()
    
    // Assessment stats
    go func() {
        defer wg.Done()
        count, _ := s.assessmentStats.CountCompletedAssessments(ctx)
        mu.Lock()
        stats.CompletedAssessments = count
        mu.Unlock()
    }()
    
    wg.Wait()
    return stats, nil
}
```

**Beneficios**:
- ✅ Solo depende de métodos de estadísticas
- ✅ Clara separación de concerns
- ✅ Tests simples (3 mocks pequeños)

---

## 🧪 Testing con Interfaces Segregadas

### Mock de Interfaz Pequeña

```go
// Mock PEQUEÑO (solo 2 métodos)
type MockUserReader struct {
    mock.Mock
}

func (m *MockUserReader) FindByID(ctx context.Context, id valueobject.UserID) (*entity.User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserReader) FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error) {
    args := m.Called(ctx, email)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.User), args.Error(1)
}

// Test
func TestUserProfileService_GetProfile(t *testing.T) {
    mockReader := new(MockUserReader)
    service := NewUserProfileService(mockReader)
    
    expectedUser := &entity.User{ID: testID, Email: "test@example.com"}
    mockReader.On("FindByID", mock.Anything, testID).Return(expectedUser, nil)
    
    user, err := service.GetProfile(context.Background(), testID)
    
    require.NoError(t, err)
    assert.Equal(t, expectedUser, user)
    mockReader.AssertExpectations(t)
}
```

**Comparación**:
- ❌ Mock completo: 9 métodos a implementar
- ✅ Mock Reader: 2 métodos (77% menos código)

---

## 💡 Mejores Prácticas

### ✅ DO:

1. **Usa la interfaz más pequeña posible**
   ```go
   // ✅ BIEN
   func NewSearchService(reader MaterialReader) *SearchService
   
   // ❌ MAL (si solo necesitas leer)
   func NewSearchService(repo MaterialRepository) *SearchService
   ```

2. **Compón interfaces cuando necesites múltiples**
   ```go
   // ✅ BIEN
   type MaterialService struct {
       reader MaterialReader
       writer MaterialWriter
   }
   ```

3. **Documenta qué interfaz necesitas**
   ```go
   // NewMaterialSearchService crea un servicio de búsqueda
   // Requiere: MaterialReader (solo lectura de materiales)
   func NewMaterialSearchService(reader repository.MaterialReader) *MaterialSearchService
   ```

### ❌ DON'T:

1. **No uses Repository completo si no lo necesitas**
   ```go
   // ❌ MAL (overkill para solo leer)
   type ViewService struct {
       repo MaterialRepository  // 9 métodos cuando solo usas 1
   }
   ```

2. **No mezcles concerns en un service**
   ```go
   // ❌ MAL
   type MixedService struct {
       materialRepo  MaterialRepository  // Todo
       progressRepo  ProgressRepository  // Todo
       assessmentRepo AssessmentRepository // Todo
   }
   
   // ✅ BIEN (específico)
   type StatsService struct {
       materialStats   MaterialStats
       progressStats   ProgressStats
       assessmentStats AssessmentStats
   }
   ```

---

## 📚 Referencias

- **Código**: `internal/domain/repository/*.go`
- **Implementaciones**: `internal/infrastructure/persistence/{postgres,mongodb}/*.go`
- **Services de ejemplo**: `internal/application/service/*.go`
- **Análisis completo**: `analisis-arquitectonico/plan-isp-segregacion/ANALISIS.md`

---

## 🎓 Resumen

**ISP en EduGo API Mobile**:
- ✅ 7/7 repositorios segregados
- ✅ Promedio 2-3 métodos por interfaz
- ✅ Interfaces Reader/Writer/Stats claramente definidas
- ✅ 95%+ cumplimiento del principio

**Beneficios logrados**:
- Dependencias mínimas en services
- Tests más simples (mocks pequeños)
- Código más mantenible
- Cumplimiento de SOLID

**Usa esta guía** para elegir la interfaz correcta al crear nuevos services. 🚀
