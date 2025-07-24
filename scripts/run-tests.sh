#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Starting URL Shortener Testing Suite${NC}"
echo "=================================="

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker first."
    exit 1
fi

# Create coverage directories
mkdir -p coverage
mkdir -p frontend/coverage

print_status "Starting test database..."
# Start test services
docker-compose --profile test up -d postgres-test

# Wait for test database to be ready
print_status "Waiting for test database to be ready..."
for i in {1..30}; do
    if docker-compose --profile test exec -T postgres-test pg_isready -U postgres -d url_shortener_test > /dev/null 2>&1; then
        print_success "Test database is ready!"
        break
    fi
    if [ $i -eq 30 ]; then
        print_error "Test database failed to start within 30 seconds"
        exit 1
    fi
    sleep 1
done

# Run backend tests
print_status "Running backend tests..."
docker-compose --profile test run --rm backend-test go test ./tests/... -v -coverprofile=coverage/coverage.out -covermode=atomic

if [ $? -eq 0 ]; then
    print_success "Backend tests completed successfully!"
    
    # Generate coverage report
    print_status "Generating backend coverage report..."
    docker-compose --profile test run --rm backend-test go tool cover -html=coverage/coverage.out -o coverage/coverage.html
    
    if [ $? -eq 0 ]; then
        print_success "Backend coverage report generated: coverage/coverage.html"
    else
        print_warning "Failed to generate backend coverage report"
    fi
else
    print_error "Backend tests failed!"
    exit 1
fi

# Run frontend tests
print_status "Running frontend tests..."
docker-compose --profile test run --rm frontend-test npm test -- --coverage --watchAll=false

if [ $? -eq 0 ]; then
    print_success "Frontend tests completed successfully!"
    
    # Check if coverage report exists
    if [ -f "frontend/coverage/lcov-report/index.html" ]; then
        print_success "Frontend coverage report generated: frontend/coverage/lcov-report/index.html"
    else
        print_warning "Frontend coverage report not found"
    fi
else
    print_error "Frontend tests failed!"
    exit 1
fi

# Run integration tests
print_status "Running integration tests..."
docker-compose --profile test run --rm backend-test go test ./tests/integration/... -v

if [ $? -eq 0 ]; then
    print_success "Integration tests completed successfully!"
else
    print_error "Integration tests failed!"
    exit 1
fi

# Stop test services
print_status "Stopping test services..."
docker-compose --profile test down

# Display summary
echo ""
echo -e "${GREEN}🎉 All tests completed successfully!${NC}"
echo "=================================="
echo -e "${BLUE}Backend Coverage:${NC} coverage/coverage.html"
echo -e "${BLUE}Frontend Coverage:${NC} frontend/coverage/lcov-report/index.html"
echo ""
echo -e "${YELLOW}To view coverage reports:${NC}"
echo "  Backend: open coverage/coverage.html"
echo "  Frontend: open frontend/coverage/lcov-report/index.html" 