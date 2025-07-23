#!/bin/bash

# Initialize Go modules using Docker (no Go required on host)
# This script creates go.sum and initializes the Go module using Docker

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

# Function to initialize Go modules using Docker
init_go_modules() {
    print_status "Initializing Go modules using Docker..."
    
    # Create a temporary Dockerfile for Go module initialization
    cat > Dockerfile.init << EOF
FROM golang:1.21-alpine

WORKDIR /app

# Copy go.mod
COPY go.mod ./

# Initialize modules and download dependencies
RUN go mod download && go mod tidy

# Copy the generated go.sum back to host
CMD ["sh", "-c", "cp go.sum /output/go.sum"]
EOF

    # Create output directory
    mkdir -p .docker-output

    # Run the initialization container
    docker build -f Dockerfile.init -t go-init .
    docker run --rm -v "$(pwd)/.docker-output:/output" go-init

    # Copy go.sum to project root
    if [ -f .docker-output/go.sum ]; then
        cp .docker-output/go.sum .
        print_success "go.sum generated successfully"
    else
        print_error "Failed to generate go.sum"
        exit 1
    fi

    # Cleanup
    rm -rf .docker-output Dockerfile.init
    docker rmi go-init > /dev/null 2>&1 || true
}

# Function to verify go.sum
verify_go_sum() {
    print_status "Verifying go.sum..."
    
    if [ -f go.sum ]; then
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
    echo "  Go Module Initialization (Docker)"
    echo "=========================================="
    echo ""
    
    # Check prerequisites
    check_docker
    
    # Check if go.mod exists
    if [ ! -f go.mod ]; then
        print_error "go.mod not found. Please create go.mod first."
        exit 1
    fi
    
    # Initialize Go modules
    init_go_modules
    
    # Verify result
    verify_go_sum
    
    echo ""
    print_success "Go modules initialized successfully!"
    print_status "You can now run: ./scripts/start.sh"
}

# Handle command line arguments
case "${1:-}" in
    "help"|"-h"|"--help")
        echo "Usage: $0"
        echo ""
        echo "This script initializes Go modules using Docker without requiring Go on the host."
        echo "It generates go.sum file needed for Docker builds."
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