package database

import (
	"context"
	"fmt"

	"github.com/EduGoGroup/edugo-api-mobile/internal/config"
	"github.com/EduGoGroup/edugo-shared/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// InitMongoDB inicializa la conexión a MongoDB con la configuración proporcionada.
// Verifica que la conexión esté activa mediante un ping.
func InitMongoDB(ctx context.Context, cfg *config.Config, log logger.Logger) (*mongo.Database, error) {
	log.Info("inicializando conexión a MongoDB",
		zap.String("database", cfg.Database.MongoDB.Database),
		zap.Duration("timeout", cfg.Database.MongoDB.Timeout),
	)

	connectCtx, cancel := context.WithTimeout(ctx, cfg.Database.MongoDB.Timeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.Database.MongoDB.URI)
	client, err := mongo.Connect(connectCtx, clientOptions)
	if err != nil {
		log.Error("error conectando a MongoDB",
			zap.Error(err),
		)
		return nil, fmt.Errorf("error connecting to mongodb: %w", err)
	}

	// Verificar conexión con ping
	if err := client.Ping(connectCtx, nil); err != nil {
		log.Error("error verificando conexión a MongoDB",
			zap.Error(err),
		)
		return nil, fmt.Errorf("error pinging mongodb: %w", err)
	}

	log.Info("MongoDB conectado exitosamente",
		zap.String("database", cfg.Database.MongoDB.Database),
	)

	return client.Database(cfg.Database.MongoDB.Database), nil
}
