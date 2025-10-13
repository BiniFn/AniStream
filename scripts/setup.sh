#!/bin/sh
set -eu

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
ENV_FILE=".env.local"
ENV_EXAMPLE=".env.example"

# Print colored output
print_status() {
    local status=$1
    local message=$2
    case $status in
        "info") printf "${BLUE}ℹ${NC} %s\n" "$message" ;;
        "success") printf "${GREEN}✓${NC} %s\n" "$message" ;;
        "warning") printf "${YELLOW}⚠${NC} %s\n" "$message" ;;
        "error") printf "${RED}✗${NC} %s\n" "$message" >&2 ;;
        "step") printf "${CYAN}▶${NC} %s\n" "$message" ;;
    esac
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check dependencies
check_dependencies() {
    print_status "info" "Checking system dependencies..."
    
    local missing_deps=()
    
    # Check Go
    if ! command_exists go; then
        missing_deps+=("Go 1.24.1+")
    else
        local go_version=$(go version | grep -o 'go[0-9]\+\.[0-9]\+' | sed 's/go//')
        local go_major=$(echo "$go_version" | cut -d. -f1)
        local go_minor=$(echo "$go_version" | cut -d. -f2)
        if [ "$go_major" -lt 1 ] || ([ "$go_major" -eq 1 ] && [ "$go_minor" -lt 24 ]); then
            missing_deps+=("Go 1.24.1+ (current: $go_version)")
        fi
    fi
    
    # Check Node.js (for Bun)
    if ! command_exists node && ! command_exists bun; then
        missing_deps+=("Node.js 18+ or Bun")
    fi
    
    # Check Docker
    if ! command_exists docker; then
        missing_deps+=("Docker")
    fi
    
    # Check Docker Compose
    if ! command_exists docker-compose && ! docker compose version >/dev/null 2>&1; then
        missing_deps+=("Docker Compose")
    fi
    
    if [ ${#missing_deps[@]} -gt 0 ]; then
        print_status "error" "Missing required dependencies:"
        for dep in "${missing_deps[@]}"; do
            printf "  • %s\n" "$dep"
        done
        printf "\nPlease install the missing dependencies and run this script again.\n"
        exit 1
    fi
    
    print_status "success" "All system dependencies found"
}

# Install Go tools
install_go_tools() {
    print_status "step" "Installing Go development tools..."
    
    local tools=(
        "github.com/air-verse/air@latest"
        "github.com/sqlc-dev/sqlc/cmd/sqlc@latest"
        "github.com/Khan/genqlient@latest"
        "github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
        "github.com/swaggo/swag/cmd/swag@latest"
    )
    
    local installed_count=0
    local total_count=${#tools[@]}
    
    for tool in "${tools[@]}"; do
        local tool_name=$(echo "$tool" | cut -d'/' -f3)
        if command_exists "$tool_name"; then
            ((installed_count++))
        else
            if go install "$tool" >/dev/null 2>&1; then
                ((installed_count++))
            else
                print_status "error" "Failed to install $tool_name"
                exit 1
            fi
        fi
    done
    
    print_status "success" "Go tools ready ($installed_count/$total_count)"
}

# Install Node.js tools
install_node_tools() {
    print_status "step" "Installing Node.js development tools..."
    
    if command_exists bun; then
        if bun install -g swagger2openapi >/dev/null 2>&1; then
            print_status "success" "Node.js tools ready"
        else
            print_status "error" "Failed to install swagger2openapi"
            exit 1
        fi
    elif command_exists npm; then
        if npm install -g swagger2openapi >/dev/null 2>&1; then
            print_status "success" "Node.js tools ready"
        else
            print_status "error" "Failed to install swagger2openapi"
            exit 1
        fi
    else
        print_status "error" "Neither Bun nor npm found"
        exit 1
    fi
}

# Setup environment file
setup_environment() {
    print_status "step" "Setting up environment configuration..."
    
    if [ -f "$ENV_FILE" ]; then
        print_status "info" "$ENV_FILE already exists"
    else
        if [ -f "$ENV_EXAMPLE" ]; then
            cp "$ENV_EXAMPLE" "$ENV_FILE"
            print_status "success" "Created $ENV_FILE from $ENV_EXAMPLE"
        else
            print_status "error" "$ENV_EXAMPLE not found"
            exit 1
        fi
    fi
}

# Generate required files
generate_files() {
    print_status "step" "Generating required files..."
    
    local generated_count=0
    local total_count=3
    
    # Generate SQLC code
    if make sqlc >/dev/null 2>&1; then
        ((generated_count++))
    else
        print_status "error" "Failed to generate SQLC code"
        exit 1
    fi
    
    # Generate GraphQL client
    if make genqlient >/dev/null 2>&1; then
        ((generated_count++))
    else
        print_status "error" "Failed to generate GraphQL client"
        exit 1
    fi
    
    # Generate OpenAPI docs
    if make openapi >/dev/null 2>&1; then
        ((generated_count++))
    else
        print_status "error" "Failed to generate OpenAPI documentation"
        exit 1
    fi
    
    print_status "success" "Code generation complete ($generated_count/$total_count)"
}

# Install frontend dependencies
install_frontend_deps() {
    print_status "step" "Installing frontend dependencies..."
    
    if [ -d "web" ]; then
        cd web
        
        if command_exists bun; then
            if bun install >/dev/null 2>&1; then
                print_status "success" "Frontend dependencies ready"
            else
                print_status "error" "Failed to install frontend dependencies"
                exit 1
            fi
        elif command_exists npm; then
            if npm install >/dev/null 2>&1; then
                print_status "success" "Frontend dependencies ready"
            else
                print_status "error" "Failed to install frontend dependencies"
                exit 1
            fi
        else
            print_status "error" "Neither Bun nor npm found for frontend dependencies"
            exit 1
        fi
        
        cd ..
    else
        print_status "warning" "web directory not found, skipping frontend dependencies"
    fi
}

# Start development services
start_services() {
    print_status "step" "Starting development services..."
    
    # Check if Docker is running
    if ! docker info >/dev/null 2>&1; then
        print_status "warning" "Docker is not running. Please start Docker and run 'make dev-docker-up' manually"
        return 0
    fi
    
    if make dev-docker-up >/dev/null 2>&1; then
        print_status "success" "Development services started"
        print_status "info" "Waiting for services to be ready..."
        sleep 3
    else
        print_status "warning" "Failed to start Docker containers automatically"
        print_status "info" "You can start them manually with: make dev-docker-up"
    fi
}

# Main installation process
main() {
    printf "${CYAN}🚀 AniWays Setup${NC}\n"
    printf "${CYAN}===============${NC}\n\n"
    
    # Check if we're in the right directory
    if [ ! -f "go.mod" ] || [ ! -f "Makefile" ]; then
        print_status "error" "Please run this script from the AniWays project root directory"
        print_status "info" "Expected files: go.mod, Makefile"
        print_status "info" "Current directory: $(pwd)"
        exit 1
    fi
    
    # Run installation steps
    check_dependencies
    install_go_tools
    install_node_tools
    setup_environment
    generate_files
    install_frontend_deps
    start_services
    
    printf "\n${GREEN}🎉 Setup complete!${NC}\n\n"
    
    print_status "info" "Next steps:"
    printf "  1. Update %s with your configuration values\n" "$ENV_FILE"
    printf "  2. Run 'make dev-api' to start the API server\n"
    printf "  3. Run 'make dev-worker' to start the background worker\n"
    printf "  4. Run 'make dev-proxy' to start the HLS proxy\n"
    printf "  5. Run 'cd web && bun run dev' to start the frontend\n"
    printf "  6. Or run 'make tmux' to start all services at once\n\n"
    
    print_status "info" "Useful commands:"
    printf "  • make help           - Show all available commands\n"
    printf "  • make dev-docker-logs - View container logs\n"
    printf "  • make dev-docker-down - Stop containers\n\n"
    
    print_status "info" "Access points:"
    printf "  • Frontend: http://localhost:3000\n"
    printf "  • API: http://localhost:8080\n"
    printf "  • Swagger UI: http://localhost:8080/swagger/\n"
    printf "  • OpenAPI spec: http://localhost:8080/openapi.yaml\n"
}

# Run main function
main "$@"
