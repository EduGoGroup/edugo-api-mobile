# Development Environment Troubleshooting

This guide helps you diagnose and fix common issues when setting up or running the EduGo API Mobile development environment.

## Table of Contents

- [Docker Issues](#docker-issues)
- [Container Issues](#container-issues)
- [Database Issues](#database-issues)
- [Network and Connectivity Issues](#network-and-connectivity-issues)
- [Port Conflicts](#port-conflicts)
- [Timeout Issues](#timeout-issues)
- [Data and Schema Issues](#data-and-schema-issues)
- [Reset Procedures](#reset-procedures)
- [Getting Help](#getting-help)

---

## Docker Issues

### Docker Daemon Not Running

**Symptoms:**
```
❌ Error: Docker daemon is not accessible
Cannot connect to the Docker daemon at unix:///var/run/docker.sock
```

**Solutions:**

1. **Start Docker Desktop** (macOS/Windows):
   - Open Docker Desktop application
   - Wait for the whale icon to stop animating
   - Verify with: `docker ps`

2. **Start Docker Service** (Linux):
   ```bash
   sudo systemctl start docker
   sudo systemctl enable docker  # Start on boot
   ```

3. **Check Docker Status**:
   ```bash
   # macOS/Linux
   docker info

   # If you see permission errors, add your user to docker group (Linux):
   sudo usermod -aG docker $USER
   newgrp docker  # Or logout and login again
   ```

### Docker Compose Not Found

**Symptoms:**
```
❌ Error: docker-compose command not found
```

**Solutions:**

1. **Install Docker Compose** (Linux):
   ```bash
   sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
   sudo chmod +x /usr/local/bin/docker-compose
   ```

2. **Verify Installation**:
   ```bash
   docker-compose --version
   ```

3. **Use Docker Compose V2** (alternative):
   ```bash
   # If you have Docker Compose V2, use:
   docker compose version

   # Update scripts to use 'docker compose' instead of 'docker-compose'
   ```

### .env File Missing

**Symptoms:**
```
❌ Error: .env file not found
```

**Solutions:**

1. **Copy from Example**:
   ```bash
   cp .env.example .env
   ```

2. **Verify File Exists**:
   ```bash
   ls -la .env
   ```

3. **Check File Permissions**:
   ```bash
   chmod 644 .env
   ```

---

## Container Issues

### Container Failed to Start

**Symptoms:**
```
❌ Error: Container edugo-postgres failed to start
```

**Diagnosis:**

```bash
# Check container status
docker ps -a | grep edugo

# View container logs
docker logs edugo-postgres
docker logs edugo-mongodb
docker logs edugo-rabbitmq
```

**Solutions:**

1. **Check for Port Conflicts** (see [Port Conflicts](#port-conflicts))

2. **Remove and Recreate Container**:
   ```bash
   docker stop edugo-postgres
   docker rm edugo-postgres
   make dev-init
   ```

3. **Check Docker Resources**:
   - Ensure Docker has enough memory (at least 4GB recommended)
   - Docker Desktop → Settings → Resources → Memory

4. **View Detailed Logs**:
   ```bash
   docker-compose -f docker-compose-local.yml logs -f
   ```

### Container Exits Immediately

**Symptoms:**
```
Container starts but exits with code 1 or 137
```

**Solutions:**

1. **Check Container Logs**:
   ```bash
   docker logs edugo-postgres --tail 50
   ```

2. **Verify Environment Variables**:
   ```bash
   # Check if .env file has correct values
   cat .env | grep POSTGRES
   ```

3. **Check Disk Space**:
   ```bash
   df -h
   # Ensure you have at least 5GB free
   ```

4. **Restart Docker**:
   ```bash
   # macOS/Windows: Restart Docker Desktop
   # Linux:
   sudo systemctl restart docker
   ```

### Container Health Check Failing

**Symptoms:**
```
❌ Error: Health check timeout after 60 seconds
Container is running but not healthy
```

**Solutions:**

1. **Check Container Health**:
   ```bash
   docker inspect edugo-postgres | grep -A 10 Health
   ```

2. **Manually Test Connection**:
   ```bash
   # PostgreSQL
   docker exec edugo-postgres pg_isready -U edugo

   # MongoDB
   docker exec edugo-mongodb mongosh --eval "db.adminCommand('ping')"

   # RabbitMQ
   curl -u edugo:edugo123 http://localhost:15672/api/overview
   ```

3. **Increase Timeout**:
   - Edit `scripts/dev-init.sh`
   - Increase timeout from 60 to 120 seconds in `wait_for_health()` function

4. **Check Container Resources**:
   ```bash
   docker stats edugo-postgres edugo-mongodb edugo-rabbitmq
   ```

---

## Database Issues

### PostgreSQL Connection Failed

**Symptoms:**
```
❌ Error: Could not connect to PostgreSQL
psql: error: connection to server at "localhost" (127.0.0.1), port 5432 failed
```

**Solutions:**

1. **Verify Container is Running**:
   ```bash
   docker ps | grep edugo-postgres
   ```

2. **Check PostgreSQL Logs**:
   ```bash
   docker logs edugo-postgres --tail 50
   ```

3. **Test Connection Manually**:
   ```bash
   psql "postgresql://edugo:edugo123@localhost:5432/edugo"
   ```

4. **Verify Port is Accessible**:
   ```bash
   nc -zv localhost 5432
   # Or
   telnet localhost 5432
   ```

5. **Check Credentials**:
   ```bash
   # Verify .env file has correct credentials
   cat .env | grep POSTGRES
   ```

### MongoDB Connection Failed

**Symptoms:**
```
❌ Error: Could not connect to MongoDB
MongoServerError: Authentication failed
```

**Solutions:**

1. **Verify Container is Running**:
   ```bash
   docker ps | grep edugo-mongodb
   ```

2. **Check MongoDB Logs**:
   ```bash
   docker logs edugo-mongodb --tail 50
   ```

3. **Test Connection with Correct Auth**:
   ```bash
   mongosh "mongodb://edugo:edugo123@localhost:27017/edugo?authSource=admin"
   ```

4. **Verify Authentication Database**:
   - MongoDB requires `authSource=admin` in connection string
   - Check that this is included in your connection attempts

### Schema Creation Failed

**Symptoms:**
```
❌ Error: Failed to execute migration script
ERROR: relation "users" already exists
```

**Solutions:**

1. **Check if Schema Already Exists**:
   ```bash
   psql "postgresql://edugo:edugo123@localhost:5432/edugo" -c "\dt"
   ```

2. **If Schema is Corrupted, Reset**:
   ```bash
   make dev-reset
   ```

3. **Manually Drop and Recreate**:
   ```bash
   # Connect to PostgreSQL
   psql "postgresql://edugo:edugo123@localhost:5432/edugo"

   # Drop all tables (careful!)
   DROP SCHEMA public CASCADE;
   CREATE SCHEMA public;

   # Exit and reinitialize
   \q
   make dev-init
   ```

### Seed Data Insertion Failed

**Symptoms:**
```
❌ Error: Failed to insert seed data
ERROR: duplicate key value violates unique constraint
```

**Solutions:**

1. **This is Usually Safe to Ignore**:
   - Seed data uses `ON CONFLICT DO NOTHING`
   - Duplicate key errors mean data already exists

2. **Verify Data Exists**:
   ```bash
   psql "postgresql://edugo:edugo123@localhost:5432/edugo" -c "SELECT COUNT(*) FROM users;"
   ```

3. **If Data is Corrupted**:
   ```bash
   # Reset environment
   make dev-reset
   ```

---

## Network and Connectivity Issues

### Cannot Reach Container from Host

**Symptoms:**
```
Connection refused when trying to connect to localhost:5432
```

**Solutions:**

1. **Check Port Mapping**:
   ```bash
   docker port edugo-postgres
   # Should show: 5432/tcp -> 0.0.0.0:5432
   ```

2. **Verify Container Network**:
   ```bash
   docker inspect edugo-postgres | grep IPAddress
   ```

3. **Try Container IP Directly**:
   ```bash
   # Get container IP
   CONTAINER_IP=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' edugo-postgres)

   # Test connection
   psql "postgresql://edugo:edugo123@${CONTAINER_IP}:5432/edugo"
   ```

4. **Check Firewall**:
   ```bash
   # macOS
   sudo /usr/libexec/ApplicationFirewall/socketfilterfw --getglobalstate

   # Linux
   sudo ufw status
   ```

### DNS Resolution Issues

**Symptoms:**
```
Could not resolve hostname
```

**Solutions:**

1. **Use IP Address Instead**:
   ```bash
   # Instead of localhost, use 127.0.0.1
   psql "postgresql://edugo:edugo123@127.0.0.1:5432/edugo"
   ```

2. **Check /etc/hosts**:
   ```bash
   cat /etc/hosts | grep localhost
   # Should contain: 127.0.0.1 localhost
   ```

3. **Restart Docker Network**:
   ```bash
   docker network prune
   docker-compose -f docker-compose-local.yml down
   docker-compose -f docker-compose-local.yml up -d
   ```

---

## Port Conflicts

### Port Already in Use

**Symptoms:**
```
❌ Error: Bind for 0.0.0.0:5432 failed: port is already allocated
```

**Solutions:**

1. **Find Process Using Port**:
   ```bash
   # macOS/Linux
   lsof -i :5432
   lsof -i :27017
   lsof -i :5672
   lsof -i :15672

   # Or use netstat
   netstat -an | grep LISTEN | grep 5432
   ```

2. **Stop Conflicting Process**:
   ```bash
   # If it's another PostgreSQL instance
   sudo systemctl stop postgresql

   # Or kill specific process
   kill -9 <PID>
   ```

3. **Change Port in Configuration**:
   - Edit `docker-compose-local.yml`
   - Change port mapping (e.g., `5433:5432` instead of `5432:5432`)
   - Update `.env` file with new port
   - Run `make dev-init`

4. **Use Different Ports**:
   ```yaml
   # Example: docker-compose-local.yml
   services:
     postgres:
       ports:
         - "5433:5432"  # Host:Container
   ```

---

## Timeout Issues

### Health Check Timeout

**Symptoms:**
```
❌ Error: Health check timeout after 60 seconds
Service did not become healthy in time
```

**Solutions:**

1. **Wait Longer**:
   - First startup can take 2-3 minutes
   - Subsequent startups are faster

2. **Check System Resources**:
   ```bash
   # Check CPU and memory usage
   docker stats

   # Check disk I/O
   iostat -x 1
   ```

3. **Increase Timeout**:
   - Edit `scripts/dev-init.sh`
   - Find `wait_for_health()` function
   - Increase timeout value

4. **Check Container Logs**:
   ```bash
   docker logs edugo-postgres --tail 100
   ```

### Connection Retry Timeout

**Symptoms:**
```
❌ Error: Connection failed after 3 retries
```

**Solutions:**

1. **Verify Service is Actually Running**:
   ```bash
   docker ps | grep edugo
   ```

2. **Check Service Logs**:
   ```bash
   docker logs edugo-postgres
   ```

3. **Manually Test Connection**:
   ```bash
   # PostgreSQL
   psql "postgresql://edugo:edugo123@localhost:5432/edugo"

   # MongoDB
   mongosh "mongodb://edugo:edugo123@localhost:27017/edugo?authSource=admin"
   ```

4. **Restart Services**:
   ```bash
   docker-compose -f docker-compose-local.yml restart
   ```

---

## Data and Schema Issues

### Duplicate Data After Multiple Runs

**Symptoms:**
```
Seeing duplicate records in database
```

**Solutions:**

1. **This Should Not Happen**:
   - Seed data uses `ON CONFLICT DO NOTHING`
   - Check if this clause is present in `02_seed_data.sql`

2. **Verify Idempotency**:
   ```bash
   # Check for duplicate users
   psql "postgresql://edugo:edugo123@localhost:5432/edugo" -c "SELECT email, COUNT(*) FROM users GROUP BY email HAVING COUNT(*) > 1;"
   ```

3. **Reset if Corrupted**:
   ```bash
   make dev-reset
   ```

### Missing Tables or Collections

**Symptoms:**
```
ERROR: relation "users" does not exist
```

**Solutions:**

1. **Run Initialization Again**:
   ```bash
   make dev-init
   ```

2. **Manually Run Migrations**:
   ```bash
   # PostgreSQL
   psql "postgresql://edugo:edugo123@localhost:5432/edugo" -f scripts/postgresql/01_create_schema.sql

   # MongoDB
   mongosh "mongodb://edugo:edugo123@localhost:27017/edugo?authSource=admin" scripts/mongodb/02_assessment_results.js
   ```

3. **Check Script Execution Logs**:
   - Review output from `make dev-init`
   - Look for error messages during schema creation

---

## Reset Procedures

### Complete Environment Reset

When things are really broken, start fresh:

```bash
# Full reset (will delete all data)
make dev-reset
```

This will:
1. Prompt for confirmation
2. Stop all containers
3. Remove all containers
4. Remove all volumes (deletes data)
5. Reinitialize environment

### Partial Reset Options

**Reset Only Containers** (keeps volumes):
```bash
docker-compose -f docker-compose-local.yml down
docker-compose -f docker-compose-local.yml up -d
```

**Reset Only PostgreSQL Data**:
```bash
docker-compose -f docker-compose-local.yml stop postgres
docker volume rm edugo-api-mobile_postgres_data
docker-compose -f docker-compose-local.yml up -d postgres
make dev-init
```

**Reset Only MongoDB Data**:
```bash
docker-compose -f docker-compose-local.yml stop mongodb
docker volume rm edugo-api-mobile_mongodb_data
docker-compose -f docker-compose-local.yml up -d mongodb
make dev-init
```

**Reset Only RabbitMQ**:
```bash
docker-compose -f docker-compose-local.yml stop rabbitmq
docker volume rm edugo-api-mobile_rabbitmq_data
docker-compose -f docker-compose-local.yml up -d rabbitmq
make dev-init
```

### Manual Cleanup

If `make dev-reset` fails:

```bash
# Stop all containers
docker stop $(docker ps -aq --filter "name=edugo")

# Remove all containers
docker rm $(docker ps -aq --filter "name=edugo")

# Remove all volumes
docker volume ls | grep edugo | awk '{print $2}' | xargs docker volume rm

# Remove network
docker network rm edugo-network

# Start fresh
make dev-init
```

---

## Getting Help

### Diagnostic Commands

Run these commands to gather information for troubleshooting:

```bash
# System information
docker version
docker-compose version
uname -a

# Container status
docker ps -a | grep edugo

# Container logs
docker logs edugo-postgres --tail 50
docker logs edugo-mongodb --tail 50
docker logs edugo-rabbitmq --tail 50

# Network information
docker network ls
docker network inspect bridge

# Volume information
docker volume ls | grep edugo

# Resource usage
docker stats --no-stream

# Environment status
make dev-status
```

### Log Locations

**Container Logs**:
```bash
# View logs
docker logs edugo-postgres
docker logs edugo-mongodb
docker logs edugo-rabbitmq

# Follow logs in real-time
docker logs -f edugo-postgres

# View last N lines
docker logs --tail 100 edugo-postgres
```

**Application Logs**:
```bash
# If running the API
tail -f logs/app.log
```

### Useful Commands Reference

```bash
# Check environment status
make dev-status

# View all containers
docker ps -a

# View container logs
docker logs <container-name>

# Execute command in container
docker exec -it edugo-postgres psql -U edugo -d edugo

# Restart specific container
docker restart edugo-postgres

# Stop all containers
docker-compose -f docker-compose-local.yml down

# Start all containers
docker-compose -f docker-compose-local.yml up -d

# View container resource usage
docker stats

# Clean up unused resources
docker system prune -a
```

### When to Ask for Help

If you've tried the solutions above and still have issues:

1. **Gather diagnostic information** using commands above
2. **Check existing issues** in the project repository
3. **Create a new issue** with:
   - Description of the problem
   - Steps to reproduce
   - Error messages
   - Output from diagnostic commands
   - Your environment (OS, Docker version, etc.)

### Additional Resources

- [Environment Setup Guide](./ENVIRONMENT_SETUP.md)
- [Credentials Reference](./CREDENTIALS.md)
- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Troubleshooting](https://docs.docker.com/compose/faq/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [MongoDB Documentation](https://docs.mongodb.com/)
- [RabbitMQ Documentation](https://www.rabbitmq.com/documentation.html)
