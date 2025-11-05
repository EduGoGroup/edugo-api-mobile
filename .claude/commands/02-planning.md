---
description: Transform architectural analysis into granular task plan
argument-hint: ""
---

# Comando: Planificación de Sprint

## 🎯 Rol: ORQUESTADOR

**Este comando NO ejecuta la planificación. Delega al agente `flow-planner`.**

Tu función:
1. Validar que existe el análisis
2. Invocar al agente especializado usando **Task tool**
3. Retornar resultado al usuario

---

## Ejecución

### 1. Validar Entrada

Verificar que existe: `sprint/current/analysis/readme.md`

Si no existe:
```
❌ Error: Análisis de sprint no encontrado
Por favor ejecuta primero: /01-analysis
```

### 2. Invocar Agente flow-planner

**USA TASK TOOL:**

```
Task(
  subagent_type: "flow-planner",
  description: "Planificación de sprint",
  prompt: "
    Transforma el análisis arquitectónico en plan de trabajo granular.

    ENTRADA: sprint/current/analysis/readme.md
    SALIDA: sprint/current/planning/readme.md

    Genera un plan con:
    - Fases organizadas con casillas de verificación
    - Tareas atómicas y granulares
    - Indicadores de dependencia entre tareas
    - Listo para ejecución modular
  "
)
```

### 3. Confirmar al Usuario

```
✅ Planificación completada

📁 Archivo: sprint/current/planning/readme.md
📌 Siguiente: /03-execution (todo) o /03-execution phase-N (específico)
```

---

## 🚨 Manejo de Errores

### Error Estructural (API, config, agente)
→ **DETENER** y reportar con formato:
```
🚨 ERROR ESTRUCTURAL
Tipo: [error]
Mensaje: [mensaje exacto]
Archivo entrada: sprint/current/analysis/readme.md
```

### Error de Ejecución (archivo faltante, contenido)
→ **EXPLICAR** y presentar opciones al usuario
