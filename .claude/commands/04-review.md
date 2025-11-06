---
description: Consolidate sprint status and generate validation guide
argument-hint: ""
---

# Comando: Revisión de Sprint

## 🎯 Rol: ORQUESTADOR

**Este comando NO ejecuta la revisión. Delega al agente `flow-review`.**

Tu función:
1. Validar que existen plan y reportes
2. Leer todos los documentos necesarios
3. Invocar al agente especializado usando **Task tool**
4. Retornar resultado al usuario

---

## Ejecución

### 1. Validar Entrada

Verificar que existe: `sprint/current/planning/readme.md`

Si no existe:
```
❌ Error: Plan de sprint no encontrado
Por favor ejecuta primero: /02-planning
```

Verificar reportes en: `sprint/current/execution/*.md`

Si no hay reportes:
```
ℹ️ Advertencia: No hay reportes de ejecución
¿Generar reporte de estado de todas formas? (útil para ver qué falta)
```

### 2. Leer Documentos

```
- Plan original: sprint/current/planning/readme.md
- Reportes: sprint/current/execution/*.md (excepto rules.md)
  Organizarlos cronológicamente
```

### 3. Invocar Agente flow-review

**USA TASK TOOL:**

```
Task(
  subagent_type: "flow-review",
  description: "Revisión consolidada del sprint",
  prompt: "
    Consolida el estado completo del sprint.

    ENTRADA:
    - Plan original: sprint/current/planning/readme.md
    - Reportes de ejecución: [lista cronológica]

    SALIDA: sprint/current/review/readme.md

    Genera documento con:
    - Plan original con tareas marcadas ✅ según reportes
    - Resumen de tareas completadas/pendientes
    - Progreso por fase (%)
    - Guía de Validación para el Usuario (CRUCIAL)
      * Pasos simples para verificar funcionalidad
      * Comandos a ejecutar
      * Resultados esperados

    PERMISOS:
    - Leer: sprint/current/planning/, sprint/current/execution/
    - Escribir: sprint/current/review/readme.md
  "
)
```

### 4. Confirmar al Usuario

```
✅ Revisión completada

📁 Archivo: sprint/current/review/readme.md

📊 Contenido:
- Plan con tareas marcadas ✅
- Resumen de pendientes
- Guía de validación práctica

📈 Resumen rápido:
├─ Fases totales: X
├─ Fases completadas: Y
├─ Tareas totales: A
├─ Tareas completadas: B
└─ Progreso: ZZ%

📌 Siguiente:
- Lee el review para ver estado completo
- Usa la Guía de Validación al final
- Si todo OK: /archive para archivar sprint
- Si falta: /03-execution [fase] para continuar
```

---

## 🚨 Manejo de Errores

### Error Estructural (API, config, agente)
→ **DETENER** y reportar con formato:
```
🚨 ERROR ESTRUCTURAL
Tipo: [error]
Mensaje: [mensaje exacto]
Documentos procesados: [lista]
```

### Error de Ejecución (documentos incompletos, inconsistencias)
→ **EXPLICAR** problema y **PRESENTAR OPCIONES**:
```
⚠️ PROBLEMA DE EJECUCIÓN
- Opción A: Generar revisión parcial con info disponible
- Opción B: Marcar solo tareas confirmadas
- Opción C: Necesito documentos adicionales
```
