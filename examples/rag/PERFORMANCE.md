# ConduitMCP RAG Performance Optimization Guide

## CPU/Memory Issues and Solutions

### ❌ Common Performance Problems

1. **High CPU during indexing** - Large documents with complex chunking
2. **Memory spikes** - Multiple concurrent embedding requests  
3. **Slow response times** - Large chunk sizes and search limits
4. **System overload** - No resource limits or timeouts

### ✅ Implemented Optimizations

#### 1. **Minimal Indexing Mode**
```bash
# Skip automatic indexing to prevent CPU spikes
SKIP_INDEXING=true go run main.go
```

#### 2. **Reduced Chunk Sizes**
- **Chunk size**: 256 characters (down from 1024+ default)
- **Overlap**: 20 characters (minimal)
- **Strategy**: Fixed size (most efficient)

#### 3. **Shorter Timeouts**
- **Embedding ping**: 10 seconds (down from 30)
- **Document indexing**: 5 seconds per document
- **Search operations**: 15 seconds max

#### 4. **Resource Management**
- **Processing delays**: 1 second between documents
- **Graceful failures**: Continue on individual document errors
- **Memory cleanup**: Proper context cancellation

#### 5. **Minimal Sample Documents**
- **Ultra-light content**: Single sentences instead of paragraphs
- **Fewer documents**: 2 instead of 5+ sample docs
- **Essential metadata only**: Reduced overhead

## Performance Configuration

### Environment Variables for Optimization

```bash
# Core performance settings
export SKIP_INDEXING=true              # Skip auto-indexing
export PORT=9091                       # Use different port
export RAG_PROVIDER=ollama             # Local embeddings

# Ollama optimization
export OLLAMA_HOST=localhost           # Local processing
export OLLAMA_MODEL=nomic-embed-text   # Lightweight model

# Database optimization  
export POSTGRES_HOST=localhost         # Local database
```

### Hardware Recommendations

#### Minimum Requirements
- **RAM**: 4GB+ available
- **CPU**: 2+ cores
- **Storage**: SSD recommended for vector operations

#### Optimal Configuration
- **RAM**: 8GB+ for larger knowledge bases
- **CPU**: 4+ cores for concurrent operations
- **GPU**: Optional for faster embedding generation

## Monitoring Performance

### 1. **System Metrics**
```bash
# Monitor CPU/memory usage
top -p $(pgrep -f "go run main.go")

# Watch memory consumption
watch -n 1 "ps aux | grep 'go run main.go'"
```

### 2. **API Response Times**
```bash
# Check response times
curl -w "@curl-format.txt" -X POST http://localhost:9091/chat \
  -H "Content-Type: application/json" \
  -d '{"message":"test"}'

# Create curl-format.txt:
echo "time_total: %{time_total}s\n" > curl-format.txt
```

### 3. **Database Performance**
```sql
-- Check database size and performance
SELECT 
  schemaname,
  tablename,
  attname,
  n_distinct,
  correlation
FROM pg_stats 
WHERE tablename LIKE '%document%' OR tablename LIKE '%chunk%';
```

## Production Deployment

### 1. **Docker Optimization**
```dockerfile
# Multi-stage build for smaller image
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -ldflags="-s -w" -o rag-server main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/rag-server .
CMD ["./rag-server"]
```

### 2. **Resource Limits**
```yaml
# Kubernetes deployment
resources:
  limits:
    memory: "2Gi"
    cpu: "1000m"
  requests:
    memory: "512Mi"
    cpu: "250m"
```

### 3. **Scaling Strategy**
```yaml
# Horizontal Pod Autoscaler
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: rag-server-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: rag-server
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

## Best Practices

### 1. **Document Management**
- **Batch uploads**: Process multiple small documents instead of large ones
- **Content preprocessing**: Clean and optimize text before indexing
- **Metadata optimization**: Use minimal but meaningful metadata
- **Incremental indexing**: Add documents gradually rather than bulk loading

### 2. **Search Optimization**
- **Limit results**: Use reasonable limits (3-5 results)
- **Cache frequent queries**: Implement query result caching
- **Async processing**: Use background workers for heavy operations
- **Connection pooling**: Optimize database connections

### 3. **Memory Management**
- **Garbage collection tuning**: Adjust `GOGC` environment variable
- **Connection limits**: Limit concurrent database connections
- **Timeout management**: Use appropriate context timeouts
- **Resource cleanup**: Ensure proper cleanup of resources

## Troubleshooting

### High CPU Usage
1. Check if indexing is running: `ps aux | grep embedding`
2. Reduce chunk size and overlap
3. Enable `SKIP_INDEXING=true`
4. Use lighter embedding model
5. Add processing delays

### High Memory Usage
1. Monitor with `htop` or `top`
2. Reduce concurrent operations
3. Lower database connection pool size
4. Implement connection timeouts
5. Use streaming for large responses

### Slow Response Times
1. Check database query performance
2. Optimize vector similarity search
3. Reduce search result limits
4. Use local embeddings (Ollama vs OpenAI API)
5. Implement response caching

### Database Issues
1. Check PostgreSQL performance: `SELECT * FROM pg_stat_activity;`
2. Monitor vector index usage
3. Optimize pgvector configuration
4. Consider index maintenance: `VACUUM ANALYZE;`

## Load Testing

### Simple Load Test
```bash
# Install apache bench
sudo apt-get install apache2-utils

# Test chat endpoint
ab -n 100 -c 10 -T application/json -p chat-payload.json \
   http://localhost:9091/chat

# Create chat-payload.json
echo '{"message":"What is ConduitMCP?"}' > chat-payload.json
```

### Advanced Load Testing
```bash
# Use wrk for more advanced testing
wrk -t12 -c400 -d30s -s chat-script.lua http://localhost:9091/chat
```

This guide ensures your ConduitMCP RAG system runs efficiently even under load!
