#!/bin/bash

# Initialize Frontend dependencies using Docker (no Node.js required on host)
# This script installs npm dependencies using Docker

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

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
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

# Function to initialize frontend dependencies using Docker
init_frontend_deps() {
    print_status "Initializing frontend dependencies using Docker..."
    
    # Check if frontend directory exists
    if [ ! -d "frontend" ]; then
        print_error "frontend directory not found"
        exit 1
    fi
    
    # Check if package.json exists
    if [ ! -f "frontend/package.json" ]; then
        print_error "frontend/package.json not found"
        exit 1
    fi
    
    # Create a temporary Dockerfile for frontend initialization
    cat > frontend/Dockerfile.init << EOF
FROM node:18-alpine

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the node_modules back to host
CMD ["sh", "-c", "cp -r node_modules /output/"]
EOF

    # Create output directory
    mkdir -p frontend/.docker-output

    # Run the initialization container
    cd frontend
    docker build -f Dockerfile.init -t frontend-init .
    docker run --rm -v "$(pwd)/.docker-output:/output" frontend-init

    # Copy node_modules to frontend directory
    if [ -d ".docker-output/node_modules" ]; then
        # Remove existing node_modules if it exists
        rm -rf node_modules 2>/dev/null || true
        
        # Copy with proper permissions
        cp -r .docker-output/node_modules .
        print_success "node_modules installed successfully"
    else
        print_error "Failed to install node_modules"
        exit 1
    fi

    # Cleanup with sudo if needed
    sudo rm -rf .docker-output Dockerfile.init 2>/dev/null || rm -rf .docker-output Dockerfile.init 2>/dev/null || true
    docker rmi frontend-init > /dev/null 2>&1 || true
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
    echo "  Frontend Dependencies (Docker)"
    echo "=========================================="
    echo ""
    
    # Check prerequisites
    check_docker
    
    # Initialize frontend dependencies
    init_frontend_deps
    
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
        echo "This script installs frontend dependencies using Docker without requiring Node.js on the host."
        echo "It creates node_modules directory needed for Docker builds."
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