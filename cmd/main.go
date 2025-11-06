package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/EduGoGroup/edugo-api-mobile/internal/config"
	"github.com/EduGoGroup/edugo-api-mobile/internal/container"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/database"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/http/handler"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/http/router"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/messaging/rabbitmq"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/storage/s3"
	"github.com/EduGoGroup/edugo-shared/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

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
	ctx := context.Background()

	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Error cargando configuración: %v", err)
	}

	// Detectar y mostrar ambiente activo
	env := getEnvironment()
	log.Printf("🌍 Ambiente: %s", env)
	log.Printf("📊 Log Level: %s, Format: %s", cfg.Logging.Level, cfg.Logging.Format)

	// Inicializar logger estructurado
	appLogger := logger.NewZapLogger(cfg.Logging.Level, cfg.Logging.Format)
	appLogger.Info("iniciando EduGo API Mobile",
		zap.String("environment", env),
		zap.String("version", "1.0.0"),
	)

	// Inicializar base de datos PostgreSQL
	db, err := database.InitPostgreSQL(ctx, cfg, appLogger)
	if err != nil {
		appLogger.Fatal("error inicializando PostgreSQL", zap.Error(err))
	}
	defer db.Close()

	// Inicializar base de datos MongoDB
	mongoDB, err := database.InitMongoDB(ctx, cfg, appLogger)
	if err != nil {
		appLogger.Fatal("error inicializando MongoDB", zap.Error(err))
	}

	// Inicializar RabbitMQ Publisher
	publisher, err := rabbitmq.NewRabbitMQPublisher(
		cfg.Messaging.RabbitMQ.URL,
		cfg.Messaging.RabbitMQ.Exchanges.Materials,
		appLogger.With(zap.String("component", "rabbitmq-publisher")),
	)
	if err != nil {
		appLogger.Fatal("error inicializando RabbitMQ Publisher", zap.Error(err))
	}
	defer publisher.Close()
	appLogger.Info("RabbitMQ Publisher inicializado correctamente",
		zap.String("exchange", cfg.Messaging.RabbitMQ.Exchanges.Materials),
	)

	// Inicializar cliente S3
	s3Config := s3.S3Config{
		Region:          cfg.Storage.S3.Region,
		BucketName:      cfg.Storage.S3.BucketName,
		AccessKeyID:     cfg.Storage.S3.AccessKeyID,
		SecretAccessKey: cfg.Storage.S3.SecretAccessKey,
		Endpoint:        cfg.Storage.S3.Endpoint,
	}
	s3Client, err := s3.NewS3Client(ctx, s3Config, appLogger.With(zap.String("component", "s3-client")))
	if err != nil {
		appLogger.Fatal("error inicializando cliente S3", zap.Error(err))
	}
	appLogger.Info("cliente S3 inicializado correctamente",
		zap.String("region", cfg.Storage.S3.Region),
		zap.String("bucket", cfg.Storage.S3.BucketName),
	)

	// Obtener JWT secret desde variable de entorno
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		appLogger.Fatal("JWT_SECRET es requerido en las variables de entorno")
	}

	// Crear container de dependencias (DI)
	c := container.NewContainer(db, mongoDB, publisher, s3Client, jwtSecret, appLogger)
	defer c.Close()
	appLogger.Info("container de dependencias inicializado correctamente")

	// Configurar modo de Gin según ambiente
	configureGinMode(env)

	// Crear handler de health check
	healthHandler := handler.NewHealthHandler(db, mongoDB)

	// Configurar router con todas las rutas
	r := router.SetupRouter(c, healthHandler)

	// Iniciar servidor HTTP
	startServer(r, cfg, appLogger)
}

// getEnvironment obtiene el ambiente de ejecución desde variables de entorno.
// Por defecto retorna "local" si no está configurado.
func getEnvironment() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}
	return env
}

// configureGinMode configura el modo de ejecución de Gin según el ambiente.
func configureGinMode(env string) {
	if env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
}

// startServer inicia el servidor HTTP en la dirección y puerto configurados.
func startServer(r *gin.Engine, cfg *config.Config, appLogger logger.Logger) {
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	appLogger.Info("servidor HTTP iniciado",
		zap.String("address", addr),
		zap.Int("port", cfg.Server.Port),
		zap.String("swagger_ui", fmt.Sprintf("http://localhost:%d/swagger/index.html", cfg.Server.Port)),
	)

	log.Printf("🚀 API Mobile ejecutándose en http://localhost:%d", cfg.Server.Port)
	log.Printf("📚 Swagger UI: http://localhost:%d/swagger/index.html", cfg.Server.Port)

	if err := r.Run(addr); err != nil {
		appLogger.Fatal("error iniciando servidor HTTP",
			zap.Error(err),
			zap.String("address", addr),
		)
	}
}
