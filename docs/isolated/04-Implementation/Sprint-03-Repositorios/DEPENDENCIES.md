# Dependencias del Sprint 03 - Repositorios

## Dependencias Técnicas Previas

### PostgreSQL 15+
- [ ] PostgreSQL corriendo en puerto 5432
- [ ] Schema de Sprint-01 ejecutado (tablas assessment*)

```bash
pg_isready -h localhost -p 5432
psql -U postgres -d edugo_test -c "\dt assessment*"
```

### MongoDB 7.0+
- [ ] MongoDB corriendo en puerto 27017
- [ ] Colección `material_assessment` con datos

```bash
mongosh --eval "db.runCommand({ ping: 1 })"
mongosh --eval "db.material_assessment.countDocuments()"
```

### Docker
- [ ] Docker corriendo (para Testcontainers)

```bash
docker --version
docker ps
```

### Sprint-02 Completado
- [ ] Entities creadas
- [ ] Repository interfaces definidas

```bash
ls -la /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/internal/domain/entities/
ls -la /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile/internal/domain/repositories/
```

---

## Dependencias de Código

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# GORM
go get gorm.io/gorm@v1.25.5
go get gorm.io/driver/postgres@v1.5.4

# MongoDB Driver
go get go.mongodb.org/mongo-driver/mongo@v1.13.1
go get go.mongodb.org/mongo-driver/bson@v1.13.1

# Testcontainers
go get github.com/testcontainers/testcontainers-go@v0.27.0
go get github.com/testcontainers/testcontainers-go/modules/postgres@v0.27.0
go get github.com/testcontainers/testcontainers-go/modules/mongodb@v0.27.0

go mod tidy
```

---

## Variables de Entorno

```bash
# PostgreSQL
export DB_HOST="localhost"
export DB_PORT="5432"
export DB_USER="postgres"
export DB_PASSWORD="postgres"
export DB_NAME="edugo_test"

# MongoDB
export MONGO_URI="mongodb://localhost:27017"
export MONGO_DATABASE="edugo_test"

# Connection Pool
export DB_MAX_OPEN_CONNS="25"
export DB_MAX_IDLE_CONNS="10"
export DB_CONN_MAX_LIFETIME="5m"
```

---

**Generado con:** Claude Code  
**Sprint:** 03/06
