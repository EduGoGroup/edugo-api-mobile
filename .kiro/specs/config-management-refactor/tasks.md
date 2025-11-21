# Implementation Plan

- [x] 1. Preparar entorno y crear rama de desarrollo
  - Crear rama `feature/config-refactor` desde main/develop
  - Asegurar que todos los tests actuales pasen antes de empezar
  - _Requirements: 1.1, 2.1_

- [ ] 2. Limpiar archivos YAML de configuración
  - [x] 2.1 Remover secretos de `config/config.yaml`
    - Eliminar valores de passwords, URIs con credenciales, API keys
    - Agregar comentarios indicando la variable ENV correspondiente
    - _Requirements: 1.3, 1.4_

  - [x] 2.2 Remover secretos de `config/config-local.yaml`
    - Eliminar password, uri, url con credenciales
    - Mantener solo configuración pública (puertos, nombres, timeouts)
    - _Requirements: 1.3, 1.4_

  - [x] 2.3 Remover secretos de `config/config-dev.yaml`
    - Aplicar misma limpieza que en local
    - _Requirements: 1.3, 1.4_

  - [x] 2.4 Remover secretos de `config/config-qa.yaml` y `config/config-prod.yaml`
    - Aplicar misma limpieza en todos los ambientes
    - _Requirements: 1.3, 1.4_

- [ ] 3. Refactorizar el Config Loader
  - [x] 3.1 Crear `internal/config/validator.go`
    - Mover función `Validate()` desde `config.go` a nuevo archivo
    - Implementar validación de campos requeridos con mensajes claros
    - Agregar validación de rangos y formatos
    - _Requirements: 2.1, 6.1, 6.2, 6.3_

  - [x] 3.2 Simplificar `internal/config/loader.go`
    - Remover todas las llamadas a `BindEnv()` manual
    - Remover todas las llamadas a `Set()` manual para forzar precedencia
    - Agregar `AutomaticEnv()` y `SetEnvKeyReplacer()`
    - Mantener lógica de carga de archivos base y específicos por ambiente
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5_

  - [x] 3.3 Actualizar comentarios en `internal/config/config.go`
    - Agregar comentarios en structs indicando qué campos vienen de ENV
    - Documentar formato esperado para cada variable ENV
    - _Requirements: 7.2, 7.3_

- [ ] 4. Actualizar archivo .env.example
  - [x] 4.1 Reorganizar y documentar variables de entorno
    - Agrupar por sección (Database, Messaging, Storage, etc.)
    - Agregar descripción clara para cada variable
    - Indicar formato esperado y ejemplos
    - Documentar cuál campo del config mapea cada ENV var
    - _Requirements: 3.5, 7.5_

  - [x] 4.2 Agregar variables faltantes
    - Asegurar que todas las variables secretas estén documentadas
    - Agregar `STORAGE_S3_ACCESS_KEY_ID` y `STORAGE_S3_SECRET_ACCESS_KEY`
    - _Requirements: 1.3, 3.5_

- [ ] 5. Actualizar configuración de herramientas de desarrollo
  - [x] 5.1 Actualizar Makefile
    - Agregar carga automática de `.env` con `include .env` y `export`
    - Remover variables hardcodeadas en el Makefile
    - Verificar que `make run` funcione correctamente
    - _Requirements: 3.2, 3.3_

  - [x] 5.2 Actualizar docker-compose.yml
    - Simplificar usando `env_file: .env`
    - Remover variables individuales duplicadas
    - Mantener solo variables específicas de Docker si es necesario
    - _Requirements: 3.4_

  - [x] 5.3 Crear/actualizar configuración de IDE
    - Actualizar `.zed/debug.json` para cargar `.env`
    - Crear `.vscode/launch.json` con `envFile` si no existe
    - Documentar cómo configurar IntelliJ/GoLand para cargar `.env`
    - _Requirements: 3.2_

- [ ] 6. Crear herramienta ConfigCTL
  - [x] 6.1 Crear estructura base del CLI
    - Crear directorio `tools/configctl/`
    - Implementar `main.go` con cobra para comandos
    - Agregar comandos: add, validate, generate-docs
    - _Requirements: 5.1_

  - [ ] 6.2 Implementar comando `add` para variables públicas
    - Parsear hierarchy path (ej: `database.postgres.pool_size`)
    - Actualizar struct en `config.go` con nuevo campo
    - Actualizar archivos YAML con valor por defecto
    - _Requirements: 5.2, 5.5_

  - [ ] 6.3 Implementar comando `add` para variables secretas
    - Actualizar struct en `config.go` con nuevo campo
    - Actualizar `.env.example` con documentación
    - Actualizar `validator.go` si la variable es requerida
    - _Requirements: 5.3, 5.4_

  - [x] 6.4 Implementar comando `validate`
    - Validar que todos los archivos YAML sean válidos
    - Validar que `.env.example` tenga todas las variables secretas
    - Validar que structs en `config.go` tengan mapstructure tags correctos
    - _Requirements: 5.5_

  - [ ] 6.5 Implementar comando `generate-docs`
    - Generar `CONFIG.md` con documentación de todas las variables
    - Incluir tipo, descripción, default, y si es requerida
    - Generar tabla de mapeo ENV var → config path
    - _Requirements: 7.1, 7.2, 7.3, 7.4_

  - [ ] 6.6 Implementar modo dry-run
    - Agregar flag `--dry-run` a comando `add`
    - Mostrar preview de cambios sin aplicarlos
    - _Requirements: 5.6_

- [ ] 7. Crear tests para el sistema de configuración
  - [x] 7.1 Tests unitarios para loader
    - Test precedencia: ENV > YAML específico > YAML base > defaults
    - Test carga con archivo YAML faltante (cloud mode)
    - Test que AutomaticEnv mapea correctamente las variables
    - _Requirements: 9.1, 9.2, 9.5_

  - [x] 7.2 Tests unitarios para validator
    - Test validación de campos requeridos faltantes
    - Test validación de rangos (puertos, conexiones)
    - Test mensajes de error claros y útiles
    - _Requirements: 9.1, 9.3_

  - [ ] 7.3 Tests de integración para carga de configuración
    - Test carga exitosa en ambiente local
    - Test carga exitosa en ambiente dev/qa/prod
    - Test fallo rápido con configuración inválida
    - _Requirements: 9.4_

- [ ] 8. Crear documentación
  - [x] 8.1 Crear CONFIG.md
    - Documentar todas las variables de configuración
    - Incluir ejemplos para cada ambiente
    - Documentar cómo agregar nuevas variables con ConfigCTL
    - _Requirements: 7.1, 7.2, 7.3, 7.4_

  - [x] 8.2 Actualizar README.md
    - Agregar sección de configuración
    - Documentar setup inicial: `cp .env.example .env`
    - Documentar cómo ejecutar en diferentes entornos
    - _Requirements: 3.5, 7.5_

  - [x] 8.3 Crear guía de migración para cloud
    - Documentar integración con AWS Secrets Manager
    - Documentar integración con Kubernetes Secrets
    - Proporcionar ejemplos de configuración
    - _Requirements: 8.4_

- [ ] 9. Validación y testing end-to-end
  - [x] 9.1 Crear archivo .env local de prueba
    - Copiar `.env.example` a `.env`
    - Llenar con valores de prueba válidos
    - _Requirements: 3.1, 3.5_

  - [ ] 9.2 Validar ejecución con IDE
    - Configurar run configuration en IDE
    - Ejecutar aplicación y verificar que carga configuración correctamente
    - Verificar logs de startup muestran configuración cargada
    - _Requirements: 3.2, 6.5_

  - [x] 9.3 Validar ejecución con Make
    - Ejecutar `make run`
    - Verificar que carga `.env` automáticamente
    - Verificar que la aplicación inicia correctamente
    - _Requirements: 3.3_

  - [x] 9.4 Validar ejecución con Docker Compose
    - Ejecutar `docker-compose up`
    - Verificar que carga `.env` correctamente
    - Verificar que todos los servicios inician y se conectan
    - _Requirements: 3.4_

  - [x] 9.5 Ejecutar todos los tests
    - Ejecutar `make test` y verificar que todos pasen
    - Ejecutar `make test-integration` si aplica
    - Verificar cobertura de tests
    - _Requirements: 9.1, 9.2, 9.3, 9.4_

- [ ] 10. Probar ConfigCTL agregando una variable de prueba
  - [ ] 10.1 Agregar variable pública de prueba
    - Ejecutar `configctl add test.feature.enabled --type bool --default false --desc "Enable test feature"`
    - Verificar que actualiza `config.go`, YAMLs correctamente
    - Compilar y verificar que no hay errores
    - _Requirements: 5.1, 5.2, 5.5_

  - [ ] 10.2 Agregar variable secreta de prueba
    - Ejecutar `configctl add test.api.token --type string --secret --desc "Test API token"`
    - Verificar que actualiza `config.go`, `.env.example`, `validator.go`
    - Agregar valor en `.env` local
    - _Requirements: 5.1, 5.3, 5.4_

  - [ ] 10.3 Validar que la aplicación carga las nuevas variables
    - Ejecutar aplicación con `make run`
    - Verificar en logs o debugger que las variables se cargan correctamente
    - Verificar precedencia: cambiar valor en ENV y verificar que sobrescribe YAML
    - _Requirements: 2.4, 9.2, 9.5_

  - [ ] 10.4 Limpiar variables de prueba
    - Remover campos de prueba de `config.go`
    - Remover de YAMLs y `.env.example`
    - Remover validación de `validator.go`
    - Verificar que aplicación sigue funcionando
    - _Requirements: 9.5_

- [ ] 11. Revisión final y merge
  - [x] 11.1 Ejecutar auditoría de código
    - Ejecutar `make audit` (fmt, vet, tests)
    - Ejecutar `make lint` si está configurado
    - Corregir cualquier issue encontrado
    - _Requirements: 9.1, 9.5_

  - [x] 11.2 Revisar documentación
    - Verificar que CONFIG.md está completo
    - Verificar que README.md está actualizado
    - Verificar que `.env.example` está bien documentado
    - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5_

  - [x] 11.3 Crear PR y solicitar revisión
    - Crear Pull Request con descripción detallada de cambios
    - Incluir antes/después de la refactorización
    - Documentar breaking changes si los hay
    - Solicitar revisión del equipo
    - _Requirements: All_
