#!/bin/bash

# URL Shortener Service - Docker Stop Script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
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

# Check if Docker Compose is available
if ! command -v docker-compose &> /dev/null; then
    print_error "Docker Compose is not installed."
    exit 1
fi

# Stop all services
print_status "Stopping URL Shortener services..."
docker-compose down

print_success "All services stopped successfully!"

# Optional: Remove volumes (uncomment if you want to clear data)
# read -p "Do you want to remove volumes (this will delete all data)? (y/N): " -n 1 -r
# echo
# if [[ $REPLY =~ ^[Yy]$ ]]; then
#     print_status "Removing volumes..."
#     docker-compose down -v
#     print_success "Volumes removed"
# fi 