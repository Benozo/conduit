#!/bin/bash

# SwarmV2 RAG Infrastructure Management Script
# This script helps manage the Docker Compose setup for RAG components

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print colored output
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
        print_error "Docker is not running. Please start Docker first."
        exit 1
    fi
}

# Function to check if docker-compose is available
check_compose() {
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        print_error "Neither 'docker-compose' nor 'docker compose' is available."
        exit 1
    fi
    
    # Use docker compose if available, fallback to docker-compose
    if docker compose version &> /dev/null; then
        COMPOSE_CMD="docker compose"
    else
        COMPOSE_CMD="docker-compose"
    fi
}

# Function to create necessary directories
create_directories() {
    print_status "Creating data directories..."
    mkdir -p data/volumes/{etcd,minio,milvus,weaviate,postgres,pgadmin,redis,ollama,ollama-webui}
    print_success "Data directories created"
}

# Function to copy environment file
setup_env() {
    if [ ! -f .env ]; then
        print_status "Creating .env file from template..."
        cp .env.example .env
        print_warning "Please edit .env file with your specific configuration"
    else
        print_status ".env file already exists"
    fi
}

# Function to start all services
start_all() {
    print_status "Starting all RAG infrastructure services..."
    $COMPOSE_CMD up -d
    print_success "All services started"
    show_urls
}

# Function to start specific services
start_service() {
    local service=$1
    if [ -z "$service" ]; then
        print_error "Please specify a service name"
        show_help
        exit 1
    fi
    
    print_status "Starting $service..."
    $COMPOSE_CMD up -d $service
    print_success "$service started"
}

# Function to stop all services
stop_all() {
    print_status "Stopping all services..."
    $COMPOSE_CMD down
    print_success "All services stopped"
}

# Function to stop specific service
stop_service() {
    local service=$1
    if [ -z "$service" ]; then
        print_error "Please specify a service name"
        show_help
        exit 1
    fi
    
    print_status "Stopping $service..."
    $COMPOSE_CMD stop $service
    print_success "$service stopped"
}

# Function to show service status
show_status() {
    print_status "Service status:"
    $COMPOSE_CMD ps
}

# Function to show logs
show_logs() {
    local service=$1
    if [ -z "$service" ]; then
        print_status "Showing logs for all services..."
        $COMPOSE_CMD logs -f
    else
        print_status "Showing logs for $service..."
        $COMPOSE_CMD logs -f $service
    fi
}

# Function to clean up everything
cleanup() {
    print_warning "This will remove all containers, volumes, and data. Are you sure? (y/N)"
    read -r response
    if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
        print_status "Cleaning up all data..."
        $COMPOSE_CMD down -v --remove-orphans
        sudo rm -rf data/volumes
        print_success "Cleanup completed"
    else
        print_status "Cleanup cancelled"
    fi
}

# Function to show access URLs
show_urls() {
    echo ""
    print_success "RAG Infrastructure Access URLs:"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "🗄️  Vector Databases:"
    echo "   • Milvus API:           http://localhost:19530"
    echo "   • Milvus UI (Attu):     http://localhost:3000"
    echo "   • Weaviate API:         http://localhost:8080"
    echo "   • Weaviate Console:     http://localhost:8081"
    echo "   • PostgreSQL:           localhost:5432 (vectordb/postgres/password)"
    echo "   • pgAdmin:              http://localhost:5050 (admin@swarmv2.com/admin123)"
    echo ""
    echo "🚀 AI & LLM Services:"
    echo "   • Ollama API:           http://localhost:11434"
    echo "   • Ollama WebUI:         http://localhost:8083"
    echo ""
    echo "⚡ Caching & Storage:"
    echo "   • Redis:                localhost:6379"
    echo "   • Redis Commander:      http://localhost:8082"
    echo "   • MinIO Console:        http://localhost:9001 (minioadmin/minioadmin)"
    echo ""
    echo "📊 Monitoring:"
    echo "   • All Services Status:  docker compose ps"
    echo "   • Logs:                 docker compose logs -f [service]"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
}

# Function to initialize Ollama models
init_ollama() {
    print_status "Initializing Ollama models..."
    
    # Wait for Ollama to be ready
    print_status "Waiting for Ollama to be ready..."
    sleep 10
    
    # Pull common models
    models=("llama3.2" "llama3.1" "qwen2.5" "nomic-embed-text")
    
    for model in "${models[@]}"; do
        print_status "Pulling $model..."
        docker exec ollama ollama pull $model || print_warning "Failed to pull $model"
    done
    
    print_success "Ollama models initialized"
}

# Function to run health checks
health_check() {
    print_status "Running health checks..."
    
    services=("milvus" "weaviate" "postgres-vector" "redis" "ollama")
    
    for service in "${services[@]}"; do
        if $COMPOSE_CMD ps $service | grep -q "Up (healthy)"; then
            print_success "$service: Healthy"
        else
            print_warning "$service: Not healthy"
        fi
    done
}

# Function to backup data
backup_data() {
    local backup_dir="backups/$(date +%Y%m%d_%H%M%S)"
    print_status "Creating backup in $backup_dir..."
    
    mkdir -p "$backup_dir"
    
    # Backup volumes
    sudo cp -r data/volumes "$backup_dir/"
    
    # Backup configurations
    cp .env "$backup_dir/" 2>/dev/null || true
    cp docker-compose.yml "$backup_dir/"
    
    print_success "Backup created in $backup_dir"
}

# Function to show help
show_help() {
    echo "SwarmV2 RAG Infrastructure Management"
    echo ""
    echo "Usage: $0 [command] [options]"
    echo ""
    echo "Commands:"
    echo "  start [service]     Start all services or specific service"
    echo "  stop [service]      Stop all services or specific service"
    echo "  restart [service]   Restart all services or specific service"
    echo "  status              Show service status"
    echo "  logs [service]      Show logs for all or specific service"
    echo "  urls                Show access URLs"
    echo "  health              Run health checks"
    echo "  init-ollama         Initialize Ollama with common models"
    echo "  backup              Backup all data"
    echo "  cleanup             Remove all containers and data"
    echo "  help                Show this help message"
    echo ""
    echo "Available services:"
    echo "  milvus, weaviate, postgres-vector, redis, ollama"
    echo "  attu, weaviate-console, pgadmin, redis-commander, ollama-webui"
    echo ""
    echo "Examples:"
    echo "  $0 start                    # Start all services"
    echo "  $0 start milvus             # Start only Milvus"
    echo "  $0 logs ollama              # Show Ollama logs"
    echo "  $0 stop weaviate            # Stop Weaviate"
}

# Main script logic
main() {
    check_docker
    check_compose
    
    case "${1:-help}" in
        "start")
            create_directories
            setup_env
            if [ -n "$2" ]; then
                start_service "$2"
            else
                start_all
            fi
            ;;
        "stop")
            if [ -n "$2" ]; then
                stop_service "$2"
            else
                stop_all
            fi
            ;;
        "restart")
            if [ -n "$2" ]; then
                stop_service "$2"
                start_service "$2"
            else
                stop_all
                start_all
            fi
            ;;
        "status")
            show_status
            ;;
        "logs")
            show_logs "$2"
            ;;
        "urls")
            show_urls
            ;;
        "health")
            health_check
            ;;
        "init-ollama")
            init_ollama
            ;;
        "backup")
            backup_data
            ;;
        "cleanup")
            cleanup
            ;;
        "help"|"--help"|"-h")
            show_help
            ;;
        *)
            print_error "Unknown command: $1"
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
