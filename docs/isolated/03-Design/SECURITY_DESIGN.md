# Diseño de Seguridad
# Sistema de Evaluaciones - EduGo

**Versión:** 1.0.0  
**Fecha:** 14 de Noviembre, 2025  
**Proyecto:** edugo-api-mobile - Sistema de Evaluaciones

---

## 1. MODELO DE AMENAZAS (STRIDE)

### 1.1 Análisis STRIDE por Componente

#### Componente: API Mobile (Endpoints de Evaluaciones)

| Amenaza | Descripción | Impacto | Probabilidad | Mitigación |
|---------|-------------|---------|--------------|------------|
| **S**poofing | Atacante se hace pasar por estudiante | 🔴 Alto | 🟡 Media | Autenticación JWT, validación de token |
| **T**ampering | Modificar respuestas correctas en tránsito | 🔴 Alto | 🟢 Baja | HTTPS, validación servidor-side |
| **R**epudiation | Estudiante niega haber realizado intento | 🟡 Medio | 🟡 Media | Logs inmutables, timestamps |
| **I**nformation Disclosure | Exponer respuestas correctas | 🔴 Crítico | 🔴 Alta | Sanitización, nunca enviar correct_answer |
| **D**enial of Service | Saturar API con requests | 🟡 Medio | 🟡 Media | Rate limiting, throttling |
| **E**levation of Privilege | Estudiante accede a intentos ajenos | 🔴 Alto | 🟢 Baja | Autorización, validar user_id |

---

### 1.2 Threat Modeling Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                        App Móvil (Cliente)                       │
│                    ⚠️ Zona No Confiable                          │
└────────────────────────┬────────────────────────────────────────┘
                         │ HTTPS/TLS 1.3
                         ↓
┌─────────────────────────────────────────────────────────────────┐
│                     API Mobile (Servidor)                        │
│                    ✅ Zona Confiable                             │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │ 1. Auth Middleware → Validar JWT                          │  │
│  │ 2. Rate Limiting → Prevenir DoS                           │  │
│  │ 3. Input Validation → Sanitizar inputs                    │  │
│  │ 4. Business Logic → NUNCA confiar en cliente              │  │
│  │ 5. Data Sanitization → Remover respuestas correctas       │  │
│  └───────────────────────────────────────────────────────────┘  │
└────────────────────────┬───────────────────┬────────────────────┘
                         │                   │
                         ↓                   ↓
            ┌─────────────────────┐ ┌────────────────────┐
            │  PostgreSQL         │ │  MongoDB           │
            │  ✅ Zona Segura     │ │  ✅ Zona Segura    │
            │  - Intentos         │ │  - Preguntas CON   │
            │  - Respuestas       │ │    respuestas      │
            └─────────────────────┘ └────────────────────┘
```

**Principios:**
1. **Trust Boundary:** Cliente NO confiable, servidor confiable
2. **Defense in Depth:** Múltiples capas de validación
3. **Least Privilege:** Cada usuario solo accede a sus datos
4. **Fail Secure:** En caso de error, denegar acceso

---

## 2. AUTENTICACIÓN Y AUTORIZACIÓN

### 2.1 Autenticación JWT

#### Token Structure
```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "sub": "01936d9a-stud-7e4c-9d3f-student11111",  // User ID
    "role": "student",
    "email": "juan.perez@edugo.com",
    "iat": 1699977600,  // Issued At
    "exp": 1699978500   // Expiration (15 min)
  },
  "signature": "HMACSHA256(...)"
}
```

#### Implementación
```go
// Middleware de autenticación
func AuthMiddleware(secretKey string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Extraer token del header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "missing authorization header"})
            c.Abort()
            return
        }
        
        // 2. Validar formato "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(401, gin.H{"error": "invalid authorization format"})
            c.Abort()
            return
        }
        
        tokenString := parts[1]
        
        // 3. Validar firma y expiración
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Verificar algoritmo
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, errors.New("invalid signing method")
            }
            return []byte(secretKey), nil
        })
        
        if err != nil || !token.Valid {
            c.JSON(401, gin.H{"error": "invalid or expired token"})
            c.Abort()
            return
        }
        
        // 4. Extraer claims
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(401, gin.H{"error": "invalid token claims"})
            c.Abort()
            return
        }
        
        // 5. Guardar en contexto
        c.Set("user_id", claims["sub"].(string))
        c.Set("role", claims["role"].(string))
        c.Set("email", claims["email"].(string))
        
        c.Next()
    }
}
```

#### Configuración de Tokens
```yaml
jwt:
  secret_key: "${JWT_SECRET_KEY}"  # Variable de entorno (256 bits)
  access_token_ttl: 15m             # Token de acceso
  refresh_token_ttl: 7d             # Token de refresco (Post-MVP)
  issuer: "edugo-api-mobile"
  audience: "edugo-app"
```

---

### 2.2 Autorización (RBAC - Role-Based Access Control)

#### Roles Definidos
```go
const (
    RoleStudent  = "student"
    RoleTeacher  = "teacher"
    RoleAdmin    = "admin"
    RoleGuardian = "guardian"
)
```

#### Matriz de Permisos

| Endpoint | Student | Teacher | Admin | Guardian |
|----------|---------|---------|-------|----------|
| GET /assessment | ✅ | ✅ | ✅ | ❌ |
| POST /attempts | ✅ | ✅ | ✅ | ❌ |
| GET /attempts/:id | ✅ (propio) | ✅ (sus estudiantes) | ✅ | ✅ (sus hijos) |
| GET /users/me/attempts | ✅ | ✅ | ✅ | ❌ |

#### Middleware de Autorización
```go
// RequireRole verifica que el usuario tenga uno de los roles permitidos
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole := c.GetString("role")
        
        for _, allowedRole := range allowedRoles {
            if userRole == allowedRole {
                c.Next()
                return
            }
        }
        
        c.JSON(403, gin.H{"error": "forbidden", "message": "insufficient permissions"})
        c.Abort()
    }
}

// Uso en routes
router.POST("/attempts/:id", 
    AuthMiddleware(secretKey),
    RequireRole(RoleStudent, RoleTeacher),
    handler.CreateAttempt)
```

#### Autorización a Nivel de Recurso
```go
// Solo el propietario puede acceder a sus intentos
func (h *AssessmentHandler) GetAttemptResults(c *gin.Context) {
    attemptID := c.Param("id")
    userID := c.GetString("user_id")
    
    // 1. Obtener intento
    attempt, err := h.service.GetAttempt(c.Request.Context(), attemptID)
    if err != nil {
        c.JSON(404, gin.H{"error": "not_found"})
        return
    }
    
    // 2. ⚠️ CRÍTICO: Verificar que el usuario es el propietario
    if attempt.StudentID != userID {
        c.JSON(403, gin.H{"error": "forbidden", "message": "access denied"})
        return
    }
    
    // 3. Retornar resultados
    c.JSON(200, attempt)
}
```

---

## 3. PROTECCIÓN CONTRA TRAMPAS

### 3.1 NUNCA Exponer Respuestas Correctas

#### ❌ INCORRECTO (INSEGURO)
```go
// ⚠️ NUNCA hacer esto
func GetAssessment(c *gin.Context) {
    questions, _ := questionRepo.FindAll()
    
    // ❌ Enviando respuestas correctas al cliente
    c.JSON(200, questions) // Incluye correct_answer
}
```

#### ✅ CORRECTO (SEGURO)
```go
// ✅ Sanitizar antes de enviar
func (s *AssessmentService) GetAssessmentForStudent(ctx context.Context, materialID string) (*AssessmentDTO, error) {
    // 1. Obtener preguntas de MongoDB (CON respuestas correctas)
    questions, err := s.questionRepo.FindByMaterialID(ctx, materialID)
    if err != nil {
        return nil, err
    }
    
    // 2. ⚠️ CRÍTICO: Sanitizar
    sanitizedQuestions := make([]QuestionDTO, len(questions))
    for i, q := range questions {
        sanitizedQuestions[i] = QuestionDTO{
            ID:      q.ID,
            Text:    q.Text,
            Type:    q.Type,
            Options: q.Options,
            // ❌ NO incluir: CorrectAnswer, Feedback
        }
    }
    
    return &AssessmentDTO{Questions: sanitizedQuestions}, nil
}
```

#### Test de Seguridad
```go
func TestAssessment_NeverExposeCorrectAnswers(t *testing.T) {
    // Setup
    handler := setupTestHandler()
    router := gin.Default()
    router.GET("/assessment/:id", handler.GetAssessment)
    
    // Request
    req := httptest.NewRequest("GET", "/assessment/uuid-1", nil)
    req.Header.Set("Authorization", "Bearer "+validToken)
    resp := httptest.NewRecorder()
    router.ServeHTTP(resp, req)
    
    // Assert
    var result map[string]interface{}
    json.Unmarshal(resp.Body.Bytes(), &result)
    
    questions := result["questions"].([]interface{})
    for _, q := range questions {
        question := q.(map[string]interface{})
        
        // ⚠️ CRÍTICO: Verificar que NO existen estos campos
        _, hasCorrectAnswer := question["correct_answer"]
        _, hasFeedback := question["feedback"]
        
        assert.False(t, hasCorrectAnswer, "correct_answer must not be exposed")
        assert.False(t, hasFeedback, "feedback must not be exposed")
    }
}
```

---

### 3.2 Validación SIEMPRE en Servidor

#### ❌ INCORRECTO (Cliente calcula puntaje)
```go
// ⚠️ NUNCA confiar en el cliente
type CreateAttemptRequest struct {
    Answers []Answer `json:"answers"`
    Score   int      `json:"score"` // ❌ Cliente puede mentir
}

func CreateAttempt(req CreateAttemptRequest) {
    // ❌ Guardar score del cliente directamente
    attempt.Score = req.Score
    db.Save(attempt)
}
```

#### ✅ CORRECTO (Servidor calcula puntaje)
```go
// ✅ Servidor calcula y valida
type CreateAttemptRequest struct {
    Answers          []Answer `json:"answers"`
    TimeSpentSeconds int      `json:"time_spent_seconds"`
    // ❌ NO incluir: Score (se calcula en servidor)
}

func (s *ScoringService) ScoreAttempt(ctx context.Context, answers []Answer) (int, error) {
    // 1. Obtener respuestas correctas de MongoDB (fuente de verdad)
    correctAnswers, err := s.questionRepo.FindCorrectAnswers(ctx, assessmentID)
    if err != nil {
        return 0, err
    }
    
    // 2. Calcular puntaje en servidor
    correctCount := 0
    for _, answer := range answers {
        correct := findCorrectAnswer(correctAnswers, answer.QuestionID)
        if answer.SelectedOption == correct {
            correctCount++
        }
    }
    
    score := (correctCount * 100) / len(answers)
    return score, nil
}
```

---

### 3.3 Detección de Intentos Sospechosos

```go
// Detectar patrones anormales
func (s *ScoringService) DetectSuspiciousAttempt(attempt *Attempt) []string {
    warnings := []string{}
    
    // 1. Tiempo sospechosamente corto (<5 seg por pregunta)
    minTime := attempt.TotalQuestions * 5
    if attempt.TimeSpentSeconds < minTime {
        warnings = append(warnings, "time_too_short")
        logger.Warn("Suspicious attempt: too fast",
            "attempt_id", attempt.ID,
            "time_spent", attempt.TimeSpentSeconds,
            "min_expected", minTime)
    }
    
    // 2. Puntaje perfecto en primer intento (puede ser legítimo, pero revisar)
    if attempt.Score == 100 && attempt.AttemptNumber == 1 {
        warnings = append(warnings, "perfect_first_attempt")
    }
    
    // 3. Patrones de respuesta (Post-MVP: ML para detectar)
    // Ej: Todas las respuestas son "a", "b", "c", "d", "a" (patrón)
    
    return warnings
}
```

---

## 4. VALIDACIÓN Y SANITIZACIÓN DE INPUTS

### 4.1 Validación de Request Body

```go
// Usar struct tags para validación automática
type CreateAttemptRequest struct {
    Answers []AnswerDTO `json:"answers" binding:"required,min=1,max=100,dive"`
    TimeSpentSeconds int `json:"time_spent_seconds" binding:"required,min=1,max=7200"`
}

type AnswerDTO struct {
    QuestionID     string `json:"question_id" binding:"required,min=1,max=100"`
    SelectedOption string `json:"selected_option" binding:"required,min=1,max=10"`
}

// Gin automáticamente valida
func (h *AssessmentHandler) CreateAttempt(c *gin.Context) {
    var req CreateAttemptRequest
    
    // Validación automática
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "validation_error", "message": err.Error()})
        return
    }
    
    // Continuar con lógica
}
```

### 4.2 Validación de UUIDs

```go
// Validar que IDs son UUIDs válidos
func ValidateUUID(id string) error {
    _, err := uuid.Parse(id)
    if err != nil {
        return errors.New("invalid UUID format")
    }
    return nil
}

// Uso
func (h *AssessmentHandler) GetAttemptResults(c *gin.Context) {
    attemptID := c.Param("id")
    
    // ⚠️ CRÍTICO: Validar antes de query
    if err := ValidateUUID(attemptID); err != nil {
        c.JSON(400, gin.H{"error": "validation_error", "message": "invalid attempt_id"})
        return
    }
    
    // Continuar...
}
```

### 4.3 Protección contra SQL Injection

```go
// ✅ CORRECTO: Usar parámetros preparados
func (r *AttemptRepository) FindByID(ctx context.Context, id uuid.UUID) (*Attempt, error) {
    var attempt Attempt
    
    // GORM previene SQL injection automáticamente
    err := r.db.WithContext(ctx).
        Where("id = ?", id).  // ✅ Parametrizado
        First(&attempt).Error
    
    return &attempt, err
}

// ❌ INCORRECTO: String concatenation (NUNCA hacer esto)
func FindByIDUnsafe(id string) (*Attempt, error) {
    query := "SELECT * FROM assessment_attempt WHERE id = '" + id + "'"  // ❌ VULNERABLE
    // Un atacante puede enviar: "'; DROP TABLE users; --"
}
```

---

## 5. PROTECCIÓN DE DATOS

### 5.1 Encriptación en Tránsito

#### HTTPS/TLS 1.3
```yaml
# nginx.conf
server {
    listen 443 ssl http2;
    server_name api.edugo.com;
    
    # Certificado SSL
    ssl_certificate /etc/letsencrypt/live/api.edugo.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.edugo.com/privkey.pem;
    
    # TLS 1.3 únicamente
    ssl_protocols TLSv1.3;
    ssl_ciphers 'TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384';
    
    # HSTS
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    
    location / {
        proxy_pass http://api-mobile:8080;
    }
}
```

### 5.2 Encriptación en Reposo

#### PostgreSQL
```sql
-- Encriptar campos sensibles con pgcrypto (Post-MVP si es necesario)
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Ejemplo: Encriptar comentarios privados (si se agregan)
INSERT INTO student_notes (student_id, note, created_at)
VALUES (
    $1,
    pgp_sym_encrypt($2, 'encryption-key'),  -- Encriptado
    NOW()
);

-- Desencriptar en query
SELECT 
    student_id,
    pgp_sym_decrypt(note::bytea, 'encryption-key') as note
FROM student_notes;
```

### 5.3 Manejo de Secretos

```bash
# .env (NUNCA commitear)
JWT_SECRET_KEY=super-secret-key-256-bits-long
DATABASE_URL=postgres://user:pass@localhost/edugo
MONGODB_URI=mongodb://localhost:27017/edugo

# Usar variables de entorno
export JWT_SECRET_KEY=$(openssl rand -base64 32)
```

```go
// Cargar desde variables de entorno
import "github.com/joho/godotenv"

func LoadConfig() (*Config, error) {
    godotenv.Load() // Cargar .env
    
    return &Config{
        JWTSecretKey: os.Getenv("JWT_SECRET_KEY"),
        DatabaseURL:  os.Getenv("DATABASE_URL"),
    }, nil
}
```

---

## 6. RATE LIMITING Y THROTTLING

### 6.1 Implementación con Gin

```go
import "github.com/ulule/limiter/v3"
import limitergin "github.com/ulule/limiter/v3/drivers/middleware/gin"
import "github.com/ulule/limiter/v3/drivers/store/memory"

func RateLimitMiddleware() gin.HandlerFunc {
    // Configurar rate limiter
    rate := limiter.Rate{
        Period: 1 * time.Minute,
        Limit:  100, // 100 req/min
    }
    
    store := memory.NewStore()
    instance := limiter.New(store, rate)
    
    return limitergin.NewMiddleware(instance)
}

// Aplicar globalmente
router := gin.Default()
router.Use(RateLimitMiddleware())
```

### 6.2 Rate Limiting por Rol

```go
func RoleBasedRateLimiter() gin.HandlerFunc {
    return func(c *gin.Context) {
        role := c.GetString("role")
        
        limits := map[string]int{
            "student": 100,  // 100 req/min
            "teacher": 200,
            "admin":   500,
        }
        
        limit, exists := limits[role]
        if !exists {
            limit = 50 // Default
        }
        
        // Verificar límite (implementación simplificada)
        key := fmt.Sprintf("ratelimit:%s:%s", role, c.GetString("user_id"))
        count := incrementCounter(key, 60) // 60 segundos
        
        if count > limit {
            c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
            c.Header("X-RateLimit-Remaining", "0")
            c.JSON(429, gin.H{
                "error": "rate_limit_exceeded",
                "message": "Too many requests. Try again later.",
                "retry_after": 60,
            })
            c.Abort()
            return
        }
        
        c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
        c.Header("X-RateLimit-Remaining", strconv.Itoa(limit-count))
        c.Next()
    }
}
```

---

## 7. LOGGING Y AUDITORÍA

### 7.1 Logs Estructurados

```go
import "github.com/edugogroup/edugo-shared/pkg/logger"

func (h *AssessmentHandler) CreateAttempt(c *gin.Context) {
    userID := c.GetString("user_id")
    materialID := c.Param("id")
    
    // Log inicio de operación
    logger.Info("Creating assessment attempt",
        "user_id", userID,
        "material_id", materialID,
        "ip", c.ClientIP(),
        "user_agent", c.Request.UserAgent())
    
    // Business logic
    attempt, err := h.service.CreateAttempt(...)
    
    if err != nil {
        // Log error
        logger.Error("Failed to create attempt",
            "user_id", userID,
            "error", err.Error())
        c.JSON(500, gin.H{"error": "internal_server_error"})
        return
    }
    
    // Log éxito
    logger.Info("Assessment attempt created successfully",
        "user_id", userID,
        "attempt_id", attempt.ID,
        "score", attempt.Score)
    
    c.JSON(201, attempt)
}
```

### 7.2 Logs de Seguridad

```go
// Loguear eventos de seguridad
func LogSecurityEvent(eventType string, details map[string]interface{}) {
    entry := logger.With(
        "event_type", eventType,
        "timestamp", time.Now().UTC(),
    )
    
    for k, v := range details {
        entry = entry.With(k, v)
    }
    
    entry.Warn("Security event")
}

// Ejemplos de uso
LogSecurityEvent("auth_failed", map[string]interface{}{
    "user_id": userID,
    "ip": clientIP,
    "reason": "invalid_token",
})

LogSecurityEvent("unauthorized_access", map[string]interface{}{
    "user_id": userID,
    "attempted_resource": attemptID,
    "owner_id": actualOwnerID,
})

LogSecurityEvent("suspicious_attempt", map[string]interface{}{
    "attempt_id": attemptID,
    "user_id": userID,
    "time_spent": timeSpent,
    "reason": "time_too_short",
})
```

### 7.3 Auditoría de Operaciones Críticas

```sql
-- Tabla de auditoría (Post-MVP)
CREATE TABLE audit_log (
    id UUID PRIMARY KEY DEFAULT gen_uuid_v7(),
    event_type VARCHAR(100) NOT NULL,
    user_id UUID NOT NULL,
    resource_type VARCHAR(100) NOT NULL,
    resource_id UUID,
    action VARCHAR(50) NOT NULL,
    ip_address INET,
    user_agent TEXT,
    metadata JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_user_id ON audit_log(user_id);
CREATE INDEX idx_audit_created_at ON audit_log(created_at DESC);
CREATE INDEX idx_audit_event_type ON audit_log(event_type);
```

```go
// Registrar en auditoría
func AuditAttemptCreated(attempt *Attempt, userIP string) {
    db.Exec(`
        INSERT INTO audit_log (event_type, user_id, resource_type, resource_id, action, ip_address, metadata)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `,
        "assessment_attempt_created",
        attempt.StudentID,
        "assessment_attempt",
        attempt.ID,
        "CREATE",
        userIP,
        map[string]interface{}{
            "score": attempt.Score,
            "assessment_id": attempt.AssessmentID,
        })
}
```

---

## 8. HEADERS DE SEGURIDAD

### 8.1 Security Headers

```go
func SecurityHeadersMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Prevenir clickjacking
        c.Header("X-Frame-Options", "DENY")
        
        // Prevenir MIME sniffing
        c.Header("X-Content-Type-Options", "nosniff")
        
        // Habilitar XSS protection
        c.Header("X-XSS-Protection", "1; mode=block")
        
        // Content Security Policy
        c.Header("Content-Security-Policy", "default-src 'self'")
        
        // Referrer Policy
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        
        // Permissions Policy
        c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
        
        c.Next()
    }
}
```

### 8.2 CORS Configuration

```go
import "github.com/gin-contrib/cors"

func CORSMiddleware() gin.HandlerFunc {
    config := cors.Config{
        AllowOrigins:     []string{"https://app.edugo.com"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Authorization", "Content-Type", "X-Request-ID"},
        ExposeHeaders:    []string{"X-RateLimit-Limit", "X-RateLimit-Remaining"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }
    
    return cors.New(config)
}
```

---

## 9. COMPLIANCE Y CUMPLIMIENTO

### 9.1 OWASP Top 10 (2021)

| # | Vulnerabilidad | Mitigación en Sistema |
|---|----------------|----------------------|
| A01 | Broken Access Control | ✅ Autorización a nivel de recurso, validar user_id |
| A02 | Cryptographic Failures | ✅ TLS 1.3, HTTPS, JWT con HMAC |
| A03 | Injection | ✅ Parámetros preparados (GORM), validación de inputs |
| A04 | Insecure Design | ✅ Threat modeling (STRIDE), validación servidor-side |
| A05 | Security Misconfiguration | ✅ Security headers, configuración por ambiente |
| A06 | Vulnerable Components | ✅ Dependencias actualizadas, Dependabot |
| A07 | Authentication Failures | ✅ JWT, rate limiting, logs de intentos fallidos |
| A08 | Software/Data Integrity | ✅ Inmutabilidad de intentos, checksums |
| A09 | Logging Failures | ✅ Logs estructurados, logs de seguridad |
| A10 | SSRF | N/A (no hay requests salientes a URLs de usuario) |

### 9.2 GDPR / Protección de Datos Personales

```go
// Anonimización de datos (Post-MVP)
func AnonymizeStudent(studentID uuid.UUID) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // Anonimizar usuario
        tx.Model(&User{}).Where("id = ?", studentID).Updates(map[string]interface{}{
            "email":      fmt.Sprintf("deleted-%s@anonymous.com", studentID),
            "first_name": "Deleted",
            "last_name":  "User",
        })
        
        // Mantener intentos pero desvinculados
        // (o eliminar según política de retención)
        
        return nil
    })
}
```

---

## 10. CHECKLIST DE SEGURIDAD

### Pre-Deployment
- [ ] Todos los endpoints requieren autenticación
- [ ] Autorizaciones a nivel de recurso implementadas
- [ ] Respuestas correctas NUNCA expuestas al cliente
- [ ] Validación de inputs en todos los endpoints
- [ ] UUIDs validados antes de queries
- [ ] Rate limiting configurado
- [ ] HTTPS/TLS 1.3 habilitado
- [ ] Security headers configurados
- [ ] CORS configurado correctamente
- [ ] Secretos en variables de entorno (no en código)
- [ ] Logs de seguridad implementados
- [ ] Tests de seguridad pasando

### Post-Deployment
- [ ] Penetration testing ejecutado
- [ ] Dependency scanning (Dependabot)
- [ ] SSL Labs test A+ rating
- [ ] Logs de seguridad siendo monitoreados
- [ ] Alertas de intentos sospechosos configuradas
- [ ] Backups encriptados
- [ ] Plan de respuesta a incidentes documentado

---

**Generado con:** Claude Code  
**Metodología:** STRIDE Threat Modeling + OWASP Top 10  
**Última actualización:** 2025-11-14
