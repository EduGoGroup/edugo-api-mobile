# Flujo de Procesos

## Descripción General

El sistema implementa un flujo de procesamiento HTTP simple basado en el patrón Request-Response. Cada petición del cliente pasa por el router de Express, se procesa según la ruta solicitada, y retorna una respuesta JSON.

El flujo es **stateless** (sin estado), lo que significa que cada petición es independiente y no requiere información de peticiones anteriores. No hay sesiones, autenticación, ni persistencia de datos.

## Proceso Principal: Manejo de Peticiones HTTP

```mermaid
flowchart TD
    Start([Cliente envía petición HTTP])
    Server[Servidor Express recibe petición]
    Router{Analizar ruta}

    HelloBasic["/api/hello"]
    HelloName["/api/hello/:name"]
    Status["/api/status"]
    NotFound[Ruta no encontrada]

    ProcessBasic[Generar mensaje: 'Hello, World!']
    ProcessName[Extraer parámetro 'name']
    ProcessStatus[Generar timestamp]

    BuildName[Construir mensaje: 'Hello, name!']
    BuildStatus[Construir objeto de estado]

    ResponseBasic[Enviar JSON con mensaje básico]
    ResponseName[Enviar JSON con mensaje personalizado]
    ResponseStatus[Enviar JSON con estado]
    Response404[Enviar error 404]

    End([Cliente recibe respuesta])

    Start --> Server
    Server --> Router

    Router -->|GET /api/hello| HelloBasic
    Router -->|GET /api/hello/:name| HelloName
    Router -->|GET /api/status| Status
    Router -->|Otra ruta| NotFound

    HelloBasic --> ProcessBasic
    ProcessBasic --> ResponseBasic
    ResponseBasic --> End

    HelloName --> ProcessName
    ProcessName --> BuildName
    BuildName --> ResponseName
    ResponseName --> End

    Status --> ProcessStatus
    ProcessStatus --> BuildStatus
    BuildStatus --> ResponseStatus
    ResponseStatus --> End

    NotFound --> Response404
    Response404 --> End

    style Start fill:#e1f5ff
    style End fill:#e1f5ff
    style Server fill:#fff4e1
    style Router fill:#ffe1e1
    style ResponseBasic fill:#e1ffe1
    style ResponseName fill:#e1ffe1
    style ResponseStatus fill:#e1ffe1
    style Response404 fill:#ffe1e1
```

## Descripción del Flujo Principal

### 1. **Recepción de Petición**
- El cliente (curl, navegador, Postman) envía una petición HTTP GET
- El servidor Express escuchando en el puerto 3000 recibe la petición
- El middleware de parsing procesa la petición

### 2. **Análisis de Ruta (Routing)**
- Express Router examina la URL de la petición
- Compara con las rutas registradas: `/api/hello`, `/api/hello/:name`, `/api/status`
- Determina qué handler debe procesar la petición

### 3. **Procesamiento según Endpoint**

#### Caso A: GET /api/hello
- No requiere procesamiento de parámetros
- Genera objeto JSON con mensaje estático: `{ "message": "Hello, World!" }`
- Envía respuesta con código HTTP 200

#### Caso B: GET /api/hello/:name
- Extrae el parámetro `name` de la URL usando `req.params.name`
- Valida que el parámetro existe (Express garantiza esto si la ruta coincide)
- Construye mensaje personalizado: `"Hello, " + name + "!"`
- Genera objeto JSON: `{ "message": "Hello, {name}!" }`
- Envía respuesta con código HTTP 200

#### Caso C: GET /api/status
- Genera timestamp actual usando `new Date().toISOString()`
- Construye objeto JSON con estado: `{ "status": "ok", "timestamp": "..." }`
- Envía respuesta con código HTTP 200

#### Caso D: Ruta No Encontrada
- Si la ruta no coincide con ninguna registrada
- Express retorna error 404
- Opcionalmente, el middleware de errores formatea la respuesta: `{ "error": "Not Found" }`

### 4. **Envío de Respuesta**
- El servidor serializa el objeto JavaScript a JSON
- Agrega headers apropiados: `Content-Type: application/json`
- Envía la respuesta al cliente
- Cierra la conexión HTTP

### 5. **Recepción por el Cliente**
- El cliente recibe el JSON
- Procesa la respuesta según su lógica
- La transacción finaliza

## Flujo de Inicio del Servidor

```mermaid
flowchart TD
    Init([Ejecutar: npm start])
    LoadDeps[Cargar dependencias: Express, etc.]
    CreateApp[Crear instancia de Express]
    ConfigMiddleware[Configurar middleware: express.json]
    RegisterRoutes[Registrar rutas desde routes/hello.js]
    StartServer[Iniciar servidor en puerto 3000]
    Listen[Escuchar peticiones]
    Ready([Servidor listo])

    Init --> LoadDeps
    LoadDeps --> CreateApp
    CreateApp --> ConfigMiddleware
    ConfigMiddleware --> RegisterRoutes
    RegisterRoutes --> StartServer
    StartServer --> Listen
    Listen --> Ready

    style Init fill:#e1f5ff
    style Ready fill:#e1ffe1
    style StartServer fill:#fff4e1
```

### Descripción del Flujo de Inicio

1. **Ejecución del comando**: `node src/index.js` o `npm start`
2. **Carga de módulos**: Node.js importa Express y módulos custom
3. **Inicialización de Express**: Se crea la aplicación Express
4. **Configuración de middleware**: Se agrega `express.json()` para parsear JSON
5. **Registro de rutas**: Se importan y registran las rutas desde `routes/hello.js`
6. **Binding del puerto**: El servidor se enlaza al puerto 3000
7. **Estado listo**: El servidor está escuchando peticiones

**Tiempo esperado**: < 3 segundos (según requerimientos)

## Flujo de Manejo de Errores

```mermaid
flowchart TD
    Request[Petición entrante]
    Process{Procesamiento}
    Success[Ejecución exitosa]
    Error[Error ocurrido]

    ErrorType{Tipo de error}
    E404[404 - Ruta no encontrada]
    E500[500 - Error del servidor]

    Format404[Formatear respuesta 404]
    Format500[Formatear respuesta 500]
    Log[Registrar error en consola]

    Send404[Enviar JSON con error 404]
    Send500[Enviar JSON con error 500]

    End([Cliente recibe error])

    Request --> Process
    Process -->|OK| Success
    Process -->|Fallo| Error

    Error --> ErrorType
    ErrorType -->|Ruta inexistente| E404
    ErrorType -->|Excepción| E500

    E404 --> Format404
    Format404 --> Send404
    Send404 --> End

    E500 --> Log
    Log --> Format500
    Format500 --> Send500
    Send500 --> End

    style Error fill:#ffe1e1
    style E404 fill:#ffe1e1
    style E500 fill:#ffe1e1
    style Log fill:#fff4e1
```

### Tipos de Errores Manejados

#### Error 404 - Not Found
- **Causa**: El cliente solicita una ruta que no existe (ej: `/api/xyz`)
- **Respuesta**:
  ```json
  {
    "error": "Not Found",
    "message": "La ruta solicitada no existe"
  }
  ```
- **Código HTTP**: 404

#### Error 500 - Internal Server Error
- **Causa**: Excepción no capturada durante el procesamiento
- **Acciones**:
  1. Registrar el error en consola para debugging
  2. Evitar exponer stack trace al cliente
  3. Retornar mensaje genérico de error
- **Respuesta**:
  ```json
  {
    "error": "Internal Server Error",
    "message": "Ocurrió un error en el servidor"
  }
  ```
- **Código HTTP**: 500

## Puntos de Integración

### Cliente HTTP → Servidor Express
- **Protocolo**: HTTP/1.1
- **Formato**: JSON
- **Métodos permitidos**: GET (únicamente)
- **Puerto**: 3000
- **Host**: localhost (desarrollo)

### Express → Routes Module
- **Mecanismo**: `require()` y `app.use()`
- **Acoplamiento**: Bajo (módulos separados)
- **Extensibilidad**: Fácil agregar nuevas rutas sin modificar index.js

### Routes → Response Builder
- **Mecanismo**: `res.json()` de Express
- **Serialización**: Automática (JavaScript Object → JSON)
- **Headers**: Gestionados por Express automáticamente

## Diagrama de Secuencia: Petición Completa

```mermaid
sequenceDiagram
    actor Cliente
    participant Express
    participant Router
    participant Handler

    Cliente->>Express: GET /api/hello/Juan
    Express->>Router: Buscar ruta coincidente
    Router->>Handler: Ejecutar handler de /api/hello/:name
    Handler->>Handler: Extraer params.name = "Juan"
    Handler->>Handler: Construir mensaje = "Hello, Juan!"
    Handler->>Express: res.json({ message: "Hello, Juan!" })
    Express->>Cliente: 200 OK + JSON Response
```

## Optimizaciones de Performance

### Decisiones que mejoran el rendimiento:
1. **Sin I/O de disco**: Todas las respuestas se generan en memoria
2. **Sin consultas a BD**: Elimina latencia de red/disco
3. **Procesamiento mínimo**: Solo concatenación de strings y generación de timestamps
4. **Sin middleware pesado**: Solo JSON parser esencial
5. **Respuestas pequeñas**: JSON mínimo (< 100 bytes)

**Resultado esperado**: Tiempo de respuesta < 100ms (según requerimientos)

## Casos de Uso en Producción

Aunque este es un ejercicio de validación, el flujo es representativo de:
- APIs de salud/status en sistemas distribuidos
- Endpoints de prueba en entornos de desarrollo
- Servicios de echo/ping para monitoreo de red
- Prototipos rápidos de APIs REST

---

**Fecha de análisis**: 2025-10-31
**Versión de documento**: 1.0
**Complejidad del flujo**: Baja (ideal para validación)
