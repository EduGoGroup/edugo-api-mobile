# Development Environment Credentials

This document contains all default credentials for the EduGo API Mobile development environment.

⚠️ **SECURITY WARNING**: These credentials are for **DEVELOPMENT ONLY**. Never use these credentials in production environments.

---

## Table of Contents

- [Service Credentials](#service-credentials)
  - [PostgreSQL](#postgresql)
  - [MongoDB](#mongodb)
  - [RabbitMQ](#rabbitmq)
- [Test User Accounts](#test-user-accounts)
  - [Teacher Accounts](#teacher-accounts)
  - [Student Accounts](#student-accounts)
- [Connection Strings](#connection-strings)
- [Security Notes](#security-notes)

---

## Service Credentials

### PostgreSQL

PostgreSQL is used as the primary relational database for user data, materials, and progress tracking.

**Connection Details:**
```
Host:     localhost
Port:     5432
Database: edugo
Username: edugo
Password: edugo123
```

**Connection String:**
```
postgresql://edugo:edugo123@localhost:5432/edugo
```

**psql Command:**
```bash
psql "postgresql://edugo:edugo123@localhost:5432/edugo"
```

**Docker Exec:**
```bash
docker exec -it edugo-postgres psql -U edugo -d edugo
```

**Environment Variables:**
```bash
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=edugo
POSTGRES_USER=edugo
POSTGRES_PASSWORD=edugo123
```

---

### MongoDB

MongoDB is used for storing assessment results and analytics data.

**Connection Details:**
```
Host:          localhost
Port:          27017
Database:      edugo
Username:      edugo
Password:      edugo123
Auth Source:   admin
```

**Connection String:**
```
mongodb://edugo:edugo123@localhost:27017/edugo?authSource=admin
```

**mongosh Command:**
```bash
mongosh "mongodb://edugo:edugo123@localhost:27017/edugo?authSource=admin"
```

**Docker Exec:**
```bash
docker exec -it edugo-mongodb mongosh -u edugo -p edugo123 --authenticationDatabase admin edugo
```

**Environment Variables:**
```bash
MONGODB_HOST=localhost
MONGODB_PORT=27017
MONGODB_DATABASE=edugo
MONGODB_USER=edugo
MONGODB_PASSWORD=edugo123
MONGODB_AUTH_SOURCE=admin
```

**Important Notes:**
- Always include `authSource=admin` in connection strings
- The authentication database is `admin`, not `edugo`

---

### RabbitMQ

RabbitMQ is used for asynchronous messaging and event-driven architecture.

**Connection Details:**
```
Host:              localhost
AMQP Port:         5672
Management Port:   15672
Username:          edugo
Password:          edugo123
Virtual Host:      /
```

**AMQP Connection String:**
```
amqp://edugo:edugo123@localhost:5672/
```

**Management UI:**
```
URL:      http://localhost:15672
Username: edugo
Password: edugo123
```

**Docker Exec:**
```bash
docker exec -it edugo-rabbitmq rabbitmqctl list_queues
```

**Environment Variables:**
```bash
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_MANAGEMENT_PORT=15672
RABBITMQ_USER=edugo
RABBITMQ_PASSWORD=edugo123
RABBITMQ_VHOST=/
```

**Management UI Access:**
1. Open browser: `http://localhost:15672`
2. Login with username `edugo` and password `edugo123`
3. View exchanges, queues, and messages

---

## Test User Accounts

The development environment includes pre-configured test user accounts for testing authentication and authorization.

### Teacher Accounts

#### Teacher 1: Juan Pérez

```
Email:    juan.perez@edugo.com
Password: password123
Role:     teacher
User ID:  11111111-1111-1111-1111-111111111111
```

**Usage Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "juan.perez@edugo.com",
    "password": "password123"
  }'
```

**Profile:**
- Full Name: Juan Pérez
- Role: Teacher
- Can create and manage materials
- Can view student progress
- Can create assessments

---

#### Teacher 2: María González

```
Email:    maria.gonzalez@edugo.com
Password: password123
Role:     teacher
User ID:  22222222-2222-2222-2222-222222222222
```

**Usage Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "maria.gonzalez@edugo.com",
    "password": "password123"
  }'
```

**Profile:**
- Full Name: María González
- Role: Teacher
- Can create and manage materials
- Can view student progress
- Can create assessments

---

### Student Accounts

#### Student 1: Carlos Rodríguez

```
Email:    carlos.rodriguez@student.edugo.com
Password: password123
Role:     student
User ID:  33333333-3333-3333-3333-333333333333
```

**Usage Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "carlos.rodriguez@student.edugo.com",
    "password": "password123"
  }'
```

**Profile:**
- Full Name: Carlos Rodríguez
- Role: Student
- Can view assigned materials
- Can complete assessments
- Can track own progress

---

#### Student 2: Ana Martínez

```
Email:    ana.martinez@student.edugo.com
Password: password123
Role:     student
User ID:  44444444-4444-4444-4444-444444444444
```

**Usage Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "ana.martinez@student.edugo.com",
    "password": "password123"
  }'
```

**Profile:**
- Full Name: Ana Martínez
- Role: Student
- Can view assigned materials
- Can complete assessments
- Can track own progress

---

#### Student 3: Luis Fernández

```
Email:    luis.fernandez@student.edugo.com
Password: password123
Role:     student
User ID:  55555555-5555-5555-5555-555555555555
```

**Usage Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "luis.fernandez@student.edugo.com",
    "password": "password123"
  }'
```

**Profile:**
- Full Name: Luis Fernández
- Role: Student
- Can view assigned materials
- Can complete assessments
- Can track own progress

---

## Connection Strings

### Quick Reference

**PostgreSQL:**
```
postgresql://edugo:edugo123@localhost:5432/edugo
```

**MongoDB:**
```
mongodb://edugo:edugo123@localhost:27017/edugo?authSource=admin
```

**RabbitMQ:**
```
amqp://edugo:edugo123@localhost:5672/
```

**RabbitMQ Management:**
```
http://localhost:15672
```

### Using in Application Code

**Go Example (PostgreSQL):**
```go
import "database/sql"

connStr := "postgresql://edugo:edugo123@localhost:5432/edugo"
db, err := sql.Open("postgres", connStr)
```

**Go Example (MongoDB):**
```go
import "go.mongodb.org/mongo-driver/mongo"

uri := "mongodb://edugo:edugo123@localhost:27017/edugo?authSource=admin"
client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
```

**Go Example (RabbitMQ):**
```go
import "github.com/streadway/amqp"

conn, err := amqp.Dial("amqp://edugo:edugo123@localhost:5672/")
```

### Using in Environment Variables

Add these to your `.env` file:

```bash
# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=edugo
POSTGRES_USER=edugo
POSTGRES_PASSWORD=edugo123

# MongoDB
MONGODB_HOST=localhost
MONGODB_PORT=27017
MONGODB_DATABASE=edugo
MONGODB_USER=edugo
MONGODB_PASSWORD=edugo123
MONGODB_AUTH_SOURCE=admin

# RabbitMQ
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=edugo
RABBITMQ_PASSWORD=edugo123

# JWT (for authentication)
JWT_SECRET=dev-secret-key-change-in-production
JWT_EXPIRATION=24h
```

---

## Security Notes

### Development vs Production

⚠️ **CRITICAL**: These credentials are **ONLY** for local development.

**DO NOT:**
- ❌ Use these credentials in production
- ❌ Commit `.env` files with real credentials
- ❌ Share production credentials in documentation
- ❌ Use simple passwords like "password123" in production
- ❌ Use default service passwords in production

**DO:**
- ✅ Use strong, unique passwords in production
- ✅ Use environment variables for configuration
- ✅ Use secrets management tools (AWS Secrets Manager, HashiCorp Vault, etc.)
- ✅ Rotate credentials regularly in production
- ✅ Use different credentials for each environment
- ✅ Enable SSL/TLS for database connections in production
- ✅ Restrict network access to databases in production

### Password Hashing

Test user passwords are hashed using bcrypt before storage:

```
Plaintext:  password123
Algorithm:  bcrypt
Cost:       10
Hash:       $2a$10$... (stored in database)
```

**Never store plaintext passwords in the database**, even in development.

### JWT Secrets

The development JWT secret is:
```
JWT_SECRET=dev-secret-key-change-in-production
```

**In production:**
- Use a cryptographically secure random string (at least 32 characters)
- Never commit JWT secrets to version control
- Rotate JWT secrets periodically
- Use different secrets for different environments

### Network Security

**Development:**
- Services bind to `localhost` (127.0.0.1)
- Only accessible from your local machine
- No SSL/TLS required

**Production:**
- Use private networks or VPCs
- Enable SSL/TLS for all connections
- Use firewalls to restrict access
- Enable authentication on all services
- Use connection pooling with authentication

### Credential Rotation

**Development:**
- Credentials can remain static
- Reset environment if compromised: `make dev-reset`

**Production:**
- Rotate credentials every 90 days minimum
- Rotate immediately if compromised
- Use automated rotation tools
- Maintain audit logs of credential access

### Best Practices

1. **Separation of Environments**
   - Use different credentials for dev, staging, and production
   - Never connect development tools to production databases

2. **Least Privilege**
   - Grant only necessary permissions
   - Use read-only accounts where possible
   - Separate admin and application accounts

3. **Monitoring**
   - Log authentication attempts
   - Monitor for unusual access patterns
   - Set up alerts for failed login attempts

4. **Backup Security**
   - Encrypt database backups
   - Secure backup credentials separately
   - Test backup restoration regularly

### Checking for Exposed Credentials

Before committing code, check for exposed credentials:

```bash
# Check for common credential patterns
grep -r "password123" .
grep -r "edugo123" .
grep -r "JWT_SECRET" .

# Use git-secrets or similar tools
git secrets --scan
```

### If Credentials Are Compromised

**Development:**
1. Run `make dev-reset` to reset environment
2. Change credentials in `.env` file
3. Restart services

**Production:**
1. Immediately rotate all credentials
2. Review access logs for unauthorized access
3. Notify security team
4. Update all applications with new credentials
5. Investigate how credentials were compromised

---

## Quick Reference Card

Print this for easy reference:

```
┌─────────────────────────────────────────────────────────────┐
│                  DEVELOPMENT CREDENTIALS                     │
├─────────────────────────────────────────────────────────────┤
│ PostgreSQL                                                   │
│   postgresql://edugo:edugo123@localhost:5432/edugo         │
│                                                              │
│ MongoDB                                                      │
│   mongodb://edugo:edugo123@localhost:27017/edugo?authSource=admin │
│                                                              │
│ RabbitMQ                                                     │
│   amqp://edugo:edugo123@localhost:5672/                    │
│   http://localhost:15672 (Management UI)                    │
│                                                              │
│ Test Users (all use password: password123)                  │
│   juan.perez@edugo.com              (teacher)               │
│   maria.gonzalez@edugo.com          (teacher)               │
│   carlos.rodriguez@student.edugo.com (student)              │
│   ana.martinez@student.edugo.com     (student)              │
│   luis.fernandez@student.edugo.com   (student)              │
│                                                              │
│ ⚠️  DEVELOPMENT ONLY - DO NOT USE IN PRODUCTION             │
└─────────────────────────────────────────────────────────────┘
```

---

## Additional Resources

- [Environment Setup Guide](./ENVIRONMENT_SETUP.md)
- [Troubleshooting Guide](./TROUBLESHOOTING.md)
- [Security Best Practices](https://owasp.org/www-project-top-ten/)
- [PostgreSQL Security](https://www.postgresql.org/docs/current/auth-methods.html)
- [MongoDB Security](https://docs.mongodb.com/manual/security/)
- [RabbitMQ Access Control](https://www.rabbitmq.com/access-control.html)
