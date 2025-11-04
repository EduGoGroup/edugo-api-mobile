---
description: Consolidate sprint status and generate validation guide
allowed-tools: Read, Task
argument-hint: ""
---

# Comando: Revisión de Sprint

## Descripción
Este comando consolida el estado completo del sprint. Lee el plan original y todos los reportes de ejecución, luego invoca al agente de revisión para generar un documento que muestre qué tareas se completaron y proporcione una guía de validación para el usuario.

## Responsabilidades del Comando
1. **Leer** el archivo `sprint/current/planning/readme.md` (plan original)
2. **Leer** todos los archivos de reporte en `sprint/current/execution/*.md`
3. **Invocar** al agente `review` con toda la información
4. **Generar** un documento consolidado con estado y guía de validación

## Instrucciones de Ejecución

Por favor, ejecuta los siguientes pasos:

### Paso 1: Validar archivos de entrada
Verifica que existe el archivo `sprint/current/planning/readme.md`. Si no existe:
```
❌ Error: Plan de sprint no encontrado

Por favor ejecuta primero: /02-planning
```

Verifica que existen archivos en `sprint/current/execution/`. Si no hay ninguno:
```
ℹ️ Advertencia: No se encontraron reportes de ejecución

El sprint no tiene tareas ejecutadas todavía.
¿Quieres generar un reporte de estado de todas formas? (útil para ver qué falta)
```

### Paso 2: Leer plan de trabajo
Lee el archivo completo `sprint/current/planning/readme.md`.

### Paso 3: Leer todos los reportes de ejecución
Lista y lee todos los archivos en `sprint/current/execution/*.md` (excepto rules.md si existe).

Organiza los reportes cronológicamente para dárselos al agente en orden.

### Paso 4: Invocar agente de revisión
Usa la herramienta Task con `subagent_type: "general-purpose"` para invocar al agente de revisión.

Pasa al agente:
- **Prompt completo**: Las instrucciones del agente (lee `.claude/agents/review.md`)
- **Plan original**: Contenido de `sprint/current/planning/readme.md`
- **Reportes de ejecución**: Todos los archivos leídos en paso 3, en orden cronológico
- **Instrucción especial**: El agente debe generar una sección final "Guía de Validación para el Usuario"

### Paso 5: Mensaje de confirmación
Una vez que el agente completa su trabajo, informa al usuario:
```
✅ Revisión completada exitosamente

📁 Archivo generado:
- sprint/current/review/readme.md

📊 Contenido del reporte:
- Plan original con tareas marcadas como completadas ✅
- Resumen de tareas pendientes
- Guía de validación para probar el sprint

📌 Siguiente paso:
- Lee sprint/current/review/readme.md para ver el estado completo
- Usa la "Guía de Validación" al final del documento para probar la aplicación
- Si todo está completo, ejecuta /archive para archivar este sprint
- Si faltan tareas, ejecuta /03-execution [fase] para continuar
```

### Paso 6: Mostrar resumen rápido (opcional)
Opcionalmente, puedes mostrar un resumen rápido en consola:
```
📈 Resumen del Sprint:
├─ Fases totales: X
├─ Fases completadas: Y
├─ Tareas totales: A
├─ Tareas completadas: B
└─ Progreso: ZZ%
```

## Notas Importantes
- Este comando da **visibilidad completa** del estado del sprint
- La **guía de validación** es crucial - debe ser simple y práctica para el usuario
- Permite tomar decisiones sobre qué hacer a continuación (continuar, archivar, o corregir)
- Útil para presentaciones/demos mostrando el progreso del trabajo
