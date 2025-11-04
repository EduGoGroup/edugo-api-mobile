---
description: Execute planned sprint tasks (all, specific phase, or task)
argument-hint: "[phase-N|task-N.M]"
---

# Comando: Ejecución de Sprint

## Descripción
Este comando ejecuta las tareas planeadas del sprint. Puede ejecutar el plan completo o fases/tareas específicas según los argumentos proporcionados. Lee el plan de ejecución y opcionalmente las reglas del proyecto, luego invoca al agente de ejecución.

## Sintaxis
```bash
/03-execution              # Ejecutar plan completo
/03-execution phase-1      # Ejecutar solo la fase 1
/03-execution task-2.3     # Ejecutar solo la tarea 3 de la fase 2
```

## Responsabilidades del Comando
1. **Leer** el archivo `sprint/current/planning/readme.md`
2. **Filtrar** contenido según argumentos (si se proporcionan)
3. **Leer** el archivo `sprint/current/execution/rules.md` (si existe)
4. **Invocar** al agente `execution` con las tareas y reglas
5. **Permitir acceso limitado** para que el agente acceda a las carpetas analysis/planning para contexto adicional

## Instrucciones de Ejecución

Por favor, ejecuta los siguientes pasos:

### Paso 1: Validar archivo de entrada
Verifica que existe el archivo `sprint/current/planning/readme.md`. Si no existe:
```
❌ Error: Plan de sprint no encontrado

Por favor ejecuta primero: /02-planning
```

### Paso 2: Leer plan de trabajo
Lee el archivo completo `sprint/current/planning/readme.md`.

### Paso 3: Procesar argumentos (si los hay)
Si el usuario proporcionó argumentos (ej: `phase-1`, `task-2.3`):
- Extrae del plan solo la sección correspondiente a esa fase/tarea
- Verifica las dependencias de esa fase/tarea
- Si hay dependencias no completadas, advierte al usuario pero permite continuar

Si NO hay argumentos:
- Usa el plan completo

### Paso 4: Verificar reglas del proyecto
Verifica si existe el archivo `sprint/current/execution/rules.md`:
```bash
Si existe → Léelo y pásalo al agente
Si NO existe → Continúa sin reglas (el agente usará mejores prácticas)
```

### Paso 5: Invocar agente de ejecución
Usa la herramienta Task con `subagent_type: "general-purpose"` para invocar al agente de ejecución.

Pasa al agente:
- **Prompt completo**: Las instrucciones del agente (lee `.claude/agents/execution.md`)
- **Tareas a ejecutar**: Plan completo o filtrado según paso 3
- **Reglas del proyecto**: Contenido de rules.md (si existe)
- **Permisos especiales**:
  - Puede leer archivos de `sprint/current/analysis/` y `sprint/current/planning/` para contexto adicional
  - Puede escribir/modificar archivos en la carpeta raíz del proyecto
  - NO PUEDE tocar la carpeta `.claude/`
  - NO PUEDE tocar la carpeta `sprint/` excepto para escribir reportes en `sprint/current/execution/`

### Paso 6: Mensaje de confirmación
Una vez que el agente completa su trabajo, informa al usuario:
```
✅ Ejecución completada exitosamente

📁 Reporte generado:
- sprint/current/execution/[phase-step]-[timestamp].md

✅ Validaciones realizadas:
- Código compiló correctamente
- Tests ejecutados (si aplica)

📌 Siguiente paso:
- Ejecuta /04-review para ver el estado consolidado del sprint
- O ejecuta /03-execution [otra-fase] para continuar con otras tareas
```

Si hubo errores de compilación o tests fallidos:
```
⚠️ Ejecución completada con advertencias

📁 Reporte generado:
- sprint/current/execution/[phase-step]-[timestamp].md

⚠️ Problemas detectados:
[Lista de problemas]

📌 Recomendación:
Revisa el reporte de ejecución y corrige los problemas antes de continuar
```

## Notas Importantes
- Este comando permite **ejecución modular** - puedes ejecutar fases/tareas específicas
- El agente **valida que el código compile** antes de marcar la tarea como completada
- Las **reglas del proyecto** son opcionales pero recomendadas para consistencia
- Cada ejecución genera un **reporte separado** con timestamp para trazabilidad
