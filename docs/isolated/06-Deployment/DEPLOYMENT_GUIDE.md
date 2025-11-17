# Guía de Deployment
# Sistema de Evaluaciones - EduGo

**Versión:** 1.0.0  
**Fecha:** 14 de Noviembre, 2025

---

## 1. PRE-REQUISITOS

### Infraestructura Requerida

- [ ] **PostgreSQL 15+** instalado y accesible
- [ ] **MongoDB 7.0+** instalado y accesible
- [ ] **Servidor/VM** con Go 1.21+ (si deployment manual)
- [ ] **Docker** (si deployment con contenedores)

### Variables de Entorno Configuradas

```bash
# PostgreSQL
export DB_HOST="db.produccion.com"
export DB_PORT="5432"
export DB_USER="edugo_app"
export DB_PASSWORD="<secret>"
export DB_NAME="edugo_prod"
export DB_SSLMODE="require"

# MongoDB
export MONGO_URI="mongodb://mongo.produccion.com:27017"
export MONGO_DATABASE="edugo_prod"
export MONGO_USERNAME="edugo_app"
export MONGO_PASSWORD="<secret>"

# API
export PORT="8080"
export GIN_MODE="release"
export JWT_SECRET_KEY="<secret>"

# Logging
export LOG_LEVEL="info"
export LOG_FORMAT="json"
```

---

## 2. PASOS DE DEPLOYMENT

### Paso 1: Ejecutar Migraciones SQL

```bash
# Conectar a PostgreSQL de producción
psql -h db.produccion.com -U edugo_app -d edugo_prod

# Ejecutar migración de assessments
\i /path/to/06_assessments.sql

# Verificar tablas creadas
\dt assessment*

# Verificar índices
\di assessment*

# Salir
\q
```

**Validación:**
```bash
psql -h db.produccion.com -U edugo_app -d edugo_prod -c "SELECT COUNT(*) FROM assessment;"
# Esperado: 0 (tabla vacía pero creada)
```

### Paso 2: Build de Aplicación

**Opción A: Build Manual**
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Build con optimizaciones
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/api-mobile \
    -ldflags="-w -s" \
    ./cmd/api

# Verificar binario
ls -lh bin/api-mobile
# Esperado: ~15-25 MB

# Ejecutar localmente para verificar
./bin/api-mobile --check-config
```

**Opción B: Build con Docker**
```bash
# Build imagen
docker build -t ghcr.io/edugogroup/edugo-api-mobile:latest .

# Verificar imagen
docker images | grep edugo-api-mobile

# Test local
docker run -d -p 8080:8080 \
    -e DB_HOST=db.produccion.com \
    -e DB_USER=edugo_app \
    --name api-mobile-test \
    ghcr.io/edugogroup/edugo-api-mobile:latest

# Health check
curl http://localhost:8080/health

# Cleanup
docker stop api-mobile-test
docker rm api-mobile-test
```

### Paso 3: Deploy del Binario/Contenedor

**Deployment Manual (VM/Server):**
```bash
# Copiar binario a servidor
scp bin/api-mobile user@servidor:/opt/edugo/api-mobile

# SSH al servidor
ssh user@servidor

# Crear servicio systemd
sudo nano /etc/systemd/system/edugo-api-mobile.service

# Contenido:
[Unit]
Description=EduGo API Mobile
After=network.target

[Service]
Type=simple
User=edugo
WorkingDirectory=/opt/edugo
ExecStart=/opt/edugo/api-mobile
Restart=always
EnvironmentFile=/opt/edugo/.env

[Install]
WantedBy=multi-user.target

# Iniciar servicio
sudo systemctl daemon-reload
sudo systemctl enable edugo-api-mobile
sudo systemctl start edugo-api-mobile
sudo systemctl status edugo-api-mobile
```

**Deployment Docker (Compose):**
```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  api-mobile:
    image: ghcr.io/edugogroup/edugo-api-mobile:latest
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - MONGO_URI=mongodb://mongo:27017
    depends_on:
      - postgres
      - mongo
    restart: always

  postgres:
    image: postgres:15-alpine
    # ... config

  mongo:
    image: mongo:7-alpine
    # ... config
```

```bash
docker-compose -f docker-compose.prod.yml up -d
```

### Paso 4: Health Checks

```bash
# Health endpoint
curl http://api.produccion.com:8080/health
# Esperado: {"status":"healthy"}

# Verificar conexión a PostgreSQL
curl http://api.produccion.com:8080/health/db
# Esperado: {"postgres":"connected","mongodb":"connected"}

# Verificar endpoints de assessments
curl -H "Authorization: Bearer $TOKEN" \
    http://api.produccion.com:8080/v1/materials/test-id/assessment
# Esperado: 200 o 404 (pero no 500)
```

### Paso 5: Smoke Tests

```bash
# Test básico de flujo
# 1. Obtener assessment
ASSESSMENT=$(curl -H "Authorization: Bearer $TOKEN" \
    http://api.produccion.com:8080/v1/materials/$MATERIAL_ID/assessment)

echo $ASSESSMENT | jq .

# 2. Crear attempt
ATTEMPT=$(curl -X POST -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"answers":[...]}' \
    http://api.produccion.com:8080/v1/materials/$MATERIAL_ID/assessment/attempts)

echo $ATTEMPT | jq .

# 3. Verificar score retornado
echo $ATTEMPT | jq .score
# Esperado: Número entre 0-100
```

---

## 3. ROLLBACK PROCEDURE

### Si Deployment Falla

**Paso 1: Revertir API**
```bash
# Systemd
sudo systemctl stop edugo-api-mobile
sudo systemctl start edugo-api-mobile-old

# Docker
docker-compose down
docker-compose -f docker-compose.old.yml up -d
```

**Paso 2: Revertir Migraciones (SI ES NECESARIO)**
```bash
# Ejecutar rollback SQL
psql -h db.produccion.com -U edugo_app -d edugo_prod < /path/to/06_assessments_rollback.sql

# ⚠️ SOLO si migración causó problemas
# ⚠️ CUIDADO: Elimina datos
```

**Paso 3: Verificar Estado**
```bash
# Verificar que API vieja responde
curl http://api.produccion.com:8080/health

# Verificar logs
sudo journalctl -u edugo-api-mobile -n 100
```

---

## 4. MONITOREO POST-DEPLOYMENT

### Logs a Revisar

```bash
# Systemd logs
sudo journalctl -u edugo-api-mobile -f

# Docker logs
docker logs -f api-mobile

# Buscar errores
sudo journalctl -u edugo-api-mobile | grep ERROR
```

### Métricas a Monitorear (Primeras 24h)

- **Error rate:** Debe ser <1%
- **Latencia p95:** <2s
- **Requests/segundo:** Según carga esperada
- **Database connections:** No debe saturar pool

---

**Generado con:** Claude Code
