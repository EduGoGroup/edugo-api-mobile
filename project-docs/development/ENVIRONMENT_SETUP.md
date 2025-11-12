# Development Environment Setup

This guide will help you set up your local development environment for the EduGo API Mobile project.

## Prerequisites

Before you begin, ensure you have the following installed on your system:

### Required Software

1. **Docker Desktop** (version 20.10 or higher)
   - macOS: [Download Docker Desktop for Mac](https://www.docker.com/products/docker-desktop)
   - Linux: [Install Docker Engine](https://docs.docker.com/engine/install/)
   - Windows: [Download Docker Desktop for Windows](https://www.docker.com/products/docker-desktop)

2. **Docker Compose** (version 2.0 or higher)
   - Included with Docker Desktop on macOS and Windows
   - Linux: Install separately following [Docker Compose installation guide](https://docs.docker.com/compose/install/)

3. **Make** (usually pre-installed on macOS and Linux)
   - macOS: Included with Xcode Command Line Tools
   - Linux: Install via package manager (`apt-get install make` or `yum install make`)
   - Windows: Install via [GnuWin32](http://gnuwin32.sourceforge.net/packages/make.htm) or use WSL2

### Verify Prerequisites

Run these commands to verify your installation:

```bash
# Check Docker
docker --version
# Expected: Docker version 20.10.x or higher

# Check Docker Compose
docker-compose --version
# Expected: Docker Compose version 2.x.x or higher

# Check Make
make --version
# Expected: GNU Make 3.x or higher

# Verify Docker daemon is running
docker ps
# Should return a list of containers (may be empty)
```

## Quick Start

The fastest way to set up your development environment is using the automated initialization script:

```bash
# Clone the repository (if you haven't already)
git clone <repository-url>
cd edugo-api-mobile

# Initialize the development environment
make dev-init
```

That's it! The script will:
- ✅ Validate all prerequisites
- ✅ Start required Docker containers (PostgreSQL, MongoDB, RabbitMQ)
- ✅ Wait for services to be healthy
- ✅ Create database schemas
- ✅ Insert seed data with test users
- ✅ Verify connectivity to all services
- ✅ Display connection strings and credentials

The entire process takes approximately 2-3 minutes on first run.

## Detailed Setup Steps

If you want to understand what happens during initialization, here's a breakdown of each step:

### Step 1: Prerequisites Validation

The script first validates that your environment is ready:

```bash
# The script checks:
- Docker daemon is running and accessible
- docker-compose command is available
- .env file exists in the project root
```

If any prerequisite is missing, you'll see a clear error message with instructions on how to fix it.

### Step 2: Container Management

The script manages Docker containers using `docker-compose-local.yml`:

```bash
# Containers started:
- edugo-postgres (PostgreSQL 15)
- edugo-mongodb (MongoDB 6.0)
- edugo-rabbitmq (RabbitMQ 3.12 with management plugin)
```

If containers already exist and are running, the script will skip this step.

### Step 3: Health Checks

The script waits for all services to be healthy before proceeding:

```bash
# Health checks with 60-second timeout for each service:
- PostgreSQL: Accepts connections on port 5432
- MongoDB: Accepts connections on port 27017
- RabbitMQ: Management API responds on port 15672
```

### Step 4: Database Schema Creation

PostgreSQL schema is created by executing migration scripts in order:

```bash
# Scripts executed (if schema doesn't exist):
1. 01_create_schema.sql      # Base tables (users, materials, etc.)
2. 03_refresh_tokens.sql     # Authentication tables
3. 04_login_attempts.sql     # Security tracking
4. 04_material_versions.sql  # Content versioning
5. 05_indexes_materials.sql  # Performance indexes
6. 05_user_progress_upsert.sql # Stored procedures
```

### Step 5: Seed Data Insertion

Test data is inserted into PostgreSQL:

```bash
# Seed data includes:
- 5 test user accounts (2 teachers, 3 students)
- Sample materials and content
- Initial progress records
```

All INSERT statements use `ON CONFLICT DO NOTHING` to ensure idempotency.

### Step 6: MongoDB Initialization

MongoDB collections and indexes are created:

```bash
# Script executed:
- 02_assessment_results.js   # Creates collections and indexes
```

### Step 7: RabbitMQ Topology Setup

RabbitMQ exchange, queues, and bindings are created:

```bash
# Resources created:
- Exchange: edugo.events (topic)
- Queues: material.created, material.updated, etc.
- Bindings: Route events to appropriate queues
```

### Step 8: Connectivity Verification

The script verifies that all services are accessible:

```bash
# Connection tests with retry logic (3 attempts, 2-second delay):
- PostgreSQL connection test
- MongoDB connection test
- RabbitMQ connection test
```

### Step 9: Summary Display

Finally, the script displays:
- Connection strings for all services
- Default credentials
- Test user accounts with passwords
- Useful commands for managing the environment
- Next steps for development

## Verification Steps

After running `make dev-init`, verify your environment is working correctly:

### 1. Check Container Status

```bash
# View running containers
docker ps

# You should see three containers running:
# - edugo-postgres
# - edugo-mongodb
# - edugo-rabbitmq
```

### 2. Test PostgreSQL Connection

```bash
# Connect using psql
psql "postgresql://edugo:edugo123@localhost:5432/edugo"

# Once connected, verify tables exist:
\dt

# You should see tables: users, materials, user_progress, etc.
```

### 3. Test MongoDB Connection

```bash
# Connect using mongosh
mongosh "mongodb://edugo:edugo123@localhost:27017/edugo?authSource=admin"

# Once connected, verify collections:
show collections

# You should see: assessment_results
```

### 4. Test RabbitMQ Management UI

Open your browser and navigate to:
```
http://localhost:15672
```

Login with:
- Username: `edugo`
- Password: `edugo123`

You should see the management dashboard with the `edugo.events` exchange and associated queues.

### 5. Verify Test Users

Try logging in with one of the test accounts:

```bash
# Example: Login as a teacher
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "juan.perez@edugo.com",
    "password": "password123"
  }'

# You should receive a JWT token in the response
```

## Next Steps After Setup

Once your environment is set up, you can:

### 1. Start the API Server

```bash
# Run the API in development mode
make run

# Or with hot-reload (if configured)
make dev
```

The API will be available at `http://localhost:8080`

### 2. View API Documentation

Open Swagger UI in your browser:
```
http://localhost:8080/swagger/index.html
```

### 3. Run Tests

```bash
# Run unit tests
make test

# Run integration tests
make test-integration

# Run all tests with coverage
make test-coverage
```

### 4. View Logs

```bash
# View logs for all containers
docker-compose -f docker-compose-local.yml logs -f

# View logs for a specific container
docker logs -f edugo-postgres
docker logs -f edugo-mongodb
docker logs -f edugo-rabbitmq
```

### 5. Check Environment Status

```bash
# View current status of all services
make dev-status
```

This command shows:
- Container status (running/stopped/missing)
- Database connectivity
- RabbitMQ queue counts

## Idempotency

The initialization script is idempotent, meaning you can run it multiple times safely:

```bash
# Running this multiple times won't cause errors
make dev-init
make dev-init
make dev-init
```

The script will:
- Skip starting containers that are already running
- Skip creating schemas that already exist
- Skip inserting data that already exists
- Display informational messages about existing resources

## Resetting Your Environment

If you need to start fresh (e.g., corrupted data, testing migrations):

```bash
# Reset the entire environment
make dev-reset
```

⚠️ **Warning**: This will:
- Stop all containers
- Remove all volumes (deletes all data)
- Reinitialize the environment from scratch

You'll be prompted for confirmation before any destructive actions.

## Environment Variables

The `.env` file contains configuration for your local environment. Key variables:

```bash
# Database Configuration
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=edugo
POSTGRES_USER=edugo
POSTGRES_PASSWORD=edugo123

MONGODB_HOST=localhost
MONGODB_PORT=27017
MONGODB_DATABASE=edugo
MONGODB_USER=edugo
MONGODB_PASSWORD=edugo123

# RabbitMQ Configuration
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=edugo
RABBITMQ_PASSWORD=edugo123

# Application Configuration
APP_ENV=local
APP_PORT=8080
JWT_SECRET=your-secret-key-here
```

⚠️ **Security Note**: Never commit the `.env` file to version control. It's included in `.gitignore`.

## Troubleshooting

If you encounter issues during setup, see the [Troubleshooting Guide](./TROUBLESHOOTING.md) for common problems and solutions.

For credential reference, see [Credentials Documentation](./CREDENTIALS.md).

## Additional Resources

- [Testing Guide](../TESTING_GUIDE.md)
- [API Documentation](../../README.md)
- [Docker Compose Reference](https://docs.docker.com/compose/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [MongoDB Documentation](https://docs.mongodb.com/)
- [RabbitMQ Documentation](https://www.rabbitmq.com/documentation.html)
