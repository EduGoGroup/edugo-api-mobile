#!/bin/bash

# Development Environment Reset Script
# This script completely resets the development environment by stopping containers,
# removing volumes, and reinitializing everything from scratch.

set -euo pipefail

# Color codes for output
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m' # No Color

# Script directory
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

# Docker compose file
readonly COMPOSE_FILE="${PROJECT_ROOT}/docker-compose-local.yml"

#######################################
# Print colored message
# Arguments:
#   $1 - Color code
#   $2 - Message
#######################################
print_message() {
    local color=$1
    shift
    echo -e "${color}$*${NC}"
}

#######################################
# Print section header
# Arguments:
#   $1 - Section name
#######################################
print_section() {
    echo ""
    print_message "${BLUE}" "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    print_message "${BLUE}" "  $1"
    print_message "${BLUE}" "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
}

#######################################
# Prompts user for confirmation before reset
# Returns:
#   0 if user confirms, 1 if user cancels
#######################################
confirm_reset() {
    print_section "⚠️  Development Environment Reset"
    
    print_message "${YELLOW}" "WARNING: This operation will:"
    echo "  • Stop all development containers"
    echo "  • Remove all containers"
    echo "  • Remove all volumes (ALL DATA WILL BE LOST)"
    echo "  • Reinitialize the environment from scratch"
    echo ""
    print_message "${RED}" "This action cannot be undone!"
    echo ""
    
    # Check if -y flag was passed
    if [[ "${1:-}" == "-y" ]] || [[ "${1:-}" == "--yes" ]]; then
        print_message "${YELLOW}" "Auto-confirmation enabled, proceeding with reset..."
        return 0
    fi
    
    read -p "Are you sure you want to continue? (yes/no): " -r
    echo ""
    
    if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
        print_message "${GREEN}" "✓ Reset cancelled"
        return 1
    fi
    
    return 0
}

#######################################
# Stops and removes all containers and volumes
# Returns:
#   0 on success
#######################################
cleanup_environment() {
    print_section "🧹 Cleaning Up Environment"
    
    # Stop containers
    print_message "${YELLOW}" "Stopping containers..."
    if docker-compose -f "${COMPOSE_FILE}" down 2>/dev/null; then
        print_message "${GREEN}" "✓ Containers stopped"
    else
        print_message "${YELLOW}" "⚠ No containers to stop or already stopped"
    fi
    
    # Remove containers and volumes
    print_message "${YELLOW}" "Removing containers and volumes..."
    if docker-compose -f "${COMPOSE_FILE}" down -v 2>/dev/null; then
        print_message "${GREEN}" "✓ Containers and volumes removed"
    else
        print_message "${YELLOW}" "⚠ No containers or volumes to remove"
    fi
    
    # Additional cleanup: remove any orphaned containers
    print_message "${YELLOW}" "Checking for orphaned containers..."
    local orphaned_containers
    orphaned_containers=$(docker ps -a --filter "name=edugo-" --format "{{.Names}}" 2>/dev/null || true)
    
    if [[ -n "$orphaned_containers" ]]; then
        print_message "${YELLOW}" "Found orphaned containers, removing..."
        echo "$orphaned_containers" | xargs -r docker rm -f 2>/dev/null || true
        print_message "${GREEN}" "✓ Orphaned containers removed"
    else
        print_message "${GREEN}" "✓ No orphaned containers found"
    fi
    
    # Additional cleanup: remove any orphaned volumes
    print_message "${YELLOW}" "Checking for orphaned volumes..."
    local orphaned_volumes
    orphaned_volumes=$(docker volume ls --filter "name=edugo" --format "{{.Name}}" 2>/dev/null || true)
    
    if [[ -n "$orphaned_volumes" ]]; then
        print_message "${YELLOW}" "Found orphaned volumes, removing..."
        echo "$orphaned_volumes" | xargs -r docker volume rm 2>/dev/null || true
        print_message "${GREEN}" "✓ Orphaned volumes removed"
    else
        print_message "${GREEN}" "✓ No orphaned volumes found"
    fi
    
    print_message "${GREEN}" "✓ Environment cleanup complete"
    return 0
}

#######################################
# Reinitializes the environment by calling dev-init.sh
# Returns:
#   Exit code from dev-init.sh
#######################################
reinitialize() {
    print_section "🚀 Reinitializing Environment"
    
    local init_script="${SCRIPT_DIR}/dev-init.sh"
    
    if [[ ! -f "$init_script" ]]; then
        print_message "${RED}" "✗ Error: dev-init.sh not found at ${init_script}"
        print_message "${YELLOW}" "Please ensure the initialization script exists"
        return 1
    fi
    
    if [[ ! -x "$init_script" ]]; then
        print_message "${YELLOW}" "Making dev-init.sh executable..."
        chmod +x "$init_script"
    fi
    
    print_message "${BLUE}" "Running dev-init.sh..."
    echo ""
    
    # Execute the initialization script
    if "$init_script"; then
        echo ""
        print_message "${GREEN}" "✓ Environment reinitialized successfully"
        return 0
    else
        echo ""
        print_message "${RED}" "✗ Environment reinitialization failed"
        print_message "${YELLOW}" "Please check the error messages above"
        return 1
    fi
}

#######################################
# Display usage information
#######################################
usage() {
    cat << EOF
Usage: $(basename "$0") [OPTIONS]

Reset the development environment by stopping containers, removing volumes,
and reinitializing everything from scratch.

OPTIONS:
    -y, --yes       Skip confirmation prompt
    -h, --help      Display this help message

EXAMPLES:
    $(basename "$0")              # Interactive mode with confirmation
    $(basename "$0") -y           # Auto-confirm and reset
    $(basename "$0") --help       # Show this help

WARNING:
    This operation will destroy all data in the development environment.
    Use with caution!

EOF
}

#######################################
# Main execution flow
#######################################
main() {
    # Parse arguments
    local auto_confirm=""
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            -y|--yes)
                auto_confirm="-y"
                shift
                ;;
            -h|--help)
                usage
                exit 0
                ;;
            *)
                print_message "${RED}" "✗ Unknown option: $1"
                echo ""
                usage
                exit 1
                ;;
        esac
    done
    
    # Change to project root
    cd "${PROJECT_ROOT}"
    
    # Step 1: Confirm reset
    if ! confirm_reset "$auto_confirm"; then
        exit 0
    fi
    
    # Step 2: Cleanup environment
    if ! cleanup_environment; then
        print_message "${RED}" "✗ Cleanup failed"
        exit 1
    fi
    
    # Step 3: Reinitialize environment
    if ! reinitialize; then
        print_message "${RED}" "✗ Reset failed during reinitialization"
        print_message "${YELLOW}" ""
        print_message "${YELLOW}" "Troubleshooting:"
        print_message "${YELLOW}" "  1. Check Docker is running: docker ps"
        print_message "${YELLOW}" "  2. Verify .env file exists: ls -la .env"
        print_message "${YELLOW}" "  3. Try manual initialization: ./scripts/dev-init.sh"
        exit 1
    fi
    
    # Success!
    print_section "✅ Reset Complete"
    print_message "${GREEN}" "Your development environment has been reset and reinitialized."
    print_message "${GREEN}" "You can now start developing!"
}

# Execute main function
main "$@"
