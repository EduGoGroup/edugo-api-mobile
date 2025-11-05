---
description: Analizar cobertura de tests del código
argument-hint: "[archivo|directorio|--recent]"
---

# Comando: Análisis de Cobertura de Tests

## 🎯 Rol: ORQUESTADOR

**Este comando NO ejecuta el análisis. Delega al agente `test-coverage-analyzer`.**

Tu función:
1. Identificar qué código analizar (archivo, directorio o cambios recientes)
2. Recopilar contexto del código a analizar
3. Invocar al agente especializado usando **Task tool**
4. Retornar resultado al usuario

---

## Sintaxis

```bash
/test-coverage                           # Analiza cambios recientes (git diff)
/test-coverage archivo.go                # Analiza archivo específico
/test-coverage internal/application/     # Analiza directorio completo
/test-coverage --recent                  # Analiza últimos cambios (git diff)
```

---

## Ejecución

### 1. Determinar Alcance del Análisis

**Si el usuario proporciona un argumento:**
- Si es archivo: Analizar ese archivo específico
- Si es directorio: Analizar todos los archivos `.go` en ese directorio
- Si es `--recent`: Analizar archivos modificados (git diff)

**Si no hay argumentos:**
- Por defecto, analizar archivos modificados recientemente (git diff)

### 2. Recopilar Contexto

**Para archivo específico:**
```bash
# Verificar que existe
ls [archivo]

# Leer contenido
Read [archivo]

# Buscar test correspondiente
Read [archivo_test.go] (si existe)
```

**Para directorio:**
```bash
# Listar archivos .go (excluir *_test.go)
Glob [directorio]/**/*.go

# Leer archivos principales (max 10 más recientes)
Read [archivos encontrados]
```

**Para cambios recientes:**
```bash
# Ver archivos modificados
git diff --name-only HEAD

# Filtrar solo archivos .go (excluir *_test.go)
# Leer contenido de archivos modificados
```

### 3. Invocar Agente test-coverage-analyzer

**USA TASK TOOL:**

```
Task(
  subagent_type: "test-coverage-analyzer",
  description: "Analizar cobertura de tests",
  prompt: "
    Analiza la cobertura de tests del siguiente código Go.

    ALCANCE: [archivo específico | directorio | cambios recientes]

    ARCHIVOS A ANALIZAR:
    [Lista de archivos con su contenido]

    TESTS EXISTENTES:
    [Lista de archivos *_test.go encontrados con su contenido]

    CONTEXTO DEL PROYECTO:
    - Arquitectura: Clean Architecture (Hexagonal)
    - Framework web: Gin
    - Testing: testcontainers para integración
    - Base de datos: PostgreSQL, MongoDB
    - Ubicación handlers actuales: internal/infrastructure/http/handler/

    ENTREGA:
    Para cada archivo analizado, proporciona:
    1. Resumen de cobertura actual
    2. Análisis detallado de funciones
    3. Tests faltantes identificados
    4. Plan de acción priorizado
    5. Recomendaciones específicas

    Sigue la estructura definida en tu prompt del agente.
  "
)
```

### 4. Confirmar al Usuario

```
✅ Análisis de cobertura completado

📁 Archivos analizados: [N archivos]

📊 Resumen:
├─ Cobertura promedio: XX%
├─ Tests existentes: Y
├─ Tests recomendados: Z
└─ Estado general: ✅ Bueno | ⚠️ Mejorable | ❌ Insuficiente

📋 Tests prioritarios a implementar: [top 3-5]

💡 Ver análisis completo arriba para detalles y recomendaciones.

📌 Siguiente:
- Implementar tests de alta prioridad
- Revisar tests existentes sugeridos para refactorización
- Volver a ejecutar después de implementar tests
```

---

## 🚨 Manejo de Errores

### Error: Archivo no encontrado
```
❌ Error: Archivo no encontrado
Archivo: [ruta]
Verifica la ruta e intenta nuevamente
```

### Error: No hay cambios recientes
```
ℹ️ No hay cambios recientes en archivos Go
├─ Última modificación: [fecha del último commit]
└─ Sugerencia: Especifica un archivo o directorio

Ejemplos:
- /test-coverage internal/application/services/
- /test-coverage internal/infrastructure/http/handler/auth_handler.go
```

### Error: Directorio vacío
```
⚠️ No se encontraron archivos Go en el directorio
Directorio: [ruta]
Verifica que el directorio contiene archivos .go
```

---

## 💡 Tips de Uso

**Después de implementar una feature:**
```bash
git add .
/test-coverage --recent
```

**Analizar un módulo completo:**
```bash
/test-coverage internal/domain/
```

**Verificar un handler específico:**
```bash
/test-coverage internal/infrastructure/http/handler/auth_handler.go
```

**Análisis antes de hacer commit:**
```bash
/test-coverage --recent
# Revisar recomendaciones
# Implementar tests críticos
# Commit con confianza
```

---

## 🎯 Objetivo

Asegurar que todo código nuevo o modificado tenga la cobertura de tests adecuada según las mejores prácticas de Clean Architecture y Go testing, antes de hacer commit.
