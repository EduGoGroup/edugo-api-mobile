#!/usr/bin/env bash
#
# docker-entrypoint.sh
#
# Entrypoint para el contenedor `api-mobile`.
# - Espera por las dependencias TCP indicadas (Postgres, MongoDB, RabbitMQ u otras definidas por WAIT_FOR_SERVICES).
# - Luego ejecuta el comando proporcionado al contenedor (si existe) o `./main` por defecto.
#
# Ejemplos:
#  - Usar defaults (postgres:5432, mongodb:27017, rabbitmq:5672):
#      docker run ... image
#  - Especificar servicios a esperar:
#      docker run -e WAIT_FOR_SERVICES="postgres:5432,mongodb:27017" ...
#  - Ejecutar comando custom después de esperar:
#      docker run ... image /bin/sh -c "echo ok; ./main"
#

set -o errexit
set -o nounset
set -o pipefail

LOG_PREFIX="[entrypoint]"

log() {
  printf "%s %s\n" "$LOG_PREFIX" "$*"
}

# Timeout por defecto en segundos (puede sobrescribirse con WAIT_TIMEOUT)
DEFAULT_TIMEOUT=${WAIT_TIMEOUT:-60}

# Construir la lista de servicios a esperar.
# Prioridad:
# 1) Variable WAIT_FOR_SERVICES (coma-separada o espacio-separada)
# 2) Valores por defecto basados en hosts/ports esperados
TARGETS=()

if [ -n "${WAIT_FOR_SERVICES:-}" ]; then
  # Reemplazar comas por espacios y leer en array
  IFS=', ' read -r -a tmp <<< "$WAIT_FOR_SERVICES"
  for t in "${tmp[@]}"; do
    t_trim="$(echo "$t" | xargs)" # trim
    if [ -n "$t_trim" ]; then
      TARGETS+=("$t_trim")
    fi
  done
else
  POSTGRES_HOST=${POSTGRES_HOST:-postgres}
  POSTGRES_PORT=${POSTGRES_PORT:-5432}
  MONGO_HOST=${MONGO_HOST:-mongodb}
  MONGO_PORT=${MONGO_PORT:-27017}
  RABBIT_HOST=${RABBIT_HOST:-rabbitmq}
  RABBIT_PORT=${RABBIT_PORT:-5672}

  TARGETS+=("${POSTGRES_HOST}:${POSTGRES_PORT}")
  TARGETS+=("${MONGO_HOST}:${MONGO_PORT}")
  TARGETS+=("${RABBIT_HOST}:${RABBIT_PORT}")
fi

# Asegurarse de que el script wait-for exista y sea ejecutable.
WAIT_SCRIPT="/scripts/wait-for.sh"
if [ -f "$WAIT_SCRIPT" ]; then
  if [ ! -x "$WAIT_SCRIPT" ]; then
    log "Haciendo ejecutable $WAIT_SCRIPT"
    chmod +x "$WAIT_SCRIPT" || true
  fi
else
  log "Aviso: $WAIT_SCRIPT no encontrado. Se intentará arrancar la aplicación sin espera."
fi

# Si tenemos targets y el wait script, esperar por ellos.
if [ "${#TARGETS[@]}" -gt 0 ] && [ -x "$WAIT_SCRIPT" ]; then
  log "Esperando servicios: ${TARGETS[*]} (timeout=${DEFAULT_TIMEOUT}s)"
  # Ejecutar wait-for con timeout y los targets
  # El script wait-for acepta --timeout y lista host:port
  "$WAIT_SCRIPT" --timeout "$DEFAULT_TIMEOUT" "${TARGETS[@]}"
  WAIT_EXIT=$?
  if [ "$WAIT_EXIT" -ne 0 ]; then
    log "Error: alguno de los servicios no respondió a tiempo (exit $WAIT_EXIT)."
    exit "$WAIT_EXIT"
  fi
  log "Servicios disponibles."
else
  if [ -x "$WAIT_SCRIPT" ]; then
    log "No hay targets configurados. Saltando espera."
  else
    log "wait-for script no disponible. Saltando espera y arrancando la aplicación."
  fi
fi

# Forward signals to the child process (simple trap)
_term() {
  log "Recibida señal de terminación, reenviando al proceso hijo..."
  if [ -n "${CHILD_PID:-}" ]; then
    kill -TERM "$CHILD_PID" 2>/dev/null || true
    wait "$CHILD_PID" 2>/dev/null || true
  fi
  exit 0
}
trap _term TERM INT

# Ejecutar el comando pasado al contenedor (si existe), o ./main por defecto.
if [ "$#" -gt 0 ]; then
  log "Ejecutando comando proporcionado: $*"
  exec "$@"
else
  log "Arrancando binario predeterminado: ./main"
  exec ./main
fi
