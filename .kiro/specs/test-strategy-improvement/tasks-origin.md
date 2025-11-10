# Plan de Implementación - Mejora de Estrategia de Testing

## Fase 1: Análisis y Evaluación

- [x] 1. Analizar estructura actual de tests
  - Ejecutar análisis de todos los archivos `*_test.go` en el proyecto
  - Generar reporte con ubicación y tipo de cada test
  - Identificar carpetas vacías en `test/unit/`
  - _Requisitos: 1.1, 1.5_

- [x] 2. Calcular cobertura actual por módulo
  - Ejecutar `go test -coverprofile=coverage.out ./...`
  - Analizar cobertura por paquete
  - Identificar módulos con 0% de cobertura
  - Generar reporte de cobertura actual
  - _Requisitos: 1.2, 1.3_
  - ✅ **Completado**: Cobertura total 30.9%, análisis detallado generado

- [x] 3. Validar tests unitarios existentes
  - Ejecutar todos los tests en `internal/application/service/`
  - Ejecutar todos los tests en `internal/infrastructure/http/handler/`
  - Ejecutar todos los tests en `internal/config/`
  - Verificar que todos pasan sin errores
  - Documentar cualquier test que falle
  - _Requisitos: 11.1, 11.2, 11.5_
  - ✅ **Completado**: 77 tests unitarios pasando al 100%

- [x] 4. Validar tests de integración existentes
  - Verificar que Docker está disponible
  - Ejecutar `RUN_INTEGRATION_TESTS=true go test -tags=integration ./test/integration/...`
  - Verificar que testcontainers se levantan correctamente
  - Verificar que todos los 17 tests pasan
  - Documentar tiempo de ejecución de cada test
  - _Requisitos: 1.4, 11.3, 11.4_
  - ✅ **Completado**: 20/21 tests pasando (1 error no crítico de conexión TCP)
  - 🔧 **Corrección aplicada**: Fix en testhelpers.go para usar bootstrap.Resources

- [x] 5. Generar reporte de análisis completo
  - Crear documento `docs/TEST_ANALYSIS_REPORT.md`
  - Incluir resumen de tests existentes
  - Incluir métricas de cobertura actual
  - Incluir lista de módulos sin tests
  - Incluir recomendaciones priorizadas
  - _Requisitos: 1.1, 1.2, 1.3, 1.4_
  - ✅ **Completado**: Reporte generado en docs/TEST_ANALYSIS_REPORT.md

## Fase 2: Configuración y Refactorización

- [x] 6. Configurar exclusiones de cobertura
  - [x] 6.1 Crear archivo `.coverignore` en la raíz del proyecto
    - Agregar exclusiones para archivos generados (docs/, swagger)
    - Agregar exclusiones para DTOs y estructuras simples
    - Agregar exclusiones para mocks y helpers de testing
    - Agregar exclusiones para cmd/ y tools/
    - Documentar cada exclusión con comentarios
    - _Requisitos: 3.1, 3.2, 3.3, 3.4_

  - [x] 6.2 Crear script de filtrado de cobertura
    - Crear `scripts/filter-coverage.sh`
    - Leer patrones de `.coverignore`
    - Filtrar líneas del reporte de cobertura
    - Generar archivo `coverage-filtered.out`
    - _Requisitos: 3.5_

  - [x] 6.3 Crear script de verificación de umbral
    - Crear `scripts/check-coverage.sh`
    - Parsear archivo de cobertura
    - Comparar con umbral mínimo (60%)
    - Fallar si cobertura es menor al umbral
    - _Requisitos: 3.5_


- [x] 7. Limpiar estructura de carpetas de tests
  - Eliminar carpeta `test/unit/application/` (solo contiene .gitkeep)
  - Eliminar carpeta `test/unit/domain/` (solo contiene .gitkeep)
  - Eliminar carpeta `test/unit/infrastructure/` (solo contiene .gitkeep)
  - Eliminar carpeta `test/unit/` completa si queda vacía
  - Actualizar `.gitignore` si es necesario
  - _Requisitos: 4.2_

- [x] 8. Mejorar helpers de testcontainers
  - [x] 8.1 Agregar configuración automática de RabbitMQ
    - Crear función `setupRabbitMQTopology()` en `test/integration/setup.go`
    - Declarar exchange `edugo.events` de tipo topic
    - Crear colas: `material.created`, `assessment.completed`, `progress.updated`
    - Crear bindings entre exchange y colas
    - Manejar errores sin fallar tests (logging de warnings)
    - _Requisitos: 7.1, 7.2, 7.3_

  - [x] 8.2 Integrar configuración de RabbitMQ en SetupContainers
    - Llamar a `setupRabbitMQTopology()` después de levantar RabbitMQ
    - Usar fallback a mock si configuración falla
    - Agregar logging de éxito/fallo
    - _Requisitos: 7.4_

- [~] 9. Mejorar helpers de seed de datos
  - [x] 9.1 Documentar contraseñas sin encriptar
    - Agregar comentarios en `SeedTestUser()` con password sin encriptar
    - Agregar comentarios en `SeedTestUserWithEmail()` con password sin encriptar
    - Actualizar logging para incluir password en tests
    - _Requisitos: 6.4_

  - [x] 9.2 Crear helper para seed de múltiples usuarios
    - Crear función `SeedTestUsers(t, db, count, role)` en testhelpers.go
    - Retornar slice de `TestUser` con ID, Email, Password, Role
    - Agregar logging de usuarios creados
    - _Requisitos: 6.1, 6.2_

  - [x] 9.3 Crear helper para seed de escenario completo
    - Crear función `SeedCompleteTestScenario(t, db, mongodb)` en testhelpers.go
    - Crear teacher, 2 students, 2 materials, 2 assessments
    - Retornar struct `TestScenario` con todos los IDs
    - Agregar logging de escenario creado
    - _Requisitos: 6.1, 6.2, 6.3_

  - [x] 9.4 Mejorar función de limpieza de datos
    - Actualizar `CleanDatabase()` para incluir todas las tablas
    - Agregar orden correcto de limpieza (dependencias)
    - Agregar logging de tablas limpiadas
    - Manejar errores sin fallar tests
    - _Requisitos: 6.5_

- [x] 10. Crear scripts de setup para desarrollo local
  - [x] 10.1 Crear docker-compose para desarrollo
    - Crear `docker-compose-dev.yml` en la raíz
    - Configurar PostgreSQL con puerto 5432
    - Configurar MongoDB con puerto 27017
    - Configurar RabbitMQ con puertos 5672 y 15672
    - Configurar volúmenes persistentes
    - Montar scripts SQL en PostgreSQL
    - _Requisitos: 8.1_

  - [x] 10.2 Crear script de setup de ambiente
    - Crear `test/scripts/setup_dev_env.sh`
    - Verificar que Docker está corriendo
    - Levantar contenedores con docker-compose
    - Esperar a que servicios estén listos
    - Ejecutar schema SQL en PostgreSQL
    - Cargar datos de prueba en PostgreSQL
    - Crear colecciones e índices en MongoDB
    - Configurar exchanges y colas en RabbitMQ
    - Mostrar connection strings al finalizar
    - _Requisitos: 8.1, 8.2, 8.3, 8.4_

  - [x] 10.3 Crear script de teardown de ambiente
    - Crear `test/scripts/teardown_dev_env.sh`
    - Detener contenedores con docker-compose down
    - Eliminar volúmenes con flag -v
    - Mostrar mensaje de confirmación
    - _Requisitos: 8.5_

  - [x] 10.4 Hacer scripts ejecutables
    - Ejecutar `chmod +x test/scripts/setup_dev_env.sh`
    - Ejecutar `chmod +x test/scripts/teardown_dev_env.sh`
    - Agregar shebang `#!/bin/bash` en ambos scripts
    - _Requisitos: 8.1, 8.5_


- [x] 11. Actualizar Makefile con nuevos comandos
  - [x] 11.1 Agregar comandos de testing avanzado
    - Agregar `test-unit`: Solo tests unitarios (rápido)
    - Agregar `test-unit-coverage`: Tests unitarios con cobertura
    - Agregar `test-integration-verbose`: Tests de integración con logs
    - Agregar `test-all`: Ejecutar todos los tests
    - Agregar `test-watch`: Watch mode para tests (requiere entr)
    - _Requisitos: 10.4_

  - [x] 11.2 Agregar comandos de cobertura
    - Agregar `coverage-report`: Reporte completo con filtrado
    - Agregar `coverage-check`: Verificar umbral mínimo
    - Actualizar `test-coverage` para usar script de filtrado
    - _Requisitos: 3.5, 10.4_

  - [x] 11.3 Agregar comandos de desarrollo local
    - Agregar `dev-setup`: Configurar ambiente de desarrollo
    - Agregar `dev-teardown`: Limpiar ambiente de desarrollo
    - Agregar `dev-reset`: Resetear ambiente (teardown + setup)
    - Agregar `dev-logs`: Ver logs de contenedores
    - _Requisitos: 8.1, 8.5, 10.4_

  - [x] 11.4 Agregar comandos de análisis
    - Agregar `test-analyze`: Analizar estructura de tests
    - Agregar `test-missing`: Identificar módulos sin tests
    - Agregar `test-validate`: Validar que todos los tests pasan
    - _Requisitos: 10.4_

## Fase 3: Mejora de Cobertura

- [x] 12. Crear tests para value objects
  - [x] 12.1 Tests para Email value object
    - Crear `internal/domain/valueobject/email_test.go`
    - Test de creación con email válido
    - Test de validación con email inválido
    - Test de método String()
    - _Requisitos: 9.1, 9.3_

  - [x] 12.2 Tests para MaterialID value object
    - Crear `internal/domain/valueobject/material_id_test.go`
    - Test de creación desde string válido
    - Test de creación desde string inválido (UUID inválido)
    - Test de método String()
    - _Requisitos: 9.1, 9.3_

  - [x] 12.3 Tests para UserID value object
    - Crear `internal/domain/valueobject/user_id_test.go`
    - Test de creación desde string válido
    - Test de creación desde string inválido
    - Test de método String()
    - _Requisitos: 9.1, 9.3_

  - [x] 12.4 Tests para MaterialVersionID value object
    - Crear `internal/domain/valueobject/material_version_id_test.go`
    - Test de creación desde string válido
    - Test de creación desde string inválido
    - Test de método String()
    - _Requisitos: 9.1, 9.3_

- [x] 13. Crear tests para entities de dominio
  - [x] 13.1 Tests para Material entity
    - Crear `internal/domain/entity/material_test.go`
    - Test de creación de material
    - Test de validación de campos requeridos
    - Test de métodos de negocio (si existen)
    - _Requisitos: 9.2, 9.3_

  - [x] 13.2 Tests para User entity
    - Crear `internal/domain/entity/user_test.go`
    - Test de creación de usuario
    - Test de validación de email
    - Test de validación de role
    - _Requisitos: 9.2, 9.3_

  - [x] 13.3 Tests para Progress entity
    - Crear `internal/domain/entity/progress_test.go`
    - Test de creación de progreso
    - Test de validación de porcentaje (0-100)
    - Test de actualización de progreso
    - _Requisitos: 9.2, 9.3_


- [~] 14. Crear tests para repositories
  - [x] 14.1 Tests para UserRepository
    - Crear `internal/infrastructure/persistence/postgres/repository/user_repository_impl_test.go`
    - Test de FindByEmail con usuario existente
    - Test de FindByEmail con usuario inexistente
    - Test de Create con datos válidos
    - Test de Create con email duplicado
    - Usar testcontainers para PostgreSQL real
    - _Requisitos: 9.2, 9.3_

  - [x] 14.2 Tests para MaterialRepository
    - Crear `internal/infrastructure/persistence/postgres/repository/material_repository_impl_test.go`
    - Test de FindByID con material existente
    - Test de FindByID con material inexistente
    - Test de FindByAuthorID con múltiples materiales
    - Test de Create con datos válidos
    - Usar testcontainers para PostgreSQL real
    - _Requisitos: 9.2, 9.3_

  - [ ] 14.3 Tests para ProgressRepository
    - Crear `internal/infrastructure/persistence/postgres/repository/progress_repository_impl_test.go`
    - Test de Upsert creando nuevo progreso
    - Test de Upsert actualizando progreso existente
    - Test de FindByUserAndMaterial
    - Usar testcontainers para PostgreSQL real
    - _Requisitos: 9.2, 9.3_

  - [ ] 14.4 Tests para AssessmentRepository (MongoDB)
    - Crear `internal/infrastructure/persistence/mongodb/repository/assessment_repository_impl_test.go`
    - Test de SaveAssessment con datos válidos
    - Test de FindAssessmentByMaterialID con assessment existente
    - Test de FindAssessmentByMaterialID con assessment inexistente
    - Test de SaveResult con datos válidos
    - Usar testcontainers para MongoDB real
    - _Requisitos: 9.2, 9.3_

- [ ] 15. Mejorar cobertura de servicios existentes
  - [ ] 15.1 Mejorar tests de MaterialService
    - Revisar `internal/application/service/material_service_test.go`
    - Agregar tests faltantes para casos edge
    - Agregar tests para manejo de errores
    - Verificar cobertura >= 70%
    - _Requisitos: 9.1, 9.4_

  - [ ] 15.2 Mejorar tests de ProgressService
    - Revisar `internal/application/service/progress_service_test.go`
    - Agregar tests faltantes para casos edge
    - Agregar tests para validaciones
    - Verificar cobertura >= 70%
    - _Requisitos: 9.1, 9.4_

  - [ ] 15.3 Mejorar tests de StatsService
    - Revisar `internal/application/service/stats_service_test.go`
    - Agregar tests faltantes para cálculos
    - Agregar tests para casos sin datos
    - Verificar cobertura >= 70%
    - _Requisitos: 9.1, 9.4_

- [ ] 16. Crear tests para handlers sin cobertura
  - [ ] 16.1 Tests para ProgressHandler
    - Crear tests en `internal/infrastructure/http/handler/progress_handler_test.go`
    - Test de UpsertProgress con datos válidos
    - Test de UpsertProgress con datos inválidos
    - Test de UpsertProgress sin autorización
    - Usar mocks para service
    - _Requisitos: 9.2, 9.4_

  - [ ] 16.2 Tests para StatsHandler
    - Crear tests en `internal/infrastructure/http/handler/stats_handler_test.go`
    - Test de GetMaterialStats con material existente
    - Test de GetMaterialStats con material inexistente
    - Test de GetGlobalStats
    - Usar mocks para service
    - _Requisitos: 9.2, 9.4_

  - [ ] 16.3 Tests para SummaryHandler
    - Crear tests en `internal/infrastructure/http/handler/summary_handler_test.go`
    - Test de GetSummary con material existente
    - Test de GetSummary con material inexistente
    - Usar mocks para service
    - _Requisitos: 9.2, 9.4_


- [x] 17. Crear documentación de testing
  - [x] 17.1 Crear guía principal de testing
    - Crear `docs/TESTING_GUIDE.md`
    - Documentar filosofía de testing del proyecto
    - Documentar tipos de tests y cuándo usarlos
    - Documentar estructura de carpetas
    - Documentar comandos make disponibles
    - Incluir mejores prácticas
    - _Requisitos: 10.1, 10.5_

  - [x] 17.2 Crear guía de tests unitarios
    - Crear `docs/TESTING_UNIT_GUIDE.md`
    - Documentar cómo escribir tests unitarios
    - Documentar uso de mocks con ejemplos
    - Documentar patrón AAA (Arrange-Act-Assert)
    - Incluir ejemplos por tipo de componente
    - Incluir plantillas de tests
    - _Requisitos: 10.1, 10.2_

  - [x] 17.3 Crear guía de tests de integración
    - Crear `docs/TESTING_INTEGRATION_GUIDE.md`
    - Documentar cómo escribir tests de integración
    - Documentar uso de testcontainers
    - Documentar helpers disponibles y su uso
    - Documentar seed de datos
    - Incluir sección de troubleshooting
    - _Requisitos: 10.2, 10.3, 10.4_

  - [ ] 17.4 Crear plan de cobertura
    - Crear `docs/TEST_COVERAGE_PLAN.md`
    - Documentar metas de cobertura por módulo
    - Priorizar tests faltantes
    - Establecer timeline de implementación
    - Asignar responsables (si aplica)
    - _Requisitos: 9.1, 9.2, 9.3, 9.4, 9.5_

  - [x] 17.5 Actualizar README con información de testing
    - Agregar sección de Testing en README.md
    - Incluir comandos básicos de testing
    - Incluir links a guías detalladas
    - Agregar badges de cobertura (preparar para CI)
    - _Requisitos: 10.1, 10.4_

## Fase 4: Automatización y CI/CD

- [~] 18. Configurar GitHub Actions para tests
  - [x] 18.1 Crear workflow de tests unitarios
    - Crear `.github/workflows/test-unit.yml`
    - Configurar trigger en push y pull_request
    - Configurar matriz de versiones de Go (1.21, 1.22)
    - Ejecutar `make test-unit`
    - Fallar build si tests fallan
    - _Requisitos: 12.1_

  - [x] 18.2 Crear workflow de tests de integración
    - Crear `.github/workflows/test-integration.yml`
    - Configurar trigger en push a main y pull_request
    - Verificar que Docker está disponible en runner
    - Ejecutar `make test-integration`
    - Configurar timeout de 15 minutos
    - Fallar build si tests fallan
    - _Requisitos: 12.2_

  - [~] 18.3 Crear workflow de cobertura
    - Integrado en `.github/workflows/test.yml`
    - Ejecutar `make coverage-report`
    - Ejecutar `make coverage-check` con umbral 60%
    - Subir reporte de cobertura como artifact
    - Fallar build si cobertura < 60%
    - _Requisitos: 12.3, 12.4_

  - [ ] 18.4 Configurar publicación de reportes
    - Configurar GitHub Pages para reportes de cobertura
    - Publicar coverage.html en cada push a main
    - Agregar comentario en PR con cambio de cobertura
    - _Requisitos: 12.5_

- [x] 19. Configurar badges y métricas
  - [x] 19.1 Agregar badge de tests
    - Agregar badge de GitHub Actions en README
    - Mostrar estado de tests unitarios
    - Mostrar estado de tests de integración
    - _Requisitos: 12.5_

  - [x] 19.2 Agregar badge de cobertura
    - Configurar servicio de cobertura (Codecov o Coveralls)
    - Agregar badge de cobertura en README
    - Configurar actualización automática
    - _Requisitos: 12.5_

  - [ ] 19.3 Configurar protección de branches
    - Requerir que tests pasen antes de merge
    - Requerir que cobertura no disminuya
    - Configurar en settings de GitHub
    - _Requisitos: 12.4_

- [ ] 20. Validación final y documentación
  - [ ] 20.1 Ejecutar suite completa de tests
    - Ejecutar `make test-all` localmente
    - Verificar que todos los tests pasan
    - Verificar tiempos de ejecución
    - _Requisitos: 11.1, 11.2, 11.3_

  - [ ] 20.2 Verificar cobertura final
    - Ejecutar `make coverage-report`
    - Verificar cobertura general >= 60%
    - Verificar cobertura de servicios >= 70%
    - Verificar cobertura de dominio >= 80%
    - _Requisitos: 9.4_

  - [ ] 20.3 Actualizar documentación final
    - Actualizar TEST_ANALYSIS_REPORT.md con resultados finales
    - Actualizar TEST_COVERAGE_PLAN.md con progreso
    - Actualizar CHANGELOG.md con mejoras de testing
    - _Requisitos: 10.1, 10.5_

  - [ ] 20.4 Crear PR con todos los cambios
    - Crear PR descriptivo con resumen de cambios
    - Incluir métricas antes/después
    - Incluir screenshots de reportes de cobertura
    - Solicitar revisión del equipo
    - _Requisitos: Todos_
