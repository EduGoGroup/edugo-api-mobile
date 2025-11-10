# Requirements Document

## Introduction

This document defines the requirements for a robust development environment initialization system for the EduGo API Mobile project. The system will provide developers with a reliable, idempotent process to set up their local development environment, including Docker containers, database schemas, and seed data, without interfering with the API runtime logic.

## Glossary

- **DevInit System**: The development environment initialization system that validates and prepares the local environment
- **Docker Daemon**: The Docker service running on the host machine
- **Container Health**: The operational status of a Docker container (running, healthy, unhealthy)
- **Schema Migration**: The process of creating database tables and structures
- **Seed Data**: Pre-populated test data for development purposes
- **Idempotent Operation**: An operation that can be executed multiple times with the same result
- **Make Command**: A command defined in the Makefile for task automation
- **Connection String**: A URL or configuration string used to connect to a database or service

## Requirements

### Requirement 1

**User Story:** As a developer, I want to validate that Docker is running before attempting any setup, so that I receive clear error messages if prerequisites are missing

#### Acceptance Criteria

1. WHEN the DevInit System executes, THE DevInit System SHALL verify that the Docker Daemon is accessible
2. IF the Docker Daemon is not accessible, THEN THE DevInit System SHALL display an error message indicating Docker is not running
3. IF the Docker Daemon is not accessible, THEN THE DevInit System SHALL terminate with a non-zero exit code
4. WHEN the Docker Daemon is accessible, THE DevInit System SHALL proceed to container validation

### Requirement 2

**User Story:** As a developer, I want the system to automatically start required containers if they don't exist, so that I don't have to manually manage Docker Compose

#### Acceptance Criteria

1. WHEN the DevInit System validates containers, THE DevInit System SHALL check if PostgreSQL container exists and is running
2. WHEN the DevInit System validates containers, THE DevInit System SHALL check if MongoDB container exists and is running
3. WHEN the DevInit System validates containers, THE DevInit System SHALL check if RabbitMQ container exists and is running
4. IF any required container is not running, THEN THE DevInit System SHALL start the container using docker-compose-local.yml
5. WHEN containers are started, THE DevInit System SHALL wait for Container Health checks to pass before proceeding
6. THE DevInit System SHALL implement a timeout of 60 seconds for Container Health checks
7. IF Container Health checks fail after timeout, THEN THE DevInit System SHALL display an error message with container logs

### Requirement 3

**User Story:** As a developer, I want database schemas to be created automatically if they don't exist, so that I can start development immediately

#### Acceptance Criteria

1. WHEN the DevInit System validates PostgreSQL, THE DevInit System SHALL check if required tables exist
2. IF PostgreSQL tables do not exist, THEN THE DevInit System SHALL execute Schema Migration scripts in order
3. THE DevInit System SHALL execute scripts/postgresql/01_create_schema.sql
4. THE DevInit System SHALL execute scripts/postgresql/03_refresh_tokens.sql
5. THE DevInit System SHALL execute scripts/postgresql/04_login_attempts.sql
6. THE DevInit System SHALL execute scripts/postgresql/04_material_versions.sql
7. THE DevInit System SHALL execute scripts/postgresql/05_indexes_materials.sql
8. THE DevInit System SHALL execute scripts/postgresql/05_user_progress_upsert.sql
9. IF any Schema Migration fails, THEN THE DevInit System SHALL display the error and terminate
10. WHEN the DevInit System validates MongoDB, THE DevInit System SHALL check if required collections exist
11. IF MongoDB collections do not exist, THEN THE DevInit System SHALL execute scripts/mongodb/02_assessment_results.js

### Requirement 4

**User Story:** As a developer, I want seed data to be inserted automatically with clear password documentation, so that I can test authentication without guessing credentials

#### Acceptance Criteria

1. WHEN the DevInit System inserts Seed Data, THE DevInit System SHALL check if users already exist
2. IF users do not exist, THEN THE DevInit System SHALL execute scripts/postgresql/02_seed_data.sql
3. THE DevInit System SHALL update 02_seed_data.sql to include plaintext password comments for each user
4. WHEN Seed Data is inserted, THE DevInit System SHALL display a summary of created users with their credentials
5. THE DevInit System SHALL display credentials in the format "email: [email] | password: [plaintext]"
6. IF Seed Data already exists, THEN THE DevInit System SHALL skip insertion and display a message indicating data exists

### Requirement 5

**User Story:** As a developer, I want the initialization process to be idempotent, so that I can run it multiple times without breaking my environment

#### Acceptance Criteria

1. WHEN the DevInit System executes multiple times, THE DevInit System SHALL not fail if containers already exist
2. WHEN the DevInit System executes multiple times, THE DevInit System SHALL not fail if schemas already exist
3. WHEN the DevInit System executes multiple times, THE DevInit System SHALL not fail if Seed Data already exists
4. WHEN the DevInit System detects existing resources, THE DevInit System SHALL display informational messages
5. THE DevInit System SHALL use SQL "ON CONFLICT DO NOTHING" clauses for Seed Data insertion
6. THE DevInit System SHALL use "IF NOT EXISTS" clauses for schema creation where applicable

### Requirement 6

**User Story:** As a developer, I want a single Make Command to initialize my environment, so that I have a simple entry point for setup

#### Acceptance Criteria

1. THE DevInit System SHALL provide a Make Command named "dev-init"
2. WHEN a developer executes "make dev-init", THE DevInit System SHALL perform all validation and setup steps
3. THE DevInit System SHALL display progress messages for each step
4. WHEN initialization completes successfully, THE DevInit System SHALL display a success message with next steps
5. THE DevInit System SHALL display Connection String information for each service
6. IF initialization fails, THEN THE DevInit System SHALL display a clear error message indicating which step failed

### Requirement 7

**User Story:** As a developer, I want the system to verify service connectivity before declaring success, so that I know the environment is truly ready

#### Acceptance Criteria

1. WHEN the DevInit System validates PostgreSQL, THE DevInit System SHALL attempt to connect using the Connection String
2. WHEN the DevInit System validates MongoDB, THE DevInit System SHALL attempt to connect using the Connection String
3. WHEN the DevInit System validates RabbitMQ, THE DevInit System SHALL attempt to connect using the Connection String
4. THE DevInit System SHALL implement a retry mechanism with 3 attempts for each connection
5. THE DevInit System SHALL wait 2 seconds between retry attempts
6. IF any service connection fails after retries, THEN THE DevInit System SHALL display an error with troubleshooting suggestions
7. WHEN all services are connected successfully, THE DevInit System SHALL display a confirmation message

### Requirement 8

**User Story:** As a developer, I want the initialization script to be separate from the API runtime, so that setup logic doesn't interfere with application code

#### Acceptance Criteria

1. THE DevInit System SHALL be implemented as a standalone shell script
2. THE DevInit System SHALL not modify any application source code in internal/ or cmd/ directories
3. THE DevInit System SHALL be located in the scripts/ directory
4. THE DevInit System SHALL be executable independently of the API application
5. THE DevInit System SHALL not be invoked automatically when running "make run" or "make dev"

### Requirement 9

**User Story:** As a developer, I want clear documentation of what the initialization process does, so that I understand what's happening in my environment

#### Acceptance Criteria

1. THE DevInit System SHALL include a README.md file documenting the initialization process
2. THE README SHALL document all prerequisites (Docker, Docker Compose)
3. THE README SHALL document all Make Command options available
4. THE README SHALL document the order of operations performed
5. THE README SHALL document how to troubleshoot common issues
6. THE README SHALL document the default credentials for all services
7. THE README SHALL document how to reset the environment completely

### Requirement 10

**User Story:** As a developer, I want the ability to reset my environment to a clean state, so that I can start fresh if something goes wrong

#### Acceptance Criteria

1. THE DevInit System SHALL provide a Make Command named "dev-reset"
2. WHEN a developer executes "make dev-reset", THE DevInit System SHALL stop all containers
3. WHEN a developer executes "make dev-reset", THE DevInit System SHALL remove all volumes
4. WHEN a developer executes "make dev-reset", THE DevInit System SHALL execute "make dev-init" to recreate the environment
5. THE DevInit System SHALL prompt for confirmation before destroying volumes
6. IF the developer cancels the reset, THEN THE DevInit System SHALL terminate without changes
