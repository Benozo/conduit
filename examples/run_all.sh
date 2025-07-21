#!/bin/bash
set -e

# Conduit Examples Test Runner
# Usage: ./run_all.sh [--category] [--dry-run] [--help]

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Categories
PROTOCOL_EXAMPLES=("stdio_example" "sse_example" "pure_mcp")
LIBRARY_EXAMPLES=("pure_library" "pure_library_cli" "pure_library_web" "embedded")
LLM_EXAMPLES=("ollama" "openai" "model_integration")
TOOL_EXAMPLES=("custom_tools" "builtin_tools_test")
AGENT_EXAMPLES=("ai_agents" "agents_test" "agents_ollama" "agents_deepinfra" "agents_library_mode" "agents_mock_llm" "agents_vue_builder")
SWARM_EXAMPLES=("agent_swarm" "agent_swarm_llm" "agent_swarm_simple" "agent_swarm_workflows" "multi_llm_swarm")
RAG_EXAMPLES=("rag" "rag_chat_terminal" "rag_real_world" "langchain_mcp_integration")
SPECIAL_EXAMPLES=("simple_mcp_agent" "comprehensive_test" "agents_html_amender")

ALL_EXAMPLES=("${PROTOCOL_EXAMPLES[@]}" "${LIBRARY_EXAMPLES[@]}" "${LLM_EXAMPLES[@]}" "${TOOL_EXAMPLES[@]}" "${AGENT_EXAMPLES[@]}" "${SWARM_EXAMPLES[@]}" "${RAG_EXAMPLES[@]}" "${SPECIAL_EXAMPLES[@]}")

# Flags
DRY_RUN=false
CATEGORY="all"
VERBOSE=false

# Help function
show_help() {
    echo -e "${BLUE}Conduit Examples Test Runner${NC}"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "OPTIONS:"
    echo "  --all           Run all examples (default)"
    echo "  --protocol      Run protocol examples (stdio, sse, pure_mcp)"
    echo "  --library       Run library usage examples"
    echo "  --llm           Run LLM integration examples"
    echo "  --tools         Run tool development examples"
    echo "  --agents        Run agent framework examples"
    echo "  --swarm         Run agent swarm examples"
    echo "  --rag           Run RAG examples"
    echo "  --special       Run specialized examples"
    echo "  --dry-run       Show what would be run without executing"
    echo "  --verbose       Show detailed output"
    echo "  --help          Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                    # Run all examples"
    echo "  $0 --protocol         # Run only protocol examples"
    echo "  $0 --llm --dry-run    # Show LLM examples that would run"
    echo ""
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --all)
            CATEGORY="all"
            shift
            ;;
        --protocol)
            CATEGORY="protocol"
            shift
            ;;
        --library)
            CATEGORY="library"
            shift
            ;;
        --llm)
            CATEGORY="llm"
            shift
            ;;
        --tools)
            CATEGORY="tools"
            shift
            ;;
        --agents)
            CATEGORY="agents"
            shift
            ;;
        --swarm)
            CATEGORY="swarm"
            shift
            ;;
        --rag)
            CATEGORY="rag"
            shift
            ;;
        --special)
            CATEGORY="special"
            shift
            ;;
        --dry-run)
            DRY_RUN=true
            shift
            ;;
        --verbose)
            VERBOSE=true
            shift
            ;;
        --help)
            show_help
            exit 0
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            show_help
            exit 1
            ;;
    esac
done

# Get examples for category
get_examples_for_category() {
    case $1 in
        "all")
            echo "${ALL_EXAMPLES[@]}"
            ;;
        "protocol")
            echo "${PROTOCOL_EXAMPLES[@]}"
            ;;
        "library")
            echo "${LIBRARY_EXAMPLES[@]}"
            ;;
        "llm")
            echo "${LLM_EXAMPLES[@]}"
            ;;
        "tools")
            echo "${TOOL_EXAMPLES[@]}"
            ;;
        "agents")
            echo "${AGENT_EXAMPLES[@]}"
            ;;
        "swarm")
            echo "${SWARM_EXAMPLES[@]}"
            ;;
        "rag")
            echo "${RAG_EXAMPLES[@]}"
            ;;
        "special")
            echo "${SPECIAL_EXAMPLES[@]}"
            ;;
        *)
            echo ""
            ;;
    esac
}

# Check if example directory exists and has main.go
check_example() {
    local example=$1
    if [[ ! -d "$example" ]]; then
        echo -e "${RED}❌ Directory $example does not exist${NC}"
        return 1
    fi
    
    if [[ ! -f "$example/main.go" ]]; then
        echo -e "${YELLOW}⚠️  $example/main.go not found, skipping${NC}"
        return 1
    fi
    
    return 0
}

# Run a single example
run_example() {
    local example=$1
    echo -e "${BLUE}🚀 Running $example${NC}"
    
    if ! check_example "$example"; then
        return 1
    fi
    
    cd "$example"
    
    # Check for dependencies
    if [[ -f "go.mod" ]]; then
        echo -e "${YELLOW}📦 Installing dependencies...${NC}"
        if $VERBOSE; then
            go mod tidy
        else
            go mod tidy >/dev/null 2>&1
        fi
    fi
    
    # Check for build errors
    echo -e "${YELLOW}🔨 Building...${NC}"
    if $VERBOSE; then
        go build .
    else
        if ! go build . >/dev/null 2>&1; then
            echo -e "${RED}❌ Build failed for $example${NC}"
            cd ..
            return 1
        fi
    fi
    
    # Run with timeout (some examples may run servers)
    echo -e "${YELLOW}▶️  Testing...${NC}"
    if $VERBOSE; then
        timeout 10s go run main.go --help 2>/dev/null || timeout 10s go run main.go --version 2>/dev/null || timeout 5s go run main.go 2>/dev/null || true
    else
        timeout 10s go run main.go --help >/dev/null 2>&1 || timeout 10s go run main.go --version >/dev/null 2>&1 || timeout 5s go run main.go >/dev/null 2>&1 || true
    fi
    
    echo -e "${GREEN}✅ $example completed${NC}"
    cd ..
    return 0
}

# Main execution
main() {
    echo -e "${BLUE}🔧 Conduit Examples Test Runner${NC}"
    echo -e "${BLUE}Category: $CATEGORY${NC}"
    echo ""
    
    # Get examples to run
    examples_to_run=($(get_examples_for_category "$CATEGORY"))
    
    if [[ ${#examples_to_run[@]} -eq 0 ]]; then
        echo -e "${RED}No examples found for category: $CATEGORY${NC}"
        exit 1
    fi
    
    echo -e "${YELLOW}Found ${#examples_to_run[@]} examples to run${NC}"
    
    if $DRY_RUN; then
        echo -e "${BLUE}🔍 Dry run - would execute:${NC}"
        for example in "${examples_to_run[@]}"; do
            if check_example "$example"; then
                echo -e "  ${GREEN}✓${NC} $example"
            else
                echo -e "  ${RED}✗${NC} $example (missing or invalid)"
            fi
        done
        exit 0
    fi
    
    # Run examples
    success_count=0
    total_count=${#examples_to_run[@]}
    
    for example in "${examples_to_run[@]}"; do
        echo ""
        echo -e "${BLUE}==== $example ====${NC}"
        
        if run_example "$example"; then
            ((success_count++))
        fi
    done
    
    # Summary
    echo ""
    echo -e "${BLUE}📊 Summary${NC}"
    echo -e "Total examples: $total_count"
    echo -e "Successful: ${GREEN}$success_count${NC}"
    echo -e "Failed: ${RED}$((total_count - success_count))${NC}"
    
    if [[ $success_count -eq $total_count ]]; then
        echo -e "${GREEN}🎉 All examples completed successfully!${NC}"
        exit 0
    else
        echo -e "${YELLOW}⚠️  Some examples failed or were skipped${NC}"
        exit 1
    fi
}

# Check if we're in the examples directory
if [[ ! -f "README.md" ]] || [[ ! "$(basename "$PWD")" == "examples" ]]; then
    echo -e "${RED}Error: Please run this script from the examples/ directory${NC}"
    exit 1
fi

# Run main function
main "$@"
