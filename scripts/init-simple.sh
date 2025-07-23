#!/bin/bash

# Simple initialization using Docker volumes (no Go/Node.js required on host)
# This script sets up the entire project using Docker volumes to avoid permission issues

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

# Function to check if Docker Compose is available
check_docker_compose() {
    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose is not installed. Please install Docker Compose and try again."
        exit 1
    fi
    print_success "Docker Compose is available"
}

# Function to initialize backend dependencies
init_backend() {
    print_status "Initializing backend dependencies..."
    
    if [ ! -f "go.mod" ]; then
        print_error "go.mod not found. Please create go.mod first."
        exit 1
    fi
    
    # Remove existing go.sum if it exists
    rm -f go.sum 2>/dev/null || true
    
    # Create a temporary container to initialize Go modules
    print_status "Initializing Go modules..."
    docker run --rm \
        -v "$(pwd):/app" \
        -w /app \
        golang:1.21-alpine \
        sh -c "go mod download && go mod tidy && chown $(id -u):$(id -g) go.sum"
    
    if [ -f "go.sum" ]; then
        print_success "Backend dependencies initialized"
    else
        print_error "Failed to initialize backend dependencies"
        exit 1
    fi
}

# Function to initialize frontend dependencies
init_frontend() {
    print_status "Initializing frontend dependencies..."
    
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
        print_success "Frontend dependencies initialized"
    else
        print_error "Failed to install frontend dependencies"
        exit 1
    fi
    
    cd ..
}

# Function to verify all dependencies
verify_dependencies() {
    print_status "Verifying all dependencies..."
    
    # Check backend
    if [ -f "go.sum" ]; then
        print_success "✓ Backend dependencies (go.sum) ready"
    else
        print_error "✗ Backend dependencies missing"
        return 1
    fi
    
    # Check frontend
    if [ -d "frontend/node_modules" ]; then
        print_success "✓ Frontend dependencies (node_modules) ready"
    else
        print_error "✗ Frontend dependencies missing"
        return 1
    fi
    
    return 0
}

# Function to show next steps
show_next_steps() {
    echo ""
    print_success "All dependencies initialized successfully!"
    echo ""
    print_status "Next steps:"
    echo "  1. Start all services:     ./scripts/start.sh"
    echo "  2. Development mode:       ./scripts/dev-docker.sh"
    echo "  3. View logs:             ./scripts/start.sh logs"
    echo "  4. Stop services:         ./scripts/stop.sh"
    echo ""
    print_status "Service URLs (after starting):"
    echo "  Frontend:     http://localhost:80"
    echo "  Backend API:  http://localhost:8080"
    echo "  Health Check: http://localhost:8080/health"
    echo ""
}

# Main execution
main() {
    echo "=========================================="
    echo "  URL Shortener - Simple Docker Init"
    echo "=========================================="
    echo ""
    
    # Check prerequisites
    check_docker
    check_docker_compose
    
    # Initialize backend
    init_backend
    
    # Initialize frontend
    init_frontend
    
    # Verify everything
    if verify_dependencies; then
        show_next_steps
    else
        print_error "Dependency verification failed"
        exit 1
    fi
}

# Handle command line arguments
case "${1:-}" in
    "backend")
        print_status "Initializing backend only..."
        check_docker
        init_backend
        print_success "Backend initialized successfully!"
        ;;
    "frontend")
        print_status "Initializing frontend only..."
        check_docker
        init_frontend
        print_success "Frontend initialized successfully!"
        ;;
    "verify")
        print_status "Verifying dependencies..."
        if verify_dependencies; then
            print_success "All dependencies are ready!"
        else
            print_error "Some dependencies are missing"
            exit 1
        fi
        ;;
    "help"|"-h"|"--help")
        echo "Usage: $0 [command]"
        echo ""
        echo "Commands:"
        echo "  (no args)  Initialize all dependencies"
        echo "  backend     Initialize backend dependencies only"
        echo "  frontend    Initialize frontend dependencies only"
        echo "  verify      Verify all dependencies"
        echo "  help        Show this help message"
        echo ""
        echo "This script initializes all dependencies using Docker volumes."
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