#!/bin/bash

# URL Shortener Service - Docker Development Script
# This script runs backend and frontend in development mode using Docker

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

# Function to start development services
start_dev_services() {
    print_status "Starting development services..."
    
    # Create development docker-compose override
    cat > docker-compose.dev.yml << EOF
services:
  backend:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=password
      - DB_NAME=url_shortener
      - DB_SSLMODE=disable
      - GIN_MODE=debug
    volumes:
      - .:/app
    command: go run main.go
    depends_on:
      postgres:
        condition: service_healthy

  frontend:
    build: 
      context: ./frontend
      target: builder
    ports:
      - "3000:3000"
    volumes:
      - ./frontend:/app
      - /app/node_modules
    command: npm run dev -- --host 0.0.0.0
    depends_on:
      - backend

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=url_shortener
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=password
    volumes:
      - postgres_dev_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  postgres_dev_data:
EOF

    # Start services
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build -d
    
    print_success "Development services started successfully!"
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
    while ! curl -f http://localhost:3000 > /dev/null 2>&1; do
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
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml ps
    
    echo ""
    print_success "Development URLs:"
    echo -e "  ${GREEN}Frontend:${NC} http://localhost:3000"
    echo -e "  ${GREEN}Backend API:${NC} http://localhost:8080"
    echo -e "  ${GREEN}Health Check:${NC} http://localhost:8080/health"
    echo -e "  ${GREEN}Database:${NC} localhost:5432"
}

# Function to show logs
show_logs() {
    print_status "Recent logs from all services:"
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml logs --tail=20
}

# Function to stop development services
stop_dev_services() {
    print_status "Stopping development services..."
    docker-compose -f docker-compose.yml -f docker-compose.dev.yml down
    print_success "Development services stopped"
}

# Main execution
main() {
    echo "=========================================="
    echo "  URL Shortener Service - Docker Dev"
    echo "=========================================="
    echo ""
    
    # Check prerequisites
    check_docker
    check_docker_compose
    
    # Check if services are already running
    if docker-compose -f docker-compose.yml -f docker-compose.dev.yml ps 2>/dev/null | grep -q "Up"; then
        print_warning "Some development services are already running"
        read -p "Do you want to restart all services? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            stop_dev_services
        else
            print_status "Using existing services"
            show_status
            exit 0
        fi
    fi
    
    # Start services
    start_dev_services
    
    # Wait for services to be ready
    wait_for_services
    
    # Show final status
    echo ""
    show_status
    
    echo ""
    print_success "Development environment is ready!"
    echo ""
    print_status "Useful commands:"
    echo "  View logs:     ./scripts/dev-docker.sh logs"
    echo "  Stop services: ./scripts/dev-docker.sh stop"
    echo "  Restart:       ./scripts/dev-docker.sh restart"
    echo ""
}

# Handle command line arguments
case "${1:-}" in
    "logs")
        show_logs
        ;;
    "restart")
        print_status "Restarting development services..."
        stop_dev_services
        start_dev_services
        wait_for_services
        show_status
        ;;
    "stop")
        stop_dev_services
        ;;
    "status")
        show_status
        ;;
    "help"|"-h"|"--help")
        echo "Usage: $0 [command]"
        echo ""
        echo "Commands:"
        echo "  (no args)  Start development services"
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