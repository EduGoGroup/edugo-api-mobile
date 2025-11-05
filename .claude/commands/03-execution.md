---
description: Execute planned sprint tasks (all, specific phase, or task)
argument-hint: "[phase-N|task-N.M]"
---

# Comando: Ejecución de Sprint

## 🎯 Rol: ORQUESTADOR

**Este comando NO ejecuta las tareas. Delega al agente `flow-execution`.**

Tu función:
1. Validar que existe el plan
2. Procesar argumentos (all, phase-N, task-N.M)
3. Invocar al agente especializado usando **Task tool**
4. Retornar resultado al usuario

---

## Argumentos

```bash
/03-execution              # Ejecutar plan completo
/03-execution phase-1      # Solo fase 1
/03-execution task-2.3     # Solo tarea 3 de fase 2
```

---

## Ejecución

### 1. Validar Entrada

Verificar que existe: `sprint/current/planning/readme.md`

Si no existe:
```
❌ Error: Plan de sprint no encontrado
Por favor ejecuta primero: /02-planning
```

### 2. Procesar Argumentos

```
Si NO hay argumentos → Ejecutar plan completo
Si phase-N → Ejecutar solo fase N
Si task-N.M → Ejecutar solo tarea M de fase N
```

**Advertencia de dependencias:**
Si se ejecuta fase/tarea específica y hay dependencias no completadas, advertir pero permitir continuar.

### 3. Verificar Reglas (Opcional)

Verificar si existe: `sprint/current/execution/rules.md`
```
Si existe → Pasar al agente
Si NO existe → Continuar sin reglas
```

### 4. Invocar Agente flow-execution

**USA TASK TOOL:**

```
Task(
  subagent_type: "flow-execution",
  description: "Ejecución de tareas del sprint",
  prompt: "
    Ejecuta las tareas planeadas del sprint.

    ENTRADA:
    - Plan: sprint/current/planning/readme.md [completo o filtrado]
    - Reglas: sprint/current/execution/rules.md [si existe]

    ALCANCE: {todo el plan | phase-N | task-N.M}

    SALIDA: sprint/current/execution/[identificador]-[timestamp].md

    VALIDACIONES OBLIGATORIAS:
    - Código debe compilar (go build ./...)
    - Ejecutar tests si aplica
    - Marcar tareas completadas solo si validaciones pasan

    PERMISOS:
    - Leer: sprint/current/analysis/, sprint/current/planning/
    - Escribir: archivos del proyecto (excepto .claude/ y sprint/)
    - Reportar: sprint/current/execution/
  "
)
```

### 5. Confirmar al Usuario

**Si exitoso:**
```
✅ Ejecución completada

📁 Reporte: sprint/current/execution/[id]-[timestamp].md
✅ Código compiló correctamente
✅ Tests ejecutados

📌 Siguiente: /04-review (ver estado consolidado)
```

**Si hubo problemas:**
```
⚠️ Ejecución completada con advertencias

📁 Reporte: sprint/current/execution/[id]-[timestamp].md
⚠️ Problemas detectados: [lista]

📌 Revisa el reporte antes de continuar
```

---

## 🚨 Manejo de Errores

### Error Estructural (API, config, agente)
→ **DETENER** y reportar con formato:
```
🚨 ERROR ESTRUCTURAL
Tipo: [error]
Mensaje: [mensaje exacto]
Tarea ejecutando: [identificador]
```

### Error de Ejecución (compilación, tests, etc.)
→ El agente debe **REPORTAR** en el archivo de ejecución y **DETENER** esa tarea.
→ **NO continuar** con tareas dependientes.
→ **PRESENTAR OPCIONES** al usuario.
