# IDE Setup Guide - EduGo API Mobile

## 📋 Resumen de Soporte de .env por IDE

| IDE | Soporte envFile | Configuración | Recomendación |
|-----|----------------|---------------|---------------|
| **VSCode** | ✅ Nativo | `.vscode/launch.json` | ⭐ Recomendado |
| **Zed** | ✅ Nativo | `.zed/debug.json` | ⭐ Recomendado |
| **Kiro** | ⚠️ Probar | `.kiro/launch.json` | Probar envFile primero |
| **IntelliJ/GoLand** | ✅ Plugin | EnvFile plugin | Requiere plugin |
| **Terminal** | ✅ Nativo | Makefile | ⭐ Siempre funciona |

---

## 🎯 VSCode (Recomendado)

### ✅ Soporte Nativo de envFile

**Archivo**: `.vscode/launch.json`

```json
{
    "name": "Launch API",
    "type": "go",
    "request": "launch",
    "program": "${workspaceFolder}/cmd/main.go",
    "envFile": "${workspaceFolder}/.env"
}
```

### Cómo Usar:
1. Abre VSCode
2. Presiona F5 o ve a Run & Debug
3. Selecciona "Launch API"
4. ✅ Carga `.env` automáticamente

### Actualizar Variables:
```bash
# Solo edita .env
code .env
# Reinicia el debug (F5)
```

---

## 🎯 Zed (Recomendado)

### ✅ Soporte Nativo de envFile

**Archivo**: `.zed/debug.json`

```json
{
    "label": "Go: Debug main (Delve)",
    "adapter": "Delve",
    "program": "./cmd",
    "envFile": "${workspaceFolder}/.env"
}
```

### Cómo Usar:
1. Abre Zed
2. Panel de Debug
3. Selecciona "Go: Debug main (Delve)"
4. ✅ Carga `.env` automáticamente

### Actualizar Variables:
```bash
# Solo edita .env
zed .env
# Reinicia el debug
```

---

## ⚠️ Kiro (Probar)

### Configuración con 3 Opciones

**Archivo**: `.kiro/launch.json`

#### Opción 1: envFile (Probar Primero)
```json
{
    "name": "Launch API (with .env)",
    "envFile": "${workspaceFolder}/.env"
}
```

#### Opción 2: Variables Explícitas (Backup)
```json
{
    "name": "Launch API (explicit vars - backup)",
    "env": {
        "DATABASE_POSTGRES_PASSWORD": "edugo123",
        // ... todas las variables
    }
}
```

#### Opción 3: Script Wrapper
```json
{
    "name": "Launch API (with script)",
    "program": "${workspaceFolder}/scripts/load-env.sh"
}
```

### Cómo Usar:
1. Abre Kiro
2. Panel de Debug
3. **Primero intenta**: "Launch API (with .env)"
4. **Si no funciona**: "Launch API (explicit vars - backup)"

### Actualizar Variables:
- **Con envFile**: Edita `.env`
- **Con explicit vars**: Edita `.kiro/launch.json`

Ver [.kiro/README.md](../.kiro/README.md) para más detalles.

---

## 🔧 IntelliJ IDEA / GoLand

### ✅ Soporte con Plugin EnvFile

#### Paso 1: Instalar Plugin
1. `Settings/Preferences` → `Plugins`
2. Buscar "EnvFile"
3. Instalar y reiniciar

#### Paso 2: Configurar Run Configuration
1. `Run` → `Edit Configurations`
2. Selecciona tu configuración Go
3. Tab `EnvFile`
4. Click `+` → Selecciona `.env`
5. Check "Enable EnvFile"

### Alternativa Sin Plugin:
Edita manualmente las variables en `Run Configuration` → `Environment variables`

Ver [.idea/runConfigurations/README.md](../.idea/runConfigurations/README.md) para más detalles.

---

## ⭐ Terminal / Make (Siempre Funciona)

### ✅ Soporte Nativo

**Archivo**: `Makefile`

```makefile
# Carga .env automáticamente
ifneq (,$(wildcard .env))
    include .env
    export
endif
```

### Cómo Usar:
```bash
# Ejecutar aplicación
make run

# Ejecutar tests
make test

# Build
make build
```

### Ventajas:
- ✅ Siempre funciona
- ✅ No requiere configuración de IDE
- ✅ Consistente en todos los entornos
- ✅ Ideal para CI/CD

---

## 🔄 Flujo de Trabajo Recomendado

### Para Desarrollo Diario:

1. **Usa VSCode o Zed** (soporte nativo de envFile)
2. **Edita `.env`** cuando necesites cambiar variables
3. **Reinicia debug** para aplicar cambios

### Para Troubleshooting:

1. **Usa terminal**: `make run`
2. Verifica que `.env` tenga todas las variables
3. Compara con `.env.example`

### Para CI/CD:

1. **No uses `.env`**
2. Usa secrets del CI/CD (GitHub Secrets, GitLab Variables)
3. Ver [CONFIG.md](../CONFIG.md) para más detalles

---

## 📝 Archivo .env

### Estructura:
```bash
# Database
DATABASE_POSTGRES_PASSWORD=edugo123
DATABASE_MONGODB_URI=mongodb://...

# Messaging
MESSAGING_RABBITMQ_URL=amqp://...

# Storage
STORAGE_S3_ACCESS_KEY_ID=...
STORAGE_S3_SECRET_ACCESS_KEY=...
STORAGE_S3_BUCKET_NAME=edugo-materials

# Auth
JWT_SECRET=dev-secret-key

# App
APP_ENV=local
```

### Ubicación:
- ✅ Raíz del proyecto: `.env`
- ✅ Gitignored (no se commitea)
- ✅ Copia de `.env.example`

---

## 🐛 Troubleshooting

### IDE no carga .env

**VSCode/Zed**:
- Verifica que `envFile` esté en `launch.json`/`debug.json`
- Reinicia el IDE

**Kiro**:
- Prueba las 3 opciones en orden
- Usa "explicit vars" como backup

**IntelliJ**:
- Instala plugin EnvFile
- O configura variables manualmente

### Variables no se actualizan

1. Edita `.env`
2. **Reinicia el debug** (no solo reload)
3. O reinicia el IDE

### Sigue sin funcionar

```bash
# Usa terminal (siempre funciona)
make run
```

---

## 💡 Recomendaciones Finales

### Para Máxima Simplicidad:
1. **Usa VSCode o Zed** (envFile nativo)
2. **Edita solo `.env`**
3. **Nunca edites launch.json**

### Para Máxima Compatibilidad:
1. **Usa terminal**: `make run`
2. **Funciona en cualquier IDE**
3. **Consistente siempre**

### Para Equipos:
1. **Documenta qué IDE usas**
2. **Comparte `.env.example`**
3. **Nunca commitees `.env`**

---

## 📚 Más Información

- [QUICKSTART.md](../QUICKSTART.md) - Guía rápida de setup
- [CONFIG.md](../CONFIG.md) - Documentación completa de configuración
- [.env.example](../.env.example) - Template de variables
- [.kiro/README.md](../.kiro/README.md) - Específico para Kiro
- [.idea/runConfigurations/README.md](../.idea/runConfigurations/README.md) - Específico para IntelliJ
