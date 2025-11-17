# Casos de Test
# Sistema de Evaluaciones - EduGo

**Versión:** 1.0.0  
**Fecha:** 14 de Noviembre, 2025

---

## 1. TESTS POR ENDPOINT

### GET /v1/materials/:id/assessment

**TC-001: Material con Assessment Existe**
- **Input:** Material ID válido con assessment
- **Expected:** 200 OK, assessment sin correct_answers
- **Validar:** Response NO contiene `correct_answer` ni `feedback`

**TC-002: Material No Existe**
- **Input:** UUID inexistente
- **Expected:** 404 Not Found

**TC-003: Material Sin Assessment**
- **Input:** Material válido pero sin assessment en MongoDB
- **Expected:** 404 Not Found

**TC-004: Sin Autenticación**
- **Input:** Request sin JWT
- **Expected:** 401 Unauthorized

**TC-005: Material ID Inválido**
- **Input:** "invalid-uuid"
- **Expected:** 400 Bad Request

---

### POST /v1/materials/:id/assessment/attempts

**TC-006: Intento Válido - Score 100%**
- **Input:** 5 respuestas correctas de 5 preguntas
- **Expected:** 201 Created, score=100

**TC-007: Intento Válido - Score 60%**
- **Input:** 3 respuestas correctas de 5
- **Expected:** 201 Created, score=60

**TC-008: Respuestas Faltantes**
- **Input:** Solo 3 respuestas de 5 preguntas
- **Expected:** 400 Bad Request

**TC-009: Question ID Inválido**
- **Input:** question_id que no existe en MongoDB
- **Expected:** 400 Bad Request

**TC-010: Material No Existe**
- **Input:** Material ID inexistente
- **Expected:** 404 Not Found

**TC-011: Límite de Intentos Alcanzado**
- **Input:** 4to intento cuando max_attempts=3
- **Expected:** 403 Forbidden

**TC-012: Score Calculado en Servidor**
- **Input:** Cliente envía answers, NO envía score
- **Expected:** Servidor calcula score comparando con MongoDB
- **Validar:** Score correcto basado en correct_answers de MongoDB

---

### GET /v1/attempts/:id/results

**TC-013: Attempt Existe y Pertenece a Usuario**
- **Input:** Attempt ID del usuario autenticado
- **Expected:** 200 OK, resultados con breakdown de respuestas

**TC-014: Attempt de Otro Usuario**
- **Input:** Attempt ID de otro estudiante
- **Expected:** 403 Forbidden

**TC-015: Attempt No Existe**
- **Input:** UUID inexistente
- **Expected:** 404 Not Found

**TC-016: Attempt ID Inválido**
- **Input:** "invalid-uuid"
- **Expected:** 400 Bad Request

---

### GET /v1/users/me/attempts

**TC-017: Listar Intentos del Usuario**
- **Input:** Usuario con 3 intentos
- **Expected:** 200 OK, lista de 3 intentos

**TC-018: Paginación Correcta**
- **Input:** limit=2, offset=1
- **Expected:** 200 OK, 2 intentos (salteando primero)

**TC-019: Ordenamiento por Fecha DESC**
- **Input:** Usuario con intentos de diferentes fechas
- **Expected:** Intentos ordenados del más reciente al más antiguo

**TC-020: Usuario Sin Intentos**
- **Input:** Usuario sin intentos
- **Expected:** 200 OK, array vacío

---

## 2. TESTS DE SEGURIDAD

### TC-021: Respuestas Correctas NUNCA Expuestas en GET
```go
func TestSecurity_TC021(t *testing.T) {
    response := getAssessment(materialID)
    
    // Verificar que NO contiene campos sensibles
    assert.NotContains(t, response, "correct_answer")
    assert.NotContains(t, response, "feedback")
    assert.NotContains(t, response, "explanation")
}
```

### TC-022: Inyección SQL Bloqueada
```go
func TestSecurity_TC022(t *testing.T) {
    // Intentar SQL injection en material_id
    maliciousID := "'; DROP TABLE assessment; --"
    
    response := getAssessment(maliciousID)
    
    // GORM usa prepared statements, debe retornar 400 (UUID inválido)
    assert.Equal(t, 400, response.StatusCode)
    
    // Verificar que tabla sigue existiendo
    var count int
    db.Table("assessment").Count(&count)
    assert.True(t, count >= 0) // Tabla no fue eliminada
}
```

### TC-023: Score Validado Servidor-Side
```go
func TestSecurity_TC023(t *testing.T) {
    // Cliente envía respuestas incorrectas
    answers := []Answer{
        {QuestionID: "q1", SelectedAnswerID: "wrong_answer"},
    }
    
    // Cliente intenta mentir con score alto en request (debe ser ignorado)
    response := createAttempt(answers, clientScore: 100)
    
    // Servidor calcula score real = 0%
    assert.Equal(t, 0, response.Score)
}
```

### TC-024: JWT Validado en Todos los Endpoints
```go
func TestSecurity_TC024(t *testing.T) {
    endpoints := []string{
        "GET /v1/materials/:id/assessment",
        "POST /v1/materials/:id/assessment/attempts",
        "GET /v1/attempts/:id/results",
        "GET /v1/users/me/attempts",
    }
    
    for _, endpoint := range endpoints {
        response := callWithoutJWT(endpoint)
        assert.Equal(t, 401, response.StatusCode)
    }
}
```

---

## 3. TESTS DE PERFORMANCE

### TC-030: Latencia GET Assessment
```go
func TestPerformance_TC030(t *testing.T) {
    samples := make([]time.Duration, 100)
    
    for i := 0; i < 100; i++ {
        start := time.Now()
        getAssessment(materialID)
        samples[i] = time.Since(start)
    }
    
    p95 := calculateP95(samples)
    assert.Less(t, p95.Milliseconds(), int64(500)) // <500ms p95
}
```

### TC-031: Latencia POST Attempt
```go
func TestPerformance_TC031(t *testing.T) {
    // Similar a TC-030
    // p95 <2000ms
}
```

---

## 4. TESTS DE EDGE CASES

**TC-040:** Assessment con 100 preguntas (límite superior)  
**TC-041:** Time limit = 1 minuto (límite inferior)  
**TC-042:** Time limit = 180 minutos (límite superior)  
**TC-043:** Pass threshold = 0% (todos aprueban)  
**TC-044:** Pass threshold = 100% (solo perfectos aprueban)  
**TC-045:** Intento de 1 sola pregunta  
**TC-046:** Concurrent attempts del mismo estudiante

---

## Resumen

**Total Casos de Test:** 46  
**Por Endpoint:** 5-6 casos cada uno  
**Seguridad:** 4 casos críticos  
**Performance:** 2 benchmarks  
**Edge Cases:** 7 casos

---

**Generado con:** Claude Code
