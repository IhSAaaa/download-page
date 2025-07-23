#!/bin/bash

# URL Shortener Service - Docker Startup Script
# This script starts the complete application stack

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
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

# Function to build and start services
start_services() {
    print_status "Building and starting services..."
    
    # Build and start all services
    docker-compose up --build -d
    
    print_success "Services started successfully!"
}

# Function to wait for services to be ready
wait_for_services() {
    print_status "Waiting for services to be ready..."
    
    # Wait for database
    print_status "Waiting for PostgreSQL..."
    timeout=60
    while ! docker-compose exec -T postgres pg_isready -U postgres > /dev/null 2>&1; do
        if [ $timeout -le 0 ]; then
            print_error "PostgreSQL failed to start within 60 seconds"
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

# Function to show service status
show_status() {
    print_status "Service Status:"
    docker-compose ps
    
    echo ""
    print_success "Application URLs:"
    echo -e "  ${GREEN}Frontend:${NC} http://localhost"
    echo -e "  ${GREEN}Backend API:${NC} http://localhost:8080"
    echo -e "  ${GREEN}Health Check:${NC} http://localhost:8080/health"
    echo -e "  ${GREEN}Database:${NC} localhost:5432"
}

# Function to show logs
show_logs() {
    print_status "Recent logs from all services:"
    docker-compose logs --tail=20
}

# Main execution
main() {
    echo "=========================================="
    echo "  URL Shortener Service - Docker Startup"
    echo "=========================================="
    echo ""
    
    # Check prerequisites
    check_docker
    check_docker_compose
    
    # Check if services are already running
    if docker-compose ps | grep -q "Up"; then
        print_warning "Some services are already running"
        read -p "Do you want to restart all services? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            print_status "Stopping existing services..."
            docker-compose down
        else
            print_status "Using existing services"
            show_status
            exit 0
        fi
    fi
    
    # Start services
    start_services
    
    # Wait for services to be ready
    wait_for_services
    
    # Show final status
    echo ""
    show_status
    
    echo ""
    print_success "URL Shortener Service is now running!"
    echo ""
    print_status "Useful commands:"
    echo "  View logs:     ./scripts/start.sh logs"
    echo "  Stop services: docker-compose down"
    echo "  Restart:       ./scripts/start.sh restart"
    echo ""
}

# Handle command line arguments
case "${1:-}" in
    "logs")
        show_logs
        ;;
    "restart")
        print_status "Restarting services..."
        docker-compose down
        docker-compose up --build -d
        wait_for_services
        show_status
        ;;
    "stop")
        print_status "Stopping services..."
        docker-compose down
        print_success "Services stopped"
        ;;
    "status")
        show_status
        ;;
    "help"|"-h"|"--help")
        echo "Usage: $0 [command]"
        echo ""
        echo "Commands:"
        echo "  (no args)  Start all services"
        echo "  logs       Show recent logs"
        echo "  restart    Restart all services"
        echo "  stop       Stop all services"
        echo "  status     Show service status"
        echo "  help       Show this help message"
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