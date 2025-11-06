package container

import (
	"database/sql"

	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/messaging/rabbitmq"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/storage/s3"
	"github.com/EduGoGroup/edugo-shared/auth"
	"github.com/EduGoGroup/edugo-shared/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

// InfrastructureContainer encapsula todas las dependencias de infraestructura
// Responsabilidad: Gestionar recursos externos (DB, Logger, Messaging, Storage, Auth)
type InfrastructureContainer struct {
	DB               *sql.DB
	MongoDB          *mongo.Database
	Logger           logger.Logger
	JWTManager       *auth.JWTManager
	MessagePublisher rabbitmq.Publisher
	S3Client         *s3.S3Client
}

// NewInfrastructureContainer crea y configura el contenedor de infraestructura
// Parámetros:
//   - db: Conexión a PostgreSQL
//   - mongoDB: Conexión a MongoDB
//   - publisher: Cliente de RabbitMQ para mensajería
//   - s3Client: Cliente de AWS S3 para almacenamiento
//   - jwtSecret: Secret para generación de tokens JWT
//   - logger: Logger compartido de la aplicación
func NewInfrastructureContainer(
	db *sql.DB,
	mongoDB *mongo.Database,
	publisher rabbitmq.Publisher,
	s3Client *s3.S3Client,
	jwtSecret string,
	logger logger.Logger,
) *InfrastructureContainer {
	return &InfrastructureContainer{
		DB:               db,
		MongoDB:          mongoDB,
		Logger:           logger,
		JWTManager:       auth.NewJWTManager(jwtSecret, "edugo-mobile"),
		MessagePublisher: publisher,
		S3Client:         s3Client,
	}
}

// Close cierra los recursos de infraestructura
// Actualmente solo cierra la conexión a PostgreSQL
// MongoDB se gestiona externamente por el driver
func (ic *InfrastructureContainer) Close() error {
	if ic.DB != nil {
		return ic.DB.Close()
	}
	return nil
}
