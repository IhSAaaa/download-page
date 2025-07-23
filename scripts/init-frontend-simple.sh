#!/bin/bash

# Simple frontend initialization using Docker volumes
# This avoids permission issues by using Docker volumes

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker and try again."
        exit 1
    fi
    print_success "Docker is running"
}

# Function to initialize frontend using Docker volumes
init_frontend() {
    print_status "Initializing frontend dependencies using Docker volumes..."
    
    if [ ! -d "frontend" ]; then
        print_error "frontend directory not found"
        exit 1
    fi
    
    if [ ! -f "frontend/package.json" ]; then
        print_error "frontend/package.json not found"
        exit 1
    fi
    
    cd frontend
    
    # Remove existing node_modules if it exists
    rm -rf node_modules 2>/dev/null || true
    
    # Create a temporary container to install dependencies
    print_status "Installing npm dependencies..."
    docker run --rm \
        -v "$(pwd):/app" \
        -w /app \
        node:18-alpine \
        sh -c "npm install && chown -R $(id -u):$(id -g) node_modules"
    
    if [ -d "node_modules" ]; then
        print_success "Frontend dependencies installed successfully"
    else
        print_error "Failed to install frontend dependencies"
        exit 1
    fi
    
    cd ..
}

# Function to verify node_modules
verify_node_modules() {
    print_status "Verifying node_modules..."
    
    if [ -d "frontend/node_modules" ]; then
        print_success "node_modules exists and is ready"
        echo "Installed packages:"
        ls frontend/node_modules | head -10
        echo "..."
    else
        print_error "node_modules not found"
        exit 1
    fi
}

# Main execution
main() {
    echo "=========================================="
    echo "  Frontend Dependencies (Docker Volumes)"
    echo "=========================================="
    echo ""
    
    # Check prerequisites
    check_docker
    
    # Initialize frontend
    init_frontend
    
    # Verify result
    verify_node_modules
    
    echo ""
    print_success "Frontend dependencies initialized successfully!"
    print_status "You can now run: ./scripts/start.sh"
}

# Handle command line arguments
case "${1:-}" in
    "help"|"-h"|"--help")
        echo "Usage: $0"
        echo ""
        echo "This script installs frontend dependencies using Docker volumes."
        echo "It avoids permission issues by using Docker volumes instead of copying files."
        ;;
    "")
        main
        ;;
    *)
        print_error "Unknown command: $1"
        echo "Use '$0 help' for usage information"
        exit 1
        ;;
esac 