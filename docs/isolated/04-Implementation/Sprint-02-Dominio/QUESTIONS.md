# Preguntas y Decisiones del Sprint 02 - Capa de Dominio

## Q001: ¿Usar pointers o valores para entities?
**Contexto:** En Go, podemos pasar structs por valor o por puntero. Para entities grandes como Assessment y Attempt, necesitamos decidir el approach.

**Opciones:**

### 1. **Opción A: Usar Pointers (*Assessment, *Attempt)** - RECOMENDADO
- **Pros:**
  - Evita copias innecesarias de memoria (entities pueden ser grandes)
  - Permite métodos mutadores (SetMaxAttempts, SetTimeLimit)
  - Permite nil checks (verificar si entity existe)
  - Más idiomático en Go para tipos complejos
  - Facilita implementación de repositorios (retornar nil cuando no existe)
  
- **Contras:**
  - Posible nil pointer dereference si no se valida
  - Más cuidado con mutabilidad no deseada
  - Necesita documentar claramente qué métodos mutan y cuáles no

### 2. **Opción B: Usar Valores (Assessment, Attempt)**
- **Pros:**
  - No hay nil pointer dereference
  - Inmutabilidad por defecto (copias automáticas)
  - Más seguro en concurrencia
  
- **Contras:**
  - Copias costosas en memoria para structs grandes
  - No se puede representar "entity no existe" sin usar Optional pattern
  - Menos idiomático en Go para tipos de dominio
  - Complicado con métodos mutadores (requiere retornar nueva entity)

**Decisión por Defecto:** **Opción A - Usar Pointers**

**Justificación:** 
- Es el patrón estándar en Go para entities y agregados en DDD
- Los repositorios Go retornan (*Entity, error) convencionalmente
- Permite mutaciones controladas donde tiene sentido (ej: SetMaxAttempts)
- Para entities inmutables (Attempt), simplemente no exponemos setters

**Implementación:**
```go
// Constructores retornan punteros
func NewAssessment(...) (*Assessment, error) {
    return &Assessment{...}, nil
}

// Repositorios trabajan con punteros
type AssessmentRepository interface {
    FindByID(ctx context.Context, id uuid.UUID) (*Assessment, error)
    Save(ctx context.Context, assessment *Assessment) error
}

// Métodos en punteros receivers
func (a *Assessment) SetMaxAttempts(max int) error {
    // Muta el receiver
}

// Métodos de consulta en value o pointer receivers (ambos OK)
func (a Assessment) CanAttempt(count int) bool {
    // Solo lectura, no muta
}
```

---

## Q002: ¿Dónde colocar business rules: en entities o en services?
**Contexto:** Reglas como "calcular score", "validar intentos permitidos", "verificar si aprobó" - ¿dónde van?

**Opciones:**

### 1. **Opción A: Business Rules en Entities (Rich Domain Model)** - RECOMENDADO
- **Pros:**
  - Alineado con Domain-Driven Design (DDD)
  - Entities encapsulan su propia lógica
  - Más fácil de testear (test unitario de entity)
  - Lógica de negocio cerca de los datos
  - Previene "Anemic Domain Model" anti-pattern
  
- **Contras:**
  - Entities pueden volverse grandes
  - Algunas reglas requieren múltiples entities (van a services)

### 2. **Opción B: Business Rules en Services (Anemic Model)**
- **Pros:**
  - Entities más simples (solo getters/setters)
  - Lógica centralizada en services
  
- **Contras:**
  - Anti-pattern "Anemic Domain Model"
  - Lógica dispersa, difícil de encontrar
  - Services se vuelven "god objects"
  - Más difícil de testear (requiere mock de dependencies)

**Decisión por Defecto:** **Opción A - Business Rules en Entities**

**Justificación:**
- Seguimos principios de DDD y Clean Architecture
- Entities son el corazón del dominio
- Lógica de validación y cálculo pertenece a la entity

**Implementación:**
```go
// ✅ CORRECTO: Lógica en Entity
type Assessment struct {
    // ...
}

func (a *Assessment) CanAttempt(attemptCount int) bool {
    if a.MaxAttempts == nil {
        return true // Business rule: nil = ilimitado
    }
    return attemptCount < *a.MaxAttempts
}

// ✅ CORRECTO: Cálculo de score en Entity
type Attempt struct {
    // ...
}

func NewAttempt(..., answers []*Answer) (*Attempt, error) {
    // Business rule: calcular score basado en respuestas correctas
    correctCount := 0
    for _, answer := range answers {
        if answer.IsCorrect {
            correctCount++
        }
    }
    score := (correctCount * 100) / len(answers)
    
    return &Attempt{Score: score, ...}, nil
}

// ❌ INCORRECTO: Lógica en Service
type AssessmentService struct {}

func (s *AssessmentService) CanStudentAttempt(assessment *Assessment, attemptCount int) bool {
    // Esta lógica debería estar en Assessment.CanAttempt()
}
```

**Reglas para decidir dónde va la lógica:**
- **En Entity:** Validaciones, cálculos que usan SOLO datos de esa entity
- **En Service:** Orquestación que requiere múltiples entities, llamadas a repos, transacciones

---

## Q003: ¿Usar time.Time o int64 (Unix timestamp) para fechas?
**Contexto:** Campos como CreatedAt, UpdatedAt, StartedAt, CompletedAt - ¿qué tipo usar?

**Opciones:**

### 1. **Opción A: time.Time** - RECOMENDADO
- **Pros:**
  - Tipo nativo de Go, más idiomático
  - Métodos útiles: After(), Before(), Sub(), Format()
  - Type-safe (no confundir segundos con milisegundos)
  - Soporte de timezone
  - Mejor legibilidad en código
  
- **Contras:**
  - Ligeramente más bytes en memoria que int64
  - Requiere parsing para JSON (pero automático con encoding/json)

### 2. **Opción B: int64 (Unix timestamp en segundos/milisegundos)**
- **Pros:**
  - Más compacto en memoria (8 bytes vs ~24 bytes de time.Time)
  - Más simple en JSON (solo número)
  - Fácil de comparar (< > ==)
  
- **Contras:**
  - No type-safe (fácil confundir segundos/milisegundos)
  - Sin soporte de timezone
  - Requiere conversión manual para operaciones
  - Menos legible: `1699999999` vs `2023-11-14T12:00:00Z`

**Decisión por Defecto:** **Opción A - time.Time**

**Justificación:**
- Go estándar usa time.Time en stdlib
- GORM soporta time.Time nativamente
- Mejor experiencia de desarrollo (autocomplete, type safety)
- Diferencia de memoria es insignificante comparada con beneficios

**Implementación:**
```go
type Assessment struct {
    CreatedAt time.Time
    UpdatedAt time.Time
}

func NewAssessment(...) (*Assessment, error) {
    now := time.Now().UTC() // Siempre usar UTC
    return &Assessment{
        CreatedAt: now,
        UpdatedAt: now,
    }, nil
}

// Uso
assessment.CreatedAt.Before(otherTime)
assessment.UpdatedAt.After(threshold)
duration := assessment.CompletedAt.Sub(assessment.StartedAt)
```

**Convención importante:** Siempre usar **UTC** para evitar problemas de timezone:
```go
now := time.Now().UTC()  // ✅ CORRECTO
now := time.Now()         // ⚠️ EVITAR (timezone local)
```

---

## Q004: ¿Validar en constructor o en método Validate() separado?
**Contexto:** Podemos validar al crear la entity (fail-fast) o tener un método Validate() separado.

**Opciones:**

### 1. **Opción A: Validar en Constructor (Fail-Fast)** - RECOMENDADO
- **Pros:**
  - Imposible crear entity inválida
  - Fail-fast: errores detectados inmediatamente
  - No necesitas recordar llamar Validate()
  - Menos código boilerplate
  
- **Contras:**
  - Constructor puede volverse largo si hay muchas validaciones
  - No se puede crear entity "draft" para modificar después

### 2. **Opción B: Método Validate() Separado**
- **Pros:**
  - Constructor más simple
  - Flexibilidad para crear entity en estado inválido temporalmente
  - Útil para ORMs que crean entities vacías
  
- **Contras:**
  - Fácil olvidar llamar Validate()
  - Posibles entities inválidas en el sistema
  - Más código (constructor + Validate)

**Decisión por Defecto:** **Opción A + Opción B (Híbrido)**

**Justificación:**
- **Constructor valida:** Previene creación de entities inválidas desde código
- **Método Validate() también existe:** Para casos donde ORM/deserialización crea entity y necesitamos validar después

**Implementación:**
```go
// Constructor con validaciones (fail-fast)
func NewAssessment(materialID uuid.UUID, title string, totalQuestions int, ...) (*Assessment, error) {
    // Validar ANTES de crear
    if materialID == uuid.Nil {
        return nil, ErrInvalidMaterialID
    }
    if title == "" {
        return nil, ErrEmptyTitle
    }
    if totalQuestions < 1 || totalQuestions > 100 {
        return nil, ErrInvalidTotalQuestions
    }
    
    // Solo crear si validaciones pasan
    return &Assessment{
        MaterialID:     materialID,
        Title:          title,
        TotalQuestions: totalQuestions,
        // ...
    }, nil
}

// Método Validate() para validar estado actual
// Útil cuando entity viene de BD o JSON
func (a *Assessment) Validate() error {
    if a.MaterialID == uuid.Nil {
        return ErrInvalidMaterialID
    }
    if a.Title == "" {
        return ErrEmptyTitle
    }
    // ... mismas validaciones que constructor
    return nil
}

// Uso
assessment, err := NewAssessment(...) // ✅ Validado en construcción
if err != nil {
    return err // No se creó entity inválida
}

// O cuando viene de BD
var assessment Assessment
db.First(&assessment)
if err := assessment.Validate(); err != nil { // ✅ Validar después de deserializar
    return err
}
```

---

## Q005: ¿Usar errores estándar (errors.New) o custom error types?
**Contexto:** Para errores de dominio como ErrInvalidScore, ErrEmptyTitle - ¿qué approach usar?

**Opciones:**

### 1. **Opción A: Custom Error Types (structs que implementan error)** 
- **Pros:**
  - Pueden contener contexto adicional (campos, valores)
  - Type assertions para manejar específicamente
  - Más información para debugging
  
- **Contras:**
  - Más código (definir struct para cada error)
  - Overkill para errores simples de validación
  - Complica error handling

### 2. **Opción B: Errores Estándar con errors.New() o fmt.Errorf()** - RECOMENDADO
- **Pros:**
  - Simple y directo
  - Suficiente para la mayoría de casos
  - Comparable con errors.Is()
  - Menos código boilerplate
  
- **Contras:**
  - No puede contener datos adicionales
  - Solo mensaje de error

**Decisión por Defecto:** **Opción B - Errores Estándar con Sentinels**

**Justificación:**
- Errores de dominio son usualmente simples (validaciones)
- errors.New() es idiomático en Go
- errors.Is() permite comparación sin type assertions
- Si necesitamos contexto, usamos fmt.Errorf() con wrapping

**Implementación:**
```go
// internal/domain/errors/errors.go
package errors

import "errors"

// Errores sentinel (variables globales)
var (
    ErrInvalidAssessmentID    = errors.New("domain: invalid assessment ID")
    ErrInvalidMaterialID      = errors.New("domain: invalid material ID")
    ErrEmptyTitle             = errors.New("domain: assessment title cannot be empty")
    ErrInvalidTotalQuestions  = errors.New("domain: total questions must be between 1 and 100")
    // ...
)

// Uso en entity
func NewAssessment(materialID uuid.UUID, ...) (*Assessment, error) {
    if materialID == uuid.Nil {
        return nil, ErrInvalidMaterialID // ✅ Retornar error sentinel
    }
    // ...
}

// Uso en caller (comparación con errors.Is)
assessment, err := NewAssessment(uuid.Nil, ...)
if err != nil {
    if errors.Is(err, domainErrors.ErrInvalidMaterialID) { // ✅ Type-safe comparison
        // Manejar específicamente
    }
    return err
}
```

**Convención de nombres:**
- Prefijo: `Err` (ej: `ErrInvalidScore`, no `InvalidScoreError`)
- Mensaje: Comenzar con `"domain:"` para indicar origen
- Descripción clara y accionable

**Cuándo usar custom error types:**
- Cuando necesitas campos adicionales (ej: validación de múltiples campos)
- Errores que requieren datos para handling (ej: retry count)

```go
// Solo si realmente necesitas contexto adicional
type ValidationError struct {
    Field   string
    Value   interface{}
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}
```

---

## Q006: ¿Entity Attempt debe ser mutable o inmutable?
**Contexto:** Attempt representa un intento completado. ¿Debe permitir modificaciones después de creado?

**Opciones:**

### 1. **Opción A: Inmutable (No Setters)** - RECOMENDADO
- **Pros:**
  - Refleja realidad: un intento completado no cambia
  - Previene modificación fraudulenta de scores
  - Mejor para auditoría (datos históricos no cambian)
  - Thread-safe por diseño
  - Alineado con event sourcing
  
- **Contras:**
  - No se puede "corregir" un intento si hay error
  - Requiere crear nuevo intento para cualquier cambio

### 2. **Opción B: Mutable (Con Setters)**
- **Pros:**
  - Flexibilidad para correcciones
  - Menos registros en BD (actualizar existente vs crear nuevo)
  
- **Contras:**
  - Riesgo de modificación fraudulenta
  - Pérdida de información histórica
  - Requiere audit log para rastrear cambios
  - Complica lógica de concurrencia

**Decisión por Defecto:** **Opción A - Inmutable**

**Justificación:**
- Cumple requisitos de auditoría educativa
- Previene fraude (modificar scores después)
- Alineado con patrón append-only de PostgreSQL
- Si hay error, crear nuevo intento (con nota de corrección)

**Implementación:**
```go
type Attempt struct {
    ID          uuid.UUID
    Score       int
    // ... otros campos
    // ❌ NO HAY setters
}

// ✅ Constructor crea intento COMPLETO
func NewAttempt(..., answers []*Answer) (*Attempt, error) {
    // Calcular score aquí (no después)
    score := calculateScore(answers)
    
    return &Attempt{
        Score:   score,
        Answers: answers,
        // ...
    }, nil
}

// ✅ Solo métodos de consulta (getters)
func (a *Attempt) GetScore() int {
    return a.Score
}

func (a *Attempt) IsPassed(threshold int) bool {
    return a.Score >= threshold
}

// ❌ NO EXPONER métodos como:
// func (a *Attempt) SetScore(score int) { ... }
// func (a *Attempt) UpdateScore(newScore int) { ... }
```

**Para correcciones:**
```go
// Si necesitas "corregir" un intento, crear uno nuevo
newAttempt := &Attempt{
    StudentID: oldAttempt.StudentID,
    Score:     correctedScore,
    // Agregar metadata de corrección
}
```

---

## Resumen de Decisiones

| Pregunta | Decisión | Implementación |
|----------|----------|----------------|
| Q001: Pointers vs Valores | **Pointers** | `*Assessment`, `*Attempt` |
| Q002: Business rules location | **En Entities** | `assessment.CanAttempt()` |
| Q003: Tipo de fechas | **time.Time** | `CreatedAt time.Time` |
| Q004: Validación | **Constructor + Validate()** | Ambos métodos |
| Q005: Tipos de errores | **Errores Estándar** | `errors.New()` con sentinels |
| Q006: Mutabilidad Attempt | **Inmutable** | Sin setters, solo getters |

**Todas estas decisiones están alineadas con:**
- ✅ Domain-Driven Design (DDD)
- ✅ Clean Architecture
- ✅ Go idioms y best practices
- ✅ Requisitos de auditoría del sistema educativo

---

**Generado con:** Claude Code  
**Sprint:** 02/06  
**Última actualización:** 2025-11-14
