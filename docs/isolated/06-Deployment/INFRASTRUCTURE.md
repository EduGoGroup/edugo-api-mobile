# Infraestructura
# Sistema de Evaluaciones - EduGo

**Versión:** 1.0.0  
**Fecha:** 14 de Noviembre, 2025

---

## 1. COMPONENTES DE INFRAESTRUCTURA

### Stack Tecnológico

- **API:** Go 1.21+ (Gin framework)
- **DB Relacional:** PostgreSQL 15+
- **DB Documental:** MongoDB 7.0+
- **Contenedores:** Docker + Docker Compose
- **CI/CD:** GitHub Actions
- **Monitoreo:** Prometheus + Grafana (Post-MVP)

### Puertos

| Servicio | Puerto | Protocolo |
|----------|--------|-----------|
| API Mobile | 8080 | HTTP |
| PostgreSQL | 5432 | TCP |
| MongoDB | 27017 | TCP |
| Prometheus | 9090 | HTTP |
| Grafana | 3000 | HTTP |

---

## 2. DOCKER COMPOSE COMPLETO

Ver `DEPLOYMENT_GUIDE.md` para docker-compose.yml completo.

---

## 3. ESCALADO

### Horizontal Scaling

- Load Balancer (Nginx)
- Múltiples instancias API (2-3)
- PostgreSQL: Single instance (MVP), replicación después
- MongoDB: Replica set (Post-MVP)

### Vertical Scaling

**PostgreSQL:**
- Recursos mínimos: 2 CPU, 4GB RAM
- Recomendado: 4 CPU, 8GB RAM

**MongoDB:**
- Recursos mínimos: 2 CPU, 4GB RAM

**API:**
- Por instancia: 1 CPU, 2GB RAM

---

## 4. BACKUPS

Ver `DEPLOYMENT_GUIDE.md` para procedimientos de backup.

**Frecuencia:**
- PostgreSQL: Diario (2 AM)
- MongoDB: Diario (3 AM)
- Retención: 30 días

---

**Generado con:** Claude Code
