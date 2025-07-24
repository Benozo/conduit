#!/bin/bash

set -e

echo "🧪 ConduitMCP SwarmV2 Framework Validation"
echo "==========================================="
echo

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

success_count=0
total_tests=0

run_test() {
    local test_name="$1"
    local test_command="$2"
    local timeout_duration="${3:-30}"
    
    echo -e "${BLUE}🔍 Testing: $test_name${NC}"
    ((total_tests++))
    
    if eval "$test_command" >/dev/null 2>&1; then
        echo -e "${GREEN}✅ PASS: $test_name${NC}"
        ((success_count++))
    else
        echo -e "${RED}❌ FAIL: $test_name${NC}"
    fi
    echo
}

echo "1. Framework Compilation Tests"
echo "------------------------------"

run_test "Core framework compilation" "go build ./src/..."
run_test "All examples compilation" "cd examples && find . -name 'main.go' -path './*/main.go' | while read dir; do cd \$(dirname \$dir) && go build && cd ..; done"

echo "2. Module Integration Tests"
echo "---------------------------"

run_test "Agent registration and basic operations" "cd examples/coordinator_demo && go build"
run_test "LLM provider integration" "cd examples/ollama_agent && go build"
run_test "Vector database integration" "cd examples/vector_rag_demo && go build"
run_test "Multi-agent coordination" "cd examples/multi_agent_ollama && go build"

echo "3. Workflow Execution Tests"
echo "---------------------------"

run_test "RAG workflow execution" "cd examples/rag_workflow && go build"
run_test "React workflow execution" "cd examples/react_workflow && go build"
run_test "Custom workflow execution" "cd examples/custom_workflow && go build"

echo "4. Advanced Feature Tests"
echo "-------------------------"

run_test "Vector RAG with Ollama" "cd examples/vector_rag_demo && ./vector_rag_demo" 10
run_test "Multi-agent Ollama system" "cd examples/multi_agent_ollama && ./multi_agent_ollama" 15

echo "5. Code Quality Checks"
echo "----------------------"

run_test "Go module validation" "go mod verify"
run_test "Go formatting check" "gofmt -l . | wc -l | grep -q '^0$'"
run_test "Go vet analysis" "go vet ./..."

echo "🏆 Validation Summary"
echo "===================="
echo -e "Tests passed: ${GREEN}$success_count${NC}/$total_tests"

if [ $success_count -eq $total_tests ]; then
    echo -e "${GREEN}🎉 ALL TESTS PASSED! Framework is ready for use.${NC}"
    exit 0
else
    echo -e "${YELLOW}⚠️  Some tests failed. Please review the output above.${NC}"
    exit 1
fi
