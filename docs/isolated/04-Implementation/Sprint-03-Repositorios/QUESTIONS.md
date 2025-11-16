# Preguntas y Decisiones del Sprint 03 - Repositorios

## Q001: ¿Usar GORM o SQL puro para PostgreSQL?
**Contexto:** Podemos usar ORM (GORM) o escribir SQL manualmente.

**Decisión por Defecto:** **GORM**

**Justificación:**
- Reduce boilerplate (no escribir SELECT/INSERT manualmente)
- Migrations automáticas con AutoMigrate()
- Type-safe (menos errores de SQL)
- EduGo ya usa GORM en otros módulos

**Implementación:**
```go
// Con GORM
db.Where("material_id = ?", materialID).First(&assessment)

// vs SQL puro
row := db.QueryRow("SELECT * FROM assessment WHERE material_id = $1", materialID)
```

---

## Q002: ¿Transacciones explícitas o automáticas?
**Contexto:** Al guardar Attempt+Answers, necesitamos atomicidad.

**Decisión por Defecto:** **Transacciones Explícitas con tx.Transaction()**

**Justificación:**
- Control total de atomicidad
- Rollback automático si falla
- Más claro en código (explicit > implicit)

**Implementación:**
```go
db.Transaction(func(tx *gorm.DB) error {
    tx.Create(&attempt)
    for _, answer := range answers {
        tx.Create(answer)
    }
    return nil // Commit automático
})
```

---

## Q003: ¿Testcontainers o mocks para tests?
**Contexto:** Tests de repositorios pueden usar mocks o BD real.

**Decisión por Defecto:** **Testcontainers (BD real)**

**Justificación:**
- Tests más realistas (queries reales)
- Detecta problemas de SQL que mocks no detectan
- Verifica constraints de BD

---

## Q004: ¿Retornar nil o error cuando no encuentra registro?
**Contexto:** FindByID puede no encontrar el registro.

**Decisión por Defecto:** **Retornar (nil, nil) cuando no existe**

**Justificación:**
- "No encontrado" no es error, es caso válido
- Caller puede verificar `if entity == nil`
- Alineado con convención Go de repositorios

**Implementación:**
```go
func (r *Repository) FindByID(id uuid.UUID) (*Entity, error) {
    var entity Entity
    result := db.First(&entity, id)
    
    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return nil, nil // ✅ No es error
    }
    
    return &entity, result.Error
}
```

---

## Q005: ¿Proyección en MongoDB para ocultar correct_answer?
**Contexto:** MongoDB tiene campo `correct_answer` que NO debe exponerse.

**Decisión por Defecto:** **Proyección con bson.M al hacer FindOne**

**Implementación:**
```go
projection := bson.M{
    "questions.correct_answer": 0,  // Excluir
    "questions.feedback": 0,
}
collection.FindOne(ctx, filter).Decode(&assessment)
```

---

**Generado con:** Claude Code  
**Sprint:** 03/06
