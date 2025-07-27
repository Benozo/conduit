#!/bin/bash

# SwarmV2 RAG Infrastructure Demo
# This script demonstrates the RAG infrastructure setup and basic usage

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

print_header() {
    echo -e "${PURPLE}╔════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${PURPLE}║$1${NC}"
    echo -e "${PURPLE}╚════════════════════════════════════════════════════════════════╝${NC}"
}

print_step() {
    echo -e "${CYAN}🔹 $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Check if infrastructure is running
check_infrastructure() {
    if ! docker compose ps | grep -q "Up"; then
        print_warning "RAG infrastructure doesn't seem to be running"
        echo "Start it with: ./rag-infrastructure.sh start"
        return 1
    fi
    return 0
}

# Function to wait for service to be ready
wait_for_service() {
    local service_name=$1
    local url=$2
    local max_attempts=${3:-30}
    local attempt=1
    
    print_step "Waiting for $service_name to be ready..."
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s "$url" > /dev/null 2>&1; then
            print_success "$service_name is ready!"
            return 0
        fi
        
        echo -n "."
        sleep 2
        ((attempt++))
    done
    
    print_warning "$service_name not ready after $max_attempts attempts"
    return 1
}

# Function to demonstrate Milvus
demo_milvus() {
    print_header "🗄️  MILVUS VECTOR DATABASE DEMO                                    "
    
    print_step "Checking Milvus API..."
    if curl -s "http://localhost:19530/health" > /dev/null 2>&1; then
        print_success "Milvus API is accessible"
    else
        print_warning "Milvus API not accessible"
        return 1
    fi
    
    print_info "Milvus UI (Attu) available at: http://localhost:3000"
    print_info "API endpoint: http://localhost:19530"
    print_info "Use Attu to:"
    echo "  • Create collections"
    echo "  • Insert vectors"
    echo "  • Perform similarity searches"
    echo "  • Monitor performance"
    echo ""
}

# Function to demonstrate Weaviate
demo_weaviate() {
    print_header "🧠 WEAVIATE VECTOR DATABASE DEMO                                   "
    
    print_step "Checking Weaviate API..."
    if curl -s "http://localhost:8080/v1/.well-known/ready" > /dev/null 2>&1; then
        print_success "Weaviate API is accessible"
        
        # Get Weaviate version info
        local version=$(curl -s "http://localhost:8080/v1/meta" | grep -o '"version":"[^"]*"' | cut -d'"' -f4 2>/dev/null || echo "unknown")
        print_info "Weaviate version: $version"
    else
        print_warning "Weaviate API not accessible"
        return 1
    fi
    
    print_info "Weaviate Console available at: http://localhost:8081"
    print_info "API endpoint: http://localhost:8080"
    print_info "Use Weaviate Console to:"
    echo "  • Create schemas"
    echo "  • Import data"
    echo "  • Execute GraphQL queries"
    echo "  • Configure vectorizers"
    echo ""
}

# Function to demonstrate PostgreSQL + pgvector
demo_pgvector() {
    print_header "🐘 POSTGRESQL + PGVECTOR DEMO                                     "
    
    print_step "Checking PostgreSQL connection..."
    if docker compose exec -T postgres-vector pg_isready -U postgres -d vectordb > /dev/null 2>&1; then
        print_success "PostgreSQL is ready"
        
        # Check pgvector extension
        local version=$(docker compose exec -T postgres-vector psql -U postgres -d vectordb -c "SELECT extversion FROM pg_extension WHERE extname = 'vector';" -t 2>/dev/null | tr -d ' \n' || echo "unknown")
        print_info "pgvector extension version: $version"
        
        # Show table structure
        print_step "Database structure:"
        docker compose exec -T postgres-vector psql -U postgres -d vectordb -c "\\dt" 2>/dev/null || true
    else
        print_warning "PostgreSQL not accessible"
        return 1
    fi
    
    print_info "pgAdmin available at: http://localhost:5050"
    print_info "Connection: localhost:5432, postgres/password, vectordb"
    print_info "Use pgAdmin to:"
    echo "  • Manage vector tables"
    echo "  • Execute vector similarity queries"
    echo "  • Monitor performance"
    echo "  • Backup/restore data"
    echo ""
}

# Function to demonstrate Ollama
demo_ollama() {
    print_header "🤖 OLLAMA LOCAL LLM DEMO                                          "
    
    print_step "Checking Ollama API..."
    if curl -s "http://localhost:11434/api/version" > /dev/null 2>&1; then
        print_success "Ollama API is accessible"
        
        # Get version
        local version=$(curl -s "http://localhost:11434/api/version" | grep -o '"version":"[^"]*"' | cut -d'"' -f4 2>/dev/null || echo "unknown")
        print_info "Ollama version: $version"
        
        # List models
        print_step "Available models:"
        curl -s "http://localhost:11434/api/tags" | grep -o '"name":"[^"]*"' | cut -d'"' -f4 2>/dev/null || echo "  No models found"
    else
        print_warning "Ollama API not accessible"
        return 1
    fi
    
    print_info "Ollama WebUI available at: http://localhost:8083"
    print_info "API endpoint: http://localhost:11434"
    print_info "Use Ollama WebUI to:"
    echo "  • Chat with models"
    echo "  • Download new models"
    echo "  • Manage model settings"
    echo "  • View model information"
    echo ""
}

# Function to demonstrate Redis
demo_redis() {
    print_header "⚡ REDIS CACHING DEMO                                             "
    
    print_step "Checking Redis connection..."
    if docker compose exec -T redis redis-cli ping > /dev/null 2>&1; then
        print_success "Redis is accessible"
        
        # Get Redis info
        local version=$(docker compose exec -T redis redis-cli info server | grep "redis_version" | cut -d':' -f2 | tr -d '\r' 2>/dev/null || echo "unknown")
        print_info "Redis version: $version"
        
        # Set and get a test value
        docker compose exec -T redis redis-cli set "swarmv2:demo" "Hello from SwarmV2!" > /dev/null 2>&1
        local test_value=$(docker compose exec -T redis redis-cli get "swarmv2:demo" 2>/dev/null | tr -d '\r' || echo "")
        if [ "$test_value" = "Hello from SwarmV2!" ]; then
            print_success "Redis test operation successful"
        fi
    else
        print_warning "Redis not accessible"
        return 1
    fi
    
    print_info "Redis Commander available at: http://localhost:8082"
    print_info "Connection: localhost:6379"
    print_info "Use Redis Commander to:"
    echo "  • Browse keys and values"
    echo "  • Monitor memory usage"
    echo "  • Execute Redis commands"
    echo "  • View connection statistics"
    echo ""
}

# Function to show integration example
show_integration_example() {
    print_header "🔗 SWARMV2 INTEGRATION EXAMPLE                                   "
    
    cat << 'EOF'
Here's how to use the RAG infrastructure in your SwarmV2 applications:

```go
package main

import (
    "context"
    "log"
    
    "github.com/benozo/neuron/src/vectordb/providers"
    "github.com/benozo/neuron/src/vectordb"
    "github.com/benozo/neuron/src/llm/providers"
    "github.com/benozo/neuron/src/workflows"
)

func main() {
    ctx := context.Background()
    
    // Choose your vector database
    
    // Option 1: Milvus
    milvusDB := providers.NewMilvusProvider("localhost:19530")
    milvusDB.Connect(ctx)
    
    // Option 2: Weaviate  
    weaviateDB := providers.NewWeaviateProvider("http://localhost:8080", "")
    weaviateDB.Connect(ctx)
    
    // Option 3: PostgreSQL + pgvector
    pgDB := providers.NewPGVectorProvider(
        "host=localhost port=5432 user=postgres password=password dbname=vectordb sslmode=disable",
    )
    pgDB.Connect(ctx)
    
    // Create RAG store (use any of the above)
    ragStore := vectordb.NewRAGStore(milvusDB) // or weaviateDB, pgDB
    
    // Setup LLM provider
    ollama := providers.NewOllamaProvider("http://localhost:11434", "llama3.2")
    
    // Create RAG workflow
    ragWorkflow := workflows.NewRAGWorkflow(ragStore, ollama)
    
    // Add documents
    documents := []vectordb.Document{
        {
            ID: "doc1",
            Content: "SwarmV2 is a powerful multi-agent framework for AI applications...",
            Type: vectordb.DocumentTypeText,
        },
    }
    
    ragStore.AddDocuments(ctx, "knowledge_base", documents)
    
    // Execute RAG query
    result, err := ragWorkflow.Execute(ctx, "What is SwarmV2?")
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Generated response: %s", result)
}
```

📚 For more examples, see:
  • examples/vector_rag_demo/
  • examples/cloudflare_ai/
  • docs/RAG_INFRASTRUCTURE.md
EOF
    echo ""
}

# Function to show next steps
show_next_steps() {
    print_header "🚀 NEXT STEPS                                                     "
    
    print_info "Explore the services:"
    echo "  1. Open Milvus UI (Attu): http://localhost:3000"
    echo "  2. Open Weaviate Console: http://localhost:8081"
    echo "  3. Open pgAdmin: http://localhost:5050"
    echo "  4. Open Ollama WebUI: http://localhost:8083"
    echo "  5. Open Redis Commander: http://localhost:8082"
    echo ""
    
    print_info "Try the examples:"
    echo "  • cd examples/vector_rag_demo && go run main.go"
    echo "  • cd examples/cloudflare_ai && go run main.go"
    echo "  • cd examples/multi_agent_ollama && go run main.go"
    echo ""
    
    print_info "Management commands:"
    echo "  • ./rag-infrastructure.sh status    # Check service status"
    echo "  • ./rag-infrastructure.sh logs      # View logs"
    echo "  • ./rag-infrastructure.sh health    # Run health checks"
    echo "  • ./rag-infrastructure.sh stop      # Stop all services"
    echo ""
    
    print_info "Learn more:"
    echo "  • Read docs/RAG_INFRASTRUCTURE.md"
    echo "  • Explore examples/ directory"
    echo "  • Check GitHub repository for updates"
    echo ""
}

# Main demo function
main() {
    clear
    print_header "🌟 SWARMV2 RAG INFRASTRUCTURE DEMO                               "
    echo ""
    print_info "This demo showcases the comprehensive RAG infrastructure setup"
    print_info "including vector databases, AI services, and management UIs."
    echo ""
    
    # Check if infrastructure is running
    if ! check_infrastructure; then
        echo ""
        print_info "To start the infrastructure, run:"
        echo "  ./rag-infrastructure.sh start"
        echo ""
        exit 1
    fi
    
    # Run demos
    demo_milvus
    demo_weaviate  
    demo_pgvector
    demo_ollama
    demo_redis
    
    show_integration_example
    show_next_steps
    
    print_header "🎉 DEMO COMPLETED SUCCESSFULLY!                                   "
    print_success "Your SwarmV2 RAG infrastructure is ready for AI applications!"
}

# Run main function
main "$@"
