# Cloudflare Workers AI Swarm Example

This example demonstrates how to integrate **Cloudflare Workers AI** with the SwarmV2 framework to create distributed AI-powered agent swarms running on Cloudflare's global edge network.

## 🌟 Features

- **Edge Computing AI**: Leverage Cloudflare's global network for low-latency AI inference
- **Multiple AI Models**: Support for various Cloudflare Workers AI models (Llama, Mistral, Phi, etc.)
- **Collaborative Agents**: Demonstrate multi-agent workflows with specialized AI agents
- **Global Distribution**: AI processing happens close to your users worldwide
- **Cost Effective**: Pay-per-use model with competitive pricing
- **Easy Integration**: Simple API integration with robust error handling

## 🚀 Quick Start

### Prerequisites

1. **Cloudflare Account**: Sign up at [cloudflare.com](https://cloudflare.com)
2. **Workers AI Access**: Enable Workers AI in your Cloudflare dashboard
3. **API Token**: Create an API token with Workers AI permissions (Standard API)
   OR
3. **Custom Endpoint**: Use your custom Cloudflare Workers AI gateway

### Setup Instructions

#### Option 1: Custom Cloudflare Endpoint (Recommended)

For custom endpoints using your own domain:

1. **Set Environment Variables**:
   ```bash
   export CLOUDFLARE_CUSTOM_URL="https://example.com"
   export CLOUDFLARE_CUSTOM_API_KEY="XX-XX-XXXX-XXX-XXX"
   export CLOUDFLARE_MODEL="@cf/meta/llama-4-scout-17b-16e-instruct"
   ```

2. **Run the Example**:
   ```bash
   cd examples/cloudflare_ai
   go run main.go
   ```

#### Option 2: Standard Cloudflare Workers AI API

1. **Get your Cloudflare Account ID**:
   - Go to Cloudflare Dashboard → Workers & Pages → Overview
   - Copy your Account ID from the right sidebar

2. **Create an API Token**:
   - Go to Cloudflare Dashboard → My Profile → API Tokens
   - Click "Create Token" → "Custom token"
   - Add permissions: `Zone:Zone:Read`, `Zone:Zone Settings:Edit`, `Account:Cloudflare Workers:Edit`
   - Add account resources: Include your account
   - Click "Continue to summary" → "Create Token"

3. **Set Environment Variables**:
   ```bash
   export CLOUDFLARE_ACCOUNT_ID="your_account_id_here"
   export CLOUDFLARE_API_TOKEN="your_api_token_here"
   export CLOUDFLARE_MODEL="@cf/meta/llama-3.1-8b-instruct"  # Optional
   ```

4. **Run the Example**:
   ```bash
   cd examples/cloudflare_ai
   go run main.go
   ```

## 🤖 Available Models

The example supports various Cloudflare Workers AI models:

### **Meta Llama Models**
- `@cf/meta/llama-3.1-8b-instruct` - Latest Llama 3.1 8B (Recommended)
- `@cf/meta/llama-3.1-70b-instruct` - Llama 3.1 70B (Most capable)
- `@cf/meta/llama-3-8b-instruct` - Llama 3 8B
- `@cf/meta/llama-4-scout-17b-16e-instruct` - Llama 4 Scout 17B (New model)
- `@cf/meta/llama-2-7b-chat-int8` - Llama 2 7B

### **Other Popular Models**
- `@cf/mistral/mistral-7b-instruct-v0.1` - Mistral 7B
- `@cf/microsoft/phi-2` - Microsoft Phi-2
- `@cf/qwen/qwen1.5-7b-chat-awq` - Qwen 1.5 7B
- `@cf/google/gemma-7b-it` - Google Gemma 7B
- `@cf/openchat/openchat-3.5-0106` - OpenChat 3.5

## 📋 Example Workflow

The demo creates three specialized AI agents:

1. **CloudflareDataAnalyst**: Market analysis and data interpretation
2. **CloudflareContentCreator**: Marketing content and documentation
3. **CloudflareStrategicAdvisor**: Business strategy and planning

### Workflow Steps:
```
📊 Market Analysis → ✍️ Content Strategy → 🎯 Strategic Recommendations
```

## 🏗️ Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Swarm Agent   │    │  Cloudflare AI   │    │  Global Edge    │
│   (Local)       │───▶│   Provider       │───▶│   Network       │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │   AI Models      │
                    │   • Llama 3.1    │
                    │   • Mistral      │
                    │   • Phi-2        │
                    │   • Gemma        │
                    └──────────────────┘
```

## 💡 Key Benefits

### **Performance**
- **Low Latency**: AI processing on Cloudflare's edge network
- **Global Scale**: 200+ data centers worldwide
- **High Availability**: Built-in redundancy and failover

### **Cost Efficiency**
- **Pay-per-use**: No minimum commitments
- **Competitive Pricing**: Cost-effective AI inference
- **No Infrastructure**: Serverless AI without managing servers

### **Developer Experience**
- **Simple API**: RESTful HTTP API
- **Multiple Models**: Choose the right model for your use case
- **Easy Integration**: Drop-in replacement for other AI providers

## 🔧 Configuration

### Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `CLOUDFLARE_ACCOUNT_ID` | Your Cloudflare Account ID | ✅ Yes | - |
| `CLOUDFLARE_API_TOKEN` | Your Cloudflare API Token | ✅ Yes | - |
| `CLOUDFLARE_MODEL` | AI model to use | ❌ No | `@cf/meta/llama-3.1-8b-instruct` |

### Model Selection

Choose your model based on your needs:

- **Fast & Efficient**: `@cf/meta/llama-3.1-8b-instruct`
- **Most Capable**: `@cf/meta/llama-3.1-70b-instruct`
- **Specialized**: `@cf/mistral/mistral-7b-instruct-v0.1`
- **Lightweight**: `@cf/microsoft/phi-2`

## 📊 Example Output

```
=== Cloudflare Workers AI Swarm Demo ===
🌐 Using Cloudflare Workers AI
   Account ID: abcd***wxyz
   Model: @cf/meta/llama-3.1-8b-instruct

🔍 Testing connection to Cloudflare Workers AI...
✅ Successfully connected to Cloudflare Workers AI!

📋 Available Cloudflare AI Models:
   1. @cf/meta/llama-3.1-8b-instruct ← Current
   2. @cf/meta/llama-3.1-70b-instruct
   3. @cf/mistral/mistral-7b-instruct-v0.1
   ...

🤖 Cloudflare AI-Powered Agents:
   - CloudflareDataAnalyst: Model: @cf/meta/llama-3.1-8b-instruct | Provider: Cloudflare Workers AI
   - CloudflareContentCreator: Model: @cf/meta/llama-3.1-8b-instruct | Provider: Cloudflare Workers AI
   - CloudflareStrategicAdvisor: Model: @cf/meta/llama-3.1-8b-instruct | Provider: Cloudflare Workers AI

🔄 Demonstrating Cloudflare AI Swarm Workflow...
📊 Step 1: CloudflareDataAnalyst analyzing market data...
📈 Analysis Result: The AI productivity app market is experiencing rapid growth...

✍️ Step 2: CloudflareContentCreator creating content strategy...
📝 Content Strategy: Based on the market analysis, a multi-channel approach...

🎯 Step 3: CloudflareStrategicAdvisor providing strategic recommendations...
🎯 Strategic Recommendations: To successfully launch the AI productivity app...

🎉 Cloudflare Workers AI Swarm Demo completed!
```

## 🔒 Security Best Practices

1. **API Token Security**:
   - Never commit API tokens to version control
   - Use environment variables for sensitive data
   - Rotate tokens regularly

2. **Least Privilege**:
   - Create tokens with minimal required permissions
   - Scope tokens to specific accounts/zones

3. **Monitoring**:
   - Monitor API usage in Cloudflare dashboard
   - Set up alerts for unusual activity

## 🌐 Global Edge Locations

Cloudflare Workers AI runs on 200+ locations worldwide:

- **Americas**: USA, Canada, Brazil, Mexico, Argentina
- **Europe**: UK, Germany, France, Netherlands, Sweden
- **Asia-Pacific**: Japan, Singapore, Australia, India, Hong Kong
- **Africa & Middle East**: South Africa, UAE, Israel

## 💰 Pricing

Cloudflare Workers AI offers competitive pricing:

- **Free Tier**: 10,000 requests per day
- **Paid Usage**: Starting at $0.012 per 1,000 requests
- **No Minimums**: Pay only for what you use

*Prices subject to change. Check [Cloudflare pricing](https://www.cloudflare.com/plans/) for current rates.*

## 🔗 Related Examples

- **[`ollama_agent/`](../ollama_agent/)** - Local AI with Ollama
- **[`multi_agent_ollama/`](../multi_agent_ollama/)** - Multi-agent Ollama workflows
- **[`coordinator_demo/`](../coordinator_demo/)** - Agent coordination patterns
- **[`react_workflow/`](../react_workflow/)** - ReAct reasoning workflows

## 📚 Resources

- [Cloudflare Workers AI Documentation](https://developers.cloudflare.com/workers-ai/)
- [Cloudflare AI Models](https://developers.cloudflare.com/workers-ai/models/)
- [Cloudflare API Documentation](https://api.cloudflare.com/)
- [Workers AI Pricing](https://www.cloudflare.com/plans/)

## 🤝 Contributing

This example demonstrates the integration patterns. To extend:

1. Add new model providers
2. Implement streaming responses
3. Add function calling capabilities
4. Create specialized agent types

---

**Note**: This example requires a Cloudflare account with Workers AI enabled. Free tier includes 10,000 daily requests.
