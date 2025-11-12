# Refactorización de main.go

## 📋 Resumen

Se ha refactorizado el archivo `cmd/main.go` para seguir los principios de **Clean Architecture** y mejorar la **separación de responsabilidades**. El código que antes tenía más de 200 líneas ahora tiene aproximadamente 135 líneas, con responsabilidades claras delegadas a módulos especializados.

## 🎯 Problemas Identificados (Antes)

El `main.go` original tenía las siguientes responsabilidades:

1. ❌ Inicialización de PostgreSQL (lógica de conexión)
2. ❌ Inicialización de MongoDB (lógica de conexión)
3. ❌ Configuración de middleware CORS
4. ❌ Configuración completa de rutas HTTP
5. ❌ Health check handler con lógica de negocio
6. ❌ Configuración de Swagger
7. ❌ Inicialización del servidor

**Total**: ~230 líneas con múltiples responsabilidades mezcladas.

## ✅ Solución Implementada

### Nueva Estructura de Archivos

```
internal/
├── infrastructure/
│   ├── database/
│   │   ├── postgres.go          # ✨ NUEVO: Inicialización PostgreSQL
│   │   └── mongodb.go            # ✨ NUEVO: Inicialización MongoDB
│   └── http/
│       ├── handler/
│       │   └── health_handler.go # ✨ NUEVO: Health check como handler
│       ├── middleware/
│       │   └── cors.go           # ✨ NUEVO: Middleware CORS reutilizable
│       └── router/
│           └── router.go         # ✨ NUEVO: Configuración de rutas
```

### Responsabilidades por Módulo

#### 1. `internal/infrastructure/database/postgres.go`

**Responsabilidad**: Inicialización y configuración de PostgreSQL.

```go
func InitPostgreSQL(ctx context.Context, cfg *config.Config, log logger.Logger) (*sql.DB, error)
```

**Características**:
- ✅ Manejo de pool de conexiones
- ✅ Verificación de conexión con timeout
- ✅ Logging estructurado con Zap
- ✅ Manejo de errores con contexto

---

#### 2. `internal/infrastructure/database/mongodb.go`

**Responsabilidad**: Inicialización y configuración de MongoDB.

```go
func InitMongoDB(ctx context.Context, cfg *config.Config, log logger.Logger) (*mongo.Database, error)
```

**Características**:
- ✅ Conexión con timeout configurable
- ✅ Ping de verificación
- ✅ Logging estructurado
- ✅ Retorna `*mongo.Database` listo para usar

---

#### 3. `internal/infrastructure/http/middleware/cors.go`

**Responsabilidad**: Configuración de CORS para peticiones cross-origin.

```go
func CORS() gin.HandlerFunc
```

**Características**:
- ✅ Middleware reutilizable
- ✅ Manejo de preflight requests (OPTIONS)
- ✅ Headers configurables
- ⚠️ **TODO**: En producción, restringir orígenes permitidos

---

#### 4. `internal/infrastructure/http/handler/health_handler.go`

**Responsabilidad**: Health check del sistema y sus dependencias.

```go
type HealthHandler struct {
    db      *sql.DB
    mongoDB *mongo.Database
}

func (h *HealthHandler) Check(c *gin.Context)
```

**Características**:
- ✅ Implementado como handler propio (no función suelta)
- ✅ Verifica estado de PostgreSQL y MongoDB
- ✅ Retorna estado general del sistema (`healthy`, `degraded`)
- ✅ Incluye timestamp y versión
- ✅ Anotaciones Swagger incluidas

**Respuesta**:
```json
{
  "status": "healthy",
  "service": "edugo-api-mobile",
  "version": "1.0.0",
  "postgres": "healthy",
  "mongodb": "healthy",
  "timestamp": "2025-01-15T10:30:00Z"
}
```

---

#### 5. `internal/infrastructure/http/router/router.go`

**Responsabilidad**: Configuración centralizada de todas las rutas HTTP.

```go
func SetupRouter(c *container.Container, healthHandler *handler.HealthHandler) *gin.Engine
```

**Organización**:
```
SetupRouter()
├── setupAuthPublicRoutes()      # POST /auth/login, /auth/refresh
├── setupProtectedRoutes()
    ├── setupAuthProtectedRoutes()  # POST /auth/logout, /auth/revoke-all
    └── setupMaterialRoutes()       # GET/POST /materials/...
```

**Ventajas**:
- ✅ Rutas organizadas por recurso
- ✅ Separación clara de rutas públicas vs protegidas
- ✅ Fácil de extender (agregar nuevo recurso = nueva función)
- ✅ Middleware JWT aplicado solo a rutas protegidas

---

#### 6. `cmd/main.go` (Refactorizado)

**Responsabilidad**: Orquestación del inicio de la aplicación.

```go
func main() {
    // 1. Cargar configuración
    // 2. Inicializar logger
    // 3. Conectar bases de datos
    // 4. Crear container DI
    // 5. Configurar router
    // 6. Iniciar servidor
}
```

**Funciones auxiliares**:
- `getEnvironment()` - Detecta ambiente (local/dev/qa/prod)
- `configureGinMode()` - Configura modo Release en producción
- `startServer()` - Inicia servidor HTTP

**Líneas de código**: ~135 (vs 230 original) = **41% de reducción** ✅

---

## 📊 Comparación Antes vs Después

| Aspecto | Antes | Después |
|---------|-------|---------|
| **Líneas en main.go** | ~230 | ~135 |
| **Responsabilidades** | 7 mezcladas | 3 claras |
| **Funciones en main.go** | 4 + main | 3 + main |
| **Módulos creados** | 0 | 5 nuevos |
| **Testabilidad** | Baja | Alta |
| **Logging estructurado** | Parcial | Completo |
| **Separación de capas** | No | Sí (Clean Arch) |

---

## 🧪 Beneficios de la Refactorización

### 1. **Testabilidad**
Ahora cada módulo se puede testear de forma independiente:

```go
// Test de health check
func TestHealthHandler_Check_AllHealthy(t *testing.T) {
    // Mock de DB y MongoDB
    handler := handler.NewHealthHandler(mockDB, mockMongo)
    // Ejecutar test...
}

// Test de inicialización de PostgreSQL
func TestInitPostgreSQL_Success(t *testing.T) {
    // Usar testcontainers para PostgreSQL real
}
```

### 2. **Mantenibilidad**
- ✅ Cambios en CORS solo afectan `middleware/cors.go`
- ✅ Nuevas rutas se agregan en `router/router.go`
- ✅ Cambios en DB no afectan el `main.go`

### 3. **Reutilización**
Los módulos son reutilizables en otros proyectos:
- `database/postgres.go` - Útil en cualquier proyecto con PostgreSQL
- `middleware/cors.go` - Reutilizable en cualquier API HTTP

### 4. **Logging Mejorado**
Ahora todos los módulos usan logging estructurado:

```go
appLogger.Info(ctx, "PostgreSQL conectado exitosamente",
    zap.String("host", cfg.Database.Postgres.Host),
    zap.Int("port", cfg.Database.Postgres.Port),
)
```

### 5. **Separación de Capas (Clean Architecture)**

```
cmd/main.go (punto de entrada)
    ↓
internal/infrastructure/http/router (HTTP layer)
    ↓
internal/container (DI)
    ↓
internal/application/service (business logic)
    ↓
internal/domain (entities, interfaces)
```

---

## 🚀 Cómo Usar

### Agregar Nueva Ruta

**Antes**: Modificar `main.go` directamente (mezclado con todo).

**Ahora**: Agregar función en `router/router.go`:

```go
// router/router.go
func setupNotificationRoutes(rg *gin.RouterGroup, c *container.Container) {
    notifications := rg.Group("/notifications")
    {
        notifications.GET("", c.NotificationHandler.List)
        notifications.POST("", c.NotificationHandler.Create)
    }
}

// Llamar desde setupProtectedRoutes()
func setupProtectedRoutes(rg *gin.RouterGroup, c *container.Container) {
    protected := rg.Group("")
    protected.Use(ginmiddleware.JWTAuthMiddleware(c.JWTManager))
    {
        setupAuthProtectedRoutes(protected, c)
        setupMaterialRoutes(protected, c)
        setupNotificationRoutes(protected, c) // ✅ Nueva ruta
    }
}
```

### Cambiar Configuración CORS

**Antes**: Buscar en 230 líneas de `main.go`.

**Ahora**: Modificar `middleware/cors.go`:

```go
// middleware/cors.go
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Cambiar de "*" a dominios específicos en producción
        allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
        c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigins)
        // ...
    }
}
```

### Agregar Verificación al Health Check

**Antes**: Modificar función en `main.go`.

**Ahora**: Modificar `handler/health_handler.go`:

```go
func (h *HealthHandler) Check(c *gin.Context) {
    // Verificar Redis (ejemplo)
    redisStatus := "healthy"
    if err := h.redis.Ping(ctx).Err(); err != nil {
        redisStatus = "unhealthy"
    }

    response := HealthResponse{
        // ... campos existentes
        Redis: redisStatus, // ✅ Nueva verificación
    }
}
```

---

## 📝 Checklist de Migración

- [x] Crear módulo `database/postgres.go`
- [x] Crear módulo `database/mongodb.go`
- [x] Crear módulo `middleware/cors.go`
- [x] Crear módulo `handler/health_handler.go`
- [x] Crear módulo `router/router.go`
- [x] Refactorizar `cmd/main.go`
- [x] Mantener anotaciones Swagger en health check
- [x] Logging estructurado en todos los módulos
- [x] Context propagation en funciones de inicialización
- [ ] **TODO**: Tests unitarios para nuevos módulos
- [ ] **TODO**: Tests de integración con testcontainers
- [ ] **TODO**: Configurar CORS dinámico según ambiente

---

## 🎓 Lecciones Aprendidas

### ✅ Buenas Prácticas Aplicadas

1. **Single Responsibility Principle (SRP)**
   - Cada módulo tiene una responsabilidad única y bien definida

2. **Dependency Injection**
   - Dependencias se inyectan mediante constructores
   - Facilita testing con mocks

3. **Context Propagation**
   - Todas las funciones de inicialización reciben `context.Context`
   - Permite timeouts y cancelación

4. **Structured Logging**
   - Uso consistente de Zap logger con campos estructurados
   - Facilita troubleshooting en producción

5. **Error Wrapping**
   - Errores se envuelven con contexto usando `fmt.Errorf(..., %w, err)`

### ⚠️ TODOs Futuros

1. **Tests**: Agregar tests para cada módulo nuevo
2. **CORS Dinámico**: Leer orígenes permitidos desde config
3. **Graceful Shutdown**: Implementar cierre ordenado del servidor
4. **Metrics**: Agregar métricas de Prometheus al health check
5. **Tracing**: Integrar OpenTelemetry para distributed tracing

---

## 📚 Referencias

- [Clean Architecture - Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Gin Framework Documentation](https://gin-gonic.com/docs/)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [EduGo Shared Library](https://github.com/EduGoGroup/edugo-shared)

---

**Última actualización**: 2025-01-15  
**Autor**: Equipo EduGo  
**Versión**: 1.0.0