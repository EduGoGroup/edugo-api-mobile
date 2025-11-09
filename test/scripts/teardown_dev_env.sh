#!/bin/bash
# teardown_dev_env.sh - Limpiar ambiente de desarrollo
#
# Este script detiene y elimina todos los contenedores Docker del ambiente
# de desarrollo, incluyendo volúmenes de datos.

set -euo pipefail

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧹 Limpiando ambiente de desarrollo...${NC}"
echo ""

# Preguntar confirmación si hay flag -v (eliminar volúmenes)
REMOVE_VOLUMES=""
if [[ "${1:-}" == "-v" ]] || [[ "${1:-}" == "--volumes" ]]; then
    echo -e "${YELLOW}⚠️  Advertencia: Se eliminarán todos los datos de los volúmenes${NC}"
    read -p "¿Estás seguro? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        REMOVE_VOLUMES="-v"
    else
        echo "Operación cancelada"
        exit 0
    fi
fi

# Detener y eliminar contenedores
echo -e "${BLUE}🐳 Deteniendo contenedores...${NC}"
docker-compose -f docker-compose-dev.yml down $REMOVE_VOLUMES

if [[ -n "$REMOVE_VOLUMES" ]]; then
    echo -e "${GREEN}✓ Contenedores y volúmenes eliminados${NC}"
else
    echo -e "${GREEN}✓ Contenedores detenidos (volúmenes preservados)${NC}"
fi

echo ""
echo -e "${GREEN}╔══════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                                                          ║${NC}"
echo -e "${GREEN}║     ✅ Ambiente de desarrollo limpiado                   ║${NC}"
echo -e "${GREEN}║                                                          ║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════════════════════════╝${NC}"
echo ""

if [[ -z "$REMOVE_VOLUMES" ]]; then
    echo -e "${BLUE}💡 Tip:${NC}"
    echo "  Para eliminar también los datos, usa:"
    echo "  ./test/scripts/teardown_dev_env.sh --volumes"
    echo ""
fi
