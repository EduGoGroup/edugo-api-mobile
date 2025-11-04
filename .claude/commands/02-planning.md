---
description: Transform architectural analysis into granular task plan
argument-hint: ""
---

# Comando: Planificación de Sprint

## Descripción
Este comando transforma el análisis arquitectónico en un plan de trabajo granular. Lee el resumen ejecutivo generado por el análisis y lo pasa al agente planificador para crear un documento estructurado con fases, tareas atómicas y dependencias.

## Responsabilidades del Comando
1. **Leer** el archivo `sprint/current/analysis/readme.md` (salida del análisis)
2. **Validar** que el análisis se ha completado
3. **Invocar** al agente `planner` pasándole el contenido
4. **Mantener aislamiento** del agente

## Instrucciones de Ejecución

Por favor, ejecuta los siguientes pasos:

### Paso 1: Validar archivo de entrada
Verifica que existe el archivo `sprint/current/analysis/readme.md`. Si no existe, informa al usuario:
```
❌ Error: Análisis de sprint no encontrado

Por favor ejecuta primero: /01-analysis
```

### Paso 2: Leer contenido del análisis
Lee el archivo completo `sprint/current/analysis/readme.md` y mantenlo en contexto.

**Opcional pero recomendado**: También lee los otros archivos de análisis para más contexto:
- `sprint/current/analysis/architecture.md`
- `sprint/current/analysis/data-model.md`
- `sprint/current/analysis/process-diagram.md`

### Paso 3: Invocar agente planificador
Usa la herramienta Task con `subagent_type: "general-purpose"` para invocar al agente planificador.

Pasa al agente:
- **Prompt completo**: Las instrucciones del agente (lee `.claude/agents/planner.md`)
- **Contexto del análisis**: El contenido del readme.md del análisis (y opcionalmente otros documentos)
- **Restricción explícita**: El agente solo trabaja con lo que recibe del comando

### Paso 4: Mensaje de confirmación
Una vez que el agente completa su trabajo, informa al usuario:
```
✅ Planificación completada exitosamente

📁 Archivo generado:
- sprint/current/planning/readme.md

📋 Contenido del plan:
- Fases organizadas con casillas de verificación
- Tareas atómicas y granulares
- Indicadores de dependencia entre tareas
- Listo para ejecución modular

📌 Siguiente paso:
- Ejecuta /03-execution para implementar todo el plan
- O ejecuta /03-execution phase-1 para una fase específica
```

## Notas Importantes
- Este comando es el **puente** entre el análisis y la ejecución
- El plan generado debe ser lo suficientemente granular para permitir la ejecución aislada de cada tarea
- Las dependencias claramente marcadas ayudan a tomar decisiones sobre el orden de ejecución
