# Design Document

## Overview

The Development Environment Initialization System (DevInit) provides a robust, idempotent solution for setting up the EduGo API Mobile local development environment. The system validates prerequisites, manages Docker containers, creates database schemas, and populates seed data with clear credential documentation.

### Key Design Principles

1. **Idempotency**: All operations can be executed multiple times safely
2. **Separation of Concerns**: Setup logic is completely separate from application runtime
3. **Clear Feedback**: Detailed progress messages and error reporting
4. **Fail-Fast**: Validate prerequisites before attempting setup
5. **Self-Documenting**: Credentials and connection strings displayed after setup

## Architecture

### Component Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                        Developer                             │
└────────────────────────┬────────────────────────────────────┘
                         │
                         │ make dev-init
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                      Makefile                                │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ dev-init: Invoke scripts/dev-init.sh                 │   │
│  │ dev-reset: Invoke scripts/dev-reset.sh               │   │
│  │ dev-status: Invoke scripts/dev-status.sh             │   │
│  └──────────────────────────────────────────────────────┘   │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                  scripts/dev-init.sh                         │
│  ┌──────────────────────────────────────────────────────┐   │
│  │ 1. validate_prerequisites()                          │   │
│  │ 2. manage_containers()                               │   │
│  │ 3. wait_for_health()                                 │   │
│  │ 4. setup_postgres()                                  │   │
│  │ 5. setup_mongodb()                                   │   │
│  │ 6. setup_rabbitmq()                                  │   │
│  │ 7. verify_connectivity()                             │   │
│  │ 8. display_summary()                                 │   │
│  └──────────────────────────────────────────────────────┘   │
└────────────────────────┬────────────────────────────────────┘
                         │
         ┌───────────────┼───────────────┐
         │               │               │
         ▼               ▼               ▼
┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│  PostgreSQL │  │   MongoDB   │  │  RabbitMQ   │
│  Container  │  │  Container  │  │  Container  │
└─────────────┘  └─────────────┘  └─────────────┘
```

### Execution Flow

```
START
  │
  ├─► validate_prerequisites()
  │     ├─► Check Docker daemon
  │     ├─► Check docker-compose
  │     └─► Check .env file
  │
  ├─► manage_containers()
  │     ├─► Check if containers exist
  │     ├─► Start containers if needed
  │     └─► Use docker-compose-local.yml
  │
  ├─► wait_for_health()
  │     ├─► PostgreSQL health check (60s timeout)
  │     ├─► MongoDB health check (60s timeout)
  │     └─► RabbitMQ health check (60s timeout)
  │
  ├─► setup_postgres()
  │     ├─► Check if schema exists
  │     ├─► Execute migration scripts if needed
  │     │     ├─► 01_create_schema.sql
  │     │     ├─► 03_refresh_tokens.sql
  │     │     ├─► 04_login_attempts.sql
  │     │     ├─► 04_material_versions.sql
  │     │     ├─► 05_indexes_materials.sql
  │     │     └─► 05_user_progress_upsert.sql
  │     └─► Insert seed data (02_seed_data.sql)
  │
  ├─► setup_mongodb()
  │     ├─► Check if collections exist
  │     └─► Execute 02_assessment_results.js if needed
  │
  ├─► setup_rabbitmq()
  │     ├─► Create exchange (edugo.events)
  │     ├─► Create queues
  │     └─► Create bindings
  │
  ├─► verify_connectivity()
  │     ├─► Test PostgreSQL connection (3 retries)
  │     ├─► Test MongoDB connection (3 retries)
  │     └─► Test RabbitMQ connection (3 retries)
  │
  └─► display_summary()
        ├─► Show connection strings
        ├─► Show default credentials
        └─► Show next steps
END
```

## Components and Interfaces

### 1. Main Initialization Script (`scripts/dev-init.sh`)

**Purpose**: Orchestrates the entire initialization process

**Functions**:

```bash
# Validates that all prerequisites are met
validate_prerequisites() {
    # Check Docker daemon is running
    # Check docker-compose is installed
    # Check .env file exists
    # Return: 0 on success, 1 on failure
}

# Manages Docker containers lifecycle
manage_containers() {
    # Check if containers exist and are running
    # Start containers using docker-compose-local.yml
    # Return: 0 on success, 1 on failure
}

# Waits for container health checks to pass
wait_for_health() {
    # Poll PostgreSQL health endpoint (60s timeout)
    # Poll MongoDB health endpoint (60s timeout)
    # Poll RabbitMQ health endpoint (60s timeout)
    # Return: 0 on success, 1 on timeout
}

# Sets up PostgreSQL schema and data
setup_postgres() {
    # Check if tables exist
    # Execute migration scripts in order
    # Insert seed data with ON CONFLICT DO NOTHING
    # Return: 0 on success, 1 on failure
}

# Sets up MongoDB collections and indexes
setup_mongodb() {
    # Check if collections exist
    # Execute JavaScript initialization script
    # Return: 0 on success, 1 on failure
}

# Sets up RabbitMQ topology
setup_rabbitmq() {
    # Create exchange
    # Create queues
    # Create bindings
    # Return: 0 on success, 1 on failure
}

# Verifies connectivity to all services
verify_connectivity() {
    # Test PostgreSQL connection with retry
    # Test MongoDB connection with retry
    # Test RabbitMQ connection with retry
    # Return: 0 on success, 1 on failure
}

# Displays summary and next steps
display_summary() {
    # Show connection strings
    # Show default credentials
    # Show useful commands
}
```

### 2. Reset Script (`scripts/dev-reset.sh`)

**Purpose**: Completely resets the development environment

**Functions**:

```bash
# Prompts user for confirmation
confirm_reset() {
    # Display warning about data loss
    # Wait for user confirmation (y/n)
    # Return: 0 to proceed, 1 to cancel
}

# Stops and removes all containers and volumes
cleanup_environment() {
    # Stop containers
    # Remove containers
    # Remove volumes
    # Return: 0 on success
}

# Reinitializes the environment
reinitialize() {
    # Call dev-init.sh
    # Return: exit code from dev-init.sh
}
```

### 3. Status Script (`scripts/dev-status.sh`)

**Purpose**: Displays current status of development environment

**Functions**:

```bash
# Checks status of all containers
check_container_status() {
    # Query Docker for container status
    # Display running/stopped/missing status
}

# Checks database connectivity
check_database_status() {
    # Test PostgreSQL connection
    # Test MongoDB connection
    # Display connection status
}

# Checks RabbitMQ status
check_rabbitmq_status() {
    # Test RabbitMQ connection
    # Display queue counts
}
```

### 4. Updated Seed Data Script (`scripts/postgresql/02_seed_data.sql`)

**Changes**:
- Add plaintext password comments for each user
- Ensure all INSERT statements use `ON CONFLICT DO NOTHING`
- Add summary output showing created users and credentials

**Format**:
```sql
-- User: juan.perez@edugo.com
-- Password: password123 (for development only)
INSERT INTO users (id, email, password_hash, ...)
VALUES (
    '11111111-1111-1111-1111-111111111111',
    'juan.perez@edugo.com',
    '$2a$10$...',  -- bcrypt hash of "password123"
    ...
)
ON CONFLICT (id) DO NOTHING;
```

### 5. Makefile Targets

```makefile
dev-init: ## Initialize development environment
	@./scripts/dev-init.sh

dev-reset: ## Reset development environment (destructive)
	@./scripts/dev-reset.sh

dev-status: ## Show development environment status
	@./scripts/dev-status.sh
```

## Data Models

### Environment Configuration

The system reads configuration from multiple sources:

1. **`.env` file**: Contains sensitive credentials
2. **`docker-compose-local.yml`**: Defines container configuration
3. **SQL migration scripts**: Define database schema
4. **Seed data scripts**: Populate initial data

### Default Credentials

```yaml
PostgreSQL:
  host: localhost
  port: 5432
  database: edugo
  username: edugo
  password: edugo123

MongoDB:
  host: localhost
  port: 27017
  database: edugo
  username: edugo
  password: edugo123
  auth_source: admin

RabbitMQ:
  host: localhost
  port: 5672
  management_port: 15672
  username: edugo
  password: edugo123

Test Users:
  - email: juan.perez@edugo.com
    password: password123
    role: teacher

  - email: maria.gonzalez@edugo.com
    password: password123
    role: teacher

  - email: carlos.rodriguez@student.edugo.com
    password: password123
    role: student

  - email: ana.martinez@student.edugo.com
    password: password123
    role: student

  - email: luis.fernandez@student.edugo.com
    password: password123
    role: student
```

## Error Handling

### Error Categories

1. **Prerequisites Errors**
   - Docker not running
   - docker-compose not installed
   - .env file missing

2. **Container Errors**
   - Failed to start container
   - Health check timeout
   - Port already in use

3. **Database Errors**
   - Connection failed
   - Schema creation failed
   - Seed data insertion failed

4. **Network Errors**
   - Cannot reach container
   - DNS resolution failed
   - Connection timeout

### Error Handling Strategy

```bash
# Example error handling pattern
function setup_postgres() {
    local max_retries=3
    local retry_count=0

    while [ $retry_count -lt $max_retries ]; do
        if execute_postgres_setup; then
            return 0
        fi

        retry_count=$((retry_count + 1))
        if [ $retry_count -lt $max_retries ]; then
            echo "Retry $retry_count/$max_retries..."
            sleep 2
        fi
    done

    echo "ERROR: PostgreSQL setup failed after $max_retries attempts"
    echo "Troubleshooting:"
    echo "  1. Check container logs: docker logs edugo-postgres"
    echo "  2. Verify connection: psql -h localhost -U edugo -d edugo"
    echo "  3. Reset environment: make dev-reset"
    return 1
}
```

### Error Messages

All error messages follow this format:

```
❌ Error: [Component] [Action] failed
   Reason: [Specific error message]

   Troubleshooting:
   1. [First suggestion]
   2. [Second suggestion]
   3. [Third suggestion]

   For more help: [Documentation link or command]
```

## Testing Strategy

### Manual Testing Checklist

1. **Fresh Installation**
   - Clone repository
   - Run `make dev-init`
   - Verify all services start
   - Verify seed data is inserted
   - Verify credentials work

2. **Idempotency Testing**
   - Run `make dev-init` twice
   - Verify no errors on second run
   - Verify no duplicate data

3. **Reset Testing**
   - Run `make dev-reset`
   - Verify containers are removed
   - Verify volumes are removed
   - Verify environment is recreated

4. **Error Handling Testing**
   - Stop Docker and run `make dev-init`
   - Verify clear error message
   - Kill a container mid-setup
   - Verify recovery or clear error

5. **Connectivity Testing**
   - Use displayed connection strings
   - Connect to PostgreSQL with psql
   - Connect to MongoDB with mongosh
   - Access RabbitMQ management UI

### Automated Testing

While the initialization scripts are primarily for development, we can add basic smoke tests:

```bash
# scripts/test-dev-init.sh
#!/bin/bash

# Test 1: Prerequisites check
test_prerequisites() {
    ./scripts/dev-init.sh --check-only
}

# Test 2: Container startup
test_container_startup() {
    make dev-reset -y
    make dev-init
    docker ps | grep edugo-postgres
    docker ps | grep edugo-mongodb
    docker ps | grep edugo-rabbitmq
}

# Test 3: Database connectivity
test_database_connectivity() {
    psql "postgresql://edugo:edugo123@localhost:5432/edugo" -c "SELECT 1"
    mongosh "mongodb://edugo:edugo123@localhost:27017/edugo" --eval "db.adminCommand('ping')"
}

# Test 4: Seed data verification
test_seed_data() {
    psql "postgresql://edugo:edugo123@localhost:5432/edugo" \
        -c "SELECT COUNT(*) FROM users" | grep "5"
}
```

## Implementation Notes

### Shell Script Best Practices

1. **Use `set -euo pipefail`**: Fail fast on errors
2. **Color-coded output**: Use ANSI colors for better readability
3. **Progress indicators**: Show what's happening at each step
4. **Timeouts**: Don't wait forever for services
5. **Cleanup on failure**: Provide clear instructions for recovery

### Docker Compose Strategy

- Use `docker-compose-local.yml` for development
- Don't modify `docker-compose.yml` (may be used for other purposes)
- Use health checks in compose file
- Use named volumes for data persistence

### SQL Script Execution Order

The order matters for foreign key constraints:

1. `01_create_schema.sql` - Base tables
2. `03_refresh_tokens.sql` - Auth tables
3. `04_login_attempts.sql` - Security tables
4. `04_material_versions.sql` - Versioning tables
5. `05_indexes_materials.sql` - Performance indexes
6. `05_user_progress_upsert.sql` - Stored procedures
7. `02_seed_data.sql` - Test data (last)

### Retry Logic

All connectivity checks use exponential backoff:

```bash
retry_with_backoff() {
    local max_attempts=$1
    shift
    local command="$@"
    local attempt=1

    while [ $attempt -le $max_attempts ]; do
        if eval "$command"; then
            return 0
        fi

        if [ $attempt -lt $max_attempts ]; then
            local wait_time=$((2 ** attempt))
            echo "Attempt $attempt failed, waiting ${wait_time}s..."
            sleep $wait_time
        fi

        attempt=$((attempt + 1))
    done

    return 1
}
```

## Security Considerations

1. **Development-Only Credentials**: All default credentials are clearly marked as development-only
2. **No Production Secrets**: The scripts never touch production credentials
3. **Local-Only Access**: Containers bind to localhost by default
4. **Clear Documentation**: Seed data comments clearly show plaintext passwords
5. **Gitignore Protection**: `.env` file is gitignored to prevent credential leaks

## Performance Considerations

1. **Parallel Health Checks**: Check all services concurrently where possible
2. **Skip Unnecessary Work**: Check if setup is already done before executing
3. **Minimal Container Images**: Use Alpine-based images for faster startup
4. **Connection Pooling**: Reuse database connections during setup
5. **Efficient SQL**: Use bulk inserts and transactions for seed data

## Documentation Structure

```
docs/
└── development/
    ├── ENVIRONMENT_SETUP.md       # Main setup guide
    ├── TROUBLESHOOTING.md         # Common issues and solutions
    └── CREDENTIALS.md             # Default credentials reference
```

### ENVIRONMENT_SETUP.md Contents

- Prerequisites
- Quick start guide
- Detailed setup steps
- Verification steps
- Next steps after setup

### TROUBLESHOOTING.md Contents

- Docker not running
- Port conflicts
- Connection timeouts
- Schema creation failures
- Seed data issues
- Reset procedures

### CREDENTIALS.md Contents

- PostgreSQL credentials
- MongoDB credentials
- RabbitMQ credentials
- Test user accounts
- Security notes
