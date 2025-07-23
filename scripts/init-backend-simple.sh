#!/bin/bash

# Simple backend initialization using Docker volumes
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

# Function to initialize backend using Docker volumes
init_backend() {
    print_status "Initializing backend dependencies using Docker volumes..."
    
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
        print_success "Backend dependencies initialized successfully"
    else
        print_error "Failed to initialize backend dependencies"
        exit 1
    fi
}

# Function to verify go.sum
verify_go_sum() {
    print_status "Verifying go.sum..."
    
    if [ -f "go.sum" ]; then
        print_success "go.sum exists and is ready"
        echo "Dependencies in go.sum:"
        grep "^github.com\|^golang.org\|^google.golang.org" go.sum | head -10
        echo "..."
    else
        print_error "go.sum not found"
        exit 1
    fi
}

# Main execution
main() {
    echo "=========================================="
    echo "  Backend Dependencies (Docker Volumes)"
    echo "=========================================="
    echo ""
    
    # Check prerequisites
    check_docker
    
    # Initialize backend
    init_backend
    
    # Verify result
    verify_go_sum
    
    echo ""
    print_success "Backend dependencies initialized successfully!"
    print_status "You can now run: ./scripts/start.sh"
}

# Handle command line arguments
case "${1:-}" in
    "help"|"-h"|"--help")
        echo "Usage: $0"
        echo ""
        echo "This script initializes Go modules using Docker volumes."
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