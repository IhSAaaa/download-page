#!/bin/bash

# Complete setup script - Initialize, build, and start everything in one command
# No Go/Node.js required on host - everything runs in Docker

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

# Function to cleanup permission issues
cleanup_permissions() {
    print_status "Cleaning up permission issues..."
    
    # Fix permissions for node_modules if it exists
    if [ -d "frontend/node_modules" ]; then
        sudo chown -R $(id -u):$(id -g) frontend/node_modules 2>/dev/null || true
    fi
    
    # Fix permissions for .docker-output if it exists
    if [ -d ".docker-output" ]; then
        sudo chown -R $(id -u):$(id -g) .docker-output 2>/dev/null || true
        sudo rm -rf .docker-output 2>/dev/null || true
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
    sudo rm -rf node_modules 2>/dev/null || true
    
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

# Function to build and start services
build_and_start() {
    print_status "Building and starting all services..."
    
    # Build and start services
    docker-compose up --build -d
    
    print_success "Services started successfully!"
}

# Function to wait for services to be ready
wait_for_services() {
    print_status "Waiting for services to be ready..."
    
    # Wait for database
    print_status "Waiting for PostgreSQL..."
    timeout=120
    while ! docker-compose exec -T url-shortener-postgres pg_isready -U postgres -d url_shortener > /dev/null 2>&1; do
        if [ $timeout -le 0 ]; then
            print_error "PostgreSQL failed to start within 120 seconds"
            exit 1
        fi
        sleep 2
        timeout=$((timeout - 2))
    done
    print_success "PostgreSQL is ready"
    
    # Wait for backend
    print_status "Waiting for Backend API..."
    timeout=60
    while ! curl -f http://localhost:8080/health > /dev/null 2>&1; do
        if [ $timeout -le 0 ]; then
            print_error "Backend API failed to start within 60 seconds"
            exit 1
        fi
        sleep 2
        timeout=$((timeout - 2))
    done
    print_success "Backend API is ready"
    
    # Wait for frontend
    print_status "Waiting for Frontend..."
    timeout=60
    while ! curl -f http://localhost:80 > /dev/null 2>&1; do
        if [ $timeout -le 0 ]; then
            print_error "Frontend failed to start within 60 seconds"
            exit 1
        fi
        sleep 2
        timeout=$((timeout - 2))
    done
    print_success "Frontend is ready"
}

# Function to show final status
show_final_status() {
    echo ""
    print_success "🎉 URL Shortener Service is now running!"
    echo ""
    print_status "Service URLs:"
    echo -e "  ${GREEN}Frontend:${NC}     http://localhost:80"
    echo -e "  ${GREEN}Backend API:${NC}  http://localhost:8080"
    echo -e "  ${GREEN}Health Check:${NC} http://localhost:8080/health"
    echo ""
    print_status "Useful commands:"
    echo "  View logs:     ./scripts/start.sh logs"
    echo "  Stop services: ./scripts/stop.sh"
    echo "  Restart:       ./scripts/start.sh restart"
    echo "  Development:   ./scripts/dev-docker.sh"
    echo ""
}

# Main execution
main() {
    echo "=========================================="
    echo "  URL Shortener - Complete Setup"
    echo "=========================================="
    echo ""
    
    # Check prerequisites
    check_docker
    check_docker_compose
    
    # Check if services are already running
    if docker-compose ps 2>/dev/null | grep -q "Up"; then
        print_warning "Some services are already running"
        read -p "Do you want to restart all services? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            print_status "Stopping existing services..."
            docker-compose down
        else
            print_status "Using existing services"
            show_final_status
            exit 0
        fi
    fi
    
    # Cleanup permission issues first
    cleanup_permissions
    
    # Initialize dependencies
    init_backend
    init_frontend
    
    # Verify dependencies
    if ! verify_dependencies; then
        print_error "Dependency verification failed"
        exit 1
    fi
    
    # Build and start services
    build_and_start
    
    # Wait for services to be ready
    wait_for_services
    
    # Show final status
    show_final_status
    
    # Final cleanup
    cleanup_permissions
}

# Handle command line arguments
case "${1:-}" in
    "init")
        print_status "Initializing dependencies only..."
        check_docker
        init_backend
        init_frontend
        if verify_dependencies; then
            print_success "Dependencies initialized successfully!"
        else
            print_error "Dependency verification failed"
            exit 1
        fi
        ;;
    "start")
        print_status "Starting services only..."
        check_docker
        check_docker_compose
        build_and_start
        wait_for_services
        show_final_status
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
        echo "  (no args)  Complete setup (init + start)"
        echo "  init        Initialize dependencies only"
        echo "  start       Start services only"
        echo "  verify      Verify dependencies"
        echo "  help        Show this help message"
        echo ""
        echo "This script does everything in one command:"
        echo "  1. Initialize backend dependencies (go.sum)"
        echo "  2. Initialize frontend dependencies (node_modules)"
        echo "  3. Build and start all services"
        echo "  4. Wait for services to be ready"
        echo ""
        echo "No Go or Node.js required on host - everything runs in Docker!"
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