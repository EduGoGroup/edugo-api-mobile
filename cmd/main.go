package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/EduGoGroup/edugo-api-mobile/internal/bootstrap"
	"github.com/EduGoGroup/edugo-api-mobile/internal/config"
	"github.com/EduGoGroup/edugo-api-mobile/internal/container"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/http/handler"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/http/router"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	_ "github.com/EduGoGroup/edugo-api-mobile/docs" // Swagger docs generados por swag init
)

func main() {
	ctx := context.Background()

	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Error cargando configuración: %v", err)
	}

	// Inicializar infraestructura usando bootstrap system
	b := bootstrap.New(cfg)
	resources, cleanup, err := b.InitializeInfrastructure(ctx)
	if err != nil {
		log.Fatalf("❌ Error inicializando infraestructura: %v", err)
	}
	defer cleanup()

	// Regenerar documentación Swagger
	if err := regenerateSwagger(resources.Logger); err != nil {
		resources.Logger.Warn("continuando inicio de aplicación con documentación Swagger existente",
			zap.Error(err),
		)
	}

	// Crear container de dependencias (DI)
	c := container.NewContainer(resources)
	resources.Logger.Info("container de dependencias inicializado correctamente")

	// Configurar modo de Gin según ambiente
	configureGinMode(cfg.Environment)

	// Crear handler de health check
	healthHandler := handler.NewHealthHandler(resources.PostgreSQL, resources.MongoDB)

	// Configurar router con todas las rutas
	r := router.SetupRouter(c, healthHandler)

	// Iniciar servidor HTTP
	startServer(r, cfg, resources.Logger)
}

// @title EduGo API Mobile
// @version 1.0
// @description API para operaciones frecuentes de docentes y estudiantes en EduGo
// @termsOfService http://edugo.com/terms/

// @contact.name Equipo EduGo
// @contact.email soporte@edugo.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @BasePath /v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT token con formato: Bearer {token}

// regenerateSwagger ejecuta swag init para actualizar la documentación Swagger.
// Captura stdout/stderr para logging y maneja errores con degradación graciosa.
// Retorna error si la regeneración falla, pero no debe detener el inicio de la aplicación.
func regenerateSwagger(log interface{ Info(string, ...interface{}) }) error {
	log.Info("iniciando regeneración de documentación Swagger")

	// Ejecutar comando swag init
	cmd := exec.Command("swag", "init", "-g", "cmd/main.go")

	// Capturar stdout y stderr
	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))

	if err != nil {
		// Verificar si swag no está instalado
		if strings.Contains(err.Error(), "executable file not found") || strings.Contains(err.Error(), "not found") {
			log.Info("swag CLI no encontrado, omitiendo regeneración. Instalar con: go install github.com/swaggo/swag/cmd/swag@latest")
			return fmt.Errorf("swag CLI no encontrado: %w", err)
		}

		// Otro tipo de error
		return fmt.Errorf("error ejecutando swag init: %w (output: %s)", err, outputStr)
	}

	// Regeneración exitosa
	log.Info("documentación Swagger regenerada exitosamente")
	return nil
}

// configureGinMode configura el modo de ejecución de Gin según el ambiente.
func configureGinMode(env string) {
	if env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
}

// startServer inicia el servidor HTTP en la dirección y puerto configurados.
func startServer(r *gin.Engine, cfg *config.Config, appLogger interface {
	Info(string, ...interface{})
	Fatal(string, ...interface{})
}) {
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
