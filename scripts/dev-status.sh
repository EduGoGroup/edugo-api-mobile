#!/bin/bash

# Development Environment Status Script
# This script displays the current status of the development environment

set -euo pipefail

# Color constants for output formatting
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly CYAN='\033[0;36m'
readonly GRAY='\033[0;90m'
readonly NC='\033[0m' # No Color

# Script directory
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

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

print_section() {
    echo -e "\n${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${CYAN}  $1${NC}"
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"
}

# Function to display usage information
usage() {
    cat << EOF
${CYAN}EduGo Development Environment Status${NC}

${YELLOW}Usage:${NC}
    $0 [OPTIONS]

${YELLOW}Options:${NC}
    -h, --help          Show this help message
    -v, --verbose       Show detailed information

${YELLOW}Description:${NC}
    This script displays the current status of the development environment:
    - Container status (running/stopped/missing)
    - Database connectivity (PostgreSQL, MongoDB)
    - RabbitMQ status and queue information

${YELLOW}Examples:${NC}
    $0                  # Show status
    $0 -v               # Show detailed status

EOF
}

# Check container status
check_container_status() {
    print_section "🐳 Container Status"
    
    # Define required containers
    local containers=("edugo-postgres" "edugo-mongodb" "edugo-rabbitmq")
    local all_running=true
    
    for container in "${containers[@]}"; do
        # Check if container exists and get its status
        if docker ps --format '{{.Names}}' | grep -q "^${container}$"; then
            # Container is running
            local uptime=$(docker ps --filter "name=${container}" --format "{{.Status}}")
            print_success "${container} - ${GREEN}Running${NC} (${uptime})"
        elif docker ps -a --format '{{.Names}}' | grep -q "^${container}$"; then
            # Container exists but is not running
            local status=$(docker ps -a --filter "name=${container}" --format "{{.Status}}")
            print_error "${container} - ${RED}Stopped${NC} (${status})"
            all_running=false
        else
            # Container doesn't exist
            print_error "${container} - ${GRAY}Missing${NC}"
            all_running=false
        fi
    done
    
    echo ""
    if [ "$all_running" = true ]; then
        print_success "All containers are running"
    else
        print_warning "Some containers are not running"
        echo -e "  ${YELLOW}Tip:${NC} Run ${CYAN}make dev-init${NC} to start all containers"
    fi
    
    return 0
}

# Check database status
check_database_status() {
    print_section "💾 Database Status"
    
    local all_connected=true
    
    # Check PostgreSQL
    echo -e "${CYAN}PostgreSQL:${NC}"
    if docker ps --format '{{.Names}}' | grep -q "^edugo-postgres$"; then
        # Container is running, test connection
        if docker exec edugo-postgres pg_isready -U edugo >/dev/null 2>&1; then
            print_success "Connected and ready"
            
            # Get database info
            local db_version=$(docker exec edugo-postgres psql -U edugo -d edugo -t -c "SELECT version();" 2>/dev/null | head -n1 | xargs)
            if [ -n "$db_version" ]; then
                echo -e "  ${GRAY}Version: ${db_version:0:50}...${NC}"
            fi
            
            # Get table count
            local table_count=$(docker exec edugo-postgres psql -U edugo -d edugo -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';" 2>/dev/null | tr -d ' ')
            if [ -n "$table_count" ]; then
                echo -e "  ${GRAY}Tables: ${table_count}${NC}"
            fi
            
            # Get user count
            local user_count=$(docker exec edugo-postgres psql -U edugo -d edugo -t -c "SELECT COUNT(*) FROM users;" 2>/dev/null | tr -d ' ')
            if [ -n "$user_count" ]; then
                echo -e "  ${GRAY}Users: ${user_count}${NC}"
            fi
            
            echo -e "  ${GRAY}Connection: postgresql://edugo:***@localhost:5432/edugo${NC}"
        else
            print_error "Container running but not accepting connections"
            all_connected=false
        fi
    else
        print_error "Container not running"
        all_connected=false
    fi
    
    echo ""
    
    # Check MongoDB
    echo -e "${CYAN}MongoDB:${NC}"
    if docker ps --format '{{.Names}}' | grep -q "^edugo-mongodb$"; then
        # Container is running, test connection
        if docker exec edugo-mongodb mongosh --quiet --eval "db.adminCommand('ping')" >/dev/null 2>&1; then
            print_success "Connected and ready"
            
            # Get MongoDB version
            local mongo_version=$(docker exec edugo-mongodb mongosh --quiet --eval "db.version()" 2>/dev/null | tr -d '"')
            if [ -n "$mongo_version" ]; then
                echo -e "  ${GRAY}Version: ${mongo_version}${NC}"
            fi
            
            # Get collection count
            local collection_count=$(docker exec edugo-mongodb mongosh -u edugo -p edugo123 --authenticationDatabase admin edugo --quiet --eval "db.getCollectionNames().length" 2>/dev/null | tr -d ' ')
            if [ -n "$collection_count" ]; then
                echo -e "  ${GRAY}Collections: ${collection_count}${NC}"
            fi
            
            # Get document count in assessment_results
            local doc_count=$(docker exec edugo-mongodb mongosh -u edugo -p edugo123 --authenticationDatabase admin edugo --quiet --eval "db.assessment_results.countDocuments()" 2>/dev/null | tr -d ' ')
            if [ -n "$doc_count" ]; then
                echo -e "  ${GRAY}Assessment results: ${doc_count}${NC}"
            fi
            
            echo -e "  ${GRAY}Connection: mongodb://edugo:***@localhost:27017/edugo${NC}"
        else
            print_error "Container running but not accepting connections"
            all_connected=false
        fi
    else
        print_error "Container not running"
        all_connected=false
    fi
    
    echo ""
    if [ "$all_connected" = true ]; then
        print_success "All databases are accessible"
    else
        print_warning "Some databases are not accessible"
        echo -e "  ${YELLOW}Tip:${NC} Run ${CYAN}make dev-init${NC} to initialize databases"
    fi
    
    return 0
}

# Check RabbitMQ status
check_rabbitmq_status() {
    print_section "🐰 RabbitMQ Status"
    
    if docker ps --format '{{.Names}}' | grep -q "^edugo-rabbitmq$"; then
        # Container is running, test connection
        if docker exec edugo-rabbitmq rabbitmq-diagnostics -q ping >/dev/null 2>&1; then
            print_success "Connected and ready"
            
            # Get RabbitMQ version
            local rabbitmq_version=$(docker exec edugo-rabbitmq rabbitmqctl version 2>/dev/null | head -n1)
            if [ -n "$rabbitmq_version" ]; then
                echo -e "  ${GRAY}Version: ${rabbitmq_version}${NC}"
            fi
            
            # Get node status
            local node_status=$(docker exec edugo-rabbitmq rabbitmqctl status 2>/dev/null | grep -A1 "Status of node" | tail -n1 | xargs)
            
            # Get queue information
            echo -e "\n  ${CYAN}Queues:${NC}"
            local queue_info=$(docker exec edugo-rabbitmq rabbitmqctl list_queues name messages consumers -q 2>/dev/null)
            
            if [ -n "$queue_info" ]; then
                local total_queues=0
                local total_messages=0
                
                while IFS=$'\t' read -r queue_name messages consumers; do
                    if [ -n "$queue_name" ]; then
                        total_queues=$((total_queues + 1))
                        total_messages=$((total_messages + messages))
                        
                        # Format queue status
                        local queue_status=""
                        if [ "$messages" -gt 0 ]; then
                            queue_status="${YELLOW}${messages} messages${NC}"
                        else
                            queue_status="${GRAY}empty${NC}"
                        fi
                        
                        if [ "$consumers" -gt 0 ]; then
                            queue_status="${queue_status}, ${GREEN}${consumers} consumers${NC}"
                        fi
                        
                        echo -e "    ${GREEN}•${NC} ${queue_name} - ${queue_status}"
                    fi
                done <<< "$queue_info"
                
                echo -e "\n  ${GRAY}Total queues: ${total_queues}${NC}"
                echo -e "  ${GRAY}Total messages: ${total_messages}${NC}"
            else
                echo -e "    ${GRAY}No queues found${NC}"
            fi
            
            # Get exchange information
            echo -e "\n  ${CYAN}Exchanges:${NC}"
            local exchange_info=$(docker exec edugo-rabbitmq rabbitmqctl list_exchanges name type -q 2>/dev/null | grep -v "^amq\." | grep -v "^$")
            
            if [ -n "$exchange_info" ]; then
                while IFS=$'\t' read -r exchange_name exchange_type; do
                    if [ -n "$exchange_name" ] && [ "$exchange_name" != "" ]; then
                        echo -e "    ${GREEN}•${NC} ${exchange_name} (${exchange_type})"
                    fi
                done <<< "$exchange_info"
            fi
            
            echo -e "\n  ${GRAY}Connection: amqp://edugo:***@localhost:5672/${NC}"
            echo -e "  ${GRAY}Management UI: http://localhost:15672${NC}"
        else
            print_error "Container running but not accepting connections"
        fi
    else
        print_error "Container not running"
        echo -e "  ${YELLOW}Tip:${NC} Run ${CYAN}make dev-init${NC} to start RabbitMQ"
    fi
    
    return 0
}

# Display formatted output with status indicators
display_summary() {
    print_section "📊 Summary"
    
    # Count running containers
    local running_count=0
    local total_count=3
    local containers=("edugo-postgres" "edugo-mongodb" "edugo-rabbitmq")
    
    for container in "${containers[@]}"; do
        if docker ps --format '{{.Names}}' | grep -q "^${container}$"; then
            running_count=$((running_count + 1))
        fi
    done
    
    # Display summary
    if [ $running_count -eq $total_count ]; then
        print_success "Environment is fully operational (${running_count}/${total_count} containers running)"
        echo -e "\n  ${GREEN}✨ Ready for development!${NC}"
        echo -e "  ${GRAY}Start the API: ${CYAN}make run${NC}"
    elif [ $running_count -gt 0 ]; then
        print_warning "Environment is partially operational (${running_count}/${total_count} containers running)"
        echo -e "\n  ${YELLOW}Some services are not running${NC}"
        echo -e "  ${GRAY}Initialize environment: ${CYAN}make dev-init${NC}"
    else
        print_error "Environment is not operational (0/${total_count} containers running)"
        echo -e "\n  ${RED}No services are running${NC}"
        echo -e "  ${GRAY}Initialize environment: ${CYAN}make dev-init${NC}"
    fi
    
    echo -e "\n  ${CYAN}Useful commands:${NC}"
    echo -e "    ${CYAN}make dev-init${NC}   - Initialize/start environment"
    echo -e "    ${CYAN}make dev-reset${NC}  - Reset environment completely"
    echo -e "    ${CYAN}make dev-status${NC} - Show this status"
    echo ""
}

# Main function
main() {
    local verbose=false
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                usage
                exit 0
                ;;
            -v|--verbose)
                verbose=true
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
    echo -e "${CYAN}║${NC}  ${YELLOW}EduGo Development Environment Status${NC}                  ${CYAN}║${NC}"
    echo -e "${CYAN}╚════════════════════════════════════════════════════════════╝${NC}"
    
    # Check container status
    check_container_status
    
    # Check database status
    check_database_status
    
    # Check RabbitMQ status
    check_rabbitmq_status
    
    # Display summary
    display_summary
}

# Execute main function
main "$@"
