# Requirements Document

## Introduction

Este documento define los requisitos para mejorar la configuración de Swagger en la aplicación EduGo API Mobile. El objetivo es garantizar que la documentación Swagger esté siempre actualizada, que se regenere automáticamente al iniciar la aplicación, y que el Swagger UI utilice dinámicamente el puerto y host correctos en tiempo de ejecución, permitiendo probar endpoints directamente desde la interfaz sin problemas de configuración.

## Glossary

- **Swagger**: Framework de documentación de APIs que utiliza la especificación OpenAPI
- **swag**: Herramienta CLI de Go para generar documentación Swagger desde anotaciones en código
- **Swagger UI**: Interfaz web interactiva para visualizar y probar APIs documentadas con Swagger
- **Host dinámico**: Configuración que permite que Swagger UI detecte automáticamente el servidor y puerto en tiempo de ejecución
- **Anotaciones Swagger**: Comentarios especiales en código Go que swag utiliza para generar documentación
- **Regeneración automática**: Proceso de ejecutar `swag init` automáticamente al iniciar la aplicación

## Requirements

### Requirement 1

**User Story:** Como desarrollador, quiero que la documentación Swagger se regenere automáticamente al iniciar la aplicación, para que siempre refleje los cambios más recientes en los endpoints sin intervención manual.

#### Acceptance Criteria

1. WHEN THE application starts, THE System SHALL execute the swag init command before initializing the HTTP server
2. WHEN THE swag init command completes successfully, THE System SHALL log a confirmation message indicating Swagger documentation was regenerated
3. IF THE swag init command fails, THEN THE System SHALL log an error message with details and continue application startup
4. THE System SHALL ensure the docs directory contains the updated swagger.json and swagger.yaml files after regeneration

### Requirement 2

**User Story:** Como desarrollador o tester, quiero que Swagger UI detecte automáticamente el puerto y host con el que se levantó la aplicación, para poder probar endpoints directamente desde la interfaz sin errores de conexión.

#### Acceptance Criteria

1. THE System SHALL configure Swagger to use dynamic host and basePath values instead of hardcoded values
2. WHEN THE Swagger UI loads in a browser, THE System SHALL inject the current window location host into the Swagger configuration
3. THE System SHALL remove hardcoded @host annotations from the main.go file
4. WHEN THE application runs on any port, THE Swagger UI SHALL construct API request URLs using the detected host and port
5. THE System SHALL ensure the /health endpoint can be successfully invoked from Swagger UI regardless of the configured port

### Requirement 3

**User Story:** Como desarrollador, quiero que todas las anotaciones Swagger en los handlers estén correctamente configuradas, para que la documentación generada sea completa y precisa.

#### Acceptance Criteria

1. THE System SHALL include proper Swagger annotations for all HTTP handlers with tags, summaries, descriptions, parameters, and responses
2. THE System SHALL define security requirements using @Security annotations for protected endpoints
3. THE System SHALL document all request and response DTOs with proper @Param and @Success annotations
4. THE System SHALL include example values in all DTO definitions for better API documentation
5. THE System SHALL validate that all endpoints defined in router.go have corresponding Swagger annotations in their handlers

### Requirement 4

**User Story:** Como desarrollador, quiero que la configuración de Swagger sea mantenible y esté centralizada, para facilitar cambios futuros en la documentación de la API.

#### Acceptance Criteria

1. THE System SHALL maintain all global Swagger configuration annotations in the main.go file
2. THE System SHALL use consistent naming conventions for tags across all endpoints
3. THE System SHALL document the Bearer token authentication scheme with clear instructions
4. THE System SHALL include API metadata such as title, version, description, contact, and license information
5. THE System SHALL organize Swagger annotations logically to improve code readability
