---
description: Quick sprint analysis without diagrams (shortcut for --mode=quick)
argument-hint: "[--source=sprint|current] [--phase=N]"
---

# Comando: Análisis Rápido de Sprint

## 🎯 Rol: ORQUESTADOR

**Este comando NO ejecuta el análisis. Delega al agente `flow-analysissis`.**

Tu función:
1. Parsear argumentos del usuario
2. Invocar al agente especializado usando **Task tool**
3. Retornar resultado al usuario

---

## Argumentos

```bash
--source=sprint|current   # Default: current
--phase=N                 # Default: todas las fases
```

---

## Ejecución

### 1. Procesar Argumentos

```
MODE = "quick"  (forzado)
SOURCE = "current"  (o "sprint" si --source=sprint)
PHASE = null  (o N si --phase=N)
```

### 2. Invocar Agente flow-analysis

**USA TASK TOOL:**

```
Task(
  subagent_type: "flow-analysis",
  description: "Análisis rápido de sprint",
  prompt: "
    Genera análisis arquitectónico del sprint sin diagramas.

    PARÁMETROS:
    - MODE: quick
    - SOURCE: {SOURCE}
    - PHASE: {PHASE o 'todas'}

    ARCHIVO: sprint/{SOURCE}/readme.md
    SALIDA: sprint/current/analysis/readme.md (o readme-phase-{N}.md)
  "
)
```

### 3. Confirmar al Usuario

```
✅ Análisis rápido completado

📁 Archivo: sprint/current/analysis/readme.md
📌 Siguiente: /02-planning
```

---

## 🚨 Manejo de Errores

### Error Estructural (API, config, agente)
→ **DETENER** y reportar con formato:
```
🚨 ERROR ESTRUCTURAL
Tipo: [error]
Mensaje: [mensaje exacto]
Parámetros: MODE=quick, SOURCE=X, PHASE=Y
```

### Error de Ejecución (archivo faltante, contenido)
→ **EXPLICAR** y presentar opciones al usuario
