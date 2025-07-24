# React Workflow Example with Ollama Integration

This example demonstrates the React (Reason, Act, Observe) pattern using the swarm framework, enhanced with Ollama for AI-powered reasoning and decision-making.

## What it demonstrates

- **AI-Powered Reasoner**: Uses Ollama to analyze complex situations and make intelligent decisions
- **AI-Guided Actor**: Generates detailed action plans using AI reasoning
- **AI-Enhanced Observer**: Provides intelligent monitoring and assessment using AI analysis
- **React Workflow**: Orchestrates the reason-act-observe loop with real AI capabilities
- **Multiple Scenarios**: Demonstrates React cycles across different problem domains
- **Fallback Handling**: Gracefully handles Ollama connection issues

## Features

- **Real-Time AI Analysis**: Uses Ollama's llama3.2 model for situation analysis
- **Structured Decision Making**: AI provides analysis, root causes, and recommended actions
- **Action Planning**: Generates specific, actionable plans with timelines and metrics
- **Intelligent Monitoring**: AI-powered observation with progress assessment
- **Multi-Scenario Testing**: Three different scenarios to showcase versatility

## Prerequisites

- Go 1.19 or higher
- **Ollama server running** with llama3.2 model available
- Access to http://192.168.10.10:11434 (or modify the host in main.go)

## How to run

```bash
# Make sure Ollama is running with llama3.2 model
# ollama serve  # (if not already running)
# ollama pull llama3.2  # (if model not already downloaded)

cd /home/engineone/Projects/AI/ConduitMCP/swarmV2/examples/react_workflow
go run main.go
```

## Configuration

You can modify the Ollama configuration in `main.go`:

```go
ollamaHost := "http://192.168.10.10:11434"  // Change to your Ollama server
ollamaModel := "llama3.2"                    // Change to your preferred model
```

## Expected output

The example will demonstrate three React cycles for different scenarios:

1. **Performance Issue**: AI analyzes slow application performance
2. **Anomaly Detection**: AI investigates unusual user behavior patterns  
3. **Feature Integration**: AI plans external API integration

For each scenario, you'll see:
- 🧠 **AI Reasoning Phase**: Detailed situation analysis with root causes and recommendations
- ⚡ **Action Phase**: Specific action plans with timelines and success metrics
- 👁️ **Observation Phase**: Monitoring assessment with progress indicators

The example will:
1. Initialize React agents (Reasoner, Actor, Observer)
2. Create a coordinator agent
3. Set up a React workflow
4. Execute the complete React loop
5. Display the workflow completion status

## Sample Output

```
=== React Loop Workflow Demo with Ollama ===
🔍 Setting up Ollama integration...
🔍 Testing connection to Ollama at http://192.168.10.10:11434...
✅ Successfully connected to Ollama!
🤖 Using model: llama3.2

🚀 Starting AI-enhanced React loop workflow...
Coordinator: ReactCoordinator
Reasoner: AI-Reasoner (AI-powered)
Actor: AI-Actor (AI-guided)
Observer: AI-Observer (AI-monitoring)

🧪 Scenario 1: A user reports slow application performance during peak hours
------------------------------------------------------------
🧠 AI Reasoning Phase...
✅ AI Analysis completed
📋 Decision: **Analysis of the Situation**
The user has reported slow application performance during peak hours...

⚡ Action Phase...
✅ Action plan generated
🎯 Plan: **Action Plan: Addressing Slow Application Performance**
Immediate Steps (Week 1-2):
1. Conduct a Thorough Performance Analysis...

👁️  Observation Phase...
✅ Monitoring assessment completed
📊 Observation: **Monitoring Report**
Implementation Progress:
* Immediate Steps (Week 1-2): Conducted thorough performance analysis...
🔄 React cycle completed for this scenario
```

## Architecture

The React workflow follows this pattern:

1. **Reasoning**: AI analyzes the situation using Ollama
2. **Acting**: AI generates specific action plans  
3. **Observing**: AI monitors and assesses implementation

Each phase uses structured prompts to ensure comprehensive analysis and actionable outputs.

## Error Handling

- **Connection Failures**: Gracefully falls back to simulated responses if Ollama is unavailable
- **AI Errors**: Provides fallback decisions and actions if AI generation fails
- **Timeout Protection**: Uses context timeouts to prevent hanging requests

## Related Examples

- `vector_rag_demo/`: For vector database-enhanced RAG workflows
- `multi_agent_ollama/`: For multi-agent coordination with Ollama
- `ollama_agent/`: For basic Ollama integration patterns
