# Changelog

Todos los cambios notables en edugo-api-mobile serán documentados en este archivo.

El formato está basado en [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
y este proyecto adhiere a [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.11.0] - 2025-11-25

### Tipo de Release: patch

- release: merge dev to main - Auth centralizada Sprint 3

---

## [0.10.0] - 2025-11-22

### Tipo de Release: patch

- fix(ci): Deshabilitar tests de integración y corregir issues (#72)
- refactor(sprint-entities): Migración completa a infrastructure entities v0.10.0 (#70)
- feat(sprint-4): migrar job lint a workflow reusable
- feat: Sprint 2 - Completar tareas pendientes (2.8-2.13) (#66)
- docs(sprint-2): corregir estado - paralelismo ya existía
- docs(sprint-2): actualizar SPRINT-STATUS.md post-merge PR #65
- feat: Sprint 2 FASE 2 - Migración Go 1.25 validada (#65)
- docs: actualizar documentación CI/CD desde análisis centralizado
- docs: actualizar documentación CI/CD desde análisis centralizado
- docs: agregar documentación de CI/CD desde análisis centralizado
- chore: actualizar a edugo-infrastructure postgres@v0.9.0
- chore: actualizar a edugo-infrastructure postgres/v0.8.0 (#62)
- fix: convertir test de score a informativo para evitar panic en CI (#60)
- docs: Sprint 05-B - Decisión de Posponer Testing Avanzado hasta Post-MVP
- docs: Sprint 05-A - Reorganización de Testing y Tests de Seguridad
- refactor: migrar tests legacy a suite compartida (eliminar 16 contenedores individuales) (#55)
- fix: retener sqlDB en CreateConnection
- chore: actualizar dependencies (bootstrap v0.9.0, infrastructure v0.7.0) (#53)
- Sprint 03: Repositorios con BD Real + Tests de Integración (#52)
- feat: Sprint 02 - Implementación completa de la capa de dominio (Clean Architecture + DDD) (#51)


---

## [0.4.2] - 2025-11-19

### Tipo de Release: patch

- chore: actualizar a edugo-infrastructure postgres@v0.9.0 (#64)

---

## [0.4.1] - 2025-11-18

### Tipo de Release: patch

- chore: sync dev to main - Infrastructure v0.8.0 (#63)

---

## [0.4.0] - 2025-11-18

### Tipo de Release: minor

- chore: sync dev to main - Sprint 05 (Testing + Docs)
- feat: Sprint 04 - Services y API REST del Sistema de Evaluaciones

---

## [0.3.2] - 2025-11-17

### Tipo de Release: patch



---

## [0.2.1] - 2025-11-17

### Tipo de Release: patch



---

## [0.1.12] - 2025-11-13

### Tipo de Release: patch

- release: v0.2.0 - shared v0.5.0 actualizado (#44)

---

## [0.1.11] - 2025-11-12

### Tipo de Release: patch

- chore: actualizar shared a v0.4.0 + integraciones completadas (#43)

---

## [0.1.10] - 2025-11-12

### Tipo de Release: patch

- Release: Merge dev to main - v0.1.9 (#41)

---

## [0.1.9] - 2025-11-10

### Tipo de Release: patch

- Add dev-init, reset and status scripts with docs (#39)
- 🧪 Mejora Integral de Estrategia de Testing (#37)

---

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
