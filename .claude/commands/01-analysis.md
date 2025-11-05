---
description: Analyze sprint requirements and generate architectural documentation
argument-hint: "[--source=sprint|current] [--phase=N] [--mode=full|quick]"
---

# Comando: Análisis de Sprint

## 🎯 Rol: ORQUESTADOR

**Este comando NO ejecuta el análisis. Delega al agente `flow-analysis`.**

Tu función:
1. Parsear argumentos del usuario
2. Invocar al agente especializado usando **Task tool**
3. Retornar resultado al usuario

---

## Argumentos

```bash
--source=sprint|current   # Default: current
--phase=N                 # Default: todas las fases
--mode=full|quick         # Default: full
```

**Modos:**
- `full`: Genera diagramas (architecture.md, data-model.md, process-diagram.md, readme.md)
- `quick`: Solo análisis ejecutivo (readme.md)

---

## Ejecución

### 1. Procesar Argumentos

```
MODE = "full" (o "quick" si --mode=quick)
SOURCE = "current" (o "sprint" si --source=sprint)
PHASE = null (o N si --phase=N)
```

### 2. Invocar Agente flow-analysis

**USA TASK TOOL:**

```
Task(
  subagent_type: "flow-analysis",
  description: "Análisis arquitectónico de sprint",
  prompt: "
    Genera análisis arquitectónico del sprint.

    PARÁMETROS:
    - MODE: {MODE}
    - SOURCE: {SOURCE}
    - PHASE: {PHASE o 'todas'}

    ARCHIVO: sprint/{SOURCE}/readme.md

    SALIDA:
    - Si MODE=full: architecture.md, data-model.md, process-diagram.md, readme.md
    - Si MODE=quick: solo readme.md
    - Si PHASE=N: agregar sufijo -phase-{N} a los archivos

    UBICACIÓN: sprint/current/analysis/
  "
)
```

### 3. Confirmar al Usuario

**Si MODE=full:**
```
✅ Análisis completo exitoso

📁 Archivos en sprint/current/analysis/:
- architecture.md
- data-model.md
- process-diagram.md
- readme.md

📌 Siguiente: /02-planning
```

**Si MODE=quick:**
```
✅ Análisis rápido exitoso

📁 Archivo: sprint/current/analysis/readme.md
💡 Para análisis completo: /01-analysis --mode=full
📌 Siguiente: /02-planning
```

**Si PHASE=N:**
```
✅ Análisis de Fase {N} exitoso

📁 Archivos: *-phase-{N}.md en sprint/current/analysis/
📌 Analizar otra fase: /01-analysis --phase=X
```

---

## 🚨 Manejo de Errores

### Error Estructural (API, config, agente)
→ **DETENER** y reportar con formato:
```
🚨 ERROR ESTRUCTURAL
Tipo: [error]
Mensaje: [mensaje exacto]
Parámetros: MODE=X, SOURCE=Y, PHASE=Z
```

### Error de Ejecución (archivo faltante, contenido)
→ **EXPLICAR** y presentar opciones al usuario
