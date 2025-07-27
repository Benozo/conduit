#!/bin/bash

# SwarmV2 RAG Infrastructure Validation Script
# This script validates the Docker Compose setup and infrastructure

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

echo "🔍 SwarmV2 RAG Infrastructure Validation"
echo "========================================"

# Check Docker
print_status "Checking Docker..."
if command -v docker &> /dev/null; then
    if docker info > /dev/null 2>&1; then
        print_success "Docker is running"
        DOCKER_VERSION=$(docker --version | cut -d' ' -f3 | cut -d',' -f1)
        echo "   Version: $DOCKER_VERSION"
    else
        print_error "Docker is installed but not running"
        exit 1
    fi
else
    print_error "Docker is not installed"
    exit 1
fi

# Check Docker Compose
print_status "Checking Docker Compose..."
if docker compose version &> /dev/null; then
    COMPOSE_CMD="docker compose"
    COMPOSE_VERSION=$(docker compose version --short 2>/dev/null || echo "unknown")
    print_success "Docker Compose (v2) is available"
    echo "   Version: $COMPOSE_VERSION"
elif command -v docker-compose &> /dev/null; then
    COMPOSE_CMD="docker-compose"
    COMPOSE_VERSION=$(docker-compose version --short 2>/dev/null || echo "unknown")
    print_success "Docker Compose (v1) is available"
    echo "   Version: $COMPOSE_VERSION"
else
    print_error "Docker Compose is not available"
    exit 1
fi

# Check docker-compose.yml
print_status "Validating docker-compose.yml..."
if [ -f "docker-compose.yml" ]; then
    if $COMPOSE_CMD config > /dev/null 2>&1; then
        print_success "docker-compose.yml is valid"
    else
        print_error "docker-compose.yml has validation errors"
        $COMPOSE_CMD config
        exit 1
    fi
else
    print_error "docker-compose.yml not found"
    exit 1
fi

# Check environment setup
print_status "Checking environment setup..."
if [ -f ".env.example" ]; then
    print_success ".env.example template found"
    if [ -f ".env" ]; then
        print_success ".env file exists"
    else
        print_warning ".env file not found - will use defaults"
        echo "   Run: cp .env.example .env"
    fi
else
    print_error ".env.example template not found"
fi

# Check required directories
print_status "Checking directory structure..."
REQUIRED_DIRS=("config" "docs")
for dir in "${REQUIRED_DIRS[@]}"; do
    if [ -d "$dir" ]; then
        print_success "$dir/ directory exists"
    else
        print_warning "$dir/ directory missing"
    fi
done

# Check configuration files
print_status "Checking configuration files..."
CONFIG_FILES=("config/init-vector-db.sql" "docs/RAG_INFRASTRUCTURE.md")
for file in "${CONFIG_FILES[@]}"; do
    if [ -f "$file" ]; then
        print_success "$file exists"
    else
        print_error "$file missing"
    fi
done

# Check management script
print_status "Checking management script..."
if [ -f "rag-infrastructure.sh" ]; then
    if [ -x "rag-infrastructure.sh" ]; then
        print_success "rag-infrastructure.sh is executable"
    else
        print_warning "rag-infrastructure.sh is not executable"
        echo "   Run: chmod +x rag-infrastructure.sh"
    fi
else
    print_error "rag-infrastructure.sh not found"
fi

# Check available disk space
print_status "Checking available disk space..."
AVAILABLE_SPACE=$(df . | tail -1 | awk '{print $4}')
AVAILABLE_GB=$((AVAILABLE_SPACE / 1024 / 1024))

if [ $AVAILABLE_GB -gt 10 ]; then
    print_success "Sufficient disk space available (${AVAILABLE_GB}GB)"
else
    print_warning "Low disk space (${AVAILABLE_GB}GB) - recommend at least 10GB"
fi

# Check for port conflicts
print_status "Checking for potential port conflicts..."
REQUIRED_PORTS=(19530 8080 8081 5432 6379 11434 3000 5050 8082 8083 9000 9001)
CONFLICTS=()

for port in "${REQUIRED_PORTS[@]}"; do
    if netstat -tuln 2>/dev/null | grep -q ":$port "; then
        CONFLICTS+=($port)
    fi
done

if [ ${#CONFLICTS[@]} -eq 0 ]; then
    print_success "No port conflicts detected"
else
    print_warning "Potential port conflicts detected: ${CONFLICTS[*]}"
    echo "   These ports are already in use. You may need to stop other services or modify docker-compose.yml"
fi

# Summary
echo ""
echo "📋 Validation Summary"
echo "===================="

TOTAL_CHECKS=8
PASSED_CHECKS=0

# Count successful checks (simplified)
if command -v docker &> /dev/null && docker info > /dev/null 2>&1; then
    ((PASSED_CHECKS++))
fi

if docker compose version &> /dev/null || command -v docker-compose &> /dev/null; then
    ((PASSED_CHECKS++))
fi

if [ -f "docker-compose.yml" ] && $COMPOSE_CMD config > /dev/null 2>&1; then
    ((PASSED_CHECKS++))
fi

if [ -f ".env.example" ]; then
    ((PASSED_CHECKS++))
fi

if [ -d "config" ] && [ -d "docs" ]; then
    ((PASSED_CHECKS++))
fi

if [ -f "config/init-vector-db.sql" ] && [ -f "docs/RAG_INFRASTRUCTURE.md" ]; then
    ((PASSED_CHECKS++))
fi

if [ -f "rag-infrastructure.sh" ]; then
    ((PASSED_CHECKS++))
fi

if [ $AVAILABLE_GB -gt 5 ]; then
    ((PASSED_CHECKS++))
fi

if [ $PASSED_CHECKS -eq $TOTAL_CHECKS ]; then
    print_success "All checks passed! ($PASSED_CHECKS/$TOTAL_CHECKS)"
    echo ""
    echo "🚀 Ready to start RAG infrastructure:"
    echo "   ./rag-infrastructure.sh start"
elif [ $PASSED_CHECKS -gt $((TOTAL_CHECKS * 3 / 4)) ]; then
    print_warning "Most checks passed ($PASSED_CHECKS/$TOTAL_CHECKS) - infrastructure should work with minor issues"
    echo ""
    echo "🚀 You can try starting the infrastructure:"
    echo "   ./rag-infrastructure.sh start"
else
    print_error "Several checks failed ($PASSED_CHECKS/$TOTAL_CHECKS) - please fix issues before proceeding"
    echo ""
    echo "📖 See docs/RAG_INFRASTRUCTURE.md for detailed setup instructions"
fi

echo ""
echo "📚 Next Steps:"
echo "   1. Review and customize .env file"
echo "   2. Start infrastructure: ./rag-infrastructure.sh start"
echo "   3. View access URLs: ./rag-infrastructure.sh urls"
echo "   4. Run health checks: ./rag-infrastructure.sh health"
