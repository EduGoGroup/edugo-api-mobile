package container

import (
	"database/sql"

	"github.com/EduGoGroup/edugo-api-mobile/internal/bootstrap"
	"github.com/EduGoGroup/edugo-api-mobile/internal/client"
	"github.com/EduGoGroup/edugo-api-mobile/internal/config"
	"github.com/EduGoGroup/edugo-api-mobile/internal/infrastructure/messaging/rabbitmq"
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
	AuthClient       *client.AuthClient // Cliente para validar tokens JWT (local + fallback remoto opcional)
	MessagePublisher rabbitmq.Publisher
	S3Client         bootstrap.S3Storage
}

// NewInfrastructureContainer crea y configura el contenedor de infraestructura
// Parámetros:
//   - db: Conexión a PostgreSQL
//   - mongoDB: Conexión a MongoDB
//   - publisher: Cliente de RabbitMQ para mensajería
//   - s3Client: Cliente de AWS S3 para almacenamiento (interfaz S3Storage)
//   - jwtSecret: Secret para generación de tokens JWT (DEBE ser el mismo que api-admin)
//   - authConfig: Configuración de autenticación (incluye JWT issuer y api-admin opcional)
//   - logger: Logger compartido de la aplicación
func NewInfrastructureContainer(
	db *sql.DB,
	mongoDB *mongo.Database,
	publisher rabbitmq.Publisher,
	s3Client bootstrap.S3Storage,
	jwtSecret string,
	authConfig config.AuthConfig,
	logger logger.Logger,
) *InfrastructureContainer {
	// Determinar el issuer JWT - debe coincidir con api-admin ("edugo-central")
	jwtIssuer := authConfig.JWT.Issuer
	if jwtIssuer == "" {
		jwtIssuer = "edugo-central" // Issuer por defecto compatible con api-admin
	}

	// Crear AuthClient con validación LOCAL preferida
	// La validación remota es opcional y se usa como fallback si está habilitada
	authClient := client.NewAuthClient(client.AuthClientConfig{
		// Validación LOCAL (preferida)
		JWTSecret: jwtSecret,
		JWTIssuer: jwtIssuer,

		// Validación REMOTA (fallback opcional)
		BaseURL:         authConfig.APIAdmin.BaseURL,
		Timeout:         authConfig.APIAdmin.Timeout,
		RemoteEnabled:   authConfig.APIAdmin.RemoteEnabled,
		FallbackEnabled: authConfig.APIAdmin.FallbackEnabled,

		// Cache
		CacheTTL:     authConfig.APIAdmin.CacheTTL,
		CacheEnabled: authConfig.APIAdmin.CacheEnabled,
	})

	return &InfrastructureContainer{
		DB:               db,
		MongoDB:          mongoDB,
		Logger:           logger,
		JWTManager:       auth.NewJWTManager(jwtSecret, jwtIssuer),
		AuthClient:       authClient,
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
