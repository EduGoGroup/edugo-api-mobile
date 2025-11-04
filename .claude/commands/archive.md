---
description: Archive current sprint and prepare project for new cycle
allowed-tools: Bash, Read, Write
argument-hint: ""
---

# Comando: Archivar Sprint

## Descripción
Este comando archiva el sprint actual y prepara el proyecto para un nuevo ciclo. Renombra la carpeta sprint/current con un timestamp, la mueve a la carpeta sprint/archived, y crea una nueva carpeta sprint/current limpia.

## Responsabilidades del Comando
1. **Validar** que existe contenido para archivar
2. **Generar** nombre único basado en fecha/hora
3. **Mover** la carpeta sprint/current actual a sprint/archived
4. **Crear** nueva carpeta sprint/current limpia con estructura inicial
5. **Reportar** el proceso al usuario

## Formato de Nombre de Archivo
```
sprint-YYYY-MM-DD-HHMM
Ejemplo: sprint-2025-10-31-1430
```

## Instrucciones de Ejecución

Por favor, ejecuta los siguientes pasos:

### Paso 1: Validar carpeta sprint/current
Verifica que la carpeta `sprint/current/` existe y contiene archivos. Si está vacía o no existe:
```
ℹ️ Advertencia: La carpeta sprint/current está vacía o no existe

No hay nada que archivar. ¿Quieres crear una nueva carpeta sprint/current de todas formas?
```

### Paso 2: Generar nombre del archivo
Genera el nombre usando el formato de fecha/hora actual:
```bash
# Formato: sprint-YYYY-MM-DD-HHMM
# Ejemplo: sprint-2025-10-31-1430

# Usa el comando date para generar el timestamp
TIMESTAMP=$(date +"%Y-%m-%d-%H%M")
ARCHIVE_NAME="sprint-${TIMESTAMP}"
```

### Paso 3: Verificar que el nombre no existe
Verifica que no existe ya una carpeta con ese nombre en `sprint/archived/`:
```bash
# Si existe, agrega un sufijo numérico
# sprint-2025-10-31-1430-1
# sprint-2025-10-31-1430-2
# etc.
```

### Paso 4: Mover carpeta sprint/current al archivo
```bash
# Asegurar que existe sprint/archived/
mkdir -p sprint/archived

# Mover sprint/current a sprint/archived con nuevo nombre
mv sprint/current "sprint/archived/${ARCHIVE_NAME}"
```

### Paso 5: Crear nueva carpeta sprint/current limpia
```bash
# Crear estructura limpia
mkdir -p sprint/current/analysis
mkdir -p sprint/current/planning
mkdir -p sprint/current/execution
mkdir -p sprint/current/review
```

### Paso 6: Mensaje de confirmación
```
✅ Sprint archivado exitosamente

📦 Archivo creado:
- sprint/archived/sprint-2025-10-31-1430/

📁 Nueva carpeta sprint/current creada con estructura limpia:
- sprint/current/
  ├─ analysis/
  ├─ planning/
  ├─ execution/
  └─ review/

📌 Siguiente paso:
1. Crea un nuevo archivo sprint/current/readme.md con los requisitos del nuevo sprint
2. Ejecuta /01-analysis para iniciar el nuevo ciclo

💡 El sprint anterior está disponible en:
   sprint/archived/sprint-2025-10-31-1430/
```

### Paso 7: Sugerencia de creación de plantilla (opcional)
Opcionalmente, pregunta al usuario si quiere crear un archivo de plantilla readme.md vacío:
```
¿Quieres crear un archivo de plantilla sprint/current/readme.md para el nuevo sprint?

Si el usuario acepta, crea:
```markdown
# Sprint: [Nombre del Sprint]

## Descripción
[Breve descripción del objetivo del sprint]

## Requisitos
- [ ] Requisito 1
- [ ] Requisito 2
- [ ] Requisito 3

## Contexto
[Información adicional relevante]

## Entregables Esperados
1. [Entregable 1]
2. [Entregable 2]

## Restricciones/Consideraciones
- [Restricción 1]
- [Restricción 2]
```
```

## Notas Importantes
- Este comando tiene **permisos completos de lectura/escritura**
- Usa **timestamp** para evitar sobrescribir archivos previos
- Preserva el **historial completo** de sprints en la carpeta sprint/archived
- Garantiza que siempre hay una carpeta sprint/current limpia lista para un nuevo ciclo
- NO archivar si el usuario está en medio de un sprint importante (pedir confirmación si sprint/current tiene contenido reciente)
