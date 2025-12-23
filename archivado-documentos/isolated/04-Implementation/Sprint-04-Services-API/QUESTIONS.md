# Preguntas y Decisiones del Sprint 04

## Q001: ¿Usar DTOs o entities directamente en handlers?
**Decisión por Defecto:** **DTOs (Data Transfer Objects)**

**Justificación:**
- Separar representación externa (JSON) de dominio interno
- Control total sobre qué campos exponer
- Evitar exponer campos internos (ej: UpdatedAt en algunos casos)

**Implementación:**
```go
// DTO para response
type AssessmentResponse struct {
    ID             uuid.UUID `json:"assessment_id"`
    MaterialID     uuid.UUID `json:"material_id"`
    Title          string    `json:"title"`
    TotalQuestions int       `json:"total_questions"`
    Questions      []QuestionDTO `json:"questions"`
    // NO exponer: mongo_document_id, created_at, etc.
}

// Convertir entity a DTO
func toDTO(entity *entities.Assessment) AssessmentResponse {
    return AssessmentResponse{
        ID: entity.ID,
        // ...
    }
}
```

---

## Q002: ¿Validación con struct tags o manual?
**Decisión por Defecto:** **Struct tags con go-playground/validator**

**Implementación:**
```go
type CreateAttemptRequest struct {
    Answers []AnswerDTO `json:"answers" binding:"required,min=1,dive"`
}

// Gin valida automáticamente
if err := c.ShouldBindJSON(&req); err != nil {
    return err
}
```

---

## Q003: ¿Dónde validar score: cliente o servidor?
**Decisión por Defecto:** **SIEMPRE en servidor (CRÍTICO)**

**Justificación:**
- Seguridad: cliente puede mentir
- Score calculado comparando con MongoDB
- NUNCA aceptar score del cliente

---

**Sprint:** 04/06
