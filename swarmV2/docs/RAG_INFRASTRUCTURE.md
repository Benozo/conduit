# RAG Infrastructure Guide

This guide provides comprehensive setup and usage instructions for the SwarmV2 RAG (Retrieve-Augment-Generate) infrastructure using Docker Compose.

## 🚀 Quick Start

1. **Clone and navigate to the project**:
   ```bash
   cd swarmV2
   ```

2. **Start all RAG services**:
   ```bash
   ./rag-infrastructure.sh start
   ```

3. **View access URLs**:
   ```bash
   ./rag-infrastructure.sh urls
   ```

## 📋 What's Included

### Vector Databases
- **Milvus**: High-performance vector database with GPU acceleration support
- **Weaviate**: Modern vector database with built-in vectorization modules
- **PostgreSQL + pgvector**: Traditional SQL database with vector extensions

### Management UIs
- **Attu** (Milvus UI): http://localhost:3000
- **Weaviate Console**: http://localhost:8081
- **pgAdmin** (PostgreSQL): http://localhost:5050

### AI & LLM Services
- **Ollama**: Local LLM inference server
- **Open WebUI**: Modern web interface for Ollama

### Supporting Services
- **Redis**: Caching and session management
- **MinIO**: S3-compatible object storage for Milvus
- **etcd**: Distributed configuration store for Milvus

## 🛠️ Management Commands

### Basic Operations
```bash
# Start all services
./rag-infrastructure.sh start

# Stop all services
./rag-infrastructure.sh stop

# Check service status
./rag-infrastructure.sh status

# Show access URLs
./rag-infrastructure.sh urls
```

### Service-Specific Operations
```bash
# Start only Milvus
./rag-infrastructure.sh start milvus

# Stop only Weaviate
./rag-infrastructure.sh stop weaviate

# View Ollama logs
./rag-infrastructure.sh logs ollama
```

### Maintenance
```bash
# Run health checks
./rag-infrastructure.sh health

# Initialize Ollama with common models
./rag-infrastructure.sh init-ollama

# Backup all data
./rag-infrastructure.sh backup

# Clean up everything (DESTRUCTIVE)
./rag-infrastructure.sh cleanup
```

## 🔧 Configuration

### Environment Setup
1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

2. Edit `.env` with your specific configuration:
   ```bash
   nano .env
   ```

### Key Configuration Options
- **DOCKER_VOLUME_DIRECTORY**: Where to store persistent data
- **LLM API Keys**: OpenAI, Anthropic, Cohere, etc.
- **Database Credentials**: Custom passwords and users
- **Service Ports**: Modify if you have conflicts

## 📊 Service Access Information

| Service | URL | Credentials | Purpose |
|---------|-----|-------------|---------|
| Milvus API | http://localhost:19530 | - | Vector operations |
| Attu (Milvus UI) | http://localhost:3000 | - | Milvus management |
| Weaviate API | http://localhost:8080 | - | Vector operations |
| Weaviate Console | http://localhost:8081 | - | Weaviate management |
| PostgreSQL | localhost:5432 | postgres/password | SQL + Vector |
| pgAdmin | http://localhost:5050 | admin@swarmv2.com/admin123 | PostgreSQL UI |
| Redis | localhost:6379 | - | Caching |
| Redis Commander | http://localhost:8082 | - | Redis UI |
| Ollama API | http://localhost:11434 | - | Local LLM |
| Ollama WebUI | http://localhost:8083 | - | Ollama interface |
| MinIO Console | http://localhost:9001 | minioadmin/minioadmin | Object storage |

## 🚀 Using with SwarmV2

### Example: Milvus Integration
```go
package main

import (
    "context"
    "log"
    
    "github.com/benozo/neuron/src/vectordb/providers"
    "github.com/benozo/neuron/src/vectordb"
)

func main() {
    ctx := context.Background()
    
    // Connect to Milvus
    milvusDB := providers.NewMilvusProvider("localhost:19530")
    err := milvusDB.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    // Create RAG store
    ragStore := vectordb.NewRAGStore(milvusDB)
    
    // Add documents
    documents := []vectordb.Document{
        {
            ID: "doc1",
            Content: "SwarmV2 is a powerful multi-agent framework...",
            Type: vectordb.DocumentTypeText,
        },
    }
    
    err = ragStore.AddDocuments(ctx, "knowledge_base", documents)
    if err != nil {
        log.Fatal(err)
    }
    
    // Search similar documents
    results, err := ragStore.SearchDocuments(ctx, "knowledge_base", "multi-agent systems", vectordb.SearchOptions{
        Limit: 5,
        MinSimilarity: 0.7,
    })
    if err != nil {
        log.Fatal(err)
    }
    
    for _, result := range results {
        log.Printf("Found: %s (similarity: %.3f)", result.Document.Content, result.Similarity)
    }
}
```

### Example: Weaviate Integration
```go
package main

import (
    "context"
    "log"
    
    "github.com/benozo/neuron/src/vectordb/providers"
)

func main() {
    ctx := context.Background()
    
    // Connect to Weaviate
    weaviateDB := providers.NewWeaviateProvider("http://localhost:8080", "")
    err := weaviateDB.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    // Create collection
    err = weaviateDB.CreateCollection(ctx, "Documents", map[string]interface{}{
        "vectorizer": "text2vec-openai",
        "properties": map[string]interface{}{
            "content": map[string]interface{}{
                "dataType": []string{"text"},
            },
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println("Weaviate collection created successfully")
}
```

### Example: PostgreSQL + pgvector
```go
package main

import (
    "context"
    "log"
    
    "github.com/benozo/neuron/src/vectordb/providers"
)

func main() {
    ctx := context.Background()
    
    // Connect to PostgreSQL with pgvector
    pgDB := providers.NewPGVectorProvider(
        "host=localhost port=5432 user=postgres password=password dbname=vectordb sslmode=disable",
    )
    err := pgDB.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println("Connected to PostgreSQL with pgvector")
}
```

## 🔍 Monitoring and Debugging

### Check Service Health
```bash
# Overall status
./rag-infrastructure.sh status

# Health checks
./rag-infrastructure.sh health

# View logs
./rag-infrastructure.sh logs [service-name]
```

### Common Issues

#### Port Conflicts
If you get port binding errors:
1. Check what's using the ports: `netstat -tulpn | grep [port]`
2. Modify ports in `docker-compose.yml`
3. Update your application configuration

#### Memory Issues
For large datasets:
1. Increase Docker memory allocation
2. Tune vector database settings
3. Consider using external vector databases

#### Storage Issues
Monitor disk usage:
```bash
# Check volume sizes
docker system df

# Clean up unused resources
docker system prune -a --volumes
```

## 📈 Performance Optimization

### Milvus Optimization
- Use GPU acceleration when available
- Tune index parameters for your use case
- Configure memory limits appropriately

### Weaviate Optimization
- Choose appropriate vectorizer modules
- Configure backup and replication
- Tune HNSW parameters

### PostgreSQL Optimization
- Adjust `shared_buffers` and `work_mem`
- Configure appropriate `maintenance_work_mem`
- Tune vector index parameters

## 🔐 Security Considerations

### Production Deployment
1. **Change default passwords**
2. **Enable authentication** for all services
3. **Use TLS/SSL** for network communication
4. **Implement proper access controls**
5. **Regular security updates**

### Network Security
- Use Docker networks for service isolation
- Implement proper firewall rules
- Consider VPN for remote access

## 📋 Backup and Recovery

### Automated Backups
```bash
# Create backup
./rag-infrastructure.sh backup

# Backups are stored in: backups/YYYYMMDD_HHMMSS/
```

### Manual Backups
```bash
# Backup specific service data
docker compose exec postgres-vector pg_dump -U postgres vectordb > backup.sql

# Backup Milvus data
docker compose exec milvus tar -czf /tmp/milvus-backup.tar.gz /var/lib/milvus
```

## 🚀 Scaling and Production

### Horizontal Scaling
- Use Kubernetes for container orchestration
- Implement load balancing for APIs
- Consider managed vector database services

### Monitoring
- Use Prometheus + Grafana for metrics
- Implement health checks and alerting
- Monitor resource usage and performance

### High Availability
- Configure database replication
- Implement backup strategies
- Use multiple availability zones

## 📞 Support and Troubleshooting

### Getting Help
1. Check service logs: `./rag-infrastructure.sh logs [service]`
2. Verify service health: `./rag-infrastructure.sh health`
3. Review configuration: Check `.env` and `docker-compose.yml`
4. Community support: GitHub Issues and Discussions

### Common Solutions
- **Service won't start**: Check logs and port conflicts
- **Connection refused**: Verify service is running and accessible
- **Performance issues**: Check resource allocation and tuning
- **Data corruption**: Restore from backup and check disk space
