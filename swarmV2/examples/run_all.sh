#!/bin/bash

# Run all SwarmV2 examples
echo "Running all SwarmV2 examples..."
echo "================================"

echo
echo "1. Running Coordinator Demo..."
echo "----------------------------"
cd coordinator_demo && go run main.go
echo

echo "2. Running RAG Workflow..."
echo "-------------------------"
cd ../rag_workflow && go run main.go
echo

echo "3. Running React Workflow..."
echo "---------------------------"
cd ../react_workflow && go run main.go
echo

echo "4. Running Custom Workflow..."
echo "----------------------------"
cd ../custom_workflow && go run main.go
echo

echo "5. Running Ollama Agent..."
echo "-------------------------"
cd ../ollama_agent && go run main.go
echo

echo "================================"
echo "All examples completed successfully!"
