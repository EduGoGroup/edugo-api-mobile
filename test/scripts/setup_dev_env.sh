#!/bin/bash
# setup_dev_env.sh - Configurar ambiente de desarrollo con Docker
#
# Este script levanta contenedores Docker para PostgreSQL, MongoDB y RabbitMQ,
# carga el schema de base de datos y configura la topología de RabbitMQ.

set -euo pipefail

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔══════════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║                                                          ║${NC}"
echo -e "${BLUE}║     🚀 EduGo API Mobile - Development Environment       ║${NC}"
echo -e "${BLUE}║                                                          ║${NC}"
echo -e "${BLUE}╚══════════════════════════════════════════════════════════╝${NC}"
echo ""

# 1. Verificar que Docker está corriendo
echo -e "${BLUE}🐳 Verificando Docker...${NC}"
if ! docker ps > /dev/null 2>&1; then
    echo -e "${RED}❌ Error: Docker no está corriendo${NC}"
    echo "Por favor inicia Docker Desktop y vuelve a intentar"
    exit 1
fi
echo -e "${GREEN}✓ Docker está corriendo${NC}"
echo ""

# 2. Levantar contenedores con docker-compose
echo -e "${BLUE}🐳 Levantando contenedores...${NC}"
docker-compose -f docker-compose-dev.yml up -d

# 3. Esperar a que los servicios estén listos
echo -e "${BLUE}⏳ Esperando a que los servicios estén listos...${NC}"
echo -n "  PostgreSQL: "
timeout 60 bash -c 'until docker exec edugo-postgres-dev pg_isready -U edugo_user -d edugo > /dev/null 2>&1; do sleep 1; done' && echo -e "${GREEN}✓${NC}" || (echo -e "${RED}✗${NC}" && exit 1)

echo -n "  MongoDB: "
timeout 60 bash -c 'until docker exec edugo-mongo-dev mongosh --eval "db.adminCommand({ping: 1})" > /dev/null 2>&1; do sleep 1; done' && echo -e "${GREEN}✓${NC}" || (echo -e "${RED}✗${NC}" && exit 1)

echo -n "  RabbitMQ: "
timeout 60 bash -c 'until docker exec edugo-rabbitmq-dev rabbitmq-diagnostics -q ping > /dev/null 2>&1; do sleep 1; done' && echo -e "${GREEN}✓${NC}" || (echo -e "${RED}✗${NC}" && exit 1)

echo ""

# 4. Ejecutar schema SQL en PostgreSQL (si existe)
if [ -f "scripts/postgresql/schema.sql" ]; then
    echo -e "${BLUE}🗄️  Ejecutando schema SQL en PostgreSQL...${NC}"
    docker exec -i edugo-postgres-dev psql -U edugo_user -d edugo < scripts/postgresql/schema.sql
    echo -e "${GREEN}✓ Schema creado${NC}"
else
    echo -e "${YELLOW}⚠️  Advertencia: scripts/postgresql/schema.sql no encontrado${NC}"
    echo "   El schema debe crearse manualmente o mediante migraciones"
fi
echo ""

# 5. Crear colecciones e índices en MongoDB
echo -e "${BLUE}🍃 Configurando MongoDB...${NC}"
docker exec edugo-mongo-dev mongosh edugo --eval "
    // Crear colecciones
    db.createCollection('material_assessments');
    db.createCollection('assessment_results');
    db.createCollection('assessment_attempts');

    // Crear índices
    db.material_assessments.createIndex({ material_id: 1 }, { unique: true });
    db.assessment_results.createIndex({ assessment_id: 1, user_id: 1 }, { unique: true });
    db.assessment_attempts.createIndex({ user_id: 1, assessment_id: 1 });

    print('✓ Colecciones e índices creados');
" > /dev/null 2>&1
echo -e "${GREEN}✓ MongoDB configurado${NC}"
echo ""

# 6. Configurar topología de RabbitMQ
echo -e "${BLUE}🐰 Configurando topología de RabbitMQ...${NC}"

# Esperar un poco más para que RabbitMQ esté completamente listo
sleep 5

# Crear exchange
docker exec edugo-rabbitmq-dev rabbitmqadmin declare exchange name=edugo.events type=topic durable=true > /dev/null 2>&1 || true

# Crear colas
declare -a queues=("material.created" "material.updated" "material.deleted" "assessment.completed" "progress.updated" "user.registered")

for queue in "${queues[@]}"; do
    docker exec edugo-rabbitmq-dev rabbitmqadmin declare queue name="$queue" durable=true > /dev/null 2>&1 || true
    docker exec edugo-rabbitmq-dev rabbitmqadmin declare binding source=edugo.events destination="$queue" routing_key="$queue" > /dev/null 2>&1 || true
done

echo -e "${GREEN}✓ RabbitMQ configurado (exchange + 6 colas)${NC}"
echo ""

# 7. Mostrar información de conexión
echo -e "${GREEN}╔══════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║                                                          ║${NC}"
echo -e "${GREEN}║     ✅ Ambiente de desarrollo listo!                     ║${NC}"
echo -e "${GREEN}║                                                          ║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${BLUE}📝 Connection Strings:${NC}"
echo ""
echo -e "  ${YELLOW}PostgreSQL:${NC}"
echo "    postgresql://edugo_user:edugo_pass@localhost:5432/edugo?sslmode=disable"
echo ""
echo -e "  ${YELLOW}MongoDB:${NC}"
echo "    mongodb://edugo_admin:edugo_pass@localhost:27017/edugo"
echo ""
echo -e "  ${YELLOW}RabbitMQ AMQP:${NC}"
echo "    amqp://edugo_user:edugo_pass@localhost:5672/"
echo ""
echo -e "${BLUE}🌐 Web Interfaces:${NC}"
echo ""
echo -e "  ${YELLOW}RabbitMQ Management:${NC} http://localhost:15672"
echo "    Usuario: edugo_user"
echo "    Password: edugo_pass"
echo ""
echo -e "${BLUE}🛠️  Comandos útiles:${NC}"
echo ""
echo "  Ver logs:          docker-compose -f docker-compose-dev.yml logs -f"
echo "  Detener:           docker-compose -f docker-compose-dev.yml stop"
echo "  Reiniciar:         docker-compose -f docker-compose-dev.yml restart"
echo "  Limpiar todo:      ./test/scripts/teardown_dev_env.sh"
echo ""
