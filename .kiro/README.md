# Kiro IDE Configuration

## Debug Configurations

El archivo `launch.json` contiene **3 configuraciones** para debuggear en Kiro:

### 1. 🎯 Launch API (with .env) - **RECOMENDADO**

```json
{
    "name": "Launch API (with .env)",
    "envFile": "${workspaceFolder}/.env"
}
```

**Intenta esta primero**. Si Kiro soporta `envFile`, cargará automáticamente todas las variables del archivo `.env`.

### 2. 🔧 Launch API (explicit vars - backup)

```json
{
    "name": "Launch API (explicit vars - backup)",
    "env": { /* todas las variables explícitas */ }
}
```

**Usa esta si la opción 1 no funciona**. Tiene todas las variables definidas explícitamente.

### 3. 📜 Launch API (with script)

```json
{
    "name": "Launch API (with script)",
    "program": "${workspaceFolder}/scripts/load-env.sh"
}
```

**Alternativa usando script**. Usa `scripts/load-env.sh` que carga `.env` antes de ejecutar.

---

## 🚀 Cómo Usar

### Opción A: Probar envFile (Más Simple)

1. Abre el panel de Debug en Kiro
2. Selecciona **"Launch API (with .env)"**
3. Presiona F5
4. ✅ Si funciona, ¡perfecto! Todas las variables se cargan de `.env`
5. ❌ Si no funciona, prueba la Opción B

### Opción B: Usar Variables Explícitas

1. Selecciona **"Launch API (explicit vars - backup)"**
2. Presiona F5
3. ✅ Debería funcionar siempre

### Opción C: Usar Script Wrapper

1. Selecciona **"Launch API (with script)"**
2. Presiona F5
3. El script `load-env.sh` carga `.env` automáticamente

---

## 🔄 Actualizar Variables

### Si usas envFile (Opción A):
```bash
# Solo edita .env
nano .env
# Los cambios se aplican automáticamente
```

### Si usas variables explícitas (Opción B):
```bash
# Edita .kiro/launch.json
# Actualiza el bloque "env"
```

### Si usas el script (Opción C):
```bash
# Solo edita .env
nano .env
# El script lo carga automáticamente
```

---

## 📊 Comparación de Opciones

| Opción | Ventaja | Desventaja |
|--------|---------|------------|
| **envFile** | Más simple, editas solo .env | Puede no funcionar en Kiro |
| **Explicit vars** | Siempre funciona | Hay que duplicar variables |
| **Script** | Flexible, usa .env | Requiere script adicional |

---

## 🐛 Troubleshooting

### Error: "envFile not supported"
- Usa la configuración **"Launch API (explicit vars - backup)"**
- O usa **"Launch API (with script)"**

### Error: "Configuration validation failed"
- Verifica que `.env` tenga todas las variables requeridas
- Compara con `.env.example`

### Error: "connection refused"
- Asegúrate de tener PostgreSQL, MongoDB y RabbitMQ corriendo
- O usa Docker Compose: `docker-compose up`

### Quiero cambiar valores
- **Con envFile**: Edita `.env` y reinicia debug
- **Con explicit vars**: Edita `.kiro/launch.json`
- **Con script**: Edita `.env` y reinicia debug

---

## 💡 Recomendación

1. **Primero intenta**: "Launch API (with .env)"
2. **Si no funciona**: "Launch API (explicit vars - backup)"
3. **Alternativa**: Usa terminal con `make run` (siempre funciona)

---

## 📚 Más Información

- Ver [QUICKSTART.md](../QUICKSTART.md) para guía completa
- Ver [CONFIG.md](../CONFIG.md) para documentación de configuración
- Ver [.env.example](../.env.example) para todas las variables disponibles

---

## 🎯 Zed IDE (Comparación)

Zed **SÍ soporta** `envFile` nativamente:

```json
{
    "envFile": "${workspaceFolder}/.env"
}
```

Si Kiro no lo soporta, considera usar Zed o VSCode para debug, o usa `make run` desde terminal.
