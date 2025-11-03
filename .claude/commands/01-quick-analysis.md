---
description: Quick sprint analysis without diagrams (shortcut for --mode=quick)
allowed-tools: SlashCommand
argument-hint: "[--source=sprint|current] [--phase=N]"
---

# Comando: Análisis Rápido de Sprint

## Descripción
Atajo conveniente para ejecutar un análisis rápido sin diagramas. Internamente llama a `/01-analysis --mode=quick` con los argumentos adicionales que proporciones.

## Argumentos Soportados

```bash
--source=sprint|current   # De dónde leer (default: current)
--phase=N                # Analizar solo fase N (default: todas)
```

## Ejemplos de Uso

```bash
/01-quick-analysis                    # Equivale a: /01-analysis --mode=quick
/01-quick-analysis --source=sprint    # Equivale a: /01-analysis --source=sprint --mode=quick
/01-quick-analysis --phase=3          # Equivale a: /01-analysis --phase=3 --mode=quick
```

## Instrucciones de Ejecución

Este comando simplemente redirige al comando principal con el modo forzado:

### Paso 1: Construir Comando Completo

```bash
# Tomar argumentos del usuario
USER_ARGS="$ARGUMENTS"

# Construir comando con --mode=quick forzado
FULL_COMMAND="/01-analysis ${USER_ARGS} --mode=quick"
```

### Paso 2: Ejecutar Comando Principal

```bash
# Invocar el comando /01-analysis con los argumentos preparados
${FULL_COMMAND}
```

## Notas

- Este es un **atajo de conveniencia** para el caso de uso más común
- Internamente usa el comando `/01-analysis` con `--mode=quick`
- Útil cuando necesitas análisis rápido frecuentemente
- Si necesitas diagramas, usa `/01-analysis --mode=full` directamente
