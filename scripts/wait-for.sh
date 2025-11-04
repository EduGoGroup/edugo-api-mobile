#!/usr/bin/env bash
#
# wait-for.sh
#
# Script para esperar a que uno o varios servicios TCP estén disponibles antes de
# ejecutar una acción (útil en entrypoints de Docker).
#
# Uso:
#   ./wait-for.sh host:port [host:port ...] [--timeout SECONDS] [--quiet] [-- command args...]
#
# También puede leerse desde la variable de entorno WAIT_FOR_SERVICES con servicios
# separados por comas: WAIT_FOR_SERVICES="postgres:5432,mongo:27017" ./wait-for.sh
#
# Ejemplos:
#   ./wait-for.sh postgres:5432 rabbitmq:5672 --timeout 60 -- echo "Servicios listos"
#   WAIT_FOR_SERVICES="postgres:5432" ./wait-for.sh -- ./start-server.sh
#
# Requisitos:
#  - Recomendado: netcat-openbsd (nc). Si no está disponible intenta usar /dev/tcp.
#
set -o errexit
set -o nounset
set -o pipefail

PROG_NAME="$(basename "$0")"

DEFAULT_TIMEOUT=60
SLEEP_INTERVAL=1

TIMEOUT=${TIMEOUT:-$DEFAULT_TIMEOUT}
QUIET=0
TARGETS=()

print_help() {
  cat <<EOF
Uso: $PROG_NAME host:port [host:port ...] [--timeout SECONDS] [--quiet] [-- command args...]

Espera hasta que todos los pares host:port estén accesibles por TCP. Si se proporciona
\"--\" seguido de un comando, ese comando se ejecutará (con exec) al terminar correctamente.

Opciones:
  -t, --timeout SECONDS   Tiempo máximo a esperar por cada servicio (por defecto: ${DEFAULT_TIMEOUT}s)
  -q, --quiet             Modo silencioso (solo códigos de salida, sin logs)
  -h, --help              Mostrar esta ayuda
EOF
}

log() {
  if [ "$QUIET" -eq 0 ]; then
    echo "$@"
  fi
}

# Parse arguments
ARGS=()
while [ $# -gt 0 ]; do
  case "$1" in
    -t|--timeout)
      if [ -z "${2:-}" ]; then
        echo "Error: falta valor para $1" >&2
        exit 2
      fi
      TIMEOUT="$2"
      shift 2
      ;;
    -q|--quiet)
      QUIET=1
      shift
      ;;
    -h|--help)
      print_help
      exit 0
      ;;
    --)
      shift
      # Remaining args after -- are the command to run
      ARGS=("$@")
      break
      ;;
    *)
      # Posibles targets host:port o variable WAIT_FOR_SERVICES si no hay argumentos
      TARGETS+=("$1")
      shift
      ;;
  esac
done

# If TARGETS empty, check WAIT_FOR_SERVICES env var
if [ "${#TARGETS[@]}" -eq 0 ] && [ -n "${WAIT_FOR_SERVICES:-}" ]; then
  IFS=',' read -r -a TMP <<< "$WAIT_FOR_SERVICES"
  for t in "${TMP[@]}"; do
    t_trimmed="$(echo "$t" | xargs)"
    if [ -n "$t_trimmed" ]; then
      TARGETS+=("$t_trimmed")
    fi
  done
fi

if [ "${#TARGETS[@]}" -eq 0 ]; then
  echo "Error: no se especificaron servicios a esperar." >&2
  print_help
  exit 2
fi

# Check if nc is available
NC_CMD=""
if command -v nc >/dev/null 2>&1; then
  NC_CMD="nc"
fi

# Helper para testear conexión TCP
test_tcp() {
  local host="$1"
  local port="$2"

  if [ -n "$NC_CMD" ]; then
    # netcat openbsd/gnu: use -z for scan, -w for timeout (seconds)
    if nc -z -w 1 "$host" "$port" >/dev/null 2>&1; then
      return 0
    fi
    return 1
  else
    # Fallback a /dev/tcp (bash-specific)
    if (echo > "/dev/tcp/${host}/${port}") >/dev/null 2>&1; then
      return 0
    fi
    return 1
  fi
}

# Esperar por un único host:port con timeout
wait_for_target() {
  local target="$1"
  local timeout_seconds="$2"

  # Soporta host:port (si hay más ':' (IPv6) se maneja parcialmente)
  local host="${target%%:*}"
  local port="${target#*:}"

  if [ -z "$host" ] || [ -z "$port" ] || [ "$host" = "$port" ]; then
    echo "Formato inválido para target: '$target'. Debe ser host:port" >&2
    return 2
  fi

  log "Esperando por $host:$port (timeout: ${timeout_seconds}s)..."
  local start_time
  start_time=$(date +%s)
  while true; do
    if test_tcp "$host" "$port"; then
      log "✓ $host:$port está disponible"
      return 0
    fi

    local now
    now=$(date +%s)
    local elapsed=$(( now - start_time ))
    if [ "$elapsed" -ge "$timeout_seconds" ]; then
      echo "⨯ Timeout esperando $host:$port después de ${timeout_seconds}s" >&2
      return 1
    fi

    sleep "$SLEEP_INTERVAL"
  done
}

# Ejecutar espera para todos los targets
EXIT_CODE=0
for t in "${TARGETS[@]}"; do
  if ! wait_for_target "$t" "$TIMEOUT"; then
    EXIT_CODE=1
    break
  fi
done

if [ "$EXIT_CODE" -ne 0 ]; then
  echo "Al menos un servicio no respondió a tiempo." >&2
  exit "$EXIT_CODE"
fi

# Si hay comando a ejecutar, lo reemplazamos con exec
if [ "${#ARGS[@]}" -gt 0 ]; then
  log "Ejecutando comando: ${ARGS[*]}"
  exec "${ARGS[@]}"
else
  log "Todos los servicios están disponibles."
  exit 0
fi
