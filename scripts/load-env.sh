#!/bin/bash
# Script para cargar variables de .env y ejecutar un comando
# Uso: ./scripts/load-env.sh <comando>

set -e

# Cargar .env si existe
if [ -f .env ]; then
    echo "📝 Cargando variables de .env..."
    set -a  # Exportar automáticamente todas las variables
    source .env
    set +a
    echo "✅ Variables cargadas"
else
    echo "⚠️  Archivo .env no encontrado"
    exit 1
fi

# Ejecutar el comando pasado como argumentos
exec "$@"
