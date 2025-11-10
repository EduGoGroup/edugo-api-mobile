# Implementation Plan

- [x] 1. Update seed data script with password documentation
  - Modify scripts/postgresql/02_seed_data.sql to add plaintext password comments above each user INSERT
  - Ensure all INSERT statements use ON CONFLICT DO NOTHING for idempotency
  - Add a summary section at the end that displays created users and their credentials
  - _Requirements: 4.3, 4.4, 4.5, 4.6_

- [x] 2. Create main initialization script
- [x] 2.1 Create scripts/dev-init.sh with basic structure
  - Create the file with proper shebang and set -euo pipefail
  - Define color constants for output formatting
  - Add main() function that orchestrates the initialization flow
  - Add usage/help function
  - _Requirements: 6.1, 6.2, 8.1, 8.3, 8.4_

- [x] 2.2 Implement validate_prerequisites() function
  - Check if Docker daemon is accessible using docker ps
  - Check if docker-compose command is available
  - Check if .env file exists in project root
  - Display clear error messages for missing prerequisites
  - Return appropriate exit codes
  - _Requirements: 1.1, 1.2, 1.3, 1.4_

- [x] 2.3 Implement manage_containers() function
  - Check if required containers exist and are running
  - Start containers using docker-compose -f docker-compose-local.yml up -d
  - Handle case where containers already exist
  - Display informational messages about container status
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 5.1_

- [x] 2.4 Implement wait_for_health() function
  - Poll PostgreSQL health check with 60 second timeout
  - Poll MongoDB health check with 60 second timeout
  - Poll RabbitMQ health check with 60 second timeout
  - Display progress indicators while waiting
  - Show container logs if health checks fail
  - _Requirements: 2.5, 2.6, 2.7_

- [x] 2.5 Implement setup_postgres() function
  - Check if PostgreSQL tables exist using psql query
  - Execute migration scripts in correct order if tables don't exist
  - Execute 01_create_schema.sql, 03_refresh_tokens.sql, 04_login_attempts.sql, 04_material_versions.sql, 05_indexes_materials.sql, 05_user_progress_upsert.sql
  - Execute 02_seed_data.sql for test data
  - Handle errors and display troubleshooting information
  - Skip if schema already exists
  - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6, 3.7, 3.8, 3.9, 4.1, 4.2, 5.2_

- [x] 2.6 Implement setup_mongodb() function
  - Check if MongoDB collections exist using mongosh
  - Execute scripts/mongodb/02_assessment_results.js if collections don't exist
  - Handle errors and display troubleshooting information
  - Skip if collections already exist
  - _Requirements: 3.10, 3.11, 5.2_

- [x] 2.7 Implement setup_rabbitmq() function
  - Create exchange "edugo.events" using rabbitmqadmin or amqp client
  - Create required queues (material.created, material.updated, etc.)
  - Create bindings between exchange and queues
  - Handle errors gracefully
  - Skip if topology already exists
  - _Requirements: 5.2_

- [x] 2.8 Implement verify_connectivity() function
  - Test PostgreSQL connection with retry mechanism (3 attempts, 2 second delay)
  - Test MongoDB connection with retry mechanism
  - Test RabbitMQ connection with retry mechanism
  - Display connection status for each service
  - Show troubleshooting suggestions if connections fail
  - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5, 7.6, 7.7_

- [x] 2.9 Implement display_summary() function
  - Display connection strings for PostgreSQL, MongoDB, and RabbitMQ
  - Display default credentials for all services
  - Display test user accounts with plaintext passwords
  - Display useful commands (logs, stop, restart)
  - Display next steps for the developer
  - _Requirements: 6.3, 6.4, 6.5_

- [x] 3. Create reset script
- [x] 3.1 Create scripts/dev-reset.sh
  - Create file with proper shebang and error handling
  - Implement confirm_reset() function that prompts user for confirmation
  - Implement cleanup_environment() function that stops containers and removes volumes
  - Implement reinitialize() function that calls dev-init.sh
  - Add main flow that confirms, cleans, and reinitializes
  - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.5, 10.6_

- [x] 4. Create status script
- [x] 4.1 Create scripts/dev-status.sh
  - Implement check_container_status() to show running/stopped/missing containers
  - Implement check_database_status() to test PostgreSQL and MongoDB connectivity
  - Implement check_rabbitmq_status() to test RabbitMQ and show queue counts
  - Display formatted output with status indicators
  - _Requirements: 6.5_

- [x] 5. Update Makefile
- [x] 5.1 Add dev-init target
  - Add target that calls ./scripts/dev-init.sh
  - Add help text describing the command
  - _Requirements: 6.1, 6.2_

- [x] 5.2 Update dev-setup target
  - Modify existing dev-setup to call new dev-init script
  - Ensure backward compatibility
  - _Requirements: 8.5_

- [x] 5.3 Add dev-reset target
  - Add target that calls ./scripts/dev-reset.sh
  - Add help text with warning about data loss
  - _Requirements: 10.1_

- [x] 5.4 Add dev-status target
  - Add target that calls ./scripts/dev-status.sh
  - Add help text describing the command
  - _Requirements: 6.5_

- [x] 6. Create documentation
- [x] 6.1 Create docs/development/ENVIRONMENT_SETUP.md
  - Document prerequisites (Docker, Docker Compose)
  - Provide quick start guide with make dev-init
  - Document detailed setup steps
  - Document verification steps
  - Document next steps after setup
  - _Requirements: 9.1, 9.2, 9.3, 9.4_

- [x] 6.2 Create docs/development/TROUBLESHOOTING.md
  - Document common issues (Docker not running, port conflicts, timeouts)
  - Provide solutions for each issue
  - Document reset procedures
  - Include links to relevant logs and commands
  - _Requirements: 9.5_

- [x] 6.3 Create docs/development/CREDENTIALS.md
  - Document PostgreSQL default credentials
  - Document MongoDB default credentials
  - Document RabbitMQ default credentials
  - Document test user accounts with passwords
  - Add security notes about development-only usage
  - _Requirements: 9.6_

- [x] 7. Make scripts executable
  - Run chmod +x scripts/dev-init.sh
  - Run chmod +x scripts/dev-reset.sh
  - Run chmod +x scripts/dev-status.sh
  - _Requirements: 8.4_

- [x] 8. Verify implementation
- [x] 8.1 Test fresh installation
  - Clone repository to clean directory
  - Run make dev-init
  - Verify all services start successfully
  - Verify seed data is inserted
  - Verify credentials work for all services
  - _Requirements: 6.2, 6.3, 6.4, 6.5_

- [x] 8.2 Test idempotency
  - Run make dev-init twice consecutively
  - Verify no errors on second run
  - Verify no duplicate data in databases
  - Verify informational messages about existing resources
  - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, 5.6_

- [x] 8.3 Test reset functionality
  - Run make dev-reset
  - Verify confirmation prompt appears
  - Confirm reset
  - Verify containers are stopped and removed
  - Verify volumes are removed
  - Verify environment is recreated successfully
  - _Requirements: 10.1, 10.2, 10.3, 10.4, 10.5, 10.6_

- [x] 8.4 Test error handling
  - Stop Docker daemon and run make dev-init
  - Verify clear error message about Docker not running
  - Start Docker and verify recovery
  - Test with missing .env file
  - Verify appropriate error messages
  - _Requirements: 1.2, 1.3, 6.6_

- [x] 8.5 Test connectivity verification
  - Use displayed connection strings to connect to each service
  - Connect to PostgreSQL using psql
  - Connect to MongoDB using mongosh
  - Access RabbitMQ management UI in browser
  - Verify all connections work with provided credentials
  - _Requirements: 7.1, 7.2, 7.3, 7.7_
