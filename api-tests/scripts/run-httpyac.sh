#!/usr/bin/env bash
set -euo pipefail
# Ejecuta login primero y luego todas las demás requests en la misma invocación de httpyac
# Uso: ./run-httpyac.sh [path-to-dotenv]

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$SCRIPT_DIR/.."
REQUESTS_DIR="$ROOT_DIR/requests"

# Flags
LOGIN_ONLY=false
PERSIST_TOKEN=false

# Parse args: optional dotenv and flags
DOTENV_FILE_DEFAULT="$ROOT_DIR/.env.local"
DOTENV_FILE=""
for arg in "$@"; do
  case "$arg" in
    --login-only)
      LOGIN_ONLY=true ;;
    --persist-token|-p)
      PERSIST_TOKEN=true ;;
    --dotenv=*)
      DOTENV_FILE="${arg#--dotenv=}" ;;
    *)
      # if arg looks like a file path to dotenv
      if [ -z "$DOTENV_FILE" ] && [ -f "$arg" ]; then
        DOTENV_FILE="$arg"
      fi ;;
  esac
done

if [ -z "$DOTENV_FILE" ]; then
  DOTENV_FILE="$DOTENV_FILE_DEFAULT"
fi

if ! command -v httpyac >/dev/null 2>&1; then
  echo "httpyac no está instalado. Instálalo con: npm install -g httpyac"
  exit 1
fi

if [ ! -f "$DOTENV_FILE" ]; then
  echo "Advertencia: dotenv no encontrado en $DOTENV_FILE. Puedes pasar otro archivo como primer parámetro."
fi

# Construir lista de ficheros: auth primero si existe, luego el resto (excluyendo auth)
AUTH_FILE="$REQUESTS_DIR/auth.http"

if [ "$PERSIST_TOKEN" = true ] || [ "$LOGIN_ONLY" = true ]; then
  if [ -f "$AUTH_FILE" ]; then
    echo "Ejecutando login para capturar tokens..."
    # Capturar stdout para parsear el JSON de tokens impreso por httpyac
    TMP_OUT=$(mktemp)
    httpyac "$AUTH_FILE" --dotenv "$DOTENV_FILE" > "$TMP_OUT" 2>&1 || true

    # Buscar línea JSON con el marcador __httpyac_tokens
    TOK_JSON=$(grep -a -m1 '__httpyac_tokens' "$TMP_OUT" | sed -n 's/.*\({.*"__httpyac_tokens".*}\).*/\1/p' || true)
    if [ -n "$TOK_JSON" ]; then
      # Extraer tokens usando jq (si disponible)
      if command -v jq >/dev/null 2>&1; then
        ACCESS_TOKEN=$(echo "$TOK_JSON" | jq -r '.access_token')
        REFRESH_TOKEN=$(echo "$TOK_JSON" | jq -r '.refresh_token')
      else
        ACCESS_TOKEN=$(echo "$TOK_JSON" | sed -n 's/.*"access_token":"\([^"]*\)".*/\1/p')
        REFRESH_TOKEN=$(echo "$TOK_JSON" | sed -n 's/.*"refresh_token":"\([^"]*\)".*/\1/p')
      fi

      if [ -n "$ACCESS_TOKEN" ]; then
        if [ "$PERSIST_TOKEN" = true ]; then
          TOKEN_FILE="$ROOT_DIR/.current_token"
          cat > "$TOKEN_FILE" <<EOF
ACCESS_TOKEN=$ACCESS_TOKEN
REFRESH_TOKEN=$REFRESH_TOKEN
SAVED_AT=$(date -u +%Y-%m-%dT%H:%M:%SZ)
EOF
          chmod 600 "$TOKEN_FILE"
          echo "Tokens persistidos en $TOKEN_FILE"
        fi
      else
        echo "No se pudo extraer access_token del output de httpyac"
      fi
    else
      echo "No se encontró salida JSON de tokens en el output de httpyac. Revisa auth.http post-request script."
    fi

    # Si solo pedimos login, terminamos aquí
    if [ "$LOGIN_ONLY" = true ]; then
      echo "Login-only completado."
      rm -f "$TMP_OUT"
      exit 0
    fi

    rm -f "$TMP_OUT"
  else
    echo "No se encontró $AUTH_FILE. Continuando sin login."
  fi
fi

# Si llegamos aquí, ejecutamos la suite completa (auth ya pudo haberse ejecutado arriba)
ARGS=()
for f in "$REQUESTS_DIR"/*.http; do
  ARGS+=("$f")
done

echo "Ejecutando httpyac para todos los archivos .http en $REQUESTS_DIR"
echo "DOTENV: $DOTENV_FILE"

httpyac "${ARGS[@]}" --dotenv "$DOTENV_FILE"

echo "Ejecución completada."
