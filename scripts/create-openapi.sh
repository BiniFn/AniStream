#!/bin/sh
set -eu

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Run command with status indicator
run_with_status() {
    local msg=$1
    shift
    
    printf "${BLUE}⏳${NC} %s... " "$msg"
    
    if "$@" >/dev/null 2>&1; then
        printf "${GREEN}✓${NC}\n"
    else
        printf "${RED}✗${NC}\n"
        exit 1
    fi
}

# Check dependencies
check_dependency() {
    if ! command -v $1 >/dev/null 2>&1; then
        printf "${RED}✗${NC} %s is not installed. Please install %s and try again.\n" "$1" "$1" >&2
        exit 1
    fi
}

printf "${YELLOW}🔧 Checking dependencies...${NC}\n"
check_dependency "swag"
check_dependency "swagger2openapi"
check_dependency "bun"
printf "${GREEN}✓${NC} All dependencies found\n\n"

ROOT="${PWD}"
cd "$ROOT"

printf "${YELLOW}📋 Generating OpenAPI specification${NC}\n"

run_with_status "Generating Swagger documentation" swag init -g ./cmd/api/main.go -o ./docs/swagger
run_with_status "Converting to OpenAPI 3.0" swagger2openapi --patch ./docs/swagger/swagger.yaml -o ./docs/openapi.yaml
run_with_status "Cleaning up temporary files" rm -rf ./docs/swagger

printf "\n${YELLOW}🎯 Generating TypeScript client${NC}\n"
cd "$ROOT/web"

run_with_status "Generating TypeScript types" bun run openapi

printf "\n${GREEN}🎉 OpenAPI generation complete!${NC}\n"
printf "${GREEN}   • OpenAPI spec:${NC} docs/openapi.yaml\n"
printf "${GREEN}   • TypeScript types:${NC} web/src/lib/api/openapi.ts\n"

exit 0
