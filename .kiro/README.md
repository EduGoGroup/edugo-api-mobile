# Kiro IDE Configuration

## Debug Configuration

El archivo `launch.json` contiene las configuraciones de debug para Kiro IDE.

### ⚠️ Nota Importante sobre Variables de Entorno

Kiro IDE actualmente no soporta el campo `envFile` como VSCode. Por lo tanto, las variables de entorno están definidas **explícitamente** en cada configuración de debug.

### Configuraciones Disponibles

1. **Launch API** - Ejecuta la aplicación completa
2. **Debug Current File** - Debuggea el archivo actual
3. **Debug Tests** - Ejecuta tests en modo debug

### 🔧 Actualizar Variables

Si necesitas cambiar las variables de entorno (por ejemplo, cambiar passwords o URIs):

**Opción 1: Editar launch.json directamente**
```json
{
    "env": {
        "DATABASE_POSTGRES_PASSWORD": "tu-nuevo-password",
        "DATABASE_MONGODB_URI": "mongodb://...",
        // ... otras variables
    }
}
```

**Opción 2: Usar el archivo .env con un script wrapper**

Crea un script `run-with-env.sh`:
```bash
#!/bin/bash
set -a
source .env
set +a
exec "$@"
```

Luego modifica `launch.json`:
```json
{
    "program": "./run-with-env.sh go run ${workspaceFolder}/cmd/main.go"
}
```

### 📝 Variables Actuales

Las variables están sincronizadas con `.env`:

```bash
APP_ENV=local
DATABASE_POSTGRES_PASSWORD=edugo123
DATABASE_MONGODB_URI=mongodb://edugo:edugo123@localhost:27017/edugo?authSource=admin
MESSAGING_RABBITMQ_URL=amqp://edugo:edugo123@localhost:5672/
STORAGE_S3_ACCESS_KEY_ID=test-access-key-id
STORAGE_S3_SECRET_ACCESS_KEY=test-secret-access-key
```

### 🚀 Cómo Usar

1. Abre el panel de Debug en Kiro
2. Selecciona "Launch API" en el dropdown
3. Presiona el botón de Play o F5
4. La aplicación iniciará con las variables configuradas

### 🐛 Troubleshooting

**Error: "Configuration validation failed"**
- Verifica que todas las variables en `launch.json` estén correctas
- Compara con el archivo `.env` para asegurar que estén sincronizadas

**Error: "connection refused"**
- Asegúrate de tener PostgreSQL, MongoDB y RabbitMQ corriendo en localhost
- O usa Docker Compose: `docker-compose up`

**Quiero usar diferentes valores**
- Edita el bloque `env` en `.kiro/launch.json`
- O corre desde terminal: `make run` (carga `.env` automáticamente)

### 💡 Alternativa: Usar Terminal

Si prefieres no editar `launch.json`, puedes correr desde terminal:

```bash
# Carga .env automáticamente
make run

# O directamente
go run cmd/main.go
```

El Makefile está configurado para cargar `.env` automáticamente.

### 📚 Más Información

- Ver [QUICKSTART.md](../QUICKSTART.md) para guía completa
- Ver [CONFIG.md](../CONFIG.md) para documentación de configuración
- Ver [.env.example](../.env.example) para todas las variables disponibles
