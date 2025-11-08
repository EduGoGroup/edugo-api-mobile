package container

import (
	"github.com/EduGoGroup/edugo-api-mobile/internal/bootstrap"
)

// Container es el contenedor raíz de dependencias de API Mobile
// Implementa el patrón Dependency Injection con segregación por capas
// Resuelve el God Object pattern dividiendo responsabilidades en sub-containers
//
// Arquitectura de contenedores:
//   - Infrastructure: Gestiona recursos externos (DB, Logger, Messaging, Storage)
//   - Repositories: Gestiona acceso a datos (PostgreSQL y MongoDB)
//   - Services: Gestiona lógica de negocio (casos de uso)
//   - Handlers: Gestiona presentación HTTP (REST API)
//
// Beneficios:
//   - SRP: Cada sub-container tiene una responsabilidad clara
//   - Testabilidad: Se pueden mockear sub-containers completos
//   - Mantenibilidad: Cambios localizados por capa
//   - Extensibilidad: Nuevas features se agregan al sub-container correspondiente
type Container struct {
	Infrastructure *InfrastructureContainer
	Repositories   *RepositoryContainer
	Services       *ServiceContainer
	Handlers       *HandlerContainer
}

// NewContainer crea un nuevo contenedor e inicializa todas las dependencias
// de forma jerárquica siguiendo la arquitectura de capas
//
// Flujo de inicialización:
//  1. Infrastructure → Recursos externos (DB, Logger, Messaging, Storage, Auth)
//  2. Repositories   → Persistencia (depende de Infrastructure)
//  3. Services       → Lógica de negocio (depende de Repositories e Infrastructure)
//  4. Handlers       → Presentación HTTP (depende de Services e Infrastructure)
//
// Parámetros:
//   - resources: Recursos de infraestructura inicializados por el bootstrap system
//
// Retorna el contenedor raíz con todos los sub-containers inicializados
func NewContainer(resources *bootstrap.Resources) *Container {
	// Paso 1: Inicializar infraestructura (capa más baja - recursos externos)
	infra := NewInfrastructureContainer(
		resources.PostgreSQL,
		resources.MongoDB,
		resources.RabbitMQPublisher,
		resources.S3Client,
		resources.JWTSecret,
		resources.Logger,
	)

	// Paso 2: Inicializar repositorios (dependen de infraestructura)
	repos := NewRepositoryContainer(infra)

	// Paso 3: Inicializar servicios (dependen de repositorios e infraestructura)
	services := NewServiceContainer(infra, repos)

	// Paso 4: Inicializar handlers (dependen de servicios e infraestructura)
	handlers := NewHandlerContainer(infra, services)

	// Retornar contenedor raíz con todos los sub-containers
	return &Container{
		Infrastructure: infra,
		Repositories:   repos,
		Services:       services,
		Handlers:       handlers,
	}
}

// Close cierra los recursos del contenedor
// Delega el cierre al InfrastructureContainer que gestiona las conexiones
func (c *Container) Close() error {
	return c.Infrastructure.Close()
}
