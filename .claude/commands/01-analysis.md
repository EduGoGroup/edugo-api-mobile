---
description: Analyze sprint requirements and generate architectural documentation
allowed-tools: Read, Task, Bash
argument-hint: "[--source=sprint|current] [--phase=N] [--mode=full|quick]"
---

# Comando: Análisis de Sprint

## Descripción
Este comando inicia el proceso de análisis del sprint con múltiples opciones de configuración. Lee el archivo readme del sprint, prepara el contenido según los parámetros, y lo pasa al agente especializado de análisis.

## Argumentos Soportados

```bash
--source=sprint|current   # De dónde leer (default: current)
                         # sprint  = sprint/readme.md
                         # current = sprint/current/readme.md

--phase=N                # Analizar solo fase N (default: todas)
                         # Genera resumen general + detalle de fase

--mode=full|quick        # Modo de análisis (default: full)
                         # full  = Con todos los diagramas
                         # quick = Solo análisis ejecutivo
```

## Ejemplos de Uso

```bash
/01-analysis                                    # current, todas, full
/01-analysis --mode=quick                       # current, todas, quick
/01-analysis --source=sprint                    # sprint root, todas, full
/01-analysis --source=sprint --phase=3          # sprint root, fase 3, full
/01-analysis --phase=2 --mode=quick             # current, fase 2, quick
/01-analysis --source=current --phase=1 --mode=full  # Explícito
```

## Responsabilidades del Comando

1. **Parsear** argumentos del usuario
2. **Leer** el archivo readme correspondiente (completo)
3. **Preparar** el contenido según `--phase`:
   - Si no hay `--phase`: pasar contenido completo
   - Si hay `--phase=N`: crear resumen general + detalle de fase N
4. **Invocar** al agente `analysis` con:
   - Contenido preparado
   - Parámetros: MODE y SCOPE
5. **No permitir** que el agente lea archivos por sí mismo

## Instrucciones de Ejecución

### Paso 1: Parsear Argumentos

Extraer los argumentos proporcionados:

```javascript
// Valores por defecto
let source = "current";  // sprint | current
let phase = null;        // null | número
let mode = "full";       // full | quick

// Parsear $ARGUMENTS
// Ejemplo: "--source=sprint --phase=3 --mode=quick"
```

### Paso 2: Determinar Ruta del Archivo

```javascript
let filePath;
if (source === "sprint") {
  filePath = "sprint/readme.md";
} else {
  filePath = "sprint/current/readme.md";
}
```

### Paso 3: Validar que el Archivo Existe

```bash
if [ ! -f "${filePath}" ]; then
  echo "❌ Error: No se encontró el archivo ${filePath}"
  echo ""
  echo "Verifica que el archivo existe."
  exit 1
fi
```

### Paso 4: Leer Contenido Completo

Lee el archivo completo y mantenlo en memoria:

```bash
CONTENT=$(cat "${filePath}")
```

### Paso 5: Preparar Contexto según --phase

#### Caso A: Sin --phase (analizar todo)

```markdown
PREPARED_CONTENT = CONTENT completo
SCOPE = "complete"
```

#### Caso B: Con --phase=N

1. **Extraer estructura general** del documento:
   - Objetivo del sprint
   - Lista de todas las fases (solo títulos)

2. **Extraer contenido completo de Fase N**

3. **Construir contexto preparado**:

```markdown
# CONTEXTO GENERAL DEL SPRINT

**Fuente**: ${filePath}

**Objetivo**: [extraído del inicio del readme]

**Estructura completa**:
- Fase 1: [título] (X tareas)
- Fase 2: [título] (Y tareas)
- Fase N: [título] (Z tareas) ← **FOCO DE ANÁLISIS**
- ...

---

# FASE ${N}: [TÍTULO DE LA FASE] (ANÁLISIS DETALLADO)

[Contenido COMPLETO de la fase ${N}, incluyendo:
- Objetivo de la fase
- Todas las tareas con descripciones
- Criterios de aceptación
- Dependencias
- Todo lo que esté en esa sección]
```

**Importante**: El agente NO debe filtrar nada. El comando ya hizo el trabajo de preparación.

### Paso 6: Determinar Parámetros para el Agente

```javascript
let MODE = mode;           // "full" o "quick"
let SCOPE = phase ? `phase-${phase}` : "complete";
let SOURCE_FILE = filePath;
```

### Paso 7: Invocar Agente de Análisis

Usar la herramienta Task con `subagent_type: "general-purpose"`:

**Prompt a pasar al agente**:

```markdown
[Contenido del archivo .claude/agents/analysis.md]

---

# CONFIGURACIÓN DE ESTA EJECUCIÓN

**MODE**: ${MODE}
**SCOPE**: ${SCOPE}
**SOURCE**: ${SOURCE_FILE}

---

# CONTENIDO A ANALIZAR

${PREPARED_CONTENT}

---

**INSTRUCCIONES IMPORTANTES**:

1. El contenido anterior ya está preparado y filtrado
2. NO debes leer ningún archivo adicional
3. Trabaja SOLO con el contenido proporcionado arriba
4. Genera archivos según MODE:
   - MODE=full: architecture.md, data-model.md, process-diagram.md, readme.md
   - MODE=quick: solo readme.md (análisis ejecutivo detallado)
5. Si SCOPE=phase-N, nombra archivos con sufijo: architecture-phase-N.md
6. Si SCOPE=complete, usa nombres normales
```

### Paso 8: Mensaje de Confirmación

Una vez que el agente complete su trabajo:

#### Si MODE=full y SCOPE=complete:
```
✅ Análisis completado exitosamente

📁 Archivos generados en sprint/current/analysis/:
- architecture.md (Diagrama de arquitectura con Mermaid validado)
- data-model.md (Diagrama ER y estructura de base de datos)
- process-diagram.md (Flujo de procesos del sistema)
- readme.md (Resumen ejecutivo)

📌 Siguiente paso: Ejecuta /02-planning para generar el plan de tareas
```

#### Si MODE=quick y SCOPE=complete:
```
✅ Análisis rápido completado exitosamente

📁 Archivo generado en sprint/current/analysis/:
- readme.md (Análisis ejecutivo detallado)

💡 Para análisis completo con diagramas, ejecuta:
   /01-analysis --mode=full

📌 Siguiente paso: Ejecuta /02-planning para generar el plan de tareas
```

#### Si SCOPE=phase-N:
```
✅ Análisis de Fase ${N} completado exitosamente

📁 Archivos generados en sprint/current/analysis/:
${MODE === 'full' ? 
  `- architecture-phase-${N}.md
- data-model-phase-${N}.md
- process-diagram-phase-${N}.md
- readme-phase-${N}.md` : 
  `- readme-phase-${N}.md`}

📌 Puedes analizar otras fases con:
   /01-analysis --phase=X
```

## Notas Importantes

- Este comando actúa como **orquestador inteligente**
- **El comando prepara el contexto** - el agente NO filtra
- El agente está **completamente aislado** - solo recibe lo que el comando le pasa
- Esto garantiza **control total** sobre qué información procesa el agente
- La preparación de contexto permite **análisis enfocado** sin perder visión general

## Algoritmo de Extracción de Fases (Referencia)

Para extraer la fase N del contenido:

```bash
# Buscar líneas que indican inicio de fase
# Patrones comunes:
# - "### Fase N:"
# - "## Fase N:"
# - "# Fase N:"

# Extraer desde inicio de Fase N hasta inicio de Fase N+1 (o fin de archivo)
```

Puedes usar herramientas como `sed`, `awk`, o simplemente parseo de markdown para extraer secciones.
