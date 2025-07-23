#!/bin/bash

# URL Shortener Service - Development Script
# This script runs backend and frontend separately for development

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

# Function to check prerequisites
check_prerequisites() {
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go 1.21+ and try again."
        exit 1
    fi
    
    # Check if Node.js is installed
    if ! command -v node &> /dev/null; then
        print_error "Node.js is not installed. Please install Node.js 18+ and try again."
        exit 1
    fi
    
    # Check if PostgreSQL is running (optional)
    if ! pg_isready -h localhost -p 5432 > /dev/null 2>&1; then
        print_warning "PostgreSQL is not running on localhost:5432"
        print_status "You can start PostgreSQL with Docker:"
        echo "  docker run --name postgres-dev -e POSTGRES_DB=url_shortener -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres:15-alpine"
        echo ""
    else
        print_success "PostgreSQL is running"
    fi
}

# Function to start backend
start_backend() {
    print_status "Starting Backend (Go)..."
    
    # Check if .env file exists
    if [ ! -f .env ]; then
        print_warning ".env file not found, creating from example..."
        cp env.example .env
    fi
    
    # Install Go dependencies
    print_status "Installing Go dependencies..."
    go mod download
    
    # Start backend
    print_status "Starting backend server on http://localhost:8080"
    go run main.go &
    BACKEND_PID=$!
    
    # Wait for backend to be ready
    print_status "Waiting for backend to be ready..."
    timeout=30
    while ! curl -f http://localhost:8080/health > /dev/null 2>&1; do
        if [ $timeout -le 0 ]; then
            print_error "Backend failed to start within 30 seconds"
            kill $BACKEND_PID 2>/dev/null || true
            exit 1
        fi
        sleep 1
        timeout=$((timeout - 1))
    done
    print_success "Backend is ready!"
}

# Function to start frontend
start_frontend() {
    print_status "Starting Frontend (React)..."
    
    # Navigate to frontend directory
    cd frontend
    
    # Check if node_modules exists
    if [ ! -d "node_modules" ]; then
        print_status "Installing Node.js dependencies..."
        npm install
    fi
    
    # Start frontend
    print_status "Starting frontend development server on http://localhost:3000"
    npm run dev &
    FRONTEND_PID=$!
    
    # Wait for frontend to be ready
    print_status "Waiting for frontend to be ready..."
    timeout=30
    while ! curl -f http://localhost:3000 > /dev/null 2>&1; do
        if [ $timeout -le 0 ]; then
            print_error "Frontend failed to start within 30 seconds"
            kill $FRONTEND_PID 2>/dev/null || true
            exit 1
        fi
        sleep 1
        timeout=$((timeout - 1))
    done
    print_success "Frontend is ready!"
    
    # Go back to root directory
    cd ..
}

# Function to handle cleanup
cleanup() {
    print_status "Shutting down services..."
    
    if [ ! -z "$BACKEND_PID" ]; then
        kill $BACKEND_PID 2>/dev/null || true
        print_status "Backend stopped"
    fi
    
    if [ ! -z "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null || true
        print_status "Frontend stopped"
    fi
    
    print_success "All services stopped"
    exit 0
}

# Set up signal handlers
trap cleanup SIGINT SIGTERM

# Main execution
main() {
    echo "=========================================="
    echo "  URL Shortener Service - Development"
    echo "=========================================="
    echo ""
    
    # Check prerequisites
    check_prerequisites
    
    # Start services
    start_backend
    start_frontend
    
    # Show status
    echo ""
    print_success "Development environment is ready!"
    echo ""
    echo "  Frontend: http://localhost:3000"
    echo "  Backend:  http://localhost:8080"
    echo "  API Docs: http://localhost:8080/health"
    echo ""
    print_status "Press Ctrl+C to stop all services"
    echo ""
    
    # Keep script running
    wait
}

# Handle command line arguments
case "${1:-}" in
    "backend")
        check_prerequisites
        start_backend
        print_status "Backend running. Press Ctrl+C to stop."
        wait
        ;;
    "frontend")
        check_prerequisites
        start_frontend
        print_status "Frontend running. Press Ctrl+C to stop."
        wait
        ;;
    "help"|"-h"|"--help")
        echo "Usage: $0 [command]"
        echo ""
        echo "Commands:"
        echo "  (no args)  Start both backend and frontend"
        echo "  backend    Start only backend"
        echo "  frontend   Start only frontend"
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