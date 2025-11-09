# Changelog

Todos los cambios notables en edugo-api-mobile serán documentados en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
y este proyecto adhiere a [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.8] - 2025-11-09

### Tipo de Release: patch

- fix(ci): ajustar umbral de cobertura a 33% (cobertura real del proyecto)
- fix(ci): ajustar umbral de cobertura y eliminar comentario problemático
- fix(ci): corregir sintaxis JavaScript en comentario de cobertura
- fix(ci): corregir sintaxis YAML en test.yml
- docs: agregar reporte de verificación pre-merge
- fix: usar rabbitmq/amqp091-go en lugar de streadway/amqp deprecado
- chore: sincronizar go.mod y go.sum
- style: formatear código con gofmt
- docs: agregar resumen ejecutivo de completación del proyecto
- feat(ci-cd): completar Fase 4 - Automatización con Controles ON/OFF (Tareas 18-19)
- docs(testing): completar Fase 3 - Documentación y Tests (Tareas 14-17)
- docs(test-strategy): agregar resumen final de progreso
- test(repositories): agregar tests para UserRepository (Tarea 14 - parcial)
- test(domain): completar Tareas 12-13 - Tests para Value Objects y Entities
- feat(test-strategy): completar Fase 2 - Configuración y Refactorización
- feat(test-strategy): completar Tareas 6-8 - Configuración y Refactorización (Fase 2 parcial)
- docs(test-strategy): completar Fase 1 - Análisis y Evaluación

---

## [0.1.7] - 2025-11-09

### Tipo de Release: patch

- fix(swagger): corregir ruta de endpoint /health en documentación Swagger
- feat(workflows): implementar estrategia fast-forward para sincronización main ↔ dev

---

## [0.1.6] - 2025-11-08

### Tipo de Release: patch

- the main (#32)

---

## [0.1.5] - 2025-11-07

### Tipo de Release: patch

- Config (#29)

---

## [0.1.4] - 2025-11-07

### Tipo de Release: patch

- Dev (#26)

---

## [0.1.3] - 2025-11-07

### Tipo de Release: patch

- fix(docker): permitir scripts/ en contexto de Docker build
- fix(docker): copiar scripts desde contexto de build en lugar de stage builder
- chore: release v0.1.3
- Dev (#25)
- feat(services): completar queries complejas y análisis arquitectónico - FASE 2.3 (#17)
- test: agregar tests unitarios y refactorizar main.go (#14)
- fix: verificar commits en sync-main-to-dev
- docs: documentar GitHub App y actualizar a v2.1.4
- feat: implementar GitHub App Token para sincronización automática

---

## [0.1.3] - 2025-11-07

### Tipo de Release: patch

- Dev (#25)
- feat(services): completar queries complejas y análisis arquitectónico - FASE 2.3 (#17)
- test: agregar tests unitarios y refactorizar main.go (#14)
- fix: verificar commits en sync-main-to-dev
- docs: documentar GitHub App y actualizar a v2.1.4
- feat: implementar GitHub App Token para sincronización automática

---

## [0.1.2] - 2025-11-03

### Tipo de Release: patch

- fix: corregir workflow sync-main-to-dev y actualizar documentación (#9)
- Dev (#8)
- docs: actualizar plan CI/CD con workflow manual-release TODO-EN-UNO

---

## [0.1.1] - 2025-11-01

### Tipo de Release: patch

- feat: workflow TODO-EN-UNO para release completo (#7)

---

## [0.1.0] - 2025-11-01

### Tipo de Release: minor

- feat: reemplazar auto-version con workflow manual controlado (#6)

---

## [Unreleased]

## [0.0.1] - 2025-11-01

### ⚠️ BREAKING CHANGES

- fix: resetear version.txt a 0.0.0 (#5)

### Fixed

- resetear version.txt a 0.0.0 (#5)

## [1.0.0] - 2025-10-31

### Added
- Sistema GitFlow profesional implementado
- Workflows de CI/CD automatizados:
  - CI Pipeline con tests y validaciones
  - Tests con cobertura y servicios de infraestructura
  - Build y push automático de Docker images
  - Release automático con versionado semántico
  - Sincronización automática main ↔ dev
- Auto-versionado basado en Conventional Commits
- Migración a edugo-shared v2.0.5 con arquitectura modular
- Submódulos: common, logger, database/postgres
- .gitignore completo para Go
- Documentación completa de workflows

### Changed
- Actualizado a Go 1.25.3
- Optimización de dependencias (reducción ~80%)

### Fixed
- Corrección de errores de linter (errcheck)
- Permisos de GitHub Container Registry configurados

[Unreleased]: https://github.com/EduGoGroup/edugo-api-mobile/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/EduGoGroup/edugo-api-mobile/releases/tag/v1.0.0
