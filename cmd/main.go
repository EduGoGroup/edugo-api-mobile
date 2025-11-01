package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/config"
	"github.com/EduGoGroup/edugo-api-mobile/internal/container"
	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "github.com/EduGoGroup/edugo-api-mobile/docs" // Swagger docs generados por swag init
)

// @title EduGo API Mobile
// @version 1.0
// @description API para operaciones frecuentes de docentes y estudiantes en EduGo
// @termsOfService http://edugo.com/terms/

// @contact.name Equipo EduGo
// @contact.email soporte@edugo.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT token con formato: Bearer {token}

func main() {
	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Error loading configuration: %v", err)
	}

	// Mostrar ambiente activo
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}
	log.Printf("🌍 Environment: %s", env)
	log.Printf("📊 Log Level: %s, Format: %s", cfg.Logging.Level, cfg.Logging.Format)

	// Inicializar logger
	appLogger := initLogger(cfg.Logging.Level, cfg.Logging.Format)
	appLogger.Info("Iniciando EduGo API Mobile...")

	// Inicializar PostgreSQL
	db, err := initPostgreSQL(cfg)
	if err != nil {
		log.Fatalf("❌ Error connecting to PostgreSQL: %v", err)
	}
	defer db.Close()
	log.Printf("✅ PostgreSQL conectado: %s:%d/%s", cfg.Database.Postgres.Host, cfg.Database.Postgres.Port, cfg.Database.Postgres.Database)

	// Inicializar MongoDB
	mongoDB, err := initMongoDB(cfg)
	if err != nil {
		log.Fatalf("❌ Error connecting to MongoDB: %v", err)
	}
	log.Printf("✅ MongoDB conectado: %s", cfg.Database.MongoDB.Database)

	// Obtener JWT secret desde variable de entorno
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("❌ JWT_SECRET environment variable is required")
	}

	// Crear container de dependencias
	c := container.NewContainer(db, mongoDB, jwtSecret, appLogger)
	defer c.Close()
	log.Printf("✅ Container de dependencias inicializado")

	// Configurar Gin
	if env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// Middleware global (usando los básicos de Gin por ahora)
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	// Health check
	r.GET("/health", healthCheckHandler(db, mongoDB))

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Rutas públicas
	v1 := r.Group("/v1")
	{
		// Autenticación
		v1.POST("/auth/login", c.AuthHandler.Login)
	}

	// Rutas protegidas (requieren JWT)
	protected := v1.Group("")
	protected.Use(jwtAuthMiddleware(c.JWTManager))
	{
		// Materials
		materials := protected.Group("/materials")
		{
			materials.GET("", c.MaterialHandler.ListMaterials)
			materials.POST("", c.MaterialHandler.CreateMaterial)
			materials.GET("/:id", c.MaterialHandler.GetMaterial)
			materials.POST("/:id/upload-complete", c.MaterialHandler.NotifyUploadComplete)
			materials.GET("/:id/summary", c.SummaryHandler.GetSummary)
			materials.GET("/:id/assessment", c.AssessmentHandler.GetAssessment)
			materials.POST("/:id/assessment/attempts", c.AssessmentHandler.RecordAttempt)
			materials.PATCH("/:id/progress", c.ProgressHandler.UpdateProgress)
			materials.GET("/:id/stats", c.StatsHandler.GetMaterialStats)
		}
	}

	// Iniciar servidor usando configuración
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("🚀 API Mobile running on http://localhost:%d", cfg.Server.Port)
	log.Printf("📚 Swagger UI: http://localhost:%d/swagger/index.html", cfg.Server.Port)

	if err := r.Run(addr); err != nil {
		log.Fatalf("❌ Error starting server: %v", err)
	}
}

// initLogger inicializa el logger compartido
func initLogger(level, format string) logger.Logger {
	return logger.NewZapLogger(level, format)
}

// initPostgreSQL inicializa la conexión a PostgreSQL
func initPostgreSQL(cfg *config.Config) (*sql.DB, error) {
	connStr := cfg.Database.Postgres.GetConnectionString()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening postgres connection: %w", err)
	}

	// Configurar pool de conexiones
	if cfg.Database.Postgres.MaxConnections > 0 {
		db.SetMaxOpenConns(cfg.Database.Postgres.MaxConnections)
	}
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Hour)

	// Verificar conexión
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error pinging postgres: %w", err)
	}

	return db, nil
}

// initMongoDB inicializa la conexión a MongoDB
func initMongoDB(cfg *config.Config) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Database.MongoDB.Timeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.Database.MongoDB.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error connecting to mongodb: %w", err)
	}

	// Verificar conexión
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("error pinging mongodb: %w", err)
	}

	return client.Database(cfg.Database.MongoDB.Database), nil
}

// corsMiddleware configura CORS para la aplicación
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// jwtAuthMiddleware valida el token JWT
func jwtAuthMiddleware(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// Extraer token del header "Bearer {token}"
		tokenString := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		// Validar token
		claims, err := jwtManager.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// Guardar claims en el contexto para uso en handlers
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// healthCheckHandler godoc
// @Summary Health check
// @Description Verifica que la API está funcionando
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func healthCheckHandler(db *sql.DB, mongoDB *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Verificar PostgreSQL
		pgStatus := "healthy"
		if err := db.Ping(); err != nil {
			pgStatus = "unhealthy"
		}

		// Verificar MongoDB
		mongoStatus := "healthy"
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := mongoDB.Client().Ping(ctx, nil); err != nil {
			mongoStatus = "unhealthy"
		}

		status := "healthy"
		if pgStatus == "unhealthy" || mongoStatus == "unhealthy" {
			status = "degraded"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    status,
			"service":   "edugo-api-mobile",
			"version":   "1.0.0",
			"postgres":  pgStatus,
			"mongodb":   mongoStatus,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	}
}
