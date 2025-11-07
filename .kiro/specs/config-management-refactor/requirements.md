# Requirements Document

## Introduction

El sistema actual de gestión de configuración en la API Mobile de EduGo presenta varios problemas que dificultan el desarrollo y despliegue:

1. **Duplicación de configuración**: Valores sensibles aparecen tanto en archivos YAML como en variables de entorno, creando confusión sobre cuál tiene precedencia
2. **Lógica innecesaria**: El código en `loader.go` tiene lógica manual para forzar precedencia de ENV vars que debería ser automática
3. **Complejidad en desarrollo local**: Múltiples formas de ejecutar la aplicación (IDE, editor, Make, Docker) requieren configuraciones diferentes y repetitivas
4. **Secretos mal manejados**: No hay una estrategia clara para separar configuración pública de secretos
5. **Falta de herramientas**: No existe un script para agregar nuevas variables de configuración de manera consistente

Este spec define un sistema simplificado que:
- Separa claramente configuración pública (YAML) de secretos (ENV)
- Elimina lógica innecesaria en el loader
- Proporciona una experiencia consistente en todos los entornos de desarrollo
- Incluye herramientas para gestionar la configuración fácilmente

## Glossary

- **ConfigSystem**: El sistema completo de gestión de configuración que incluye archivos YAML, variables de entorno, y el código Go que las carga
- **PublicConfig**: Configuración no sensible que puede estar en archivos YAML versionados (puertos, timeouts, nombres de colas)
- **SecretConfig**: Configuración sensible que NUNCA debe estar en YAML (passwords, API keys, tokens)
- **EnvironmentFile**: Archivo `.env` que contiene variables de entorno para desarrollo local
- **ConfigLoader**: Componente Go que carga y valida la configuración usando Viper
- **ConfigTemplate**: Archivo de plantilla que usa marcadores `${VAR_NAME}` para valores que deben ser sustituidos
- **ConfigCLI**: Herramienta de línea de comandos para gestionar variables de configuración
- **HierarchyPath**: Ruta en notación de puntos que representa la ubicación de una variable en la estructura YAML (ej: `database.postgres.password`)

## Requirements

### Requirement 1: Separación Clara de Configuración Pública y Secretos

**User Story:** Como desarrollador, quiero que la configuración pública esté en YAML y los secretos solo en variables de entorno, para que sea obvio qué información es sensible y no haya duplicación.

#### Acceptance Criteria

1. WHEN THE ConfigSystem loads configuration, THE ConfigLoader SHALL read public configuration from YAML files
2. WHEN THE ConfigSystem loads configuration, THE ConfigLoader SHALL read secret configuration exclusively from environment variables
3. THE PublicConfig files SHALL NOT contain any password, API key, token, or connection string with credentials
4. WHERE a configuration value is sensitive, THE PublicConfig files SHALL use placeholder comments indicating the environment variable name
5. THE ConfigLoader SHALL validate that all required SecretConfig variables are present before starting the application

### Requirement 2: Simplificación del Loader de Configuración

**User Story:** Como desarrollador, quiero que Viper maneje automáticamente la precedencia de configuración, para eliminar código manual innecesario y reducir bugs.

#### Acceptance Criteria

1. THE ConfigLoader SHALL use Viper's built-in precedence mechanism without manual overrides
2. THE ConfigLoader SHALL follow this precedence order: environment variables > environment-specific YAML > base YAML > defaults
3. THE ConfigLoader SHALL NOT contain manual `Set()` calls to force environment variable precedence
4. WHEN an environment variable is defined, THE ConfigLoader SHALL automatically override any YAML value for that key
5. THE ConfigLoader SHALL bind environment variables using a consistent naming convention without requiring manual binding for each variable

### Requirement 3: Experiencia Consistente en Desarrollo Local

**User Story:** Como desarrollador, quiero una forma única y simple de configurar mi entorno local, para que funcione igual en IDE, editor de texto, Make, y Docker.

#### Acceptance Criteria

1. THE ConfigSystem SHALL provide a single `.env` file that works for all local development scenarios
2. WHEN running via IDE, THE ConfigSystem SHALL load configuration from `.env` file automatically
3. WHEN running via Make, THE ConfigSystem SHALL load configuration from `.env` file automatically
4. WHEN running via Docker Compose, THE ConfigSystem SHALL load configuration from `.env` file automatically
5. THE ConfigSystem SHALL provide a `.env.example` file with all required variables documented with clear descriptions

### Requirement 4: Gestión de Múltiples Ambientes

**User Story:** Como DevOps, quiero que cada ambiente (local, dev, qa, prod) tenga su propia configuración, para que los valores específicos del ambiente estén claramente separados.

#### Acceptance Criteria

1. THE ConfigSystem SHALL support environment-specific configuration files named `config-{env}.yaml`
2. WHEN THE ConfigSystem determines the environment, THE ConfigLoader SHALL read the `APP_ENV` environment variable
3. WHERE `APP_ENV` is not set, THE ConfigLoader SHALL default to "local" environment
4. THE ConfigLoader SHALL merge environment-specific configuration over base configuration
5. THE ConfigSystem SHALL provide separate configuration files for local, dev, qa, and prod environments

### Requirement 5: Herramienta CLI para Gestión de Configuración

**User Story:** Como desarrollador, quiero una herramienta CLI que me permita agregar nuevas variables de configuración, para que el proceso sea consistente y no tenga que editar múltiples archivos manualmente.

#### Acceptance Criteria

1. THE ConfigCLI SHALL accept a variable name, value type, and hierarchy path as input
2. WHEN adding a public configuration variable, THE ConfigCLI SHALL update all environment-specific YAML files with the new variable
3. WHEN adding a secret configuration variable, THE ConfigCLI SHALL update the `.env.example` file with a placeholder and documentation
4. THE ConfigCLI SHALL update the Go `Config` struct with the new field using proper mapstructure tags
5. THE ConfigCLI SHALL validate that the hierarchy path is valid before making changes
6. THE ConfigCLI SHALL provide a dry-run mode to preview changes before applying them

### Requirement 6: Validación Robusta de Configuración

**User Story:** Como desarrollador, quiero que la aplicación falle rápido con mensajes claros si falta configuración requerida, para detectar problemas antes de que la aplicación inicie.

#### Acceptance Criteria

1. THE ConfigLoader SHALL validate all required configuration fields after loading
2. WHEN a required SecretConfig variable is missing, THE ConfigLoader SHALL return a clear error message indicating which variable is missing
3. WHEN a configuration value has an invalid format, THE ConfigLoader SHALL return a clear error message with the expected format
4. THE ConfigLoader SHALL validate configuration before any application component initializes
5. THE ConfigLoader SHALL log all loaded configuration keys (without values) at startup for debugging

### Requirement 7: Documentación de Variables de Configuración

**User Story:** Como nuevo desarrollador en el equipo, quiero documentación clara de todas las variables de configuración, para entender qué hace cada una y cómo configurarlas.

#### Acceptance Criteria

1. THE ConfigSystem SHALL provide a `CONFIG.md` file documenting all configuration variables
2. THE documentation SHALL include for each variable: name, type, description, default value, and whether it's required
3. THE documentation SHALL indicate which variables are secrets and must be provided via environment variables
4. THE documentation SHALL provide examples for each environment (local, dev, qa, prod)
5. THE `.env.example` file SHALL contain inline comments explaining each variable

### Requirement 8: Compatibilidad con Secretos en la Nube

**User Story:** Como DevOps, quiero que el sistema de configuración sea compatible con servicios de secretos en la nube (AWS Secrets Manager, etc.), para poder desplegar en producción de manera segura.

#### Acceptance Criteria

1. THE ConfigSystem SHALL support loading secrets from environment variables without requiring file-based configuration
2. THE ConfigLoader SHALL NOT require YAML files to be present if all configuration is provided via environment variables
3. WHEN deploying to cloud environments, THE ConfigSystem SHALL allow all configuration to be provided via environment variables
4. THE ConfigSystem SHALL provide documentation on how to integrate with AWS Secrets Manager and similar services
5. THE ConfigLoader SHALL support a "cloud mode" where YAML files are optional

### Requirement 9: Testing de Configuración

**User Story:** Como desarrollador, quiero tests automatizados para el sistema de configuración, para asegurar que los cambios no rompan la carga de configuración.

#### Acceptance Criteria

1. THE ConfigSystem SHALL include unit tests for the ConfigLoader
2. THE tests SHALL verify that environment variables override YAML values correctly
3. THE tests SHALL verify that validation catches missing required fields
4. THE tests SHALL verify that all environment-specific config files can be loaded successfully
5. THE tests SHALL verify that the precedence order works as expected
