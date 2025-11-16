# Contratos de API
# Sistema de Evaluaciones - EduGo

**Versión:** 1.0.0  
**Fecha:** 14 de Noviembre, 2025  
**Proyecto:** edugo-api-mobile - Sistema de Evaluaciones

---

## 1. ESPECIFICACIÓN OPENAPI 3.0

### 1.1 Información General

```yaml
openapi: 3.0.3
info:
  title: EduGo API Mobile - Sistema de Evaluaciones
  description: |
    API REST para gestión de evaluaciones automáticas en la plataforma EduGo.
    
    Permite a estudiantes:
    - Obtener cuestionarios de materiales educativos
    - Enviar respuestas y recibir calificación automática
    - Consultar historial de intentos
    
    **Seguridad:** Autenticación JWT requerida en todos los endpoints.
  version: 1.0.0
  contact:
    name: Equipo EduGo
    email: dev@edugo.com
  license:
    name: Privado
    
servers:
  - url: http://localhost:8080
    description: Desarrollo local
  - url: https://api-dev.edugo.com
    description: Desarrollo
  - url: https://api-qa.edugo.com
    description: QA
  - url: https://api.edugo.com
    description: Producción

tags:
  - name: Assessments
    description: Operaciones de evaluaciones
  - name: Attempts
    description: Intentos de evaluación
  - name: Results
    description: Resultados y estadísticas

security:
  - BearerAuth: []

components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
      description: |
        Token JWT obtenido del endpoint /v1/auth/login.
        Formato: `Authorization: Bearer <token>`
```

---

## 2. ENDPOINTS DETALLADOS

### 2.1 GET /v1/materials/:id/assessment

#### Descripción
Obtiene el cuestionario asociado a un material educativo específico.

**⚠️ CRÍTICO:** Las respuestas correctas NUNCA son incluidas en el response.

#### OpenAPI Spec
```yaml
/v1/materials/{materialId}/assessment:
  get:
    tags:
      - Assessments
    summary: Obtener cuestionario de un material
    description: |
      Retorna las preguntas de evaluación de un material sin exponer las respuestas correctas.
      
      **Validaciones:**
      - Material debe existir
      - Material debe tener processing_status = 'completed'
      - Assessment debe estar disponible en MongoDB
      
      **Seguridad:**
      - Respuestas correctas sanitizadas antes de enviar
      - Solo usuarios autenticados
    operationId: getAssessment
    parameters:
      - name: materialId
        in: path
        required: true
        description: UUID del material educativo
        schema:
          type: string
          format: uuid
        example: "01936d9a-7f8e-7e4c-9d3f-987654321cba"
    
    responses:
      '200':
        description: Cuestionario obtenido exitosamente
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/AssessmentResponse'
            examples:
              success:
                summary: Cuestionario de Pascal
                value:
                  assessment_id: "01936d9a-0001-7e4c-9d3f-111111111111"
                  material_id: "01936d9a-7f8e-7e4c-9d3f-987654321cba"
                  title: "Cuestionario: Introducción a Pascal"
                  total_questions: 5
                  estimated_time_minutes: 10
                  questions:
                    - id: "q1"
                      text: "¿Qué es un compilador?"
                      type: "multiple_choice"
                      options:
                        - id: "a"
                          text: "Un programa que traduce código fuente"
                        - id: "b"
                          text: "Un tipo de variable"
      
      '400':
        $ref: '#/components/responses/BadRequest'
      '401':
        $ref: '#/components/responses/Unauthorized'
      '404':
        description: Material o assessment no encontrado
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ErrorResponse'
            examples:
              materialNotFound:
                value:
                  error: "not_found"
                  message: "Material not found"
              assessmentNotAvailable:
                value:
                  error: "not_found"
                  message: "Assessment not available for this material"
      '500':
        $ref: '#/components/responses/InternalServerError'
```

#### Request Example
```bash
curl -X GET \
  'http://localhost:8080/v1/materials/01936d9a-7f8e-7e4c-9d3f-987654321cba/assessment' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...'
```

#### Response Example
```json
{
  "assessment_id": "01936d9a-0001-7e4c-9d3f-111111111111",
  "material_id": "01936d9a-7f8e-7e4c-9d3f-987654321cba",
  "title": "Cuestionario: Introducción a Pascal",
  "total_questions": 5,
  "estimated_time_minutes": 10,
  "questions": [
    {
      "id": "q1",
      "text": "¿Qué es un compilador?",
      "type": "multiple_choice",
      "options": [
        {
          "id": "a",
          "text": "Un programa que traduce código fuente a código máquina"
        },
        {
          "id": "b",
          "text": "Un tipo de variable en Pascal"
        },
        {
          "id": "c",
          "text": "Una estructura de control"
        },
        {
          "id": "d",
          "text": "Un editor de texto"
        }
      ]
    }
  ]
}
```

---

### 2.2 POST /v1/materials/:id/assessment/attempts

#### Descripción
Crea un nuevo intento de evaluación enviando las respuestas del estudiante. El sistema valida, califica y retorna resultados inmediatamente.

#### OpenAPI Spec
```yaml
/v1/materials/{materialId}/assessment/attempts:
  post:
    tags:
      - Attempts
    summary: Crear intento de evaluación
    description: |
      Envía respuestas de un cuestionario, calcula puntaje automáticamente y retorna feedback educativo.
      
      **Proceso:**
      1. Validar que todas las preguntas tienen respuesta
      2. Obtener respuestas correctas de MongoDB
      3. Calcular puntaje
      4. Generar feedback educativo
      5. Persistir intento y respuestas en PostgreSQL (transacción ACID)
      6. Retornar resultados
      
      **Performance:** <2 segundos (p95)
    operationId: createAttempt
    parameters:
      - name: materialId
        in: path
        required: true
        schema:
          type: string
          format: uuid
    
    requestBody:
      required: true
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/CreateAttemptRequest'
          examples:
            validAttempt:
              summary: Intento válido con 5 respuestas
              value:
                answers:
                  - question_id: "q1"
                    selected_option: "a"
                  - question_id: "q2"
                    selected_option: "b"
                  - question_id: "q3"
                    selected_option: "c"
                  - question_id: "q4"
                    selected_option: "d"
                  - question_id: "q5"
                    selected_option: "a"
                time_spent_seconds: 420
    
    responses:
      '201':
        description: Intento creado exitosamente
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/AttemptResultResponse'
            examples:
              passed:
                summary: Intento aprobado (80%)
                value:
                  attempt_id: "01936d9b-1234-7e4c-9d3f-abcdef123456"
                  score: 80
                  max_score: 100
                  correct_answers: 4
                  total_questions: 5
                  pass_threshold: 70
                  passed: true
                  feedback:
                    - question_id: "q1"
                      question_text: "¿Qué es un compilador?"
                      selected_option: "a"
                      correct_answer: "a"
                      is_correct: true
                      message: "¡Correcto! Un compilador traduce código fuente..."
                  can_retake: true
                  previous_best_score: null
      
      '400':
        description: Request inválido
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ErrorResponse'
            examples:
              incompleteAnswers:
                value:
                  error: "validation_error"
                  message: "incomplete answers: expected 5, got 3"
              invalidQuestionId:
                value:
                  error: "validation_error"
                  message: "invalid question_id: q99"
      
      '401':
        $ref: '#/components/responses/Unauthorized'
      '404':
        $ref: '#/components/responses/NotFound'
      '500':
        $ref: '#/components/responses/InternalServerError'
```

#### Request Example
```bash
curl -X POST \
  'http://localhost:8080/v1/materials/01936d9a-7f8e-7e4c-9d3f-987654321cba/assessment/attempts' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...' \
  -H 'Content-Type: application/json' \
  -d '{
    "answers": [
      {"question_id": "q1", "selected_option": "a"},
      {"question_id": "q2", "selected_option": "b"},
      {"question_id": "q3", "selected_option": "c"},
      {"question_id": "q4", "selected_option": "d"},
      {"question_id": "q5", "selected_option": "a"}
    ],
    "time_spent_seconds": 420
  }'
```

#### Response Example
```json
{
  "attempt_id": "01936d9b-1234-7e4c-9d3f-abcdef123456",
  "score": 80,
  "max_score": 100,
  "correct_answers": 4,
  "total_questions": 5,
  "pass_threshold": 70,
  "passed": true,
  "feedback": [
    {
      "question_id": "q1",
      "question_text": "¿Qué es un compilador?",
      "selected_option": "a",
      "correct_answer": "a",
      "is_correct": true,
      "message": "¡Correcto! Un compilador traduce código fuente a código máquina ejecutable."
    },
    {
      "question_id": "q2",
      "question_text": "¿Cuál es la principal ventaja de la tipificación fuerte?",
      "selected_option": "c",
      "correct_answer": "b",
      "is_correct": false,
      "message": "Incorrecto. Piensa en qué sucede durante la compilación con tipos estrictos."
    }
  ],
  "can_retake": true,
  "previous_best_score": null
}
```

---

### 2.3 GET /v1/attempts/:id/results

#### OpenAPI Spec
```yaml
/v1/attempts/{attemptId}/results:
  get:
    tags:
      - Results
    summary: Obtener resultados de un intento
    description: |
      Retorna los resultados detallados de un intento específico.
      
      **Autorización:**
      - Solo el propietario del intento puede acceder
      - Retorna 403 si el usuario no es el propietario
    operationId: getAttemptResults
    parameters:
      - name: attemptId
        in: path
        required: true
        schema:
          type: string
          format: uuid
    
    responses:
      '200':
        description: Resultados obtenidos exitosamente
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/AttemptResultResponse'
      '401':
        $ref: '#/components/responses/Unauthorized'
      '403':
        description: Acceso denegado
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ErrorResponse'
            example:
              error: "forbidden"
              message: "You don't have permission to access this attempt"
      '404':
        $ref: '#/components/responses/NotFound'
```

---

### 2.4 GET /v1/users/me/attempts

#### OpenAPI Spec
```yaml
/v1/users/me/attempts:
  get:
    tags:
      - Attempts
    summary: Obtener historial de intentos del usuario
    description: |
      Retorna todos los intentos de evaluación del usuario autenticado, ordenados por fecha descendente.
      
      **Paginación:** Soporta limit y offset
    operationId: getUserAttemptHistory
    parameters:
      - name: limit
        in: query
        description: Número máximo de resultados
        schema:
          type: integer
          minimum: 1
          maximum: 100
          default: 10
      - name: offset
        in: query
        description: Número de resultados a saltar
        schema:
          type: integer
          minimum: 0
          default: 0
    
    responses:
      '200':
        description: Historial obtenido exitosamente
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/AttemptHistoryResponse'
            example:
              attempts:
                - attempt_id: "01936d9b-0003-7e4c-9d3f-333333333333"
                  material_id: "01936d9a-mat2-7e4c-9d3f-987654321xyz"
                  material_title: "Estructuras de Datos"
                  score: 90
                  max_score: 100
                  passed: true
                  completed_at: "2025-11-14T14:30:00Z"
                - attempt_id: "01936d9b-0001-7e4c-9d3f-111111111111"
                  material_id: "01936d9a-mat1-7e4c-9d3f-987654321cba"
                  material_title: "Introducción a Pascal"
                  score: 80
                  max_score: 100
                  passed: true
                  completed_at: "2025-11-14T10:07:00Z"
              total_count: 2
              page: 1
              limit: 10
```

---

## 3. SCHEMAS (MODELOS DE DATOS)

### 3.1 Request Schemas

```yaml
components:
  schemas:
    CreateAttemptRequest:
      type: object
      required:
        - answers
        - time_spent_seconds
      properties:
        answers:
          type: array
          minItems: 1
          maxItems: 100
          items:
            $ref: '#/components/schemas/AnswerInput'
        time_spent_seconds:
          type: integer
          minimum: 1
          maximum: 7200
          description: Tiempo total del intento en segundos (máx 2 horas)
      example:
        answers:
          - question_id: "q1"
            selected_option: "a"
        time_spent_seconds: 420
    
    AnswerInput:
      type: object
      required:
        - question_id
        - selected_option
      properties:
        question_id:
          type: string
          minLength: 1
          maxLength: 100
          description: ID de la pregunta
          example: "q1"
        selected_option:
          type: string
          minLength: 1
          maxLength: 10
          description: Opción seleccionada
          example: "a"
```

### 3.2 Response Schemas

```yaml
components:
  schemas:
    AssessmentResponse:
      type: object
      required:
        - assessment_id
        - material_id
        - title
        - total_questions
        - estimated_time_minutes
        - questions
      properties:
        assessment_id:
          type: string
          format: uuid
        material_id:
          type: string
          format: uuid
        title:
          type: string
          maxLength: 255
        total_questions:
          type: integer
          minimum: 1
        estimated_time_minutes:
          type: integer
          minimum: 1
        questions:
          type: array
          items:
            $ref: '#/components/schemas/QuestionDTO'
    
    QuestionDTO:
      type: object
      required:
        - id
        - text
        - type
        - options
      properties:
        id:
          type: string
          description: ID de la pregunta
          example: "q1"
        text:
          type: string
          minLength: 10
          maxLength: 500
          example: "¿Qué es un compilador?"
        type:
          type: string
          enum: [multiple_choice, true_false, short_answer]
        options:
          type: array
          items:
            $ref: '#/components/schemas/OptionDTO'
      description: |
        ⚠️ IMPORTANTE: Este DTO NUNCA incluye:
        - correct_answer
        - feedback
        
        Estos campos se sanitizan antes de enviar al cliente.
    
    OptionDTO:
      type: object
      required:
        - id
        - text
      properties:
        id:
          type: string
          example: "a"
        text:
          type: string
          example: "Un programa que traduce código fuente"
    
    AttemptResultResponse:
      type: object
      required:
        - attempt_id
        - score
        - max_score
        - correct_answers
        - total_questions
        - pass_threshold
        - passed
        - feedback
        - can_retake
      properties:
        attempt_id:
          type: string
          format: uuid
        score:
          type: integer
          minimum: 0
          maximum: 100
        max_score:
          type: integer
          default: 100
        correct_answers:
          type: integer
          minimum: 0
        total_questions:
          type: integer
          minimum: 1
        pass_threshold:
          type: integer
          minimum: 0
          maximum: 100
        passed:
          type: boolean
        feedback:
          type: array
          items:
            $ref: '#/components/schemas/FeedbackDTO'
        can_retake:
          type: boolean
        previous_best_score:
          type: integer
          nullable: true
          minimum: 0
          maximum: 100
    
    FeedbackDTO:
      type: object
      required:
        - question_id
        - question_text
        - selected_option
        - correct_answer
        - is_correct
        - message
      properties:
        question_id:
          type: string
        question_text:
          type: string
        selected_option:
          type: string
        correct_answer:
          type: string
        is_correct:
          type: boolean
        message:
          type: string
          description: Feedback educativo personalizado
    
    AttemptHistoryResponse:
      type: object
      required:
        - attempts
        - total_count
        - page
        - limit
      properties:
        attempts:
          type: array
          items:
            $ref: '#/components/schemas/AttemptSummaryDTO'
        total_count:
          type: integer
          minimum: 0
        page:
          type: integer
          minimum: 1
        limit:
          type: integer
          minimum: 1
          maximum: 100
    
    AttemptSummaryDTO:
      type: object
      required:
        - attempt_id
        - material_id
        - material_title
        - score
        - max_score
        - passed
        - completed_at
      properties:
        attempt_id:
          type: string
          format: uuid
        material_id:
          type: string
          format: uuid
        material_title:
          type: string
        score:
          type: integer
        max_score:
          type: integer
        passed:
          type: boolean
        completed_at:
          type: string
          format: date-time
```

### 3.3 Error Schemas

```yaml
components:
  schemas:
    ErrorResponse:
      type: object
      required:
        - error
        - message
      properties:
        error:
          type: string
          description: Código de error
          enum:
            - validation_error
            - not_found
            - unauthorized
            - forbidden
            - internal_server_error
        message:
          type: string
          description: Mensaje descriptivo del error
        details:
          type: object
          description: Detalles adicionales del error (opcional)
          additionalProperties: true
      example:
        error: "validation_error"
        message: "incomplete answers: expected 5, got 3"
        details:
          expected: 5
          received: 3
```

---

## 4. RESPONSES REUTILIZABLES

```yaml
components:
  responses:
    BadRequest:
      description: Request inválido
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ErrorResponse'
          example:
            error: "validation_error"
            message: "Invalid request parameters"
    
    Unauthorized:
      description: No autenticado
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ErrorResponse'
          example:
            error: "unauthorized"
            message: "Missing or invalid authentication token"
    
    Forbidden:
      description: Sin permisos
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ErrorResponse'
          example:
            error: "forbidden"
            message: "You don't have permission to access this resource"
    
    NotFound:
      description: Recurso no encontrado
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ErrorResponse'
          example:
            error: "not_found"
            message: "Resource not found"
    
    InternalServerError:
      description: Error interno del servidor
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/ErrorResponse'
          example:
            error: "internal_server_error"
            message: "An unexpected error occurred"
```

---

## 5. CÓDIGOS HTTP Y SIGNIFICADO

| Código | Nombre | Uso en API | Ejemplo |
|--------|--------|------------|---------|
| **200** | OK | Operación exitosa (GET) | Obtener assessment |
| **201** | Created | Recurso creado (POST) | Crear intento |
| **400** | Bad Request | Validación fallida | Respuestas incompletas |
| **401** | Unauthorized | Token faltante o inválido | Sin header Authorization |
| **403** | Forbidden | Sin permisos | Acceder a intento ajeno |
| **404** | Not Found | Recurso no existe | Material sin assessment |
| **500** | Internal Server Error | Error del servidor | DB caída |

---

## 6. HEADERS REQUERIDOS

### 6.1 Request Headers

```yaml
headers:
  Authorization:
    description: Token JWT de autenticación
    required: true
    schema:
      type: string
      pattern: "^Bearer [A-Za-z0-9-_=]+\.[A-Za-z0-9-_=]+\.?[A-Za-z0-9-_.+/=]*$"
    example: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  
  Content-Type:
    description: Tipo de contenido (solo para POST/PUT)
    required: true
    schema:
      type: string
      enum: [application/json]
    example: "application/json"
  
  Idempotency-Key:
    description: Clave de idempotencia para POST (Post-MVP)
    required: false
    schema:
      type: string
      format: uuid
    example: "01936d9b-1234-7e4c-9d3f-idempotency123"
```

### 6.2 Response Headers

```yaml
headers:
  Content-Type:
    schema:
      type: string
      enum: [application/json]
  
  X-Request-ID:
    description: ID único de la petición (para tracing)
    schema:
      type: string
      format: uuid
    example: "01936d9c-req1-7e4c-9d3f-requestid1234"
  
  X-RateLimit-Limit:
    description: Límite de requests por minuto
    schema:
      type: integer
    example: 100
  
  X-RateLimit-Remaining:
    description: Requests restantes en la ventana actual
    schema:
      type: integer
    example: 95
  
  X-RateLimit-Reset:
    description: Timestamp cuando se resetea el contador
    schema:
      type: integer
      format: unix-timestamp
    example: 1699977600
```

---

## 7. VERSIONADO DE API

### 7.1 Estrategia de Versionado

**Método:** URL Path Versioning

```
https://api.edugo.com/v1/materials/:id/assessment
                      ^^^ versión
```

**Razones:**
- ✅ Simple y explícito
- ✅ Fácil de cachear
- ✅ Compatible con Swagger/OpenAPI

**Alternativas rechazadas:**
- ❌ Header versioning (`Accept: application/vnd.edugo.v1+json`)
- ❌ Query parameter (`?version=1`)

### 7.2 Política de Deprecación

1. Anunciar deprecación con 6 meses de anticipación
2. Agregar header `Deprecation: true` en responses
3. Mantener versión antigua funcionando mínimo 1 año
4. Documentar migración en changelog

---

## 8. RATE LIMITING

### 8.1 Límites por Rol

| Rol | Requests/Minuto | Requests/Hora |
|-----|-----------------|---------------|
| **Student** | 100 | 3,000 |
| **Teacher** | 200 | 6,000 |
| **Admin** | 500 | 15,000 |

### 8.2 Response cuando se excede límite

```http
HTTP/1.1 429 Too Many Requests
Content-Type: application/json
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1699977660
Retry-After: 60

{
  "error": "rate_limit_exceeded",
  "message": "You have exceeded the rate limit. Try again in 60 seconds.",
  "retry_after": 60
}
```

---

## 9. EJEMPLOS DE USO COMPLETOS

### 9.1 Flujo Completo: Estudiante Realiza Quiz

```bash
# Paso 1: Obtener cuestionario
curl -X GET \
  'http://localhost:8080/v1/materials/01936d9a-mat1-7e4c-9d3f-987654321cba/assessment' \
  -H 'Authorization: Bearer <token>' \
  | jq .

# Response:
# {
#   "assessment_id": "...",
#   "questions": [...]
# }

# Paso 2: Estudiante responde en la app (5-10 minutos)

# Paso 3: Enviar respuestas
curl -X POST \
  'http://localhost:8080/v1/materials/01936d9a-mat1-7e4c-9d3f-987654321cba/assessment/attempts' \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "answers": [
      {"question_id": "q1", "selected_option": "a"},
      {"question_id": "q2", "selected_option": "b"},
      {"question_id": "q3", "selected_option": "c"},
      {"question_id": "q4", "selected_option": "d"},
      {"question_id": "q5", "selected_option": "a"}
    ],
    "time_spent_seconds": 420
  }' \
  | jq .

# Response:
# {
#   "attempt_id": "...",
#   "score": 80,
#   "passed": true,
#   "feedback": [...]
# }

# Paso 4: Consultar historial
curl -X GET \
  'http://localhost:8080/v1/users/me/attempts?limit=10' \
  -H 'Authorization: Bearer <token>' \
  | jq .
```

---

## 10. GENERACIÓN DE CÓDIGO CLIENTE

### 10.1 Usando OpenAPI Generator

```bash
# Generar cliente TypeScript
openapi-generator-cli generate \
  -i openapi.yaml \
  -g typescript-axios \
  -o ./generated/typescript-client

# Generar cliente Swift (iOS)
openapi-generator-cli generate \
  -i openapi.yaml \
  -g swift5 \
  -o ./generated/swift-client

# Generar cliente Kotlin (Android)
openapi-generator-cli generate \
  -i openapi.yaml \
  -g kotlin \
  -o ./generated/kotlin-client
```

### 10.2 Ejemplo de Uso del Cliente Generado

```typescript
// TypeScript Client
import { AssessmentsApi, Configuration } from './generated/typescript-client';

const config = new Configuration({
  basePath: 'https://api.edugo.com',
  accessToken: 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...'
});

const api = new AssessmentsApi(config);

// Obtener assessment
const assessment = await api.getAssessment('01936d9a-mat1-7e4c-9d3f-987654321cba');

// Crear intento
const result = await api.createAttempt('01936d9a-mat1-7e4c-9d3f-987654321cba', {
  answers: [
    { question_id: 'q1', selected_option: 'a' }
  ],
  time_spent_seconds: 420
});

console.log(`Score: ${result.score}/${result.max_score}`);
```

---

## 11. TESTING DE CONTRATOS

### 11.1 Validación con Spectral

```yaml
# .spectral.yaml
extends: ["spectral:oas"]
rules:
  oas3-api-servers: error
  operation-operationId-unique: error
  operation-description: warn
  operation-tags: error
  
  # Custom rules
  no-$ref-siblings: error
  paths-kebab-case: warn
```

```bash
# Validar spec
spectral lint openapi.yaml
```

### 11.2 Contract Testing con Pact

```go
// contract_test.go
func TestAssessmentContract(t *testing.T) {
    pact := &dsl.Pact{
        Consumer: "mobile-app",
        Provider: "api-mobile",
    }
    
    pact.AddInteraction().
        Given("Material uuid-1 has assessment").
        UponReceiving("GET /v1/materials/uuid-1/assessment").
        WithRequest(dsl.Request{
            Method: "GET",
            Path:   "/v1/materials/uuid-1/assessment",
            Headers: dsl.MapMatcher{
                "Authorization": "Bearer token",
            },
        }).
        WillRespondWith(dsl.Response{
            Status: 200,
            Body:   assessmentMatcher,
        })
}
```

---

**Generado con:** Claude Code  
**Total Endpoints:** 4  
**Spec Version:** OpenAPI 3.0.3  
**Última actualización:** 2025-11-14
