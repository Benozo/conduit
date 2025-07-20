# AI Agents with DeepInfra LLM Integration

This example demonstrates how to use ConduitMCP's AI agents with **real LLM integration** through DeepInfra's API. DeepInfra provides access to state-of-the-art language models through an OpenAI-compatible API.

## What This Example Shows

- 🧠 **Real LLM Reasoning**: Agents use actual language models for intelligent problem solving
- 🔗 **DeepInfra Integration**: OpenAI-compatible API with bearer token authentication
- 🚀 **Cloud-Powered Intelligence**: Access to powerful models without local infrastructure
- 🔧 **Tool Integration**: LLM-guided usage of MCP tools for complex tasks
- ⚡ **Production Ready**: Scalable cloud infrastructure for real applications

## Prerequisites

1. **DeepInfra Account**: Sign up at [deepinfra.com](https://deepinfra.com)
2. **API Token**: Get your bearer token from the DeepInfra dashboard
3. **Environment Variables**: Set required configuration

## Setup

### 1. Get DeepInfra API Token
```bash
# Visit https://deepinfra.com
# Sign up/login and get your API token
```

### 2. Set Environment Variables
```bash
# Required: Your DeepInfra API token
export DEEPINFRA_TOKEN="your_deepinfra_token_here"

# Optional: Specify model (defaults to Meta-Llama-3.1-8B-Instruct)
export DEEPINFRA_MODEL="meta-llama/Meta-Llama-3.1-8B-Instruct"
```

### 3. Run the Example
```bash
cd /path/to/ConduitMCP
go run examples/agents_deepinfra/main.go
```

## Available Models

DeepInfra supports many popular models:

- `meta-llama/Meta-Llama-3.1-8B-Instruct` (default)
- `meta-llama/Meta-Llama-3.1-70B-Instruct`
- `microsoft/WizardLM-2-8x22B`
- `mistralai/Mixtral-8x7B-Instruct-v0.1`
- `google/gemma-2-27b-it`
- And many more...

Set your preferred model:
```bash
export DEEPINFRA_MODEL="meta-llama/Meta-Llama-3.1-70B-Instruct"
```

## Example Demonstrations

### 1. Mathematical Problem Solving
```
Problem: A restaurant serves 8 tables per hour. If each table pays $25 on average, how much revenue per hour?

🧠 DeepInfra LLM is analyzing the problem...
   🧮 Tool executed: 8.0 × 25.0 = 200.0
✅ Problem solved in 2.3s!
```

The LLM analyzes the word problem and determines it needs to multiply tables per hour by average payment.

### 2. Text Analysis
```
Text: DeepInfra provides access to cutting-edge language models through a simple API...

🧠 DeepInfra LLM is analyzing the text...
✅ Text analyzed in 1.8s!
```

The LLM analyzes text content, counts words, and stores insights about the platform.

### 3. Complex Multi-Step Problems
```
Problem: A software team of 5 developers works 40 hours/week at $50/hour. What's the monthly cost?

🧠 DeepInfra LLM is working on the complex problem...
✅ Complex problem solved in 3.1s!
```

The LLM breaks down complex problems into multiple calculation steps.

## Code Structure

### API Integration
```go
// Create DeepInfra model function
deepInfraModel := conduit.CreateDeepInfraModel(bearerToken)

// Create LLM-powered agent manager
llmAgentManager := agents.NewLLMAgentManager(server, deepInfraModel, modelName)
```

### Agent Creation
```go
agent, err := llmAgentManager.CreateLLMAgent(
    "deepinfra_agent",
    "DeepInfra Problem Solver",
    "An agent powered by DeepInfra's LLM for intelligent reasoning",
    systemPrompt, // Detailed instructions for the LLM
    []string{"add", "multiply", "word_count", "remember", "recall", "uuid"},
    &agents.AgentConfig{
        MaxTokens:     1000,
        Temperature:   0.3,
        EnableMemory:  true,
        EnableLogging: true,
    },
)
```

### Task Execution
```go
task, _ := llmAgentManager.CreateTask(
    "deepinfra_agent",
    "Problem Title",
    "Problem Description",
    map[string]interface{}{
        "problem": "The problem statement",
        "param1": value1,
        "param2": value2,
    },
)

// Execute with real LLM reasoning
err := llmAgentManager.ExecuteTaskWithLLM(task.ID)
```

## Expected Output

```
🧠 AI Agents with DeepInfra LLM Integration
==========================================
Using real LLM from DeepInfra for intelligent agent reasoning

🔗 Using DeepInfra API
🤖 Model: meta-llama/Meta-Llama-3.1-8B-Instruct
✅ Created: DeepInfra Problem Solver

🧮 Demo 1: Mathematical Problem Solving
Problem: A restaurant serves 8 tables per hour. If each table pays $25 on average, how much revenue per hour?
🧠 DeepInfra LLM is analyzing the problem...
   🧮 Tool executed: 8.0 × 25.0 = 200.0
✅ Problem solved in 2.3s!
  📋 Task: Restaurant Revenue Calculation
  📊 Status: completed
    Step 1: llm_reasoning (completed)
    Step 2: calculate_revenue (completed)
      ✅ Result: 200

📝 Demo 2: Text Analysis with Real LLM
Text: DeepInfra provides access to cutting-edge language models...
🧠 DeepInfra LLM is analyzing the text...
✅ Text analyzed in 1.8s!

🔢 Demo 3: Complex Multi-Step Problem
Problem: A software team of 5 developers works 40 hours/week at $50/hour. What's the monthly cost?
🧠 DeepInfra LLM is working on the complex problem...
✅ Complex problem solved in 3.1s!

🎓 DeepInfra Integration Benefits:
=================================
✅ Real LLM reasoning and intelligence
✅ OpenAI-compatible API integration
✅ High-quality model responses
✅ Bearer token authentication
✅ Scalable cloud infrastructure
✅ Multiple model options available

📊 Session Summary:
Total Agents: 1
Total Tasks: 3
Available Tools: 19

🎉 DeepInfra demonstration completed successfully!
```

## DeepInfra vs Mock LLM vs Local Ollama

| Feature | DeepInfra | Mock LLM | Local Ollama |
|---------|-----------|----------|--------------|
| **Intelligence** | 🧠 Real AI | 🤖 Simulated | 🧠 Real AI |
| **Setup** | ☁️ Cloud API | ⚡ Instant | 🔧 Local install |
| **Performance** | 🚀 Fast cloud | ⚡ Instant | 🐌 Depends on hardware |
| **Cost** | 💰 Pay per use | 🆓 Free | 🆓 Free + hardware |
| **Reliability** | 🔒 Production | 🧪 Demo only | 🏠 Local dependency |
| **Models** | 🎯 Many options | 📝 Fixed responses | 🔄 Model dependent |

## Troubleshooting

### Common Issues

1. **Missing Token Error**
   ```
   ❌ DEEPINFRA_TOKEN environment variable is required
   ```
   **Solution**: Set your DeepInfra API token:
   ```bash
   export DEEPINFRA_TOKEN="your_token_here"
   ```

2. **Authentication Error**
   ```
   API returned status 401: Unauthorized
   ```
   **Solution**: Verify your token is correct and has sufficient credits

3. **Model Not Found**
   ```
   API returned status 404: Model not found
   ```
   **Solution**: Check available models at [DeepInfra Models](https://deepinfra.com/models)

4. **Rate Limiting**
   ```
   API returned status 429: Too many requests
   ```
   **Solution**: Add delays between requests or upgrade your plan

### Performance Tips

1. **Choose Appropriate Models**
   - 8B models: Faster, good for simple tasks
   - 70B models: Slower, better for complex reasoning

2. **Optimize Temperature**
   - 0.1-0.3: More deterministic responses
   - 0.7-0.9: More creative responses

3. **Manage Token Limits**
   - Adjust MaxTokens based on task complexity
   - Monitor usage to control costs

## Integration with Your Applications

### Basic Integration
```go
package main

import (
    "os"
    "github.com/benozo/conduit/agents"
    conduit "github.com/benozo/conduit/lib"
)

func main() {
    // Get token from environment
    token := os.Getenv("DEEPINFRA_TOKEN")
    
    // Create DeepInfra-powered agent
    server := conduit.NewEnhancedServer(conduit.DefaultConfig())
    model := conduit.CreateDeepInfraModel(token)
    manager := agents.NewLLMAgentManager(server, model, "meta-llama/Meta-Llama-3.1-8B-Instruct")
    
    // Use in your application...
}
```

### Production Considerations

1. **Error Handling**: Implement robust error handling for API failures
2. **Caching**: Cache responses for repeated queries
3. **Monitoring**: Track API usage and costs
4. **Fallbacks**: Have backup strategies for API downtime

## Next Steps

1. **Explore Models**: Try different models for various use cases
2. **Custom Prompts**: Optimize system prompts for your domain
3. **Add Tools**: Register custom tools for your specific needs
4. **Scale Up**: Implement multi-agent systems for complex workflows

This example provides a production-ready foundation for building intelligent applications with real LLM capabilities through DeepInfra's API!
