# 🧪 API Tests - httpyac

Sistema de testing HTTP para EduGo API usando **httpyac**. Este directorio contiene archivos `.http` que pueden ejecutarse desde la línea de comandos o desde VSCode.

---

## 📋 Propósito

Esta carpeta es **exclusivamente** para ejecutar endpoints del proyecto. **No** está integrada con CI/CD ni otros sistemas. Su función es facilitar el testing manual y local de la API.

---

## ⚡ Inicio Rápido

### 1. Instalar httpyac

```bash
# Instalar httpyac globalmente
npm install -g httpyac
```

### 2. Configurar Variables de Entorno

Crea un archivo `.env.local` basado en el ejemplo:

```bash
cp api-tests/.env.example api-tests/.env.local
```

Edita `.env.local` con tus credenciales:

```bash
BASE_URL=http://localhost:9090
EMAIL=tu-email@ejemplo.com
PASSWORD=tu-password
```

### 3. Ejecutar Tests

```bash
# Ejecutar todos los archivos .http
httpyac api-tests/requests --dotenv api-tests/.env.local

# Ejecutar un archivo específico
httpyac api-tests/requests/auth.http --dotenv api-tests/.env.local

# Ejecutar un request individual con nombre
httpyac api-tests/requests/auth.http --request login --dotenv api-tests/.env.local
```

---

## 📁 Estructura

```
api-tests/
├── README.md                    # Este archivo
├── .env.example                 # Plantilla de variables
├── .env.local                   # Variables locales (no versionado)
├── .gitignore                   # Ignora .env.local y datos sensibles
├── requests/                    # Archivos .http con peticiones
│   ├── auth.http               # Autenticación (login, refresh, logout)
│   ├── materials.http          # CRUD de materiales
│   ├── assessments.http        # Evaluaciones
│   ├── progress.http           # Progreso de estudiantes
│   ├── stats.http              # Estadísticas
│   └── health.http             # Health check
├── scripts/                     # Scripts de utilidad
│   ├── run-httpyac.sh          # Ejecuta httpyac con login automático
│   ├── list-requests           # Lista todas las peticiones disponibles
│   └── update-token.sh         # (Deprecado - usa httpyac directamente)
└── data/                        # Datos de prueba (ignorado en git)
```

---

## 🎯 Archivos HTTP Disponibles

| Archivo | Descripción | Requests |
|---------|-------------|----------|
| [auth.http](requests/auth.http) | Login, refresh, logout, revoke-all | 4 |
| [materials.http](requests/materials.http) | CRUD, S3, versiones, resumen, stats | 10 |
| [assessments.http](requests/assessments.http) | Evaluaciones | 3 |
| [progress.http](requests/progress.http) | Progreso de estudiantes | 2 |
| [stats.http](requests/stats.http) | Estadísticas | 2 |
| [health.http](requests/health.http) | Health check (sin auth) | 1 |

---

## 🔐 Manejo de Tokens

httpyac maneja los tokens automáticamente usando **post-request scripts** en JavaScript.

### Flujo de Autenticación

1. **Login**: Ejecuta `auth.http` → el script post-request exporta `access_token` y `refresh_token`
2. **Uso**: Los demás archivos usan `{{access_token}}` automáticamente
3. **Refresh**: Cuando expire, ejecuta la request de refresh

### Ejemplo: auth.http

```http
# @name login
POST {{baseUrl}}/v1/auth/login
Content-Type: application/json

{
  "email": "{{email}}",
  "password": "{{password}}"
}

{{
  // Post-request script: exporta tokens
  exports.access_token = response.parsedBody.access_token;
  exports.refresh_token = response.parsedBody.refresh_token;
}}
```

### Uso en Otros Archivos

```http
GET {{baseUrl}}/v1/materials
Authorization: Bearer {{access_token}}
```

---

## 🚀 Scripts de Utilidad

### `run-httpyac.sh`

Script mejorado que ejecuta httpyac con login automático:

```bash
# Ejecutar con login automático
./api-tests/scripts/run-httpyac.sh

# Solo hacer login (obtener tokens)
./api-tests/scripts/run-httpyac.sh --login-only

# Login y guardar tokens en archivo
./api-tests/scripts/run-httpyac.sh --persist-token

# Usar otro archivo .env
./api-tests/scripts/run-httpyac.sh --dotenv=api-tests/.env.dev
```

### `list-requests`

Lista todas las peticiones disponibles:

```bash
./api-tests/scripts/list-requests
```

---

## 🌍 Variables de Entorno

Las variables se cargan desde `.env.local`:

| Variable | Descripción | Ejemplo |
|----------|-------------|---------|
| `BASE_URL` | URL base del API | `http://localhost:9090` |
| `EMAIL` | Email de autenticación | `test@edugo.com` |
| `PASSWORD` | Password de autenticación | `Test123!` |

### Múltiples Ambientes

Puedes crear archivos `.env` para diferentes ambientes:

```bash
.env.local    # Local development
.env.dev      # Dev server
.env.qa       # QA server
.env.prod     # Production (¡cuidado!)
```

Luego especifica cuál usar:

```bash
httpyac api-tests/requests --dotenv api-tests/.env.dev
```

---

## 🔍 Características de httpyac

### 1. Post-Request Scripts

Ejecuta código JavaScript después de cada request:

```http
POST {{baseUrl}}/v1/materials
...

{{
  // Guardar el ID del material creado
  exports.material_id = response.parsedBody.data.id;
}}
```

### 2. Assertions

Valida respuestas con asserts:

```http
GET {{baseUrl}}/v1/health

{{
  test("Status is 200", () => {
    expect(response.statusCode).toBe(200);
  });

  test("Response has status field", () => {
    expect(response.parsedBody.status).toBe("ok");
  });
}}
```

### 3. Variables Dinámicas

Usa variables del sistema:

```http
POST {{baseUrl}}/v1/materials
Content-Type: application/json

{
  "title": "Material {{$timestamp}}",
  "created_at": "{{$datetime iso8601}}"
}
```

Variables disponibles:
- `{{$timestamp}}` - Unix timestamp
- `{{$datetime iso8601}}` - Fecha ISO 8601
- `{{$guid}}` - GUID/UUID aleatorio
- `{{$randomInt min max}}` - Número aleatorio

---

## 📊 Reporting

### Generar Reportes JUnit

```bash
httpyac api-tests/requests \
  --dotenv api-tests/.env.local \
  --report junit:reports/junit.xml
```

### Salida en JSON

```bash
httpyac api-tests/requests \
  --dotenv api-tests/.env.local \
  --output json > results.json
```

---

## 🛠️ Troubleshooting

### Error: "httpyac: command not found"

**Solución**: Instala httpyac globalmente

```bash
npm install -g httpyac
```

### Error: "baseUrl is not found"

**Solución**: Verifica que `.env.local` existe y contiene `BASE_URL`

```bash
cat api-tests/.env.local
```

### Tokens no se comparten entre archivos

**Causa**: Los tokens se exportan con `exports` en el post-request script de auth.http

**Solución**: Ejecuta primero `auth.http` para generar los tokens, luego ejecuta otros archivos en la misma sesión de httpyac

### Error 401 Unauthorized

**Solución**: El token expiró. Ejecuta nuevamente la request de login:

```bash
httpyac api-tests/requests/auth.http --request login --dotenv api-tests/.env.local
```

---

## 🎓 Uso con VSCode

Puedes instalar la extensión de httpyac para VSCode:

```bash
code --install-extension anweber.vscode-httpyac
```

**Características**:
- ✅ Ejecuta requests directamente desde el editor
- ✅ Autocompletado de variables
- ✅ Visualización de respuestas
- ✅ Debugging de scripts

---

## ⚠️ Importante

- **NO** versiones archivos `.env.local` con credenciales reales
- **NO** commities tokens o datos sensibles
- Esta carpeta es **solo** para testing manual
- **NO** integrar con CI/CD (usa tests unitarios/integración para eso)

---

## 📚 Documentación

- **httpyac**: https://httpyac.github.io/
- **Formato .http**: https://httpyac.github.io/guide/request.html
- **Scripting**: https://httpyac.github.io/guide/scripting.html

---

**Última actualización**: 11 de noviembre de 2025  
**Versión**: 3.0 (httpyac)  
**Responsable**: Claude Code + Jhoan Medina
