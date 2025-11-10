#!/bin/bash

# Development Environment Initialization Script
# This script sets up the local development environment for EduGo API Mobile

set -euo pipefail

# Color constants for output formatting
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly CYAN='\033[0;36m'
readonly NC='\033[0m' # No Color

# Script directory
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

# Function to display usage information
usage() {
    cat << EOF
${CYAN}EduGo Development Environment Initialization${NC}

${YELLOW}Usage:${NC}
    $0 [OPTIONS]

${YELLOW}Options:${NC}
    -h, --help          Show this help message
    --check-only        Only validate prerequisites without setup

${YELLOW}Description:${NC}
    This script initializes the local development environment by:
    1. Validating prerequisites (Docker, docker-compose, .env)
    2. Starting required containers (PostgreSQL, MongoDB, RabbitMQ)
    3. Creating database schemas
    4. Inserting seed data
    5. Verifying connectivity to all services

${YELLOW}Examples:${NC}
    $0                  # Run full initialization
    $0 --check-only     # Only check prerequisites

${YELLOW}For more information:${NC}
    See docs/development/ENVIRONMENT_SETUP.md

EOF
}

# Function to print colored messages
print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_step() {
    echo -e "\n${CYAN}==>${NC} ${YELLOW}$1${NC}\n"
}

# Main orchestration function
main() {
    local check_only=false
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                usage
                exit 0
                ;;
            --check-only)
                check_only=true
                shift
                ;;
            *)
                print_error "Unknown option: $1"
                usage
                exit 1
                ;;
        esac
    done
    
    echo -e "${CYAN}╔════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║${NC}  ${YELLOW}EduGo Development Environment Initialization${NC}          ${CYAN}║${NC}"
    echo -e "${CYAN}╚════════════════════════════════════════════════════════════╝${NC}\n"
    
    # Step 1: Validate prerequisites
    print_step "Step 1: Validating prerequisites"
    if ! validate_prerequisites; then
        print_error "Prerequisites validation failed"
        exit 1
    fi
    print_success "All prerequisites validated"
    
    # Exit if only checking prerequisites
    if [ "$check_only" = true ]; then
        print_success "Prerequisites check completed successfully"
        exit 0
    fi
    
    # Step 2: Manage containers
    print_step "Step 2: Managing Docker containers"
    if ! manage_containers; then
        print_error "Container management failed"
        exit 1
    fi
    print_success "Containers are running"
    
    # Step 3: Wait for health checks
    print_step "Step 3: Waiting for services to be healthy"
    if ! wait_for_health; then
        print_error "Health checks failed"
        exit 1
    fi
    print_success "All services are healthy"
    
    # Step 4: Setup PostgreSQL
    print_step "Step 4: Setting up PostgreSQL"
    if ! setup_postgres; then
        print_error "PostgreSQL setup failed"
        exit 1
    fi
    print_success "PostgreSQL setup completed"
    
    # Step 5: Setup MongoDB
    print_step "Step 5: Setting up MongoDB"
    if ! setup_mongodb; then
        print_error "MongoDB setup failed"
        exit 1
    fi
    print_success "MongoDB setup completed"
    
    # Step 6: Setup RabbitMQ
    print_step "Step 6: Setting up RabbitMQ"
    if ! setup_rabbitmq; then
        print_error "RabbitMQ setup failed"
        exit 1
    fi
    print_success "RabbitMQ setup completed"
    
    # Step 7: Verify connectivity
    print_step "Step 7: Verifying service connectivity"
    if ! verify_connectivity; then
        print_error "Connectivity verification failed"
        exit 1
    fi
    print_success "All services are accessible"
    
    # Step 8: Display summary
    print_step "Step 8: Environment ready!"
    display_summary
    
    echo -e "\n${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║${NC}  ${GREEN}✓ Development environment initialized successfully!${NC}    ${GREEN}║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}\n"
}

# Validates that all prerequisites are met
validate_prerequisites() {
    local all_valid=true
    
    # Check if Docker daemon is accessible
    print_info "Checking Docker daemon..."
    if ! docker ps >/dev/null 2>&1; then
        print_error "Docker daemon is not accessible"
        echo -e "  ${RED}Reason:${NC} Docker is not running or you don't have permission to access it"
        echo -e "\n  ${YELLOW}Troubleshooting:${NC}"
        echo -e "    1. Start Docker Desktop or Docker daemon"
        echo -e "    2. Verify Docker is running: ${CYAN}docker ps${NC}"
        echo -e "    3. Check Docker permissions for your user"
        all_valid=false
    else
        print_success "Docker daemon is accessible"
    fi
    
    # Check if docker-compose command is available
    print_info "Checking docker-compose..."
    if ! command -v docker-compose >/dev/null 2>&1; then
        print_error "docker-compose command not found"
        echo -e "  ${RED}Reason:${NC} docker-compose is not installed or not in PATH"
        echo -e "\n  ${YELLOW}Troubleshooting:${NC}"
        echo -e "    1. Install docker-compose: ${CYAN}https://docs.docker.com/compose/install/${NC}"
        echo -e "    2. Verify installation: ${CYAN}docker-compose --version${NC}"
        all_valid=false
    else
        local compose_version=$(docker-compose --version 2>&1)
        print_success "docker-compose is available (${compose_version})"
    fi
    
    # Check if .env file exists in project root
    print_info "Checking .env file..."
    if [ ! -f "${PROJECT_ROOT}/.env" ]; then
        print_error ".env file not found in project root"
        echo -e "  ${RED}Reason:${NC} Configuration file .env is missing"
        echo -e "\n  ${YELLOW}Troubleshooting:${NC}"
        echo -e "    1. Copy the example file: ${CYAN}cp .env.example .env${NC}"
        echo -e "    2. Edit .env with your local configuration"
        echo -e "    3. Ensure .env is in the project root: ${CYAN}${PROJECT_ROOT}${NC}"
        all_valid=false
    else
        print_success ".env file exists"
    fi
    
    if [ "$all_valid" = false ]; then
        echo -e "\n${RED}Prerequisites validation failed. Please fix the issues above and try again.${NC}\n"
        return 1
    fi
    
    return 0
}

manage_containers() {
    local compose_file="${PROJECT_ROOT}/docker-compose-local.yml"
    
    # Check if docker-compose-local.yml exists
    if [ ! -f "$compose_file" ]; then
        print_error "docker-compose-local.yml not found"
        echo -e "  ${RED}Reason:${NC} Required compose file is missing"
        echo -e "\n  ${YELLOW}Expected location:${NC} ${compose_file}"
        return 1
    fi
    
    # Define required containers
    local required_containers=("edugo-postgres" "edugo-mongodb" "edugo-rabbitmq")
    local all_running=true
    
    # Check status of each required container
    print_info "Checking container status..."
    for container in "${required_containers[@]}"; do
        if docker ps --format '{{.Names}}' | grep -q "^${container}$"; then
            print_success "${container} is running"
        elif docker ps -a --format '{{.Names}}' | grep -q "^${container}$"; then
            print_warning "${container} exists but is not running"
            all_running=false
        else
            print_info "${container} does not exist"
            all_running=false
        fi
    done
    
    # Start containers if needed
    if [ "$all_running" = false ]; then
        print_info "Starting containers with docker-compose..."
        cd "${PROJECT_ROOT}"
        if docker-compose -f docker-compose-local.yml up -d; then
            print_success "Containers started successfully"
        else
            print_error "Failed to start containers"
            echo -e "\n  ${YELLOW}Troubleshooting:${NC}"
            echo -e "    1. Check Docker logs: ${CYAN}docker-compose -f docker-compose-local.yml logs${NC}"
            echo -e "    2. Check for port conflicts: ${CYAN}docker ps -a${NC}"
            echo -e "    3. Try stopping existing containers: ${CYAN}docker-compose -f docker-compose-local.yml down${NC}"
            return 1
        fi
    else
        print_success "All required containers are already running"
    fi
    
    return 0
}

wait_for_health() {
    local timeout=60
    local check_interval=2
    
    # Function to check PostgreSQL health
    check_postgres_health() {
        docker exec edugo-postgres pg_isready -U edugo >/dev/null 2>&1
    }
    
    # Function to check MongoDB health
    check_mongodb_health() {
        docker exec edugo-mongodb mongosh --quiet --eval "db.adminCommand('ping')" >/dev/null 2>&1
    }
    
    # Function to check RabbitMQ health
    check_rabbitmq_health() {
        docker exec edugo-rabbitmq rabbitmq-diagnostics -q ping >/dev/null 2>&1
    }
    
    # Wait for PostgreSQL
    print_info "Waiting for PostgreSQL to be healthy..."
    local elapsed=0
    while [ $elapsed -lt $timeout ]; do
        if check_postgres_health; then
            print_success "PostgreSQL is healthy"
            break
        fi
        echo -n "."
        sleep $check_interval
        elapsed=$((elapsed + check_interval))
    done
    
    if [ $elapsed -ge $timeout ]; then
        print_error "PostgreSQL health check timed out after ${timeout}s"
        echo -e "\n  ${YELLOW}Container logs:${NC}"
        docker logs --tail 20 edugo-postgres 2>&1 | sed 's/^/    /'
        echo -e "\n  ${YELLOW}Troubleshooting:${NC}"
        echo -e "    1. Check full logs: ${CYAN}docker logs edugo-postgres${NC}"
        echo -e "    2. Restart container: ${CYAN}docker restart edugo-postgres${NC}"
        return 1
    fi
    
    # Wait for MongoDB
    print_info "Waiting for MongoDB to be healthy..."
    elapsed=0
    while [ $elapsed -lt $timeout ]; do
        if check_mongodb_health; then
            print_success "MongoDB is healthy"
            break
        fi
        echo -n "."
        sleep $check_interval
        elapsed=$((elapsed + check_interval))
    done
    
    if [ $elapsed -ge $timeout ]; then
        print_error "MongoDB health check timed out after ${timeout}s"
        echo -e "\n  ${YELLOW}Container logs:${NC}"
        docker logs --tail 20 edugo-mongodb 2>&1 | sed 's/^/    /'
        echo -e "\n  ${YELLOW}Troubleshooting:${NC}"
        echo -e "    1. Check full logs: ${CYAN}docker logs edugo-mongodb${NC}"
        echo -e "    2. Restart container: ${CYAN}docker restart edugo-mongodb${NC}"
        return 1
    fi
    
    # Wait for RabbitMQ
    print_info "Waiting for RabbitMQ to be healthy..."
    elapsed=0
    while [ $elapsed -lt $timeout ]; do
        if check_rabbitmq_health; then
            print_success "RabbitMQ is healthy"
            break
        fi
        echo -n "."
        sleep $check_interval
        elapsed=$((elapsed + check_interval))
    done
    
    if [ $elapsed -ge $timeout ]; then
        print_error "RabbitMQ health check timed out after ${timeout}s"
        echo -e "\n  ${YELLOW}Container logs:${NC}"
        docker logs --tail 20 edugo-rabbitmq 2>&1 | sed 's/^/    /'
        echo -e "\n  ${YELLOW}Troubleshooting:${NC}"
        echo -e "    1. Check full logs: ${CYAN}docker logs edugo-rabbitmq${NC}"
        echo -e "    2. Restart container: ${CYAN}docker restart edugo-rabbitmq${NC}"
        return 1
    fi
    
    return 0
}

setup_postgres() {
    local pg_host="localhost"
    local pg_port="5432"
    local pg_user="edugo"
    local pg_db="edugo"
    local pg_password="edugo123"
    
    # Connection string for psql
    local psql_cmd="docker exec -i edugo-postgres psql -U ${pg_user} -d ${pg_db}"
    
    # Check if tables already exist
    print_info "Checking if PostgreSQL schema exists..."
    local table_count=$(echo "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';" | $psql_cmd -t 2>/dev/null | tr -d ' ')
    
    if [ -n "$table_count" ] && [ "$table_count" -gt 0 ]; then
        print_success "PostgreSQL schema already exists (${table_count} tables found)"
        print_info "Skipping schema creation (idempotent)"
    else
        print_info "Creating PostgreSQL schema..."
        
        # Execute migration scripts in order
        local migration_scripts=(
            "01_create_schema.sql"
            "03_refresh_tokens.sql"
            "04_login_attempts.sql"
            "04_material_versions.sql"
            "05_indexes_materials.sql"
            "05_user_progress_upsert.sql"
        )
        
        for script in "${migration_scripts[@]}"; do
            local script_path="${PROJECT_ROOT}/scripts/postgresql/${script}"
            if [ ! -f "$script_path" ]; then
                print_error "Migration script not found: ${script}"
                return 1
            fi
            
            print_info "Executing ${script}..."
            if cat "$script_path" | $psql_cmd >/dev/null 2>&1; then
                print_success "${script} executed successfully"
            else
                print_error "Failed to execute ${script}"
                echo -e "\n  ${YELLOW}Troubleshooting:${NC}"
                echo -e "    1. Check script syntax: ${CYAN}cat ${script_path}${NC}"
                echo -e "    2. Test manually: ${CYAN}docker exec -it edugo-postgres psql -U edugo -d edugo${NC}"
                echo -e "    3. Check PostgreSQL logs: ${CYAN}docker logs edugo-postgres${NC}"
                return 1
            fi
        done
        
        print_success "PostgreSQL schema created successfully"
    fi
    
    # Insert seed data
    print_info "Checking if seed data exists..."
    local user_count=$(echo "SELECT COUNT(*) FROM users;" | $psql_cmd -t 2>/dev/null | tr -d ' ')
    
    if [ -n "$user_count" ] && [ "$user_count" -gt 0 ]; then
        print_success "Seed data already exists (${user_count} users found)"
        print_info "Skipping seed data insertion (idempotent)"
    else
        print_info "Inserting seed data..."
        local seed_script="${PROJECT_ROOT}/scripts/postgresql/02_seed_data.sql"
        
        if [ ! -f "$seed_script" ]; then
            print_error "Seed data script not found: 02_seed_data.sql"
            return 1
        fi
        
        if cat "$seed_script" | $psql_cmd >/dev/null 2>&1; then
            print_success "Seed data inserted successfully"
            
            # Display created users
            echo -e "\n  ${CYAN}Test users created:${NC}"
            echo "SELECT email, 'password123' as password FROM users ORDER BY email;" | $psql_cmd -t 2>/dev/null | while read -r line; do
                if [ -n "$line" ]; then
                    echo -e "    ${GREEN}•${NC} ${line}"
                fi
            done
        else
            print_error "Failed to insert seed data"
            echo -e "\n  ${YELLOW}Troubleshooting:${NC}"
            echo -e "    1. Check script syntax: ${CYAN}cat ${seed_script}${NC}"
            echo -e "    2. Check PostgreSQL logs: ${CYAN}docker logs edugo-postgres${NC}"
            return 1
        fi
    fi
    
    return 0
}

setup_mongodb() {
    local mongo_user="edugo"
    local mongo_password="edugo123"
    local mongo_db="edugo"
    local mongo_auth_db="admin"
    
    # Connection command for mongosh
    local mongosh_cmd="docker exec -i edugo-mongodb mongosh -u ${mongo_user} -p ${mongo_password} --authenticationDatabase ${mongo_auth_db} ${mongo_db}"
    
    # Check if collections already exist
    print_info "Checking if MongoDB collections exist..."
    local collection_count=$(echo "db.getCollectionNames().length" | $mongosh_cmd --quiet 2>/dev/null | tr -d ' ')
    
    if [ -n "$collection_count" ] && [ "$collection_count" -gt 0 ]; then
        print_success "MongoDB collections already exist (${collection_count} collections found)"
        print_info "Skipping collection creation (idempotent)"
    else
        print_info "Creating MongoDB collections..."
        
        local mongo_script="${PROJECT_ROOT}/scripts/mongodb/02_assessment_results.js"
        
        if [ ! -f "$mongo_script" ]; then
            print_error "MongoDB script not found: 02_assessment_results.js"
            return 1
        fi
        
        print_info "Executing 02_assessment_results.js..."
        if cat "$mongo_script" | $mongosh_cmd >/dev/null 2>&1; then
            print_success "MongoDB collections created successfully"
            
            # Verify collections were created
            local created_collections=$(echo "db.getCollectionNames()" | $mongosh_cmd --quiet 2>/dev/null)
            if [ -n "$created_collections" ]; then
                echo -e "\n  ${CYAN}Collections created:${NC}"
                echo "$created_collections" | grep -o '"[^"]*"' | tr -d '"' | while read -r collection; do
                    if [ -n "$collection" ]; then
                        echo -e "    ${GREEN}•${NC} ${collection}"
                    fi
                done
            fi
        else
            print_error "Failed to execute MongoDB script"
            echo -e "\n  ${YELLOW}Troubleshooting:${NC}"
            echo -e "    1. Check script syntax: ${CYAN}cat ${mongo_script}${NC}"
            echo -e "    2. Test manually: ${CYAN}docker exec -it edugo-mongodb mongosh -u edugo -p edugo123${NC}"
            echo -e "    3. Check MongoDB logs: ${CYAN}docker logs edugo-mongodb${NC}"
            return 1
        fi
    fi
    
    return 0
}

setup_rabbitmq() {
    local rabbitmq_user="edugo"
    local rabbitmq_password="edugo123"
    local exchange_name="edugo.events"
    
    # Define queues to create
    local queues=(
        "material.created"
        "material.updated"
        "material.deleted"
        "assessment.completed"
        "progress.updated"
    )
    
    # Check if exchange already exists
    print_info "Checking if RabbitMQ exchange exists..."
    local exchange_exists=$(docker exec edugo-rabbitmq rabbitmqctl list_exchanges -q 2>/dev/null | grep -c "^${exchange_name}")
    
    if [ "$exchange_exists" -gt 0 ]; then
        print_success "RabbitMQ exchange '${exchange_name}' already exists"
        print_info "Skipping RabbitMQ setup (idempotent)"
    else
        print_info "Setting up RabbitMQ topology..."
        
        # Create exchange using rabbitmqadmin (if available) or direct AMQP commands
        print_info "Creating exchange '${exchange_name}'..."
        if docker exec edugo-rabbitmq rabbitmqadmin declare exchange name="${exchange_name}" type=topic durable=true -u "${rabbitmq_user}" -p "${rabbitmq_password}" >/dev/null 2>&1; then
            print_success "Exchange '${exchange_name}' created"
        else
            # Fallback: try using rabbitmqctl
            print_warning "rabbitmqadmin not available, using alternative method"
            # Note: Exchange will be created automatically when first queue binds to it
            print_info "Exchange will be created on first queue binding"
        fi
        
        # Create queues and bindings
        for queue in "${queues[@]}"; do
            print_info "Creating queue '${queue}'..."
            
            # Declare queue
            if docker exec edugo-rabbitmq rabbitmqadmin declare queue name="${queue}" durable=true -u "${rabbitmq_user}" -p "${rabbitmq_password}" >/dev/null 2>&1; then
                print_success "Queue '${queue}' created"
            else
                print_warning "Could not create queue '${queue}' (may already exist or rabbitmqadmin unavailable)"
            fi
            
            # Create binding (routing key is the queue name)
            print_info "Creating binding for '${queue}'..."
            if docker exec edugo-rabbitmq rabbitmqadmin declare binding source="${exchange_name}" destination="${queue}" routing_key="${queue}" -u "${rabbitmq_user}" -p "${rabbitmq_password}" >/dev/null 2>&1; then
                print_success "Binding created for '${queue}'"
            else
                print_warning "Could not create binding for '${queue}' (may already exist or rabbitmqadmin unavailable)"
            fi
        done
        
        print_success "RabbitMQ topology setup completed"
    fi
    
    # Display queue summary
    print_info "RabbitMQ queues:"
    docker exec edugo-rabbitmq rabbitmqctl list_queues name messages -q 2>/dev/null | while read -r line; do
        if [ -n "$line" ]; then
            echo -e "    ${GREEN}•${NC} ${line}"
        fi
    done
    
    return 0
}

verify_connectivity() {
    local max_retries=3
    local retry_delay=2
    local all_connected=true
    
    # Function to retry a command
    retry_command() {
        local description=$1
        local command=$2
        local attempt=1
        
        while [ $attempt -le $max_retries ]; do
            if eval "$command" >/dev/null 2>&1; then
                return 0
            fi
            
            if [ $attempt -lt $max_retries ]; then
                print_warning "${description} failed (attempt ${attempt}/${max_retries}), retrying in ${retry_delay}s..."
                sleep $retry_delay
            fi
            
            attempt=$((attempt + 1))
        done
        
        return 1
    }
    
    # Test PostgreSQL connection
    print_info "Testing PostgreSQL connection..."
    local pg_test_cmd="docker exec edugo-postgres psql -U edugo -d edugo -c 'SELECT 1'"
    
    if retry_command "PostgreSQL connection" "$pg_test_cmd"; then
        print_success "PostgreSQL connection successful"
        echo -e "    ${CYAN}Connection:${NC} postgresql://edugo:edugo123@localhost:5432/edugo"
    else
        print_error "PostgreSQL connection failed after ${max_retries} attempts"
        echo -e "\n  ${YELLOW}Troubleshooting:${NC}"
        echo -e "    1. Check container status: ${CYAN}docker ps | grep edugo-postgres${NC}"
        echo -e "    2. Check logs: ${CYAN}docker logs edugo-postgres${NC}"
        echo -e "    3. Test manually: ${CYAN}docker exec -it edugo-postgres psql -U edugo -d edugo${NC}"
        echo -e "    4. Verify credentials in .env file"
        all_connected=false
    fi
    
    # Test MongoDB connection
    print_info "Testing MongoDB connection..."
    local mongo_test_cmd="docker exec edugo-mongodb mongosh -u edugo -p edugo123 --authenticationDatabase admin edugo --eval 'db.adminCommand({ping: 1})'"
    
    if retry_command "MongoDB connection" "$mongo_test_cmd"; then
        print_success "MongoDB connection successful"
        echo -e "    ${CYAN}Connection:${NC} mongodb://edugo:edugo123@localhost:27017/edugo?authSource=admin"
    else
        print_error "MongoDB connection failed after ${max_retries} attempts"
        echo -e "\n  ${YELLOW}Troubleshooting:${NC}"
        echo -e "    1. Check container status: ${CYAN}docker ps | grep edugo-mongodb${NC}"
        echo -e "    2. Check logs: ${CYAN}docker logs edugo-mongodb${NC}"
        echo -e "    3. Test manually: ${CYAN}docker exec -it edugo-mongodb mongosh -u edugo -p edugo123${NC}"
        echo -e "    4. Verify credentials in .env file"
        all_connected=false
    fi
    
    # Test RabbitMQ connection
    print_info "Testing RabbitMQ connection..."
    local rabbitmq_test_cmd="docker exec edugo-rabbitmq rabbitmqctl status"
    
    if retry_command "RabbitMQ connection" "$rabbitmq_test_cmd"; then
        print_success "RabbitMQ connection successful"
        echo -e "    ${CYAN}Connection:${NC} amqp://edugo:edugo123@localhost:5672/"
        echo -e "    ${CYAN}Management UI:${NC} http://localhost:15672 (user: edugo, password: edugo123)"
    else
        print_error "RabbitMQ connection failed after ${max_retries} attempts"
        echo -e "\n  ${YELLOW}Troubleshooting:${NC}"
        echo -e "    1. Check container status: ${CYAN}docker ps | grep edugo-rabbitmq${NC}"
        echo -e "    2. Check logs: ${CYAN}docker logs edugo-rabbitmq${NC}"
        echo -e "    3. Test manually: ${CYAN}docker exec -it edugo-rabbitmq rabbitmqctl status${NC}"
        echo -e "    4. Verify credentials in .env file"
        all_connected=false
    fi
    
    if [ "$all_connected" = false ]; then
        echo -e "\n${RED}Some services are not accessible. Please fix the issues above.${NC}\n"
        return 1
    fi
    
    return 0
}

display_summary() {
    echo -e "${CYAN}╔════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║${NC}  ${YELLOW}Development Environment Summary${NC}                        ${CYAN}║${NC}"
    echo -e "${CYAN}╚════════════════════════════════════════════════════════════╝${NC}\n"
    
    # Connection Strings
    echo -e "${YELLOW}📡 Connection Strings:${NC}\n"
    echo -e "  ${CYAN}PostgreSQL:${NC}"
    echo -e "    postgresql://edugo:edugo123@localhost:5432/edugo\n"
    echo -e "  ${CYAN}MongoDB:${NC}"
    echo -e "    mongodb://edugo:edugo123@localhost:27017/edugo?authSource=admin\n"
    echo -e "  ${CYAN}RabbitMQ:${NC}"
    echo -e "    amqp://edugo:edugo123@localhost:5672/\n"
    
    # Default Credentials
    echo -e "${YELLOW}🔑 Default Credentials:${NC}\n"
    echo -e "  ${CYAN}PostgreSQL:${NC}"
    echo -e "    Host:     localhost"
    echo -e "    Port:     5432"
    echo -e "    Database: edugo"
    echo -e "    User:     edugo"
    echo -e "    Password: edugo123\n"
    
    echo -e "  ${CYAN}MongoDB:${NC}"
    echo -e "    Host:     localhost"
    echo -e "    Port:     27017"
    echo -e "    Database: edugo"
    echo -e "    User:     edugo"
    echo -e "    Password: edugo123"
    echo -e "    AuthDB:   admin\n"
    
    echo -e "  ${CYAN}RabbitMQ:${NC}"
    echo -e "    Host:        localhost"
    echo -e "    Port:        5672"
    echo -e "    Management:  http://localhost:15672"
    echo -e "    User:        edugo"
    echo -e "    Password:    edugo123\n"
    
    # Test User Accounts
    echo -e "${YELLOW}👥 Test User Accounts:${NC}\n"
    echo -e "  ${GREEN}Teachers:${NC}"
    echo -e "    • juan.perez@edugo.com | password: ${CYAN}password123${NC}"
    echo -e "    • maria.gonzalez@edugo.com | password: ${CYAN}password123${NC}\n"
    
    echo -e "  ${GREEN}Students:${NC}"
    echo -e "    • carlos.rodriguez@student.edugo.com | password: ${CYAN}password123${NC}"
    echo -e "    • ana.martinez@student.edugo.com | password: ${CYAN}password123${NC}"
    echo -e "    • luis.fernandez@student.edugo.com | password: ${CYAN}password123${NC}\n"
    
    # Useful Commands
    echo -e "${YELLOW}🛠️  Useful Commands:${NC}\n"
    echo -e "  ${CYAN}View logs:${NC}"
    echo -e "    docker-compose -f docker-compose-local.yml logs -f\n"
    
    echo -e "  ${CYAN}Stop services:${NC}"
    echo -e "    docker-compose -f docker-compose-local.yml stop\n"
    
    echo -e "  ${CYAN}Restart services:${NC}"
    echo -e "    docker-compose -f docker-compose-local.yml restart\n"
    
    echo -e "  ${CYAN}Reset environment:${NC}"
    echo -e "    make dev-reset\n"
    
    echo -e "  ${CYAN}Check status:${NC}"
    echo -e "    make dev-status\n"
    
    # Next Steps
    echo -e "${YELLOW}🚀 Next Steps:${NC}\n"
    echo -e "  1. Start the API server:"
    echo -e "     ${CYAN}make run${NC} or ${CYAN}go run cmd/main.go${NC}\n"
    
    echo -e "  2. Access the API documentation:"
    echo -e "     ${CYAN}http://localhost:8080/swagger/index.html${NC}\n"
    
    echo -e "  3. Test authentication with any test user above\n"
    
    echo -e "  4. View logs to monitor activity:"
    echo -e "     ${CYAN}docker-compose -f docker-compose-local.yml logs -f${NC}\n"
    
    echo -e "${GREEN}✨ Happy coding!${NC}\n"
}

# Execute main function
main "$@"
