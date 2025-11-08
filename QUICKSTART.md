# 🚀 Quick Start - EduGo API Mobile

## Configuración Rápida (5 minutos)

### Opción 1: Desarrollo Local (sin Docker)

**Requisitos**: PostgreSQL, MongoDB, RabbitMQ corriendo en tu máquina

```bash
# 1. El archivo .env ya está configurado para localhost
cat .env  # Verifica que las URIs usen 'localhost'

# 2. Asegúrate de tener los servicios corriendo
# PostgreSQL en puerto 5432
# MongoDB en puerto 27017
# RabbitMQ en puerto 5672

# 3. Ejecutar la aplicación
make run
# o
go run cmd/main.go
```

**Configuración de IDEs**:
- ✅ **VSCode**: Ya configurado en `.vscode/launch.json`
- ✅ **Zed**: Ya configurado en `.zed/debug.json`
- ✅ **Kiro**: Ya configurado en `.kiro/launch.json`
- ✅ **IntelliJ/GoLand**: Ver `.idea/runConfigurations/README.md`

Todos cargan automáticamente el archivo `.env`

---

### Opción 2: Docker Compose (Recomendado)

**Requisitos**: Solo Docker Desktop

```bash
# 1. Usar el archivo .env.docker (tiene hosts correctos para Docker)
cp .env.docker .env

# 2. Levantar todos los servicios
docker-compose up

# La API estará en: http://localhost:9090
# Swagger en: http://localhost:9090/swagger/index.html
```

**Para detener**:
```bash
docker-compose down
```

---

### Opción 3: Solo la App en Docker, Servicios Locales

```bash
# 1. Asegúrate de que .env use 'localhost'
# 2. Tener PostgreSQL, MongoDB, RabbitMQ corriendo localmente
# 3. Construir y correr solo la API

docker build -t edugo-api-mobile .
docker run --env-file .env -p 9090:8080 edugo-api-mobile
```

---

## 🧪 Ejecutar Tests

### Tests Unitarios
```bash
make test
# o
go test ./... -short
```

### Tests de Integración (con Testcontainers)
```bash
# Los testcontainers crean sus propios contenedores temporales
# NO necesitan .env ni servicios corriendo
make test-integration

# O manualmente:
RUN_INTEGRATION_TESTS=true go test -tags=integration ./test/integration/... -v
```

### Testing con Bootstrap

El sistema de bootstrap facilita el testing con inyección de mocks:

```go
func TestMyFeature(t *testing.T) {
    cfg := testConfig()
    
    // Inyectar mocks
    b := bootstrap.New(cfg,
        bootstrap.WithLogger(mockLogger),
        bootstrap.WithPostgreSQL(mockDB),
        bootstrap.WithMongoDB(mockMongoDB),
        bootstrap.WithRabbitMQ(mockPublisher),
        bootstrap.WithS3Client(mockS3),
    )
    
    resources, cleanup, err := b.InitializeInfrastructure(context.Background())
    require.NoError(t, err)
    defer cleanup()
    
    // Usar resources en tus tests
}
```

**📖 Guía completa**: [internal/bootstrap/INTEGRATION_TESTS.md](internal/bootstrap/INTEGRATION_TESTS.md)

---

## 🔧 Validar Configuración

```bash
# Validar que todos los archivos de configuración sean válidos
make config-validate

# Ver qué variables se están cargando (sin valores sensibles)
go run cmd/main.go --help
```

---

## 🧩 Desarrollo sin Infraestructura Completa

### Recursos Opcionales

Si no tienes RabbitMQ o S3 disponibles, puedes marcarlos como opcionales:

**Opción 1: Configuración YAML**

Edita `config/config-local.yaml`:
```yaml
infrastructure:
  optional_resources:
    - rabbitmq
    - s3
```

**Opción 2: Variables de Entorno**

Agrega a tu `.env`:
```bash
INFRASTRUCTURE_OPTIONAL_RESOURCES=rabbitmq,s3
```

La aplicación iniciará con implementaciones noop para estos recursos. Las funcionalidades que dependen de ellos registrarán advertencias pero no fallarán.

### ¿Qué Recursos Puedo Hacer Opcionales?

- ✅ **RabbitMQ**: Eventos no se publicarán (solo logs)
- ✅ **S3**: URLs de archivos no estarán disponibles
- ❌ **PostgreSQL**: Siempre requerido (datos principales)
- ❌ **MongoDB**: Siempre requerido (evaluaciones y resúmenes)

**📖 Más información**: 
- [config/OPTIONAL_RESOURCES.md](config/OPTIONAL_RESOURCES.md) - Configuración de recursos opcionales
- [docs/BOOTSTRAP_USAGE.md](docs/BOOTSTRAP_USAGE.md) - Guía completa del sistema de bootstrap

---

## 📝 Archivos de Configuración

```
.env              ← Tu configuración local (localhost)
.env.docker       ← Configuración para Docker Compose
.env.example      ← Template con documentación completa

config/
├── config.yaml         ← Base (todos los ambientes)
├── config-local.yaml   ← Local (puerto 9090, logs debug)
├── config-dev.yaml     ← Development server
├── config-qa.yaml      ← QA/Staging
└── config-prod.yaml    ← Production
```

---

## 🐛 Troubleshooting

### Error: "Configuration validation failed"
```bash
# Verifica que todas las variables requeridas estén en .env
cat .env

# Compara con .env.example para ver qué falta
diff .env .env.example
```

### Error: "connection refused" (PostgreSQL/MongoDB/RabbitMQ)
```bash
# Opción A: Verifica que los servicios estén corriendo
docker ps  # Si usas Docker
# o
lsof -i :5432  # PostgreSQL
lsof -i :27017 # MongoDB
lsof -i :5672  # RabbitMQ

# Opción B: Usa Docker Compose (más fácil)
cp .env.docker .env
docker-compose up
```

### La app no carga el .env
```bash
# Verifica que el archivo existe
ls -la .env

# Verifica que tu IDE esté configurado
# VSCode: .vscode/launch.json debe tener "envFile"
# Zed: .zed/debug.json debe tener "envFile"

# Para Make, debería funcionar automáticamente
make run
```

---

## 📚 Documentación Completa

- **[CONFIG.md](CONFIG.md)** - Guía completa de configuración
- **[README.md](README.md)** - Documentación general del proyecto
- **[.env.example](.env.example)** - Todas las variables disponibles

---

## 🎯 Comandos Útiles

```bash
# Desarrollo
make run              # Ejecutar la aplicación
make build            # Compilar binario
make test             # Ejecutar tests
make config-validate  # Validar configuración

# Docker
docker-compose up     # Levantar todo
docker-compose down   # Detener todo
docker-compose logs   # Ver logs

# Limpieza
make clean            # Limpiar binarios y cache
```

---

## ✅ Checklist de Setup

- [ ] Archivo `.env` existe y tiene valores correctos
- [ ] Servicios corriendo (PostgreSQL, MongoDB, RabbitMQ)
- [ ] `make test` pasa sin errores
- [ ] `make config-validate` pasa sin errores
- [ ] La aplicación inicia sin errores
- [ ] Swagger accesible en http://localhost:9090/swagger/index.html

---

**¿Problemas?** Revisa [CONFIG.md](CONFIG.md) o pregunta al equipo.
