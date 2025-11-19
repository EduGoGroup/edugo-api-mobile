# Scripts de Implementación - edugo-api-mobile

Colección de scripts bash listos para ejecutar durante los Sprints 2 y 4.

---

## 📋 Scripts Disponibles

### Sprint 2: Migración + Optimización

1. **prepare-sprint-2.sh**
   - Duración: ~5 min
   - Función: Setup inicial, crear backups, verificar herramientas
   - Prerequisitos: Ninguno
   - Genera: Rama de trabajo, backups

2. **migrate-to-go-1.25.sh**
   - Duración: ~10 min
   - Función: Migrar proyecto a Go 1.25
   - Prerequisitos: Script 1 ejecutado
   - Genera: Cambios en go.mod, workflows, Dockerfile

3. **validate-go-1.25-local.sh**
   - Duración: ~15 min
   - Función: Validación exhaustiva local
   - Prerequisitos: Script 2 ejecutado
   - Genera: Reporte de validación

4. **validate-go-1.25-ci.sh**
   - Duración: ~15-30 min (incluye espera de CI)
   - Función: Crear PR y monitorear CI
   - Prerequisitos: Script 3 exitoso
   - Genera: PR draft en GitHub

5. **implement-parallelism-pr-to-dev.sh**
   - Duración: ~5 min
   - Función: Implementar paralelismo en workflow
   - Prerequisitos: Go 1.25 validado en CI
   - Genera: Workflow optimizado

### Sprint 4: Workflows Reusables

6. **setup-infrastructure-reusables.sh**
   - Duración: ~5 min
   - Función: Preparar infrastructure para workflows reusables
   - Prerequisitos: Sprint 2 completado
   - Genera: Estructura en infrastructure

7. **create-pr-validation-reusable.sh**
   - Duración: ~10 min
   - Función: Crear workflow reusable de validación
   - Prerequisitos: Script 6 ejecutado
   - Genera: pr-validation.yml reusable

8. **create-sync-branches-reusable.sh**
   - Duración: ~5 min
   - Función: Crear workflow reusable de sincronización
   - Prerequisitos: Script 6 ejecutado
   - Genera: sync-branches.yml reusable

---

## 🚀 Cómo Usar

### Opción A: Ejecución Manual

```bash
# 1. Navegar al directorio de scripts
cd /Users/jhoanmedina/source/EduGo/Analisys/00-Projects-Isolated/cicd-analysis/implementation-plans/03-api-mobile/SCRIPTS

# 2. Dar permisos de ejecución (primera vez)
chmod +x *.sh

# 3. Ejecutar scripts en orden
./prepare-sprint-2.sh
./migrate-to-go-1.25.sh
./validate-go-1.25-local.sh
./validate-go-1.25-ci.sh
./implement-parallelism-pr-to-dev.sh
```

### Opción B: Copiar/Pegar desde SPRINT-*-TASKS.md

Los scripts están completamente documentados en los archivos de tareas. Puedes copiarlos directamente desde allí.

---

## ⚠️ IMPORTANTE

1. **Lee la tarea completa antes de ejecutar el script**
   - Cada script tiene contexto y validaciones en SPRINT-*-TASKS.md
   - Entiende qué hace antes de correr

2. **Ejecuta EN ORDEN**
   - Los scripts tienen dependencias entre sí
   - No saltar pasos

3. **Valida después de cada script**
   - Usa los "Checkpoints" de cada tarea
   - Confirma que funcionó antes de continuar

4. **Ten un plan de rollback**
   - Los scripts crean backups automáticamente
   - Sabe cómo revertir si algo falla

---

## 📊 Estado de Generación

| Script | Estado | Ubicación |
|--------|--------|-----------|
| prepare-sprint-2.sh | ⏳ Por crear | Código en SPRINT-2-TASKS.md Tarea 2.1 |
| migrate-to-go-1.25.sh | ⏳ Por crear | Código en SPRINT-2-TASKS.md Tarea 2.2 |
| validate-go-1.25-local.sh | ⏳ Por crear | Código en SPRINT-2-TASKS.md Tarea 2.3 |
| validate-go-1.25-ci.sh | ⏳ Por crear | Código en SPRINT-2-TASKS.md Tarea 2.4 |
| implement-parallelism-pr-to-dev.sh | ⏳ Por crear | Código en SPRINT-2-TASKS.md Tarea 2.5 |
| setup-infrastructure-reusables.sh | ⏳ Por crear | Código en SPRINT-4-TASKS.md Tarea 4.1 |
| create-pr-validation-reusable.sh | ⏳ Por crear | Código en SPRINT-4-TASKS.md Tarea 4.2 |
| create-sync-branches-reusable.sh | ⏳ Por crear | Código en SPRINT-4-TASKS.md Tarea 4.3 |

**Nota:** Los scripts están documentados en los archivos SPRINT-*-TASKS.md. Puedes copiarlos de allí y guardarlos aquí cuando los necesites.

---

## 🛠️ Crear Scripts Desde Documentación

Si prefieres tener todos los scripts como archivos antes de comenzar:

```bash
# Este comando extraerá todos los scripts de SPRINT-2-TASKS.md
# y los guardará en este directorio

# TODO: Agregar script extractor cuando sea necesario
```

---

**Última actualización:** 19 de Noviembre, 2025  
**Generado por:** Claude Code
